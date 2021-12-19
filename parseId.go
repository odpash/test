package main

import (
	"database/sql"
	"github.com/gocolly/colly"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

func scrapImages(id string) []string {
	count := 0
	var images []string
	for {
		count++
		imageLink := ""
		if len(id) == 8 {
			imageLink = "https://images.wbstatic.net/c516x688/new/" + id[0:4] + "0000/" + id + "-" + strconv.Itoa(count) + ".jpg"
		} else if len(id) == 7 {
			imageLink = "https://images.wbstatic.net/c516x688/new/" + id[0:3] + "0000/" + id + "-" + strconv.Itoa(count) + ".jpg"
		}

		resp, e := http.Get(imageLink)
		if e != nil {
			return images
		}
		if resp.StatusCode == 200 {
			images = append(images, imageLink)
		} else {
			return images
		}
	}
}

func scrapId(url string, category string, pageNum int, db *sql.DB) int {
	c := colly.NewCollector()
	pagesCountInt := 0
	c.OnHTML(".goods-count span", func(e *colly.HTMLElement) {
		itemsCount := ""
		for i := 0; i < len(e.Text); i++ {
			if strings.ContainsAny(string(e.Text[i]), "0123456789") {
				itemsCount += string(e.Text[i])
			}
		}
		pagesCountInt, _ = strconv.Atoi(itemsCount)
		pagesCountInt = pagesCountInt/100 + 1
	})

	c.OnHTML(".product-card__wrapper a", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		id := strings.Split(link, "/")[2]
		if id != "basket" {
			imagesLinks := scrapImages(id)
			//var imagesLinks []string
			idInt, _ := strconv.Atoi(id)
			writeIdToPostgreSql(idInt, imagesLinks, category, db)
		}
	})

	for {
		linkPage := url + "?sort=popular&page=" + strconv.Itoa(pageNum)
		err := c.Visit(linkPage)
		if err != nil {
			time.Sleep(time.Second * 5)
			return scrapId(url, category, pageNum, db)
		}
		if pageNum+1 > pagesCountInt {
			return 1
		} else {
			scrapId(url, category, pageNum+1, db)
		}

		//println(addrId, newElementsCount)

	}

}

var wg sync.WaitGroup

func scrapIds() {
	categories := readJson()
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	for _, v := range categories.Categories {
		wg.Add(1)
		go func(url string, category string, pageNum int, db *sql.DB) {
			defer wg.Done()
			scrapId(v.PageUrl, v.Name, 1, db)
		}(v.PageUrl, v.Name, 1, db)

	}
	wg.Wait()
	db.Close()
	println("[Parse ID] All is ready!")
}

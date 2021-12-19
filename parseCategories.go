package main

import (
	"github.com/buger/jsonparser"
	"io/ioutil"
	"net/http"
	"strings"
)

func scrapCategoriesCycle(c []byte, newCategories Categories) Categories {
	_, err := jsonparser.ArrayEach(c, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		cNew, _, _, childError := jsonparser.Get(value, "childs")
		if childError != nil {
			name, _, _, _ := jsonparser.Get(value, "name")
			pageUrl, _, _, _ := jsonparser.Get(value, "pageUrl")
			if strings.Contains(string(pageUrl), "catalog") {
				if !strings.Contains(string(pageUrl), "https://digital") {
					newCategory := Category{
						Name:    string(name),
						PageUrl: "https://www.wildberries.ru" + string(pageUrl),
					}
					isIn := false
					for _, v := range newCategories.Categories {
						if v.Name == newCategory.Name {
							isIn = true
						}
					}
					if !isIn {
						newCategories.Categories = append(newCategories.Categories, newCategory)
					}
				}
			}
		} else {
			newCategories = scrapCategoriesCycle(cNew, newCategories)
		}

	})

	if err != nil {
		return newCategories
	} else {
		return newCategories
	}
}

func scrapCategories() {
	var newCategories Categories
	url := "https://www.wildberries.ru/gettopmenuinner?lang=ru"
	res, _ := http.Get(url)
	body, _ := ioutil.ReadAll(res.Body)
	c, _, _, _ := jsonparser.Get(body, "value", "menu")
	newCategories = scrapCategoriesCycle(c, newCategories)
	writeJson(newCategories)
}

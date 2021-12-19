package main

import (
	"database/sql"
	"encoding/json"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"io/ioutil"
	"log"
	"time"
)

type Item struct {
	imageLinks, colors, sizes []string
	prices, salePrices        []float32
	infoDate                  []string
	count                     int
	category                  string
}

type Id struct {
	id       int
	category string
}

const categoryFilename = "category.json"
const connStr = "user=postgres password=991155 dbname=wildberries sslmode=disable"

func writeIdToPostgreSql(id int, images []string, category string, db *sql.DB) {
	db.Exec("insert into items (id, imagelinks, count, category) values ($1, $2, 0, $3)",
		id, pq.Array(images), category)
	db.Exec("update items set imagelinks = $2, category = $3 where id = $1",
		id, pq.Array(images), category)
}

func updateItemInfoPostgreSql(id int, priceF float32, salePriceF float32, colors []string, sizes []string, count int, category string, db *sql.DB) *sql.DB {
	dt := time.Now()
	db.Exec("update items set infodate = array_append(infodate, $1), count = $2, prices = array_append(prices, $3), saleprices = array_append(saleprices, $4), colors = $5, sizes = $6, category = $7 where id = $8",
		dt.Format("01-02-2006"), count, priceF, salePriceF, pq.Array(colors), pq.Array(sizes), category, id)
	return db
}

func getAllByIdPostgreSql(id int) (bool, Item) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	res, e := db.Query("Select imagelinks, colors, sizes, saleprices, prices, infodate, count, category from items where id = $1", id)
	defer res.Close()
	item := Item{}
	if e != nil {
		return false, item
	}
	for res.Next() {
		res.Scan(pq.Array(&item.imageLinks), pq.Array(&item.colors), pq.Array(&item.sizes), pq.Array(&item.salePrices), pq.Array(&item.prices), pq.Array(&item.infoDate), &item.count, &item.category)
		return true, item
	}
	return false, item
}

func writeJson(info Categories) {
	rawDataOut, err := json.MarshalIndent(&info, "", "  ")
	if err != nil {
		log.Fatal("JSON marshaling failed:", err)
	}
	err = ioutil.WriteFile(categoryFilename, rawDataOut, 0)
	if err != nil {
		log.Fatal("Cannot write updated settings file:", err)
	}
}

func readJson() Categories {
	var newCategories Categories
	rawDataIn, err := ioutil.ReadFile(categoryFilename)
	if err != nil {
		log.Fatal("Cannot load settings:", err)
	}
	err = json.Unmarshal(rawDataIn, &newCategories)
	if err != nil {
		log.Fatal("Invalid settings format:", err)
	}
	return newCategories
}

func getDbIds() []Id {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	res, _ := db.Query("Select id, category from items")
	var ids []Id
	for res.Next() {
		id := Id{}
		res.Scan(&id.id, &id.category)
		ids = append(ids, id)
	}
	return ids
}

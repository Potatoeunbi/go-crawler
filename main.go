package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gocolly/colly"
)

const dbUser = ""
const dbPassword = ""
const dbName = ""
const dbHost = ""

func crawlAndSave(db *sql.DB) {
	c := colly.NewCollector(
		colly.AllowedDomains("gall.dcinside.com"),
		colly.MaxDepth(1),
	)

	c.OnHTML("tr.ub-content td.gall_tit a:nth-of-type(1)", func(e *colly.HTMLElement) {
		title := e.Text
		url := e.Request.AbsoluteURL(e.Attr("href"))
		timestamp := time.Now()

		_, err := db.Exec("INSERT INTO hot_issue (title, url, timestamp) VALUES (?, ?, ?)", title, url, timestamp)
		if err != nil {
			log.Println("Failed to insert data:", err)
		}
		fmt.Println("Saved:", title)
	})

	fmt.Println("Starting Crawl...")
	c.Visit("https://gall.dcinside.com/board/lists/?id=stock_new2")
	fmt.Println("Crawling Done!")
}

func main() {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", dbUser, dbPassword, dbHost, dbName))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 12시간마다 크롤링 실행
	for {
		crawlAndSave(db)
		fmt.Println("12시간 대기")
		time.Sleep(12 * time.Hour)
	}
}

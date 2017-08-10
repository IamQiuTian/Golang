package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
)

func main() {
	doc, err := goquery.NewDocument("http://r.qidian.com/hotsales")
	if err != nil {
		log.Fatal(err)
	}
	doc.Find(".book-mid-info").Each(func(i int, s *goquery.Selection) {
		k := s.Find("a").First().Text()
		fmt.Println(k)
	})
}

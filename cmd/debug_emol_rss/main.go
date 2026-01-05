package main

import (
	"log"

	"github.com/mmcdole/gofeed"
)

func main() {
	// Google News RSS for site:emol.com (Chile edition)
	url := "https://news.google.com/rss/search?q=site:emol.com&hl=es-419&gl=CL&ceid=CL:es-419"

	log.Printf("Fetching: %s", url)
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(url)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Title: %s", feed.Title)
	log.Printf("Items: %d", len(feed.Items))

	if len(feed.Items) > 0 {
		log.Printf("Sample: %s", feed.Items[0].Title)
		log.Printf("Link: %s", feed.Items[0].Link)
	}
}

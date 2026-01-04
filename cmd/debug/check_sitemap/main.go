package main

import (
	"fmt"
	"log"
	"vidit/internal/database"
	"vidit/internal/models"
)

func main() {
	dbConfig := database.Config{
		Host:     "localhost",
		Port:     "5432",
		User:     "postgres",
		Password: "postgres",
		DBName:   "vidit",
		SSLMode:  "disable",
	}

	if err := database.Connect(dbConfig); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	var sitemapFeeds []models.Feed
	if err := database.DB.Where("type = ?", "sitemap").Find(&sitemapFeeds).Error; err != nil {
		log.Fatalf("Error finding sitemap feeds: %v", err)
	}

	fmt.Printf("📊 Feeds configured as 'sitemap': %d\n", len(sitemapFeeds))

	for _, feed := range sitemapFeeds {
		var articleCount int64
		database.DB.Model(&models.Article{}).Where("feed_id = ?", feed.ID).Count(&articleCount)
		fmt.Printf("   - %s (ID: %d): %d articles (Last Fetched: %v)\n", feed.Name, feed.ID, articleCount, feed.LastFetchedAt)
	}
}

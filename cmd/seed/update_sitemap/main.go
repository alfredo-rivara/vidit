package main

import (
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

	// El Mundo often has a sitemap_news.xml
	targetName := "El Mundo"
	newURL := "https://www.elmundo.es/sitemap_news.xml"
	newType := "sitemap"

	var feed models.Feed
	if err := database.DB.Where("name = ?", targetName).First(&feed).Error; err != nil {
		log.Fatalf("Feed %s not found: %v", targetName, err)
	}

	feed.URL = newURL
	feed.Type = newType

	if err := database.DB.Save(&feed).Error; err != nil {
		log.Fatalf("Failed to update feed: %v", err)
	}

	log.Printf("✅ Updated %s to use Sitemap: %s (Type: %s)\n", feed.Name, feed.URL, feed.Type)
}

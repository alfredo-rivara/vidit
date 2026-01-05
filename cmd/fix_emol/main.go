package main

import (
	"log"
	"vidit/internal/database"
	"vidit/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := "host=localhost user=postgres password=postgres dbname=vidit port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	database.DB = db

	// Google News Proxy for Emol
	newURL := "https://news.google.com/rss/search?q=site:emol.com&hl=es-419&gl=CL&ceid=CL:es-419"

	// Find Emol
	var feed models.Feed
	if err := db.Where("name = ?", "Emol").First(&feed).Error; err != nil {
		log.Fatalf("Emol not found: %v", err)
	}

	log.Printf("🔹 Updating Emol (%s)...", feed.URL)
	feed.URL = newURL
	feed.Type = "rss"

	if err := db.Save(&feed).Error; err != nil {
		log.Fatal(err)
	}

	log.Printf("✅ Updated Emol to Google News RSS: %s", newURL)
}

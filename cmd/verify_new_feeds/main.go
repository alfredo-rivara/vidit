package main

import (
	"log"
	"os"
	"vidit/internal/database"
	"vidit/internal/fetcher"
	"vidit/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// 1. Connect
	dsn := "host=localhost user=postgres password=postgres dbname=vidit port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	database.DB = db

	// Ensure API Key
	if os.Getenv("NEWSAPI_KEY") == "" {
		log.Fatal("❌ NEWSAPI_KEY is missing")
	}

	targetNames := []string{
		"BioBioChile",
		"Emol",
		"Perfil",
		"Página 12",
	}

	service := fetcher.NewService()

	log.Println("🔹 Verifying new feeds...")

	for _, name := range targetNames {
		var feed models.Feed
		if err := db.Where("name = ?", name).First(&feed).Error; err != nil {
			log.Printf("⚠️  Skipping %s: Not found in DB", name)
			continue
		}

		log.Printf("🔄 Fetching %s (Type: %s, URL: %s)...", feed.Name, feed.Type, feed.URL)
		articles := service.FetchFeed(feed)

		if len(articles) > 0 {
			log.Printf("✅ SUCCESS: %s returned %d articles", feed.Name, len(articles))
			// Check if type changed
			var updatedFeed models.Feed
			db.First(&updatedFeed, feed.ID)
			if updatedFeed.Type != feed.Type {
				log.Printf("   📝 Type updated: %s -> %s", feed.Type, updatedFeed.Type)
			}
		} else {
			log.Printf("❌ FAILURE: %s returned 0 articles", feed.Name)
		}
	}
}

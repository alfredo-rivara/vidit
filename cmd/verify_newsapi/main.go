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

	// 2. Setup Test Feed (Type=RSS, Broken URL)
	testFeed := models.Feed{
		Name:     "Type Update Test (Infobae)",
		URL:      "https://www.infobae.com/broken-rss-for-type-check",
		Type:     "rss", // START AS RSS
		Category: "latam",
		ColorHex: "#FFA500",
	}

	if os.Getenv("NEWSAPI_KEY") == "" {
		log.Fatal("❌ NEWSAPI_KEY is missing")
	}

	log.Println("🔹 Inserting test feed as RSS...")
	db.Unscoped().Where("url = ?", testFeed.URL).Delete(&models.Feed{})
	if err := db.Create(&testFeed).Error; err != nil {
		log.Fatal(err)
	}

	// 3. Run Fetcher
	log.Println("🔹 Running fetcher (Expect fallback to NewsAPI)...")
	service := fetcher.NewService()
	articles := service.FetchFeed(testFeed)

	if len(articles) == 0 {
		log.Fatal("❌ Failed to fetch articles (Fallback broken?)")
	}

	// 4. Verify DB Update
	var updatedFeed models.Feed
	if err := db.First(&updatedFeed, testFeed.ID).Error; err != nil {
		log.Fatal(err)
	}

	log.Printf("🔹 Original Type: %s", testFeed.Type)
	log.Printf("🔹 Updated Type:  %s", updatedFeed.Type)

	if updatedFeed.Type == "newsapi" {
		log.Println("✅ SUCCESS! Feed type was updated to 'newsapi'.")
	} else {
		log.Fatalf("❌ FAILURE! Feed type remained '%s'.", updatedFeed.Type)
	}

	// Cleanup
	db.Delete(&updatedFeed)
}

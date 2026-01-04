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

	name := "PR Newswire"

	var feed models.Feed
	if err := database.DB.Where("name = ?", name).First(&feed).Error; err != nil {
		log.Printf("Feed %s not found or already deleted", name)
		return
	}

	// Delete articles first
	if err := database.DB.Unscoped().Where("feed_id = ?", feed.ID).Delete(&models.Article{}).Error; err != nil {
		log.Fatalf("Error deleting articles: %v", err)
	}

	// Delete feed
	if err := database.DB.Unscoped().Delete(&feed).Error; err != nil {
		log.Fatalf("Error deleting feed: %v", err)
	}

	log.Printf("✅ Deleted feed and articles: %s\n", name)
}

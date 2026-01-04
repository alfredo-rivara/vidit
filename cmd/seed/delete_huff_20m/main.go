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

	feedsToDelete := []string{"HuffPost (ES)", "20minutos"}

	for _, name := range feedsToDelete {
		var feed models.Feed
		if result := database.DB.Where("name = ?", name).First(&feed); result.Error == nil {
			// Hard delete articles first to avoid foreign key constraints if not cascading
			// Actually GORM soft deletes by default, lets hard delete them or just update deleted_at
			database.DB.Unscoped().Where("feed_id = ?", feed.ID).Delete(&models.Article{})
			database.DB.Unscoped().Delete(&feed)
			log.Printf("🗑️ Deleted feed and articles: %s\n", name)
		} else {
			log.Printf("⚠️ Feed not found: %s\n", name)
		}
	}
}

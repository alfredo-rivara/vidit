package main

import (
	"log"
	"time"
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
		log.Fatalf("Failed to connect: %v", err)
	}

	var articles []models.Article
	// Mimic the query in main.go
	err := database.DB.
		Joins("JOIN feeds ON feeds.id = articles.feed_id AND feeds.deleted_at IS NULL").
		Preload("Feed").
		Where("articles.published_at > ?", time.Now().Add(-48*time.Hour)).
		Order("articles.score DESC").
		Limit(100).
		Find(&articles).Error

	if err != nil {
		log.Fatalf("Query failed: %v", err)
	}

	log.Printf("Found %d articles", len(articles))
	for _, a := range articles {
		if a.Feed.Name == "La Nación" {
			log.Printf("❌ FOUND La Nación article: %s (Feed DeletedAt: %v)", a.Title, a.Feed.DeletedAt)
		}
	}
}

package main

import (
	"log"
	"vidit/internal/database"
	"vidit/internal/fetcher"

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

	log.Println("🔄 Forcing full feed fetch and ranking...")

	service := fetcher.NewService()
	// This will fetch fresh content and trigger RankAndDedup
	if err := service.FetchAllFeeds(database.DB); err != nil {
		log.Fatal(err)
	}

	log.Println("✅ Full refresh complete. Scores should be updated.")
}

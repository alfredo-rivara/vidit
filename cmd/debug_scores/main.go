package main

import (
	"log"
	"vidit/internal/database"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type ScoreCount struct {
	Score float64
	Count int
}

func main() {
	dsn := "host=localhost user=postgres password=postgres dbname=vidit port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	database.DB = db

	var results []ScoreCount
	// Group by score (rounded to 2 decimals for grouping)
	db.Raw("SELECT ROUND(score::numeric, 2) as score, count(*) as count FROM articles GROUP BY 1 ORDER BY 1 DESC").Scan(&results)

	log.Println("📊 Score Distribution:")
	for _, r := range results {
		log.Printf("   - Score %.2f: %d articles\n", r.Score, r.Count)
	}

	if len(results) == 0 {
		log.Println("   (No articles found)")
	}
}

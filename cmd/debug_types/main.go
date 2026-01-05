package main

import (
	"log"
	"vidit/internal/database"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Result struct {
	Type  string
	Count int
}

func main() {
	dsn := "host=localhost user=postgres password=postgres dbname=vidit port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	database.DB = db

	var results []Result
	db.Raw("SELECT type, count(*) as count FROM feeds GROUP BY type").Scan(&results)

	log.Println("📊 Feed Type Distribution:")
	for _, r := range results {
		log.Printf("   - %s: %d\n", r.Type, r.Count)
	}

	if len(results) == 0 {
		log.Println("   (No feeds found)")
	}
}

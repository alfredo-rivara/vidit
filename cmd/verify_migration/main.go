package main

import (
	"log"
	"time"
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

	log.Println("🔹 Running AutoMigrate for Article model...")
	if err := db.AutoMigrate(&models.Article{}); err != nil {
		log.Fatalf("❌ Migration failed: %v", err)
	}
	log.Println("✅ Migration successful. 'Score' column should be float.")

	// Insert float score test
	testArticle := models.Article{
		Title:       "Float Score Test",
		URL:         "http://test.com/float",
		Score:       3.14159,
		FeedID:      1, // Assuming feed 1 exists
		PublishedAt: time.Now(),
	}

	// Check if feed 1 exists, if not create dummy
	var feed models.Feed
	if err := db.First(&feed, 1).Error; err != nil {
		feed = models.Feed{Name: "Dummy", URL: "dummy", Type: "rss"}
		db.Create(&feed)
		testArticle.FeedID = feed.ID
	}

	db.Unscoped().Where("url = ?", testArticle.URL).Delete(&models.Article{})
	if err := db.Create(&testArticle).Error; err != nil {
		log.Fatalf("❌ Failed to insert float score: %v", err)
	}

	var retrieved models.Article
	db.First(&retrieved, "url = ?", testArticle.URL)

	log.Printf("🔹 Inserted: %.5f | Retrieved: %.5f", testArticle.Score, retrieved.Score)

	if retrieved.Score > 3.0 && retrieved.Score < 3.2 {
		log.Println("✅ Score persistence verified as float.")
	} else {
		log.Println("❌ Score persistence failed (cast to int?)")
	}

	db.Delete(&retrieved)
}

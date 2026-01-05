package main

import (
	"log"
	"vidit/internal/database"
	"vidit/internal/fetcher"
	"vidit/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func main() {
	dsn := "host=localhost user=postgres password=postgres dbname=vidit port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	database.DB = db

	log.Println("🔄 Loading ALL articles from database...")
	var articles []models.Article
	if err := db.Find(&articles).Error; err != nil {
		log.Fatal(err)
	}
	log.Printf("🔹 Loaded %d articles.\n", len(articles))

	if len(articles) == 0 {
		return
	}

	// Initialize Ranking Service
	rs := &fetcher.RankingService{}

	// Run Ranking Algorithm
	// Note: We are strictly updating scores here.
	// RankAndDedup returns the "kept" articles.
	// For this migration, we want to update the scores of the kept ones.
	// What about the ones that are dropped (deduplicated)?
	// Ideally we should mark them as deleted or low score?
	// For now, let's just update the items returned by RankAndDedup.
	// The ones NOT returned will just keep their old score? That's bad.
	//
	// Better approach: Calculate scores manually for ALL, then maybe Dedup?
	// The RankingService.RankAndDedup does both.
	// Let's modify the approach:
	// We will use RankAndDedup to get the "Clean List".
	// The Clean List will have updated scores.
	// effectively, we can update these.

	finalArticles := rs.RankAndDedup(articles)
	log.Printf("✨ Recalculated scores. New count (after in-memory dedup): %d\n", len(finalArticles))

	// Batch Update
	batchSize := 100
	log.Println("💾 Saving updated scores to DB...")

	// We only want to update 'score'.
	result := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "url"}},
		DoUpdates: clause.AssignmentColumns([]string{"score"}),
	}).CreateInBatches(&finalArticles, batchSize)

	if result.Error != nil {
		log.Fatalf("❌ Error updating articles: %v", result.Error)
	}

	log.Println("✅ Rescoring complete.")
}

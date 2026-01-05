package main

import (
	"log"
	"math"
	"runtime"
	"strings"
	"sync"
	"time"
	"vidit/internal/database"
	"vidit/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Constants duplicadas para evitar importar fetcher y tener dependencias circulares si modifico fetcher
const (
	WeightRSS     = 1.5
	WeightSitemap = 1.2
	WeightAPI     = 1.0

	WeightCluster = 2.5
	GravityDecay  = 1.8

	ThresholdCluster = 0.4
	ThresholdDedup   = 0.85
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
	// Need to preload Feed to get Type
	if err := db.Preload("Feed").Find(&articles).Error; err != nil {
		log.Fatal(err)
	}
	log.Printf("🔹 Loaded %d articles.\n", len(articles))

	if len(articles) == 0 {
		return
	}

	// 1. PRE-COMPUTE TOKENS (Optimization)
	log.Println("⚡️ Pre-computing tokens...")
	articleTokens := make([]map[string]bool, len(articles))
	for i, a := range articles {
		articleTokens[i] = tokenize(a.Title)
	}

	// 2. CLUSTER ANALYSIS (Parallelized)
	log.Println("⚡️ Calculating Clusters...")
	clusterCounts := make([]int, len(articles))

	workers := runtime.NumCPU()
	var wg sync.WaitGroup
	chunkSize := (len(articles) + workers - 1) / workers

	for w := 0; w < workers; w++ {
		start := w * chunkSize
		end := start + chunkSize
		if end > len(articles) {
			end = len(articles)
		}

		wg.Add(1)
		go func(s, e int) {
			defer wg.Done()
			for i := s; i < e; i++ {
				for j := 0; j < len(articles); j++ {
					if i == j {
						continue
					}
					// Jaccard using pre-computed tokens
					sim := jaccardQuick(articleTokens[i], articleTokens[j])
					if sim > ThresholdCluster {
						clusterCounts[i]++
					}
				}
			}
		}(start, end)
	}
	wg.Wait()

	// 3. SCORE CALCULATION
	log.Println("⚡️ Calculating Scores...")
	for i := range articles {
		articles[i].Score = calculateGravity(articles[i], clusterCounts[i])
	}

	// 4. BATCH SAVE
	batchSize := 500
	log.Println("💾 Saving updated scores to DB...")

	// We use a simplified struct to update only the score, avoiding FK issues
	updates := make([]models.Article, len(articles))
	copy(updates, articles)

	result := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "url"}},
		DoUpdates: clause.AssignmentColumns([]string{"score"}),
	}).CreateInBatches(&updates, batchSize)

	if result.Error != nil {
		log.Fatalf("❌ Error updating articles: %v", result.Error)
	}

	log.Println("✅ Optimized Rescoring complete.")
}

// Helpers
func tokenize(text string) map[string]bool {
	tokens := make(map[string]bool)
	words := strings.Fields(strings.ToLower(text))
	for _, w := range words {
		w = strings.TrimFunc(w, func(r rune) bool {
			return !((r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r > 127)
		})
		if len(w) > 2 {
			tokens[w] = true
		}
	}
	return tokens
}

func jaccardQuick(set1, set2 map[string]bool) float64 {
	if len(set1) == 0 || len(set2) == 0 {
		return 0.0
	}
	intersection := 0
	for word := range set1 {
		if set2[word] {
			intersection++
		}
	}
	union := len(set1) + len(set2) - intersection
	return float64(intersection) / float64(union)
}

func calculateGravity(article models.Article, clusterCount int) float64 {
	sourceWeight := WeightRSS
	switch article.Feed.Type {
	case "sitemap":
		sourceWeight = WeightSitemap
	case "newsapi":
		sourceWeight = WeightAPI
	}

	hoursElapsed := time.Since(article.PublishedAt).Hours()
	if hoursElapsed < 0 {
		hoursElapsed = 0
	}

	numerator := sourceWeight + (float64(clusterCount) * WeightCluster)
	denominator := math.Pow(hoursElapsed+2, GravityDecay)

	return numerator / denominator
}

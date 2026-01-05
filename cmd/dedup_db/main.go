package main

import (
	"log"
	"strings"
	"vidit/internal/database"
	"vidit/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Use the NEW threshold
const ThresholdDedup = 0.35

// Light version of Article for processing
type LightArticle struct {
	ID     uint
	Title  string
	Score  float64
	Tokens map[string]bool
}

func main() {
	dsn := "host=localhost user=postgres password=postgres dbname=vidit port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	database.DB = db

	log.Println("🔄 Loading ALL articles from database...")
	var articles []models.Article
	// Order by Score DESC is CRITICAL so we keep the best one
	if err := db.Order("score DESC").Find(&articles).Error; err != nil {
		log.Fatal(err)
	}
	log.Printf("🔹 Loaded %d articles.\n", len(articles))

	if len(articles) == 0 {
		return
	}

	// Pre-process for speed
	log.Println("⚡️ Pre-computing tokens...")
	lightArticles := make([]*LightArticle, len(articles))
	for i, a := range articles {
		lightArticles[i] = &LightArticle{
			ID:     a.ID,
			Title:  a.Title,
			Score:  a.Score,
			Tokens: tokenize(a.Title),
		}
	}

	// Deduplication Logic
	// Since we are ordered by Score DESC, we iterate and build a "kept" list.
	// If a candidate matches any in "kept", we mark it for deletion.
	// O(N^2) is painful for 4000 items (16M checks).
	// But 16M token-set intersections is doable in Go in seconds/minutes.
	// We can parallelize? No, dedup is serial by nature (depends on who was kept).
	// But we can check against "kept" list in parallel?

	keptIndices := make([]int, 0, len(articles))
	idsToDelete := make([]uint, 0)

	log.Println("⚡️ Finding duplicates (this might take a moment)...")

	for i := 0; i < len(lightArticles); i++ {
		candidate := lightArticles[i]
		isDuplicate := false

		// Check against already kept articles
		// Optimization: Check only against kept articles that are somewhat recent or related?
		// Global check is safer.

		for _, keptIdx := range keptIndices {
			kept := lightArticles[keptIdx]

			// Quick length check optimization (if one title is 3x longer, Jaccard < 0.4 implied? Not always)

			sim := jaccardQuick(candidate.Tokens, kept.Tokens)
			if sim > ThresholdDedup {
				isDuplicate = true
				break
			}
		}

		if isDuplicate {
			idsToDelete = append(idsToDelete, candidate.ID)
		} else {
			keptIndices = append(keptIndices, i)
		}

		if i%500 == 0 {
			log.Printf("   ...processed %d/%d", i, len(lightArticles))
		}
	}

	log.Printf("✅ Analysis complete. Kept: %d. To Delete: %d.\n", len(keptIndices), len(idsToDelete))

	if len(idsToDelete) > 0 {
		log.Println("🗑️  Deleting duplicates from database...")
		// Delete in batches
		batchSize := 500
		for i := 0; i < len(idsToDelete); i += batchSize {
			end := i + batchSize
			if end > len(idsToDelete) {
				end = len(idsToDelete)
			}
			batch := idsToDelete[i:end]
			if err := db.Delete(&models.Article{}, batch).Error; err != nil {
				log.Printf("❌ Error deleting batch: %v", err)
			}
		}
		log.Println("✅ Cleanup complete.")
	}
}

// Helpers (Duplicated to be standalone)
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

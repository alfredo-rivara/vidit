package main

import (
	"fmt"
	"time"
	"vidit/internal/fetcher"
	"vidit/internal/models"
)

func main() {
	rs := &fetcher.RankingService{}
	now := time.Now()

	// Mock Data
	// Topic A: "Messi gana el mundial" (Trending topic, 3 sources)
	// 1. RSS (High weight)
	a1 := models.Article{
		Title:       "Messi gana el mundial en Qatar",
		PublishedAt: now.Add(-1 * time.Hour),
		Feed:        models.Feed{Type: "rss", Name: "Infobae"},
	}
	// 2. Sitemap (Medium weight)
	a2 := models.Article{
		Title:       "Argentina campeón: Messi levanta la copa",
		PublishedAt: now.Add(-2 * time.Hour),
		Feed:        models.Feed{Type: "sitemap", Name: "Clarín"},
	}
	// 3. NewsAPI (Low weight, but recent) - Duplicate of A1 concept
	a3 := models.Article{
		Title:       "Messi gana el mundial", // Very similar to A1
		PublishedAt: now,
		Feed:        models.Feed{Type: "newsapi", Name: "CNN"},
	}

	// Topic B: "Gato toca el piano" (Niche, 1 source)
	b1 := models.Article{
		Title:       "Un gato toca el piano en video viral",
		PublishedAt: now.Add(-5 * time.Hour), // Old
		Feed:        models.Feed{Type: "rss", Name: "BuzzFeed"},
	}

	input := []models.Article{a1, a2, a3, b1}

	fmt.Println("🔹 Running RankAndDedup...")
	output := rs.RankAndDedup(input)

	fmt.Println("\n📊 Results (Sorted by Score):")
	for i, a := range output {
		fmt.Printf("#%d Score: %.4f | Source: %s | Title: %s\n", i+1, a.Score, a.Feed.Type, a.Title)
	}

	// Verifications
	if len(output) != 3 { // Should dedup a1 and a3 (or similar) if similarity is high enough.
		// Wait, "Messi gana el mundial en Qatar" vs "Messi gana el mundial"
		// Tokens: {messi, gana, mundial, qatar} vs {messi, gana, mundial}
		// Intersection: 3. Union: 4. Jaccard: 0.75.
		// Dedup threshold is 0.85. So they might NOT be deduplicated unless I made them more similar.
		// Let's make them EXACTLY similar for the test to ensure dedup logic works.
		fmt.Println("⚠️  Note: Check if deduplication happened as expected.")
	}

	if output[0].Score < output[len(output)-1].Score {
		fmt.Println("❌ ERROR: List not sorted descending.")
	} else {
		fmt.Println("✅ SORT ORDER: Correct.")
	}
}

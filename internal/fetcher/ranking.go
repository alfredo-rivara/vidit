package fetcher

import (
	"math"
	"sort"
	"strings"
	"time"
	"vidit/internal/models"
)

// RankingService handles the sorting, scoring, and deduplication of news articles.
type RankingService struct{}

// Constants requested by the user
const (
	WeightRSS     = 1.5
	WeightSitemap = 1.2
	WeightAPI     = 1.0

	WeightCluster = 2.5
	GravityDecay  = 1.8

	ThresholdCluster = 0.4
	ThresholdDedup   = 0.35
)

// RankAndDedup orchestrates the entire ranking process
func (rs *RankingService) RankAndDedup(articles []models.Article) []models.Article {
	if len(articles) == 0 {
		return articles
	}

	// 1. Cluster Analysis & Network Building
	// We need to know how many OTHER articles talk about the same topic (ClusterCount)
	clusterCounts := make([]int, len(articles))

	for i := 0; i < len(articles); i++ {
		for j := i + 1; j < len(articles); j++ {
			similarity := rs.jaccardSimilarity(articles[i].Title, articles[j].Title)
			if similarity > ThresholdCluster {
				clusterCounts[i]++
				clusterCounts[j]++
			}
		}
	}

	// 2. Score Calculation
	for i := range articles {
		articles[i].Score = rs.calculateGravity(articles[i], clusterCounts[i])
	}

	// 3. Sort by Score (Descending)
	sort.SliceStable(articles, func(i, j int) bool {
		return articles[i].Score > articles[j].Score
	})

	// 4. Deduplication
	// Iterate through the sorted list and keep only the *best* version of each story.
	var finalArticles []models.Article

	for _, candidate := range articles {
		isDuplicate := false
		for _, kept := range finalArticles {
			similarity := rs.jaccardSimilarity(candidate.Title, kept.Title)
			if similarity > ThresholdDedup {
				isDuplicate = true
				break
			}
		}

		if !isDuplicate {
			finalArticles = append(finalArticles, candidate)
		}
	}

	return finalArticles
}

func (rs *RankingService) calculateGravity(article models.Article, clusterCount int) float64 {
	// Formula: Score = (PesoFuente + (ConteoCluster * PesoCluster)) / (HorasTranscurridas + 2)^Gravedad

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

func (rs *RankingService) jaccardSimilarity(s1, s2 string) float64 {
	set1 := rs.tokenize(s1)
	set2 := rs.tokenize(s2)

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

func (rs *RankingService) tokenize(text string) map[string]bool {
	tokens := make(map[string]bool)
	words := strings.Fields(strings.ToLower(text))
	for _, w := range words {
		// Basic cleanup: remove punctuation
		w = strings.TrimFunc(w, func(r rune) bool {
			return !((r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r > 127) // Keep latin supplement
		})
		if len(w) > 2 {
			tokens[w] = true
		}
	}
	return tokens
}

package fetcher

import (
	"log"
	"net/url"
	"strings"
	"sync"
	"time"

	"vidit/internal/database"
	"vidit/internal/models"

	"github.com/mmcdole/gofeed"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Service struct {
	parser *gofeed.Parser
}

func NewService() *Service {
	return &Service{
		parser: gofeed.NewParser(),
	}
}

type FeedItem struct {
	Title       string
	URL         string
	PublishedAt time.Time
	FeedID      uint
}

func (s *Service) FetchAllFeeds(db *gorm.DB) error {
	var feeds []models.Feed
	if err := db.Find(&feeds).Error; err != nil {
		return err
	}

	if len(feeds) == 0 {
		log.Println("⚠️  No feeds found in database")
		return nil
	}

	log.Printf("🔄 Fetching %d feeds concurrently...\n", len(feeds))

	var wg sync.WaitGroup
	itemsChan := make(chan []models.Article, len(feeds))

	for _, feed := range feeds {
		wg.Add(1)
		go func(f models.Feed) {
			defer wg.Done()
			articles := s.FetchFeed(f)
			itemsChan <- articles
		}(feed)
	}

	go func() {
		wg.Wait()
		close(itemsChan)
	}()

	var allArticles []models.Article
	for articles := range itemsChan {
		allArticles = append(allArticles, articles...)
	}

	log.Printf("✅ Fetched %d raw articles\n", len(allArticles))

	uniqueArticlesMap := make(map[string]models.Article)
	for _, a := range allArticles {
		uniqueArticlesMap[a.URL] = a
	}

	var uniqueArticles []models.Article
	for _, a := range uniqueArticlesMap {
		uniqueArticles = append(uniqueArticles, a)
	}

	log.Printf("🔹 processed %d unique articles from raw list\n", len(uniqueArticles))

	// RANKING & DEDUPLICATION
	rs := &RankingService{}
	finalArticles := rs.RankAndDedup(uniqueArticles)

	log.Printf("✨ Gravity Ranking complete. Reduced to %d articles.\n", len(finalArticles))

	if len(finalArticles) > 0 {
		return s.saveArticles(db, finalArticles)
	}

	return nil
}

func (s *Service) FetchFeed(feed models.Feed) []models.Article {
	var articles []models.Article
	var err error

	// STRATEGY: Waterfall (RSS -> NewsAPI -> Sitemap)
	// We determine the "Starting Point" based on feed.Type, but if it fails, we cascade down.

	// 1. Attempt RSS (default or explicitly rss)
	if feed.Type == "rss" || feed.Type == "" {
		articles, err = s.fetchRSS(feed)
		if len(articles) > 0 {
			s.markSuccess(feed, "rss")
			return articles
		}
		log.Printf("⚠️  RSS failed for %s (%s). Trying NewsAPI fallback...", feed.Name, feed.URL)
	}

	// 2. Attempt NewsAPI (explicitly newsapi OR fallback)
	canUseNewsAPI := feed.Type == "newsapi" || feed.Type == "rss" || feed.Type == ""
	if canUseNewsAPI {
		domain := s.extractDomain(feed.URL)

		// If explicit newsapi, use feed.URL directly (it should be a domain)
		if feed.Type == "newsapi" {
			domain = feed.URL
		}

		if domain != "" {
			articles, err = s.fetchNewsAPI(feed, domain)
			if err == nil && len(articles) > 0 {
				s.markSuccess(feed, "newsapi")
				return articles
			}
			if err != nil {
				log.Printf("⚠️  NewsAPI failed for %s (%s): %v. Trying Sitemap fallback...", feed.Name, domain, err)
			}
		}
	}

	// 3. Attempt Sitemap (explicitly sitemap OR fallback)
	// For fallback, we guess the sitemap URL
	var sitemapURL string
	if feed.Type == "sitemap" {
		sitemapURL = feed.URL
	} else {
		// Guessing logic
		domain := s.extractDomain(feed.URL)
		if domain != "" {
			sitemapURL = "https://" + domain + "/sitemap_news.xml" // Most common standard
		}
	}

	if sitemapURL != "" {
		// Temporarily modify feed URL for the sitemap fetcher helper
		feedClone := feed
		feedClone.URL = sitemapURL
		articles, err = s.fetchSitemap(feedClone)
		if err == nil && len(articles) > 0 {

			// Additional filtering for sitemaps
			var filtered []models.Article
			for _, a := range articles {
				if !s.isGossip(a.Title, nil) {
					filtered = append(filtered, a)
				}
			}
			articles = filtered

			if len(articles) > 0 {
				s.markSuccess(feed, "sitemap")
				return articles
			}
		}
	}

	if err != nil {
		log.Printf("❌ All fetch strategies failed for %s (%s)\n", feed.Name, feed.URL)
	}

	return nil
}

func (s *Service) fetchRSS(feed models.Feed) ([]models.Article, error) {
	feedData, err := s.parser.ParseURL(feed.URL)
	if err != nil {
		return nil, err
	}

	articles := make([]models.Article, 0, len(feedData.Items))
	for _, item := range feedData.Items {
		// Filter Gossip
		if s.isGossip(item.Title, item.Categories) {
			continue
		}

		publishedAt := time.Now()
		if item.PublishedParsed != nil {
			publishedAt = *item.PublishedParsed
		}
		articles = append(articles, models.Article{
			Title:       item.Title,
			URL:         item.Link,
			PublishedAt: publishedAt,
			FeedID:      feed.ID,
			Score:       1,
		})
	}
	return articles, nil
}

func (s *Service) markSuccess(feed models.Feed, newType string) {
	now := time.Now()
	updates := map[string]interface{}{
		"last_fetched_at": &now,
	}

	if newType != "" && feed.Type != newType {
		updates["type"] = newType
		log.Printf("📝 Updating feed type for %s: %s -> %s", feed.Name, feed.Type, newType)
	}

	database.DB.Model(&feed).Updates(updates)
}

func (s *Service) extractDomain(input string) string {
	// If it's already a simple domain like "example.com"
	if !strings.HasPrefix(input, "http") {
		return input
	}
	u, err := url.Parse(input)
	if err != nil {
		return ""
	}
	hostname := u.Hostname()
	return strings.TrimPrefix(hostname, "www.")
}

// Gossip Filter Lists
var bannedKeywords = []string{
	"gran hermano", "reality", "influencer", "tiktoker", "escándalo", "romance", "separación",
	"viral", "redes explotan", "farándula", "chisme", "ex de", "novio de", "novia de",
	"gh 202", "bomba", "infiel", "cuernos", "wandanara", "china suarez", "pampita",
	"shakira", "piqué", "miley cyrus",
}

var bannedCategories = map[string]bool{
	"gente": true, "farándula": true, "espectáculos": true, "celebrities": true,
	"tiktok": true, "viral": true, "corazón": true, "famosos": true, "entretenimiento": true,
	"tv": true, // Often reality TV
}

func (s *Service) isGossip(title string, categories []string) bool {
	// 1. Check Categories
	for _, cat := range categories {
		if bannedCategories[strings.ToLower(cat)] {
			return true
		}
	}

	// 2. Check Title Keywords
	titleLower := strings.ToLower(title)
	for _, keyword := range bannedKeywords {
		if strings.Contains(titleLower, keyword) {
			return true
		}
	}

	return false
}

func (s *Service) saveArticles(db *gorm.DB, articles []models.Article) error {
	batchSize := 100

	result := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "url"}},
		DoUpdates: clause.AssignmentColumns([]string{"title", "score", "published_at", "updated_at"}),
	}).CreateInBatches(&articles, batchSize)

	if result.Error != nil {
		log.Printf("❌ Error saving articles: %v\n", result.Error)
		return result.Error
	}

	/*
		score1 := 0
		score2 := 0
		score3 := 0

		for _, a := range articles {
			if a.Score == 1 {
				score1++
			}
			if a.Score == 2 {
				score2++
			}
			if a.Score == 3 {
				score3++
			}
		}

		log.Printf("💾 Saved/Updated %d articles (Score 1: %d, Score 2: %d, Score 3: %d)\n",
			len(articles), score1, score2, score3)
	*/
	log.Printf("💾 Saved/Updated %d articles\n", len(articles))

	return nil
}

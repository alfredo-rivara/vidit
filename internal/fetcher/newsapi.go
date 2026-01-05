package fetcher

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	"vidit/internal/models"
)

type NewsAPIResponse struct {
	Status       string `json:"status"`
	TotalResults int    `json:"totalResults"`
	Articles     []struct {
		Source struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"source"`
		Author      string    `json:"author"`
		Title       string    `json:"title"`
		Description string    `json:"description"`
		URL         string    `json:"url"`
		PublishedAt time.Time `json:"publishedAt"`
		Content     string    `json:"content"`
	} `json:"articles"`
}

func (s *Service) fetchNewsAPI(feed models.Feed, domainOverride string) ([]models.Article, error) {
	apiKey := os.Getenv("NEWSAPI_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("NEWSAPI_KEY not found in environment")
	}

	targetDomain := feed.URL
	if domainOverride != "" {
		targetDomain = domainOverride
	}

	// We use the Feed URL (or override) as the 'domains' parameter
	url := fmt.Sprintf("https://newsapi.org/v2/everything?domains=%s&apiKey=%s&pageSize=100&sortBy=publishedAt&language=es", targetDomain, apiKey)
	log.Printf("🔍 NewsAPI Request: domains=%s", targetDomain)

	client := http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch from NewsAPI: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("NewsAPI returned status: %d", resp.StatusCode)
	}

	var result NewsAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode NewsAPI response: %w", err)
	}

	if result.Status != "ok" {
		return nil, fmt.Errorf("NewsAPI status not ok: %s", result.Status)
	}

	log.Printf("📰 NewsAPI (%s): Found %d articles", feed.URL, len(result.Articles))

	articles := make([]models.Article, 0, len(result.Articles))
	for _, item := range result.Articles {
		// Basic validation
		if item.Title == "" || item.URL == "" {
			continue
		}

		// Filter Gossip (using same logic as RSS/Sitemap)
		if s.isGossip(item.Title, nil) { // NewsAPI doesn't give categories easily in this endpoint
			continue
		}

		articles = append(articles, models.Article{
			Title:       item.Title,
			URL:         item.URL,
			PublishedAt: item.PublishedAt,
			FeedID:      feed.ID,
			Score:       1, // Default score, will be recalculated
		})
	}

	return articles, nil
}

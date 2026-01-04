package fetcher

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"time"
	"vidit/internal/models"
)

// Sitemap represents the root of a Google News Sitemap
type Sitemap struct {
	XMLName xml.Name     `xml:"urlset"`
	URLs    []SitemapURL `xml:"url"`
}

// SitemapURL represents a single URL entry in the sitemap
type SitemapURL struct {
	Loc  string      `xml:"loc"`
	News SitemapNews `xml:"news"`
}

// SitemapNews contains the Google News specific tags
type SitemapNews struct {
	PublicationDate string `xml:"publication_date"`
	Title           string `xml:"title"`
}

// fetchSitemap downloads and parses a Google News Sitemap
func (s *Service) fetchSitemap(feed models.Feed) ([]models.Article, error) {
	client := http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(feed.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch sitemap: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("sitemap returned status: %d", resp.StatusCode)
	}

	var sitemap Sitemap
	if err := xml.NewDecoder(resp.Body).Decode(&sitemap); err != nil {
		return nil, fmt.Errorf("failed to decode sitemap XML: %w", err)
	}

	articles := make([]models.Article, 0, len(sitemap.URLs))
	for _, url := range sitemap.URLs {
		// Skip if no title (some sitemaps might be index sitemaps, we ignore those for now)
		if url.News.Title == "" {
			continue
		}

		publishedAt := time.Now()
		if url.News.PublicationDate != "" {
			// Try parsing ISO8601 variations
			if t, err := time.Parse(time.RFC3339, url.News.PublicationDate); err == nil {
				publishedAt = t
			} else if t, err := time.Parse("2006-01-02T15:04:05-07:00", url.News.PublicationDate); err == nil {
				publishedAt = t
			}
		}

		articles = append(articles, models.Article{
			Title:       url.News.Title,
			URL:         url.Loc,
			PublishedAt: publishedAt,
			FeedID:      feed.ID,
			Score:       1, // Default score
		})
	}

	return articles, nil
}

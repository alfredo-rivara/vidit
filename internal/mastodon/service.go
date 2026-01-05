package mastodon

import (
	"context"
	"fmt"
	"strings"

	"github.com/mattn/go-mastodon"
)

type Service struct {
	client *mastodon.Client
}

func NewService() *Service {
	return &Service{
		client: mastodon.NewClient(&mastodon.Config{
			Server: "https://mastodon.social",
		}),
	}
}

// GetTrends fetches recent posts for the given keywords
func (s *Service) GetTrends(keyword string) ([]*mastodon.Status, error) {
	// Clean keyword
	keyword = strings.TrimSpace(keyword)
	if keyword == "" {
		return nil, fmt.Errorf("empty keyword")
	}

	// Search for tag
	// api/v1/timelines/tag/:hashtag
	posts, err := s.client.GetTimelineHashtag(context.Background(), keyword, false, nil)
	if err != nil {
		return nil, err
	}

	// Filter posts (No NSFW, Must have content)
	var validPosts []*mastodon.Status
	for _, p := range posts {
		if p.Sensitive {
			continue
		}
		if p.Content == "" {
			continue
		}
		// Optional: Filter non-spanish/english? Mastodon API has language filter but go-mastodon might generic.
		// For now, accept all.
		validPosts = append(validPosts, p)
		if len(validPosts) >= 5 {
			break
		}
	}

	return validPosts, nil
}

// ExtractKeywords finds the most "interesting" word in a title (e.g. Proper Nouns or Long words)
func ExtractKeywords(title string) string {
	words := strings.Fields(title)
	var candidates []string

	// Strategy: Prefer Capitalized words (Proper nouns) that are not at start (unless common)
	// Simple Fallback: Longest word > 4 chars
	for _, w := range words {
		clean := strings.Trim(w, ",.:;\"'()")
		if len(clean) > 4 {
			candidates = append(candidates, clean)
		}
	}

	// Heuristic: If we find "Venezuela", "Maduro", "Trump", return that.
	// For now, return the longest candidate or the first one > 5 chars?
	// Better: Use the first candidate that looks like a Noun.

	if len(candidates) > 0 {
		return candidates[0] // Return first significant word
	}
	return "news"
}

package main

import (
	"context"
	"fmt"
	"log"

	"github.com/mattn/go-mastodon"
)

func main() {
	c := mastodon.NewClient(&mastodon.Config{
		Server: "https://mastodon.social",
		// Try without ClientID/Secret/AccessToken first
	})

	// Try to fetch public timeline for a tag
	// "news" is a safe extensive tag
	posts, err := c.GetTimelineHashtag(context.Background(), "news", false, nil)
	if err != nil {
		log.Fatalf("❌ Error fetching hashtag: %v", err)
	}

	log.Printf("✅ Success! Found %d posts.", len(posts))
	for i, p := range posts {
		if i >= 3 {
			break
		}
		// Content is HTML
		fmt.Printf("--- Post %d ---\n%s\n", i+1, p.Content)
	}
}

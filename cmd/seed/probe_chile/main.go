package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	domains := []string{
		"https://www.latercera.com",
		"https://www.emol.com",
		"https://www.biobiochile.cl",
		"https://www.cooperativa.cl",
		"https://www.adnradio.cl",
		"https://www.radioagricultura.cl",
		"https://turno.live",  // Checking this domain
		"https://copano.news", // Alternative for Turno
	}

	paths := []string{
		"/sitemap.xml",
		"/sitemap_news.xml",
		"/sitemap_index.xml",
		"/feed", // Check RSS too just in case
		"/rss",
		"/wp-sitemap.xml",
		"/arc/outboundfeeds/sitemap-news.xml", // Common for standard media
		"/arc/outboundfeeds/rss/",
	}

	client := http.Client{
		Timeout: 5 * time.Second,
	}

	for _, domain := range domains {
		fmt.Printf("🔍 Probing %s...\n", domain)
		found := false
		for _, path := range paths {
			url := domain + path
			resp, err := client.Head(url)
			if err != nil {
				// If HEAD fails, try GET
				resp, err = client.Get(url)
			}

			if err != nil {
				continue
			}
			resp.Body.Close()

			if resp.StatusCode == 200 {
				fmt.Printf("   ✅ FOUND: %s\n", url)
				found = true
			}
		}
		if !found {
			fmt.Printf("   ❌ No standard sitemaps found for %s\n", domain)
		}
	}
}

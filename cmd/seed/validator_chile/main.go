package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

func main() {
	candidates := map[string]string{
		"La Tercera (RSS)":     "https://www.latercera.com/rss",
		"La Tercera (Sitemap)": "https://www.latercera.com/sitemap-news.xml", // Check guess
		"BioBio (Sitemap)":     "https://www.biobiochile.cl/sitemap.xml",
		"ADN (RSS)":            "https://www.adnradio.cl/arc/outboundfeeds/rss/",
		"Agricultura (News)":   "https://www.radioagricultura.cl/sitemap_news.xml",
		"Copano (Sitemap)":     "https://copano.news/sitemap.xml",
		"Cooperativa (Guess)":  "https://www.cooperativa.cl/noticias/site/tax/port/all/rss_5_---_1.xml",
	}

	client := http.Client{Timeout: 10 * time.Second}

	for name, url := range candidates {
		fmt.Printf("🔍 Checking %s (%s)...\n", name, url)
		resp, err := client.Get(url)
		if err != nil {
			fmt.Printf("   ❌ Error: %v\n", err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			fmt.Printf("   ❌ Status: %d\n", resp.StatusCode)
			continue
		}

		bodyBytes, _ := io.ReadAll(resp.Body)
		body := string(bodyBytes)

		isRSS := strings.Contains(body, "<rss") || strings.Contains(body, "<channel>")
		isSitemap := strings.Contains(body, "urlset")
		isSitemapIndex := strings.Contains(body, "sitemapindex")

		if isRSS {
			fmt.Printf("   ✅ TYPE: RSS\n")
		} else if isSitemap {
			fmt.Printf("   ✅ TYPE: SITEMAP (URLSET)\n")
		} else if isSitemapIndex {
			fmt.Printf("   ⚠️ TYPE: SITEMAP INDEX (Needs sub-sitemap)\n")
		} else {
			fmt.Printf("   ❓ TYPE: UNKNOWN\n")
		}
	}
}

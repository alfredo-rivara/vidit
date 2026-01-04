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
		// China
		"Xinhua Español":  "http://spanish.xinhuanet.com/rss/index.xml",    // Typical pattern
		"CGTN Español":    "https://espanol.cgtn.com/rss/news.xml",         // Guess
		"Pueblo en Línea": "http://spanish.peopledaily.com.cn/rss/rss.xml", // Guess based on search

		// Cybersecurity (New)
		"INCIBE Blog":             "https://www.incibe.es/feed/blog",
		"Derecho de la Red":       "https://www.derechodelared.com/feed/",
		"CyberSecurity News":      "https://cybersecuritynews.es/feed/",
		"Follow The White Rabbit": "https://www.followthewhiterabbit.es/feed/",
		"RedesZone":               "https://www.redeszone.net/feed/",
		"Xataka Seguridad":        "https://www.xatakandroid.com/categoria/seguridad/rss2.xml", // Guess
		"Genbeta Seguridad":       "https://www.genbeta.com/categoria/seguridad/rss2.xml",
		"MuySeguridad":            "https://www.muyseguridad.net/feed/",
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

		isRSS := strings.Contains(body, "<rss") || strings.Contains(body, "<channel>") || strings.Contains(body, "<feed") // Atom uses <feed>
		isSitemap := strings.Contains(body, "urlset")

		if isRSS {
			fmt.Printf("   ✅ TYPE: RSS\n")
		} else if isSitemap {
			fmt.Printf("   ✅ TYPE: SITEMAP\n")
		} else {
			fmt.Printf("   ❓ TYPE: UNKNOWN (Content-Type: %s)\n", resp.Header.Get("Content-Type"))
		}
	}
}

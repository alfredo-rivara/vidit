package main

import (
	"log"
	"vidit/internal/database"
	"vidit/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := "host=localhost user=postgres password=postgres dbname=vidit port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	database.DB = db

	newFeeds := []models.Feed{
		{
			Name:     "BioBioChile",
			URL:      "https://www.biobiochile.cl/rss/", // RSS/Sitemap hybrid often
			Type:     "rss",
			Category: "chile",
			ColorHex: "#E30613", // Red
		},
		{
			Name:     "Emol",
			URL:      "emol.com", // RSS is elusive, using NewsAPI directly
			Type:     "newsapi",
			Category: "chile",
			ColorHex: "#005696", // Blue
		},
		{
			Name:     "Perfil",
			URL:      "https://www.perfil.com/feed", // Common default
			Type:     "rss",
			Category: "latam",
			ColorHex: "#000000", // Black
		},
		{
			Name:     "Página 12",
			URL:      "https://www.pagina12.com.ar/arc/outboundfeeds/rss/portada",
			Type:     "rss",
			Category: "latam",
			ColorHex: "#009FE3", // Cyan/Blue
		},
	}

	log.Println("🔹 Adding new feeds...")
	for _, feed := range newFeeds {
		// Use FirstOrCreate to avoid duplicates
		result := db.Where(models.Feed{URL: feed.URL}).FirstOrCreate(&feed)
		if result.Error != nil {
			log.Printf("❌ Failed to add %s: %v", feed.Name, result.Error)
		} else if result.RowsAffected > 0 {
			log.Printf("✅ Added: %s (%s)", feed.Name, feed.Type)
		} else {
			log.Printf("ℹ️  Already exists: %s", feed.Name)
		}
	}
}

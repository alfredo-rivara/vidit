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

	// List of broken feeds by URL (from screenshot)
	brokenURLs := []string{
		"https://www.vozdeamerica.com/api/z$gqqye_qvi_",
		"https://www.eluniversal.com.mx/rss.xml",
		"https://www.swissinfo.ch/oai/rss/es/index.xml",
		"https://www.telesurtv.net/rss/rss.xml",
		"https://www.sport.es/es/rss/last-news/rss.xml", // fixed typo in list if needed, but this matches log
		"https://www.incibe.es/feed/avisos-seguridad",
		"https://www.agenciasinc.es/Sindicacion/Noticias",
		"https://www.elcomercio.pe/rss/",
	}

	log.Println("🗑️  Removing broken feeds...")

	for _, url := range brokenURLs {
		// Soft delete using GORM's DeletedAt
		result := db.Where("url = ?", url).Delete(&models.Feed{})
		if result.Error != nil {
			log.Printf("❌ Error deleting %s: %v", url, result.Error)
		} else if result.RowsAffected == 0 {
			log.Printf("⚠️  Feed not found (already deleted?): %s", url)
		} else {
			log.Printf("✅ Deleted feed: %s", url)
		}
	}

	log.Println("✅ Feed cleanup complete.")
}

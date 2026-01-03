package main

import (
	"log"
	"vidit/internal/database"
	"vidit/internal/models"
)

func main() {
	dbConfig := database.Config{
		Host:     "localhost",
		Port:     "5432",
		User:     "postgres",
		Password: "postgres",
		DBName:   "vidit",
		SSLMode:  "disable",
	}

	if err := database.Connect(dbConfig); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := database.Migrate(); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	feeds := []models.Feed{
		{Name: "El Mostrador", URL: "https://www.elmostrador.cl/feed/", ColorHex: "#E63946"},
		{Name: "El Desconcierto", URL: "https://www.eldesconcierto.cl/feed/", ColorHex: "#F77F00"},
		{Name: "El Ciudadano", URL: "https://www.elciudadano.com/feed/", ColorHex: "#06AED5"},
		{Name: "Interferencia", URL: "https://interferencia.cl/feed", ColorHex: "#1D3557"},
		{Name: "The Clinic", URL: "https://www.theclinic.cl/feed/", ColorHex: "#9D4EDD"},
		{Name: "CIPER Chile", URL: "https://www.ciperchile.cl/feed/", ColorHex: "#2A9D8F"},
	}

	globalFeeds := []models.Feed{
		{Name: "El País", URL: "https://feeds.elpais.com/mrss-s/pages/ep/site/elpais.com/portada", ColorHex: "#004481"},
		{Name: "BBC Mundo", URL: "https://feeds.bbci.co.uk/mundo/rss.xml", ColorHex: "#BB1919"},
		{Name: "El Mundo", URL: "https://e00-elmundo.uecdn.es/elmundo/rss/portada.xml", ColorHex: "#2E6D9D"},
		{Name: "Clarín", URL: "https://www.clarin.com/rss/lo-ultimo/", ColorHex: "#FF0000"},
		{Name: "CNN Español", URL: "https://cnnespanol.cnn.com/feed/", ColorHex: "#CC0000"},
		{Name: "La Nación", URL: "https://www.lanacion.com.ar/arc/outboundfeeds/rss/?outputType=xml", ColorHex: "#0060A9"},
		{Name: "El Universal", URL: "https://www.eluniversal.com.mx/rss.xml", ColorHex: "#231F20"},
		{Name: "Infobae", URL: "https://www.infobae.com/feeds/rss/", ColorHex: "#FA6900"},
		{Name: "DW Español", URL: "https://rss.dw.com/xml/rss-sp-all", ColorHex: "#002D5A"},
		{Name: "El Tiempo", URL: "https://www.eltiempo.com/rss/mundo.xml", ColorHex: "#17479E"},
	}

	feeds = append(feeds, globalFeeds...)

	agencyFeeds := []models.Feed{
		{Name: "Europa Press", URL: "https://www.europapress.es/rss/rss.aspx", ColorHex: "#F78F1E"},
		{Name: "ONU News", URL: "https://news.un.org/feed/subscribe/es/news/all/rss.xml", ColorHex: "#009EDB"},
		{Name: "France 24", URL: "https://www.france24.com/es/rss", ColorHex: "#00A6EB"},
		{Name: "Euronews", URL: "https://es.euronews.com/rss", ColorHex: "#003D8F"},
		{Name: "SwissInfo", URL: "https://www.swissinfo.ch/oai/rss/es/index.xml", ColorHex: "#FF0000"},
		{Name: "Agencia SINC", URL: "https://www.agenciasinc.es/Sindicacion/Noticias", ColorHex: "#9DC63F"},
		{Name: "PR Newswire", URL: "https://www.prnewswire.com/rss/news-releases-list.rss?lang=es", ColorHex: "#ED1C24"},
		{Name: "Voz de América", URL: "https://www.vozdeamerica.com/api/z$gqqye_qvi_", ColorHex: "#003366"}, // URL updated
		{Name: "Democracy Now", URL: "https://www.democracynow.org/es/datos/noticias.xml", ColorHex: "#B41F25"},
		{Name: "Global Voices", URL: "https://es.globalvoices.org/feed/", ColorHex: "#2F3C44"},
	}
	feeds = append(feeds, agencyFeeds...)

	latamFeeds := []models.Feed{
		{Name: "RT en Español", URL: "https://actualidad.rt.com/feeds/all.rss", ColorHex: "#66CC00"},
		{Name: "Somos Télam", URL: "https://somostelam.com.ar/feed/", ColorHex: "#3399FF"},
		{Name: "Telesur", URL: "https://www.telesurtv.net/rss/rss.xml", ColorHex: "#DA251D"},
		{Name: "Infobae (LatAm)", URL: "https://www.infobae.com/feeds/rss/", ColorHex: "#FA6900"},
	}
	feeds = append(feeds, latamFeeds...)

	for _, feed := range feeds {
		result := database.DB.FirstOrCreate(&feed, models.Feed{URL: feed.URL})
		if result.Error != nil {
			log.Printf("Error inserting feed %s: %v\n", feed.Name, result.Error)
		} else {
			log.Printf("✅ Added feed: %s\n", feed.Name)
		}
	}

	log.Println("\n🎉 Database seeded successfully!")
	log.Println("Run the server with: go run cmd/server/main.go cmd/server/renderer.go")
}

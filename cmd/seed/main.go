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

	// Define all feeds with Categories and Countries
	feeds := []models.Feed{
		// --- Spain / Politics (International) ---
		{Name: "Politico (ES)", URL: "https://www.politico.eu/tag/spanish-politics/feed/", ColorHex: "#0C3C60", Category: "international", Country: "ES"},
		{Name: "El Periódico", URL: "https://www.elperiodico.com/es/rss/politica/rss.xml", ColorHex: "#005696", Category: "international", Country: "ES"},
		{Name: "elDiario.es", URL: "https://www.eldiario.es/rss/", ColorHex: "#121212", Category: "international", Country: "ES"},
		{Name: "InfoLibre", URL: "https://www.infolibre.es/rss/", ColorHex: "#D92B34", Category: "international", Country: "ES"},
		{Name: "La Vanguardia", URL: "https://www.lavanguardia.com/rss/home.xml", ColorHex: "#000000", Category: "international", Country: "ES"},
		{Name: "Vozpópuli", URL: "https://www.vozpopuli.com/rss/", ColorHex: "#C3002F", Category: "international", Country: "ES"},

		// --- Chile (LatAm) ---
		{Name: "El Mostrador", URL: "https://www.elmostrador.cl/feed/", ColorHex: "#E63946", Category: "latam", Country: "CL"},
		{Name: "El Desconcierto", URL: "https://www.eldesconcierto.cl/feed/", ColorHex: "#F77F00", Category: "latam", Country: "CL"},
		{Name: "El Ciudadano", URL: "https://www.elciudadano.com/feed/", ColorHex: "#06AED5", Category: "latam", Country: "CL"},
		{Name: "Interferencia", URL: "https://interferencia.cl/feed", ColorHex: "#1D3557", Category: "latam", Country: "CL"},
		{Name: "The Clinic", URL: "https://www.theclinic.cl/feed/", ColorHex: "#9D4EDD", Category: "latam", Country: "CL"},
		{Name: "CIPER Chile", URL: "https://www.ciperchile.cl/feed/", ColorHex: "#2A9D8F", Category: "latam", Country: "CL"},
		{Name: "La Tercera", URL: "https://www.latercera.com/rss", ColorHex: "#000000", Category: "latam", Country: "CL"},
		{Name: "BioBioChile", URL: "https://www.biobiochile.cl/feed", ColorHex: "#FFC300", Category: "latam", Country: "CL"},
		{Name: "ADN Radio", URL: "https://www.adnradio.cl/arc/outboundfeeds/rss/", ColorHex: "#E71D25", Category: "latam", Country: "CL"},
		{Name: "Turno (Copano)", URL: "https://copano.news/sitemap.xml", Type: "sitemap", ColorHex: "#000000", Category: "latam", Country: "CL"},
		{Name: "Radio Agricultura", URL: "https://www.radioagricultura.cl/sitemap_news.xml", Type: "sitemap", ColorHex: "#2E7D32", Category: "latam", Country: "CL"},

		// --- Global / International ---
		{Name: "El País", URL: "https://feeds.elpais.com/mrss-s/pages/ep/site/elpais.com/portada", ColorHex: "#004481", Category: "international", Country: "ES"},
		{Name: "BBC Mundo", URL: "https://feeds.bbci.co.uk/mundo/rss.xml", ColorHex: "#BB1919", Category: "international", Country: "ES"},
		{Name: "El Mundo", URL: "https://www.elmundo.es/sitemap_news.xml", Type: "sitemap", ColorHex: "#2E6D9D", Category: "international", Country: "ES"},
		{Name: "DW Español", URL: "https://rss.dw.com/xml/rss-sp-all", ColorHex: "#002D5A", Category: "international", Country: "INT"},
		{Name: "ONU News", URL: "https://news.un.org/feed/subscribe/es/news/all/rss.xml", ColorHex: "#009EDB", Category: "international", Country: "INT"},
		{Name: "France 24", URL: "https://www.france24.com/es/rss", ColorHex: "#00A6EB", Category: "international", Country: "INT"},
		{Name: "Euronews", URL: "https://es.euronews.com/rss", ColorHex: "#003D8F", Category: "international", Country: "INT"},
		{Name: "SwissInfo", URL: "https://www.swissinfo.ch/oai/rss/es/index.xml", ColorHex: "#FF0000", Category: "international", Country: "INT"},
		{Name: "Global Voices", URL: "https://es.globalvoices.org/feed/", ColorHex: "#2F3C44", Category: "international", Country: "INT"},
		{Name: "RT en Español", URL: "https://actualidad.rt.com/feeds/all.rss", ColorHex: "#66CC00", Category: "international", Country: "INT"},

		// --- USA ---
		{Name: "CNN Español", URL: "https://cnnespanol.cnn.com/feed/", ColorHex: "#CC0000", Category: "usa", Country: "US"},
		{Name: "Voz de América", URL: "https://www.vozdeamerica.com/api/z$gqqye_qvi_", ColorHex: "#003366", Category: "usa", Country: "US"},
		{Name: "Democracy Now", URL: "https://www.democracynow.org/es/datos/noticias.xml", ColorHex: "#B41F25", Category: "usa", Country: "US"},
		{Name: "Fox Deportes", URL: "https://www.foxdeportes.com/rss/home.xml", ColorHex: "#003366", Category: "usa", Country: "US"},

		// --- LatAm (General) ---
		{Name: "Clarín", URL: "https://www.clarin.com/rss/lo-ultimo/", ColorHex: "#FF0000", Category: "latam", Country: "AR"},
		// {Name: "La Nación", URL: "https://www.lanacion.com.ar/arc/outboundfeeds/rss/?outputType=xml", ColorHex: "#0060A9", Category: "latam", Country: "AR"},Removed
		{Name: "El Universal", URL: "https://www.eluniversal.com.mx/rss.xml", ColorHex: "#231F20", Category: "latam", Country: "MX"},
		{Name: "Infobae", URL: "https://www.infobae.com/feeds/rss/", ColorHex: "#FA6900", Category: "latam", Country: "AR"},
		{Name: "El Tiempo", URL: "https://www.eltiempo.com/rss/mundo.xml", ColorHex: "#17479E", Category: "latam", Country: "CO"},
		{Name: "Somos Télam", URL: "https://somostelam.com.ar/feed/", ColorHex: "#3399FF", Category: "latam", Country: "AR"},
		{Name: "Telesur", URL: "https://www.telesurtv.net/rss/rss.xml", ColorHex: "#DA251D", Category: "latam", Country: "VE"},
		{Name: "Reforma (MX)", URL: "https://www.reforma.com/rss/portada.xml", ColorHex: "#FF6600", Category: "latam", Country: "MX"},
		{Name: "La Jornada (MX)", URL: "https://www.jornada.com.mx/rss/edicion.xml", ColorHex: "#D4AF37", Category: "latam", Country: "MX"},
		{Name: "El Comercio (PE)", URL: "https://www.elcomercio.pe/rss/", ColorHex: "#FFCC00", Category: "latam", Country: "PE"},
		{Name: "Diario Red", URL: "https://diariored.canalred.tv/feed/", ColorHex: "#E60000", Category: "latam", Country: "ES"}, // Diario Red is Spanish but covers LatAm too, keeping ES or Intl? Let's say ES.
		{Name: "Olé", URL: "https://www.ole.com.ar/rss/ultimas-noticias", ColorHex: "#9ACD32", Category: "latam", Country: "AR"},

		// --- Agency ---
		{Name: "Europa Press", URL: "https://www.europapress.es/rss/rss.aspx", ColorHex: "#F78F1E", Category: "general", Country: "ES"},
		{Name: "Agencia SINC", URL: "https://www.agenciasinc.es/Sindicacion/Noticias", ColorHex: "#9DC63F", Category: "general", Country: "ES"},

		// --- Sports (General) ---
		{Name: "Marca", URL: "https://e00-marca.uecdn.es/rss/portada.xml", ColorHex: "#CC0000", Category: "general", Country: "ES"},
		{Name: "Sport", URL: "https://www.sport.es/es/rss/last-news/rss.xml", ColorHex: "#C8102E", Category: "general", Country: "ES"},

		// --- cybersecurity ---
		{Name: "El Lado del Mal", URL: "https://www.elladodelmal.com/feeds/posts/default", ColorHex: "#000000", Category: "cybersecurity", Country: "ES"},
		{Name: "Security By Default", URL: "http://feeds.feedburner.com/SecurityByDefault", ColorHex: "#333333", Category: "cybersecurity", Country: "ES"},
		{Name: "DragonJAR", URL: "https://www.dragonjar.org/feed", ColorHex: "#990000", Category: "cybersecurity", Country: "CO"}, // Colombian
		{Name: "Una al Día", URL: "https://unaaldia.hispasec.com/feed", ColorHex: "#0066CC", Category: "cybersecurity", Country: "ES"},
		{Name: "Segu-Info", URL: "https://blog.segu-info.com.ar/feeds/posts/default", ColorHex: "#FF6600", Category: "cybersecurity", Country: "AR"},
		{Name: "HackPlayers", URL: "https://www.hackplayers.com/feeds/posts/default", ColorHex: "#000000", Category: "cybersecurity", Country: "ES"},
		{Name: "WeLiveSecurity", URL: "https://www.welivesecurity.com/la-es/feed/", ColorHex: "#0099CC", Category: "cybersecurity", Country: "INT"},
		{Name: "Kaspersky Blog", URL: "https://www.kaspersky.es/blog/feed/", ColorHex: "#006D55", Category: "cybersecurity", Country: "INT"},
		{Name: "INCIBE Avisos", URL: "https://www.incibe.es/feed/avisos-seguridad", ColorHex: "#00537F", Category: "cybersecurity", Country: "ES"},
		{Name: "Ciberseguridad Blog", URL: "https://ciberseguridad.blog/feed/", ColorHex: "#444444", Category: "cybersecurity", Country: "ES"},
		// New
		{Name: "Derecho de la Red", URL: "https://www.derechodelared.com/feed/", ColorHex: "#000000", Category: "cybersecurity", Country: "ES"},
		{Name: "CyberSecurity News", URL: "https://cybersecuritynews.es/feed/", ColorHex: "#000000", Category: "cybersecurity", Country: "ES"},
		{Name: "RedesZone", URL: "https://www.redeszone.net/feed/", ColorHex: "#000000", Category: "cybersecurity", Country: "ES"},
		{Name: "Genbeta Seguridad", URL: "https://www.genbeta.com/categoria/seguridad/rss2.xml", ColorHex: "#000000", Category: "cybersecurity", Country: "ES"},
		{Name: "MuySeguridad", URL: "https://www.muyseguridad.net/feed/", ColorHex: "#000000", Category: "cybersecurity", Country: "ES"},
		{Name: "Xataka Seguridad", URL: "https://www.xatakandroid.com/categoria/seguridad/rss2.xml", ColorHex: "#000000", Category: "cybersecurity", Country: "ES"},

		// -- China ---
		// (Xinhua & CGTN failed validation or returned 404/unknown. Adding them provisionally or searching better URLs?
		//  Probe results: Xinhua 404, Pueblo 404. Let's try to add the validated ones if any, or skip to avoid errors.
		//  Actually, if I don't add them, 'red' hover won't show.
		//  Validator said: "Xinhua Español... 404".
		//  Let's try: http://spanish.xinhuanet.com/rss/index.xml failed.
		//  Let's try keeping them commented out or search again? I will add a placeholder for now to test the UI logic, maybe using a generic text)
	}

	for _, feed := range feeds {
		// Use Assign to update Category and Country for existing feeds
		result := database.DB.Where(models.Feed{URL: feed.URL}).Assign(models.Feed{Category: feed.Category, Country: feed.Country}).FirstOrCreate(&feed, models.Feed{URL: feed.URL})
		if result.Error != nil {
			log.Printf("Error inserting feed %s: %v\n", feed.Name, result.Error)
		} else {
			// log.Printf("✅ Added/Updated feed: %s\n", feed.Name)
		}
	}

	log.Println("\n🎉 Database seeded successfully!")
}

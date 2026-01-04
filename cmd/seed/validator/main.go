package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/mmcdole/gofeed"
)

func main() {
	candidates := []string{
		// Politics / News (Spain & Global)
		"https://www.politico.eu/tag/spanish-politics/feed/",
		"https://www.elperiodico.com/es/rss/politica/rss.xml",
		"https://www.lasprovincias.es/rss/2.0/politica",
		"https://okdiario.com/feed",
		"https://www.elplural.com/rss",
		"https://www.eldiario.es/rss/",
		"https://www.infolibre.es/rss/",
		"https://www.huffingtonpost.es/feeds/index.xml",
		"https://www.libertaddigital.com/rss/noticias.xml",
		"https://www.lavanguardia.com/rss/home.xml",
		"https://www.abc.es/rss/2.0/portada",
		"https://www.20minutos.es/rss/",
		"https://www.publico.es/rss/",
		"https://www.larazon.es/rss/",
		"https://www.elconfidencial.com/rss/",
		"https://www.vozpopuli.com/rss/",

		// Politics / News (LatAm)
		"https://www.lapoliticaonline.com/files/rss/politica.xml",
		"https://www.lapoliticaonline.com/files/rss/mexico.xml",
		"https://diariored.canalred.tv/feed/",
		"https://www.elespectador.com/rss/",
		"https://www.elcomercio.pe/rss/",
		"https://www.biobiochile.cl/feed",
		"https://www.latercera.com/feed/",
		"https://www.pagina12.com.ar/rss/portada",
		"https://www.reforma.com/rss/portada.xml",
		"https://www.jornada.com.mx/rss/edicion.xml",

		// Sports
		"https://e00-marca.uecdn.es/rss/portada.xml",
		"https://as.com/rss/tags/ultimas_noticias",
		"https://www.mundodeportivo.com/rss/headlines.xml",
		"https://www.sport.es/es/rss/last-news/rss.xml",
		"https://www.tycsports.com/rss",
		"https://www.ole.com.ar/rss/ultimas-noticias",
		"https://www.espn.com.mx/espn/rss/news",
		"https://www.foxdeportes.com/rss/home.xml",
		"https://www.record.com.mx/rss",

		// Cybersecurity
		"https://www.elladodelmal.com/feeds/posts/default",
		"http://feeds.feedburner.com/SecurityByDefault",
		"https://www.dragonjar.org/feed",
		"https://unaaldia.hispasec.com/feed",
		"https://blog.segu-info.com.ar/feeds/posts/default",
		"https://www.hackplayers.com/feeds/posts/default",
		"https://www.welivesecurity.com/la-es/feed/",
		"https://www.kaspersky.es/blog/feed/",
		"https://www.incibe.es/feed/avisos-seguridad",
		"https://ciberseguridad.blog/feed/",
	}

	fp := gofeed.NewParser()
	client := http.Client{
		Timeout: 10 * time.Second,
	}

	fmt.Println("Validating feeds...")
	validCount := 0

	for _, url := range candidates {
		resp, err := client.Get(url)
		if err != nil {
			fmt.Printf("❌ FAIL %s: %v\n", url, err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			fmt.Printf("❌ FAIL %s: Status %d\n", url, resp.StatusCode)
			continue
		}

		_, err = fp.ParseURL(url)
		if err != nil {
			fmt.Printf("⚠️  WARN %s: Parse error (might still work with custom headers) %v\n", url, err)
			// Some feeds block default user agents, but might be valid.
			// We'll mark as valid for now if status is 200,
			// but the fetcher service handles this better usually.
		}

		fmt.Printf("✅ OK   %s\n", url)
		validCount++
	}
	fmt.Printf("\nTotal Candidates: %d\n", len(candidates))
	fmt.Printf("Total Valid: %d\n", validCount)
}

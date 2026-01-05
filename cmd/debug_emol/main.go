package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type NewsAPIResponse struct {
	Status       string `json:"status"`
	TotalResults int    `json:"totalResults"`
	Articles     []struct {
		Title string `json:"title"`
		URL   string `json:"url"`
	} `json:"articles"`
}

func main() {
	apiKey := os.Getenv("NEWSAPI_KEY")
	if apiKey == "" {
		log.Fatal("NEWSAPI_KEY missing")
	}

	domains := []string{"emol.com", "www.emol.com"}

	for _, d := range domains {
		url := fmt.Sprintf("https://newsapi.org/v2/everything?domains=%s&apiKey=%s&pageSize=5", d, apiKey)
		log.Printf("Testing domain: %s", d)

		resp, err := http.Get(url)
		if err != nil {
			log.Printf("Error: %v", err)
			continue
		}
		defer resp.Body.Close()

		var result NewsAPIResponse
		json.NewDecoder(resp.Body).Decode(&result)

		log.Printf("Found %d articles for %s", result.TotalResults, d)
		if len(result.Articles) > 0 {
			log.Printf(" - Sample: %s", result.Articles[0].Title)
		}
	}
}

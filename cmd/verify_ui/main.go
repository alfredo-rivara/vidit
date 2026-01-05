package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"time"
	"vidit/internal/models"
)

func main() {
	// Mock FuncMap (same as main.go)
	funcMap := template.FuncMap{
		"spanishDate": func(t time.Time) string {
			return "01 Ene 12:00"
		},
		"formatScore": func(score float64) string {
			return fmt.Sprintf("%.2f", score)
		},
	}

	tmpl, err := template.New("").Funcs(funcMap).ParseGlob("views/*.html")
	if err != nil {
		log.Fatalf("❌ Template parsing failed: %v", err)
	}

	// Mock Data
	data := map[string]interface{}{
		"Articles": []models.Article{
			{Title: "Test Article", Score: 1.23456, Feed: models.Feed{Type: "rss", Name: "TestFeed"}},
		},
		"Count": 1,
	}

	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, "index.html", data); err != nil {
		log.Fatalf("❌ Template execution failed: %v", err)
	}

	output := buf.String()
	// Check if score is rendered
	expected := "1.23"
	if bytes.Contains(buf.Bytes(), []byte(expected)) {
		log.Println("✅ Success: Score 1.23 found in rendered HTML.")
	} else {
		log.Printf("❌ Failure: Score 1.23 not found in rendered output.\nSample output: %s", output[:500])
	}
}

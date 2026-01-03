package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"time"
	"vidit/internal/database"
	"vidit/internal/fetcher"
	"vidit/internal/models"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// TemplateRenderer is a custom html/template renderer for Echo
type TemplateRenderer struct {
	templates *template.Template
}

// Render renders a template document
func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {

	dbConfig := database.Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "postgres"),
		DBName:   getEnv("DB_NAME", "vidit"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}

	if err := database.Connect(dbConfig); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := database.Migrate(); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	e := echo.New()
	e.HideBanner = true

	funcMap := template.FuncMap{
		"spanishDate": func(t time.Time) string {
			months := map[time.Month]string{
				time.January:   "Ene",
				time.February:  "Feb",
				time.March:     "Mar",
				time.April:     "Abr",
				time.May:       "May",
				time.June:      "Jun",
				time.July:      "Jul",
				time.August:    "Ago",
				time.September: "Sep",
				time.October:   "Oct",
				time.November:  "Nov",
				time.December:  "Dic",
			}
			return fmt.Sprintf("%02d %s %02d:%02d", t.Day(), months[t.Month()], t.Hour(), t.Minute())
		},
	}

	e.Renderer = &TemplateRenderer{
		templates: template.Must(template.New("").Funcs(funcMap).ParseGlob("views/*.html")),
	}

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Static("/css", "public/css")

	e.GET("/", handleHome)
	e.POST("/fetch", handleFetch)

	port := getEnv("PORT", "3000")
	log.Printf("🚀 Vidit server starting on http://localhost:%s\n", port)
	e.Logger.Fatal(e.Start(":" + port))
}

func handleHome(c echo.Context) error {
	var articles []models.Article

	result := database.DB.
		Preload("Feed").
		Order("score DESC, published_at DESC").
		Limit(500).
		Find(&articles)

	if result.Error != nil {
		return c.String(http.StatusInternalServerError, "Error loading articles")
	}

	return c.Render(http.StatusOK, "index.html", map[string]interface{}{
		"Articles": articles,
		"Count":    len(articles),
	})
}

// function getTrendingTopics deleted

func handleFetch(c echo.Context) error {
	service := fetcher.NewService()

	if err := service.FetchAllFeeds(database.DB); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"status": "Feeds fetched successfully",
	})
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

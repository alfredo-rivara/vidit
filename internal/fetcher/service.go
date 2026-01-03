package fetcher

import (
	"log"
	"strings"
	"sync"
	"time"
	"unicode"
	"vidit/internal/database"
	"vidit/internal/models"

	"github.com/mmcdole/gofeed"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Service struct {
	parser *gofeed.Parser
}

func NewService() *Service {
	return &Service{
		parser: gofeed.NewParser(),
	}
}

var spanishStopWords = map[string]bool{
	"el": true, "la": true, "de": true, "en": true, "y": true,
	"a": true, "los": true, "las": true, "un": true, "una": true,
	"del": true, "por": true, "con": true, "para": true, "es": true,
	"al": true, "lo": true, "su": true, "se": true, "como": true,
	"que": true, "más": true, "o": true, "ha": true, "fue": true,
	"le": true, "tras": true, "ante": true, "sobre": true, "sin": true,
	"entre": true, "sus": true, "no": true,
}

type FeedItem struct {
	Title       string
	URL         string
	PublishedAt time.Time
	FeedID      uint
}

func (s *Service) FetchAllFeeds(db *gorm.DB) error {
	var feeds []models.Feed
	if err := db.Find(&feeds).Error; err != nil {
		return err
	}

	if len(feeds) == 0 {
		log.Println("⚠️  No feeds found in database")
		return nil
	}

	log.Printf("🔄 Fetching %d feeds concurrently...\n", len(feeds))

	var wg sync.WaitGroup
	itemsChan := make(chan []models.Article, len(feeds))

	for _, feed := range feeds {
		wg.Add(1)
		go func(f models.Feed) {
			defer wg.Done()
			articles := s.fetchFeed(f)
			itemsChan <- articles
		}(feed)
	}

	go func() {
		wg.Wait()
		close(itemsChan)
	}()

	var allArticles []models.Article
	for articles := range itemsChan {
		allArticles = append(allArticles, articles...)
	}

	log.Printf("✅ Fetched %d raw articles\n", len(allArticles))

	uniqueArticlesMap := make(map[string]models.Article)
	for _, a := range allArticles {
		uniqueArticlesMap[a.URL] = a
	}

	var uniqueArticles []models.Article
	for _, a := range uniqueArticlesMap {
		uniqueArticles = append(uniqueArticles, a)
	}

	log.Printf("🔹 processed %d unique articles from raw list\n", len(uniqueArticles))

	scoredArticles := s.calculateScores(uniqueArticles)

	if len(scoredArticles) > 0 {
		return s.saveArticles(db, scoredArticles)
	}

	return nil
}

func (s *Service) fetchFeed(feed models.Feed) []models.Article {
	feedData, err := s.parser.ParseURL(feed.URL)
	if err != nil {
		log.Printf("❌ Error fetching %s (%s): %v\n", feed.Name, feed.URL, err)
		return nil
	}

	articles := make([]models.Article, 0, len(feedData.Items))
	for _, item := range feedData.Items {
		publishedAt := time.Now()
		if item.PublishedParsed != nil {
			publishedAt = *item.PublishedParsed
		}

		articles = append(articles, models.Article{
			Title:       item.Title,
			URL:         item.Link,
			PublishedAt: publishedAt,
			FeedID:      feed.ID,
			Score:       1, // Default score
		})
	}

	now := time.Now()
	database.DB.Model(&feed).Update("last_fetched_at", &now)

	return articles
}

func (s *Service) calculateScores(articles []models.Article) []models.Article {
	keywordFeedMap := make(map[string]map[uint]bool)

	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	normalize := func(input string) string {
		output, _, _ := transform.String(t, input)
		return strings.ToLower(output)
	}

	for _, article := range articles {
		normalizedTitle := normalize(article.Title)
		words := strings.Fields(normalizedTitle)

		for _, word := range words {
			word = strings.TrimFunc(word, func(r rune) bool {
				return !unicode.IsLetter(r) && !unicode.IsNumber(r)
			})

			if len(word) < 3 || spanishStopWords[word] {
				continue
			}

			if keywordFeedMap[word] == nil {
				keywordFeedMap[word] = make(map[uint]bool)
			}
			keywordFeedMap[word][article.FeedID] = true
		}
	}

	scoredArticles := make([]models.Article, len(articles))
	copy(scoredArticles, articles)

	for i, article := range scoredArticles {
		normalizedTitle := normalize(article.Title)
		words := strings.Fields(normalizedTitle)

		currentMaxScore := 1

		for _, word := range words {
			word = strings.TrimFunc(word, func(r rune) bool {
				return !unicode.IsLetter(r) && !unicode.IsNumber(r)
			})

			if len(word) < 3 || spanishStopWords[word] {
				continue
			}

			feedCount := len(keywordFeedMap[word])

			if feedCount >= 3 {
				currentMaxScore = 3 // Giant
				break               // Already max score, no need to check other words
			} else if feedCount >= 2 && currentMaxScore < 2 {
				currentMaxScore = 2 // Large
			}
		}

		scoredArticles[i].Score = currentMaxScore
	}

	return scoredArticles
}

func (s *Service) saveArticles(db *gorm.DB, articles []models.Article) error {
	batchSize := 100

	result := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "url"}},
		DoUpdates: clause.AssignmentColumns([]string{"title", "score", "published_at", "updated_at"}),
	}).CreateInBatches(&articles, batchSize)

	if result.Error != nil {
		log.Printf("❌ Error saving articles: %v\n", result.Error)
		return result.Error
	}

	score1 := 0
	score2 := 0
	score3 := 0

	for _, a := range articles {
		if a.Score == 1 {
			score1++
		}
		if a.Score == 2 {
			score2++
		}
		if a.Score == 3 {
			score3++
		}
	}

	log.Printf("💾 Saved/Updated %d articles (Score 1: %d, Score 2: %d, Score 3: %d)\n",
		len(articles), score1, score2, score3)

	return nil
}

# Vidit - Living News Mosaic

A news aggregator with a unique **keyword clustering algorithm** that replaces traditional engagement metrics. Instead of retweets or likes, relevance is calculated by how many RSS feeds mention the same keywords concurrently.

## 🎯 Core Concept

News cards in a **living mosaic** where size represents relevance:
- **Giant Cards (Score 3)**: Topics mentioned in 3+ feeds
- **Large Cards (Score 2)**: Topics mentioned in 2 feeds  
- **Normal Cards (Score 1)**: Unique news from single feeds

## 🛠️ Tech Stack

- **Backend**: Go 1.25+
- **Web Framework**: Echo v4
- **ORM**: GORM with PostgreSQL
- **RSS Parsing**: gofeed
- **Frontend**: HTML templates + CSS Grid (Masonry)

## 📁 Project Structure

```
vidit/
├── cmd/
│   ├── server/           # Main application entry
│   │   ├── main.go
│   │   └── renderer.go
│   └── seed/             # Database seeder
│       └── main.go
├── internal/
│   ├── models/           # GORM models
│   │   ├── feed.go
│   │   └── article.go
│   ├── database/         # DB connection
│   │   └── database.go
│   └── fetcher/          # RSS fetching & scoring
│       └── service.go
├── views/                # HTML templates
│   └── index.html
├── public/
│   └── css/
│       └── style.css
└── .env.example
```

## 🚀 Quick Start

### 1. Prerequisites

- Go 1.25+
- PostgreSQL running locally

### 2. Database Setup

```bash
# Create database
createdb vidit

# Optional: Copy and configure environment
cp .env.example .env
# Edit .env with your PostgreSQL credentials
```

### 3. Seed Sample Feeds

```bash
go run cmd/seed/main.go
```

This will populate the database with popular Spanish/Argentine news sources.

### 4. Run the Server

```bash
go run cmd/server/main.go cmd/server/renderer.go
```

Visit: **http://localhost:3000**

### 5. Fetch News

Click the **"🔄 Refresh Feeds"** button in the UI or use:

```bash
curl -X POST http://localhost:3000/fetch
```

## 🧠 The "Vidit" Algorithm

1. **Concurrent Fetching**: All RSS feeds are fetched in parallel using goroutines
2. **Keyword Extraction**: Titles are normalized (lowercase, Spanish stop-words removed)
3. **Clustering**: Keywords are mapped across all feeds to count occurrences
4. **Scoring**:
   - If a keyword appears in 3+ feeds → Score 3 (Giant)
   - If a keyword appears in 2 feeds → Score 2 (Large)
   - Unique news → Score 1 (Normal)
5. **Upsert**: Articles are saved with conflict resolution on URL

## 🎨 Frontend Features

- **CSS Grid Masonry Layout** with `grid-auto-flow: dense`
- Cards span 1-2 columns and 1-2 rows based on score
- Feed color-coded badges
- Responsive design
- Smooth hover animations

## 📝 Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `DB_HOST` | localhost | PostgreSQL host |
| `DB_PORT` | 5432 | PostgreSQL port |
| `DB_USER` | postgres | Database user |
| `DB_PASSWORD` | postgres | Database password |
| `DB_NAME` | vidit | Database name |
| `DB_SSLMODE` | disable | SSL mode |
| `PORT` | 3000 | Server port |

## 🔧 Development

### Add New Feed

```go
feed := models.Feed{
    Name:     "Source Name",
    URL:      "https://example.com/rss",
    ColorHex: "#ff5733",
}
database.DB.Create(&feed)
```

### Manual Fetch via Code

```go
service := fetcher.NewService()
service.FetchAllFeeds()
```

## 📦 Dependencies

All dependencies are managed via Go modules:

```bash
go get github.com/labstack/echo/v4
go get gorm.io/gorm
go get gorm.io/driver/postgres
go get github.com/mmcdole/gofeed
```

## 🏗️ Architecture Highlights

- **Concurrent Processing**: Uses `sync.WaitGroup` for parallel RSS fetching
- **Stop-word Filtering**: Spanish stop-words are removed during keyword extraction
- **Upsert Logic**: Prevents duplicate articles using GORM's `OnConflict` clause
- **Score-based Layout**: CSS Grid dynamically sizes cards based on relevance score

## 📜 License

MIT

---

**Built with ❤️ using Go and Echo**

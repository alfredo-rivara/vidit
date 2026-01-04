package models

import (
	"time"

	"gorm.io/gorm"
)

// Feed represents an RSS feed source
type Feed struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Name          string     `gorm:"not null" json:"name"`
	URL           string     `gorm:"not null;unique" json:"url"`
	Type          string     `gorm:"default:'rss'" json:"type"`         // rss or sitemap
	Category      string     `gorm:"default:'general'" json:"category"` // cybersecurity, international, latam, usa, china
	ColorHex      string     `gorm:"default:#3b82f6" json:"color_hex"`  // Default blue
	LastFetchedAt *time.Time `json:"last_fetched_at"`

	// Relationships
	Articles []Article `gorm:"foreignKey:FeedID" json:"-"`
}

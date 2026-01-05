package models

import (
	"time"

	"gorm.io/gorm"
)

// Article represents a news article from an RSS feed
type Article struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Title       string    `gorm:"not null" json:"title"`
	URL         string    `gorm:"not null;unique" json:"url"`
	PublishedAt time.Time `gorm:"index" json:"published_at"`
	Score       float64   `gorm:"default:1.0;index" json:"score"`

	// Foreign Key
	FeedID uint `gorm:"not null;index" json:"feed_id"`
	Feed   Feed `gorm:"foreignKey:FeedID" json:"feed"`
}

package database

import (
	"fmt"
	"log"
	"vidit/internal/models"
	
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// Config holds database connection parameters
type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// Connect initializes the database connection
func Connect(cfg Config) error {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode,
	)
	
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	
	log.Println("✅ Database connection established")
	return nil
}

// Migrate runs auto-migration for all models
func Migrate() error {
	err := DB.AutoMigrate(
		&models.Feed{},
		&models.Article{},
	)
	
	if err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}
	
	log.Println("✅ Database migration completed")
	return nil
}

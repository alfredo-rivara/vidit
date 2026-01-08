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

	// Soft delete La Nación
	result := database.DB.Where("name = ?", "La Nación").Delete(&models.Feed{})
	if result.Error != nil {
		log.Fatalf("Error deleting feed: %v", result.Error)
	}

	if result.RowsAffected > 0 {
		log.Printf("Successfully removed La Nación (Rows affected: %d)", result.RowsAffected)
	} else {
		log.Println("Feed 'La Nación' not found or already deleted.")
	}
}

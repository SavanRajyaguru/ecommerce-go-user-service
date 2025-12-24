package migrations

import (
	"log"

	"github.com/SavanRajyaguru/ecommerce-go-user-service/internal/database"
	"github.com/SavanRajyaguru/ecommerce-go-user-service/models"
)

func RunMigrations() {
	if database.DB == nil {
		log.Fatal("Database connection not initialized")
	}

	err := database.DB.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
	log.Println("Database Migration Completed")
}

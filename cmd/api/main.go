package main

import (
	"log"

	"github.com/SavanRajyaguru/ecommerce-go-user-service/api"
	"github.com/SavanRajyaguru/ecommerce-go-user-service/config"
	"github.com/SavanRajyaguru/ecommerce-go-user-service/internal/cache"
	"github.com/SavanRajyaguru/ecommerce-go-user-service/internal/database"
	"github.com/SavanRajyaguru/ecommerce-go-user-service/migrations"
	"github.com/SavanRajyaguru/ecommerce-go-user-service/pkg/logger"
)

func main() {
	// 1. Initialize Logger
	logger.InitLogger()
	log.Println("Starting User Service...")

	// 2. Load Configuration (Env + Config Service)
	config.LoadConfig()

	// 3. Initialize Database
	database.ConnectDB()

	// 4. Run Migrations
	// In production, this might be a separate logic/flag.
	migrations.RunMigrations()

	// 5. Initialize Redis
	cache.ConnectRedis()

	// 6. Start HTTP Server
	router := api.SetupRouter()

	port := config.AppConfig.AppPort
	if port == "" {
		port = "8080"
	}

	log.Printf("Server running on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}

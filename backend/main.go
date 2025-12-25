package main

import (
	"log"

	"health-tracker/config"
	"health-tracker/database"
	"health-tracker/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	config.LoadConfig()

	// Set Gin mode
	gin.SetMode(config.AppConfig.GinMode)

	// Initialize database
	database.InitDatabase()

	// Create Gin router
	r := gin.Default()

	// Setup routes
	routes.SetupRoutes(r)

	// Start server
	log.Printf("🚀 Health Tracker API starting on port %s", config.AppConfig.Port)
	log.Printf("📚 API endpoints available at http://localhost:%s/api", config.AppConfig.Port)
	log.Printf("❤️  Health check at http://localhost:%s/health", config.AppConfig.Port)

	if err := r.Run(":" + config.AppConfig.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

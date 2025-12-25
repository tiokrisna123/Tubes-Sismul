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

	// --- TAMBAHAN: PENGATURAN CORS (Supaya Netlify tidak diblokir) ---
	r.Use(func(c *gin.Context) {
		// Izinkan semua domain mengakses (*)
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		// Izinkan header-header penting (Authorization untuk login token, dll)
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		// Izinkan method standar
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		// Handle preflight request (OPTIONS)
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})
	// ------------------------------------------------------------------

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

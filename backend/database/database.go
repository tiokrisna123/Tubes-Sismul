package database

import (
	"health-tracker/config"
	"health-tracker/models"
	"log"

	"gorm.io/driver/postgres" // IMPORT BARU: Driver Postgres
	"gorm.io/driver/sqlite"   // Driver SQLite tetap disimpan untuk backup
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDatabase() {
	var err error
	var dialector gorm.Dialector

	// LOGIKA BARU: Cek apakah kita di Render (punya DB_URL) atau di Laptop
	if config.AppConfig.DatabaseURL != "" {
		// Jika ada DB_URL, gunakan PostgreSQL (Untuk Render/Neon)
		log.Println("🌍 Detected DB_URL, connecting to PostgreSQL...")
		dialector = postgres.Open(config.AppConfig.DatabaseURL)
	} else {
		// Jika tidak ada, gunakan SQLite (Untuk Localhost)
		log.Println("💻 No DB_URL found, using local SQLite...")
		dialector = sqlite.Open(config.AppConfig.DatabasePath)
	}

	DB, err = gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Database connected successfully")

	// Auto-migrate models
	err = DB.AutoMigrate(
		&models.User{},
		&models.HealthData{},
		&models.Symptom{},
		&models.SymptomTemplate{},
		&models.FamilyMember{},
		&models.Recommendation{},
		&models.Article{},
		&models.Post{},
		&models.Comment{},
		&models.Like{},
		&models.WaterIntake{},
		&models.Goal{},
		&models.Reminder{},
	)

	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	log.Println("Database migration completed")

	// Seed initial data
	SeedData()
}

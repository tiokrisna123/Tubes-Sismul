package database

import (
	"health-tracker/config"
	"health-tracker/models"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDatabase() {
	var err error
	
	DB, err = gorm.Open(sqlite.Open(config.AppConfig.DatabasePath), &gorm.Config{
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

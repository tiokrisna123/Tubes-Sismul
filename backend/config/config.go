package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port           string
	GinMode        string
	JWTSecret      string
	JWTExpiryHours int
	DatabasePath   string // Untuk SQLite (Local)
	DatabaseURL    string // Untuk PostgreSQL (Render/Neon)
}

var AppConfig *Config

func LoadConfig() {
	godotenv.Load()

	expiryHours, _ := strconv.Atoi(getEnv("JWT_EXPIRY_HOURS", "24"))

	AppConfig = &Config{
		Port:           getEnv("PORT", "8080"),
		GinMode:        getEnv("GIN_MODE", "debug"),
		JWTSecret:      getEnv("JWT_SECRET", "default-secret-key"),
		JWTExpiryHours: expiryHours,
		DatabasePath:   getEnv("DATABASE_PATH", "./health_tracker.db"),
		// INI YANG BARU: Membaca Environment Variable DB_URL dari Render
		DatabaseURL: getEnv("DB_URL", ""),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

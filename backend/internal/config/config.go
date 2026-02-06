package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	DatabaseURL string
}

func Load() *Config {
	_ = godotenv.Load() // Load .env if it exists, ignore if not

	return &Config{
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://gable_user:gable_password@localhost:5434/gable_db?sslmode=disable"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

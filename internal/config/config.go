package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL   string
	SessionSecret string
	Port          string
	GinMode       string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	return &Config{
		DatabaseURL:   getEnv("DATABASE_URL", "postgres://taskflow:taskflow@localhost:5432/taskflow?sslmode=disable"),
		SessionSecret: getEnv("SESSION_SECRET", "dev-secret"),
		Port:          getEnv("PORT", "3000"),
		GinMode:       getEnv("GIN_MODE", "debug"),
	}, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

package config

import (
	"os"
	"strconv"
)

type Config struct {
	Env           string
	Port          string
	DatabaseURL   string
	JWTSecret     string
	TelegramToken string
	SessionSecret string
	SessionTTL    int
	FrontendDist  string
}

func Load() *Config {
	return &Config{
		Env:           getEnv("ENV", "development"),
		Port:          getEnv("PORT", "8080"),
		DatabaseURL:   getEnv("DATABASE_URL", "postgres://kanban:kanban@localhost:5433/kanban?sslmode=disable"),
		JWTSecret:     getEnv("JWT_SECRET", "dev-secret-change-in-production"),
		TelegramToken: getEnv("TELEGRAM_BOT_TOKEN", ""),
		SessionSecret: getEnv("SESSION_SECRET", "dev-session-secret"),
		SessionTTL:    getEnvInt("SESSION_TTL", 604800), // 7 days
		FrontendDist:  getEnv("FRONTEND_DIST", "dist"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
	}
	return defaultValue
}

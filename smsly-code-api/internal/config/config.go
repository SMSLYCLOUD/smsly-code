package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
	Port       string
}

func LoadConfig() Config {
	// Try to load from .env file, but don't fail if it doesn't exist (e.g. Docker)
	_ = godotenv.Load()

	return Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "smsly"),
		DBPassword: getEnv("DB_PASSWORD", "smsly_secret"),
		DBName:     getEnv("DB_NAME", "smsly_code"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),
		Port:       getEnv("API_PORT", "8080"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

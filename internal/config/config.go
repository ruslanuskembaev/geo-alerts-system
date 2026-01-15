package config

import (
	"os"
)

type Config struct {
	ServerPort string
	APIKey     string
}

func Load() *Config {
	return &Config{
		ServerPort: getEnv("SERVER_PORT", "8080"),
		APIKey:     getEnv("API_KEY", "default_api_key"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

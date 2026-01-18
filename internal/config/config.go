package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	ServerPort string
	APIKey     string

	// Database
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	// Redis
	RedisHost     string
	RedisPort     string
	RedisPassword string
	RedisDB       int

	// Webhook
	WebhookURL           string
	WebhookRetryAttempts int
	WebhookRetryDelay    time.Duration

	// Stats
	StatsTimeWindow time.Duration

	// Cache
	CacheTTL time.Duration
}

func Load() *Config {
	return &Config{
		ServerPort: getEnv("SERVER_PORT", "8080"),
		APIKey:     getEnv("API_KEY", "dev_api_key_12345"),

		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "geoalerts"),
		DBPassword: getEnv("DB_PASSWORD", "password"),
		DBName:     getEnv("DB_NAME", "geoalerts_db"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),

		RedisHost:     getEnv("REDIS_HOST", "localhost"),
		RedisPort:     getEnv("REDIS_PORT", "6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       getEnvAsInt("REDIS_DB", 0),

		WebhookURL:           getEnv("WEBHOOK_URL", "http://localhost:9090/webhook"),
		WebhookRetryAttempts: getEnvAsInt("WEBHOOK_RETRY_ATTEMPTS", 3),
		WebhookRetryDelay:    time.Duration(getEnvAsInt("WEBHOOK_RETRY_DELAY_SECONDS", 5)) * time.Second,

		StatsTimeWindow: time.Duration(getEnvAsInt("STATS_TIME_WINDOW_MINUTES", 60)) * time.Minute,

		CacheTTL: time.Duration(getEnvAsInt("CACHE_TTL_SECONDS", 300)) * time.Second,
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

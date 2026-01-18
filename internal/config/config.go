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
	DBMaxConns int
	DBMinConns int

	// DB pool timeouts
	DBMaxConnLifetime time.Duration
	DBMaxConnIdleTime time.Duration

	// Redis
	RedisHost     string
	RedisPort     string
	RedisPassword string
	RedisDB       int

	// Webhook
	WebhookURL           string
	WebhookRetryAttempts int
	WebhookRetryDelay    time.Duration
	WebhookTimeout       time.Duration

	// Stats
	StatsTimeWindow time.Duration

	// Cache
	CacheTTL time.Duration

	// HTTP server
	HTTPReadTimeout  time.Duration
	HTTPWriteTimeout time.Duration
	HTTPIdleTimeout  time.Duration
	ShutdownTimeout  time.Duration

	// Health
	HealthTimeout time.Duration
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
		DBMaxConns: getEnvAsInt("DB_MAX_CONNS", 10),
		DBMinConns: getEnvAsInt("DB_MIN_CONNS", 2),

		DBMaxConnLifetime: getEnvAsDuration("DB_MAX_CONN_LIFETIME_SECONDS", 1800),
		DBMaxConnIdleTime: getEnvAsDuration("DB_MAX_CONN_IDLE_SECONDS", 600),

		RedisHost:     getEnv("REDIS_HOST", "localhost"),
		RedisPort:     getEnv("REDIS_PORT", "6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       getEnvAsInt("REDIS_DB", 0),

		WebhookURL:           getEnv("WEBHOOK_URL", "http://localhost:9090/webhook"),
		WebhookRetryAttempts: getEnvAsInt("WEBHOOK_RETRY_ATTEMPTS", 3),
		WebhookRetryDelay:    getEnvAsDuration("WEBHOOK_RETRY_DELAY_SECONDS", 5),
		WebhookTimeout:       getEnvAsDuration("WEBHOOK_TIMEOUT_SECONDS", 5),

		StatsTimeWindow: time.Duration(getEnvAsInt("STATS_TIME_WINDOW_MINUTES", 60)) * time.Minute,

		CacheTTL: getEnvAsDuration("CACHE_TTL_SECONDS", 300),

		HTTPReadTimeout:  getEnvAsDuration("HTTP_READ_TIMEOUT_SECONDS", 5),
		HTTPWriteTimeout: getEnvAsDuration("HTTP_WRITE_TIMEOUT_SECONDS", 10),
		HTTPIdleTimeout:  getEnvAsDuration("HTTP_IDLE_TIMEOUT_SECONDS", 60),
		ShutdownTimeout:  getEnvAsDuration("SHUTDOWN_TIMEOUT_SECONDS", 10),

		HealthTimeout: getEnvAsDuration("HEALTH_TIMEOUT_SECONDS", 2),
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

func getEnvAsDuration(key string, defaultSeconds int) time.Duration {
	return time.Duration(getEnvAsInt(key, defaultSeconds)) * time.Second
}

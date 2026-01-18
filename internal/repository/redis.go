package repository

import (
	"fmt"

	"github.com/redis/go-redis/v9"

	"github.com/ruslanuskembaev/geo-alerts-system/internal/config"
)

// NewRedisClient создаёт клиент Redis
func NewRedisClient(cfg *config.Config) *redis.Client {
	addr := fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort)
	return redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})
}

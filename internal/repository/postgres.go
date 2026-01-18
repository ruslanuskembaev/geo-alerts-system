package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ruslanuskembaev/geo-alerts-system/internal/config"
)

// NewPostgresPool создаёт пул подключений к PostgreSQL
func NewPostgresPool(cfg *config.Config) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
		cfg.DBSSLMode,
	)

	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}
	if cfg.DBMaxConns > 0 {
		poolConfig.MaxConns = int32(cfg.DBMaxConns)
	}
	if cfg.DBMinConns > 0 {
		poolConfig.MinConns = int32(cfg.DBMinConns)
	}
	if cfg.DBMaxConnLifetime > 0 {
		poolConfig.MaxConnLifetime = cfg.DBMaxConnLifetime
	}
	if cfg.DBMaxConnIdleTime > 0 {
		poolConfig.MaxConnIdleTime = cfg.DBMaxConnIdleTime
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}

	return pool, nil
}

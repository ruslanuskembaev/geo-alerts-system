package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

// HealthRepository defines dependency health checks.
type HealthRepository interface {
	PingDB(ctx context.Context) error
	PingRedis(ctx context.Context) error
}

// SystemRepository implements health checks for DB and Redis.
type SystemRepository struct {
	db    *pgxpool.Pool
	redis *redis.Client
}

func NewSystemRepository(db *pgxpool.Pool, redis *redis.Client) *SystemRepository {
	return &SystemRepository{
		db:    db,
		redis: redis,
	}
}

func (r *SystemRepository) PingDB(ctx context.Context) error {
	return r.db.Ping(ctx)
}

func (r *SystemRepository) PingRedis(ctx context.Context) error {
	return r.redis.Ping(ctx).Err()
}

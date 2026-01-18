package repository

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/ruslanuskembaev/geo-alerts-system/internal/domain"
)

const activeIncidentsCacheKey = "geoalerts:active_incidents"

// IncidentCache defines active incidents cache behavior.
type IncidentCache interface {
	GetActive(ctx context.Context) ([]*domain.Incident, bool, error)
	SetActive(ctx context.Context, incidents []*domain.Incident) error
	Invalidate(ctx context.Context) error
}

// RedisIncidentCache implements IncidentCache using Redis.
type RedisIncidentCache struct {
	client *redis.Client
	ttl    time.Duration
}

func NewIncidentCache(client *redis.Client, ttl time.Duration) *RedisIncidentCache {
	return &RedisIncidentCache{
		client: client,
		ttl:    ttl,
	}
}

func (c *RedisIncidentCache) GetActive(ctx context.Context) ([]*domain.Incident, bool, error) {
	raw, err := c.client.Get(ctx, activeIncidentsCacheKey).Result()
	if err == redis.Nil {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}

	var incidents []*domain.Incident
	if err := json.Unmarshal([]byte(raw), &incidents); err != nil {
		return nil, false, err
	}

	return incidents, true, nil
}

func (c *RedisIncidentCache) SetActive(ctx context.Context, incidents []*domain.Incident) error {
	raw, err := json.Marshal(incidents)
	if err != nil {
		return err
	}

	return c.client.Set(ctx, activeIncidentsCacheKey, raw, c.ttl).Err()
}

func (c *RedisIncidentCache) Invalidate(ctx context.Context) error {
	return c.client.Del(ctx, activeIncidentsCacheKey).Err()
}

package repository

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/ruslanuskembaev/geo-alerts-system/internal/domain"
)

const webhookQueueKey = "geoalerts:webhook_queue"

// WebhookQueue defines enqueue/dequeue operations for webhook jobs.
type WebhookQueue interface {
	Enqueue(ctx context.Context, job domain.WebhookJob) error
	Dequeue(ctx context.Context, timeout time.Duration) (*domain.WebhookJob, bool, error)
}

// RedisWebhookQueue implements WebhookQueue using Redis lists.
type RedisWebhookQueue struct {
	client *redis.Client
}

func NewWebhookQueue(client *redis.Client) *RedisWebhookQueue {
	return &RedisWebhookQueue{client: client}
}

func (q *RedisWebhookQueue) Enqueue(ctx context.Context, job domain.WebhookJob) error {
	raw, err := json.Marshal(job)
	if err != nil {
		return err
	}
	return q.client.LPush(ctx, webhookQueueKey, raw).Err()
}

func (q *RedisWebhookQueue) Dequeue(ctx context.Context, timeout time.Duration) (*domain.WebhookJob, bool, error) {
	result, err := q.client.BRPop(ctx, timeout, webhookQueueKey).Result()
	if err == redis.Nil {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}
	if len(result) < 2 {
		return nil, false, nil
	}

	var job domain.WebhookJob
	if err := json.Unmarshal([]byte(result[1]), &job); err != nil {
		return nil, false, err
	}

	return &job, true, nil
}

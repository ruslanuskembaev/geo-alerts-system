//go:build integration

package integration

import (
	"context"
	"testing"
	"time"

	"github.com/ruslanuskembaev/geo-alerts-system/internal/domain"
	"github.com/ruslanuskembaev/geo-alerts-system/internal/repository"
)

func TestRedisQueueAndCache(t *testing.T) {
	client := testRedis(t)
	defer func() {
		_ = client.Close()
	}()

	cache := repository.NewIncidentCache(client, time.Minute)
	queue := repository.NewWebhookQueue(client)

	ctx := context.Background()

	incidents := []*domain.Incident{
		{
			ID:           "incident-1",
			Title:        "Test",
			Severity:     domain.SeverityLow,
			Latitude:     1,
			Longitude:    2,
			RadiusMeters: 100,
			IsActive:     true,
			CreatedAt:    time.Now().UTC(),
			UpdatedAt:    time.Now().UTC(),
		},
	}

	if err := cache.SetActive(ctx, incidents); err != nil {
		t.Fatalf("cache set failed: %v", err)
	}
	got, ok, err := cache.GetActive(ctx)
	if err != nil {
		t.Fatalf("cache get failed: %v", err)
	}
	if !ok || len(got) != 1 {
		t.Fatalf("expected cache hit with one incident")
	}

	job := domain.WebhookJob{
		Payload: domain.WebhookPayload{
			CheckID:        "check-1",
			UserID:         "user-1",
			Latitude:       1,
			Longitude:      2,
			IsInDangerZone: true,
			CheckedAt:      time.Now().UTC(),
			Incidents: []domain.NearbyIncident{
				{ID: "incident-1", Title: "Test", Severity: domain.SeverityLow},
			},
		},
		Attempt:   0,
		CreatedAt: time.Now().UTC(),
	}

	if err := queue.Enqueue(ctx, job); err != nil {
		t.Fatalf("enqueue failed: %v", err)
	}

	dequeued, ok, err := queue.Dequeue(ctx, 2*time.Second)
	if err != nil {
		t.Fatalf("dequeue failed: %v", err)
	}
	if !ok || dequeued == nil {
		t.Fatalf("expected a dequeued job")
	}
	if dequeued.Payload.CheckID != job.Payload.CheckID {
		t.Fatalf("unexpected dequeued job payload")
	}
}

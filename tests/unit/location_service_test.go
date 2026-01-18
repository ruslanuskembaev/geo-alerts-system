package unit

import (
	"context"
	"testing"
	"time"

	"github.com/ruslanuskembaev/geo-alerts-system/internal/domain"
	svc "github.com/ruslanuskembaev/geo-alerts-system/internal/service"
)

func TestLocationService_CheckLocation_CacheHit_EnqueuesWebhook(t *testing.T) {
	incidents := []*domain.Incident{
		{
			ID:           "incident-1",
			Title:        "Near",
			Severity:     domain.SeverityHigh,
			Latitude:     0,
			Longitude:    0,
			RadiusMeters: 1000,
			IsActive:     true,
		},
		{
			ID:           "incident-2",
			Title:        "Farther",
			Severity:     domain.SeverityMedium,
			Latitude:     0,
			Longitude:    0.01,
			RadiusMeters: 2000,
			IsActive:     true,
		},
	}

	repo := &fakeIncidentRepo{
		listActiveFn: func(ctx context.Context) ([]*domain.Incident, error) {
			return nil, nil
		},
	}
	cache := &fakeIncidentCache{
		getFn: func(ctx context.Context) ([]*domain.Incident, bool, error) {
			return incidents, true, nil
		},
	}
	checkRepo := &fakeCheckRepo{}
	queue := &fakeQueue{}

	service := svc.NewLocationService(repo, cache, checkRepo, queue)

	resp, err := service.CheckLocation(context.Background(), domain.LocationCheckRequest{
		UserID:    "user-1",
		Latitude:  0,
		Longitude: 0,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !resp.IsInDangerZone {
		t.Fatalf("expected IsInDangerZone=true")
	}
	if len(resp.Incidents) != 2 {
		t.Fatalf("expected 2 incidents, got %d", len(resp.Incidents))
	}
	if resp.Incidents[0].ID != "incident-1" {
		t.Fatalf("expected nearest incident first")
	}
	if resp.Incidents[0].DistanceMeters > resp.Incidents[1].DistanceMeters {
		t.Fatalf("expected incidents sorted by distance")
	}
	if repo.listActiveCalls != 0 {
		t.Fatalf("expected cache hit to skip ListActive")
	}
	if cache.setCalls != 0 {
		t.Fatalf("expected cache SetActive not called on hit")
	}
	if checkRepo.createCalls != 1 {
		t.Fatalf("expected location check to be stored")
	}
	if len(checkRepo.lastIncidentIDs) != 2 {
		t.Fatalf("expected incident IDs stored for check")
	}
	if len(queue.enqueued) != 1 {
		t.Fatalf("expected webhook job enqueued")
	}
	job := queue.enqueued[0]
	if job.Payload.CheckID != resp.CheckID {
		t.Fatalf("expected webhook payload to include check ID")
	}
	if len(job.Payload.Incidents) != 2 {
		t.Fatalf("expected webhook payload incidents")
	}
	if job.CreatedAt.IsZero() {
		t.Fatalf("expected job CreatedAt to be set")
	}
}

func TestLocationService_CheckLocation_CacheMiss_NoMatches(t *testing.T) {
	incidents := []*domain.Incident{
		{
			ID:           "incident-1",
			Title:        "Far",
			Severity:     domain.SeverityLow,
			Latitude:     10,
			Longitude:    10,
			RadiusMeters: 100,
			IsActive:     true,
		},
	}

	repo := &fakeIncidentRepo{
		listActiveFn: func(ctx context.Context) ([]*domain.Incident, error) {
			return incidents, nil
		},
	}
	cache := &fakeIncidentCache{
		getFn: func(ctx context.Context) ([]*domain.Incident, bool, error) {
			return nil, false, nil
		},
		setFn: func(ctx context.Context, incidents []*domain.Incident) error {
			if len(incidents) != 1 {
				t.Fatalf("expected cache to store incidents")
			}
			return nil
		},
	}
	checkRepo := &fakeCheckRepo{}
	queue := &fakeQueue{
		enqueueFn: func(ctx context.Context, job domain.WebhookJob) error {
			t.Fatalf("expected no webhook enqueue when no matches")
			return nil
		},
	}

	service := svc.NewLocationService(repo, cache, checkRepo, queue)

	resp, err := service.CheckLocation(context.Background(), domain.LocationCheckRequest{
		UserID:    "user-2",
		Latitude:  0,
		Longitude: 0,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.IsInDangerZone {
		t.Fatalf("expected IsInDangerZone=false")
	}
	if len(resp.Incidents) != 0 {
		t.Fatalf("expected no incidents")
	}
	if repo.listActiveCalls != 1 {
		t.Fatalf("expected ListActive called on cache miss")
	}
	if cache.setCalls != 1 {
		t.Fatalf("expected cache SetActive called on miss")
	}
	if checkRepo.createCalls != 1 {
		t.Fatalf("expected location check to be stored")
	}
	if len(queue.enqueued) != 0 {
		t.Fatalf("expected no webhook jobs enqueued")
	}
	if checkRepo.lastCheck.CheckedAt.After(time.Now().UTC().Add(1 * time.Second)) {
		t.Fatalf("unexpected check timestamp")
	}
}

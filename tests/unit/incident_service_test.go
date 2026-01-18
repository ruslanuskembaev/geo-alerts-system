package unit

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/ruslanuskembaev/geo-alerts-system/internal/domain"
	svc "github.com/ruslanuskembaev/geo-alerts-system/internal/service"
)

func TestIncidentService_InvalidatesCache(t *testing.T) {
	repo := &fakeIncidentRepo{
		createFn: func(ctx context.Context, req domain.CreateIncidentRequest) (*domain.Incident, error) {
			return &domain.Incident{ID: "incident-1"}, nil
		},
		updateFn: func(ctx context.Context, id string, req domain.UpdateIncidentRequest) (*domain.Incident, error) {
			return &domain.Incident{ID: id}, nil
		},
		deactivateFn: func(ctx context.Context, id string) error {
			return nil
		},
	}
	cache := &fakeIncidentCache{}
	checkRepo := &fakeCheckRepo{}

	service := svc.NewIncidentService(repo, cache, checkRepo)

	if _, err := service.Create(context.Background(), domain.CreateIncidentRequest{
		Title:        "Test",
		Severity:     domain.SeverityLow,
		Latitude:     0,
		Longitude:    0,
		RadiusMeters: 100,
	}); err != nil {
		t.Fatalf("unexpected create error: %v", err)
	}

	title := "Updated"
	if _, err := service.Update(context.Background(), "incident-1", domain.UpdateIncidentRequest{
		Title: &title,
	}); err != nil {
		t.Fatalf("unexpected update error: %v", err)
	}

	if err := service.Deactivate(context.Background(), "incident-1"); err != nil {
		t.Fatalf("unexpected deactivate error: %v", err)
	}

	if cache.invalidateCalls != 3 {
		t.Fatalf("expected cache invalidation on create/update/deactivate, got %d", cache.invalidateCalls)
	}
}

func TestIncidentService_StatsByIncident(t *testing.T) {
	expected := []domain.IncidentStats{
		{IncidentID: "incident-1", Title: "Test", UserCount: 2},
	}

	checkRepo := &fakeCheckRepo{
		statsFn: func(ctx context.Context, since time.Time) ([]domain.IncidentStats, error) {
			return expected, nil
		},
	}

	service := svc.NewIncidentService(&fakeIncidentRepo{}, &fakeIncidentCache{}, checkRepo)

	stats, err := service.StatsByIncident(context.Background(), time.Now().Add(-time.Hour))
	if err != nil {
		t.Fatalf("unexpected stats error: %v", err)
	}
	if !reflect.DeepEqual(stats, expected) {
		t.Fatalf("unexpected stats result")
	}
}

//go:build integration

package integration

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/ruslanuskembaev/geo-alerts-system/internal/domain"
	"github.com/ruslanuskembaev/geo-alerts-system/internal/repository"
)

func TestLocationCheckRepository_Stats(t *testing.T) {
	pool := testDB(t)
	defer pool.Close()
	truncateTables(t, pool)

	incidentRepo := repository.NewIncidentRepository(pool)
	checkRepo := repository.NewLocationCheckRepository(pool)

	incident, err := incidentRepo.Create(context.Background(), domain.CreateIncidentRequest{
		Title:        "Test Incident",
		Description:  "Test",
		Severity:     domain.SeverityLow,
		Latitude:     10,
		Longitude:    10,
		RadiusMeters: 1000,
	})
	if err != nil {
		t.Fatalf("create incident failed: %v", err)
	}

	check := domain.LocationCheck{
		ID:             uuid.New().String(),
		UserID:         "user-1",
		Latitude:       10,
		Longitude:      10,
		IsInDangerZone: true,
		CheckedAt:      time.Now().UTC(),
	}
	if err := checkRepo.Create(context.Background(), check, []string{incident.ID}); err != nil {
		t.Fatalf("create check failed: %v", err)
	}

	another := domain.LocationCheck{
		ID:             uuid.New().String(),
		UserID:         "user-1",
		Latitude:       10,
		Longitude:      10,
		IsInDangerZone: true,
		CheckedAt:      time.Now().UTC(),
	}
	if err := checkRepo.Create(context.Background(), another, []string{incident.ID}); err != nil {
		t.Fatalf("create second check failed: %v", err)
	}

	stats, err := checkRepo.StatsByIncident(context.Background(), time.Now().Add(-time.Hour))
	if err != nil {
		t.Fatalf("stats failed: %v", err)
	}
	if len(stats) != 1 {
		t.Fatalf("expected one incident stat")
	}
	if stats[0].UserCount != 1 {
		t.Fatalf("expected unique user_count=1")
	}
}

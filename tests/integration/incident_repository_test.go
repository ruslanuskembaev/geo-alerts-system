//go:build integration

package integration

import (
	"context"
	"testing"
	"time"

	"github.com/ruslanuskembaev/geo-alerts-system/internal/domain"
	"github.com/ruslanuskembaev/geo-alerts-system/internal/repository"
)

func TestIncidentRepository_CRUD(t *testing.T) {
	pool := testDB(t)
	defer pool.Close()
	truncateTables(t, pool)

	repo := repository.NewIncidentRepository(pool)

	created, err := repo.Create(context.Background(), domain.CreateIncidentRequest{
		Title:        "Test Incident",
		Description:  "Smoke near park",
		Severity:     domain.SeverityHigh,
		Latitude:     55.7,
		Longitude:    37.6,
		RadiusMeters: 500,
	})
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}

	got, err := repo.GetByID(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("get failed: %v", err)
	}
	if got.Title != created.Title {
		t.Fatalf("unexpected title")
	}

	list, total, err := repo.List(context.Background(), 10, 0)
	if err != nil {
		t.Fatalf("list failed: %v", err)
	}
	if total != 1 || len(list) != 1 {
		t.Fatalf("expected total=1 list=1")
	}

	newTitle := "Updated"
	newRadius := 750
	updated, err := repo.Update(context.Background(), created.ID, domain.UpdateIncidentRequest{
		Title:        &newTitle,
		RadiusMeters: &newRadius,
	})
	if err != nil {
		t.Fatalf("update failed: %v", err)
	}
	if updated.Title != newTitle || updated.RadiusMeters != newRadius {
		t.Fatalf("update values not applied")
	}

	if err := repo.Deactivate(context.Background(), created.ID); err != nil {
		t.Fatalf("deactivate failed: %v", err)
	}

	active, err := repo.ListActive(context.Background())
	if err != nil {
		t.Fatalf("list active failed: %v", err)
	}
	if len(active) != 0 {
		t.Fatalf("expected no active incidents after deactivation")
	}

	found, err := repo.GetByID(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("get after deactivate failed: %v", err)
	}
	if found.IsActive {
		t.Fatalf("expected incident to be inactive")
	}
	if time.Since(found.UpdatedAt) > time.Minute {
		t.Fatalf("expected updated_at to be recent")
	}
}

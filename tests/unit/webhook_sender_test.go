package unit

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ruslanuskembaev/geo-alerts-system/internal/domain"
	svc "github.com/ruslanuskembaev/geo-alerts-system/internal/service"
)

func TestWebhookSender_Send_OK(t *testing.T) {
	expected := domain.WebhookPayload{
		CheckID:        "check-1",
		UserID:         "user-1",
		Latitude:       1.2,
		Longitude:      3.4,
		IsInDangerZone: true,
		CheckedAt:      time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC),
		Incidents: []domain.NearbyIncident{
			{ID: "incident-1", Title: "Test", Severity: domain.SeverityHigh},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Fatalf("expected application/json content type")
		}

		var received domain.WebhookPayload
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Fatalf("failed to decode payload: %v", err)
		}

		if received.CheckID != expected.CheckID || received.UserID != expected.UserID {
			t.Fatalf("unexpected payload identifiers")
		}
		if !received.CheckedAt.Equal(expected.CheckedAt) {
			t.Fatalf("unexpected payload timestamp")
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	sender := svc.NewWebhookSender(server.URL, 2*time.Second)
	if err := sender.Send(context.Background(), expected); err != nil {
		t.Fatalf("unexpected send error: %v", err)
	}
}

func TestWebhookSender_Send_ErrorStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	sender := svc.NewWebhookSender(server.URL, 2*time.Second)
	if err := sender.Send(context.Background(), domain.WebhookPayload{}); err == nil {
		t.Fatalf("expected error on non-2xx response")
	}
}

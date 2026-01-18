package service

import (
	"context"
	"math"
	"sort"
	"time"

	"github.com/google/uuid"

	"github.com/ruslanuskembaev/geo-alerts-system/internal/domain"
	"github.com/ruslanuskembaev/geo-alerts-system/internal/repository"
)

// LocationService сервис проверки координат
type LocationService struct {
	incidentRepo repository.IncidentRepository
	cache        repository.IncidentCache
	checkRepo    repository.LocationCheckRepository
	queue        repository.WebhookQueue
}

func NewLocationService(
	incidentRepo repository.IncidentRepository,
	cache repository.IncidentCache,
	checkRepo repository.LocationCheckRepository,
	queue repository.WebhookQueue,
) *LocationService {
	return &LocationService{
		incidentRepo: incidentRepo,
		cache:        cache,
		checkRepo:    checkRepo,
		queue:        queue,
	}
}

func (s *LocationService) CheckLocation(ctx context.Context, req domain.LocationCheckRequest) (*domain.LocationCheckResponse, error) {
	incidents, ok, err := s.cache.GetActive(ctx)
	if err != nil {
		return nil, err
	}
	if !ok {
		incidents, err = s.incidentRepo.ListActive(ctx)
		if err != nil {
			return nil, err
		}
		_ = s.cache.SetActive(ctx, incidents)
	}

	matched := make([]domain.NearbyIncident, 0)
	incidentIDs := make([]string, 0)
	for _, incident := range incidents {
		distance := distanceMeters(req.Latitude, req.Longitude, incident.Latitude, incident.Longitude)
		if distance <= float64(incident.RadiusMeters) {
			matched = append(matched, domain.NearbyIncident{
				ID:             incident.ID,
				Title:          incident.Title,
				Severity:       incident.Severity,
				Latitude:       incident.Latitude,
				Longitude:      incident.Longitude,
				RadiusMeters:   incident.RadiusMeters,
				DistanceMeters: distance,
			})
			incidentIDs = append(incidentIDs, incident.ID)
		}
	}

	sort.Slice(matched, func(i, j int) bool {
		return matched[i].DistanceMeters < matched[j].DistanceMeters
	})

	now := time.Now().UTC()
	check := domain.LocationCheck{
		ID:             uuid.New().String(),
		UserID:         req.UserID,
		Latitude:       req.Latitude,
		Longitude:      req.Longitude,
		IsInDangerZone: len(matched) > 0,
		CheckedAt:      now,
	}

	if err := s.checkRepo.Create(ctx, check, incidentIDs); err != nil {
		return nil, err
	}

	if len(matched) > 0 {
		job := domain.WebhookJob{
			Payload: domain.WebhookPayload{
				CheckID:        check.ID,
				UserID:         check.UserID,
				Latitude:       check.Latitude,
				Longitude:      check.Longitude,
				IsInDangerZone: check.IsInDangerZone,
				CheckedAt:      check.CheckedAt,
				Incidents:      matched,
			},
			Attempt:   0,
			CreatedAt: now,
		}

		if err := s.queue.Enqueue(ctx, job); err != nil {
			return nil, err
		}
	}

	return &domain.LocationCheckResponse{
		CheckID:        check.ID,
		IsInDangerZone: check.IsInDangerZone,
		CheckedAt:      check.CheckedAt,
		Incidents:      matched,
	}, nil
}

func distanceMeters(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadius = 6371000

	lat1Rad := lat1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	deltaLat := (lat2 - lat1) * math.Pi / 180
	deltaLon := (lon2 - lon1) * math.Pi / 180

	sinLat := math.Sin(deltaLat / 2)
	sinLon := math.Sin(deltaLon / 2)
	a := sinLat*sinLat + math.Cos(lat1Rad)*math.Cos(lat2Rad)*sinLon*sinLon
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadius * c
}

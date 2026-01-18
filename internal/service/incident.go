package service

import (
	"context"
	"time"

	"github.com/ruslanuskembaev/geo-alerts-system/internal/domain"
	"github.com/ruslanuskembaev/geo-alerts-system/internal/repository"
)

// IncidentService сервис для работы с инцидентами
type IncidentService struct {
	repo      repository.IncidentRepository
	cache     repository.IncidentCache
	checkRepo repository.LocationCheckRepository
}

func NewIncidentService(
	repo repository.IncidentRepository,
	cache repository.IncidentCache,
	checkRepo repository.LocationCheckRepository,
) *IncidentService {
	return &IncidentService{
		repo:      repo,
		cache:     cache,
		checkRepo: checkRepo,
	}
}

func (s *IncidentService) Create(ctx context.Context, req domain.CreateIncidentRequest) (*domain.Incident, error) {
	incident, err := s.repo.Create(ctx, req)
	if err != nil {
		return nil, err
	}
	_ = s.cache.Invalidate(ctx)
	return incident, nil
}

func (s *IncidentService) GetByID(ctx context.Context, id string) (*domain.Incident, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *IncidentService) List(ctx context.Context, limit, offset int) ([]*domain.Incident, int, error) {
	return s.repo.List(ctx, limit, offset)
}

func (s *IncidentService) Update(ctx context.Context, id string, req domain.UpdateIncidentRequest) (*domain.Incident, error) {
	incident, err := s.repo.Update(ctx, id, req)
	if err != nil {
		return nil, err
	}
	_ = s.cache.Invalidate(ctx)
	return incident, nil
}

func (s *IncidentService) Deactivate(ctx context.Context, id string) error {
	if err := s.repo.Deactivate(ctx, id); err != nil {
		return err
	}
	_ = s.cache.Invalidate(ctx)
	return nil
}

func (s *IncidentService) StatsByIncident(ctx context.Context, since time.Time) ([]domain.IncidentStats, error) {
	return s.checkRepo.StatsByIncident(ctx, since)
}

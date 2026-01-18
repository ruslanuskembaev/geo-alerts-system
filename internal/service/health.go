package service

import (
	"context"
	"time"

	"github.com/ruslanuskembaev/geo-alerts-system/internal/domain"
	"github.com/ruslanuskembaev/geo-alerts-system/internal/repository"
)

// HealthService checks system dependencies.
type HealthService struct {
	repo    repository.HealthRepository
	timeout time.Duration
}

func NewHealthService(repo repository.HealthRepository, timeout time.Duration) *HealthService {
	return &HealthService{
		repo:    repo,
		timeout: timeout,
	}
}

func (s *HealthService) Check(ctx context.Context) domain.SystemHealth {
	status := domain.SystemHealth{
		Status:       "healthy",
		Service:      "geo-alerts-system",
		Timestamp:    time.Now().UTC(),
		Dependencies: map[string]domain.DependencyHealth{},
	}

	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	if err := s.repo.PingDB(ctx); err != nil {
		status.Status = "degraded"
		status.Dependencies["postgres"] = domain.DependencyHealth{
			Status: "error",
			Error:  err.Error(),
		}
	} else {
		status.Dependencies["postgres"] = domain.DependencyHealth{Status: "ok"}
	}

	if err := s.repo.PingRedis(ctx); err != nil {
		status.Status = "degraded"
		status.Dependencies["redis"] = domain.DependencyHealth{
			Status: "error",
			Error:  err.Error(),
		}
	} else {
		status.Dependencies["redis"] = domain.DependencyHealth{Status: "ok"}
	}

	return status
}

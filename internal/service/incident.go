package service

import (
	"github.com/ruslanuskembaev/geo-alerts-system/internal/domain"
	"github.com/ruslanuskembaev/geo-alerts-system/internal/repository"
)

// IncidentService сервис для работы с инцидентами
type IncidentService struct {
	repo *repository.IncidentRepository
}

func NewIncidentService(repo *repository.IncidentRepository) *IncidentService {
	return &IncidentService{repo: repo}
}

func (s *IncidentService) Create(req domain.CreateIncidentRequest) (*domain.Incident, error) {
	return s.repo.Create(req)
}

func (s *IncidentService) GetByID(id string) (*domain.Incident, error) {
	return s.repo.GetByID(id)
}

func (s *IncidentService) List() ([]*domain.Incident, error) {
	return s.repo.List()
}

func (s *IncidentService) Update(id string, req domain.UpdateIncidentRequest) (*domain.Incident, error) {
	return s.repo.Update(id, req)
}

func (s *IncidentService) Delete(id string) error {
	return s.repo.Delete(id)
}

func (s *IncidentService) GetActiveIncidents() ([]*domain.Incident, error) {
	return s.repo.GetActiveIncidents()
}

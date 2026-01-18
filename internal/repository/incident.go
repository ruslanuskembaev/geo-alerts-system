package repository

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/ruslanuskembaev/geo-alerts-system/internal/domain"
)

// IncidentRepository репозиторий для работы с инцидентами
type IncidentRepository struct {
	mu        sync.RWMutex
	incidents map[string]*domain.Incident
}

func NewIncidentRepository() *IncidentRepository {
	return &IncidentRepository{
		incidents: make(map[string]*domain.Incident),
	}
}

// Create создаёт новый инцидент
func (r *IncidentRepository) Create(req domain.CreateIncidentRequest) (*domain.Incident, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	incident := &domain.Incident{
		ID:           uuid.New().String(),
		Title:        req.Title,
		Description:  req.Description,
		Severity:     req.Severity,
		Latitude:     req.Latitude,
		Longitude:    req.Longitude,
		RadiusMeters: req.RadiusMeters,
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	r.incidents[incident.ID] = incident
	return incident, nil
}

// GetByID получает инцидент по ID
func (r *IncidentRepository) GetByID(id string) (*domain.Incident, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	incident, exists := r.incidents[id]
	if !exists {
		return nil, fmt.Errorf("incident not found")
	}

	return incident, nil
}

// List возвращает список всех активных инцидентов
func (r *IncidentRepository) List() ([]*domain.Incident, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	incidents := make([]*domain.Incident, 0, len(r.incidents))
	for _, incident := range r.incidents {
		if incident.IsActive {
			incidents = append(incidents, incident)
		}
	}

	return incidents, nil
}

// Update обновляет инцидент
func (r *IncidentRepository) Update(id string, req domain.UpdateIncidentRequest) (*domain.Incident, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	incident, exists := r.incidents[id]
	if !exists {
		return nil, fmt.Errorf("incident not found")
	}

	// Применяем изменения только если поля переданы
	if req.Title != nil {
		incident.Title = *req.Title
	}
	if req.Description != nil {
		incident.Description = *req.Description
	}
	if req.Severity != nil {
		incident.Severity = *req.Severity
	}
	if req.Latitude != nil {
		incident.Latitude = *req.Latitude
	}
	if req.Longitude != nil {
		incident.Longitude = *req.Longitude
	}
	if req.RadiusMeters != nil {
		incident.RadiusMeters = *req.RadiusMeters
	}

	incident.UpdatedAt = time.Now()
	return incident, nil
}

// Delete деактивирует инцидент (soft delete)
func (r *IncidentRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	incident, exists := r.incidents[id]
	if !exists {
		return fmt.Errorf("incident not found")
	}

	incident.IsActive = false
	incident.UpdatedAt = time.Now()
	return nil
}

// GetActiveIncidents получает все активные инциденты
func (r *IncidentRepository) GetActiveIncidents() ([]*domain.Incident, error) {
	return r.List()
}

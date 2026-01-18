package domain

import "time"

// Severity уровни опасности
type Severity string

const (
	SeverityLow    Severity = "low"
	SeverityMedium Severity = "medium"
	SeverityHigh   Severity = "high"
)

// Incident инцидент/опасная зона
type Incident struct {
	ID           string    `json:"id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	Severity     Severity  `json:"severity"`
	Latitude     float64   `json:"latitude"`
	Longitude    float64   `json:"longitude"`
	RadiusMeters int       `json:"radius_meters"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// CreateIncidentRequest запрос на создание инцидента
type CreateIncidentRequest struct {
	Title        string   `json:"title" binding:"required,min=3,max=200"`
	Description  string   `json:"description" binding:"max=1000"`
	Severity     Severity `json:"severity" binding:"required,oneof=low medium high"`
	Latitude     float64  `json:"latitude" binding:"required,min=-90,max=90"`
	Longitude    float64  `json:"longitude" binding:"required,min=-180,max=180"`
	RadiusMeters int      `json:"radius_meters" binding:"required,min=10,max=100000"`
}

// UpdateIncidentRequest запрос на обновление инцидента
type UpdateIncidentRequest struct {
	Title        *string   `json:"title" binding:"omitempty,min=3,max=200"`
	Description  *string   `json:"description" binding:"omitempty,max=1000"`
	Severity     *Severity `json:"severity" binding:"omitempty,oneof=low medium high"`
	Latitude     *float64  `json:"latitude" binding:"omitempty,min=-90,max=90"`
	Longitude    *float64  `json:"longitude" binding:"omitempty,min=-180,max=180"`
	RadiusMeters *int      `json:"radius_meters" binding:"omitempty,min=10,max=100000"`
}

// LocationCheckRequest запрос на проверку локации
type LocationCheckRequest struct {
	UserID    string  `json:"user_id" binding:"required,min=1,max=100"`
	Latitude  float64 `json:"latitude" binding:"required,min=-90,max=90"`
	Longitude float64 `json:"longitude" binding:"required,min=-180,max=180"`
}

// NearbyIncident инцидент рядом с локацией
type NearbyIncident struct {
	ID             string   `json:"id"`
	Title          string   `json:"title"`
	Severity       Severity `json:"severity"`
	Latitude       float64  `json:"latitude"`
	Longitude      float64  `json:"longitude"`
	RadiusMeters   int      `json:"radius_meters"`
	DistanceMeters float64  `json:"distance_meters"`
}

// LocationCheckResponse ответ на проверку локации
type LocationCheckResponse struct {
	CheckID        string           `json:"check_id"`
	IsInDangerZone bool             `json:"is_in_danger_zone"`
	CheckedAt      time.Time        `json:"checked_at"`
	Incidents      []NearbyIncident `json:"incidents"`
}

// LocationCheck запись о проверке локации
type LocationCheck struct {
	ID             string    `json:"id"`
	UserID         string    `json:"user_id"`
	Latitude       float64   `json:"latitude"`
	Longitude      float64   `json:"longitude"`
	IsInDangerZone bool      `json:"is_in_danger_zone"`
	CheckedAt      time.Time `json:"checked_at"`
}

// IncidentStats статистика по инциденту
type IncidentStats struct {
	IncidentID string `json:"incident_id"`
	Title      string `json:"title"`
	UserCount  int    `json:"user_count"`
}

package domain

import "time"

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

// CreateIncidentRequest запрос на создание
type CreateIncidentRequest struct {
	Title        string   `json:"title" binding:"required"`
	Description  string   `json:"description"`
	Severity     Severity `json:"severity" binding:"required"`
	Latitude     float64  `json:"latitude" binding:"required"`
	Longitude    float64  `json:"longitude" binding:"required"`
	RadiusMeters int      `json:"radius_meters" binding:"required"`
}

// LocationCheckRequest проверка локации
type LocationCheckRequest struct {
	UserID    string  `json:"user_id" binding:"required"`
	Latitude  float64 `json:"latitude" binding:"required"`
	Longitude float64 `json:"longitude" binding:"required"`
}

// LocationCheckResponse ответ
type LocationCheckResponse struct {
	CheckID        string    `json:"check_id"`
	IsInDangerZone bool      `json:"is_in_danger_zone"`
	CheckedAt      time.Time `json:"checked_at"`
}

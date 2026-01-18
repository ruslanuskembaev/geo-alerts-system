package domain

import "time"

// WebhookPayload тело вебхука
type WebhookPayload struct {
	CheckID        string           `json:"check_id"`
	UserID         string           `json:"user_id"`
	Latitude       float64          `json:"latitude"`
	Longitude      float64          `json:"longitude"`
	IsInDangerZone bool             `json:"is_in_danger_zone"`
	CheckedAt      time.Time        `json:"checked_at"`
	Incidents      []NearbyIncident `json:"incidents"`
}

// WebhookJob задача для очереди
type WebhookJob struct {
	Payload   WebhookPayload `json:"payload"`
	Attempt   int            `json:"attempt"`
	CreatedAt time.Time      `json:"created_at"`
}

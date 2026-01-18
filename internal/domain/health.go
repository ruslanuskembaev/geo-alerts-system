package domain

import "time"

// DependencyHealth represents a single dependency status.
type DependencyHealth struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

// SystemHealth represents overall health status.
type SystemHealth struct {
	Status       string                      `json:"status"`
	Service      string                      `json:"service"`
	Timestamp    time.Time                   `json:"timestamp"`
	Dependencies map[string]DependencyHealth `json:"dependencies,omitempty"`
}

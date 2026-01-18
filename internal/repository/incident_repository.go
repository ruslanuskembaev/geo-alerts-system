package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ruslanuskembaev/geo-alerts-system/internal/domain"
)

var ErrNotFound = errors.New("not found")

// IncidentRepository defines incident storage operations.
type IncidentRepository interface {
	Create(ctx context.Context, req domain.CreateIncidentRequest) (*domain.Incident, error)
	GetByID(ctx context.Context, id string) (*domain.Incident, error)
	List(ctx context.Context, limit, offset int) ([]*domain.Incident, int, error)
	Update(ctx context.Context, id string, req domain.UpdateIncidentRequest) (*domain.Incident, error)
	Deactivate(ctx context.Context, id string) error
	ListActive(ctx context.Context) ([]*domain.Incident, error)
}

// PostgresIncidentRepository implements IncidentRepository using PostgreSQL.
type PostgresIncidentRepository struct {
	db *pgxpool.Pool
}

func NewIncidentRepository(db *pgxpool.Pool) *PostgresIncidentRepository {
	return &PostgresIncidentRepository{db: db}
}

func (r *PostgresIncidentRepository) Create(ctx context.Context, req domain.CreateIncidentRequest) (*domain.Incident, error) {
	id := uuid.New().String()
	now := time.Now().UTC()

	incident := &domain.Incident{
		ID:           id,
		Title:        req.Title,
		Description:  req.Description,
		Severity:     req.Severity,
		Latitude:     req.Latitude,
		Longitude:    req.Longitude,
		RadiusMeters: req.RadiusMeters,
		IsActive:     true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	_, err := r.db.Exec(ctx, `
		INSERT INTO incidents (
			id, title, description, severity, latitude, longitude, radius_meters,
			is_active, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`, incident.ID, incident.Title, incident.Description, incident.Severity, incident.Latitude, incident.Longitude, incident.RadiusMeters, incident.IsActive, incident.CreatedAt, incident.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return incident, nil
}

func (r *PostgresIncidentRepository) GetByID(ctx context.Context, id string) (*domain.Incident, error) {
	var incident domain.Incident

	err := r.db.QueryRow(ctx, `
		SELECT id, title, description, severity, latitude, longitude, radius_meters,
		       is_active, created_at, updated_at
		FROM incidents
		WHERE id = $1
	`, id).Scan(
		&incident.ID,
		&incident.Title,
		&incident.Description,
		&incident.Severity,
		&incident.Latitude,
		&incident.Longitude,
		&incident.RadiusMeters,
		&incident.IsActive,
		&incident.CreatedAt,
		&incident.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return &incident, nil
}

func (r *PostgresIncidentRepository) List(ctx context.Context, limit, offset int) ([]*domain.Incident, int, error) {
	var total int
	if err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM incidents`).Scan(&total); err != nil {
		return nil, 0, err
	}

	rows, err := r.db.Query(ctx, `
		SELECT id, title, description, severity, latitude, longitude, radius_meters,
		       is_active, created_at, updated_at
		FROM incidents
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	incidents := make([]*domain.Incident, 0)
	for rows.Next() {
		var incident domain.Incident
		if err := rows.Scan(
			&incident.ID,
			&incident.Title,
			&incident.Description,
			&incident.Severity,
			&incident.Latitude,
			&incident.Longitude,
			&incident.RadiusMeters,
			&incident.IsActive,
			&incident.CreatedAt,
			&incident.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		incidents = append(incidents, &incident)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return incidents, total, nil
}

func (r *PostgresIncidentRepository) Update(ctx context.Context, id string, req domain.UpdateIncidentRequest) (*domain.Incident, error) {
	existing, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Title != nil {
		existing.Title = *req.Title
	}
	if req.Description != nil {
		existing.Description = *req.Description
	}
	if req.Severity != nil {
		existing.Severity = *req.Severity
	}
	if req.Latitude != nil {
		existing.Latitude = *req.Latitude
	}
	if req.Longitude != nil {
		existing.Longitude = *req.Longitude
	}
	if req.RadiusMeters != nil {
		existing.RadiusMeters = *req.RadiusMeters
	}

	existing.UpdatedAt = time.Now().UTC()

	_, err = r.db.Exec(ctx, `
		UPDATE incidents
		SET title = $2,
			description = $3,
			severity = $4,
			latitude = $5,
			longitude = $6,
			radius_meters = $7,
			is_active = $8,
			updated_at = $9
		WHERE id = $1
	`, existing.ID, existing.Title, existing.Description, existing.Severity, existing.Latitude, existing.Longitude, existing.RadiusMeters, existing.IsActive, existing.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return existing, nil
}

func (r *PostgresIncidentRepository) Deactivate(ctx context.Context, id string) error {
	tag, err := r.db.Exec(ctx, `
		UPDATE incidents
		SET is_active = false, updated_at = $2
		WHERE id = $1 AND is_active = true
	`, id, time.Now().UTC())
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *PostgresIncidentRepository) ListActive(ctx context.Context) ([]*domain.Incident, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, title, description, severity, latitude, longitude, radius_meters,
		       is_active, created_at, updated_at
		FROM incidents
		WHERE is_active = true
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	incidents := make([]*domain.Incident, 0)
	for rows.Next() {
		var incident domain.Incident
		if err := rows.Scan(
			&incident.ID,
			&incident.Title,
			&incident.Description,
			&incident.Severity,
			&incident.Latitude,
			&incident.Longitude,
			&incident.RadiusMeters,
			&incident.IsActive,
			&incident.CreatedAt,
			&incident.UpdatedAt,
		); err != nil {
			return nil, err
		}
		incidents = append(incidents, &incident)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return incidents, nil
}

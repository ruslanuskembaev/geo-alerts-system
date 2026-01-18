package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ruslanuskembaev/geo-alerts-system/internal/domain"
)

// LocationCheckRepository defines storage operations for location checks.
type LocationCheckRepository interface {
	Create(ctx context.Context, check domain.LocationCheck, incidentIDs []string) error
	StatsByIncident(ctx context.Context, since time.Time) ([]domain.IncidentStats, error)
}

// PostgresLocationCheckRepository implements LocationCheckRepository using PostgreSQL.
type PostgresLocationCheckRepository struct {
	db *pgxpool.Pool
}

func NewLocationCheckRepository(db *pgxpool.Pool) *PostgresLocationCheckRepository {
	return &PostgresLocationCheckRepository{db: db}
}

func (r *PostgresLocationCheckRepository) Create(ctx context.Context, check domain.LocationCheck, incidentIDs []string) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	_, err = tx.Exec(ctx, `
		INSERT INTO location_checks (
			id, user_id, latitude, longitude, is_in_danger_zone, checked_at
		) VALUES ($1, $2, $3, $4, $5, $6)
	`, check.ID, check.UserID, check.Latitude, check.Longitude, check.IsInDangerZone, check.CheckedAt)
	if err != nil {
		return err
	}

	if len(incidentIDs) > 0 {
		uuids := make([]uuid.UUID, 0, len(incidentIDs))
		for _, id := range incidentIDs {
			parsed, parseErr := uuid.Parse(id)
			if parseErr != nil {
				return parseErr
			}
			uuids = append(uuids, parsed)
		}

		_, err = tx.Exec(ctx, `
			INSERT INTO location_check_incidents (check_id, incident_id)
			SELECT $1, UNNEST($2::uuid[])
		`, check.ID, uuids)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *PostgresLocationCheckRepository) StatsByIncident(ctx context.Context, since time.Time) ([]domain.IncidentStats, error) {
	rows, err := r.db.Query(ctx, `
		SELECT i.id,
		       i.title,
		       COALESCE(COUNT(DISTINCT lc.user_id), 0) AS user_count
		FROM incidents i
		LEFT JOIN location_check_incidents lci ON i.id = lci.incident_id
		LEFT JOIN location_checks lc ON lc.id = lci.check_id AND lc.checked_at >= $1
		WHERE i.is_active = true
		GROUP BY i.id, i.title
		ORDER BY i.created_at DESC
	`, since)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats := make([]domain.IncidentStats, 0)
	for rows.Next() {
		var item domain.IncidentStats
		if err := rows.Scan(&item.IncidentID, &item.Title, &item.UserCount); err != nil {
			return nil, err
		}
		stats = append(stats, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return stats, nil
}

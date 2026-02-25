package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/rakhshon-mirzoev/gps-tracker/pkg/db"
)

// GpsAlertStatesRepository управляет состоянием "активного GPS-алерта" на сотрудника.
// Один сотрудник — максимум одна запись (PRIMARY KEY = staff_id).
type GpsAlertStatesRepository interface {
	// IsAlerted возвращает true, если алерт уже отправлен и сотрудник ещё не вернулся.
	IsAlerted(ctx context.Context, staffID uuid.UUID) (bool, error)
	// SetAlerted помечает сотрудника как "алерт отправлен".
	// Если запись уже есть — игнорирует (идемпотентно).
	SetAlerted(ctx context.Context, staffID uuid.UUID) error
	// Resolve удаляет запись: сотрудник прислал GPS, алерт сброшен.
	Resolve(ctx context.Context, staffID uuid.UUID) error
}

type gpsAlertStatesRepo struct {
	DB *db.DB
}

func NewGpsAlertStatesRepository(db *db.DB) GpsAlertStatesRepository {
	return &gpsAlertStatesRepo{DB: db}
}

func (r *gpsAlertStatesRepo) IsAlerted(ctx context.Context, staffID uuid.UUID) (bool, error) {
	var exists bool
	err := r.DB.Pool.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM gps_alert_states WHERE staff_id = $1)`,
		staffID,
	).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("gps alert states is alerted: %w", err)
	}
	return exists, nil
}

func (r *gpsAlertStatesRepo) SetAlerted(ctx context.Context, staffID uuid.UUID) error {
	_, err := r.DB.Pool.Exec(ctx,
		`INSERT INTO gps_alert_states (staff_id) VALUES ($1) ON CONFLICT DO NOTHING`,
		staffID,
	)
	if err != nil {
		return fmt.Errorf("gps alert states set alerted: %w", err)
	}
	return nil
}

func (r *gpsAlertStatesRepo) Resolve(ctx context.Context, staffID uuid.UUID) error {
	_, err := r.DB.Pool.Exec(ctx,
		`DELETE FROM gps_alert_states WHERE staff_id = $1`,
		staffID,
	)
	if err != nil {
		return fmt.Errorf("gps alert states resolve: %w", err)
	}
	return nil
}

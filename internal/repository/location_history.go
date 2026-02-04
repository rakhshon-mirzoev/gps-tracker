package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/model"
	"github.com/rakhshon-mirzoev/gps-tracker/pkg/db"
)

type LocationHistoryRepository interface {
	Create(ctx context.Context, lh *model.LocationHistory) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.LocationHistory, error)
	List(ctx context.Context, staffID uuid.UUID) ([]model.LocationHistory, error)
	Update(ctx context.Context, id uuid.UUID, dto model.UpdateLocationHistory) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetAllStaffOnlineOffline(ctx context.Context, supervisorID uuid.UUID, begin, end time.Time) (map[uuid.UUID]map[string]time.Duration, error)
}

type locationHistory struct {
	DB *db.DB
}

func NewLocationHistoryRepository(db *db.DB) LocationHistoryRepository {
	return &locationHistory{DB: db}
}

func (r *locationHistory) Create(ctx context.Context, lh *model.LocationHistory) error {
	sql := `
        INSERT INTO location_history (
            staff_id,
            latitude,
            longitude,
            battery_level,
            captured_at,
			status
        ) VALUES (
            $1, $2, $3, $4, $5, $6
        ) RETURNING id;
    `

	err := r.DB.Pool.QueryRow(ctx, sql, lh.StaffID, lh.Latitude, lh.Longitude, lh.BatteryLevel, lh.CapturedAt, lh.Status).Scan(&lh.ID)
	if err != nil {
		return fmt.Errorf("create location history: %w", err)
	}

	return nil
}

func (r *locationHistory) GetByID(ctx context.Context, id uuid.UUID) (*model.LocationHistory, error) {
	sql := `
        SELECT id, staff_id, latitude, longitude, battery_level, captured_at, status
        FROM location_history WHERE id = $1
    `

	var lh model.LocationHistory
	err := r.DB.Pool.QueryRow(ctx, sql, id).Scan(
		&lh.ID,
		&lh.StaffID,
		&lh.Latitude,
		&lh.Longitude,
		&lh.BatteryLevel,
		&lh.CapturedAt,
		&lh.Status,
	)
	if err != nil {
		return nil, fmt.Errorf("get by id location history: %w", err)
	}

	return &lh, nil
}

func (r *locationHistory) List(ctx context.Context, staffID uuid.UUID) ([]model.LocationHistory, error) {
	sql := `
        SELECT id, staff_id, latitude, longitude, battery_level, captured_at, status
        FROM location_history WHERE staff_id = $1
    `

	rows, err := r.DB.Pool.Query(ctx, sql, staffID)
	if err != nil {
		return nil, fmt.Errorf("list location history: %w", err)
	}

	var history = make([]model.LocationHistory, 0)

	for rows.Next() {
		var lh model.LocationHistory

		if err := rows.Scan(
			&lh.ID,
			&lh.StaffID,
			&lh.Latitude,
			&lh.Longitude,
			&lh.BatteryLevel,
			&lh.CapturedAt,
			&lh.Status,
		); err != nil {
			return nil, fmt.Errorf("scan location history: %w", err)
		}

		history = append(history, lh)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows location history: %w", err)
	}

	return history, nil
}

func (r *locationHistory) Update(ctx context.Context, id uuid.UUID, dto model.UpdateLocationHistory) error {
	sql := `UPDATE location_history SET `
	var args []any
	n := 1

	if dto.StaffID != nil {
		sql += fmt.Sprintf("staff_id = $%d, ", n)
		args = append(args, *dto.StaffID)
		n++
	}

	if dto.Latitude != nil {
		sql += fmt.Sprintf("latitude = $%d, ", n)
		args = append(args, *dto.Latitude)
		n++
	}

	if dto.Longitude != nil {
		sql += fmt.Sprintf("longitude = $%d, ", n)
		args = append(args, *dto.Longitude)
		n++
	}

	if dto.BatteryLevel != nil {
		sql += fmt.Sprintf("battery_level = $%d, ", n)
		args = append(args, *dto.BatteryLevel)
		n++
	}

	if dto.CapturedAt != nil {
		sql += fmt.Sprintf("captured_at = $%d, ", n)
		args = append(args, *dto.CapturedAt)
		n++
	}

	if len(args) == 0 {
		return nil
	}

	sql = strings.TrimSuffix(sql, ", ")
	sql += fmt.Sprintf(" WHERE id = $%d", n)
	args = append(args, id)

	_, err := r.DB.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("update location history: %w", err)
	}

	return nil
}

func (r *locationHistory) Delete(ctx context.Context, id uuid.UUID) error {
	sql := `
        DELETE FROM location_history WHERE id = $1
    `

	_, err := r.DB.Pool.Exec(ctx, sql, id)
	if err != nil {
		return fmt.Errorf("delete location history: %w", err)
	}

	return nil
}

func (r *locationHistory) GetAllStaffOnlineOffline(ctx context.Context, supervisorID uuid.UUID, begin, end time.Time) (map[uuid.UUID]map[string]time.Duration, error) {
	sql := `
        SELECT
			l.staff_id,
			l.status,
			COUNT(*) AS cnt
		FROM location_history l
		JOIN staff_users s ON s.id = l.staff_id
		JOIN groups g ON g.id = s.group_id
		WHERE g.supervisor_id = $1
		AND l.captured_at >= $2
		AND l.captured_at <= $3
		GROUP BY l.staff_id, l.status;
    `

	rows, err := r.DB.Pool.Query(ctx, sql, supervisorID, begin, end)
	if err != nil {
		return nil, fmt.Errorf("get all staff online offline: %w", err)
	}
	defer rows.Close()

	result := make(map[uuid.UUID]map[string]time.Duration)

	for rows.Next() {
		var staffID uuid.UUID
		var status string
		var count int

		if err := rows.Scan(&staffID, &status, &count); err != nil {
			return nil, fmt.Errorf("scan staff online offline: %w", err)
		}

		if _, ok := result[staffID]; !ok {
			result[staffID] = map[string]time.Duration{
				"online":  0,
				"offline": 0,
			}
		}

		duration := time.Duration(count*5) * time.Minute

		if status == "offline" {
			result[staffID]["offline"] += duration
		} else {
			result[staffID]["online"] += duration
		}
	}

	return result, nil
}

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
	CreateBatch(ctx context.Context, locations []model.LocationHistory) error
	LatestByGroupID(ctx context.Context, groupID uuid.UUID) ([]model.LocationHistoryWithStaff, error)
	LatestAll(ctx context.Context) ([]model.LocationHistoryWithStaff, error)
	// GetStaffMissingGPS возвращает активных сотрудников, чей последний GPS
	// старше threshold (или которые вообще не отправляли GPS).
	GetStaffMissingGPS(ctx context.Context, threshold time.Duration) ([]model.LocationHistoryWithStaff, error)
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

func (r *locationHistory) CreateBatch(ctx context.Context, locations []model.LocationHistory) error {
	tx, err := r.DB.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

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

	for i := range locations {
		err := tx.QueryRow(
			ctx,
			sql,
			locations[i].StaffID,
			locations[i].Latitude,
			locations[i].Longitude,
			locations[i].BatteryLevel,
			locations[i].CapturedAt,
			locations[i].Status,
		).Scan(&locations[i].ID)

		if err != nil {
			return fmt.Errorf("insert location[%d]: %w", i, err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit tx: %w", err)
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
        FROM location_history WHERE staff_id = $1 ORDER BY captured_at DESC
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

func (r *locationHistory) LatestByGroupID(ctx context.Context, groupID uuid.UUID) ([]model.LocationHistoryWithStaff, error) {
	sql := `
		SELECT DISTINCT ON (s.id)
			s.id, s.group_id, s.full_name,
			s.phone_number,
			l.latitude, l.longitude, l.battery_level, l.captured_at, l.status
		FROM staff_users s
		JOIN location_history l ON l.staff_id = s.id
		WHERE s.group_id = $1 AND s.is_active = true
		ORDER BY s.id, l.captured_at DESC
	`

	rows, err := r.DB.Pool.Query(ctx, sql, groupID)
	if err != nil {
		return nil, fmt.Errorf("latest by group id: %w", err)
	}
	defer rows.Close()

	var result []model.LocationHistoryWithStaff
	for rows.Next() {
		var loc model.LocationHistoryWithStaff
		if err := rows.Scan(
			&loc.StaffID, &loc.GroupID, &loc.FullName,
			&loc.PhoneNumber,
			&loc.Latitude, &loc.Longitude, &loc.BatteryLevel, &loc.CapturedAt, &loc.Status,
		); err != nil {
			return nil, fmt.Errorf("scan latest by group: %w", err)
		}
		result = append(result, loc)
	}

	return result, rows.Err()
}

func (r *locationHistory) LatestAll(ctx context.Context) ([]model.LocationHistoryWithStaff, error) {
	sql := `
		SELECT DISTINCT ON (s.id)
			s.id, s.group_id, s.full_name,
			s.phone_number,
			l.latitude, l.longitude, l.battery_level, l.captured_at, l.status
		FROM staff_users s
		JOIN location_history l ON l.staff_id = s.id
		WHERE s.is_active = true
		ORDER BY s.id, l.captured_at DESC
	`

	rows, err := r.DB.Pool.Query(ctx, sql)
	if err != nil {
		return nil, fmt.Errorf("latest all: %w", err)
	}
	defer rows.Close()

	var result []model.LocationHistoryWithStaff
	for rows.Next() {
		var loc model.LocationHistoryWithStaff
		if err := rows.Scan(
			&loc.StaffID, &loc.GroupID, &loc.FullName,
			&loc.PhoneNumber,
			&loc.Latitude, &loc.Longitude, &loc.BatteryLevel, &loc.CapturedAt, &loc.Status,
		); err != nil {
			return nil, fmt.Errorf("scan latest all: %w", err)
		}
		result = append(result, loc)
	}

	return result, rows.Err()
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

// GetStaffMissingGPS возвращает сотрудников у которых:
//   - день начат (day_started = TRUE)
//   - они активны (не деактивированы)
//   - последний GPS старше `threshold` (или GPS не было вообще)
func (r *locationHistory) GetStaffMissingGPS(ctx context.Context, threshold time.Duration) ([]model.LocationHistoryWithStaff, error) {
	sql := `
		SELECT DISTINCT ON (s.id)
			s.id, s.group_id, s.full_name, s.phone_number,
			COALESCE(l.latitude, 0),
			COALESCE(l.longitude, 0),
			COALESCE(l.battery_level, 0),
			COALESCE(l.captured_at, s.start_at),
			COALESCE(l.status, 'offline')
		FROM staff_users s
		LEFT JOIN location_history l ON l.staff_id = s.id
		WHERE s.is_active          = TRUE
		  AND s.deactivated_by_id IS NULL
		  AND s.day_started        = TRUE
		ORDER BY s.id, l.captured_at DESC NULLS LAST
	`

	rows, err := r.DB.Pool.Query(ctx, sql)
	if err != nil {
		return nil, fmt.Errorf("get staff missing gps: %w", err)
	}
	defer rows.Close()

	cutoff := time.Now().Add(-threshold)

	var result []model.LocationHistoryWithStaff
	for rows.Next() {
		var loc model.LocationHistoryWithStaff
		if err := rows.Scan(
			&loc.StaffID, &loc.GroupID, &loc.FullName, &loc.PhoneNumber,
			&loc.Latitude, &loc.Longitude, &loc.BatteryLevel, &loc.CapturedAt, &loc.Status,
		); err != nil {
			return nil, fmt.Errorf("scan staff missing gps: %w", err)
		}
		// Фильтруем в Go: только те, чей последний GPS старше порога
		if loc.CapturedAt.Before(cutoff) {
			result = append(result, loc)
		}
	}

	return result, rows.Err()
}

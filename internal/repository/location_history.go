package repository

import (
    "context"

    "github.com/google/uuid"
    "github.com/rakhshon-mirzoev/gps-tracker/internal/model"
    "github.com/rakhshon-mirzoev/gps-tracker/pkg/db"
)

type LocationHistoryRepository interface {
    Create(ctx context.Context, lh *model.LocationHistory) error
    GetByID(ctx context.Context, id uuid.UUID) (*model.LocationHistory, error)
    List(ctx context.Context, staffID uuid.UUID) ([]model.LocationHistory, error)
    Update(ctx context.Context, lh *model.LocationHistory) error
    Delete(ctx context.Context, id uuid.UUID) error
}

type locationHistory struct {
    DB *db.DB
}

func InitLocationHistoryRepository(db *db.DB) LocationHistoryRepository {
    return &locationHistory{DB: db}
}

func (r *locationHistory) Create(ctx context.Context, lh *model.LocationHistory) error {
    sql := `
        INSERT INTO location_history (
            staff_id,
            location,
            battery_level,
            captured_at
        ) VALUES (
            $1, ST_GeomFromText($2, 4326), $3, $4
        ) RETURNING id;
    `

    err := r.DB.Pool.QueryRow(ctx, sql, lh.StaffID, lh.Location, lh.BatteryLevel, lh.CapturedAt).Scan(&lh.ID)
    if err != nil {
        return err
    }

    return nil
}

func (r *locationHistory) GetByID(ctx context.Context, id uuid.UUID) (*model.LocationHistory, error) {
    sql := `
        SELECT id, staff_id, ST_AsText(location) AS location, battery_level, captured_at
        FROM location_history WHERE id = $1
    `

    var lh model.LocationHistory
    err := r.DB.Pool.QueryRow(ctx, sql, id).Scan(
        &lh.ID,
        &lh.StaffID,
        &lh.Location,
        &lh.BatteryLevel,
        &lh.CapturedAt,
    )
    if err != nil {
        return nil, err
    }

    return &lh, nil
}

func (r *locationHistory) List(ctx context.Context, staffID uuid.UUID) ([]model.LocationHistory, error) {
    sql := `
        SELECT id, staff_id, ST_AsText(location) AS location, battery_level, captured_at
        FROM location_history WHERE staff_id = $1
    `

    rows, err := r.DB.Pool.Query(ctx, sql, staffID)
    if err != nil {
        return nil, err
    }

    var history []model.LocationHistory

    for rows.Next() {
        var lh model.LocationHistory

        if err := rows.Scan(
            &lh.ID,
            &lh.StaffID,
            &lh.Location,
            &lh.BatteryLevel,
            &lh.CapturedAt,
        ); err != nil {
            return nil, err
        }

        history = append(history, lh)
    }

    if err := rows.Err(); err != nil {
        return nil, err
    }

    return history, nil
}

func (r *locationHistory) Update(ctx context.Context, lh *model.LocationHistory) error {
    sql := `
        UPDATE location_history
        SET staff_id = $1, location = ST_GeomFromText($2, 4326), battery_level = $3, captured_at = $4
        WHERE id = $5
    `

    _, err := r.DB.Pool.Exec(ctx, sql, lh.StaffID, lh.Location, lh.BatteryLevel, lh.CapturedAt, lh.ID)
    if err != nil {
        return err
    }

    return nil
}

func (r *locationHistory) Delete(ctx context.Context, id uuid.UUID) error {
    sql := `
        DELETE FROM location_history WHERE id = $1
    `

    _, err := r.DB.Pool.Exec(ctx, sql, id)
    if err != nil {
        return err
    }

    return nil
}
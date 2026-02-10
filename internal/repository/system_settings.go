package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/model"
	"github.com/rakhshon-mirzoev/gps-tracker/pkg/db"
	errs "github.com/rakhshon-mirzoev/gps-tracker/pkg/error"
)

type SystemSettingsRepository interface {
	Create(ctx context.Context, ss *model.SystemSettings) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.SystemSettings, error)
	List(ctx context.Context, limit, offset int) ([]model.SystemSettings, int, error)
	Update(ctx context.Context, id uuid.UUID, dto model.UpdateSystemSettings) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type systemSettings struct {
	DB *db.DB
}

func NewSystemSettingsRepository(db *db.DB) SystemSettingsRepository {
	return &systemSettings{DB: db}
}

func (r *systemSettings) Create(ctx context.Context, ss *model.SystemSettings) error {
	sql := `
        INSERT INTO system_settings (
            key,
            value,
            description
        ) VALUES (
            $1, $2, $3
        ) RETURNING id;
    `

	err := r.DB.Pool.QueryRow(ctx, sql, ss.Key, ss.Value, ss.Description).Scan(&ss.ID)
	if err != nil {
		return fmt.Errorf("create system settings: %w", err)
	}

	return nil
}

func (r *systemSettings) GetByID(ctx context.Context, id uuid.UUID) (*model.SystemSettings, error) {
	sql := `
        SELECT id, key, value, description
        FROM system_settings WHERE id = $1
    `

	var ss model.SystemSettings
	err := r.DB.Pool.QueryRow(ctx, sql, id).Scan(
		&ss.ID,
		&ss.Key,
		&ss.Value,
		&ss.Description,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.ErrorNotFound
		}
		return nil, fmt.Errorf("get by id system settings: %w", err)
	}
	return &ss, nil
}

func (r *systemSettings) List(ctx context.Context, limit, offset int) ([]model.SystemSettings, int, error) {
	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	var total int
	countSQL := `SELECT COUNT(*) FROM system_settings`
	if err := r.DB.Pool.QueryRow(ctx, countSQL).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count system settings: %w", err)
	}

	sql := `
        SELECT id, key, value, description
        FROM system_settings
        ORDER BY key
        LIMIT $1 OFFSET $2
    `

	rows, err := r.DB.Pool.Query(ctx, sql, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("list system settings: %w", err)
	}

	var settings = make([]model.SystemSettings, 0)

	for rows.Next() {
		var ss model.SystemSettings

		if err := rows.Scan(
			&ss.ID,
			&ss.Key,
			&ss.Value,
			&ss.Description,
		); err != nil {
			return nil, 0, fmt.Errorf("scan system settings: %w", err)
		}

		settings = append(settings, ss)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("rows system settings: %w", err)
	}

	return settings, total, nil
}

func (r *systemSettings) Update(ctx context.Context, id uuid.UUID, dto model.UpdateSystemSettings) error {
	sql := `UPDATE system_settings SET `
	var args []any
	n := 1

	if dto.Key != nil {
		sql += fmt.Sprintf("key = $%d, ", n)
		args = append(args, *dto.Key)
		n++
	}

	if dto.Value != nil {
		sql += fmt.Sprintf("value = $%d, ", n)
		args = append(args, *dto.Value)
		n++
	}

	if dto.Description != nil {
		sql += fmt.Sprintf("description = $%d, ", n)
		args = append(args, *dto.Description)
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
		return fmt.Errorf("update system settings: %w", err)
	}

	return nil
}

func (r *systemSettings) Delete(ctx context.Context, id uuid.UUID) error {
	sql := `
        DELETE FROM system_settings WHERE id = $1
    `

	_, err := r.DB.Pool.Exec(ctx, sql, id)
	if err != nil {
		return fmt.Errorf("delete system settings: %w", err)
	}

	return nil
}

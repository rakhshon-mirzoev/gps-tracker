package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/model"
	"github.com/rakhshon-mirzoev/gps-tracker/pkg/db"
)

type SystemSettingsRepository interface {
	Create(ctx context.Context, ss *model.SystemSettings) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.SystemSettings, error)
	List(ctx context.Context) ([]model.SystemSettings, error)
	Update(ctx context.Context, ss *model.SystemSettings) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type systemSettings struct {
	DB *db.DB
}

func InitSystemSettingsRepository(db *db.DB) SystemSettingsRepository {
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
		return err
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
		return nil, err
	}

	return &ss, nil
}

func (r *systemSettings) List(ctx context.Context) ([]model.SystemSettings, error) {
	sql := `
        SELECT id, key, value, description
        FROM system_settings
    `

	rows, err := r.DB.Pool.Query(ctx, sql)
	if err != nil {
		return nil, err
	}

	var settings []model.SystemSettings

	for rows.Next() {
		var ss model.SystemSettings

		if err := rows.Scan(
			&ss.ID,
			&ss.Key,
			&ss.Value,
			&ss.Description,
		); err != nil {
			return nil, err
		}

		settings = append(settings, ss)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return settings, nil
}

func (r *systemSettings) Update(ctx context.Context, ss *model.SystemSettings) error {
	sql := `
        UPDATE system_settings
        SET key = $1, value = $2, description = $3
        WHERE id = $4
    `

	_, err := r.DB.Pool.Exec(ctx, sql, ss.Key, ss.Value, ss.Description, ss.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *systemSettings) Delete(ctx context.Context, id uuid.UUID) error {
	sql := `
        DELETE FROM system_settings WHERE id = $1
    `

	_, err := r.DB.Pool.Exec(ctx, sql, id)
	if err != nil {
		return err
	}

	return nil
}

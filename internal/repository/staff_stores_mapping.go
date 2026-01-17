package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/model"
	"github.com/rakhshon-mirzoev/gps-tracker/pkg/db"
)

type StaffStoresMappingRepository interface {
	Create(ctx context.Context, ssm *model.StaffStoresMapping) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.StaffStoresMapping, error)
	List(ctx context.Context, staffID uuid.UUID) ([]model.StaffStoresMapping, error)
	Update(ctx context.Context, ssm *model.StaffStoresMapping) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type staffStoresMapping struct {
	DB *db.DB
}

func InitStaffStoresMappingRepository(db *db.DB) StaffStoresMappingRepository {
	return &staffStoresMapping{DB: db}
}

func (r *staffStoresMapping) Create(ctx context.Context, ssm *model.StaffStoresMapping) error {
	sql := `
        INSERT INTO staff_stores_mapping (
            staff_id,
            store_id
        ) VALUES (
            $1, $2
        ) RETURNING id;
    `

	err := r.DB.Pool.QueryRow(ctx, sql, ssm.StaffID, ssm.StoreID).Scan(&ssm.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *staffStoresMapping) GetByID(ctx context.Context, id uuid.UUID) (*model.StaffStoresMapping, error) {
	sql := `
        SELECT id, staff_id, store_id
        FROM staff_stores_mapping WHERE id = $1
    `

	var ssm model.StaffStoresMapping
	err := r.DB.Pool.QueryRow(ctx, sql, id).Scan(
		&ssm.ID,
		&ssm.StaffID,
		&ssm.StoreID,
	)
	if err != nil {
		return nil, err
	}

	return &ssm, nil
}

func (r *staffStoresMapping) List(ctx context.Context, staffID uuid.UUID) ([]model.StaffStoresMapping, error) {
	sql := `
        SELECT id, staff_id, store_id
        FROM staff_stores_mapping WHERE staff_id = $1
    `

	rows, err := r.DB.Pool.Query(ctx, sql, staffID)
	if err != nil {
		return nil, err
	}

	var mappings []model.StaffStoresMapping

	for rows.Next() {
		var ssm model.StaffStoresMapping

		if err := rows.Scan(
			&ssm.ID,
			&ssm.StaffID,
			&ssm.StoreID,
		); err != nil {
			return nil, err
		}

		mappings = append(mappings, ssm)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return mappings, nil
}

func (r *staffStoresMapping) Update(ctx context.Context, ssm *model.StaffStoresMapping) error {
	sql := `
        UPDATE staff_stores_mapping
        SET staff_id = $1, store_id = $2
        WHERE id = $3
    `

	_, err := r.DB.Pool.Exec(ctx, sql, ssm.StaffID, ssm.StoreID, ssm.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *staffStoresMapping) Delete(ctx context.Context, id uuid.UUID) error {
	sql := `
        DELETE FROM staff_stores_mapping WHERE id = $1
    `

	_, err := r.DB.Pool.Exec(ctx, sql, id)
	if err != nil {
		return err
	}

	return nil
}

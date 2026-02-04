package repository

import (
	"context"
	"fmt"

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
	ListByStaffID(ctx context.Context, id uuid.UUID, lim, off int) ([]model.Store, error)
}

type staffStoresMapping struct {
	DB *db.DB
}

func NewStaffStoresMappingRepository(db *db.DB) StaffStoresMappingRepository {
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
		return fmt.Errorf("create staff stores mapping: %w", err)
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
		return nil, fmt.Errorf("get by id staff stores mapping: %w", err)
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
		return nil, fmt.Errorf("list staff stores mapping: %w", err)
	}

	var mappings = make([]model.StaffStoresMapping, 0)

	for rows.Next() {
		var ssm model.StaffStoresMapping

		if err := rows.Scan(
			&ssm.ID,
			&ssm.StaffID,
			&ssm.StoreID,
		); err != nil {
			return nil, fmt.Errorf("scan staff stores mapping: %w", err)
		}

		mappings = append(mappings, ssm)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows staff stores mapping: %w", err)
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
		return fmt.Errorf("update staff stores mapping: %w", err)
	}

	return nil
}

func (r *staffStoresMapping) Delete(ctx context.Context, id uuid.UUID) error {
	sql := `
        DELETE FROM staff_stores_mapping WHERE id = $1
    `

	_, err := r.DB.Pool.Exec(ctx, sql, id)
	if err != nil {
		return fmt.Errorf("delete staff stores mapping: %w", err)
	}

	return nil
}

func (r *staffStoresMapping) ListByStaffID(ctx context.Context, id uuid.UUID, lim, off int) ([]model.Store, error) {
	defaultLim := 5

	if lim <= 0 {
		lim = defaultLim
	}

	if off < 0 {
		off = 0
	}

	sqlStores := `
		SELECT s.id, s.name, s.address, s.latitude, s.longitude, s.owner_full_name, s.owner_inn, 
		       s.owner_passport, s.status, s.rejection_comment, s.created_by_staff_id, s.created_by_admin_id, 
		       s.is_active, s.is_frozen, s.deactivated_at, s.deactivated_by_id, s.deactivation_reason, s.created_at
		FROM stores s
		JOIN staff_stores_mapping ssm ON s.id = ssm.store_id
		WHERE ssm.staff_id = $1
		ORDER BY s.created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.DB.Pool.Query(ctx, sqlStores, id, lim, off)
	if err != nil {
		return nil, fmt.Errorf("list stores by staff id: %w", err)
	}
	defer rows.Close()

	var stores = make([]model.Store, 0)

	for rows.Next() {
		var store model.Store

		err := rows.Scan(
			&store.ID,
			&store.Name,
			&store.Address,
			&store.Latitude,
			&store.Longitude,
			&store.OwnerFullName,
			&store.OwnerINN,
			&store.OwnerPassport,
			&store.Status,
			&store.RejectionComment,
			&store.CreatedByStaffID,
			&store.CreatedByAdminID,
			&store.IsActive,
			&store.IsFrozen,
			&store.DeactivatedAt,
			&store.DeactivatedByID,
			&store.DeactivationReason,
			&store.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan store: %w", err)
		}

		stores = append(stores, store)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows stores: %w", err)
	}

	return stores, nil
}

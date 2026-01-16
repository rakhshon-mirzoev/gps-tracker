package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/model"
	"github.com/rakhshon-mirzoev/gps-tracker/pkg/db"
	errs "github.com/rakhshon-mirzoev/gps-tracker/pkg/error"
	"github.com/rakhshon-mirzoev/gps-tracker/pkg/utils"
)

type stores struct {
	DB *db.DB
}

type StoreRepository interface {
	Create(ctx context.Context, store *model.Store) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Store, error)
	Update(ctx context.Context, store *model.Store) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, filter model.StoreListFilter) ([]*model.Store, int, error)
	Deactivate(ctx context.Context, id, deactivatedBy uuid.UUID, reason string) error
}

func NewStoreRepository(db *db.DB) StoreRepository {
	return &stores{DB: db}
}

func (r *stores) Create(ctx context.Context, store *model.Store) error {
	sql := `
		INSERT INTO stores (
			name,
			address,
			location,
			owner_full_name,
			owner_inn,
			owner_passport,
			status,
			created_by_staff_id,
			created_by_admin_id,
			is_active,
			is_frozen,
			created_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, NOW())
		RETURNING id, created_at;
	`
	err := r.DB.Pool.QueryRow(
		ctx,
		sql,
		store.Name,
		store.Address,
		store.Location,
		store.OwnerFullName,
		store.OwnerINN,
		store.OwnerPassport,
		store.Status,
		store.CreatedByStaffID,
		store.CreatedByAdminID,
		store.IsActive,
		store.IsFrozen,
	).Scan(&store.ID, &store.CreatedAt)

	if err != nil {
		return fmt.Errorf("create store: %w", err)
	}

	return nil
}

func (r *stores) GetByID(ctx context.Context, id uuid.UUID) (*model.Store, error) {
	store := &model.Store{}

	sql := `
		SELECT 
			id, name, address, location,
			owner_full_name, owner_inn, owner_passport,
			status, rejection_comment,
			created_by_staff_id, created_by_admin_id,
			is_active, is_frozen,
			deactivated_at, deactivated_by_id, deactivation_reason,
			created_at
		FROM stores 
		WHERE id = $1`

	err := r.DB.Pool.QueryRow(ctx, sql, id).Scan(
		&store.ID,
		&store.Name,
		&store.Address,
		&store.Location,
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
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.ErrorNotFound
		}
		return nil, fmt.Errorf("get store by id: %w", err)
	}

	return store, nil
}

func (r *stores) Update(ctx context.Context, store *model.Store) error {
	sql := `
		UPDATE stores
		SET 
			name = $1,
			address = $2,
			owner_full_name = $3,
			owner_inn = $4,
			owner_passport = $5,
			status = $6,
			rejection_comment = $7,
			is_active = $8,
			is_frozen = $9
		WHERE id = $10`

	ct, err := r.DB.Pool.Exec(ctx, sql,
		store.Name,
		store.Address,
		store.OwnerFullName,
		store.OwnerINN,
		store.OwnerPassport,
		store.Status,
		store.RejectionComment,
		store.IsActive,
		store.IsFrozen,
		store.ID,
	)

	if err != nil {
		return fmt.Errorf("update store: %w", err)
	}

	if ct.RowsAffected() == 0 {
		return errs.ErrorNotFound
	}

	return nil
}

func (r *stores) Delete(ctx context.Context, id uuid.UUID) error {
	sql := `DELETE FROM stores WHERE id = $1`

	ct, err := r.DB.Pool.Exec(ctx, sql, id)
	if err != nil {
		return fmt.Errorf("delete store: %w", err)
	}

	if ct.RowsAffected() == 0 {
		return errs.ErrorNotFound
	}

	return nil
}

func (r *stores) Deactivate(ctx context.Context, id, deactivatedBy uuid.UUID, reason string) error {
	sql := `
		UPDATE stores
		SET 
			is_active = false,
			deactivated_at = NOW(),
			deactivated_by_id = $1,
			deactivation_reason = $2
		WHERE id = $3`

	ct, err := r.DB.Pool.Exec(ctx, sql, deactivatedBy, reason, id)
	if err != nil {
		return fmt.Errorf("deactivate store: %w", err)
	}

	if ct.RowsAffected() == 0 {
		return errs.ErrorNotFound
	}

	return nil
}

func (r *stores) List(ctx context.Context, f model.StoreListFilter) ([]*model.Store, int, error) {
	var conditions []string
	var args []interface{}
	argNum := 1

	if f.Status != nil {
		conditions = append(conditions, fmt.Sprintf("status = $%d", argNum))
		args = append(args, *f.Status)
		argNum++
	}

	if f.IsActive != nil {
		conditions = append(conditions, fmt.Sprintf("is_active = $%d", argNum))
		args = append(args, *f.IsActive)
		argNum++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + utils.Join(conditions, " AND ")
	}

	orderBy := "created_at DESC"
	if f.SortBy != "" {
		orderBy = f.SortBy
		if f.SortOrder == "ASC" {
			orderBy += " ASC"
		} else {
			orderBy += " DESC"
		}
	}

	limit := 20
	offset := 0
	if f.Limit > 0 {
		limit = f.Limit
	}
	if f.Offset > 0 {
		offset = f.Offset
	}

	baseSQL := `
		SELECT 
			id, name, address, location,
			owner_full_name, owner_inn, owner_passport,
			status, rejection_comment,
			created_by_staff_id, created_by_admin_id,
			is_active, is_frozen,
			created_at
		FROM stores %s
		ORDER BY %s
		LIMIT $%d OFFSET $%d`

	countSQL := `SELECT COUNT(*) FROM stores ` + whereClause

	// Считаем общее количество
	var total int
	err := r.DB.Pool.QueryRow(ctx, countSQL, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("count stores: %w", err)
	}

	// Добавляем limit/offset
	args = append(args, limit, offset)

	query := fmt.Sprintf(baseSQL, whereClause, orderBy, argNum, argNum+1)

	rows, err := r.DB.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("list stores: %w", err)
	}
	defer rows.Close()

	var storesList []*model.Store

	for rows.Next() {
		s := &model.Store{}
		if err := rows.Scan(
			&s.ID, &s.Name, &s.Address, &s.Location,
			&s.OwnerFullName, &s.OwnerINN, &s.OwnerPassport,
			&s.Status, &s.RejectionComment,
			&s.CreatedByStaffID, &s.CreatedByAdminID,
			&s.IsActive, &s.IsFrozen,
			&s.CreatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("scan store: %w", err)
		}
		storesList = append(storesList, s)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("rows error: %w", err)
	}

	return storesList, total, nil
}

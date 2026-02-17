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
	Update(ctx context.Context, store model.UpdateStore, id uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, isActive, status, sortBy, sortOrder string, lim, off int) ([]model.Store, int, error)
	Deactivate(ctx context.Context, id, deactivatedBy uuid.UUID, reason string) error
	IfExists(ctx context.Context, name string) (bool, error)
	BeginTx(ctx context.Context) (pgx.Tx, error)
	IfExistsTx(ctx context.Context, tx pgx.Tx, name string) (bool, error)
	CreateTx(ctx context.Context, tx pgx.Tx, st *model.Store) error
}

func NewStoreRepository(db *db.DB) StoreRepository {
	return &stores{DB: db}
}

func (r *stores) BeginTx(ctx context.Context) (pgx.Tx, error) {
	return r.DB.Pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.ReadCommitted,
	})

}

func (s *stores) IfExistsTx(ctx context.Context, tx pgx.Tx, name string) (bool, error) {
	sql := `SELECT EXISTS(SELECT 1 FROM stores WHERE name=$1)`
	var exists bool
	if err := tx.QueryRow(ctx, sql, name).Scan(&exists); err != nil {
		return false, fmt.Errorf("check exists: %w", err)
	}
	return exists, nil
}

func (r *stores) CreateTx(ctx context.Context, tx pgx.Tx, st *model.Store) error {
	sql := `
		INSERT INTO stores (
			name,
			address,
			latitude,
			longitude,
			owner_full_name,
			owner_inn,
			owner_passport,
			owner_phone,
			additional_phone,
			comment,
			created_by_staff_id
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id;
	`
	if err := r.DB.Pool.QueryRow(
		ctx,
		sql,
		st.Name,
		st.Address,
		st.Latitude,
		st.Longitude,
		st.OwnerFullName,
		st.OwnerINN,
		st.OwnerPassport,
		st.OwnerPhone,
		st.AdditionalPhone,
		st.Comment,
		st.CreatedByStaffID,
	).Scan(&st.ID); err != nil {
		return fmt.Errorf("create store tx: %w", err)
	}
	return nil
}

func (r *stores) Create(ctx context.Context, store *model.Store) error {
	sql := `
		INSERT INTO stores (
			name,
			address,
			latitude,
			longitude,
			owner_full_name,
			owner_inn,
			owner_passport,
			owner_phone,
			additional_phone,
			comment,
			created_by_staff_id
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id;
	`
	err := r.DB.Pool.QueryRow(
		ctx,
		sql,
		store.Name,
		store.Address,
		store.Latitude,
		store.Longitude,
		store.OwnerFullName,
		store.OwnerINN,
		store.OwnerPassport,
		store.OwnerPhone,
		store.AdditionalPhone,
		store.Comment,
		store.CreatedByStaffID,
	).Scan(&store.ID)

	if err != nil {
		return fmt.Errorf("create store: %w", err)
	}

	return nil
}

func (r *stores) GetByID(ctx context.Context, id uuid.UUID) (*model.Store, error) {
	store := &model.Store{}

	sql := `
		SELECT 
			id, name, address, latitude, longitude,
			owner_full_name, owner_inn, owner_passport,
			owner_phone, additional_phone, comment,
			status, rejection_comment,
			created_by_staff_id,
			is_active, is_frozen,
			deactivated_at, deactivated_by_id, deactivation_reason,
			created_at
		FROM stores 
		WHERE id = $1`

	err := r.DB.Pool.QueryRow(ctx, sql, id).Scan(
		&store.ID,
		&store.Name,
		&store.Address,
		&store.Latitude,
		&store.Longitude,
		&store.OwnerFullName,
		&store.OwnerINN,
		&store.OwnerPassport,
		&store.OwnerPhone,
		&store.AdditionalPhone,

		&store.Comment,
		&store.Status,
		&store.RejectionComment,
		&store.CreatedByStaffID,
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
		return nil, err
	}

	return store, nil
}

func (r *stores) Update(ctx context.Context, store model.UpdateStore, id uuid.UUID) error {
	var updates []string
	var args []any
	argNum := 1

	if store.Name != nil {
		updates = append(updates, fmt.Sprintf("name = $%d", argNum))
		args = append(args, *store.Name)
		argNum++
	}

	if store.Address != nil {
		updates = append(updates, fmt.Sprintf("address = $%d", argNum))
		args = append(args, *store.Address)
		argNum++
	}

	if store.Latitude != nil {
		updates = append(updates, fmt.Sprintf("latitude = $%d", argNum))
		args = append(args, *store.Latitude)
		argNum++
	}

	if store.Longitude != nil {
		updates = append(updates, fmt.Sprintf("longitude = $%d", argNum))
		args = append(args, *store.Longitude)
		argNum++
	}

	if store.OwnerFullName != nil {
		updates = append(updates, fmt.Sprintf("owner_full_name = $%d", argNum))
		args = append(args, *store.OwnerFullName)
		argNum++
	}

	if store.OwnerINN != nil {
		updates = append(updates, fmt.Sprintf("owner_inn = $%d", argNum))
		args = append(args, *store.OwnerINN)
		argNum++
	}

	if store.OwnerPassport != nil {
		updates = append(updates, fmt.Sprintf("owner_passport = $%d", argNum))
		args = append(args, *store.OwnerPassport)
		argNum++
	}

	if store.OwnerPhone != nil {
		updates = append(updates, fmt.Sprintf("owner_phone = $%d", argNum))
		args = append(args, *store.OwnerPhone)
		argNum++
	}

	if store.AdditionalPhone != nil {
		updates = append(updates, fmt.Sprintf("additional_phone = $%d", argNum))
		args = append(args, *store.AdditionalPhone)
		argNum++
	}

	if store.Comment != nil {
		updates = append(updates, fmt.Sprintf("comment = $%d", argNum))
		args = append(args, *store.Comment)
		argNum++
	}

	if store.Status != nil {
		updates = append(updates, fmt.Sprintf("status = $%d", argNum))
		args = append(args, *store.Status)
		argNum++
	}

	if store.RejectionComment != nil {
		updates = append(updates, fmt.Sprintf("rejection_comment = $%d", argNum))
		args = append(args, *store.RejectionComment)
		argNum++
	}

	if store.IsActive != nil {
		updates = append(updates, fmt.Sprintf("is_active = $%d", argNum))
		args = append(args, *store.IsActive)
		argNum++
	}

	if store.IsFrozen != nil {
		updates = append(updates, fmt.Sprintf("is_frozen = $%d", argNum))
		args = append(args, *store.IsFrozen)
		argNum++
	}

	if store.CreatedByStaffID != nil {
		staffID, err := uuid.Parse(*store.CreatedByStaffID)
		if err != nil {
			return fmt.Errorf("invalid created_by_staff_id: %w", err)
		}
		updates = append(updates, fmt.Sprintf("created_by_staff_id = $%d", argNum))
		args = append(args, staffID)
		argNum++
	}

	if store.DeactivatedByID != nil {
		deactivatedByID, err := uuid.Parse(*store.DeactivatedByID)
		if err != nil {
			return fmt.Errorf("invalid deactivated_by_id: %w", err)
		}
		updates = append(updates, fmt.Sprintf("deactivated_by_id = $%d", argNum))
		args = append(args, deactivatedByID)
		argNum++
	}

	if store.DeactivatedAt != nil {
		updates = append(updates, fmt.Sprintf("deactivated_at = $%d", argNum))
		args = append(args, *store.DeactivatedAt)
		argNum++
	}

	if store.DeactivationReason != nil {
		updates = append(updates, fmt.Sprintf("deactivation_reason = $%d", argNum))
		args = append(args, *store.DeactivationReason)
		argNum++
	}

	if len(updates) == 0 {
		return nil
	}

	sql := fmt.Sprintf(`
		UPDATE stores
		SET %s
		WHERE id = $%d`,
		utils.Join(updates, ", "),
		argNum,
	)

	args = append(args, id)

	ct, err := r.DB.Pool.Exec(ctx, sql, args...)
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

func (r *stores) List(ctx context.Context, isActive, status, sortBy, sortOrder string, lim, off int) ([]model.Store, int, error) {
	var conditions []string
	var args []interface{}
	argNum := 1

	if status != "" {
		conditions = append(conditions, fmt.Sprintf("status = $%d", argNum))
		args = append(args, status)
		argNum++
	}

	if isActive != "" {
		conditions = append(conditions, fmt.Sprintf("is_active = $%d", argNum))
		args = append(args, isActive)
		argNum++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + utils.Join(conditions, " AND ")
	}

	orderBy := "created_at DESC"
	if sortBy != "" {
		orderBy = sortBy
		if sortOrder == "ASC" {
			orderBy += " ASC"
		} else {
			orderBy += " DESC"
		}
	}

	limit := 5
	offset := 0
	if lim > 0 {
		limit = lim
	}
	if off > 0 {
		offset = off
	}

	baseSQL := `
		SELECT 
			id, name, address, latitude, longitude,
			owner_full_name, owner_inn, owner_passport,
			owner_phone, additional_phone, comment,
			status, rejection_comment,
			created_by_staff_id,
			is_active, is_frozen,
			created_at
		FROM stores %s
		ORDER BY %s
		LIMIT $%d OFFSET $%d`

	countSQL := `SELECT COUNT(*) FROM stores ` + whereClause

	var total int
	err := r.DB.Pool.QueryRow(ctx, countSQL, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("count stores: %w", err)
	}

	args = append(args, limit, offset)

	query := fmt.Sprintf(baseSQL, whereClause, orderBy, argNum, argNum+1)

	rows, err := r.DB.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("list stores: %w", err)
	}
	defer rows.Close()

	var storesList = make([]model.Store, 0)

	for rows.Next() {
		s := model.Store{}
		if err := rows.Scan(
			&s.ID, &s.Name, &s.Address, &s.Latitude, &s.Longitude,
			&s.OwnerFullName, &s.OwnerINN, &s.OwnerPassport,
			&s.OwnerPhone, &s.AdditionalPhone, &s.Comment,
			&s.Status, &s.RejectionComment,
			&s.CreatedByStaffID,
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

func (r *stores) IfExists(ctx context.Context, name string) (bool, error) {
	var s model.Store

	sql := `
		select id, owner_full_name from stores where name = $1
	`

	if err := r.DB.Pool.QueryRow(ctx, sql, name).Scan(&s.ID, &s.OwnerFullName); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("scan error: %w", err)
	}

	return true, nil
}

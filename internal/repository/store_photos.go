package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/model"
	"github.com/rakhshon-mirzoev/gps-tracker/pkg/db"
)

type StorePhotosRepository interface {
	Create(ctx context.Context, sp *model.StorePhotos) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.StorePhotos, error)
	List(ctx context.Context, storeID uuid.UUID) ([]model.StorePhotos, error)
	Delete(ctx context.Context, id uuid.UUID) error
	CreateTx(ctx context.Context, tx pgx.Tx, sp *model.StorePhotos) error
}

type storePhotos struct {
	DB *db.DB
}

func NewStorePhotosRepository(db *db.DB) StorePhotosRepository {
	return &storePhotos{DB: db}
}

func (s *storePhotos) CreateTx(ctx context.Context, tx pgx.Tx, sp *model.StorePhotos) error {
	sql := `
		INSERT INTO store_photos (store_id, object_key, type)
		VALUES ($1,$2,$3)
		RETURNING id, created_at
	`
	if err := tx.QueryRow(ctx, sql, sp.StoreID, sp.ObjectKey, sp.Type).Scan(&sp.ID, &sp.CreatedAt); err != nil {
		return fmt.Errorf("create store photo tx: %w", err)
	}
	return nil
}

func (r *storePhotos) Create(ctx context.Context, sp *model.StorePhotos) error {
	sql := `
		INSERT INTO store_photos (store_id, object_key, type)
		VALUES ($1, $2, $3)
	`

	_, err := r.DB.Pool.Exec(ctx, sql, sp.StoreID, sp.ObjectKey, sp.Type)
	if err != nil {
		return fmt.Errorf("add store photo: %w", err)
	}

	return nil
}

func (r *storePhotos) GetByID(ctx context.Context, id uuid.UUID) (*model.StorePhotos, error) {
	sql := `
        SELECT id, store_id, object_key, type, created_at
        FROM store_photos WHERE id = $1
    `

	var sp model.StorePhotos
	err := r.DB.Pool.QueryRow(ctx, sql, id).Scan(
		&sp.ID,
		&sp.StoreID,
		&sp.ObjectKey,
		&sp.Type,
		&sp.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("get by id store photos: %w", err)
	}

	return &sp, nil
}

func (r *storePhotos) List(ctx context.Context, storeID uuid.UUID) ([]model.StorePhotos, error) {
	sql := `
        SELECT id, store_id, object_key, type, created_at
        FROM store_photos WHERE store_id = $1 ORDER BY created_at
    `

	rows, err := r.DB.Pool.Query(ctx, sql, storeID)
	if err != nil {
		return nil, fmt.Errorf("list store photos: %w", err)
	}

	var photos = make([]model.StorePhotos, 0)

	for rows.Next() {
		var sp model.StorePhotos

		if err := rows.Scan(
			&sp.ID,
			&sp.StoreID,
			&sp.ObjectKey,
			&sp.Type,
			&sp.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan store photos: %w", err)
		}

		photos = append(photos, sp)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows store photos: %w", err)
	}

	return photos, nil
}

func (r *storePhotos) Delete(ctx context.Context, id uuid.UUID) error {
	sql := `
        DELETE FROM store_photos WHERE id = $1
    `

	_, err := r.DB.Pool.Exec(ctx, sql, id)
	if err != nil {
		return fmt.Errorf("delete store photos: %w", err)
	}

	return nil
}

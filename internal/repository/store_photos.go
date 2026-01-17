package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/model"
	"github.com/rakhshon-mirzoev/gps-tracker/pkg/db"
)

type StorePhotosRepository interface {
	Create(ctx context.Context, sp *model.StorePhotos) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.StorePhotos, error)
	List(ctx context.Context, storeID uuid.UUID) ([]model.StorePhotos, error)
	Update(ctx context.Context, sp *model.StorePhotos) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type storePhotos struct {
	DB *db.DB
}

func InitStorePhotosRepository(db *db.DB) StorePhotosRepository {
	return &storePhotos{DB: db}
}

func (r *storePhotos) Create(ctx context.Context, sp *model.StorePhotos) error {
	sql := `
        INSERT INTO store_photos (
            store_id,
            photo_url,
            type
        ) VALUES (
            $1, $2, $3
        ) RETURNING id, created_at;
    `

	err := r.DB.Pool.QueryRow(ctx, sql, sp.StoreID, sp.PhotoURL, sp.Type).Scan(&sp.ID, &sp.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (r *storePhotos) GetByID(ctx context.Context, id uuid.UUID) (*model.StorePhotos, error) {
	sql := `
        SELECT id, store_id, photo_url, type, created_at
        FROM store_photos WHERE id = $1
    `

	var sp model.StorePhotos
	err := r.DB.Pool.QueryRow(ctx, sql, id).Scan(
		&sp.ID,
		&sp.StoreID,
		&sp.PhotoURL,
		&sp.Type,
		&sp.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &sp, nil
}

func (r *storePhotos) List(ctx context.Context, storeID uuid.UUID) ([]model.StorePhotos, error) {
	sql := `
        SELECT id, store_id, photo_url, type, created_at
        FROM store_photos WHERE store_id = $1
    `

	rows, err := r.DB.Pool.Query(ctx, sql, storeID)
	if err != nil {
		return nil, err
	}

	var photos []model.StorePhotos

	for rows.Next() {
		var sp model.StorePhotos

		if err := rows.Scan(
			&sp.ID,
			&sp.StoreID,
			&sp.PhotoURL,
			&sp.Type,
			&sp.CreatedAt,
		); err != nil {
			return nil, err
		}

		photos = append(photos, sp)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return photos, nil
}

func (r *storePhotos) Update(ctx context.Context, sp *model.StorePhotos) error {
	sql := `
        UPDATE store_photos
        SET store_id = $1, photo_url = $2, type = $3
        WHERE id = $4
    `

	_, err := r.DB.Pool.Exec(ctx, sql, sp.StoreID, sp.PhotoURL, sp.Type, sp.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *storePhotos) Delete(ctx context.Context, id uuid.UUID) error {
	sql := `
        DELETE FROM store_photos WHERE id = $1
    `

	_, err := r.DB.Pool.Exec(ctx, sql, id)
	if err != nil {
		return err
	}

	return nil
}

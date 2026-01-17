package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/model"
	"github.com/rakhshon-mirzoev/gps-tracker/pkg/db"
)

type VisitPhotosRepository interface {
	Create(ctx context.Context, v *model.VisitPhotos) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.VisitPhotos, error)
	List(ctx context.Context, visitID uuid.UUID) ([]model.VisitPhotos, error)
	Update(ctx context.Context, v *model.VisitPhotos) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type visitPhotos struct {
	DB *db.DB
}

func InitVisitPhotosRepository(db *db.DB) VisitPhotosRepository {
	return &visitPhotos{DB: db}
}

func (r *visitPhotos) Create(ctx context.Context, v *model.VisitPhotos) error {
	sql := `
		INSERT INTO visit_photos (
			visit_id,
			photo_url,
			captured_at
		) VALUES (
			$1,
			$2,
			$3 
		) RETURNING id;
	`

	err := r.DB.Pool.QueryRow(ctx, sql, v.VisitID, v.PhotoURL, v.CapturedAt).Scan(&v.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *visitPhotos) List(ctx context.Context, id uuid.UUID) ([]model.VisitPhotos, error) {
	sql := `
		SELECT id, visit_id, photo_url, captured_at
		FROM visit_photos WHERE visit_id = $1
	`

	rows, err := r.DB.Pool.Query(ctx, sql, id)
	if err != nil {
		return nil, err
	}

	var v []model.VisitPhotos

	for rows.Next() {
		var vPh model.VisitPhotos

		if err := rows.Scan(
			&vPh.ID,
			&vPh.VisitID,
			&vPh.PhotoURL,
			&vPh.CapturedAt,
		); err != nil {
			return nil, err
		}

		v = append(v, vPh)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return v, nil
}

func (r *visitPhotos) GetByID(ctx context.Context, id uuid.UUID) (*model.VisitPhotos, error) {
	sql := `
        SELECT id, visit_id, photo_url, captured_at
        FROM visit_photos WHERE id = $1
    `

	var v model.VisitPhotos
	err := r.DB.Pool.QueryRow(ctx, sql, id).Scan(
		&v.ID,
		&v.VisitID,
		&v.PhotoURL,
		&v.CapturedAt,
	)
	if err != nil {
		return nil, err
	}

	return &v, nil
}

func (r *visitPhotos) Update(ctx context.Context, v *model.VisitPhotos) error {
	sql := `
        UPDATE visit_photos
        SET visit_id = $1, photo_url = $2, captured_at = $3
        WHERE id = $4
    `

	_, err := r.DB.Pool.Exec(ctx, sql, v.VisitID, v.PhotoURL, v.CapturedAt, v.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *visitPhotos) Delete(ctx context.Context, id uuid.UUID) error {
	sql := `
        DELETE FROM visit_photos WHERE id = $1
    `

	_, err := r.DB.Pool.Exec(ctx, sql, id)
	if err != nil {
		return err
	}

	return nil
}

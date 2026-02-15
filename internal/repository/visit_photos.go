package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/model"
	"github.com/rakhshon-mirzoev/gps-tracker/pkg/db"
)

type VisitPhotosRepository interface {
	Create(ctx context.Context, v *model.VisitPhotos) error
	GetByVisitID(ctx context.Context, id uuid.UUID) ([]model.VisitPhotos, error)
	List(ctx context.Context, visitID uuid.UUID) ([]model.VisitPhotos, error)
	Update(ctx context.Context, id uuid.UUID, dto model.UpdateVisitPhotos) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.VisitPhotos, error)
}

type visitPhotos struct {
	DB *db.DB
}

func NewVisitPhotosRepository(db *db.DB) VisitPhotosRepository {
	return &visitPhotos{DB: db}
}

func (r *visitPhotos) Create(ctx context.Context, v *model.VisitPhotos) error {
	sql := `
		INSERT INTO visit_photos (
			visit_id,
			photo_key,
			original_name,
			captured_at
		) VALUES (
			$1,
			$2,
			$3,
			$4
		) RETURNING id;
	`

	err := r.DB.Pool.QueryRow(ctx, sql, v.VisitID, v.PhotoKey, v.OriginalName, v.CapturedAt).Scan(&v.ID)
	if err != nil {
		return fmt.Errorf("create visit photos: %w", err)
	}

	return nil
}

func (r *visitPhotos) List(ctx context.Context, id uuid.UUID) ([]model.VisitPhotos, error) {
	sql := `
		SELECT id, visit_id, photo_key, original_name, captured_at
		FROM visit_photos WHERE visit_id = $1
	`

	rows, err := r.DB.Pool.Query(ctx, sql, id)
	if err != nil {
		return nil, fmt.Errorf("list visit photos: %w", err)
	}

	var v = make([]model.VisitPhotos, 0)

	for rows.Next() {
		var vPh model.VisitPhotos

		if err := rows.Scan(
			&vPh.ID,
			&vPh.VisitID,
			&vPh.PhotoKey,
			&vPh.OriginalName,
			&vPh.CapturedAt,
		); err != nil {
			return nil, fmt.Errorf("scan visit photos: %w", err)
		}

		v = append(v, vPh)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows visit photos: %w", err)
	}

	return v, nil
}

func (r *visitPhotos) GetByVisitID(ctx context.Context, id uuid.UUID) ([]model.VisitPhotos, error) {
	sql := `
        SELECT id, visit_id, photo_key, original_name, captured_at
        FROM visit_photos WHERE visit_id = $1
    `

	v := make([]model.VisitPhotos, 0)
	rows, err := r.DB.Pool.Query(ctx, sql, id)
	if err != nil {
		return nil, fmt.Errorf("get by id visit photos rows: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var vp model.VisitPhotos
		if err := rows.Scan(
			&vp.ID,
			&vp.VisitID,
			&vp.PhotoKey,
			&vp.OriginalName,
			&vp.CapturedAt,
		); err != nil {
			return nil, fmt.Errorf("rows scan: %w", err)
		}

		v = append(v, vp)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows.Err: %w", err)
	}

	return v, nil
}

func (r *visitPhotos) GetByID(ctx context.Context, id uuid.UUID) (*model.VisitPhotos, error) {
	sql := `
        SELECT id, visit_id, photo_key, original_name, captured_at
        FROM visit_photos WHERE id = $1
    `

	var v model.VisitPhotos
	if err := r.DB.Pool.QueryRow(ctx, sql, id).Scan(
		&v.ID,
		&v.VisitID,
		&v.PhotoKey,
		&v.OriginalName,
		&v.CapturedAt,
	); err != nil {
		return nil, fmt.Errorf("get by id visit photos: %w", err)
	}

	return &v, nil
}

func (r *visitPhotos) Update(ctx context.Context, id uuid.UUID, dto model.UpdateVisitPhotos) error {
	sql := `UPDATE visit_photos SET `
	var args []any
	n := 1

	if dto.VisitID != nil {
		sql += fmt.Sprintf("visit_id = $%d, ", n)
		args = append(args, *dto.VisitID)
		n++
	}

	if dto.PhotoURL != nil {
		sql += fmt.Sprintf("photo_url = $%d, ", n)
		args = append(args, *dto.PhotoURL)
		n++
	}

	if dto.CapturedAt != nil {
		sql += fmt.Sprintf("captured_at = $%d, ", n)
		args = append(args, *dto.CapturedAt)
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
		return fmt.Errorf("update visit photos: %w", err)
	}

	return nil
}

func (r *visitPhotos) Delete(ctx context.Context, id uuid.UUID) error {
	sql := `
        DELETE FROM visit_photos WHERE id = $1
    `

	_, err := r.DB.Pool.Exec(ctx, sql, id)
	if err != nil {
		return fmt.Errorf("delete visit photos: %w", err)
	}

	return nil
}

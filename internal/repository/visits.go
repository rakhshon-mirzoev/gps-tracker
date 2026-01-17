package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/model"
	"github.com/rakhshon-mirzoev/gps-tracker/pkg/db"
)

type VisitRepository interface {
	Create(ctx context.Context, v *model.Visit) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Visit, error)
	List(ctx context.Context, scheduleID uuid.UUID) ([]model.Visit, error)
	Update(ctx context.Context, v *model.Visit) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type visit struct {
	DB *db.DB
}

func InitVisitRepository(db *db.DB) VisitRepository {
	return &visit{DB: db}
}

func (r *visit) Create(ctx context.Context, v *model.Visit) error {
	sql := `
        INSERT INTO visits (
            schedule_id,
            staff_id,
            store_id,
            check_in_at,
            check_out_at,
            check_in_location,
            distance_deviation,
            selfie_url,
            comment,
            is_offline,
            review_status,
            review_comment,
            reviewer_id,
            reviewed_at
        ) VALUES (
            $1, $2, $3, $4, $5, ST_GeomFromText($6, 4326), $7, $8, $9, $10, $11, $12, $13, $14
        ) RETURNING id, created_at;
    `

	err := r.DB.Pool.QueryRow(ctx, sql,
		v.ScheduleID, v.StaffID, v.StoreID, v.CheckInAt, v.CheckOutAt,
		v.CheckInLocation, v.DistanceDeviation, v.SelfieURL, v.Comment,
		v.IsOffline, v.ReviewStatus, v.ReviewComment, v.ReviewerID, v.ReviewedAt,
	).Scan(&v.ID, &v.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (r *visit) GetByID(ctx context.Context, id uuid.UUID) (*model.Visit, error) {
	sql := `
        SELECT id, schedule_id, staff_id, store_id, check_in_at, check_out_at,
               ST_AsText(check_in_location) AS check_in_location, distance_deviation, selfie_url, comment,
               is_offline, review_status, review_comment, reviewer_id,
               reviewed_at, created_at
        FROM visits WHERE id = $1
    `

	var v model.Visit
	err := r.DB.Pool.QueryRow(ctx, sql, id).Scan(
		&v.ID, &v.ScheduleID, &v.StaffID, &v.StoreID, &v.CheckInAt, &v.CheckOutAt,
		&v.CheckInLocation, &v.DistanceDeviation, &v.SelfieURL, &v.Comment,
		&v.IsOffline, &v.ReviewStatus, &v.ReviewComment, &v.ReviewerID,
		&v.ReviewedAt, &v.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &v, nil
}

func (r *visit) List(ctx context.Context, scheduleID uuid.UUID) ([]model.Visit, error) {
	sql := `
        SELECT id, schedule_id, staff_id, store_id, check_in_at, check_out_at,
               ST_AsText(check_in_location) AS check_in_location, distance_deviation, selfie_url, comment,
               is_offline, review_status, review_comment, reviewer_id,
               reviewed_at, created_at
        FROM visits WHERE schedule_id = $1
    `

	rows, err := r.DB.Pool.Query(ctx, sql, scheduleID)
	if err != nil {
		return nil, err
	}

	var visits []model.Visit

	for rows.Next() {
		var v model.Visit

		if err := rows.Scan(
			&v.ID, &v.ScheduleID, &v.StaffID, &v.StoreID, &v.CheckInAt, &v.CheckOutAt,
			&v.CheckInLocation, &v.DistanceDeviation, &v.SelfieURL, &v.Comment,
			&v.IsOffline, &v.ReviewStatus, &v.ReviewComment, &v.ReviewerID,
			&v.ReviewedAt, &v.CreatedAt,
		); err != nil {
			return nil, err
		}

		visits = append(visits, v)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return visits, nil
}

func (r *visit) Update(ctx context.Context, v *model.Visit) error {
	sql := `
        UPDATE visits
        SET schedule_id = $1, staff_id = $2, store_id = $3, check_in_at = $4,
            check_out_at = $5, check_in_location = ST_GeomFromText($6, 4326), distance_deviation = $7,
            selfie_url = $8, comment = $9, is_offline = $10, review_status = $11,
            review_comment = $12, reviewer_id = $13, reviewed_at = $14
        WHERE id = $15
    `

	_, err := r.DB.Pool.Exec(ctx, sql,
		v.ScheduleID, v.StaffID, v.StoreID, v.CheckInAt, v.CheckOutAt,
		v.CheckInLocation, v.DistanceDeviation, v.SelfieURL, v.Comment,
		v.IsOffline, v.ReviewStatus, v.ReviewComment, v.ReviewerID, v.ReviewedAt, v.ID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *visit) Delete(ctx context.Context, id uuid.UUID) error {
	sql := `
        DELETE FROM visits WHERE id = $1
    `

	_, err := r.DB.Pool.Exec(ctx, sql, id)
	if err != nil {
		return err
	}

	return nil
}

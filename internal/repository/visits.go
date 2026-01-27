package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/model"
	"github.com/rakhshon-mirzoev/gps-tracker/pkg/db"
)

type VisitRepository interface {
	Create(ctx context.Context, v *model.Visit) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Visit, error)
	List(ctx context.Context, scheduleID uuid.UUID) ([]model.Visit, error)
	Update(ctx context.Context, id uuid.UUID, dto model.UpdateVisit) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type visit struct {
	DB *db.DB
}

func NewVisitRepository(db *db.DB) VisitRepository {
	return &visit{DB: db}
}

func (r *visit) Create(ctx context.Context, v *model.Visit) error {
	sql := `
        INSERT INTO visits (
            schedule_id,
            staff_id,
            store_id,
            check_in_at,
            check_in_lat,
			check_in_lon,
        ) VALUES (
            $1, $2, $3, $4, $5, $6
    	) RETURNING id, created_at;
    `

	err := r.DB.Pool.QueryRow(ctx, sql,
		v.ScheduleID, v.StaffID, v.StoreID, v.CheckInAt,
		v.CheckInLat, v.CheckInLon,
	).Scan(&v.ID, &v.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (r *visit) GetByID(ctx context.Context, id uuid.UUID) (*model.Visit, error) {
	sql := `
        SELECT id, schedule_id, staff_id, store_id, check_in_at, check_out_at, 
			   check_in_lat, check_in_lon,, distance_deviation, selfie_url, comment,
               is_offline, review_status, review_comment, reviewer_id,
               reviewed_at, created_at
        FROM visits WHERE id = $1
    `

	var v model.Visit
	err := r.DB.Pool.QueryRow(ctx, sql, id).Scan(
		&v.ID, &v.ScheduleID, &v.StaffID, &v.StoreID, &v.CheckInAt, &v.CheckOutAt,
		&v.CheckInLat, &v.CheckInLon, &v.DistanceDeviation, &v.SelfieURL, &v.Comment,
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
               check_in_lat, check_in_lon, distance_deviation, selfie_url, comment,
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
			&v.CheckInLat, &v.CheckInLon, &v.DistanceDeviation, &v.SelfieURL, &v.Comment,
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

func (r *visit) Update(ctx context.Context, id uuid.UUID, dto model.UpdateVisit) error {
	sql := `UPDATE visits SET `
	args := []interface{}{}
	n := 1

	if dto.ScheduleID != nil {
		sql += fmt.Sprintf("schedule_id = $%d, ", n)
		args = append(args, *dto.ScheduleID)
		n++
	}

	if dto.StaffID != nil {
		sql += fmt.Sprintf("staff_id = $%d, ", n)
		args = append(args, *dto.StaffID)
		n++
	}

	if dto.StoreID != nil {
		sql += fmt.Sprintf("store_id = $%d, ", n)
		args = append(args, *dto.StoreID)
		n++
	}

	if dto.CheckInAt != nil {
		sql += fmt.Sprintf("check_in_at = $%d, ", n)
		args = append(args, *dto.CheckInAt)
		n++
	}

	if dto.CheckOutAt != nil {
		sql += fmt.Sprintf("check_out_at = $%d, ", n)
		args = append(args, *dto.CheckOutAt)
		n++
	}

	if dto.CheckInLat != nil {
		sql += fmt.Sprintf("check_in_lat = $%d, ", n)
		args = append(args, *dto.CheckInLat)
		n++
	}

	if dto.CheckInLon != nil {
		sql += fmt.Sprintf("check_in_lon = $%d, ", n)
		args = append(args, *dto.CheckInLon)
		n++
	}

	if dto.DistanceDeviation != nil {
		sql += fmt.Sprintf("distance_deviation = $%d, ", n)
		args = append(args, *dto.DistanceDeviation)
		n++
	}

	if dto.SelfieURL != nil {
		sql += fmt.Sprintf("selfie_url = $%d, ", n)
		args = append(args, *dto.SelfieURL)
		n++
	}

	if dto.Comment != nil {
		sql += fmt.Sprintf("comment = $%d, ", n)
		args = append(args, *dto.Comment)
		n++
	}

	if dto.IsOffline != nil {
		sql += fmt.Sprintf("is_offline = $%d, ", n)
		args = append(args, *dto.IsOffline)
		n++
	}

	if dto.ReviewStatus != nil {
		sql += fmt.Sprintf("review_status = $%d, ", n)
		args = append(args, *dto.ReviewStatus)
		n++
	}

	if dto.ReviewComment != nil {
		sql += fmt.Sprintf("review_comment = $%d, ", n)
		args = append(args, *dto.ReviewComment)
		n++
	}

	if dto.ReviewerID != nil {
		sql += fmt.Sprintf("reviewer_id = $%d, ", n)
		args = append(args, *dto.ReviewerID)
		n++
	}

	if dto.ReviewedAt != nil {
		sql += fmt.Sprintf("reviewed_at = $%d, ", n)
		args = append(args, *dto.ReviewedAt)
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

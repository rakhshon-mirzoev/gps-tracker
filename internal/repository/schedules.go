package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/IkbolNazarov/gps-tracker/internal/model"
	"github.com/IkbolNazarov/gps-tracker/pkg/db"
	errs "github.com/IkbolNazarov/gps-tracker/pkg/error"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type schedules struct {
	DB *db.DB
}

type ScheduleRepository interface {
	Create(ctx context.Context, schedule *model.Schedule) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Schedule, error)
	Update(ctx context.Context, id uuid.UUID, dto model.UpdateSchedule) error
	Delete(ctx context.Context, id uuid.UUID) error
	ListByStaffID(ctx context.Context, staffID uuid.UUID, limit, offset int) ([]model.Schedule, int, error)
	ListAll(ctx context.Context, staffID, storeID *uuid.UUID, dateFrom, dateTo string, limit, offset int) ([]model.Schedule, int, error)
}

func NewScheduleRepository(db *db.DB) ScheduleRepository {
	return &schedules{DB: db}
}

func (r *schedules) Create(ctx context.Context, schedule *model.Schedule) error {
	sql := `
		INSERT INTO schedules (
			staff_id,
			store_id,
			planned_date,
			description,
			window_open,
			window_close,
			is_photo_required,
			is_selfie_required,
			is_comment_required,
			created_by_admin_id
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id `

	err := r.DB.Pool.QueryRow(
		ctx,
		sql,
		schedule.StaffID,
		schedule.StoreID,
		schedule.PlannedDate,
		schedule.Description,
		schedule.WindowOpen,
		schedule.WindowClose,
		schedule.IsPhotoRequired,
		schedule.IsSelfieRequired,
		schedule.IsCommentRequired,
		schedule.CreatedByAdminID,
	).Scan(&schedule.ID)

	if err != nil {
		return fmt.Errorf("create schedule: %w", err)
	}

	return nil
}

func (r *schedules) GetByID(ctx context.Context, id uuid.UUID) (*model.Schedule, error) {
	schedule := &model.Schedule{}

	sql := `
		SELECT 
			id, staff_id, store_id, planned_date, description,
			window_open, window_close,
			is_photo_required, is_selfie_required, is_comment_required,
			status, created_by_admin_id, review_comment
		FROM schedules 
		WHERE id = $1 AND is_deleted = false`

	err := r.DB.Pool.QueryRow(ctx, sql, id).Scan(
		&schedule.ID,
		&schedule.StaffID,
		&schedule.StoreID,
		&schedule.PlannedDate,
		&schedule.Description,
		&schedule.WindowOpen,
		&schedule.WindowClose,
		&schedule.IsPhotoRequired,
		&schedule.IsSelfieRequired,
		&schedule.IsCommentRequired,
		&schedule.Status,
		&schedule.CreatedByAdminID,
		&schedule.ReviewComment,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.ErrorNotFound
		}
		return nil, err
	}

	return schedule, nil
}

func (r *schedules) Update(ctx context.Context, id uuid.UUID, dto model.UpdateSchedule) error {
	sql := `UPDATE schedules SET `
	var args []any
	n := 1

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

	if dto.PlannedDate != nil {
		sql += fmt.Sprintf("planned_date = $%d, ", n)
		args = append(args, *dto.PlannedDate)
		n++
	}

	if dto.Description != nil {
		sql += fmt.Sprintf("description = $%d, ", n)
		args = append(args, *dto.Description)
		n++
	}

	if dto.WindowOpen != nil {
		sql += fmt.Sprintf("window_open = $%d, ", n)
		args = append(args, *dto.WindowOpen)
		n++
	}

	if dto.WindowClose != nil {
		sql += fmt.Sprintf("window_close = $%d, ", n)
		args = append(args, *dto.WindowClose)
		n++
	}

	if dto.IsPhotoRequired != nil {
		sql += fmt.Sprintf("is_photo_required = $%d, ", n)
		args = append(args, *dto.IsPhotoRequired)
		n++
	}

	if dto.IsSelfieRequired != nil {
		sql += fmt.Sprintf("is_selfie_required = $%d, ", n)
		args = append(args, *dto.IsSelfieRequired)
		n++
	}

	if dto.IsCommentRequired != nil {
		sql += fmt.Sprintf("is_comment_required = $%d, ", n)
		args = append(args, *dto.IsCommentRequired)
		n++
	}

	if dto.Status != nil {
		sql += fmt.Sprintf("status = $%d, ", n)
		args = append(args, *dto.Status)
		n++
	}

	if dto.ReviewComment != nil {
		sql += fmt.Sprintf("review_comment = $%d, ", n)
		args = append(args, *dto.ReviewComment)
		n++
	}

	if len(args) == 0 {
		return nil
	}

	sql = strings.TrimSuffix(sql, ", ")
	sql += fmt.Sprintf(" WHERE id = $%d", n)
	args = append(args, id)

	ct, err := r.DB.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("update schedule: %w", err)
	}

	if ct.RowsAffected() == 0 {
		return errs.ErrorNotFound
	}

	return nil
}

func (r *schedules) Delete(ctx context.Context, id uuid.UUID) error {
	sql := `UPDATE schedules SET is_deleted = true, deleted_at = now() WHERE id = $1 AND is_deleted = false`

	ct, err := r.DB.Pool.Exec(ctx, sql, id)
	if err != nil {
		return fmt.Errorf("delete schedule: %w", err)
	}

	if ct.RowsAffected() == 0 {
		return errs.ErrorNotFound
	}

	sqlVis := `UPDATE visits set review_status = 'schedule_deleted' WHERE schedule_id = $1`

	_, err = r.DB.Pool.Exec(ctx, sqlVis, id)
	if err != nil {
		return fmt.Errorf("delete schedule: %w", err)
	}

	return nil
}

func (r *schedules) ListByStaffID(ctx context.Context, staffID uuid.UUID, limit, offset int) ([]model.Schedule, int, error) {
	const defaultLimit = 5

	if limit <= 0 {
		limit = defaultLimit
	}
	if offset < 0 {
		offset = 0
	}

	var total int
	countSQL := `SELECT COUNT(*) FROM schedules WHERE staff_id = $1`
	err := r.DB.Pool.QueryRow(ctx, countSQL, staffID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("count schedules: %w", err)
	}

	sql := `
		SELECT 
			s.id, s.staff_id, s.store_id, st.name, s.planned_date, s.description,
			s.window_open, s.window_close,
			s.is_photo_required, s.is_selfie_required, s.is_comment_required,
			s.status, s.created_by_admin_id, s.review_comment,
			v.id
		FROM schedules s
		JOIN stores st ON st.id = s.store_id
		LEFT JOIN LATERAL (
			SELECT id
			FROM visits
			WHERE schedule_id = s.id
			ORDER BY created_at DESC
			LIMIT 1
		) v ON true
		WHERE s.staff_id = $1 AND s.is_deleted = false
		ORDER BY s.planned_date DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.DB.Pool.Query(ctx, sql, staffID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("list schedules by staff: %w", err)
	}
	defer rows.Close()

	var schedulesList = make([]model.Schedule, 0)
	for rows.Next() {
		s := model.Schedule{}
		if err := rows.Scan(
			&s.ID, &s.StaffID, &s.StoreID, &s.StoreName, &s.PlannedDate, &s.Description,
			&s.WindowOpen, &s.WindowClose,
			&s.IsPhotoRequired, &s.IsSelfieRequired, &s.IsCommentRequired,
			&s.Status, &s.CreatedByAdminID, &s.ReviewComment, &s.VisitId,
		); err != nil {
			return nil, 0, fmt.Errorf("scan schedule: %w", err)
		}
		schedulesList = append(schedulesList, s)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("rows error: %w", err)
	}

	return schedulesList, total, nil
}

func (r *schedules) ListAll(ctx context.Context, staffID, storeID *uuid.UUID, dateFrom, dateTo string, limit, offset int) ([]model.Schedule, int, error) {
	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	var conditions []string
	var args []any
	argNum := 1

	if staffID != nil {
		conditions = append(conditions, fmt.Sprintf("staff_id = $%d", argNum))
		args = append(args, *staffID)
		argNum++
	}

	if storeID != nil {
		conditions = append(conditions, fmt.Sprintf("store_id = $%d", argNum))
		args = append(args, *storeID)
		argNum++
	}

	if dateFrom != "" {
		conditions = append(conditions, fmt.Sprintf("planned_date >= $%d", argNum))
		args = append(args, dateFrom)
		argNum++
	}

	if dateTo != "" {
		conditions = append(conditions, fmt.Sprintf("planned_date <= $%d", argNum))
		args = append(args, dateTo)
		argNum++
	}

	conditions = append(conditions, "is_deleted = false")

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	var total int
	countSQL := fmt.Sprintf(`SELECT COUNT(*) FROM schedules %s`, whereClause)
	if err := r.DB.Pool.QueryRow(ctx, countSQL, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count schedules: %w", err)
	}

	sql := fmt.Sprintf(`
		SELECT 
			s.id, s.staff_id, s.store_id, st.name, s.planned_date, s.description,
			s.window_open, s.window_close,
			s.is_photo_required, s.is_selfie_required, s.is_comment_required,
			s.status, s.created_by_admin_id, s.review_comment, s.is_deleted, s.deleted_at,
			v.id
		FROM schedules s
		JOIN stores st ON st.id = s.store_id
		LEFT JOIN LATERAL (
			SELECT id
			FROM visits
			WHERE schedule_id = s.id
			ORDER BY created_at DESC
			LIMIT 1
		) v ON true
		%s
		ORDER BY planned_date DESC
		LIMIT $%d OFFSET $%d`, whereClause, argNum, argNum+1)

	args = append(args, limit, offset)

	rows, err := r.DB.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("list all schedules: %w", err)
	}
	defer rows.Close()

	var schedulesList = make([]model.Schedule, 0)
	for rows.Next() {
		s := model.Schedule{}
		if err := rows.Scan(
			&s.ID, &s.StaffID, &s.StoreID, &s.StoreName, &s.PlannedDate, &s.Description,
			&s.WindowOpen, &s.WindowClose,
			&s.IsPhotoRequired, &s.IsSelfieRequired, &s.IsCommentRequired,
			&s.Status, &s.CreatedByAdminID, &s.ReviewComment, &s.IsDeleted, &s.DeletedAt, &s.VisitId,
		); err != nil {
			return nil, 0, fmt.Errorf("scan schedule: %w", err)
		}
		schedulesList = append(schedulesList, s)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("rows error: %w", err)
	}

	return schedulesList, total, nil
}

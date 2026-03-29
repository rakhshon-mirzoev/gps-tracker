package repository

import (
	"context"
	"fmt"

	"github.com/IkbolNazarov/gps-tracker/internal/model"
	"github.com/IkbolNazarov/gps-tracker/pkg/db"
	errs "github.com/IkbolNazarov/gps-tracker/pkg/error"
	"github.com/IkbolNazarov/gps-tracker/pkg/utils"
	"github.com/google/uuid"
)

type StaffScheduleRepository interface {
	CreateBatch(ctx context.Context, staffID uuid.UUID, schedules []model.ScheduleDTO) error
	List(ctx context.Context, staffID uuid.UUID) ([]model.StaffUserSchedule, error)
	Delete(ctx context.Context, schID uuid.UUID) error
	Update(ctx context.Context, schID uuid.UUID, dto model.UpdateStaffSchedule) error
}

type staffScheduleRepo struct {
	DB *db.DB
}

func NewStaffScheduleRepository(db *db.DB) StaffScheduleRepository {
	return &staffScheduleRepo{DB: db}
}

func (r *staffScheduleRepo) CreateBatch(ctx context.Context, staffID uuid.UUID, schedules []model.ScheduleDTO) error {
	tx, err := r.DB.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	query := `
		INSERT INTO staff_user_schedules (staff_id, weekday, start_time, end_time)
		VALUES ($1, $2, $3, $4)
	`

	for i := range schedules {
		if schedules[i].StartTime.After(schedules[i].EndTime) {
			return errs.ErrorUnprocessable
		}

		_, err := tx.Exec(ctx, query, staffID, schedules[i].Weekday, schedules[i].StartTime, schedules[i].EndTime)
		if err != nil {
			return fmt.Errorf("insert schedule[%d]: %w", i, err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}

	return nil
}

func (r *staffScheduleRepo) List(ctx context.Context, staffID uuid.UUID) ([]model.StaffUserSchedule, error) {
	sql := `
		SELECT id, weekday, start_time, end_time, created_at, is_deleted, deleted_at
		FROM staff_user_schedules
		WHERE staff_id = $1 ORDER BY weekday
	`

	rows, err := r.DB.Pool.Query(ctx, sql, staffID)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	defer rows.Close()

	stSch := make([]model.StaffUserSchedule, 0)

	for rows.Next() {
		var sch model.StaffUserSchedule

		if err := rows.Scan(
			&sch.ID,
			&sch.Weekday,
			&sch.StartTime,
			&sch.EndTime,
			&sch.CreatedAt,
			&sch.IsDeleted,
			&sch.DeletedAt,
		); err != nil {
			return nil, fmt.Errorf("rows scan: %w", err)
		}

		stSch = append(stSch, sch)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows next: %w", err)
	}

	return stSch, nil
}

func (r *staffScheduleRepo) Delete(ctx context.Context, schID uuid.UUID) error {
	sql := `
		UPDATE staff_user_schedules SET is_deleted = true, deleted_at = now() WHERE id = $1 AND is_deleted = false
	`

	ct, err := r.DB.Pool.Exec(ctx, sql, schID)
	if err != nil {
		return fmt.Errorf("exec delete: %w", err)
	}

	if ct.RowsAffected() == 0 {
		return errs.ErrorNotFound
	}

	return nil
}

func (r *staffScheduleRepo) Update(ctx context.Context, schID uuid.UUID, dto model.UpdateStaffSchedule) error {
	var (
		query  []string
		args   []any
		argNum = 1
	)

	if dto.Weekday != nil {
		query = append(query, fmt.Sprintf("weekday = $%d, ", argNum))
		args = append(args, *dto.Weekday)
		argNum++
	}

	if dto.StartTime != nil {
		query = append(query, fmt.Sprintf("start_time = $%d, ", argNum))
		args = append(args, *dto.StartTime)
		argNum++
	}

	if dto.EndTime != nil {
		query = append(query, fmt.Sprintf("end_time = $%d, ", argNum))
		args = append(args, *dto.EndTime)
		argNum++
	}

	sql := fmt.Sprintf(`
		UPDATE staff_user_schedules
		SET %s
		WHERE id = $%d`,
		utils.Join(query, ", "),
		argNum,
	)

	args = append(args, schID)

	ct, err := r.DB.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("update staff_user_schedules: %w", err)
	}

	if ct.RowsAffected() == 0 {
		return errs.ErrorNotFound
	}

	return nil

}

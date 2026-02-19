package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/model"
	"github.com/rakhshon-mirzoev/gps-tracker/pkg/db"
	errs "github.com/rakhshon-mirzoev/gps-tracker/pkg/error"
)

type StaffScheduleRepository interface {
	CreateBatch(ctx context.Context, staffID uuid.UUID, schedules []model.ScheduleDTO) error
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

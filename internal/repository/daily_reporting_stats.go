package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/model"
	"github.com/rakhshon-mirzoev/gps-tracker/pkg/db"
)

type DailyReportingStatsRepository interface {
	Create(ctx context.Context, drs *model.DailyReportingStats) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.DailyReportingStats, error)
	List(ctx context.Context, staffID uuid.UUID) ([]model.DailyReportingStats, error)
	Update(ctx context.Context, id uuid.UUID, dto model.UpdateDailyReportingStats) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type dailyReportingStats struct {
	DB *db.DB
}

func InitDailyReportingStatsRepository(db *db.DB) DailyReportingStatsRepository {
	return &dailyReportingStats{DB: db}
}

func (r *dailyReportingStats) Create(ctx context.Context, drs *model.DailyReportingStats) error {
	sql := `
        INSERT INTO daily_reporting_stats (
            date,
            staff_id,
            group_id,
            planned_count,
            completed_count,
            rejected_count,
            missed_count,
            work_start_at,
            work_end_at,
            total_km
        ) VALUES (
            $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
        ) RETURNING id;
    `

	err := r.DB.Pool.QueryRow(ctx, sql,
		drs.Date, drs.StaffID, drs.GroupID, drs.PlannedCount, drs.CompletedCount,
		drs.RejectedCount, drs.MissedCount, drs.WorkStartAt, drs.WorkEndAt, drs.TotalKm,
	).Scan(&drs.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *dailyReportingStats) GetByID(ctx context.Context, id uuid.UUID) (*model.DailyReportingStats, error) {
	sql := `
        SELECT id, date, staff_id, group_id, planned_count, completed_count,
               rejected_count, missed_count, work_start_at, work_end_at, total_km
        FROM daily_reporting_stats WHERE id = $1
    `

	var drs model.DailyReportingStats
	err := r.DB.Pool.QueryRow(ctx, sql, id).Scan(
		&drs.ID, &drs.Date, &drs.StaffID, &drs.GroupID, &drs.PlannedCount, &drs.CompletedCount,
		&drs.RejectedCount, &drs.MissedCount, &drs.WorkStartAt, &drs.WorkEndAt, &drs.TotalKm,
	)
	if err != nil {
		return nil, err
	}

	return &drs, nil
}

func (r *dailyReportingStats) List(ctx context.Context, staffID uuid.UUID) ([]model.DailyReportingStats, error) {
	sql := `
        SELECT id, date, staff_id, group_id, planned_count, completed_count,
               rejected_count, missed_count, work_start_at, work_end_at, total_km
        FROM daily_reporting_stats WHERE staff_id = $1
    `

	rows, err := r.DB.Pool.Query(ctx, sql, staffID)
	if err != nil {
		return nil, err
	}

	var stats = make([]model.DailyReportingStats, 0)

	for rows.Next() {
		var drs model.DailyReportingStats

		if err := rows.Scan(
			&drs.ID, &drs.Date, &drs.StaffID, &drs.GroupID, &drs.PlannedCount, &drs.CompletedCount,
			&drs.RejectedCount, &drs.MissedCount, &drs.WorkStartAt, &drs.WorkEndAt, &drs.TotalKm,
		); err != nil {
			return nil, err
		}

		stats = append(stats, drs)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return stats, nil
}

func (r *dailyReportingStats) Update(ctx context.Context, id uuid.UUID, dto model.UpdateDailyReportingStats) error {
	sql := `UPDATE daily_reporting_stats SET `
	args := []interface{}{}
	n := 1

	if dto.Date != nil {
		sql += fmt.Sprintf("date = $%d, ", n)
		args = append(args, *dto.Date)
		n++
	}

	if dto.StaffID != nil {
		sql += fmt.Sprintf("staff_id = $%d, ", n)
		args = append(args, *dto.StaffID)
		n++
	}

	if dto.GroupID != nil {
		sql += fmt.Sprintf("group_id = $%d, ", n)
		args = append(args, *dto.GroupID)
		n++
	}

	if dto.PlannedCount != nil {
		sql += fmt.Sprintf("planned_count = $%d, ", n)
		args = append(args, *dto.PlannedCount)
		n++
	}

	if dto.CompletedCount != nil {
		sql += fmt.Sprintf("completed_count = $%d, ", n)
		args = append(args, *dto.CompletedCount)
		n++
	}

	if dto.RejectedCount != nil {
		sql += fmt.Sprintf("rejected_count = $%d, ", n)
		args = append(args, *dto.RejectedCount)
		n++
	}

	if dto.MissedCount != nil {
		sql += fmt.Sprintf("missed_count = $%d, ", n)
		args = append(args, *dto.MissedCount)
		n++
	}

	if dto.WorkStartAt != nil {
		sql += fmt.Sprintf("work_start_at = $%d, ", n)
		args = append(args, *dto.WorkStartAt)
		n++
	}

	if dto.WorkEndAt != nil {
		sql += fmt.Sprintf("work_end_at = $%d, ", n)
		args = append(args, *dto.WorkEndAt)
		n++
	}

	if dto.TotalKm != nil {
		sql += fmt.Sprintf("total_km = $%d, ", n)
		args = append(args, *dto.TotalKm)
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

func (r *dailyReportingStats) Delete(ctx context.Context, id uuid.UUID) error {
	sql := `
        DELETE FROM daily_reporting_stats WHERE id = $1
    `

	_, err := r.DB.Pool.Exec(ctx, sql, id)
	if err != nil {
		return err
	}

	return nil
}

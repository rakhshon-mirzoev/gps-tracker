package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/model"
	"github.com/rakhshon-mirzoev/gps-tracker/pkg/db"
)

type DailyReportingStatsRepository interface {
	Create(ctx context.Context, drs *model.DailyReportingStats) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.DailyReportingStats, error)
	List(ctx context.Context, staffID uuid.UUID) ([]model.DailyReportingStats, error)
	Update(ctx context.Context, drs *model.DailyReportingStats) error
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

	var stats []model.DailyReportingStats

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

func (r *dailyReportingStats) Update(ctx context.Context, drs *model.DailyReportingStats) error {
	sql := `
        UPDATE daily_reporting_stats
        SET date = $1, staff_id = $2, group_id = $3, planned_count = $4,
            completed_count = $5, rejected_count = $6, missed_count = $7,
            work_start_at = $8, work_end_at = $9, total_km = $10
        WHERE id = $11
    `

	_, err := r.DB.Pool.Exec(ctx, sql,
		drs.Date, drs.StaffID, drs.GroupID, drs.PlannedCount, drs.CompletedCount,
		drs.RejectedCount, drs.MissedCount, drs.WorkStartAt, drs.WorkEndAt, drs.TotalKm, drs.ID,
	)
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

package service

import (
	"context"
	"errors"

	"github.com/IkbolNazarov/gps-tracker/internal/model"
	"github.com/IkbolNazarov/gps-tracker/internal/repository"
	"github.com/google/uuid"
)

type DailyReportingStats interface {
	Create(ctx context.Context, drs *model.DailyReportingStats) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.DailyReportingStats, error)
	List(ctx context.Context, staffID uuid.UUID) ([]model.DailyReportingStats, error)
	Update(ctx context.Context, id uuid.UUID, dto model.UpdateDailyReportingStats) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type dailyReportingStatsService struct {
	repo repository.DailyReportingStatsRepository
}

func NewDailyReportingStatsService(repo repository.DailyReportingStatsRepository) DailyReportingStats {
	return &dailyReportingStatsService{repo: repo}
}

func (s *dailyReportingStatsService) Create(ctx context.Context, drs *model.DailyReportingStats) error {
	if drs.StaffID == uuid.Nil {
		return errors.New("staff_id is required")
	}
	if drs.GroupID == uuid.Nil {
		return errors.New("group_id is required")
	}
	if drs.PlannedCount < 0 {
		return errors.New("planned_count must be non-negative")
	}
	if drs.CompletedCount < 0 {
		return errors.New("completed_count must be non-negative")
	}
	if drs.RejectedCount < 0 {
		return errors.New("rejected_count must be non-negative")
	}
	if drs.MissedCount < 0 {
		return errors.New("missed_count must be non-negative")
	}

	return s.repo.Create(ctx, drs)
}

func (s *dailyReportingStatsService) GetByID(ctx context.Context, id uuid.UUID) (*model.DailyReportingStats, error) {
	if id == uuid.Nil {
		return nil, errors.New("id is required")
	}

	return s.repo.GetByID(ctx, id)
}

func (s *dailyReportingStatsService) List(ctx context.Context, staffID uuid.UUID) ([]model.DailyReportingStats, error) {
	if staffID == uuid.Nil {
		return nil, errors.New("staff_id is required")
	}

	return s.repo.List(ctx, staffID)
}

func (s *dailyReportingStatsService) Update(ctx context.Context, id uuid.UUID, dto model.UpdateDailyReportingStats) error {
	return s.repo.Update(ctx, id, dto)
}

func (s *dailyReportingStatsService) Delete(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("id is required")
	}

	return s.repo.Delete(ctx, id)
}

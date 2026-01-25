package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/model"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/repository"
)

type Schedule interface {
	Create(ctx context.Context, schedule *model.Schedule) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Schedule, error)
	Update(ctx context.Context, id uuid.UUID, dto model.UpdateSchedule) error
	Delete(ctx context.Context, id uuid.UUID) error
	ListByStaffID(ctx context.Context, staffID uuid.UUID, limit, offset int) ([]*model.Schedule, int, error)
	ListByStoreID(ctx context.Context, storeID uuid.UUID, limit, offset int) ([]*model.Schedule, int, error)
	ListByDate(ctx context.Context, plannedDate string, limit, offset int) ([]*model.Schedule, int, error)
	ListByStatus(ctx context.Context, status string, limit, offset int) ([]*model.Schedule, int, error)
	GetProfile(ctx context.Context, id uuid.UUID) (*model.StaffUserProfile, error)
}

type scheduleService struct {
	r     repository.ScheduleRepository
	staff StaffUsers
}

func NewScheduleService(r repository.ScheduleRepository) Schedule {
	return &scheduleService{r: r}
}

func (s *scheduleService) Create(ctx context.Context, schedule *model.Schedule) error {
	return s.r.Create(ctx, schedule)
}

func (s *scheduleService) GetProfile(ctx context.Context, id uuid.UUID) (*model.StaffUserProfile, error) {
	if id == uuid.Nil {
		return nil, fmt.Errorf("id is required")
	}

	user, err := s.staff.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	schedules, total, err := s.r.ListByStaffID(ctx, id, 50, 0)
	if err != nil {
		return nil, fmt.Errorf("get schedules: %w", err)
	}

	return &model.StaffUserProfile{
		User:       user,
		Schedules:  schedules,
		TotalCount: total,
	}, nil
}

func (s *scheduleService) GetByID(ctx context.Context, id uuid.UUID) (*model.Schedule, error) {
	if id == uuid.Nil {
		return nil, fmt.Errorf("id is required")
	}

	return s.r.GetByID(ctx, id)
}

func (s *scheduleService) Update(ctx context.Context, id uuid.UUID, dto model.UpdateSchedule) error {
	if id == uuid.Nil {
		return fmt.Errorf("id is required")
	}

	_, err := s.r.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("schedule not found: %w", err)
	}

	return s.r.Update(ctx, id, dto)
}

func (s *scheduleService) Delete(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return fmt.Errorf("id is required")
	}

	return s.r.Delete(ctx, id)
}

func (s *scheduleService) ListByStaffID(ctx context.Context, staffID uuid.UUID, limit, offset int) ([]*model.Schedule, int, error) {
	if staffID == uuid.Nil {
		return nil, 0, fmt.Errorf("staff_id is required")
	}

	return s.r.ListByStaffID(ctx, staffID, limit, offset)
}

func (s *scheduleService) ListByStoreID(ctx context.Context, storeID uuid.UUID, limit, offset int) ([]*model.Schedule, int, error) {
	if storeID == uuid.Nil {
		return nil, 0, fmt.Errorf("store_id is required")
	}

	return s.r.ListByStoreID(ctx, storeID, limit, offset)
}

func (s *scheduleService) ListByDate(ctx context.Context, plannedDate string, limit, offset int) ([]*model.Schedule, int, error) {
	if plannedDate == "" {
		return nil, 0, fmt.Errorf("planned_date is required")
	}

	return s.r.ListByDate(ctx, plannedDate, limit, offset)
}

func (s *scheduleService) ListByStatus(ctx context.Context, status string, limit, offset int) ([]*model.Schedule, int, error) {
	if status == "" {
		return nil, 0, fmt.Errorf("status is required")
	}

	return s.r.ListByStatus(ctx, status, limit, offset)
}

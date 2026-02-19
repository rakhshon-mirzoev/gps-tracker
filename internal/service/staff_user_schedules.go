package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/model"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/repository"
)

type StaffSchedule interface {
	CreateBatch(ctx context.Context, staffID uuid.UUID, schedule []model.ScheduleDTO) error
	List(ctx context.Context, staffID uuid.UUID) ([]model.StaffUserSchedule, error)
	Delete(ctx context.Context, schID uuid.UUID) error
	Update(ctx context.Context, schID uuid.UUID, dto model.UpdateStaffSchedule) error
}

type staffScheduleService struct {
	r repository.StaffScheduleRepository
}

func NewStaffScheduleService(r repository.StaffScheduleRepository) StaffSchedule {
	return &staffScheduleService{r: r}
}

func (s *staffScheduleService) CreateBatch(ctx context.Context, staffID uuid.UUID, schedules []model.ScheduleDTO) error {
	return s.r.CreateBatch(ctx, staffID, schedules)
}

func (s *staffScheduleService) List(ctx context.Context, staffID uuid.UUID) ([]model.StaffUserSchedule, error) {
	return s.r.List(ctx, staffID)
}
func (s *staffScheduleService) Delete(ctx context.Context, schID uuid.UUID) error {
	return s.r.Delete(ctx, schID)
}
func (s *staffScheduleService) Update(ctx context.Context, schID uuid.UUID, dto model.UpdateStaffSchedule) error {
	return s.r.Update(ctx, schID, dto)
}

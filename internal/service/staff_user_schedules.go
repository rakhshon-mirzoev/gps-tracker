package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/model"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/repository"
)

type StaffSchedule interface {
	CreateBatch(ctx context.Context, staffID uuid.UUID, schedule []model.ScheduleDTO) error
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

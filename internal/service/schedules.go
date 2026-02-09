package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/model"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/repository"
)

type Schedule interface {
	Create(ctx context.Context, dto *model.CreateSchedule, adminID uuid.UUID, ip, schedule, role, action string) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Schedule, error)
	Update(ctx context.Context, id uuid.UUID, dto model.UpdateSchedule) error
	Delete(ctx context.Context, id uuid.UUID) error
	ListByStaffID(ctx context.Context, staffID uuid.UUID, limit, offset int) ([]model.Schedule, int, error)
	ListByStoreID(ctx context.Context, storeID uuid.UUID, limit, offset int) ([]model.Schedule, int, error)
	ListByDate(ctx context.Context, plannedDate string, limit, offset int) ([]model.Schedule, int, error)
	ListByStatus(ctx context.Context, status string, limit, offset int) ([]model.Schedule, int, error)
}

type scheduleService struct {
	r     repository.ScheduleRepository
	audit AuditLogs
}

func NewScheduleService(r repository.ScheduleRepository, audit AuditLogs) Schedule {
	return &scheduleService{r: r, audit: audit}
}

func (s *scheduleService) Create(ctx context.Context, dto *model.CreateSchedule, adminID uuid.UUID, ip, scheduleConst, role, action string) error {
	staffID, err := uuid.Parse(dto.StaffID)
	if err != nil {
		return err
	}

	storeID, err := uuid.Parse(dto.StoreID)
	if err != nil {
		return err
	}

	schedule := &model.Schedule{
		StaffID:           staffID,
		StoreID:           storeID,
		PlannedDate:       dto.PlannedDate,
		Description:       dto.Description,
		WindowOpen:        dto.WindowOpen,
		WindowClose:       dto.WindowClose,
		IsPhotoRequired:   dto.IsPhotoRequired,
		IsSelfieRequired:  dto.IsSelfieRequired,
		IsCommentRequired: dto.IsCommentRequired,
		CreatedByAdminID:  adminID,
	}

	if err := s.r.Create(ctx, schedule); err != nil {
		return err
	}

	if err := s.audit.Create(ctx, &model.AuditLogs{
		ActorID:     adminID,
		ActorType:   role,
		Action:      "AssignSchedule",
		TargetTable: &scheduleConst,
		IPAddress:   &ip,
	}); err != nil {
		return err
	}

	return nil
}

func (s *scheduleService) GetByID(ctx context.Context, id uuid.UUID) (*model.Schedule, error) {
	return s.r.GetByID(ctx, id)
}

func (s *scheduleService) Update(ctx context.Context, id uuid.UUID, dto model.UpdateSchedule) error {
	_, err := s.r.GetByID(ctx, id)
	if err != nil {
		return err
	}

	updateDto := dto
	if dto.StaffID != nil {
		staffID, err := uuid.Parse(*dto.StaffID)
		if err != nil {
			return err
		}
		staffIDStr := staffID.String()
		updateDto.StaffID = &staffIDStr
	}

	if dto.StoreID != nil {
		storeID, err := uuid.Parse(*dto.StoreID)
		if err != nil {
			return err
		}
		storeIDStr := storeID.String()
		updateDto.StoreID = &storeIDStr
	}
	return s.r.Update(ctx, id, updateDto)
}

func (s *scheduleService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.r.Delete(ctx, id)
}

func (s *scheduleService) ListByStaffID(ctx context.Context, staffID uuid.UUID, limit, offset int) ([]model.Schedule, int, error) {
	return s.r.ListByStaffID(ctx, staffID, limit, offset)
}

func (s *scheduleService) ListByStoreID(ctx context.Context, storeID uuid.UUID, limit, offset int) ([]model.Schedule, int, error) {
	return s.r.ListByStoreID(ctx, storeID, limit, offset)
}

func (s *scheduleService) ListByDate(ctx context.Context, plannedDate string, limit, offset int) ([]model.Schedule, int, error) {
	return s.r.ListByDate(ctx, plannedDate, limit, offset)
}

func (s *scheduleService) ListByStatus(ctx context.Context, status string, limit, offset int) ([]model.Schedule, int, error) {
	return s.r.ListByStatus(ctx, status, limit, offset)
}

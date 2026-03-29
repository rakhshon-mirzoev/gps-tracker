package service

import (
	"context"
	"fmt"
	"time"

	"github.com/IkbolNazarov/gps-tracker/internal/model"
	"github.com/IkbolNazarov/gps-tracker/internal/repository"
	"github.com/IkbolNazarov/gps-tracker/pkg/timezone"
	"github.com/google/uuid"
)

type Schedule interface {
	Create(ctx context.Context, dto *model.CreateSchedule, adminID uuid.UUID, ip, schedule, role, action string) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Schedule, error)
	Update(ctx context.Context, id uuid.UUID, dto model.UpdateSchedule) error
	Delete(ctx context.Context, id uuid.UUID) error
	ListByStaffID(ctx context.Context, staffID uuid.UUID, limit, offset int) ([]model.Schedule, int, error)
	ListAll(ctx context.Context, staffIDStr, storeIDStr, dateFrom, dateTo string, limit, offset int) ([]model.Schedule, int, error)
}

type scheduleService struct {
	r            repository.ScheduleRepository
	audit        AuditLogs
	notification Notification
	staff        repository.StaffUsersRepository
	store        Stores
}

func NewScheduleService(r repository.ScheduleRepository, audit AuditLogs, notification Notification, staff repository.StaffUsersRepository, store Stores) Schedule {
	return &scheduleService{r: r, audit: audit, notification: notification, staff: staff, store: store}
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

	staff, err := s.staff.GetByID(ctx, staffID)
	if err != nil {
		return err
	}

	store, err := s.store.GetByID(ctx, storeID)
	if err != nil {
		return err
	}

	if staff.FcmToken != nil {
		plannedDate := dto.PlannedDate.Format("02.01.2006, 15:04")

		open := dto.WindowOpen.Format("02.01.2006, 15:04")
		close := dto.WindowClose.Format("02.01.2006, 15:04")

		body := fmt.Sprintf(
			"Магазин: %s\nПлановая дата: %s\nОписание: %s\nОкно визита: %s-%s",
			store.Name,
			plannedDate,
			dto.Description,
			open,
			close,
		)
		if err := s.notification.CreateAndPushToUser(ctx, staff.ID, "Назначена новая задача", body, "schedule_assign", map[string]any{"type": "schedule_assign"}, staff.FcmToken); err != nil {
			return err
		}
	}

	if err := s.audit.Create(ctx, &model.AuditLogs{
		ActorID:     adminID,
		ActorType:   role,
		Action:      "AssignSchedule",
		TargetTable: &scheduleConst,
		IPAddress:   &ip,
		CreatedAt:   time.Now().In(timezone.Dushanbe),
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

func (s *scheduleService) ListAll(ctx context.Context, staffIDStr, storeIDStr, dateFrom, dateTo string, limit, offset int) ([]model.Schedule, int, error) {
	var staffID, storeID *uuid.UUID

	if staffIDStr != "" {
		id, err := uuid.Parse(staffIDStr)
		if err != nil {
			return nil, 0, err
		}
		staffID = &id
	}

	if storeIDStr != "" {
		id, err := uuid.Parse(storeIDStr)
		if err != nil {
			return nil, 0, err
		}
		storeID = &id
	}

	return s.r.ListAll(ctx, staffID, storeID, dateFrom, dateTo, limit, offset)
}

package service

import (
	"context"
	"fmt"
	"time"

	"github.com/IkbolNazarov/gps-tracker/internal/model"
	"github.com/IkbolNazarov/gps-tracker/internal/repository"
	errs "github.com/IkbolNazarov/gps-tracker/pkg/error"
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
	CreateBulk(ctx context.Context, dto *model.CreateSchedulesBulk, adminID uuid.UUID, ip, scheduleConst, role string) (int, error)
	InitSchedule(ctx context.Context, dto *model.CreateSchedule) error
}

type scheduleService struct {
	r               repository.ScheduleRepository
	audit           AuditLogs
	notification    Notification
	sysNotification SystemNotification
	staff           repository.StaffUsersRepository
	store           Stores
	group           Group
}

func NewScheduleService(r repository.ScheduleRepository, audit AuditLogs, notification Notification, sysNotif SystemNotification, staff repository.StaffUsersRepository, store Stores, group Group) Schedule {
	return &scheduleService{r: r, audit: audit, notification: notification, sysNotification: sysNotif, staff: staff, store: store, group: group}
}

func (s *scheduleService) Create(ctx context.Context, dto *model.CreateSchedule, adminID uuid.UUID, ip, scheduleConst, role, action string) error {
	schedule := &model.Schedule{
		StaffID:           dto.StaffID,
		StoreID:           dto.StoreID,
		PlannedDate:       dto.PlannedDate,
		Description:       dto.Description,
		WindowOpen:        dto.WindowOpen,
		WindowClose:       dto.WindowClose,
		IsPhotoRequired:   &dto.IsPhotoRequired,
		IsSelfieRequired:  &dto.IsSelfieRequired,
		IsCommentRequired: &dto.IsCommentRequired,
		CreatedByAdminID:  adminID,
	}

	if err := s.r.Create(ctx, schedule); err != nil {
		return err
	}

	staff, err := s.staff.GetByID(ctx, dto.StaffID)
	if err != nil {
		return err
	}

	store, err := s.store.GetByID(ctx, dto.StoreID)
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

func (s *scheduleService) InitSchedule(ctx context.Context, dto *model.CreateSchedule) error {
	schedule := model.Schedule{
		StaffID:          dto.StaffID,
		StoreID:          dto.StoreID,
		PlannedDate:      dto.PlannedDate,
		WindowOpen:       dto.WindowOpen,
		WindowClose:      dto.WindowClose,
		Description:      dto.Description,
		Status:           "initialized",
		CreatedByAdminID: dto.StaffID,
		IsUnplanned:      true,
		CreatedAtLat:     dto.Lat,
		CreatedAtLon:     dto.Lon,
	}

	if err := s.r.InitSchedule(ctx, &schedule); err != nil {
		return err
	}

	// Get staff and store information for logging and notification
	staff, _ := s.staff.GetByID(ctx, dto.StaffID)
	store, _ := s.store.GetByID(ctx, dto.StoreID)

	// Get supervisor for notification
	var supervisorID *uuid.UUID
	if staff != nil && staff.GroupID != uuid.Nil {
		group, err := s.group.GetByID(ctx, staff.GroupID)
		if err == nil && group != nil {
			supervisorID = group.SupervisorID
		}
	}

	// Log the unplanned schedule creation to audit
	_ = s.audit.Create(ctx, &model.AuditLogs{
		ActorID:     dto.StaffID,
		ActorType:   "staff",
		Action:      "InitSchedule",
		TargetTable: nil,
		TargetID:    &schedule.ID,
		IPAddress:   nil,
		CreatedAt:   time.Now().In(timezone.Dushanbe),
	})

	// Send system notification to supervisor if available
	if supervisorID != nil && store != nil && staff != nil {
		locMsg := ""
		if dto.Lat != nil && dto.Lon != nil {
			locMsg = fmt.Sprintf(" (Геолокация: %.4f, %.4f)", *dto.Lat, *dto.Lon)
		}
		body := fmt.Sprintf(
			"Сотрудник %s создал внеплановый визит\nМагазин: %s\nПричина: %s%s",
			staff.FullName, store.Name, dto.Description, locMsg,
		)
		sysNotif := &model.SystemNotification{
			RecipientID: *supervisorID,
			StaffID:     &staff.ID,
			Type:        "unplanned_schedule",
			Title:       "Внеплановый визит",
			Body:        body,
			IsRead:      false,
		}
		_ = s.sysNotification.Create(ctx, sysNotif)
	}

	return nil
}

func (s *scheduleService) CreateBulk(ctx context.Context, dto *model.CreateSchedulesBulk, adminID uuid.UUID, ip, scheduleConst, role string) (int, error) {
	staffID, err := uuid.Parse(dto.StaffID)
	if err != nil {
		return 0, err
	}

	selectedStores := make(map[uuid.UUID]struct{})

	if dto.GroupID != nil {
		groupID, err := uuid.Parse(*dto.GroupID)
		if err != nil {
			return 0, err
		}

		stores, err := s.store.ListByGroupID(ctx, groupID)
		if err != nil {
			return 0, err
		}

		for _, st := range stores {
			selectedStores[st.Store.ID] = struct{}{}
		}
	}

	for _, sID := range dto.StoreIDs {
		storeID, err := uuid.Parse(sID)
		if err != nil {
			return 0, fmt.Errorf("parse store id %s: %w", sID, err)
		}
		selectedStores[storeID] = struct{}{}
	}

	for _, excludeID := range dto.ExcludeStoreIDs {
		excludedStoreID, err := uuid.Parse(excludeID)
		if err != nil {
			return 0, fmt.Errorf("parse exclude_store_id %s: %w", excludeID, err)
		}
		delete(selectedStores, excludedStoreID)
	}

	if len(selectedStores) == 0 {
		return 0, errs.ErrorInvalidInput
	}

	schedules := make([]model.Schedule, 0, len(selectedStores))
	for storeID := range selectedStores {
		schedules = append(schedules, model.Schedule{
			StaffID:           staffID,
			StoreID:           storeID,
			PlannedDate:       dto.PlannedDate,
			Description:       dto.Description,
			WindowOpen:        dto.WindowOpen,
			WindowClose:       dto.WindowClose,
			IsPhotoRequired:   &dto.IsPhotoRequired,
			IsSelfieRequired:  &dto.IsSelfieRequired,
			IsCommentRequired: &dto.IsCommentRequired,
			CreatedByAdminID:  adminID,
		})
	}

	if err := s.r.CreateBulk(ctx, schedules); err != nil {
		return 0, err
	}

	staff, _ := s.staff.GetByID(ctx, staffID)
	if staff != nil && staff.FcmToken != nil {
		body := fmt.Sprintf("Вам назначено несколько новых задач (%d шт.) на %s",
			len(schedules), dto.PlannedDate.Format("02.01.2006"))

		_ = s.notification.CreateAndPushToUser(ctx, staff.ID, "Новые задачи", body, "schedule_assign_bulk", map[string]any{"type": "schedule_assign"}, staff.FcmToken)
	}

	_ = s.audit.Create(ctx, &model.AuditLogs{
		ActorID:     adminID,
		ActorType:   role,
		Action:      "AssignScheduleBulk",
		TargetTable: &scheduleConst,
		IPAddress:   &ip,
		CreatedAt:   time.Now().In(timezone.Dushanbe),
	})

	return len(schedules), nil
}

func (s *scheduleService) GetByID(ctx context.Context, id uuid.UUID) (*model.Schedule, error) {
	return s.r.GetByID(ctx, id)
}

func (s *scheduleService) Update(ctx context.Context, id uuid.UUID, dto model.UpdateSchedule) error {
	return s.r.Update(ctx, id, dto)
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

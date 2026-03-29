package service

import (
	"context"
	"fmt"
	"mime/multipart"
	"time"

	"github.com/IkbolNazarov/gps-tracker/internal/logger"
	"github.com/IkbolNazarov/gps-tracker/internal/model"
	"github.com/IkbolNazarov/gps-tracker/internal/repository"
	"github.com/IkbolNazarov/gps-tracker/internal/validator"
	errs "github.com/IkbolNazarov/gps-tracker/pkg/error"
	"github.com/IkbolNazarov/gps-tracker/pkg/fcm"
	minioservice "github.com/IkbolNazarov/gps-tracker/pkg/storage/service"
	"github.com/IkbolNazarov/gps-tracker/pkg/timezone"
	"github.com/IkbolNazarov/gps-tracker/pkg/utils"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type StaffUsers interface {
	Create(ctx context.Context, dto model.CreateStaffUser, admin_id uuid.UUID, ip, role, st, action string) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.StaffUser, error)
	Delete(ctx context.Context, id uuid.UUID, user_id uuid.UUID, ip, role, st, action string) error
	Update(ctx context.Context, id uuid.UUID, dto model.UpdateStaffUser, user_id uuid.UUID, ip, role, st, action string) error
	List(ctx context.Context, limit, offset int, groupID *uuid.UUID, status string) ([]model.StaffUser, int, error)
	Deactivate(ctx context.Context, id uuid.UUID, d model.DeactivateStaffUser, user_id uuid.UUID, ip, role, st, action string) error
	GetProfile(ctx context.Context, id uuid.UUID, lim, off int) (*model.StaffUserProfile, error)
	ListByGroupID(ctx context.Context, user_id *uuid.UUID) (staff []model.StaffUserLocation, err error)
	GetLocationByID(ctx context.Context, id uuid.UUID, begin, end time.Time) (*model.StaffUserLocationHistory, error)
	UpdateAvatar(ctx context.Context, userID uuid.UUID, file *multipart.FileHeader) error
	GetWorkTimeStaff(ctx context.Context, id uuid.UUID) (float64, error)
	DayStarted(ctx context.Context, id uuid.UUID) (bool, error)
}
type staffUserService struct {
	r             repository.StaffUsersRepository
	visitRepo     repository.VisitRepository
	schedules     Schedule
	audit         AuditLogs
	notifications Notification
	alertRepo     repository.GpsAlertStatesRepository
	minio         minioservice.Minio
	staffSchedule StaffSchedule
}

func NewStaffUserService(r repository.StaffUsersRepository, visitRepo repository.VisitRepository, sch Schedule, sender fcm.PushSender, notifications Notification, audit AuditLogs, minio minioservice.Minio, staffSchedule StaffSchedule, alertRepo repository.GpsAlertStatesRepository) StaffUsers {
	return &staffUserService{r: r, visitRepo: visitRepo, schedules: sch, notifications: notifications, audit: audit, minio: minio, staffSchedule: staffSchedule, alertRepo: alertRepo}
}

func (s *staffUserService) GetWorkTimeStaff(ctx context.Context, id uuid.UUID) (float64, error) {
	user, err := s.r.GetByID(ctx, id)
	if err != nil {
		return 0, err
	}

	if user.StartAt == nil {
		return 0, errs.ErrorNotStartedDay
	}

	if user.EndAt == nil {
		return 0, errs.ErrorNotFinishedDay
	}

	if user.StartAt.After(*user.EndAt) {
		return 0, errs.ErrorUnprocessable
	}

	return user.EndAt.Sub(*user.StartAt).Hours(), nil
}

func (s *staffUserService) DayStarted(ctx context.Context, id uuid.UUID) (bool, error) {
	user, err := s.r.GetByID(ctx, id)
	if err != nil {
		return false, err
	}

	if user.DayStarted == nil {
		return false, nil
	}

	return *user.DayStarted, nil
}

func (s *staffUserService) GetProfile(ctx context.Context, id uuid.UUID, lim, off int) (*model.StaffUserProfile, error) {

	user, err := s.r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if user.DeactivatedByID != nil {
		return nil, errs.ErrorUserDeactivated
	}

	if user.AvatarKey != nil {
		url, err := s.minio.PresignedDownload(ctx, *user.AvatarKey, 15*time.Minute)
		if err != nil {
			return nil, err
		}

		user.AvatarUrl = &url
	}

	visits, total, err := s.visitRepo.ListByStaffID(ctx, id, lim, off, time.Now().In(timezone.Dushanbe).AddDate(0, -1, 0), time.Now().In(timezone.Dushanbe))
	if err != nil {
		return nil, err
	}

	/* now := time.Now().In(timezone.Dushanbe)
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()) - если календарный месяц*/

	/* weekday := int(now.Weekday())
	if weekday == 0 {
	    weekday = 7
	}
	weekStart := time.Date(now.Year(), now.Month(), now.Day()-weekday+1, 0, 0, 0, 0, now.Location())  - неделя с понедельника*/

	doneMonth, totalMonth, err := s.visitRepo.CountVisit(ctx, id, time.Now().In(timezone.Dushanbe).AddDate(0, -1, 0), time.Now().In(timezone.Dushanbe))
	if err != nil {
		return nil, err
	}

	doneWeek, totalWeek, err := s.visitRepo.CountVisit(ctx, id, time.Now().In(timezone.Dushanbe).AddDate(0, 0, -7), time.Now().In(timezone.Dushanbe))
	if err != nil {
		return nil, err
	}

	schedules, err := s.staffSchedule.List(ctx, id)
	if err != nil {
		return nil, err
	}

	return &model.StaffUserProfile{
		User:       user,
		Visits:     visits,
		Schedules:  schedules,
		TotalCount: total,
		TotalMonth: totalMonth,
		DoneMonth:  doneMonth,
		TotalWeek:  totalWeek,
		DoneWeek:   doneWeek,
	}, nil
}

func (s *staffUserService) UpdateAvatar(ctx context.Context, userID uuid.UUID, file *multipart.FileHeader) error {

	user, err := s.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	if user.DeactivatedByID != nil {
		return errs.ErrorUserDeactivated
	}

	ext := utils.GetFileExtension(file.Filename)
	objectName := fmt.Sprintf("staff/avatars/%s%s", userID, ext)

	key, err := s.minio.UploadFile(ctx, file, objectName)
	if err != nil {
		return err
	}

	return s.r.Update(ctx, userID, model.UpdateStaffUser{
		AvatarKey: &key,
	})
}

func (s *staffUserService) Create(ctx context.Context, dto model.CreateStaffUser, admin_id uuid.UUID, ip, role, st, action string) error {
	if dto.Password != dto.ConfirmPassword {
		return errs.ErrorPasswordsDoNotMatch
	}

	if !validator.IsStrongPassword(dto.Password) {
		return errs.ErrorWeakPassword
	}

	normalizedPhone, ok := validator.NormalizeTajikPhone(dto.PhoneNumber)
	if !ok {
		return errs.ErrorInvalidPhone
	}

	normalizedINN, ok := validator.NormalizeAndValidateINN(dto.Inn)
	if !ok {
		return errs.ErrorInvalidINN
	}

	groupID, err := uuid.Parse(dto.GroupID)
	if err != nil {
		return errs.ErrorInvalidGroup
	}

	hash, err := utils.HashPassword(dto.Password)
	if err != nil {
		return errs.ErrorPasswordHashError
	}

	user := &model.StaffUser{
		GroupID:             groupID,
		FullName:            dto.FullName,
		Login:               dto.Login,
		PhoneNumber:         normalizedPhone,
		Inn:                 normalizedINN,
		Birthday:            dto.Birthday,
		PasswordHash:        hash,
		IsTemporaryPassword: true,
	}

	if err := s.r.Create(ctx, user); err != nil {
		return err
	}

	if len(dto.Schedule) > 0 {
		if err := s.staffSchedule.CreateBatch(ctx, user.ID, dto.Schedule); err != nil {
			return err
		}
	}

	if len(dto.Schedule) == 0 {
		defaultSchedule := make([]model.ScheduleDTO, 0, 5)

		for i := int16(1); i <= 5; i++ {
			start := time.Date(0, 0, 0, 9, 0, 0, 0, timezone.Dushanbe)
			end := time.Date(0, 0, 0, 18, 0, 0, 0, timezone.Dushanbe)

			defaultSchedule = append(defaultSchedule, model.ScheduleDTO{
				Weekday:   i,
				StartTime: start,
				EndTime:   end,
			})
		}

		if err := s.staffSchedule.CreateBatch(ctx, user.ID, defaultSchedule); err != nil {
			return err
		}
	}

	if admin_id != uuid.Nil {
		if err := s.audit.Create(ctx, &model.AuditLogs{
			ActorID:     admin_id,
			ActorType:   role,
			TargetTable: &st,
			IPAddress:   &ip,
			Action:      action,
			CreatedAt:   time.Now().In(timezone.Dushanbe),
		}); err != nil {
			return err
		}
	}

	return nil
}

func (s *staffUserService) GetByID(ctx context.Context, id uuid.UUID) (*model.StaffUser, error) {

	user, err := s.r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if user.DeactivatedByID != nil {
		return nil, errs.ErrorUserDeactivated
	}

	return user, nil
}

func (s *staffUserService) List(ctx context.Context, limit, offset int, groupID *uuid.UUID, status string) ([]model.StaffUser, int, error) {
	users, total, err := s.r.List(ctx, limit, offset, groupID, status)
	if err != nil {
		return nil, 0, err
	}

	for i := range users {
		if users[i].AvatarKey != nil {
			url, err := s.minio.PresignedDownload(ctx, *users[i].AvatarKey, time.Minute*15)
			if err != nil {
				return nil, 0, err
			}

			users[i].AvatarUrl = &url
		}
	}

	return users, total, nil
}

func (s *staffUserService) Update(ctx context.Context, id uuid.UUID, dto model.UpdateStaffUser, user_id uuid.UUID, ip, role, st, action string) error {
	user, err := s.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if user.DeactivatedByID != nil {
		return errs.ErrorUserDeactivated
	}

	if err := s.r.Update(ctx, id, dto); err != nil {
		return err
	}

	// Если сотрудник заканчивает день — сбрасываем активный GPS-алерт.
	// Это нужно чтобы на следующий рабочий день алерт снова мог сработать.
	if dto.EndAt != nil {
		if err := s.alertRepo.Resolve(ctx, id); err != nil {
			logger.Log.Warn("staff update: resolve gps alert on day end",
				zap.Stringer("staff_id", id),
				zap.Error(err),
			)
		}
	}

	if user_id != uuid.Nil {
		if err := s.audit.Create(ctx, &model.AuditLogs{
			ActorID:     user_id,
			ActorType:   role,
			TargetID:    &id,
			TargetTable: &st,
			Action:      action,
			IPAddress:   &ip,
			CreatedAt:   time.Now().In(timezone.Dushanbe),
		}); err != nil {
			return err
		}
	}

	return nil
}

func (s *staffUserService) Delete(ctx context.Context, id uuid.UUID, user_id uuid.UUID, ip, role, st, action string) error {
	if err := s.r.Delete(ctx, id); err != nil {
		return err
	}

	if err := s.audit.Create(ctx, &model.AuditLogs{
		ActorID:     user_id,
		ActorType:   role,
		TargetID:    &id,
		TargetTable: &st,
		Action:      action,
		IPAddress:   &ip,
		CreatedAt:   time.Now().In(timezone.Dushanbe),
	}); err != nil {
		return err
	}

	return nil
}

func (s *staffUserService) Deactivate(ctx context.Context, id uuid.UUID, d model.DeactivateStaffUser, user_id uuid.UUID, ip, role, st, action string) error {

	user, err := s.r.GetByID(ctx, id)
	if err != nil || user.FcmToken == nil {
		return nil
	}

	if user.DeactivatedByID != nil {
		return errs.ErrorUserDeactivated
	}

	if err := s.r.Deactivate(ctx, id, d); err != nil {
		return err
	}

	if err := s.notifications.CreateAndPushToUser(ctx, id, "access deny", "account is deactivated", "account_blocked", map[string]any{"type": "account_blocked"}, user.FcmToken); err != nil {
		logger.Log.Error(fmt.Sprintf("failed to create/push notification to %s", user.ID), zap.Error(err))
	}

	if err := s.audit.Create(ctx, &model.AuditLogs{
		ActorID:     user_id,
		ActorType:   role,
		TargetID:    &id,
		TargetTable: &st,
		Action:      action,
		IPAddress:   &ip,
		CreatedAt:   time.Now().In(timezone.Dushanbe),
	}); err != nil {
		return err
	}

	return nil
}

func (s *staffUserService) ListByGroupID(ctx context.Context, supervisorID *uuid.UUID) ([]model.StaffUserLocation, error) {
	users, err := s.r.ListByGroupID(ctx, supervisorID)
	if err != nil {
		return nil, err
	}

	for i := range users {
		if users[i].AvatarKey != nil {
			url, err := s.minio.PresignedDownload(ctx, *users[i].AvatarKey, time.Minute*15)
			if err != nil {
				return nil, err
			}
			users[i].AvatarUrl = &url
		}
	}

	return users, nil
}

func (s *staffUserService) GetLocationByID(ctx context.Context, id uuid.UUID, begin, end time.Time) (*model.StaffUserLocationHistory, error) {
	return s.r.GetLocationByID(ctx, id, begin, end)
}

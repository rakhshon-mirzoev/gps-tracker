package service

import (
	"context"
	"fmt"
	"mime/multipart"
	"time"

	"github.com/google/uuid"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/logger"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/model"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/repository"
	errors "github.com/rakhshon-mirzoev/gps-tracker/pkg/error"
	errs "github.com/rakhshon-mirzoev/gps-tracker/pkg/error"
	"github.com/rakhshon-mirzoev/gps-tracker/pkg/fcm"
	minioservice "github.com/rakhshon-mirzoev/gps-tracker/pkg/storage/service"
	"github.com/rakhshon-mirzoev/gps-tracker/pkg/utils"
	"go.uber.org/zap"
)

type StaffUsers interface {
	Create(ctx context.Context, dto model.CreateStaffUser, admin_id uuid.UUID, ip, role, st, action string) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.StaffUser, error)
	Delete(ctx context.Context, id uuid.UUID, user_id uuid.UUID, ip, role, st, action string) error
	Update(ctx context.Context, id uuid.UUID, dto model.UpdateStaffUser, user_id uuid.UUID, ip, role, st, action string) error
	List(ctx context.Context, limit, offset int) ([]model.StaffUser, int, error)
	Deactivate(ctx context.Context, id uuid.UUID, d model.DeactivateStaffUser, user_id uuid.UUID, ip, role, st, action string) error
	GetProfile(ctx context.Context, id uuid.UUID, lim, off int) (*model.StaffUserProfile, error)
	ListByGroupID(ctx context.Context, user_id uuid.UUID) (staff []model.StaffUserLocation, err error)
	GetLocationByID(ctx context.Context, id uuid.UUID, begin, end time.Time) (*model.StaffUserLocationHistory, error)
	UpdateAvatar(ctx context.Context, userID uuid.UUID, file *multipart.FileHeader) error
}
type staffUserService struct {
	r             repository.StaffUsersRepository
	schedules     Schedule
	audit         AuditLogs
	notifications Notification
	minio         minioservice.Minio
}

func NewStaffUserService(r repository.StaffUsersRepository, sch Schedule, sender fcm.PushSender, notifications Notification, audit AuditLogs, minio minioservice.Minio) StaffUsers {
	return &staffUserService{r: r, schedules: sch, notifications: notifications, audit: audit, minio: minio}
}

func (s *staffUserService) GetProfile(ctx context.Context, id uuid.UUID, lim, off int) (*model.StaffUserProfile, error) {

	user, err := s.r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if user.AvatarKey != nil {
		url, err := s.minio.PresignedDownload(ctx, *user.AvatarKey, 15*time.Minute)
		if err != nil {
			return nil, err
		}

		user.AvatarUrl = &url
	}

	schedules, total, err := s.schedules.ListByStaffID(ctx, id, lim, off)
	if err != nil {
		return nil, err
	}

	return &model.StaffUserProfile{
		User:       user,
		Schedules:  schedules,
		TotalCount: total,
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

	objectName := fmt.Sprintf("staff/avatars/%s.jpg", userID)

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

	groupID, err := uuid.Parse(dto.GroupID)
	if err != nil {
		return errs.ErrorInvalidGroup
	}

	hash, err := utils.HashPassword(dto.Password)
	if err != nil {
		return errors.ErrorPasswordHashError
	}

	user := &model.StaffUser{
		GroupID:             groupID,
		FullName:            dto.FullName,
		Login:               dto.Login,
		PhoneNumber:         dto.PhoneNumber,
		Inn:                 dto.Inn,
		Birthday:            dto.Birthday,
		PasswordHash:        hash,
		IsTemporaryPassword: true,
	}

	if err := s.r.Create(ctx, user); err != nil {
		return err
	}

	if admin_id != uuid.Nil {
		if err := s.audit.Create(ctx, &model.AuditLogs{
			ActorID:     admin_id,
			ActorType:   role,
			TargetTable: &st,
			IPAddress:   &ip,
			Action:      action,
		}); err != nil {
			return err
		}
	}

	return nil
}

func (s *staffUserService) GetByID(ctx context.Context, id uuid.UUID) (*model.StaffUser, error) {
	return s.r.GetByID(ctx, id)
}

func (s *staffUserService) List(ctx context.Context, limit, offset int) ([]model.StaffUser, int, error) {
	return s.r.List(ctx, limit, offset)
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

	if user_id != uuid.Nil {
		if err := s.audit.Create(ctx, &model.AuditLogs{
			ActorID:     user_id,
			ActorType:   role,
			TargetID:    &id,
			TargetTable: &st,
			Action:      action,
			IPAddress:   &ip,
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
	}); err != nil {
		return err
	}

	return nil
}

func (s *staffUserService) Deactivate(ctx context.Context, id uuid.UUID, d model.DeactivateStaffUser, user_id uuid.UUID, ip, role, st, action string) error {
	if err := s.r.Deactivate(ctx, id, d); err != nil {
		return err
	}

	user, err := s.r.GetByID(ctx, id)
	if err != nil || user.FcmToken == nil {
		return nil
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
	}); err != nil {
		return err
	}

	return nil
}

func (s *staffUserService) ListByGroupID(ctx context.Context, user_id uuid.UUID) (staff []model.StaffUserLocation, err error) {
	return s.r.ListByGroupID(ctx, user_id)
}

func (s *staffUserService) GetLocationByID(ctx context.Context, id uuid.UUID, begin, end time.Time) (*model.StaffUserLocationHistory, error) {
	return s.r.GetLocationByID(ctx, id, begin, end)
}

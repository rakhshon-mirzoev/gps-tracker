package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/logger"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/model"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/repository"
	errors "github.com/rakhshon-mirzoev/gps-tracker/pkg/error"
	errs "github.com/rakhshon-mirzoev/gps-tracker/pkg/error"
	"github.com/rakhshon-mirzoev/gps-tracker/pkg/fcm"
	"github.com/rakhshon-mirzoev/gps-tracker/pkg/utils"
	"go.uber.org/zap"
)

type StaffUsers interface {
	Create(ctx context.Context, dto model.CreateStaffUser, admin_id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.StaffUser, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Update(ctx context.Context, id uuid.UUID, dto model.UpdateStaffUser, user_id uuid.UUID) error
	List(ctx context.Context) ([]model.StaffUser, error)
	Deactivate(ctx context.Context, id uuid.UUID, d model.DeactivateStaffUser) error
	GetProfile(ctx context.Context, id uuid.UUID, lim, off int) (*model.StaffUserProfile, error)
	ListByGroupID(ctx context.Context, user_id uuid.UUID) (staff []model.StaffUserLocation, err error)
	GetLocationByID(ctx context.Context, id uuid.UUID, begin, end time.Time) (*model.StaffUserLocationHistory, error)
}
type staffUserService struct {
	r             repository.StaffUsersRepository
	schedules     Schedule
	audit         AuditLogs
	notifications Notification
}

func NewStaffUserService(r repository.StaffUsersRepository, sch Schedule, sender fcm.PushSender, notifications Notification) StaffUsers {
	return &staffUserService{r: r, schedules: sch, notifications: notifications}
}

func (s *staffUserService) GetProfile(ctx context.Context, id uuid.UUID, lim, off int) (*model.StaffUserProfile, error) {

	user, err := s.r.GetByID(ctx, id)
	if err != nil {
		return nil, err
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

func (s *staffUserService) Create(ctx context.Context, dto model.CreateStaffUser, admin_id uuid.UUID) error {
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

	return nil
}

func (s *staffUserService) GetByID(ctx context.Context, id uuid.UUID) (*model.StaffUser, error) {
	return s.r.GetByID(ctx, id)
}

func (s *staffUserService) List(ctx context.Context) ([]model.StaffUser, error) {
	return s.r.List(ctx)
}

func (s *staffUserService) Update(ctx context.Context, id uuid.UUID, dto model.UpdateStaffUser, user_id uuid.UUID) error {
	return s.r.Update(ctx, id, dto)
}

func (s *staffUserService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.r.Delete(ctx, id)
}

func (s *staffUserService) Deactivate(ctx context.Context, id uuid.UUID, d model.DeactivateStaffUser) error {
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

	return nil
}

func (s *staffUserService) ListByGroupID(ctx context.Context, user_id uuid.UUID) (staff []model.StaffUserLocation, err error) {
	return s.r.ListByGroupID(ctx, user_id)
}

func (s *staffUserService) GetLocationByID(ctx context.Context, id uuid.UUID, begin, end time.Time) (*model.StaffUserLocationHistory, error) {
	return s.r.GetLocationByID(ctx, id, begin, end)
}

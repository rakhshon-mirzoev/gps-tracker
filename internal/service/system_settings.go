package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/model"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/repository"
)

type SystemSettings interface {
	Create(ctx context.Context, ss *model.SystemSettings, user_id uuid.UUID, role, ip, settings, action string) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.SystemSettings, error)
	List(ctx context.Context, limit, offset int) ([]model.SystemSettings, int, error)
	Update(ctx context.Context, id uuid.UUID, dto model.UpdateSystemSettings, user_id uuid.UUID, role, ip, settings, action string) error
	Delete(ctx context.Context, id uuid.UUID, user_id uuid.UUID, role, ip, settings, action string) error
}

type systemSettingsService struct {
	repo  repository.SystemSettingsRepository
	audit AuditLogs
}

func NewSystemSettingsService(repo repository.SystemSettingsRepository, audit AuditLogs) SystemSettings {
	return &systemSettingsService{repo: repo, audit: audit}
}

func (s *systemSettingsService) Create(ctx context.Context, ss *model.SystemSettings, user_id uuid.UUID, role, ip, settings, action string) error {

	if err := s.repo.Create(ctx, ss); err != nil {
		return err
	}

	if err := s.audit.Create(ctx, &model.AuditLogs{
		ActorID:     user_id,
		ActorType:   role,
		IPAddress:   &ip,
		Action:      action,
		TargetTable: &settings,
	}); err != nil {
		return err
	}

	return nil
}

func (s *systemSettingsService) GetByID(ctx context.Context, id uuid.UUID) (*model.SystemSettings, error) {
	if id == uuid.Nil {
		return nil, errors.New("id is required")
	}

	return s.repo.GetByID(ctx, id)
}

func (s *systemSettingsService) List(ctx context.Context, limit, offset int) ([]model.SystemSettings, int, error) {
	return s.repo.List(ctx, limit, offset)
}

func (s *systemSettingsService) Update(ctx context.Context, id uuid.UUID, dto model.UpdateSystemSettings, user_id uuid.UUID, role, ip, settings, action string) error {
	if err := s.repo.Update(ctx, id, dto); err != nil {
		return err
	}

	if err := s.audit.Create(ctx, &model.AuditLogs{
		ActorID:     user_id,
		ActorType:   role,
		IPAddress:   &ip,
		Action:      action,
		TargetTable: &settings,
	}); err != nil {
		return err
	}

	return nil
}

func (s *systemSettingsService) Delete(ctx context.Context, id uuid.UUID, user_id uuid.UUID, role, ip, settings, action string) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	if err := s.audit.Create(ctx, &model.AuditLogs{
		ActorID:     user_id,
		ActorType:   role,
		IPAddress:   &ip,
		Action:      action,
		TargetTable: &settings,
	}); err != nil {
		return err
	}

	return nil
}

package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/model"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/repository"
)

type SystemSettingsService interface {
	Create(ctx context.Context, ss *model.SystemSettings) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.SystemSettings, error)
	List(ctx context.Context) ([]model.SystemSettings, error)
	Update(ctx context.Context, id uuid.UUID, dto model.UpdateSystemSettings) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type systemSettingsService struct {
	repo repository.SystemSettingsRepository
}

func NewSystemSettingsService(repo repository.SystemSettingsRepository) SystemSettingsService {
	return &systemSettingsService{repo: repo}
}

func (s *systemSettingsService) Create(ctx context.Context, ss *model.SystemSettings) error {
	if ss.Key == "" {
		return errors.New("key is required")
	}
	if ss.Value == "" {
		return errors.New("value is required")
	}

	return s.repo.Create(ctx, ss)
}

func (s *systemSettingsService) GetByID(ctx context.Context, id uuid.UUID) (*model.SystemSettings, error) {
	if id == uuid.Nil {
		return nil, errors.New("id is required")
	}

	return s.repo.GetByID(ctx, id)
}

func (s *systemSettingsService) List(ctx context.Context) ([]model.SystemSettings, error) {
	return s.repo.List(ctx)
}

func (s *systemSettingsService) Update(ctx context.Context, id uuid.UUID, dto model.UpdateSystemSettings) error {
	return s.repo.Update(ctx, id, dto)
}

func (s *systemSettingsService) Delete(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("id is required")
	}

	return s.repo.Delete(ctx, id)
}

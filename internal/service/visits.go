package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/model"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/repository"
)

type VisitsService interface {
	Create(ctx context.Context, v *model.Visit) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Visit, error)
	List(ctx context.Context, scheduleID uuid.UUID) ([]model.Visit, error)
	Update(ctx context.Context, id uuid.UUID, dto model.UpdateVisit) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type visitsService struct {
	repo repository.VisitRepository
}

func NewVisitsService(repo repository.VisitRepository) VisitsService {
	return &visitsService{repo: repo}
}

func (s *visitsService) Create(ctx context.Context, v *model.Visit) error {
	if v.StaffID == uuid.Nil {
		return errors.New("staff_id is required")
	}
	if v.StoreID == uuid.Nil {
		return errors.New("store_id is required")
	}
	if v.ReviewStatus == "" {
		return errors.New("review_status is required")
	}

	return s.repo.Create(ctx, v)
}

func (s *visitsService) GetByID(ctx context.Context, id uuid.UUID) (*model.Visit, error) {
	if id == uuid.Nil {
		return nil, errors.New("id is required")
	}

	return s.repo.GetByID(ctx, id)
}

func (s *visitsService) List(ctx context.Context, scheduleID uuid.UUID) ([]model.Visit, error) {
	if scheduleID == uuid.Nil {
		return nil, errors.New("schedule_id is required")
	}

	return s.repo.List(ctx, scheduleID)
}

func (s *visitsService) Update(ctx context.Context, id uuid.UUID, dto model.UpdateVisit) error {
	return s.repo.Update(ctx, id, dto)
}

func (s *visitsService) Delete(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("id is required")
	}

	return s.repo.Delete(ctx, id)
}

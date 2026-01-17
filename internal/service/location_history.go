package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/model"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/repository"
)

type LocationHistoryService interface {
	Create(ctx context.Context, lh *model.LocationHistory) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.LocationHistory, error)
	List(ctx context.Context, staffID uuid.UUID) ([]model.LocationHistory, error)
	Update(ctx context.Context, lh *model.LocationHistory) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type locationHistoryService struct {
	repo repository.LocationHistoryRepository
}

func NewLocationHistoryService(repo repository.LocationHistoryRepository) LocationHistoryService {
	return &locationHistoryService{repo: repo}
}

func (s *locationHistoryService) Create(ctx context.Context, lh *model.LocationHistory) error {
	if lh.StaffID == uuid.Nil {
		return errors.New("staff_id is required")
	}
	if lh.Location == "" {
		return errors.New("location is required")
	}

	return s.repo.Create(ctx, lh)
}

func (s *locationHistoryService) GetByID(ctx context.Context, id uuid.UUID) (*model.LocationHistory, error) {
	if id == uuid.Nil {
		return nil, errors.New("id is required")
	}

	return s.repo.GetByID(ctx, id)
}

func (s *locationHistoryService) List(ctx context.Context, staffID uuid.UUID) ([]model.LocationHistory, error) {
	if staffID == uuid.Nil {
		return nil, errors.New("staff_id is required")
	}

	return s.repo.List(ctx, staffID)
}

func (s *locationHistoryService) Update(ctx context.Context, lh *model.LocationHistory) error {
	if lh.ID == uuid.Nil {
		return errors.New("id is required")
	}
	if lh.StaffID == uuid.Nil {
		return errors.New("staff_id is required")
	}
	if lh.Location == "" {
		return errors.New("location is required")
	}

	return s.repo.Update(ctx, lh)
}

func (s *locationHistoryService) Delete(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("id is required")
	}

	return s.repo.Delete(ctx, id)
}

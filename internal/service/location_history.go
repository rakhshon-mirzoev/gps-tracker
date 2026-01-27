package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/model"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/repository"
)

type LocationHistory interface {
	Create(ctx context.Context, lh *model.LocationHistory, user_id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.LocationHistory, error)
	List(ctx context.Context, staffID uuid.UUID) ([]model.LocationHistory, error)
	Update(ctx context.Context, id uuid.UUID, dto model.UpdateLocationHistory) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type locationHistoryService struct {
	r repository.LocationHistoryRepository
}

func NewLocationHistoryService(r repository.LocationHistoryRepository) LocationHistory {
	return &locationHistoryService{r: r}
}

func (s *locationHistoryService) Create(ctx context.Context, lh *model.LocationHistory, user_id uuid.UUID) error {
	lh.StaffID = user_id
	return s.r.Create(ctx, lh)
}

func (s *locationHistoryService) GetByID(ctx context.Context, id uuid.UUID) (*model.LocationHistory, error) {
	if id == uuid.Nil {
		return nil, errors.New("id is required")
	}

	return s.r.GetByID(ctx, id)
}

func (s *locationHistoryService) List(ctx context.Context, staffID uuid.UUID) ([]model.LocationHistory, error) {
	if staffID == uuid.Nil {
		return nil, errors.New("staff_id is required")
	}

	return s.r.List(ctx, staffID)
}

func (s *locationHistoryService) Update(ctx context.Context, id uuid.UUID, dto model.UpdateLocationHistory) error {
	return s.r.Update(ctx, id, dto)
}

func (s *locationHistoryService) Delete(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("id is required")
	}

	return s.r.Delete(ctx, id)
}

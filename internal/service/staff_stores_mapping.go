package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/model"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/repository"
)

type StaffStoresMapping interface {
	Create(ctx context.Context, ssm *model.StaffStoresMapping) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.StaffStoresMapping, error)
	List(ctx context.Context, staffID uuid.UUID) ([]model.StaffStoresMapping, error)
	Update(ctx context.Context, ssm *model.StaffStoresMapping) error
	Delete(ctx context.Context, id uuid.UUID) error
	ListByStaffID(ctx context.Context, id uuid.UUID, lim, off int) ([]*model.Store, error)
}

type staffStoresMappingService struct {
	r repository.StaffStoresMappingRepository
}

func NewStaffStoresMappingService(r repository.StaffStoresMappingRepository) StaffStoresMapping {
	return &staffStoresMappingService{r: r}
}

func (s *staffStoresMappingService) Create(ctx context.Context, ssm *model.StaffStoresMapping) error {
	if ssm.StaffID == uuid.Nil {
		return errors.New("staff_id is required")
	}
	if ssm.StoreID == uuid.Nil {
		return errors.New("store_id is required")
	}

	return s.r.Create(ctx, ssm)
}

func (s *staffStoresMappingService) GetByID(ctx context.Context, id uuid.UUID) (*model.StaffStoresMapping, error) {
	if id == uuid.Nil {
		return nil, errors.New("id is required")
	}

	return s.r.GetByID(ctx, id)
}

func (s *staffStoresMappingService) List(ctx context.Context, staffID uuid.UUID) ([]model.StaffStoresMapping, error) {
	if staffID == uuid.Nil {
		return nil, errors.New("staff_id is required")
	}

	return s.r.List(ctx, staffID)
}

func (s *staffStoresMappingService) Update(ctx context.Context, ssm *model.StaffStoresMapping) error {
	if ssm.ID == uuid.Nil {
		return errors.New("id is required")
	}
	if ssm.StaffID == uuid.Nil {
		return errors.New("staff_id is required")
	}
	if ssm.StoreID == uuid.Nil {
		return errors.New("store_id is required")
	}

	return s.r.Update(ctx, ssm)
}

func (s *staffStoresMappingService) Delete(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("id is required")
	}

	return s.r.Delete(ctx, id)
}

func (s *staffStoresMappingService) ListByStaffID(ctx context.Context, id uuid.UUID, lim, off int) ([]*model.Store, error) {
	return s.r.ListByStaffID(ctx, id, lim, off)
}

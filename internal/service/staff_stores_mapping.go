package service

import (
	"context"
	"errors"

	"github.com/IkbolNazarov/gps-tracker/internal/model"
	"github.com/IkbolNazarov/gps-tracker/internal/repository"
	"github.com/google/uuid"
)

type StaffStoresMapping interface {
	Create(ctx context.Context, ssm *model.StaffStoresMapping) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.StaffStoresMapping, error)
	// List(ctx context.Context, staffID uuid.UUID) ([]model.StaffStoresMapping, error)
	Update(ctx context.Context, ssm *model.StaffStoresMapping) error
	Delete(ctx context.Context, store_id, staff_id uuid.UUID) error
	ListByStaffID(ctx context.Context, id uuid.UUID, lim, off int) ([]model.Store, int, error)
}

type staffStoresMappingService struct {
	r     repository.StaffStoresMappingRepository
	staff StaffUsers
}

func NewStaffStoresMappingService(r repository.StaffStoresMappingRepository, staff StaffUsers) StaffStoresMapping {
	return &staffStoresMappingService{r: r, staff: staff}
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

func (s *staffStoresMappingService) Delete(ctx context.Context, store_id, staff_id uuid.UUID) error {

	return s.r.Delete(ctx, store_id, staff_id)
}

func (s *staffStoresMappingService) ListByStaffID(ctx context.Context, id uuid.UUID, lim, off int) ([]model.Store, int, error) {
	_, err := s.staff.GetByID(ctx, id)
	if err != nil {
		return nil, 0, err
	}

	store, total, err := s.r.ListByStaffID(ctx, id, lim, off)
	if err != nil {
		return nil, 0, err
	}

	return store, total, nil
}

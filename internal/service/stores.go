package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/model"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/repository"
)

type Stores interface {
	Create(ctx context.Context, s *model.CreateStore, user_id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Store, error)
	List(ctx context.Context, filter model.StoreListFilter) ([]model.Store, int, error)
	Update(ctx context.Context, s model.UpdateStore, id uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type storesService struct {
	r repository.StoreRepository
}

func NewStoresService(repo repository.StoreRepository) Stores {
	return &storesService{r: repo}
}

func (s *storesService) Create(ctx context.Context, store *model.CreateStore, user_id uuid.UUID) error {
	staffID, err := uuid.Parse(store.CreatedByStaffID)
	if err != nil {
		return err
	}

	st := &model.Store{
		Name:             store.Name,
		Address:          store.Address,
		Latitude:         store.Latitude,
		Longitude:        store.Longitude,
		OwnerFullName:    store.OwnerFullName,
		OwnerINN:         store.OwnerINN,
		OwnerPassport:    store.OwnerPassport,
		CreatedByStaffID: staffID,
	}

	return s.r.Create(ctx, st)
}

func (s *storesService) GetByID(ctx context.Context, id uuid.UUID) (*model.Store, error) {
	if id == uuid.Nil {
		return nil, errors.New("id is required")
	}

	return s.r.GetByID(ctx, id)
}

func (s *storesService) List(ctx context.Context, filter model.StoreListFilter) ([]model.Store, int, error) {
	return s.r.List(ctx, filter)
}

func (s *storesService) Update(ctx context.Context, store model.UpdateStore, id uuid.UUID) error {
	return s.r.Update(ctx, store, id)
}

func (s *storesService) Delete(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("id is required")
	}

	return s.r.Delete(ctx, id)
}

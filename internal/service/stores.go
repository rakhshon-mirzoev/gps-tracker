package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/model"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/repository"
)

type StoresService interface {
	Create(ctx context.Context, s *model.Store) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Store, error)
	List(ctx context.Context, filter model.StoreListFilter) ([]*model.Store, int, error)
	Update(ctx context.Context, s *model.Store) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type storesService struct {
	repo repository.StoreRepository
}

func NewStoresService(repo repository.StoreRepository) StoresService {
	return &storesService{repo: repo}
}

func (s *storesService) Create(ctx context.Context, store *model.Store) error {
	if store.Name == "" {
		return errors.New("name is required")
	}

	return s.repo.Create(ctx, store)
}

func (s *storesService) GetByID(ctx context.Context, id uuid.UUID) (*model.Store, error) {
	if id == uuid.Nil {
		return nil, errors.New("id is required")
	}

	return s.repo.GetByID(ctx, id)
}

func (s *storesService) List(ctx context.Context, filter model.StoreListFilter) ([]*model.Store, int, error) {
	return s.repo.List(ctx, filter)
}

func (s *storesService) Update(ctx context.Context, store *model.Store) error {
	if store.ID == uuid.Nil {
		return errors.New("id is required")
	}
	if store.Name == "" {
		return errors.New("name is required")
	}

	return s.repo.Update(ctx, store)
}

func (s *storesService) Delete(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("id is required")
	}

	return s.repo.Delete(ctx, id)
}

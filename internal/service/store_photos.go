package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/model"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/repository"
)

type StorePhotosService interface {
	Create(ctx context.Context, sp *model.StorePhotos) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.StorePhotos, error)
	List(ctx context.Context, storeID uuid.UUID) ([]model.StorePhotos, error)
	Update(ctx context.Context, id uuid.UUID, dto model.UpdateStorePhotos) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type storePhotosService struct {
	repo repository.StorePhotosRepository
}

func NewStorePhotosService(repo repository.StorePhotosRepository) StorePhotosService {
	return &storePhotosService{repo: repo}
}

func (s *storePhotosService) Create(ctx context.Context, sp *model.StorePhotos) error {
	if sp.StoreID == uuid.Nil {
		return errors.New("store_id is required")
	}
	if sp.PhotoURL == "" {
		return errors.New("photo_url is required")
	}
	if sp.Type == "" {
		return errors.New("type is required")
	}

	return s.repo.Create(ctx, sp)
}

func (s *storePhotosService) GetByID(ctx context.Context, id uuid.UUID) (*model.StorePhotos, error) {
	if id == uuid.Nil {
		return nil, errors.New("id is required")
	}

	return s.repo.GetByID(ctx, id)
}

func (s *storePhotosService) List(ctx context.Context, storeID uuid.UUID) ([]model.StorePhotos, error) {
	if storeID == uuid.Nil {
		return nil, errors.New("store_id is required")
	}

	return s.repo.List(ctx, storeID)
}

func (s *storePhotosService) Update(ctx context.Context, id uuid.UUID, dto model.UpdateStorePhotos) error {
	return s.repo.Update(ctx, id, dto)
}

func (s *storePhotosService) Delete(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("id is required")
	}

	return s.repo.Delete(ctx, id)
}

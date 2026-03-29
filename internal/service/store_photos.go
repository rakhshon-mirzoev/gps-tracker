package service

import (
	"context"
	"errors"

	"github.com/IkbolNazarov/gps-tracker/internal/model"
	"github.com/IkbolNazarov/gps-tracker/internal/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type StorePhotos interface {
	Create(ctx context.Context, sp *model.StorePhotos) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.StorePhotos, error)
	List(ctx context.Context, storeID uuid.UUID) ([]model.StorePhotos, error)
	Delete(ctx context.Context, id uuid.UUID) error
	CreateTx(ctx context.Context, tx pgx.Tx, sp *model.StorePhotos) error
}

type storePhotosService struct {
	repo repository.StorePhotosRepository
}

func NewStorePhotosService(repo repository.StorePhotosRepository) StorePhotos {
	return &storePhotosService{repo: repo}
}

func (s *storePhotosService) Create(ctx context.Context, sp *model.StorePhotos) error {
	return s.repo.Create(ctx, sp)
}

func (s *storePhotosService) CreateTx(ctx context.Context, tx pgx.Tx, sp *model.StorePhotos) error {
	return s.repo.CreateTx(ctx, tx, sp)
}

func (s *storePhotosService) GetByID(ctx context.Context, id uuid.UUID) (*model.StorePhotos, error) {
	if id == uuid.Nil {
		return nil, errors.New("id is required")
	}

	return s.repo.GetByID(ctx, id)
}

func (s *storePhotosService) List(ctx context.Context, storeID uuid.UUID) ([]model.StorePhotos, error) {
	return s.repo.List(ctx, storeID)
}

func (s *storePhotosService) Delete(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("id is required")
	}

	return s.repo.Delete(ctx, id)
}

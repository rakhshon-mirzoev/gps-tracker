package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/model"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/repository"
	minioservice "github.com/rakhshon-mirzoev/gps-tracker/pkg/storage/service"
)

type VisitPhotos interface {
	Create(ctx context.Context, vp *model.VisitPhotos) error
	GetByVisitID(ctx context.Context, id uuid.UUID) ([]model.VisitPhotos, error)
	List(ctx context.Context, visitID uuid.UUID) ([]model.VisitPhotos, error)
	Update(ctx context.Context, id uuid.UUID, dto model.UpdateVisitPhotos) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type visitPhotosService struct {
	repo  repository.VisitPhotosRepository
	minio minioservice.Minio
}

func NewVisitPhotosService(repo repository.VisitPhotosRepository, minio minioservice.Minio) VisitPhotos {
	return &visitPhotosService{repo: repo}
}

func (s *visitPhotosService) Create(ctx context.Context, vp *model.VisitPhotos) error {

	return s.repo.Create(ctx, vp)
}

func (s *visitPhotosService) GetByVisitID(ctx context.Context, id uuid.UUID) ([]model.VisitPhotos, error) {
	return s.repo.GetByVisitID(ctx, id)
}

func (s *visitPhotosService) List(ctx context.Context, visitID uuid.UUID) ([]model.VisitPhotos, error) {
	if visitID == uuid.Nil {
		return nil, errors.New("visit_id is required")
	}

	return s.repo.List(ctx, visitID)
}

func (s *visitPhotosService) Update(ctx context.Context, id uuid.UUID, dto model.UpdateVisitPhotos) error {

	return s.repo.Update(ctx, id, dto)
}

func (s *visitPhotosService) Delete(ctx context.Context, id uuid.UUID) error {

	vp, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if err := s.minio.RemoveObject(ctx, vp.PhotoKey); err != nil {
		return err
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	return nil
}

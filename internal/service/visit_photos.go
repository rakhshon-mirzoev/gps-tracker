package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/model"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/repository"
)

type VisitPhotosService interface {
	Create(ctx context.Context, vp *model.VisitPhotos) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.VisitPhotos, error)
	List(ctx context.Context, visitID uuid.UUID) ([]model.VisitPhotos, error)
	Update(ctx context.Context, vp *model.VisitPhotos) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type visitPhotosService struct {
	repo repository.VisitPhotosRepository
}

func NewVisitPhotosService(repo repository.VisitPhotosRepository) VisitPhotosService {
	return &visitPhotosService{repo: repo}
}

func (s *visitPhotosService) Create(ctx context.Context, vp *model.VisitPhotos) error {
	if vp.VisitID == uuid.Nil {
		return errors.New("visit_id is required")
	}
	if vp.PhotoURL == "" {
		return errors.New("photo_url is required")
	}

	return s.repo.Create(ctx, vp)
}

func (s *visitPhotosService) GetByID(ctx context.Context, id uuid.UUID) (*model.VisitPhotos, error) {
	if id == uuid.Nil {
		return nil, errors.New("id is required")
	}

	return s.repo.GetByID(ctx, id)
}

func (s *visitPhotosService) List(ctx context.Context, visitID uuid.UUID) ([]model.VisitPhotos, error) {
	if visitID == uuid.Nil {
		return nil, errors.New("visit_id is required")
	}

	return s.repo.List(ctx, visitID)
}

func (s *visitPhotosService) Update(ctx context.Context, vp *model.VisitPhotos) error {
	if vp.ID == uuid.Nil {
		return errors.New("id is required")
	}
	if vp.VisitID == uuid.Nil {
		return errors.New("visit_id is required")
	}
	if vp.PhotoURL == "" {
		return errors.New("photo_url is required")
	}

	return s.repo.Update(ctx, vp)
}

func (s *visitPhotosService) Delete(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("id is required")
	}

	return s.repo.Delete(ctx, id)
}

package service

import (
	"context"
	"errors"
	"log"

	"github.com/IkbolNazarov/gps-tracker/internal/model"
	"github.com/IkbolNazarov/gps-tracker/internal/repository"
	minioservice "github.com/IkbolNazarov/gps-tracker/pkg/storage/service"
	"github.com/google/uuid"
)

type VisitPhotos interface {
	Create(ctx context.Context, vp *model.VisitPhotos) error
	GetByVisitID(ctx context.Context, id uuid.UUID) ([]model.VisitPhotos, error)
	List(ctx context.Context, visitID uuid.UUID) ([]model.VisitPhotos, error)
	Update(ctx context.Context, id uuid.UUID, dto model.UpdateVisitPhotos) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByVisitIDs(ctx context.Context, ids []uuid.UUID) ([]model.VisitPhotos, error)
}

type visitPhotosService struct {
	repo  repository.VisitPhotosRepository
	minio minioservice.Minio
}

func NewVisitPhotosService(repo repository.VisitPhotosRepository, minio minioservice.Minio) VisitPhotos {
	return &visitPhotosService{repo: repo, minio: minio}
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

	if s.minio == nil {
		log.Fatal("s minio")
	}

	if err := s.minio.RemoveObject(ctx, vp.PhotoKey); err != nil {
		return err
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	return nil
}

func (s *visitPhotosService) GetByVisitIDs(ctx context.Context, ids []uuid.UUID) ([]model.VisitPhotos, error) {
	return s.repo.GetByVisitIDs(ctx, ids)
}

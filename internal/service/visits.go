package service

import (
	"context"
	"errors"
	"math"
	"mime/multipart"

	"github.com/google/uuid"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/model"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/repository"
	errs "github.com/rakhshon-mirzoev/gps-tracker/pkg/error"
	"github.com/rakhshon-mirzoev/gps-tracker/pkg/utils"
)

type Visits interface {
	Create(ctx context.Context, v *model.CreateVisit, user_id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Visit, error)
	List(ctx context.Context, scheduleID uuid.UUID) ([]model.Visit, error)
	Update(ctx context.Context, id string, dto model.UpdateVisit, file *multipart.Form, user_id uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type visitsService struct {
	repo  repository.VisitRepository
	store Stores
	photo VisitPhotos
}

func NewVisitsService(repo repository.VisitRepository) Visits {
	return &visitsService{repo: repo}
}

func (s *visitsService) Create(ctx context.Context, v *model.CreateVisit, user_id uuid.UUID) error {
	store, err := s.store.GetByID(ctx, v.StoreID)
	if err != nil {
		return err
	}

	m := Distance(store.Latitude, store.Longitude, v.Lat, v.Lon)
	if m > 150 {
		return errs.ErrorUnprocessable
	}

	visit := &model.Visit{
		StaffID:    user_id,
		StoreID:    v.StoreID,
		ScheduleID: &v.ScheduleID,
		CheckInLat: &v.Lat,
		CheckInLon: &v.Lon,
		CheckInAt:  &v.CheckInAt,
	}

	return s.repo.Create(ctx, visit)
}

func (s *visitsService) GetByID(ctx context.Context, id uuid.UUID) (*model.Visit, error) {
	if id == uuid.Nil {
		return nil, errors.New("id is required")
	}

	return s.repo.GetByID(ctx, id)
}

func (s *visitsService) List(ctx context.Context, scheduleID uuid.UUID) ([]model.Visit, error) {
	if scheduleID == uuid.Nil {
		return nil, errors.New("schedule_id is required")
	}

	return s.repo.List(ctx, scheduleID)
}

func (s *visitsService) Update(ctx context.Context, id string, dto model.UpdateVisit, file *multipart.Form, user_id uuid.UUID) error {
	uid, err := uuid.FromBytes([]byte(id))
	if err != nil {
		return err
	}

	selfie := file.File["selfie"]

	if len(selfie) == 0 || len(selfie) > 1 {
		return errs.ErrorInvalidFileFormat
	}

	for _, f := range selfie {
		if err := utils.ValidateFileType(f); err != nil {
			return errs.ErrorInvalidFileFormat
		}

		savePath, err := utils.PathCreator(f.Filename, "visit_photos", "selfie", uid)
		if err != nil {
			return err
		}

		dto.SelfieURL = &savePath

		vp := &model.VisitPhotos{
			VisitID:  uid,
			PhotoURL: savePath,
		}

		s.photo.Create(ctx, vp)
	}

	shelf := file.File["shelf"]
	if len(shelf) == 0 {
		return errs.ErrorEmptyRequestBody
	}

	for _, f := range shelf {
		if err := utils.ValidateFileType(f); err != nil {
			return errs.ErrorInvalidFileFormat
		}

		savePath, err := utils.PathCreator(f.Filename, "visit_photos", "shelf", uid)
		if err != nil {
			return err
		}

		vp := &model.VisitPhotos{
			VisitID:  uid,
			PhotoURL: savePath,
		}

		s.photo.Create(ctx, vp)
	}

	return s.repo.Update(ctx, uid, dto)
}

func (s *visitsService) Delete(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("id is required")
	}

	return s.repo.Delete(ctx, id)
}

const earthRadius = 6371000

func Distance(lat1, lon1, lat2, lon2 float64) float64 {
	lat1Rad := lat1 * math.Pi / 180
	lon1Rad := lon1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	lon2Rad := lon2 * math.Pi / 180

	dLat := lat2Rad - lat1Rad
	dLon := lon2Rad - lon1Rad

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(dLon/2)*math.Sin(dLon/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadius * c
}

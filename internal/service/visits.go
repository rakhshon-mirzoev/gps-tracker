package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/model"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/repository"
	"github.com/rakhshon-mirzoev/gps-tracker/pkg/constants"
	errs "github.com/rakhshon-mirzoev/gps-tracker/pkg/error"
	minioservice "github.com/rakhshon-mirzoev/gps-tracker/pkg/storage/service"
	"github.com/rakhshon-mirzoev/gps-tracker/pkg/utils"
)

type Visits interface {
	Create(ctx context.Context, v *model.CreateVisit, user_id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Visit, error)
	List(ctx context.Context, scheduleID uuid.UUID) ([]model.Visit, error)
	VisitComplete(ctx context.Context, id string, dto model.UpdateVisit) error
	VisitPhoto(ctx context.Context, id string, dto model.UpdateVisit, file *multipart.Form) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type visitsService struct {
	repo  repository.VisitRepository
	store Stores
	photo VisitPhotos
	minio minioservice.Minio
}

func NewVisitsService(repo repository.VisitRepository, minio minioservice.Minio) Visits {
	return &visitsService{repo: repo, minio: minio}
}

func (s *visitsService) Create(ctx context.Context, v *model.CreateVisit, user_id uuid.UUID) error {
	store, err := s.store.GetByID(ctx, v.StoreID)
	if err != nil {
		return err
	}

	m := utils.Distance(store.Latitude, store.Longitude, v.Lat, v.Lon)
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

func (s *visitsService) VisitComplete(ctx context.Context, id string, dto model.UpdateVisit) error {
	uid, err := uuid.FromBytes([]byte(id))
	if err != nil {
		return err
	}
	*dto.ReviewStatus = "pending"
	return s.repo.Update(ctx, uid, dto)
}

func (s *visitsService) VisitPhoto(ctx context.Context, id string, dto model.UpdateVisit, file *multipart.Form) error {
	uid, err := uuid.FromBytes([]byte(id))
	if err != nil {
		return err
	}

	selfie := file.File["selfie"]

	if len(selfie) == 0 || len(selfie) > 1 {
		return errs.ErrorInvalidFileFormat
	}

	if err := utils.ValidateFileType(selfie[0]); err != nil {
		return err
	}

	if err := utils.ValidateFileSize(selfie[0]); err != nil {
		return err
	}

	ext := utils.GetFileExtension(selfie[0].Filename)
	if ext == "" {
		return errs.ErrorFileExt
	}
	timestamp := time.Now().Unix()
	filename := fmt.Sprintf("%s_%d%s", strings.TrimSuffix(selfie[0].Filename, ext), timestamp, ext)
	objectName := fmt.Sprintf(constants.VisitDir+constants.SelfieDir+"/%s/%s", uid, filename)

	selfiePath, err := s.minio.UploadFile(ctx, selfie[0], objectName)
	if err != nil {
		return fmt.Errorf("failed to upload selfie to minio: %w", err)
	}

	localPath, err := utils.PathCreator(filename, constants.VisitDir, constants.SelfieDir, uid)
	if err == nil {
		if src, err := selfie[0].Open(); err == nil {
			defer src.Close()
			if dst, err := os.Create(localPath); err == nil {
				defer dst.Close()
				io.Copy(dst, src)
			}
		}
	}

	dto.SelfieURL = &selfiePath

	vp := &model.VisitPhotos{
		VisitID:      uid,
		PhotoURL:     selfiePath,
		OriginalName: selfie[0].Filename,
	}

	err = s.photo.Create(ctx, vp)
	if err != nil {
		return err
	}

	shelf := file.File["shelf"]
	if len(shelf) == 0 {
		return errs.ErrorEmptyRequestBody
	}

	for idx, f := range shelf {
		if err := utils.ValidateFileType(f); err != nil {
			return err
		}

		if err := utils.ValidateFileSize(f); err != nil {
			return err
		}

		ext := utils.GetFileExtension(f.Filename)
		if ext == "" {
			return errs.ErrorFileExt
		}
		timestamp := time.Now().Unix()
		filename := fmt.Sprintf("%s_%d%s", strings.TrimSuffix(f.Filename, ext), timestamp, ext)
		objectName := fmt.Sprintf(constants.VisitDir+constants.ShelfDir+"/%s/%s", uid, filename)

		photoPath, err := s.minio.UploadFile(ctx, f, objectName)
		if err != nil {
			return fmt.Errorf("failed to upload shelf photo %d to minio: %w", idx+1, err)
		}

		localPath, err := utils.PathCreator(filename, constants.VisitDir, constants.ShelfDir, uid)
		if err == nil {
			if src, err := f.Open(); err == nil {
				defer src.Close()
				if dst, err := os.Create(localPath); err == nil {
					defer dst.Close()
					io.Copy(dst, src)
				}
			}
		}

		vp = &model.VisitPhotos{
			PhotoURL:     photoPath,
			OriginalName: f.Filename,
			CapturedAt:   time.Now(),
		}

		err = s.photo.Create(ctx, vp)
		if err != nil {
			return err
		}
	}
	return s.repo.Update(ctx, uid, dto)
}

func (s *visitsService) Delete(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("id is required")
	}

	return s.repo.Delete(ctx, id)
}

package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/logger"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/model"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/repository"
	"github.com/rakhshon-mirzoev/gps-tracker/pkg/constants"
	errs "github.com/rakhshon-mirzoev/gps-tracker/pkg/error"
	minioservice "github.com/rakhshon-mirzoev/gps-tracker/pkg/storage/service"
	"github.com/rakhshon-mirzoev/gps-tracker/pkg/utils"
	"go.uber.org/zap"
)

type Visits interface {
	Create(ctx context.Context, v model.CreateVisit, user_id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Visit, error)
	List(ctx context.Context, scheduleID uuid.UUID) ([]model.Visit, error)
	ListByStaffID(ctx context.Context, staff_id uuid.UUID, lim, off int) ([]model.Visit, int, error)
	VisitComplete(ctx context.Context, id string, dto model.UpdateVisit, staff_id uuid.UUID) error
	VisitPhoto(ctx context.Context, id uuid.UUID, dto model.UpdateVisit, file *multipart.Form) error
	Delete(ctx context.Context, id uuid.UUID) error
	ListByReviewStatus(ctx context.Context, status string, limit, offset int, id uuid.UUID) ([]model.VisitResponse, int, error)
	UpdateVisit(ctx context.Context, id uuid.UUID, dto model.UpdateVisit, user_id uuid.UUID, role, ip, settings, action string) error
	GetDataForExcel(ctx context.Context, staff_id, store_id *uuid.UUID, begin, end *time.Time) ([]model.VisitExcel, error)
}
type visitsService struct {
	repo          repository.VisitRepository
	store         Stores
	photo         VisitPhotos
	minio         minioservice.Minio
	notifications Notification
	staff         StaffUsers
	audit         AuditLogs
	group         Group
	schedule      Schedule
}

func NewVisitsService(repo repository.VisitRepository, minio minioservice.Minio, store Stores, photo VisitPhotos, staff StaffUsers, notifications Notification, audit AuditLogs, schedule Schedule, group Group) Visits {
	return &visitsService{repo: repo, minio: minio, store: store, photo: photo, staff: staff, notifications: notifications, audit: audit, group: group, schedule: schedule}
}

func (s *visitsService) Create(ctx context.Context, vCreate model.CreateVisit, user_id uuid.UUID) error {
	storeID, err := uuid.Parse(vCreate.StoreID)
	if err != nil {
		return err
	}

	scheduleID, err := uuid.Parse(vCreate.ScheduleID)
	if err != nil {
		return err
	}

	store, err := s.store.GetByID(ctx, storeID)
	if err != nil {
		return err
	}

	v, err := s.repo.GetByStoreScheduleId(ctx, storeID, scheduleID, user_id)
	if err != nil {
		return err
	}

	if v != nil && v.CheckOutAt != nil {
		return errs.ErrorAlreadyReviewed
	}

	staff, err := s.staff.GetByID(ctx, user_id)
	if err != nil {
		return err
	}

	group, err := s.group.GetByID(ctx, staff.GroupID)
	if err != nil {
		return err
	}

	distance := utils.Distance(store.Latitude, store.Longitude, vCreate.Lat, vCreate.Lon)
	if distance > constants.MaxDistanceMeters {
		return errs.ErrorUnprocessable
	}

	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	now := time.Now().UTC()

	visit := model.Visit{
		StaffID:           user_id,
		StoreID:           storeID,
		ScheduleID:        &scheduleID,
		CheckInLat:        &vCreate.Lat,
		CheckInLon:        &vCreate.Lon,
		CheckInAt:         &now,
		DistanceDeviation: int(distance),
		ReviewerID:        group.SupervisorID,
	}

	if err := s.repo.CreateTx(ctx, tx, &visit); err != nil {
		return err
	}

	meta := map[string]any{"type": "visit_check_in"}

	b, err := json.Marshal(meta)
	if err != nil {
		return err
	}

	n := &model.Notification{
		UserID:   &user_id,
		Title:    "Задание принято в работу",
		Body:     fmt.Sprintf("Магазин %s, время начала %s.", store.Name, now),
		Type:     "visit_check_in",
		Meta:     b,
		IsRead:   false,
		PushSent: false,
	}
	if err := s.notifications.CreateTx(
		ctx,
		tx,
		n,
	); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
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

func (s *visitsService) VisitComplete(ctx context.Context, id string, dto model.UpdateVisit, staff_id uuid.UUID) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	pending := "pending"
	now := time.Now().UTC()
	dto.ReviewStatus = &pending
	if err := s.repo.Update(ctx, uid, dto); err != nil {
		return err
	}

	schedule_id, err := uuid.Parse(*dto.ScheduleID)
	if err != nil {
		return err
	}

	err = s.schedule.Update(ctx, schedule_id, model.UpdateSchedule{
		Status: &pending,
	})

	user, err := s.staff.GetByID(ctx, staff_id)
	if err != nil {
		return err
	}

	vSt, err := s.GetByID(ctx, uid)
	if err != nil {
		return err
	}

	store, err := s.store.GetByID(ctx, vSt.StoreID)
	if err != nil {
		return err
	}

	if err := s.notifications.CreateAndPushToUser(ctx, staff_id, "Задание отправлено", fmt.Sprintf("Магазин %s, время окончания %s.", store.Name, now), "visit_complete", map[string]any{"type": "visit_complete"}, user.FcmToken); err != nil {
		logger.Log.Error(
			fmt.Sprintf("%s: failed operation", "CreateAndPushToUser in visit complete"),
			zap.Error(err),
		)
	}

	return nil
}

func (s *visitsService) VisitPhoto(ctx context.Context, id uuid.UUID, dto model.UpdateVisit, file *multipart.Form) error {

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
	objectName := fmt.Sprintf(constants.VisitDir+constants.SelfieDir+"/%s/%s", id, filename)

	selfiePath, err := s.minio.UploadFile(ctx, selfie[0], objectName)
	if err != nil {
		return fmt.Errorf("failed to upload selfie to minio: %w", err)
	}

	localPath, err := utils.PathCreator(filename, constants.VisitDir, constants.SelfieDir, id)
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
		VisitID:      id,
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
		objectName := fmt.Sprintf(constants.VisitDir+constants.ShelfDir+"/%s/%s", id, filename)

		photoPath, err := s.minio.UploadFile(ctx, f, objectName)
		if err != nil {
			return fmt.Errorf("failed to upload shelf photo %d to minio: %w", idx+1, err)
		}

		localPath, err := utils.PathCreator(filename, constants.VisitDir, constants.ShelfDir, id)
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
	return s.repo.Update(ctx, id, dto)
}

func (s *visitsService) Delete(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("id is required")
	}

	return s.repo.Delete(ctx, id)
}

func (s *visitsService) ListByReviewStatus(ctx context.Context, status string, limit, offset int, id uuid.UUID) ([]model.VisitResponse, int, error) {
	return s.repo.ListByReviewStatus(ctx, status, limit, offset, id)
}

func (s *visitsService) UpdateVisit(ctx context.Context, id uuid.UUID, dto model.UpdateVisit, user_id uuid.UUID, role, ip, settings, action string) error {
	if err := s.repo.Update(ctx, id, dto); err != nil {
		return err
	}

	v, err := s.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if err := s.schedule.Update(ctx, *v.ScheduleID, model.UpdateSchedule{
		Status: dto.ReviewStatus,
	}); err != nil {
		return err
	}

	u, err := s.staff.GetByID(ctx, v.StaffID)
	if err != nil {
		return err
	}

	if u.FcmToken != nil {
		if err := s.notifications.CreateAndPushToUser(ctx, u.ID, "visit review result",
			fmt.Sprintf("visit in store result: %s; review comment: %s; reviewed at: %v", v.ReviewStatus, *v.ReviewComment, v.ReviewedAt), "review_result",
			map[string]any{"type": "review_result"}, u.FcmToken); err != nil {
			return err
		}
	}

	if err := s.audit.Create(ctx, &model.AuditLogs{
		ActorID:     user_id,
		ActorType:   role,
		IPAddress:   &ip,
		TargetTable: &settings,
		TargetID:    &id,
		Action:      action,
	}); err != nil {
		return err
	}

	return nil
}

func (s *visitsService) GetDataForExcel(ctx context.Context, staff_id, store_id *uuid.UUID, begin, end *time.Time) ([]model.VisitExcel, error) {
	return s.repo.GetDataForExcel(ctx, staff_id, store_id, begin, end)
}

func (s *visitsService) ListByStaffID(ctx context.Context, staff_id uuid.UUID, lim, off int) ([]model.Visit, int, error) {
	return s.repo.ListByStaffID(ctx, staff_id, lim, off)
}

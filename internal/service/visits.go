package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"strings"
	"time"

	"github.com/IkbolNazarov/gps-tracker/internal/model"
	"github.com/IkbolNazarov/gps-tracker/internal/repository"
	"github.com/IkbolNazarov/gps-tracker/pkg/constants"
	errs "github.com/IkbolNazarov/gps-tracker/pkg/error"
	minioservice "github.com/IkbolNazarov/gps-tracker/pkg/storage/service"
	"github.com/IkbolNazarov/gps-tracker/pkg/timezone"
	"github.com/IkbolNazarov/gps-tracker/pkg/utils"
	"github.com/google/uuid"
)

type Visits interface {
	Create(ctx context.Context, v model.CreateVisit, user_id uuid.UUID) (*model.Visit, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.Visit, error)
	List(ctx context.Context, scheduleID uuid.UUID) ([]model.Visit, error)
	ListByStaffID(ctx context.Context, staff_id uuid.UUID, lim, off int, begin, end time.Time) ([]model.Visit, int, error)
	VisitComplete(ctx context.Context, id uuid.UUID, dto model.UpdateVisit, staff_id uuid.UUID) error
	VisitPhoto(ctx context.Context, id uuid.UUID, file *multipart.Form) error
	Delete(ctx context.Context, id uuid.UUID) error
	ListByReviewStatus(ctx context.Context, status string, limit, offset int, id uuid.UUID, role string) ([]model.VisitResponse, int, error)
	UpdateVisit(ctx context.Context, id uuid.UUID, dto model.UpdateVisit, user_id uuid.UUID, role, ip, settings, action string) error
	GetDataForExcel(ctx context.Context, staff_id, store_id *uuid.UUID, begin, end *time.Time) ([]model.VisitExcel, error)
	GetByIDPhoto(ctx context.Context, id uuid.UUID) (*model.VisitWithPhoto, error)
	Update(ctx context.Context, id uuid.UUID, dto model.UpdateVisit) error
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

func (s *visitsService) Create(ctx context.Context, vCreate model.CreateVisit, user_id uuid.UUID) (*model.Visit, error) {
	storeID, err := uuid.Parse(vCreate.StoreID)
	if err != nil {
		return nil, err
	}

	scheduleID, err := uuid.Parse(vCreate.ScheduleID)
	if err != nil {
		return nil, err
	}

	schedule, err := s.schedule.GetByID(ctx, scheduleID)
	if err != nil {
		return nil, err
	}

	if schedule.Status == "initialized" {
		return nil, errs.ErrorScheduleIsNotApproved
	}

	store, err := s.store.GetByID(ctx, storeID)
	if err != nil {
		return nil, err
	}

	v, err := s.repo.GetByStoreScheduleId(ctx, storeID, scheduleID, user_id)
	if err != nil {
		return nil, err
	}

	if v != nil && v.CheckOutAt != nil && v.ReviewStatus != "rejected" && v.ReviewStatus != "revision" {
		return nil, errs.ErrorAlreadyReviewed
	}

	staff, err := s.staff.GetByID(ctx, user_id)
	if err != nil {
		return nil, err
	}

	group, err := s.group.GetByID(ctx, staff.GroupID)
	if err != nil {
		return nil, err
	}

	distance := utils.Distance(store.Latitude, store.Longitude, vCreate.Lat, vCreate.Lon)
	if distance > constants.MaxDistanceMeters {
		return nil, errs.ErrorUnprocessable
	}

	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	now := time.Now().In(timezone.Dushanbe)
	formatted := now.Format("2006-01-02 15:04:05")

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
		return nil, err
	}

	meta := map[string]any{"type": "visit_check_in"}

	b, err := json.Marshal(meta)
	if err != nil {
		return nil, err
	}

	n := &model.Notification{
		UserID:   &user_id,
		Title:    "Задание принято в работу",
		Body:     fmt.Sprintf("Магазин %s, время начала %s.", store.Name, formatted),
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
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &visit, nil
}

func (s *visitsService) GetByIDPhoto(ctx context.Context, id uuid.UUID) (*model.VisitWithPhoto, error) {

	photos, err := s.photo.GetByVisitID(ctx, id)
	if err != nil {
		return nil, err
	}

	for i := range photos {
		url, err := s.minio.PresignedDownload(ctx, photos[i].PhotoKey, time.Minute*15)
		if err != nil {
			return nil, err
		}
		if strings.Contains(photos[i].PhotoKey, "selfie") {
			photos[i].PhotoType = "selfie"
		} else {
			photos[i].PhotoType = "shelf"
		}

		photos[i].PhotoUrl = url
	}

	visit, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	sch, err := s.schedule.GetByID(ctx, *visit.ScheduleID)
	if err != nil {
		return nil, err
	}

	requirement := model.Required{
		IsPhotoRequired:   *sch.IsPhotoRequired,
		IsSelfieRequired:  *sch.IsSelfieRequired,
		IsCommentRequired: *sch.IsCommentRequired,
	}

	return &model.VisitWithPhoto{
		Visit:    *visit,
		Photo:    photos,
		Required: requirement,
	}, nil
}

func (s *visitsService) GetByID(ctx context.Context, id uuid.UUID) (*model.Visit, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *visitsService) List(ctx context.Context, scheduleID uuid.UUID) ([]model.Visit, error) {
	return s.repo.List(ctx, scheduleID)
}

func (s *visitsService) VisitComplete(ctx context.Context, uid uuid.UUID, dto model.UpdateVisit, staff_id uuid.UUID) error {
	pending := "pending"
	dto.ReviewStatus = &pending
	if err := s.repo.Update(ctx, uid, dto); err != nil {
		return err
	}

	vSt, err := s.GetByID(ctx, uid)
	if err != nil {
		return err
	}

	if err := s.schedule.Update(ctx, *vSt.ScheduleID, model.UpdateSchedule{
		Status: &pending,
	}); err != nil {
		return err
	}

	return nil
}

func (s *visitsService) VisitPhoto(ctx context.Context, id uuid.UUID, file *multipart.Form) error {

	now := time.Now().In(timezone.Dushanbe)

	if file == nil || file.File == nil {
		return errs.ErrorInvalidFileFormat
	}

	if len(file.File) == 0 {
		return errs.ErrorInvalidFileFormat
	}

	selfie := file.File["selfie"]

	if len(selfie) > 1 {
		return errs.ErrorInvalidFileFormat
	}

	if len(selfie) != 0 {
		if selfie[0] == nil {
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
		timestamp := time.Now().In(timezone.Dushanbe).Unix()
		filename := fmt.Sprintf("%s_%d%s", uuid.New(), timestamp, ext)
		objectName := fmt.Sprintf(constants.VisitDir+"/"+constants.SelfieDir+"/%s/%s", id, filename)

		selfieKey, err := s.minio.UploadFile(ctx, selfie[0], objectName)
		if err != nil {
			return fmt.Errorf("failed to upload selfie to minio: %w", err)
		}

		vp := &model.VisitPhotos{
			VisitID:      id,
			PhotoKey:     selfieKey,
			OriginalName: selfie[0].Filename,
			CapturedAt:   now,
		}

		err = s.photo.Create(ctx, vp)
		if err != nil {
			return err
		}
	}

	shelf := file.File["shelf"]
	if len(shelf) != 0 {
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
			timestamp := now.Unix()
			filename := fmt.Sprintf("%s_%d%s", uuid.New(), timestamp, ext)
			objectName := fmt.Sprintf(constants.VisitDir+"/"+constants.ShelfDir+"/%s/%s", id, filename)

			photoKey, err := s.minio.UploadFile(ctx, f, objectName)
			if err != nil {
				return fmt.Errorf("failed to upload shelf photo %d to minio: %w", idx+1, err)
			}

			vp := &model.VisitPhotos{
				PhotoKey:     photoKey,
				VisitID:      id,
				OriginalName: f.Filename,
				CapturedAt:   now,
			}

			err = s.photo.Create(ctx, vp)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *visitsService) Delete(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("id is required")
	}

	return s.repo.Delete(ctx, id)
}

func (s *visitsService) ListByReviewStatus(ctx context.Context, status string, limit, offset int, id uuid.UUID, role string) ([]model.VisitResponse, int, error) {
	isAdmin := false
	if role == "admin" {
		isAdmin = true
		id = uuid.Nil
	}

	visits, total, err := s.repo.ListByReviewStatus(ctx, status, limit, offset, id, isAdmin)
	if err != nil {
		return nil, 0, err
	}

	if len(visits) == 0 {
		return visits, total, nil
	}

	ids := make([]uuid.UUID, 0, len(visits))
	for _, v := range visits {
		ids = append(ids, v.ID)
	}

	photos, err := s.photo.GetByVisitIDs(ctx, ids)
	if err != nil {
		return nil, 0, err
	}

	photoMap := make(map[uuid.UUID][]string)

	for _, p := range photos {

		url, err := s.minio.PresignedDownload(ctx, p.PhotoKey, time.Minute*15)
		if err != nil {
			return nil, 0, err
		}

		if strings.Contains(p.PhotoKey, "shelf") {
			photoMap[p.VisitID] = append(photoMap[p.VisitID], url)
		}
	}

	for i := range visits {
		visits[i].ShelfPhotos = photoMap[visits[i].ID]
	}

	return visits, total, nil
}

func (s *visitsService) UpdateVisit(ctx context.Context, id uuid.UUID, dto model.UpdateVisit, user_id uuid.UUID, role, ip, settings, action string) error {
	now := time.Now().In(timezone.Dushanbe)

	dto.ReviewedAt = &now

	v, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	schedule, err := s.schedule.GetByID(ctx, *v.ScheduleID)
	if err != nil {
		return err
	}

	if schedule.IsDeleted == true {
		scheduleDeleted := "schedule_deleted"
		if err := s.repo.Update(ctx, id, model.UpdateVisit{
			ReviewStatus: &scheduleDeleted,
			ReviewedAt:   schedule.DeletedAt,
		}); err != nil {
			return err
		}

		return nil
	}

	if err := s.repo.Update(ctx, id, dto); err != nil {
		return err
	}

	scheduleStatus := dto.ReviewStatus
	if scheduleStatus != nil && *scheduleStatus == "revision" {
		pending := "pending"
		scheduleStatus = &pending
	}

	if err := s.schedule.Update(ctx, *v.ScheduleID, model.UpdateSchedule{
		Status:        scheduleStatus,
		ReviewComment: dto.ReviewComment,
	}); err != nil {
		return err
	}

	u, err := s.staff.GetByID(ctx, v.StaffID)
	if err != nil {
		return err
	}

	if u.FcmToken != nil {
		body := ""
		status := ""

		reviewedAt := ""
		if v.ReviewedAt != nil {
			reviewedAt = v.ReviewedAt.Format("02.01.2006, 15:04")
		}

		switch *dto.ReviewStatus {
		case "approved":
			status = "Одобрено"
		case "revision":
			status = "Требуется исправление"
		case "rejected":
			status = "Отклонено"
		}

		if dto.ReviewComment != nil && status != "Одобрено" {
			body = fmt.Sprintf(
				"Результат проверки визита: %s\nКомментарий: %s\nПроверено: %s",
				status,
				*dto.ReviewComment,
				reviewedAt,
			)
		} else {
			body = fmt.Sprintf(
				"Результат проверки визита: %s\nПроверено: %s",
				status,
				reviewedAt,
			)
		}

		if err := s.notifications.CreateAndPushToUser(ctx, u.ID, "Результаты проверки визита", body, "review_result", map[string]any{"type": "review_result"}, u.FcmToken); err != nil {
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
		CreatedAt:   time.Now().In(timezone.Dushanbe),
	}); err != nil {
		return err
	}

	return nil
}

func (s *visitsService) GetDataForExcel(ctx context.Context, staff_id, store_id *uuid.UUID, begin, end *time.Time) ([]model.VisitExcel, error) {
	return s.repo.GetDataForExcel(ctx, staff_id, store_id, begin, end)
}

func (s *visitsService) ListByStaffID(ctx context.Context, staff_id uuid.UUID, lim, off int, begin, end time.Time) ([]model.Visit, int, error) {
	return s.repo.ListByStaffID(ctx, staff_id, lim, off, begin, end)
}

func (s *visitsService) Update(ctx context.Context, id uuid.UUID, dto model.UpdateVisit) error {
	v, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if v.ReviewStatus == "rejected" {
		return errs.ErrorAlreadyReviewed
	}

	store, err := s.store.GetByID(ctx, v.StoreID)
	if err != nil {
		return err
	}

	distance := utils.Distance(store.Latitude, store.Longitude, *dto.CheckInLat, *dto.CheckInLon)
	if distance > constants.MaxDistanceMeters {
		return errs.ErrorUnprocessable
	}

	if v.ReviewStatus != "pending" && v.ReviewStatus != "revision" {
		return errs.ErrorAlreadyReviewed
	}

	// Если визит был на доработке — при сохранении сотрудником автоматически возвращаем в pending
	if v.ReviewStatus == "revision" {
		pending := "pending"
		dto.ReviewStatus = &pending
	}

	return s.repo.Update(ctx, id, dto)
}

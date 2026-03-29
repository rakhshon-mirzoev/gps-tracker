package service

import (
	"context"
	"errors"
	"time"

	"github.com/IkbolNazarov/gps-tracker/internal/events"
	"github.com/IkbolNazarov/gps-tracker/internal/logger"
	"github.com/IkbolNazarov/gps-tracker/internal/model"
	"github.com/IkbolNazarov/gps-tracker/internal/repository"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type LocationHistory interface {
	Create(ctx context.Context, lh *model.LocationHistory) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.LocationHistory, error)
	List(ctx context.Context, staffID uuid.UUID) ([]model.LocationHistory, error)
	Update(ctx context.Context, id uuid.UUID, dto model.UpdateLocationHistory) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetAllStaffOnlineOffline(ctx context.Context, supervisorID uuid.UUID, begin, end time.Time) (map[uuid.UUID]map[string]time.Duration, error)
	CreateBatch(ctx context.Context, locations []model.LocationHistory) error
	LatestByGroupID(ctx context.Context, groupID uuid.UUID) ([]events.LocationUpdate, error)
	LatestAll(ctx context.Context) ([]events.LocationUpdate, error)
}

type locationHistoryService struct {
	r         repository.LocationHistoryRepository
	alertRepo repository.GpsAlertStatesRepository
	staff     StaffUsers
	publisher events.LocationPublisher
}

func NewLocationHistoryService(
	r repository.LocationHistoryRepository,
	alertRepo repository.GpsAlertStatesRepository,
	staff StaffUsers,
	publisher events.LocationPublisher,
) LocationHistory {
	return &locationHistoryService{r: r, alertRepo: alertRepo, staff: staff, publisher: publisher}
}

func (s *locationHistoryService) Create(ctx context.Context, lh *model.LocationHistory) error {
	return s.r.Create(ctx, lh)
}

func (s *locationHistoryService) CreateBatch(ctx context.Context, locations []model.LocationHistory) error {
	err := s.r.CreateBatch(ctx, locations)
	if err != nil {
		return err
	}

	var latest *model.LocationHistory
	for i := len(locations) - 1; i >= 0; i-- {
		if *locations[i].Status == "online" {
			latest = &locations[i]
			break
		}
	}

	if latest != nil {
		staff, err := s.staff.GetByID(ctx, latest.StaffID)
		if err != nil {
			logger.Log.Error("err while staff get by id", zap.Error(err))
		}

		// Сотрудник прислал GPS — сбрасываем активный GPS-алерт,
		// чтобы при следующей "пропаже" уведомление отправилось снова.
		if err := s.alertRepo.Resolve(ctx, latest.StaffID); err != nil {
			logger.Log.Warn("location history: resolve gps alert",
				zap.Stringer("staff_id", latest.StaffID),
				zap.Error(err),
			)
		}

		update := events.LocationUpdate{
			StaffID:      latest.StaffID,
			GroupID:      staff.GroupID,
			FullName:     staff.FullName,
			Latitude:     latest.Latitude,
			Longitude:    latest.Longitude,
			BatteryLevel: latest.BatteryLevel,
			CapturedAt:   latest.CapturedAt,
			Status:       latest.Status,
		}

		_ = s.publisher.PublishOnlineLocation(ctx, update)
	}

	return nil
}

func (s *locationHistoryService) GetByID(ctx context.Context, id uuid.UUID) (*model.LocationHistory, error) {
	return s.r.GetByID(ctx, id)
}

func (s *locationHistoryService) List(ctx context.Context, staffID uuid.UUID) ([]model.LocationHistory, error) {
	if staffID == uuid.Nil {
		return nil, errors.New("staff_id is required")
	}

	return s.r.List(ctx, staffID)
}

func (s *locationHistoryService) Update(ctx context.Context, id uuid.UUID, dto model.UpdateLocationHistory) error {
	return s.r.Update(ctx, id, dto)
}

func (s *locationHistoryService) Delete(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("id is required")
	}

	return s.r.Delete(ctx, id)
}

func (s *locationHistoryService) LatestByGroupID(ctx context.Context, groupID uuid.UUID) ([]events.LocationUpdate, error) {
	rows, err := s.r.LatestByGroupID(ctx, groupID)
	if err != nil {
		return nil, err
	}

	return toLocationUpdates(rows), nil
}

func (s *locationHistoryService) LatestAll(ctx context.Context) ([]events.LocationUpdate, error) {
	rows, err := s.r.LatestAll(ctx)
	if err != nil {
		return nil, err
	}

	return toLocationUpdates(rows), nil
}

func toLocationUpdates(rows []model.LocationHistoryWithStaff) []events.LocationUpdate {
	updates := make([]events.LocationUpdate, 0, len(rows))
	for _, r := range rows {
		updates = append(updates, events.LocationUpdate{
			StaffID:      r.StaffID,
			GroupID:      r.GroupID,
			FullName:     r.FullName,
			PhoneNumber:  r.PhoneNumber,
			Latitude:     r.Latitude,
			Longitude:    r.Longitude,
			BatteryLevel: r.BatteryLevel,
			CapturedAt:   r.CapturedAt,
			Status:       r.Status,
		})
	}
	return updates
}

func (s *locationHistoryService) GetAllStaffOnlineOffline(ctx context.Context, supervisorID uuid.UUID, begin, end time.Time) (map[uuid.UUID]map[string]time.Duration, error) {
	return s.r.GetAllStaffOnlineOffline(ctx, supervisorID, begin, end)
}

package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/rakhshon-mirzoev/gps-tracker/pkg/db"
)

type Tx interface {
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type Repository struct {
	DB                  *db.DB
	SystemUsers         SystemUsersRepository
	Groups              GroupsRepository
	StaffUsers          StaffUsersRepository
	Notifications       NotificationsRepository
	Auth                Auth
	Schedules           ScheduleRepository
	Visit               VisitRepository
	Store               StoreRepository
	Location            LocationHistoryRepository
	SystemSettings      SystemSettingsRepository
	StaffStoresMapping  StaffStoresMappingRepository
	Audit               AuditLogsRepository
	VisitPhotos         VisitPhotosRepository
	DailyReportingStats DailyReportingStatsRepository
	StorePhotos         StorePhotosRepository
}

func New(db *db.DB) *Repository {
	return &Repository{
		DB:                  db,
		SystemUsers:         NewSystemUsersRepository(db),
		Groups:              NewGroupsRepository(db),
		StaffUsers:          NewStaffUsersRepository(db),
		Notifications:       NewNotificationsRepository(db),
		Auth:                NewAuthRepository(db),
		Schedules:           NewScheduleRepository(db),
		Visit:               NewVisitRepository(db),
		Store:               NewStoreRepository(db),
		Location:            NewLocationHistoryRepository(db),
		SystemSettings:      NewSystemSettingsRepository(db),
		StaffStoresMapping:  NewStaffStoresMappingRepository(db),
		Audit:               NewAuditLogsRepository(db),
		VisitPhotos:         NewVisitPhotosRepository(db),
		DailyReportingStats: NewDailyReportingStatsRepository(db),
		StorePhotos:         NewStorePhotosRepository(db),
	}
}

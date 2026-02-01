package repository

import "github.com/rakhshon-mirzoev/gps-tracker/pkg/db"

type Repository struct {
	DB                 *db.DB
	SystemUsers        SystemUsersRepository
	Groups             GroupsRepository
	StaffUsers         StaffUsersRepository
	Auth               Auth
	Schedules          ScheduleRepository
	Visit              VisitRepository
	Store              StoreRepository
	Location           LocationHistoryRepository
	SystemSettings     SystemSettingsRepository
	StaffStoresMapping StaffStoresMappingRepository
}

func New(db *db.DB) *Repository {
	return &Repository{
		DB:                 db,
		SystemUsers:        NewSystemUsersRepository(db),
		Groups:             NewGroupsRepository(db),
		StaffUsers:         NewStaffUsersRepository(db),
		Auth:               NewAuthRepository(db),
		Schedules:          NewScheduleRepository(db),
		Visit:              NewVisitRepository(db),
		Store:              NewStoreRepository(db),
		Location:           NewLocationHistoryRepository(db),
		SystemSettings:     NewSystemSettingsRepository(db),
		StaffStoresMapping: NewStaffStoresMappingRepository(db),
	}
}

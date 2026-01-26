package repository

import "github.com/rakhshon-mirzoev/gps-tracker/pkg/db"

type Repository struct {
	DB          *db.DB
	SystemUsers SystemUsersRepository
	Groups      GroupsRepository
	StaffUsers  StaffUsersRepository
	Auth        Auth
	Schedules   ScheduleRepository
	Visit       VisitRepository
}

func New(db *db.DB) *Repository {
	return &Repository{
		DB:          db,
		SystemUsers: NewSystemUsersRepository(db),
		Groups:      NewGroupsRepository(db),
		StaffUsers:  NewStaffUsersRepository(db),
		Auth:        NewAuthRepository(db),
		Schedules:   NewScheduleRepository(db),
		Visit:       NewVisitRepository(db),
	}
}

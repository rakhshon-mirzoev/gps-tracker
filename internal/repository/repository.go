package repository

import "github.com/rakhshon-mirzoev/gps-tracker/pkg/db"

type Repository struct {
	DB          *db.DB
	SystemUsers SystemUsersRepository
	Groups      GroupsRepository
	StaffUsers  StaffUsersRepository
}

func InitRepository(db *db.DB) *Repository {
	return &Repository{
		DB:          db,
		SystemUsers: InitSystemUsersRepository(db),
		Groups:      InitGroupsRepository(db),
		StaffUsers:  InitStaffUsersRepository(db),
	}
}

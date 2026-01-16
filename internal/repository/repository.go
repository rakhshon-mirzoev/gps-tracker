package repository

import "github.com/rakhshon-mirzoev/gps-tracker/pkg/db"

type Repository struct {
	DB          *db.DB
	SystemUsers SystemUsersRepository
	Groups      GroupsRepository
	StaffUsers  StaffUsersRepository
	Auth        Auth
}

func InitRepository(db *db.DB) *Repository {
	return &Repository{
		DB:          db,
		SystemUsers: InitSystemUsersRepository(db),
		Groups:      InitGroupsRepository(db),
		StaffUsers:  InitStaffUsersRepository(db),
		Auth:        InitAuthRepository(db),
	}
}

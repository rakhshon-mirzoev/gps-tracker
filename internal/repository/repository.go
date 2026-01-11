package repository

import "github.com/rakhshon-mirzoev/gps-tracker/pkg/db"

type Repository struct {
	DB *db.DB
}

func InitRepository(db *db.DB) *Repository {
	return &Repository{DB: db}
}

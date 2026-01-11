package service

import "github.com/rakhshon-mirzoev/gps-tracker/internal/repository"

type Service struct {
	R *repository.Repository
}

func InitService(r *repository.Repository) *Service {
	return &Service{R: r}
}

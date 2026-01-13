package service

import "github.com/rakhshon-mirzoev/gps-tracker/internal/repository"

type Service struct {
	SystemUsers  *SystemUsers
	GroupService *GroupService
	StaffUser    *StaffUserService
}

func InitService(r *repository.Repository) *Service {
	return &Service{
		SystemUsers:  InitSystemUsersService(r.SystemUsers),
		GroupService: InitGroupService(r.Groups),
		StaffUser:    InitStaffUserService(r.StaffUsers),
	}
}

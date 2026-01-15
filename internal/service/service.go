package service

import (
	"github.com/rakhshon-mirzoev/gps-tracker/internal/repository"
	"github.com/rakhshon-mirzoev/gps-tracker/pkg/token"
)

type Service struct {
	SystemUsers  SystemUsers
	GroupService Group
	StaffUser    StaffUsers
	Authorize    Authorize
}

func InitService(r *repository.Repository, t token.Token) *Service {
	return &Service{
		SystemUsers:  InitSystemUsersService(r.SystemUsers),
		GroupService: InitGroupService(r.Groups),
		StaffUser:    InitStaffUserService(r.StaffUsers),
		Authorize:    InitAuthorizeService(t),
	}
}

package service

import (
	"github.com/rakhshon-mirzoev/gps-tracker/internal/repository"
	"github.com/rakhshon-mirzoev/gps-tracker/pkg/token"
)

type Service struct {
	SystemUsers SystemUsers
	Group       Group
	StaffUser   StaffUsers
	Authorize   Authorize
}

func InitService(r *repository.Repository, t token.Token) *Service {
	return &Service{
		SystemUsers: InitSystemUsersService(r.SystemUsers),
		Group:       InitGroupService(r.Groups),
		StaffUser:   InitStaffUserService(r.StaffUsers),
		Authorize:   InitAuthorizeService(t, r.Auth, r.SystemUsers, r.StaffUsers),
	}
}

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
	sys := InitSystemUsersService(r.SystemUsers)
	group := InitGroupService(r.Groups)
	staff := InitStaffUserService(r.StaffUsers)
	authorize := InitAuthorizeService(t, r.Auth, sys, staff)
	return &Service{
		SystemUsers: sys,
		Group:       group,
		StaffUser:   staff,
		Authorize:   authorize,
	}
}

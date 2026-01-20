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

func New(r *repository.Repository, t token.Token) *Service {
	sys := NewSystemUsersService(r.SystemUsers)
	group := NewGroupService(r.Groups)
	staff := NewStaffUserService(r.StaffUsers)
	authorize := NewAuthorizeService(t, r.Auth, sys, staff)
	return &Service{
		SystemUsers: sys,
		Group:       group,
		StaffUser:   staff,
		Authorize:   authorize,
	}
}

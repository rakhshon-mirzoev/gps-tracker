package service

import (
	"github.com/rakhshon-mirzoev/gps-tracker/internal/repository"
	"github.com/rakhshon-mirzoev/gps-tracker/pkg/token"
)

type Service struct {
	SystemUsers        SystemUsers
	Group              Group
	StaffUser          StaffUsers
	Authorize          Authorize
	Schedule           Schedule
	StaffStoresMapping StaffStoresMapping
	Visit              Visits
}

func New(r *repository.Repository, t token.Token) *Service {
	sch := NewScheduleService(r.Schedules)
	sys := NewSystemUsersService(r.SystemUsers)
	group := NewGroupService(r.Groups)
	staff := NewStaffUserService(r.StaffUsers, sch)
	authorize := NewAuthorizeService(t, r.Auth, sys, staff)
	visit := NewVisitsService(r.Visit)
	return &Service{
		SystemUsers: sys,
		Group:       group,
		StaffUser:   staff,
		Authorize:   authorize,
		Schedule:    sch,
		Visit:       visit,
	}
}

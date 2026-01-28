package service

import (
	"github.com/rakhshon-mirzoev/gps-tracker/internal/repository"
	minioservice "github.com/rakhshon-mirzoev/gps-tracker/pkg/storage/service"
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
	Store              Stores
	Location           LocationHistory
	SystemSettings     SystemSettings
}

func New(r *repository.Repository, t token.Token, minio minioservice.Minio) *Service {
	sch := NewScheduleService(r.Schedules)
	sys := NewSystemUsersService(r.SystemUsers)
	group := NewGroupService(r.Groups)
	staff := NewStaffUserService(r.StaffUsers, sch)
	authorize := NewAuthorizeService(t, r.Auth, sys, staff)
	visit := NewVisitsService(r.Visit, minio)
	store := NewStoresService(r.Store)
	location := NewLocationHistoryService(r.Location)
	ss := NewSystemSettingsService(r.SystemSettings)
	return &Service{
		SystemUsers:    sys,
		Group:          group,
		StaffUser:      staff,
		Authorize:      authorize,
		Schedule:       sch,
		Visit:          visit,
		Store:          store,
		Location:       location,
		SystemSettings: ss,
	}
}

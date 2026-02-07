package service

import (
	"github.com/rakhshon-mirzoev/gps-tracker/internal/repository"
	"github.com/rakhshon-mirzoev/gps-tracker/pkg/fcm"
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
	Excel              Excel
	Audit              AuditLogs
	Notification       Notification
	VisitPhotos        VisitPhotos
	StorePhotos        StorePhotos
}

func New(r *repository.Repository, t token.Token, minio minioservice.Minio, sender fcm.PushSender) *Service {
	sch := NewScheduleService(r.Schedules)
	sys := NewSystemUsersService(r.SystemUsers)
	group := NewGroupService(r.Groups)
	notifSvc := NewNotificationService(r.Notifications, r.StaffUsers, sender)
	staff := NewStaffUserService(r.StaffUsers, sch, sender, notifSvc)
	authorize := NewAuthorizeService(t, r.Auth, sys, staff)
	visitPhotos := NewVisitPhotosService(r.VisitPhotos)
	store := NewStoresService(r.Store)
	visit := NewVisitsService(r.Visit, minio, store, visitPhotos, staff)
	location := NewLocationHistoryService(r.Location)
	ss := NewSystemSettingsService(r.SystemSettings)
	ssm := NewStaffStoresMappingService(r.StaffStoresMapping)
	excel := NewExcelService()
	audit := NewAuditLogsService(r.Audit)
	storePhotos := NewStorePhotosService(r.StorePhotos)
	return &Service{
		SystemUsers:        sys,
		Group:              group,
		StaffUser:          staff,
		Authorize:          authorize,
		Schedule:           sch,
		Visit:              visit,
		Store:              store,
		Location:           location,
		SystemSettings:     ss,
		StaffStoresMapping: ssm,
		Excel:              excel,
		Audit:              audit,
		Notification:       notifSvc,
		VisitPhotos:        visitPhotos,
		StorePhotos:        storePhotos,
	}
}

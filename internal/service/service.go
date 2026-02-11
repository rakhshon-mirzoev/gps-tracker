package service

import (
	"github.com/rakhshon-mirzoev/gps-tracker/internal/repository"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/websocket"
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

func New(r *repository.Repository, t token.Token, minio minioservice.Minio, sender fcm.PushSender, wsHub *websocket.Hub) *Service {
	audit := NewAuditLogsService(r.Audit)
	sch := NewScheduleService(r.Schedules, audit)
	sys := NewSystemUsersService(r.SystemUsers, audit)
	group := NewGroupService(r.Groups, audit)
	notifSvc := NewNotificationService(r.Notifications, r.StaffUsers, sender)
	staff := NewStaffUserService(r.StaffUsers, sch, sender, notifSvc, audit, minio)
	authorize := NewAuthorizeService(t, r.Auth, sys, staff)
	visitPhotos := NewVisitPhotosService(r.VisitPhotos)
	store := NewStoresService(r.Store, audit)
	visit := NewVisitsService(r.Visit, minio, store, visitPhotos, staff, notifSvc, audit, sch, group)
	location := NewLocationHistoryService(r.Location, staff, websocket.NewLocationPublisher(wsHub))
	ss := NewSystemSettingsService(r.SystemSettings, audit)
	ssm := NewStaffStoresMappingService(r.StaffStoresMapping)
	excel := NewExcelService()
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

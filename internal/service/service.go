package service

import (
	"github.com/rakhshon-mirzoev/gps-tracker/internal/events"
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
	SystemNotification SystemNotification
	VisitPhotos        VisitPhotos
	StorePhotos        StorePhotos
	StaffSchedule      StaffSchedule
}

func New(r *repository.Repository, t token.Token, minio minioservice.Minio, sender fcm.PushSender, publisher events.LocationPublisher) *Service {
	audit := NewAuditLogsService(r.Audit)
	storePhotos := NewStorePhotosService(r.StorePhotos)
	sch := NewScheduleService(r.Schedules, audit)
	sys := NewSystemUsersService(r.SystemUsers, audit)
	group := NewGroupService(r.Groups, audit)
	ssch := NewStaffScheduleService(r.StaffSchedule)
	notifSvc := NewNotificationService(r.Notifications, r.StaffUsers, sender)
	sysNotifSvc := NewSystemNotificationService(r.SystemNotifications, r.SystemUsers)
	staff := NewStaffUserService(r.StaffUsers, r.Visit, sch, sender, notifSvc, audit, minio, ssch)
	authorize := NewAuthorizeService(t, r.Auth, sys, staff)
	visitPhotos := NewVisitPhotosService(r.VisitPhotos, minio)
	store := NewStoresService(r.Store, audit, minio, storePhotos)
	visit := NewVisitsService(r.Visit, minio, store, visitPhotos, staff, notifSvc, audit, sch, group)
	location := NewLocationHistoryService(r.Location, r.GpsAlertStates, staff, publisher)
	ss := NewSystemSettingsService(r.SystemSettings, audit)
	ssm := NewStaffStoresMappingService(r.StaffStoresMapping, staff)
	excel := NewExcelService()

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
		SystemNotification: sysNotifSvc,
		VisitPhotos:        visitPhotos,
		StorePhotos:        storePhotos,
		StaffSchedule:      ssch,
	}
}

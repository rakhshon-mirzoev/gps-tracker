package main

import (
	"github.com/IkbolNazarov/gps-tracker/internal/handler"
	"github.com/IkbolNazarov/gps-tracker/internal/middleware"
	"github.com/IkbolNazarov/gps-tracker/internal/service"
	"github.com/IkbolNazarov/gps-tracker/internal/websocket"
	"github.com/IkbolNazarov/gps-tracker/pkg/constants"
	miniohandler "github.com/IkbolNazarov/gps-tracker/pkg/storage/handler"
	"github.com/IkbolNazarov/gps-tracker/pkg/token"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"golang.org/x/time/rate"

	_ "github.com/IkbolNazarov/gps-tracker/docs/admin"
)

func registerAdminRoutes(r *gin.Engine, mw middleware.Middleware, h *handler.Handler, m *miniohandler.MinioHandler, wsHub *websocket.Hub, nfHub *websocket.NotificationHub, svc *service.Service, tokenSvc token.Token) {
	r.ForwardedByClientIP = true
	r.SetTrustedProxies([]string{"127.0.0.1"})

	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization", "Connection", "Upgrade"},
		AllowWebSockets:  true,
		AllowCredentials: false,
	}))

	// limiter := middleware.NewIPRateLimiter(5, 10)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.InstanceName("admin")))

	api := r.Group("/api/v1")

	api.Use(
		mw.Recovery(),
		mw.RequestID(),
		mw.Logger(),
		// mw.RateLimitMiddleware(limiter),
	)
	{
		api.GET("/ping", h.Ping)
	}
	storage := api.Group("/storage")
	{
		storage.GET("download-url", m.GetDownloadURL)
		storage.GET("upload-url", m.GetUploadURL)

	}

	loginLimiter := middleware.NewIPRateLimiter(rate.Limit(1.0/60.0), 5)

	auth := api.Group(constants.AuthGroup + "/" + constants.UserTypeSystem)
	auth.Use(mw.RateLimitMiddleware(loginLimiter))
	{
		auth.POST("/login", h.Login)
		auth.POST("/register", h.RegisterSystem)
		auth.POST("/refresh", h.Refresh)
		auth.POST("/logout", h.Logout)
	}

	admin := api.Group(constants.AdminGroup)
	admin.Use(mw.Auth(), mw.Admin())

	{
		admin.GET("/reports/export", h.ExportReports)
		admin.GET("/audit/logs", h.AuditLogs)

		staff := admin.Group(constants.UserTypeStaff)
		{
			staff.GET("", h.ListStaffUsers)
			staff.POST("", h.CreateStaffUser)
			staff.PUT("/:id", h.UpdateStaffUser)
			staff.PATCH("/:id/deactivate", h.DeactivateStaffUser)
			staff.GET("/schedules", h.ListSchedulesAdmin)
			staff.POST("/schedules/assign", h.AssignSchedule)
			staff.PUT("/schedules/:id", h.UpdateSchedule)
			staff.DELETE("/schedules/:id", h.DeleteSchedule)
			staff.GET("/:id", h.GetStaffUserByID)
			staff.GET("/online-offline", h.AllStaffOnlineOffline)
			staff.GET("/workday/:id", h.GetWorkTimeStaff)
			staff.GET("/day-start/:id", h.DayStartedStaff)

			schedule := staff.Group(constants.TimeSchedule)
			{
				schedule.DELETE("/:id", h.DeleteStaffSchedule)
				schedule.PUT("/:id", h.UpdateStaffSchedule)
				schedule.GET("/:id", h.ListStaffSchedule)
			}
		}

		monitoring := staff.Group(constants.Monitoring)
		{
			monitoring.GET("/real-time", h.ListStaffUsersLocationByGroupID)
			monitoring.GET("/track/:staff_id", h.ListStaffUserLocationByDate)
		}

		store := admin.Group(constants.Stores)
		{
			store.GET("", h.ListStores)
			store.GET("/:id", h.GetStoreByID)
			store.PATCH("/:id/approve", h.ApproveStore)
			store.DELETE("/:id", h.DeleteStore)
		}

		group := admin.Group(constants.Group)
		{
			group.POST("", h.CreateGroup)
			group.DELETE("/:id", h.DeleteGroup)
			group.GET("/:id", h.GetByIdGroup)
			group.GET("", h.ListGroups)
			group.PUT("/:id", h.UpdateGroup)
		}

		staff_mapping := admin.Group(constants.StaffMapping)
		{
			staff_mapping.GET("/:staff_id", h.ListStoresByStaffIdAdmin)
			staff_mapping.POST("/:staff_id", h.MapStoreStaff)
			staff_mapping.DELETE("/:staff_id", h.UnMapStoreStaff)
		}

		visit := admin.Group(constants.VisitGroup)
		{
			visit.GET("/review-queue", h.ReviewQueue)
			visit.PATCH("/:id/review", h.ReviewVisit)
		}

		settings := admin.Group(constants.Settings)
		{
			settings.POST("", h.CreateSystemSetting)
			settings.PUT("/:id", h.UpdateSystemSetting)
			settings.GET("", h.ListSystemSettings)
			settings.DELETE("/:id", h.DeleteSystemSetting)
			settings.GET("/:id", h.GetSystemSettingByID)
		}

		administrate := admin.Group(constants.UserTypeSystem)
		administrate.Use(mw.System())
		{
			administrate.GET("", h.ListSystemUsers)
			administrate.DELETE("/:id", h.DeleteSystemUser)
			administrate.PUT("/:id", h.UpdateSystemUser)
			administrate.POST("", h.CreateSystemUser)
		}

		notifications := admin.Group(constants.Notifications)
		{
			notifications.POST("", h.CreateNotification)
			notifications.GET("", h.GetSystemNotifications)
			notifications.PATCH("/:id/read", h.MarkSystemNotificationRead)
		}

	}

	r.GET("/api/v1/admin/ws-location", websocket.Handler(websocket.Deps{
		Hub:              wsHub,
		TokenService:     tokenSvc,
		LocationProvider: svc.Location,
		GroupResolver:    svc.Group,
	}))

	r.GET("/api/v1/admin/ws-notification", websocket.NotificationHandler(nfHub, tokenSvc))
}

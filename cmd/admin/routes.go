package main

import (
	"github.com/gin-gonic/gin"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/handler"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/middleware"
	"github.com/rakhshon-mirzoev/gps-tracker/pkg/constants"
	miniohandler "github.com/rakhshon-mirzoev/gps-tracker/pkg/storage/handler"
)

func registerAdminRoutes(r *gin.Engine, mw middleware.Middleware, h *handler.Handler, m *miniohandler.MinioHandler) {
	api := r.Group("/api/v1")

	api.Use(mw.CorsMiddleWare()) /* ,
		mw.Recovery(),
		mw.RequestID(),
		mw.Logger(),
	) */
	{
		api.GET("/ping", h.Ping)
	}
	storage := api.Group("/storage")
	{
		storage.GET("download-url", m.GetDownloadURL)
		storage.GET("upload-url", m.GetUploadURL)

	}

	auth := api.Group(constants.AuthGroup + "/" + constants.UserTypeSystem)
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

		staff := admin.Group(constants.UserTypeStaff)
		{
			staff.GET("", h.ListStaffUsers)
			staff.POST("", h.CreateStaffUser)
			staff.PUT("/:id", h.UpdateStaffUser)
			staff.PATCH("/:id/deactivate", h.DeactivateStaffUser)
			staff.POST("/schedules/assign", h.AssignSchedule)
			staff.GET("/online-offline", h.AllStaffOnlineOffline)
		}

		monitoring := staff.Group(constants.Monitoring)
		{
			monitoring.GET("/real-time", h.ListStaffUsersLocationByGroupID)
			monitoring.GET("/track/:staff_id", h.ListStaffUserLocationByDate)
		}

		store := admin.Group(constants.Stores)
		{
			store.GET(constants.Stores, h.ListStores)
			store.PATCH("/:id/approve", h.ApproveStore)
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

	}

}

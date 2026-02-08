package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/handler"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/middleware"
	"github.com/rakhshon-mirzoev/gps-tracker/pkg/constants"
	miniohandler "github.com/rakhshon-mirzoev/gps-tracker/pkg/storage/handler"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/rakhshon-mirzoev/gps-tracker/docs/mobile"
)

func registerMobileRoutes(r *gin.Engine, mw middleware.Middleware, h *handler.Handler, m *miniohandler.MinioHandler) {
	r.Use(cors.Default())

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.InstanceName("mobile")))

	api := r.Group("/api/v1")

	api.Use(
		mw.Recovery(),
		mw.RequestID(),
		mw.Logger(),
	)
	{
		api.GET("/ping", h.Ping)
		api.GET("/settings", h.ListSystemSettings)
	}
	storage := api.Group("/storage")
	{
		storage.GET("download-url", m.GetDownloadURL)
		storage.GET("upload-url", m.GetUploadURL)

	}
	auth := api.Group(constants.AuthGroup + "/" + constants.UserTypeStaff)
	{
		auth.POST("/login", h.Login)
		auth.POST("/register", h.RegisterStaff)
		auth.POST("/refresh", h.Refresh)
		auth.POST("/logout", h.Logout)
	}

	mobi := api.Group(constants.MobileGroup)
	mobi.Use(mw.Auth(), mw.Mobile())
	{
		mobi.GET("/profile", h.Profile)
		mobi.PUT("/device-token", h.AttachDevice)
		mobi.GET("/notifications", h.ListNotifications)
		mobi.PATCH("/notifications/:id/read", h.MarkNotificationRead)
		mobi.GET("/schedule", h.ListShedules)
		mobi.GET("/my-stores", h.ListStoresByStaffIdStaff)
		mobi.POST("/location", h.Location)
		mobi.POST("/stores/propose", h.CreateStore)

	}
	visit := mobi.Group(constants.VisitGroup)
	{
		visit.POST("/check-in", h.VisitCheckIn)
		visit.POST("/:id/photo", h.VisitPhoto)
		visit.POST("/:id/complete", h.VisitComplete)
	}
}

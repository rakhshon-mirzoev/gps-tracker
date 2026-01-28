package main

import (
	"github.com/gin-gonic/gin"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/handler"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/middleware"
	"github.com/rakhshon-mirzoev/gps-tracker/pkg/constants"
	miniohandler "github.com/rakhshon-mirzoev/gps-tracker/pkg/storage/handler"
)

func registerMobileRoutes(r *gin.Engine, mw middleware.Middleware, h *handler.Handler, m *miniohandler.MinioHandler) {
	api := r.Group("/api/v1")

	api.Use(mw.CorsMiddleWare()) /*,
		a.middleware.Recovery(),
		a.middleware.RequestID(),
		a.middleware.Logger(),
	) */
	{
		api.GET("/ping", h.Ping)
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
	mobi.Use(mw.Auth(), mw.Role())
	{
		mobi.GET("/profile", h.Profile)
		mobi.PUT("/device-token", h.AttachDevice)
		mobi.GET("/schedule", h.ListShedules)
		mobi.GET("/my-stores", h.ListStores)
		mobi.POST("/location", h.Location)
		visit := mobi.Group(constants.VisitGroup)
		{
			visit.POST("/check-in", h.VisitCheckIn)
			visit.POST("/:id/photo", h.VisitPhoto)
			visit.POST("/:id/complete", h.VisitComplete)
		}
	}
}

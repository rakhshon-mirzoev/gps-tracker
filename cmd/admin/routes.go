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

	auth := api.Group(constants.AuthGroup + "/" + constants.UserTypeSystem)
	{
		auth.POST("/login", h.Login)
		auth.POST("/register", h.RegisterSystem)
		auth.POST("/refresh", h.Refresh)
		auth.POST("/logout", h.Logout)
	}
}

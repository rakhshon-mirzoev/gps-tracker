package main

import (
	"github.com/gin-gonic/gin"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/handler"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/middleware"
)

func registerAdminRoutes(r *gin.Engine, mw middleware.Middleware, h *handler.Handler, m *handler.MinioHandler) {
	api := r.Group("/api/v1")

	api.Use(mw.CorsMiddleWare()) /*,
		a.middleware.Recovery(),
		a.middleware.RequestID(),
		a.middleware.Logger(),
	) */
	{
		api.GET("", m.GetDownloadURL)
		api.GET("/ping", h.Ping)
		api.POST("/refresh", h.Refresh)
		api.POST("/logout", h.Logout)

	}
}

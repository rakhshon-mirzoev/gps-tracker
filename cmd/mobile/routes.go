package main

import (
	"github.com/IkbolNazarov/gps-tracker/internal/handler"
	"github.com/IkbolNazarov/gps-tracker/internal/middleware"
	"github.com/IkbolNazarov/gps-tracker/pkg/constants"
	miniohandler "github.com/IkbolNazarov/gps-tracker/pkg/storage/handler"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"golang.org/x/time/rate"

	_ "github.com/IkbolNazarov/gps-tracker/docs/mobile"
)

func registerMobileRoutes(r *gin.Engine, mw middleware.Middleware, h *handler.Handler, m *miniohandler.MinioHandler) {
	r.ForwardedByClientIP = true
	r.SetTrustedProxies([]string{"127.0.0.1"})

	r.Use(cors.Default())

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.InstanceName("mobile")))

	limiter := middleware.NewIPRateLimiter(5, 10)

	api := r.Group("/api/v1")

	api.Use(
		mw.Recovery(),
		mw.RequestID(),
		mw.Logger(),
		mw.RateLimitMiddleware(limiter),
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

	loginLimiter := middleware.NewIPRateLimiter(rate.Limit(1.0/60.0), 5)

	auth := api.Group(constants.AuthGroup + "/" + constants.UserTypeStaff)
	auth.Use(mw.RateLimitMiddleware(loginLimiter))
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
		mobi.PUT("/profile", h.UpdateProfile)
		mobi.PUT("/profile/avatar", h.UpdateAvatar)
		mobi.PUT("/device-token", h.AttachDevice)
		mobi.GET("/notifications", h.ListNotifications)
		mobi.PATCH("/notifications/:id/read", h.MarkNotificationRead)
		mobi.GET("/schedule", h.ListShedules)
		mobi.GET("/my-stores", h.ListStoresByStaffIdStaff)
		mobi.POST("/location", h.Location)
		mobi.POST("/stores/upload-photo", h.UploadStorePhoto)
		mobi.POST("/stores/propose", h.CreateStore)
		mobi.GET("/day-status", h.DayStatus)
		mobi.PUT("/start-work", h.StartWork)
		mobi.PUT("/end-work", h.EndWork)
	}

	visit := mobi.Group(constants.VisitGroup)
	{
		visit.POST("/check-in", h.VisitCheckIn)
		visit.POST("/:id/photo", h.VisitPhoto)
		visit.POST("/:id/complete", h.VisitComplete)
		visit.GET("", h.ListVisitsByStaffID)
		visit.GET("/:id", h.GetVisitByID)
		visit.DELETE("/:id/photo", h.DeleteVisitPhoto)
		visit.PUT("/:id", h.UpdateVisit)
	}
}

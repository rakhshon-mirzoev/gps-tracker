package app

import "github.com/rakhshon-mirzoev/gps-tracker/pkg/constants"

func (a *App) registerRoutes() {
	api := a.router.Group("/api/v1")

	api.Use(a.middleware.CorsMiddleWare()) /*,
		a.middleware.Recovery(),
		a.middleware.RequestID(),
		a.middleware.Logger(),
	) */
	{
		api.GET("/ping", a.h.Ping)
		api.POST("/refresh", a.h.Refresh)
		api.POST("/logout", a.h.Logout)
		api.GET("/storage/upload-url", a.h.GetUploadURL)
		api.GET("/storage/download-url", a.h.GetDownloadURL)

	}

	auth := api.Group(constants.AuthGroup)
	{
		auth.POST("/login", a.h.Login)
		auth.POST("/register/:type", a.h.Register)
	}

	// system.Use(a.middleware.Auth())
}

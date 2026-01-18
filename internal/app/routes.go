package app

import "github.com/rakhshon-mirzoev/gps-tracker/pkg/constants"

func (a *App) registerRoutes() {
	api := a.Router.Group("/api-v1")

	api.Use(a.middleware.CorsMiddleWare()) /*,
		a.middleware.Recovery(),
		a.middleware.RequestID(),
		a.middleware.Logger(),
	) */
	{
		api.GET("/ping", a.H.Ping)
		api.POST("/refresh", a.H.Refresh)
		api.POST("/logout", a.H.Logout)
	}

	auth := api.Group(constants.AuthGroup)
	{
		auth.POST("/login", a.H.Login)
		auth.POST("/register/:type", a.H.Register)
	}

	// system.Use(a.middleware.Auth())
}

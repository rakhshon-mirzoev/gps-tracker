package app

import "github.com/rakhshon-mirzoev/gps-tracker/pkg/constants"

func (a *App) registerRoutes() {
	api := a.Router.Group("/api-v1")

	api.Use(
		a.middleware.CorsMiddleWare(),
		a.middleware.Recovery(),
		a.middleware.RequestID(),
		a.middleware.Logger(),
	)

	api.GET("/ping", a.H.Ping)

	system := api.Group("/system")

	system.Use(a.middleware.Auth())

	system.POST(constants.AuthGroup+"/login", a.H.Login)
	system.POST(constants.AuthGroup+"/registrate", a.H.Registrate)
}

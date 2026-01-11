package app

func (a *App) registerRoutes() {
	api := a.Router.Group("/api-v1")

	api.GET("/ping", a.H.Ping)
}

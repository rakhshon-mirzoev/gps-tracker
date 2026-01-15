package app

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rakhshon-mirzoev/gps-tracker/config"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/handler"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/middleware"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/repository"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/service"
	"github.com/rakhshon-mirzoev/gps-tracker/pkg/db"
	"github.com/rakhshon-mirzoev/gps-tracker/pkg/token"
)

type App struct {
	DB         *db.DB
	H          *handler.Handler
	S          *service.Service
	R          *repository.Repository
	Router     *gin.Engine
	Server     *http.Server
	middleware middleware.Middleware
	token      token.Token
}

func New(cfg *config.Config) (*App, error) {
	database, err := db.ConnectToDb(cfg)
	if err != nil {
		return nil, err
	}

	app := &App{}
	app.DB = database
	app.token = token.InitTokenService()
	repository := repository.InitRepository(app.DB)
	service := service.InitService(repository, app.token)
	handler := handler.InitHandler(service)
	middleware := middleware.InitMiddleware(app.token)

	router := gin.Default()

	return &App{
		DB:         database,
		R:          repository,
		S:          service,
		H:          handler,
		Router:     router,
		middleware: middleware,
	}, nil
}

func (a *App) Run(addr string) error {
	a.registerRoutes()

	a.Server = &http.Server{
		Addr:    addr,
		Handler: a.Router,
	}

	return a.Server.ListenAndServe()
}

func (a *App) Shutdown(ctx context.Context) error {
	return a.Server.Shutdown(ctx)
}

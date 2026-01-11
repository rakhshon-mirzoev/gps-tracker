package app

import (
	"github.com/gin-gonic/gin"
	"github.com/rakhshon-mirzoev/gps-tracker/config"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/handler"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/repository"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/service"
	"github.com/rakhshon-mirzoev/gps-tracker/pkg/db"
)

type App struct {
	DB     *db.DB
	H      *handler.Handler
	S      *service.Service
	R      *repository.Repository
	Router *gin.Engine
}

func New(cfg *config.Config) (*App, error) {
	database, err := db.ConnectToDb(cfg)
	if err != nil {
		return nil, err
	}

	repo := repository.InitRepository(database)
	service := service.InitService(repo)
	handler := handler.InitHandler(service)

	router := gin.Default()

	return &App{
		DB:     database,
		R:      repo,
		S:      service,
		H:      handler,
		Router: router,
	}, nil
}

func (a *App) Run(addr string) error {
	a.registerRoutes()

	return a.Router.Run(addr)
}

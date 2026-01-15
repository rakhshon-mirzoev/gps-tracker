package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	router := gin.New()

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

	server := &http.Server{
		Addr:    addr,
		Handler: a.Router,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("server started on %s", addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen error: %v", err)
		}
	}()

	<-quit
	log.Println("shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown failed: %w", err)
	}

	a.DB.Close()

	log.Println("server exited gracefully")
	return nil
}

func (a *App) Shutdown(ctx context.Context) error {
	return a.Server.Shutdown(ctx)
}

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
	"github.com/minio/minio-go/v7"
	"github.com/rakhshon-mirzoev/gps-tracker/config"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/handler"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/middleware"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/repository"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/service"
	"github.com/rakhshon-mirzoev/gps-tracker/pkg/db"
	mino "github.com/rakhshon-mirzoev/gps-tracker/pkg/minio"
	"github.com/rakhshon-mirzoev/gps-tracker/pkg/token"
)

type App struct {
	db          *db.DB
	h           *handler.Handler
	s           *service.Service
	r           *repository.Repository
	router      *gin.Engine
	server      *http.Server
	middleware  middleware.Middleware
	token       token.Token
	minioClient *minio.Client
}

func New(cfg *config.GeneralConfig) (*App, error) {
	database, err := db.ConnectToDb(&cfg.Core)
	if err != nil {
		return nil, err
	}

	minioClient, err := mino.NewMinio(cfg.Minio)
	if err != nil {
		return nil, err
	}

	token := token.InitTokenService(os.Getenv("SECRET_KEY"), time.Duration(cfg.Core.TTLMinutes)*time.Minute)
	repository := repository.InitRepository(database)
	service := service.InitService(repository, token, minioClient, cfg.Minio.Bucket)
	handler := handler.InitHandler(service)
	middleware := middleware.InitMiddleware(token)

	router := gin.Default()

	return &App{
		db:          database,
		r:           repository,
		s:           service,
		h:           handler,
		router:      router,
		middleware:  middleware,
		token:       token,
		minioClient: minioClient,
	}, nil
}

func (a *App) Run(addr string) error {
	a.registerRoutes()

	server := &http.Server{
		Addr:    addr,
		Handler: a.router,
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

	a.db.Close()

	log.Println("server exited gracefully")
	return nil
}

func (a *App) Shutdown(ctx context.Context) error {
	return a.server.Shutdown(ctx)
}

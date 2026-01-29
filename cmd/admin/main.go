package main

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
	"github.com/rakhshon-mirzoev/gps-tracker/internal/app"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/handler"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/logger"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/middleware"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/repository"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/service"
	"github.com/rakhshon-mirzoev/gps-tracker/pkg/db"
	miniohandler "github.com/rakhshon-mirzoev/gps-tracker/pkg/storage/handler"
	minioservice "github.com/rakhshon-mirzoev/gps-tracker/pkg/storage/service"
	"github.com/rakhshon-mirzoev/gps-tracker/pkg/token"
)

func main() {
	if err := logger.Init(); err != nil {
		log.Fatalf("failed to init logger: %v", err)
	}
	defer logger.Log.Sync()

	if err := run(); err != nil {
		logger.Log.Fatalf("failed to run: %v", err)
	}
}

func run() error {
	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		return err
	}

	db, err := db.ConnectToDb(&cfg.DB)
	if err != nil {
		return err
	}

	client, err := minioservice.NewMinio(cfg.Minio)
	if err != nil {
		return err
	}

	token := token.InitTokenService(os.Getenv("SECRET_KEY"), time.Duration(cfg.Admin.TTLMinutes))

	repository := repository.New(db)
	service := service.New(repository, token, client)
	minio := miniohandler.New(client)
	handler := handler.New(service)

	mw := middleware.New(token)

	router := gin.Default()

	registerAdminRoutes(router, mw, handler, minio)

	app := app.New(router)

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	go func() {
		addr := fmt.Sprintf(":%s", cfg.Mobile.Port)
		if err := app.Run(addr); err != nil && err != http.ErrServerClosed {
			logger.Log.Fatalf("server error: %v", err)
		}
	}()

	<-ctx.Done()
	logger.Log.Info("shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := app.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("shutdown failed: %w", err)
	}

	logger.Log.Info("server exited gracefully")
	return nil
}

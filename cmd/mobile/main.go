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
	fcm "github.com/rakhshon-mirzoev/gps-tracker/pkg/fcm/app"
	sender "github.com/rakhshon-mirzoev/gps-tracker/pkg/fcm/service"
	miniohandler "github.com/rakhshon-mirzoev/gps-tracker/pkg/storage/handler"
	minioservice "github.com/rakhshon-mirzoev/gps-tracker/pkg/storage/service"
	"github.com/rakhshon-mirzoev/gps-tracker/pkg/token"
	"go.uber.org/zap"
)

func main() {
	if err := logger.Init(); err != nil {
		log.Fatalf("failed to init logger: %v", err)
	}

	defer func() {
		_ = logger.Log.Sync()
	}()

	if err := run(); err != nil {
		logger.Log.Error("application stopped with error", zap.Error(err))
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		return err
	}

	dbConn, err := db.ConnectToDb(&cfg.DB)
	if err != nil {
		return err
	}
	defer dbConn.Close()

	client, err := minioservice.NewMinio(cfg.Minio)
	if err != nil {
		return err
	}

	tokenService := token.NewTokenService(os.Getenv("SECRET_KEY"), time.Duration(cfg.Mobile.TTLMinutes)*time.Minute)

	firebaseApp, err := fcm.NewFirebaseApp(context.Background(), cfg.Mobile.FirebaseCreds)
	if err != nil {
		return err
	}

	firebaseService, err := sender.NewFirebaseService(context.Background(), firebaseApp)
	if err != nil {
		return err
	}

	repo := repository.New(dbConn)
	services := service.New(repo, tokenService, client, firebaseService)
	handlers := handler.New(services)
	minioHandler := miniohandler.New(client)

	mw := middleware.New(tokenService)

	router := gin.Default()
	registerMobileRoutes(router, mw, handlers, minioHandler)

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
			logger.Log.Fatal("failed to run mobi", zap.Error(err))
		}
	}()

	<-ctx.Done()
	log.Println("shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := app.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("shutdown failed: %w", err)
	}

	log.Println("server exited gracefully")
	return nil
}

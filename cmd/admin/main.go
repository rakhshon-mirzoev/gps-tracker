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
	"github.com/rakhshon-mirzoev/gps-tracker/internal/websocket"
	"github.com/rakhshon-mirzoev/gps-tracker/pkg/db"
	"github.com/rakhshon-mirzoev/gps-tracker/pkg/fcm"
	"github.com/rakhshon-mirzoev/gps-tracker/pkg/kafka"
	miniohandler "github.com/rakhshon-mirzoev/gps-tracker/pkg/storage/handler"
	minioservice "github.com/rakhshon-mirzoev/gps-tracker/pkg/storage/service"
	"github.com/rakhshon-mirzoev/gps-tracker/pkg/token"
	"go.uber.org/zap"
)

func main() {
	if err := logger.New(); err != nil {
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
	cfg, creds, err := config.LoadConfig("config.json")
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

	firebaseApp, err := fcm.NewFirebaseApp(context.Background(), creds)
	if err != nil {
		return err
	}

	firebaseService, err := fcm.NewFirebaseService(context.Background(), firebaseApp)
	if err != nil {
		return err
	}

	token := token.NewTokenService(os.Getenv("SECRET_KEY"), time.Duration(cfg.Admin.TTLMinutes)*time.Minute)

	router := gin.Default()

	app := app.New(router)

	// WebSocket publisher — для локальной публикации в Hub (если admin тоже генерирует локации)
	wsPublisher := websocket.NewLocationPublisher(app.WsHub())

	repository := repository.New(db)
	service := service.New(repository, token, client, firebaseService, wsPublisher)
	minio := miniohandler.New(client)
	handler := handler.New(service, client)

	mw := middleware.New(token)

	registerAdminRoutes(router, mw, handler, minio, app.WsHub(), service, token)

	// Kafka consumer — читает локации из топика и бросает в Hub
	kafkaConsumer := kafka.NewLocationConsumer(cfg.Kafka.Brokers, cfg.Kafka.Topic, app.WsHub())
	defer kafkaConsumer.Close()

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	go kafkaConsumer.Run(ctx)

	go func() {
		addr := fmt.Sprintf(":%s", cfg.Admin.Port)
		if err := app.Run(addr); err != nil && err != http.ErrServerClosed {
			logger.Log.Fatal("failed to run admin", zap.Error(err))
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

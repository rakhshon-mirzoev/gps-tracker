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

	"github.com/IkbolNazarov/gps-tracker/config"
	"github.com/IkbolNazarov/gps-tracker/internal/app"
	"github.com/IkbolNazarov/gps-tracker/internal/handler"
	"github.com/IkbolNazarov/gps-tracker/internal/logger"
	"github.com/IkbolNazarov/gps-tracker/internal/middleware"
	"github.com/IkbolNazarov/gps-tracker/internal/repository"
	"github.com/IkbolNazarov/gps-tracker/internal/service"
	"github.com/IkbolNazarov/gps-tracker/internal/websocket"
	gps_alert "github.com/IkbolNazarov/gps-tracker/internal/worker/gps_alert"
	"github.com/IkbolNazarov/gps-tracker/pkg/db"
	"github.com/IkbolNazarov/gps-tracker/pkg/fcm"
	"github.com/IkbolNazarov/gps-tracker/pkg/kafka"
	miniohandler "github.com/IkbolNazarov/gps-tracker/pkg/storage/handler"
	minioservice "github.com/IkbolNazarov/gps-tracker/pkg/storage/service"
	"github.com/IkbolNazarov/gps-tracker/pkg/timezone"
	"github.com/IkbolNazarov/gps-tracker/pkg/token"
	"github.com/gin-gonic/gin"
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

	if err := timezone.Init(); err != nil {
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

	wsPublisher := websocket.NewLocationPublisher(app.WsHub())

	// NotificationHub — адресный хаб для уведомлений admin/supervisor
	notifHub := websocket.NewNotificationHub()

	repository := repository.New(db)
	service := service.New(repository, token, client, firebaseService, wsPublisher)
	minio := miniohandler.New(client)
	handler := handler.New(service, client)

	mw := middleware.New(token)

	registerAdminRoutes(router, mw, handler, minio, app.WsHub(), notifHub, service, token)

	// Создаём топик Kafka для локаций
	if err := kafka.EnsureTopic(cfg.Kafka.Brokers, cfg.Kafka.Topic, 1, 1); err != nil {
		logger.Log.Warn("failed to ensure kafka location topic (may already exist)", zap.Error(err))
	}

	// Создаём топик Kafka для уведомлений
	notifTopic := cfg.Kafka.NotificationTopic
	if notifTopic == "" {
		notifTopic = "admin-notifications"
	}
	if err := kafka.EnsureTopic(cfg.Kafka.Brokers, notifTopic, 1, 1); err != nil {
		logger.Log.Warn("failed to ensure kafka notification topic (may already exist)", zap.Error(err))
	}

	// Kafka consumer для локаций
	kafkaConsumer := kafka.NewLocationConsumer(cfg.Kafka.Brokers, cfg.Kafka.Topic, app.WsHub())
	defer kafkaConsumer.Close()

	// Kafka consumer для уведомлений
	notifPublisher := kafka.NewNotificationPublisher(cfg.Kafka.Brokers, notifTopic)
	defer notifPublisher.Close()

	notifConsumer := kafka.NewNotificationConsumer(cfg.Kafka.Brokers, notifTopic, service.SystemNotification, notifHub)
	defer notifConsumer.Close()

	// GPS-alert worker: каждые 5 минут ищет сотрудников без GPS и кидает в Kafka
	gpsWorker := gps_alert.New(repository.Location, repository.GpsAlertStates, notifPublisher, 5*time.Minute)

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	// app.WsHub().Run(ctx) уже запущен внутри app.New()
	go notifHub.Run(ctx)
	go kafkaConsumer.Run(ctx)
	go notifConsumer.Run(ctx)
	go gpsWorker.Run()

	go func() {
		addr := fmt.Sprintf(":%s", cfg.Admin.Port)
		if err := app.Run(addr); err != nil && err != http.ErrServerClosed {
			logger.Log.Fatal("failed to run admin", zap.Error(err))
		}
	}()

	<-ctx.Done()
	logger.Log.Info("shutdown signal received")

	gpsWorker.Stop()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := app.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("shutdown failed: %w", err)
	}

	logger.Log.Info("server exited gracefully")
	return nil
}

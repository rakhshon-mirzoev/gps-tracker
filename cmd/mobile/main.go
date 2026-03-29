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
	"github.com/IkbolNazarov/gps-tracker/internal/handler"
	"github.com/IkbolNazarov/gps-tracker/internal/logger"
	"github.com/IkbolNazarov/gps-tracker/internal/middleware"
	"github.com/IkbolNazarov/gps-tracker/internal/repository"
	"github.com/IkbolNazarov/gps-tracker/internal/service"
	"github.com/IkbolNazarov/gps-tracker/internal/worker/notifications"
	"github.com/IkbolNazarov/gps-tracker/pkg/cron"
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

	firebaseApp, err := fcm.NewFirebaseApp(context.Background(), creds)
	if err != nil {
		return err
	}

	firebaseService, err := fcm.NewFirebaseService(context.Background(), firebaseApp)
	if err != nil {
		return err
	}

	cleanupCron := cron.StartCleanupCron(client)

	// Создаём топик Kafka, если он ещё не существует
	if err := kafka.EnsureTopic(cfg.Kafka.Brokers, cfg.Kafka.Topic, 1, 1); err != nil {
		logger.Log.Warn("failed to ensure kafka topic (may already exist)", zap.Error(err))
	}

	// Kafka publisher — отправляет локации в топик, admin consumer прочитает и пробросит в WebSocket
	publisher := kafka.NewLocationPublisher(cfg.Kafka.Brokers, cfg.Kafka.Topic)
	defer publisher.Close()

	router := gin.Default()

	repo := repository.New(dbConn)
	services := service.New(repo, tokenService, client, firebaseService, publisher)
	handlers := handler.New(services, client)
	minioHandler := miniohandler.New(client)

	mw := middleware.New(tokenService)

	registerMobileRoutes(router, mw, handlers, minioHandler)

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	worker := notifications.New(repo.Notifications, repo.StaffUsers, firebaseService, time.Duration(cfg.Mobile.Interval)*time.Second)

	go worker.Run()

	srv := &http.Server{
		Addr:        fmt.Sprintf(":%s", cfg.Mobile.Port),
		Handler:     router,
		ReadTimeout: 10 * time.Second,
		IdleTimeout: 120 * time.Second,
	}

	go func() {
		log.Printf("mobile server started on :%s", cfg.Mobile.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Log.Fatal("failed to run mobile", zap.Error(err))
		}
	}()

	<-ctx.Done()
	log.Println("shutdown signal received")

	cleanupCtx := cleanupCron.Stop()
	<-cleanupCtx.Done()
	logger.Log.Info("cleanupCtx cron stopped")

	worker.Stop()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("shutdown failed: %w", err)
	}

	log.Println("server exited gracefully")
	return nil
}

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
	"github.com/rakhshon-mirzoev/gps-tracker/internal/middleware"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/repository"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/service"
	"github.com/rakhshon-mirzoev/gps-tracker/pkg/db"
	miniohandler "github.com/rakhshon-mirzoev/gps-tracker/pkg/storage/handler"
	minioservice "github.com/rakhshon-mirzoev/gps-tracker/pkg/storage/service"
	"github.com/rakhshon-mirzoev/gps-tracker/pkg/token"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		return err
	}

	// init infra
	dbConn, err := db.ConnectToDb(&cfg.DB)
	if err != nil {
		return err
	}
	defer dbConn.Close()

	client, err := minioservice.NewMinio(cfg.Minio)
	if err != nil {
		return err
	}

	tokenService := token.InitTokenService(
		os.Getenv("SECRET_KEY"),
		time.Duration(cfg.Mobile.TTLMinutes)*time.Minute,
	)

	repo := repository.InitRepository(dbConn)
	services := service.InitService(repo, tokenService)
	handlers := handler.InitHandler(services)
	minioHandler := miniohandler.NewMinioHandler(client)
	mw := middleware.InitMiddleware(tokenService)

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
		addr := fmt.Sprintf("%s:%s", cfg.Mobile.Host, cfg.Mobile.Port)
		if err := app.Run(addr); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
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

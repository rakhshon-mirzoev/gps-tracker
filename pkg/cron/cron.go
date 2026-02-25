package cron

import (
	"context"
	"time"

	"github.com/rakhshon-mirzoev/gps-tracker/internal/logger"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/service"
	minioservice "github.com/rakhshon-mirzoev/gps-tracker/pkg/storage/service"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

func StartCleanupCron(minioService minioservice.Minio) *cron.Cron {
	c := cron.New(
		cron.WithChain(
			cron.Recover(cron.DefaultLogger),
		),
	)

	_, err := c.AddFunc("@every 10m", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()

		err := minioService.RemoveOldTmpObjects(ctx, "stores/tmp/", 1*time.Hour)
		if err != nil {
			logger.Log.Error("cleanup tmp objects failed", zap.Error(err))
		}
	})

	if err != nil {
		logger.Log.Error("failed to add cron job", zap.Error(err))
	}

	c.Start()
	return c
}

func MakeInactive(s service.Service) *cron.Cron {
	c := cron.New(
		cron.WithChain(
			cron.Recover(cron.DefaultLogger),
		),
	)

	_, err := c.AddFunc("@every 5m", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()

		err := s.StaffUser.MakeInactiveUsers(ctx)
		if err != nil {
			logger.Log.Error("make inactive", zap.Error(err))
		}
	})

	if err != nil {
		logger.Log.Error("failed to add cron job", zap.Error(err))
	}

	c.Start()
	return c
}

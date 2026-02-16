package notifications

import (
	"context"
	"sync"
	"time"

	"github.com/rakhshon-mirzoev/gps-tracker/internal/logger"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/repository"
	"github.com/rakhshon-mirzoev/gps-tracker/pkg/fcm"
	"go.uber.org/zap"
)

type Worker struct {
	notifications repository.NotificationsRepository
	staff         repository.StaffUsersRepository
	sender        fcm.PushSender
	interval      time.Duration
	ticker        *time.Ticker
	wg            sync.WaitGroup
	stop          chan struct{}
}

func New(notificationsRepo repository.NotificationsRepository, staffRepo repository.StaffUsersRepository, sender fcm.PushSender, interval time.Duration) *Worker {
	return &Worker{
		notifications: notificationsRepo,
		staff:         staffRepo,
		sender:        sender,
		interval:      interval,
		stop:          make(chan struct{}),
	}
}

func (w *Worker) process(ctx context.Context) {
	notifications, err := w.notifications.GetUnsent(ctx, 50)
	if err != nil {
		logger.Log.Error("get unsent notifications", zap.Error(err))
		return
	}

	for _, n := range notifications {
		if n.UserID == nil {
			continue
		}

		user, err := w.staff.GetByID(ctx, *n.UserID)
		if err != nil || user.FcmToken == nil {
			continue
		}

		if user.DeactivatedByID != nil {
			continue
		}

		err = w.sender.SendToToken(
			ctx,
			*user.FcmToken,
			n.Title,
			n.Body,
			map[string]string{
				"type": n.Type,
				"id":   n.ID.String(),
			},
		)

		if err != nil {
			logger.Log.Warn("push failed", zap.Error(err))
			continue
		}

		if err := w.notifications.UpdatePushSent(ctx, n.ID, true); err != nil {
			logger.Log.Error("mark sent", zap.Error(err))
		}
	}
}

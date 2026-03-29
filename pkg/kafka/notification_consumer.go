package kafka

import (
	"context"
	"encoding/json"
	"time"

	"github.com/IkbolNazarov/gps-tracker/internal/events"
	"github.com/IkbolNazarov/gps-tracker/internal/logger"
	"github.com/IkbolNazarov/gps-tracker/internal/service"
	"github.com/IkbolNazarov/gps-tracker/internal/websocket"
	kafkago "github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

// NotificationConsumer читает NotificationEvent из Kafka,
// сохраняет уведомления в БД и пушит активным WS-соединениям.
type NotificationConsumer struct {
	reader *kafkago.Reader
	svc    service.SystemNotification
	hub    *websocket.NotificationHub
}

func NewNotificationConsumer(
	brokers []string,
	topic string,
	svc service.SystemNotification,
	hub *websocket.NotificationHub,
) *NotificationConsumer {
	r := kafkago.NewReader(kafkago.ReaderConfig{
		Brokers:  brokers,
		Topic:    topic,
		GroupID:  "admin-notifications",
		MinBytes: 1,
		MaxBytes: 10e6,
	})
	return &NotificationConsumer{reader: r, svc: svc, hub: hub}
}

func (c *NotificationConsumer) Run(ctx context.Context) {
	backoff := initialBackoff

	for {
		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			logger.Log.Warn("notification kafka read error, retrying",
				zap.Error(err),
				zap.Duration("backoff", backoff),
			)
			select {
			case <-time.After(backoff):
			case <-ctx.Done():
				return
			}
			backoff *= backoffFactor
			if backoff > maxBackoff {
				backoff = maxBackoff
			}
			continue
		}

		backoff = initialBackoff

		var event events.NotificationEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			logger.Log.Error("notification consumer: unmarshal event", zap.Error(err))
			continue
		}

		// Сохраняем в БД и получаем список созданных уведомлений
		created, err := c.svc.Dispatch(ctx, event)
		if err != nil {
			logger.Log.Error("notification consumer: dispatch", zap.Error(err))
			continue
		}

		// Пушим в активные WS-сессии
		for _, n := range created {
			c.hub.Push(n)
		}
	}
}

func (c *NotificationConsumer) Close() error {
	return c.reader.Close()
}

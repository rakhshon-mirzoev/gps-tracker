package kafka

import (
	"context"
	"time"

	"github.com/IkbolNazarov/gps-tracker/internal/logger"
	"github.com/IkbolNazarov/gps-tracker/internal/websocket"
	kafkago "github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

const (
	initialBackoff = 1 * time.Second
	maxBackoff     = 30 * time.Second
	backoffFactor  = 2
)

type LocationConsumer struct {
	reader *kafkago.Reader
	hub    *websocket.Hub
}

func NewLocationConsumer(brokers []string, topic string, hub *websocket.Hub) *LocationConsumer {
	r := kafkago.NewReader(kafkago.ReaderConfig{
		Brokers:  brokers,
		Topic:    topic,
		GroupID:  "admin-ws",
		MinBytes: 1,
		MaxBytes: 10e6,
	})
	return &LocationConsumer{reader: r, hub: hub}
}

func (c *LocationConsumer) Run(ctx context.Context) {
	backoff := initialBackoff

	for {
		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return
			}

			logger.Log.Warn("kafka read error, retrying",
				zap.Error(err),
				zap.Duration("backoff", backoff),
			)

			select {
			case <-time.After(backoff):
			case <-ctx.Done():
				return
			}

			// Экспоненциальный backoff до maxBackoff
			backoff *= backoffFactor
			if backoff > maxBackoff {
				backoff = maxBackoff
			}

			continue
		}

		// Сброс backoff при успешном чтении
		backoff = initialBackoff

		select {
		case c.hub.Broadcast <- msg.Value:
		default:
			logger.Log.Warn("hub broadcast full, dropping kafka message")
		}
	}
}

func (c *LocationConsumer) Close() error {
	return c.reader.Close()
}

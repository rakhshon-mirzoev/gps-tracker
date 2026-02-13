package kafka

import (
	"context"

	"github.com/rakhshon-mirzoev/gps-tracker/internal/logger"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/websocket"
	kafkago "github.com/segmentio/kafka-go"
	"go.uber.org/zap"
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
	for {
		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			logger.Log.Error("kafka read error", zap.Error(err))
			continue
		}

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

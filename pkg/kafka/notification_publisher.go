package kafka

import (
	"context"
	"encoding/json"

	"github.com/IkbolNazarov/gps-tracker/internal/events"
	kafkago "github.com/segmentio/kafka-go"
)

// NotificationPublisher публикует NotificationEvent в Kafka-топик уведомлений.
type NotificationPublisher struct {
	writer *kafkago.Writer
}

func NewNotificationPublisher(brokers []string, topic string) *NotificationPublisher {
	w := &kafkago.Writer{
		Addr:     kafkago.TCP(brokers...),
		Topic:    topic,
		Balancer: &kafkago.LeastBytes{},
	}
	return &NotificationPublisher{writer: w}
}

func (p *NotificationPublisher) PublishNotification(ctx context.Context, event events.NotificationEvent) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return p.writer.WriteMessages(ctx, kafkago.Message{
		Key:   []byte(event.StaffID.String()),
		Value: payload,
	})
}

func (p *NotificationPublisher) Close() error {
	return p.writer.Close()
}

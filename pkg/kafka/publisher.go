package kafka

import (
	"context"
	"encoding/json"

	"github.com/IkbolNazarov/gps-tracker/internal/events"
	kafkago "github.com/segmentio/kafka-go"
)

type LocationPublisher struct {
	writer *kafkago.Writer
}

func NewLocationPublisher(brokers []string, topic string) *LocationPublisher {
	w := &kafkago.Writer{
		Addr:     kafkago.TCP(brokers...),
		Topic:    topic,
		Balancer: &kafkago.LeastBytes{},
	}
	return &LocationPublisher{writer: w}
}

func (p *LocationPublisher) PublishOnlineLocation(ctx context.Context, update events.LocationUpdate) error {
	payload, err := json.Marshal(update)
	if err != nil {
		return err
	}

	return p.writer.WriteMessages(ctx, kafkago.Message{
		Key:   []byte(update.StaffID.String()),
		Value: payload,
	})
}

func (p *LocationPublisher) Close() error {
	return p.writer.Close()
}

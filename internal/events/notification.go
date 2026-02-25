package events

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// NotificationEvent — событие, которое мобильный сервис кидает в Kafka
// когда сотрудник не отправлял GPS более заданного времени.
type NotificationEvent struct {
	StaffID   uuid.UUID `json:"staff_id"`
	GroupID   uuid.UUID `json:"group_id"`
	FullName  string    `json:"full_name"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"created_at"`
}

type NotificationPublisher interface {
	PublishNotification(ctx context.Context, event NotificationEvent) error
}

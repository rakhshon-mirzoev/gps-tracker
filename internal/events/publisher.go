package events

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type LocationUpdate struct {
	StaffID      uuid.UUID `json:"staff_id"`
	GroupID      uuid.UUID `json:"group_id"`
	FullName     string    `json:"full_name,omitempty"`
	Latitude     float64   `json:"latitude"`
	Longitude    float64   `json:"longitude"`
	BatteryLevel int       `json:"battery_level"`
	CapturedAt   time.Time `json:"captured_at"`
	Status       string    `json:"status"`
}

type LocationPublisher interface {
	PublishOnlineLocation(ctx context.Context, update LocationUpdate) error
}

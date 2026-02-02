package model

import (
	"time"

	"github.com/google/uuid"
)

type LocationHistory struct {
	ID           uuid.UUID `json:"id" db:"id"`
	StaffID      uuid.UUID `json:"staff_id" db:"staff_id"`
	Latitude     float64   `json:"latitude" db:"latitude"`
	Longitude    float64   `json:"longitude" db:"longitude"`
	BatteryLevel int       `json:"battery_level" db:"battery_level"`
	CapturedAt   time.Time `json:"captured_at" db:"captured_at"`
	Status       string    `json:"status" db:"status"`
}

type UpdateLocationHistory struct {
	StaffID      *uuid.UUID `json:"staff_id"`
	Latitude     *float64   `json:"latitude"`
	Longitude    *float64   `json:"longitude"`
	BatteryLevel *int       `json:"battery_level"`
	CapturedAt   *time.Time `json:"captured_at"`
}

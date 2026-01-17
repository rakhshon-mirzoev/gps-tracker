package model

import (
    "time"

    "github.com/google/uuid"
)

type LocationHistory struct {
    ID           uuid.UUID `json:"id" db:"id"`
    StaffID      uuid.UUID `json:"staff_id" db:"staff_id"`
    Location     string    `json:"location" db:"location"` // WKT-формат для PostGIS
    BatteryLevel *int      `json:"battery_level" db:"battery_level"`
    CapturedAt   time.Time `json:"captured_at" db:"captured_at"`
}
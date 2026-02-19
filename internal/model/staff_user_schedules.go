package model

import (
	"time"

	"github.com/google/uuid"
)

type StaffUserSchedule struct {
	ID        uuid.UUID
	StaffID   uuid.UUID
	Weekday   int16
	StartTime time.Time
	EndTime   time.Time
	CreatedAt time.Time
}

type ScheduleDTO struct {
	Weekday   int16     `json:"weekday"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

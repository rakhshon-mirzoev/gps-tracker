package model

import (
	"time"

	"github.com/google/uuid"
)

type StaffUserSchedule struct {
	ID        uuid.UUID `json:"id" db:"id"`
	StaffID   uuid.UUID `json:"staff_id" db:"staff_id"`
	Weekday   int16     `json:"weekday" db:"weekday"`
	StartTime time.Time `json:"start_time" db:"start_time"`
	EndTime   time.Time `json:"end_time" db:"end_time"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type ScheduleDTO struct {
	Weekday   int16     `json:"weekday"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

type UpdateStaffSchedule struct {
	Weekday   *int16     `json:"weekday"`
	StartTime *time.Time `json:"start_time"`
	EndTime   *time.Time `json:"end_time"`
}

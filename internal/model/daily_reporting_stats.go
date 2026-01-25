package model

import (
	"time"

	"github.com/google/uuid"
)

type DailyReportingStats struct {
	ID             uuid.UUID  `json:"id" db:"id"`
	Date           time.Time  `json:"date" db:"date"`
	StaffID        uuid.UUID  `json:"staff_id" db:"staff_id"`
	GroupID        uuid.UUID  `json:"group_id" db:"group_id"`
	PlannedCount   int        `json:"planned_count" db:"planned_count"`
	CompletedCount int        `json:"completed_count" db:"completed_count"`
	RejectedCount  int        `json:"rejected_count" db:"rejected_count"`
	MissedCount    int        `json:"missed_count" db:"missed_count"`
	WorkStartAt    *time.Time `json:"work_start_at" db:"work_start_at"`
	WorkEndAt      *time.Time `json:"work_end_at" db:"work_end_at"`
	TotalKm        *float64   `json:"total_km" db:"total_km"`
}

type UpdateDailyReportingStats struct {
	Date           *time.Time `json:"date"`
	StaffID        *uuid.UUID `json:"staff_id"`
	GroupID        *uuid.UUID `json:"group_id"`
	PlannedCount   *int       `json:"planned_count"`
	CompletedCount *int       `json:"completed_count"`
	RejectedCount  *int       `json:"rejected_count"`
	MissedCount    *int       `json:"missed_count"`
	WorkStartAt    *time.Time `json:"work_start_at"`
	WorkEndAt      *time.Time `json:"work_end_at"`
	TotalKm        *float64   `json:"total_km"`
}

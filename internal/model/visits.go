package model

import (
	"time"

	"github.com/google/uuid"
)

type Visit struct {
	ID                uuid.UUID  `json:"id" db:"id"`
	ScheduleID        *uuid.UUID `json:"schedule_id" db:"schedule_id"`
	StaffID           uuid.UUID  `json:"staff_id" db:"staff_id"`
	StoreID           uuid.UUID  `json:"store_id" db:"store_id"`
	CheckInAt         *time.Time `json:"check_in_at" db:"check_in_at"`
	CheckOutAt        *time.Time `json:"check_out_at" db:"check_out_at"`
	CheckInLat        *float64   `json:"check_in_lat" db:"check_in_lat"`
	CheckInLon        *float64   `json:"check_in_lon" db:"check_in_lon"`
	DistanceDeviation int        `json:"distance_deviation" db:"distance_deviation"`
	SelfieURL         *string    `json:"selfie_url" db:"selfie_url"`
	Comment           *string    `json:"comment" db:"comment"`
	IsOffline         bool       `json:"is_offline" db:"is_offline"`
	ReviewStatus      string     `json:"review_status" db:"review_status"`
	ReviewComment     *string    `json:"review_comment" db:"review_comment"`
	ReviewerID        *uuid.UUID `json:"reviewer_id" db:"reviewer_id"`
	ReviewedAt        *time.Time `json:"reviewed_at" db:"reviewed_at"`
	CreatedAt         time.Time  `json:"created_at" db:"created_at"`
}

type UpdateVisit struct {
	ScheduleID        *string    `json:"schedule_id"`
	StoreID           *string    `json:"store_id"`
	CheckInAt         *time.Time `json:"check_in_at"`
	CheckOutAt        *time.Time `json:"check_out_at"`
	CheckInLat        *float64   `json:"check_in_lat"`
	CheckInLon        *float64   `json:"check_in_lon"`
	DistanceDeviation *int       `json:"distance_deviation"`
	SelfieURL         *string    `json:"selfie_url"`
	Comment           *string    `json:"comment"`
	IsOffline         *bool      `json:"is_offline"`
	ReviewStatus      *string    `json:"review_status"`
	ReviewComment     *string    `json:"review_comment"`
	ReviewerID        *uuid.UUID `json:"reviewer_id"`
	ReviewedAt        *time.Time `json:"reviewed_at"`
}

type VisitStaffBrief struct {
	ID          uuid.UUID `json:"id"`
	FullName    string    `json:"full_name"`
	PhoneNumber string    `json:"phone_number,omitempty"`
}

type VisitStoreBrief struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Address   string    `json:"address"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
}

type VisitResponse struct {
	ID                uuid.UUID       `json:"id"`
	ScheduleID        *uuid.UUID      `json:"schedule_id"`
	StaffID           uuid.UUID       `json:"staff_id"`
	StoreID           uuid.UUID       `json:"store_id"`
	CheckInAt         *time.Time      `json:"check_in_at"`
	CheckOutAt        *time.Time      `json:"check_out_at"`
	CheckInLat        *float64        `json:"check_in_lat"`
	CheckInLon        *float64        `json:"check_in_lon"`
	DistanceDeviation int             `json:"distance_deviation"`
	SelfieURL         *string         `json:"selfie_url"`
	ShelfPhotos       []string        `json:"shelf_photos"`
	Comment           *string         `json:"comment"`
	IsOffline         bool            `json:"is_offline"`
	ReviewStatus      string          `json:"review_status"`
	ReviewComment     *string         `json:"review_comment"`
	ReviewerID        *uuid.UUID      `json:"reviewer_id"`
	ReviewedAt        *time.Time      `json:"reviewed_at"`
	CreatedAt         time.Time       `json:"created_at"`
	Staff             *VisitStaffBrief `json:"staff,omitempty"`
	Store             *VisitStoreBrief `json:"store,omitempty"`
}

type CreateVisit struct {
	ScheduleID string  `json:"schedule_id" binding:"required"`
	StoreID    string  `json:"store_id" binding:"required"`
	Lat        float64 `json:"lat" binding:"required"`
	Lon        float64 `json:"lon" binding:"required"`
}

type VisitExcel struct {
	ID                uuid.UUID  `json:"id"`
	Staff             StaffUser  `json:"staff_id"`
	Store             Store      `json:"store_id"`
	Group             Group      `json:"group"`
	CheckInAt         time.Time  `json:"check_in_at"`
	CheckOutAt        time.Time  `json:"check_out_at"`
	CheckInLat        float64    `json:"check_in_lat"`
	CheckInLon        float64    `json:"check_in_lon"`
	DistanceDeviation int        `json:"distance_deviation"`
	SelfieURL         string     `json:"selfie_url"`
	Comment           string     `json:"comment"`
	ReviewStatus      string     `json:"review_status"`
	ReviewComment     string     `json:"review_comment"`
	Reviewer          SystemUser `json:"reviewer"`
	ReviewedAt        time.Time  `json:"reviewed_at"`
	CreatedAt         time.Time  `json:"created_at"`
}

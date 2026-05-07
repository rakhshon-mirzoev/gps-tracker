package model

import (
	"time"

	"github.com/google/uuid"
)

type Schedule struct {
	ID                uuid.UUID  `json:"id" db:"id"`
	StaffID           uuid.UUID  `json:"staff_id" db:"staff_id"`
	StoreID           uuid.UUID  `json:"store_id" db:"store_id"`
	StoreName         string     `json:"store_name"`
	PlannedDate       time.Time  `json:"planned_date" db:"planned_date"`
	Description       string     `json:"description" db:"description"`
	WindowOpen        time.Time  `json:"window_open" db:"window_open"`
	WindowClose       time.Time  `json:"window_close" db:"window_close"`
	IsPhotoRequired   *bool      `json:"is_photo_required" db:"is_photo_required"`
	IsSelfieRequired  *bool      `json:"is_selfie_required" db:"is_selfie_required"`
	IsCommentRequired *bool      `json:"is_comment_required" db:"is_comment_required"`
	Status            string     `json:"status" db:"status"`
	CreatedByAdminID  uuid.UUID  `json:"created_by_admin_id" db:"created_by_admin_id"`
	ReviewComment     *string    `json:"review_comment" db:"review_comment"`
	IsDeleted         bool       `json:"is_deleted" db:"is_deleted"`
	DeletedAt         *time.Time `json:"deleted_at" db:"deleted_at"`
	VisitId           *uuid.UUID `json:"visit_id"`
	IsUnplanned       bool       `json:"is_unplanned" db:"is_unplanned"`
	CreatedAtLat      *float64   `json:"created_at_lat" db:"created_at_lat"`
	CreatedAtLon      *float64   `json:"created_at_lon" db:"created_at_lon"`
}

type CreateSchedule struct {
	StaffID           uuid.UUID `json:"staff_id"`
	StoreID           uuid.UUID `json:"store_id" binding:"required"`
	PlannedDate       time.Time `json:"planned_date" binding:"required"`
	Description       string    `json:"description" binding:"required"`
	WindowOpen        time.Time `json:"window_open" binding:"required"`
	WindowClose       time.Time `json:"window_close" binding:"required"`
	IsPhotoRequired   bool      `json:"is_photo_required"`
	IsSelfieRequired  bool      `json:"is_selfie_required"`
	IsCommentRequired bool      `json:"is_comment_required"`
	Lat               *float64  `json:"lat,omitempty"`
	Lon               *float64  `json:"lon,omitempty"`
}

type UpdateSchedule struct {
	StaffID           *uuid.UUID `json:"staff_id"`
	StoreID           *uuid.UUID `json:"store_id"`
	PlannedDate       *time.Time `json:"planned_date"`
	Description       *string    `json:"description"`
	WindowOpen        *time.Time `json:"window_open"`
	WindowClose       *time.Time `json:"window_close"`
	IsPhotoRequired   *bool      `json:"is_photo_required"`
	IsSelfieRequired  *bool      `json:"is_selfie_required"`
	IsCommentRequired *bool      `json:"is_comment_required"`
	Status            *string    `json:"status"`
	ReviewComment     *string    `json:"review_comment"`
}

type CreateSchedulesBulk struct {
	StaffID           string    `json:"staff_id"`
	GroupID           *string   `json:"group_id,omitempty"`
	StoreIDs          []string  `json:"store_ids"`
	ExcludeStoreIDs   []string  `json:"exclude_store_ids,omitempty"`
	PlannedDate       time.Time `json:"planned_date"`
	Description       string    `json:"description"`
	WindowOpen        time.Time `json:"window_open"`
	WindowClose       time.Time `json:"window_close"`
	IsPhotoRequired   bool      `json:"is_photo_required"`
	IsSelfieRequired  bool      `json:"is_selfie_required"`
	IsCommentRequired bool      `json:"is_comment_required"`
}

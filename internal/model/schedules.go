package model

import (
	"time"

	"github.com/google/uuid"
)

type Schedule struct {
	ID                uuid.UUID `json:"id" db:"id"`
	StaffID           uuid.UUID `json:"staff_id" db:"staff_id"`
	StoreID           uuid.UUID `json:"store_id" db:"store_id"`
	StoreName         string    `json:"store_name"`
	PlannedDate       time.Time `json:"planned_date" db:"planned_date"`
	Description       string    `json:"description" db:"description"`
	WindowOpen        time.Time `json:"window_open" db:"window_open"`
	WindowClose       time.Time `json:"window_close" db:"window_close"`
	IsPhotoRequired   bool      `json:"is_photo_required" db:"is_photo_required"`
	IsSelfieRequired  bool      `json:"is_selfie_required" db:"is_selfie_required"`
	IsCommentRequired bool      `json:"is_comment_required" db:"is_comment_required"`
	Status            string    `json:"status" db:"status"`
	CreatedByAdminID  uuid.UUID `json:"created_by_admin_id" db:"created_by_admin_id"`
}

type CreateSchedule struct {
	StaffID           string    `json:"staff_id" binding:"required"`
	StoreID           string    `json:"store_id" binding:"required"`
	PlannedDate       time.Time `json:"planned_date" binding:"required"`
	Description       string    `json:"description" binding:"required"`
	WindowOpen        time.Time `json:"window_open" binding:"required"`
	WindowClose       time.Time `json:"window_close" binding:"required"`
	IsPhotoRequired   bool      `json:"is_photo_required"`
	IsSelfieRequired  bool      `json:"is_selfie_required"`
	IsCommentRequired bool      `json:"is_comment_required"`
}

type UpdateSchedule struct {
	StaffID           *string    `json:"staff_id"`
	StoreID           *string    `json:"store_id"`
	PlannedDate       *time.Time `json:"planned_date"`
	Description       *string    `json:"description"`
	WindowOpen        *time.Time `json:"window_open"`
	WindowClose       *time.Time `json:"window_close"`
	IsPhotoRequired   *bool      `json:"is_photo_required"`
	IsSelfieRequired  *bool      `json:"is_selfie_required"`
	IsCommentRequired *bool      `json:"is_comment_required"`
	Status            *string    `json:"status"`
}

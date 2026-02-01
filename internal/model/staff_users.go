package model

import (
	"time"

	"github.com/google/uuid"
)

type StaffUser struct {
	ID                  uuid.UUID  `db:"id" json:"id"`
	GroupID             uuid.UUID  `db:"group_id" json:"group_id"`
	FullName            string     `db:"full_name" json:"full_name"`
	Login               string     `db:"login" json:"login"`
	PhoneNumber         string     `db:"phone_number" json:"phone_number"`
	Role                string     `json:"role" db:"role"`
	Inn                 string     `db:"inn" json:"inn"`
	Birthday            *time.Time `db:"birthday" json:"birthday"`
	PasswordHash        string     `db:"password_hash" json:"-"`
	IsTemporaryPassword bool       `db:"is_temporary_password" json:"is_temporary_password"`
	PasswordUpdatedAt   *time.Time `db:"password_updated_at" json:"password_updated_at"`
	DeviceID            *string    `db:"device_id" json:"device_id"`
	FcmToken            *string    `db:"fcm_token" json:"fcm_token"`
	AppVersion          *string    `db:"app_version" json:"app_version"`
	LastBatteryLevel    *int       `db:"last_battery_level" json:"last_battery_level"`
	LastSeenAt          *time.Time `db:"last_seen_at" json:"last_seen_at"`
	IsActive            bool       `db:"is_active" json:"is_active"`
	DeactivatedAt       *time.Time `db:"deactivated_at" json:"deactivated_at"`
	DeactivatedByID     *uuid.UUID `db:"deactivated_by_id" json:"deactivated_by_id"`
	DeactivationReason  *string    `db:"deactivation_reason" json:"deactivation_reason"`
	CreatedAt           time.Time  `db:"created_at" json:"created_at"`
}

type UpdateStaffUser struct {
	FullName           *string    `json:"full_name"`
	PhoneNumber        *string    `json:"phone_number"`
	IsActive           *bool      `json:"is_active"`
	GroupID            *uuid.UUID `json:"group_id"`
	FcmToken           *string    `json:"fcm_token"`
	DeviceId           *string    `json:"device_id"`
	DeactivatedAt      *time.Time `json:"deactivated_at"`
	DeactivatedByID    *uuid.UUID `json:"deactivated_by_id"`
	DeactivationReason *string    `json:"deactivation_reason"`
}

type CreateStaffUser struct {
	GroupID         string     `json:"group_id" binding:"required"`
	FullName        string     `json:"full_name" binding:"required"`
	Login           string     `json:"login" binding:"required"`
	PhoneNumber     string     `json:"phone_number"`
	Inn             string     `json:"inn"`
	Birthday        *time.Time `json:"birthday"`
	Password        string     `json:"password" binding:"required"`
	ConfirmPassword string     `json:"confirm_password" binding:"required"`
}

type StaffUserResponse struct {
	ID         uuid.UUID  `json:"id"`
	FullName   string     `json:"full_name"`
	Login      string     `json:"login"`
	IsActive   bool       `json:"is_active"`
	LastSeenAt *time.Time `json:"last_seen_at"`
}

type StaffUserProfile struct {
	User       *StaffUser  `json:"user"`
	Schedules  []*Schedule `json:"schedules"`
	TotalCount int         `json:"total_count"`
}

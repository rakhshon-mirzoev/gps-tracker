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
	AvatarKey           *string    `json:"-" db:"avatar_key"`
	AvatarUrl           *string    `json:"avatar_url,omitempty"`
	StartAt             *time.Time `json:"start_at" db:"start_at"`
	EndAt               *time.Time `json:"end_at" db:"end_at"`
	DayStarted          *bool      `json:"day_started" db:"day_started"`
	GroupName           string     `json:"group_name"`
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
	AvatarKey          *string    `json:"avatar_key"`
	StartAt            *time.Time `json:"start_at"`
	EndAt              *time.Time `json:"end_at"`
	DayStarted         *bool      `json:"day_started"`
}

type DeactivateStaffUser struct {
	DeactivatedAt      *time.Time `json:"deactivated_at"`
	DeactivatedByID    *uuid.UUID `json:"deactivated_by_id"`
	DeactivationReason *string    `json:"deactivation_reason"`
}

type CreateStaffUser struct {
	GroupID         string        `json:"group_id"`
	FullName        string        `json:"full_name" binding:"required"`
	Login           string        `json:"login" binding:"required"`
	PhoneNumber     string        `json:"phone_number" binding:"required"`
	Inn             string        `json:"inn" binding:"required"`
	Birthday        *time.Time    `json:"birthday,omitempty"`
	Password        string        `json:"password" binding:"required"`
	ConfirmPassword string        `json:"confirm_password" binding:"required"`
	Schedule        []ScheduleDTO `json:"schedule"`
}

type StaffUserProfile struct {
	User       *StaffUser          `json:"user"`
	Visits     []Visit             `json:"visits"`
	Schedules  []StaffUserSchedule `json:"schedule"`
	TotalCount int                 `json:"total_count"`
	DoneMonth  int                 `json:"done_count_month"`
	DoneWeek   int                 `json:"done_count_week"`
	TotalMonth int                 `json:"total_count_month"`
	TotalWeek  int                 `json:"total_count_week"`
}

type StaffUserLocation struct {
	ID               uuid.UUID        `json:"id"`
	GroupID          uuid.UUID        `json:"group_id"`
	FullName         string           `json:"full_name"`
	Login            string           `json:"login"`
	PhoneNumber      string           `json:"phone_number"`
	Inn              string           `json:"inn"`
	Birthday         *time.Time       `json:"birthday"`
	DeviceID         *string          `json:"device_id"`
	FcmToken         *string          `json:"fcm_token"`
	AppVersion       *string          `json:"app_version"`
	LastBatteryLevel *int             `json:"last_battery_level"`
	StartAt          *time.Time       `json:"start_at"`
	EndAt            *time.Time       `json:"end_at"`
	DayStarted       *bool            `json:"day_started"`
	LastSeenAt       *time.Time       `json:"last_seen_at"`
	Location         *LocationHistory `json:"location"`
	GroupName        *string          `json:"group_name"`
	AvatarKey        *string          `json:"-" db:"avatar_key"`
	AvatarUrl        *string          `json:"avatar_url"`
}

type StaffUserLocationHistory struct {
	ID               uuid.UUID         `json:"id"`
	GroupID          uuid.UUID         `json:"group_id"`
	FullName         string            `json:"full_name"`
	Login            string            `json:"login"`
	PhoneNumber      string            `json:"phone_number"`
	Inn              string            `json:"inn"`
	Birthday         *time.Time        `json:"birthday"`
	DeviceID         *string           `json:"device_id"`
	FcmToken         *string           `json:"fcm_token"`
	AppVersion       *string           `json:"app_version"`
	LastBatteryLevel *int              `json:"last_battery_level"`
	LastSeenAt       *time.Time        `json:"last_seen_at"`
	StartAt          *time.Time        `json:"start_at"`
	EndAt            *time.Time        `json:"end_at"`
	DayStarted       *bool             `json:"day_started"`
	GroupName        string            `json:"group_name"`
	Locations        []LocationHistory `json:"location"`
}

type StaffUserSchedules struct {
	StaffUser StaffUser           `json:"staff_user"`
	Schedules []StaffUserSchedule `json:"schedule"`
}

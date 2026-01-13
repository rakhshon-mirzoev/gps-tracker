package model

import (
	"time"

	"github.com/google/uuid"
)

type SystemUser struct {
	ID                 uuid.UUID  `db:"id" json:"id"`
	FullName           string     `db:"full_name" json:"full_name"`
	Login              string     `db:"login" json:"login"`
	PasswordHash       string     `db:"password_hash" json:"-"`
	Role               string     `db:"role" json:"role"`
	IsActive           bool       `db:"is_active" json:"is_active"`
	DeactivatedAt      *time.Time `db:"deactivated_at" json:"deactivated_at"`
	DeactivatedByID    *int64     `db:"deactivated_by_id" json:"deactivated_by_id"`
	DeactivationReason *string    `db:"deactivation_reason" json:"deactivation_reason"`
	CreatedAt          time.Time  `db:"created_at" json:"created_at"`
}

type CreateSystemUser struct {
	FullName string `json:"full_name" binding:"required"`
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"role" binding:"required"`
}

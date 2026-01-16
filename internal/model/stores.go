package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Store struct {
	ID uuid.UUID `json:"id"`

	Name    string `json:"name"`
	Address string `json:"address"`

	Location pgtype.Point `json:"location"`

	OwnerFullName string  `json:"owner_full_name"`
	OwnerINN      string  `json:"owner_inn"`
	OwnerPassport *string `json:"owner_passport,omitempty"`

	Status           string  `json:"status"`
	RejectionComment *string `json:"rejection_comment,omitempty"`

	CreatedByStaffID uuid.UUID `json:"created_by_staff_id"`
	CreatedByAdminID uuid.UUID `json:"created_by_admin_id"`

	IsActive bool `json:"is_active"`
	IsFrozen bool `json:"is_frozen"`

	DeactivatedAt      *time.Time `json:"deactivated_at,omitempty"`
	DeactivatedByID    *uuid.UUID `json:"deactivated_by_id,omitempty"`
	DeactivationReason *string    `json:"deactivation_reason,omitempty"`

	CreatedAt time.Time `json:"created_at"`
}

type CreateStore struct {
	Name    string `json:"name" binding:"required"`
	Address string `json:"address" binding:"required"`

	Latitude  float64 `json:"latitude" binding:"required"`
	Longitude float64 `json:"longitude" binding:"required"`

	OwnerFullName string  `json:"owner_full_name" binding:"required"`
	OwnerINN      string  `json:"owner_inn" binding:"required,len=12"`
	OwnerPassport *string `json:"owner_passport"`

	CreatedByStaffID uuid.UUID `json:"created_by_staff_id" binding:"required"`
	CreatedByAdminID uuid.UUID `json:"created_by_admin_id" binding:"required"`
}

type StoreListFilter struct {
	Status    *string
	IsActive  *bool
	Limit     int
	Offset    int
	SortBy    string
	SortOrder string
}

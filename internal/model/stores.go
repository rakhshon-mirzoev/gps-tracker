package model

import (
	"time"

	"github.com/google/uuid"
)

type Store struct {
	ID                 uuid.UUID  `json:"id" db:"id"`
	Name               string     `json:"name" db:"name"`
	Address            string     `json:"address" db:"address"`
	Latitude           float64    `json:"latitude" db:"latitude"`
	Longitude          float64    `json:"longitude" db:"longitude"`
	OwnerFullName      string     `json:"owner_full_name" db:"owner_full_name"`
	OwnerINN           string     `json:"owner_inn" db:"owner_inn"`
	OwnerPassport      *string    `json:"owner_passport,omitempty" db:"owner_passport"`
	OwnerPhone         *string    `json:"owner_phone,omitempty" db:"owner_phone"`
	AdditionalPhone    *string    `json:"additional_phone,omitempty" db:"additional_phone"`
	Comment            *string    `json:"comment,omitempty" db:"comment"`
	Status             string     `json:"status" db:"status"`
	RejectionComment   *string    `json:"rejection_comment,omitempty" db:"rejection_comment"`
	CreatedByStaffID   uuid.UUID  `json:"created_by_staff_id" db:"created_by_staff_id"`
	IsActive           bool       `json:"is_active" db:"is_active"`
	IsFrozen           bool       `json:"is_frozen" db:"is_frozen"`
	DeactivatedAt      *time.Time `json:"deactivated_at,omitempty" db:"deactivated_at"`
	DeactivatedByID    *uuid.UUID `json:"deactivated_by_id,omitempty" db:"deactivated_by_id"`
	DeactivationReason *string    `json:"deactivation_reason,omitempty" db:"deactivation_reason"`
	CreatedAt          time.Time  `json:"created_at" db:"created_at"`
	StaffFullName      string     `json:"staff_name"`
	StaffPhoneNumber   string     `json:"staff_phone"`
}

type CreateStore struct {
	Name            string   `json:"name" binding:"required"`
	Address         string   `json:"address" binding:"required"`
	Latitude        float64  `json:"latitude" binding:"required"`
	Longitude       float64  `json:"longitude" binding:"required"`
	OwnerFullName   string   `json:"owner_full_name" binding:"required"`
	OwnerINN        string   `json:"owner_inn" binding:"required"`
	OwnerPhone      string   `json:"owner_phone" binding:"required"`
	OwnerPassport   *string  `json:"owner_passport"`
	AdditionalPhone *string  `json:"additional_phone"`
	Comment         *string  `json:"comment"`
	InteriorPhotos  []string `json:"interior_photos"`
	ExteriorPhotos  []string `json:"exterior_photos"`
}

type UpdateStore struct {
	Name               *string    `json:"name"`
	Address            *string    `json:"address"`
	Latitude           *float64   `json:"latitude"`
	Longitude          *float64   `json:"longitude"`
	OwnerFullName      *string    `json:"owner_full_name"`
	OwnerINN           *string    `json:"owner_inn"`
	OwnerPassport      *string    `json:"owner_passport"`
	OwnerPhone         *string    `json:"owner_phone"`
	AdditionalPhone    *string    `json:"additional_phone"`
	Comment            *string    `json:"comment"`
	CreatedByStaffID   *string    `json:"created_by_staff_id"`
	Status             *string    `json:"status"`
	RejectionComment   *string    `json:"rejection_comment"`
	IsActive           *bool      `json:"is_active"`
	IsFrozen           *bool      `json:"is_frozen"`
	DeactivatedAt      *time.Time `json:"deactivated_at,omitempty"`
	DeactivatedByID    *string    `json:"deactivated_by_id,omitempty"`
	DeactivationReason *string    `json:"deactivation_reason,omitempty"`
}

type StoreWithPhoto struct {
	Store       Store         `json:"store"`
	StorePhotos []StorePhotos `json:"photos"`
}

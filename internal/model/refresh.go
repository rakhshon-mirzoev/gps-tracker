package model

import (
	"time"

	"github.com/google/uuid"
)

type Refresh struct {
	ID        uuid.UUID `json:"id" db:"id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	UserType  string    `json:"user_type" db:"user_type"`
	Role      string    `json:"role" db:"role"`
	Token     string    `json:"token" db:"token"`
	ExpiresAt time.Time `json:"expires_at" db:"expires_at"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

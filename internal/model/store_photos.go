package model

import (
	"time"

	"github.com/google/uuid"
)

type StorePhotos struct {
	ID        uuid.UUID `json:"id" db:"id"`
	StoreID   uuid.UUID `json:"store_id" db:"store_id"`
	PhotoURL  string    `json:"photo_url" db:"photo_url"`
	Type      string    `json:"type" db:"type"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

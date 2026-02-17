package model

import (
	"time"

	"github.com/google/uuid"
)

type StorePhotos struct {
	ID        uuid.UUID `json:"id" db:"id"`
	StoreID   uuid.UUID `json:"store_id" db:"store_id"`
	ObjectKey string    `json:"-" db:"object_key"`
	URL       string    `json:"url"`
	Type      string    `json:"type" db:"type"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

package model

import (
	"time"

	"github.com/google/uuid"
)

type VisitPhotos struct {
	ID           uuid.UUID `json:"id" db:"id"`
	VisitID      uuid.UUID `json:"visit_id" db:"visit_id"`
	PhotoUrl     string    `json:"photo_url"`
	PhotoKey     string    `json:"-" db:"photo_key"`
	CapturedAt   time.Time `json:"captured_at" db:"captured_at"`
	OriginalName string    `json:"original_name" db:"original_name"`
	PhotoType    string    `json:"photo_type"`
}

type UpdateVisitPhotos struct {
	VisitID    *uuid.UUID `json:"visit_id"`
	PhotoURL   *string    `json:"photo_url"`
	CapturedAt *time.Time `json:"captured_at"`
}

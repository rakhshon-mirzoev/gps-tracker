package model

import (
	"time"

	"github.com/google/uuid"
)

type VisitPhotos struct {
	ID         uuid.UUID `json:"id" db:"id"`
	VisitID    uuid.UUID `json:"visit_id" db:"visit_id"`
	PhotoURL   string    `json:"photo_url" db:"photo_url"`
	CapturedAt time.Time `json:"captured_at" db:"captured_at"`
}

type UpdateVisitPhotos struct {
	VisitID    *uuid.UUID `json:"visit_id"`
	PhotoURL   *string    `json:"photo_url"`
	CapturedAt *time.Time `json:"captured_at"`
}

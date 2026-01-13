package model

import (
	"time"

	"github.com/google/uuid"
)

type Group struct {
	ID             uuid.UUID  `db:"id" json:"id"`
	Name           string     `db:"name" json:"name"`
	SupervisorID   *uuid.UUID `db:"supervisor_id" json:"supervisor_id,omitempty"`
	SupervisorName *string    `db:"supervisor_name" json:"supervisor_name,omitempty"`
	CreatedAt      time.Time  `db:"created_at" json:"created_at"`
}

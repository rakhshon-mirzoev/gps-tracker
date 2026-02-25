package model

import (
	"time"

	"github.com/google/uuid"
)

// SystemNotification — уведомление для admin/supervisor в панели управления.
type SystemNotification struct {
	ID          uuid.UUID  `json:"id"`
	RecipientID uuid.UUID  `json:"recipient_id"`
	StaffID     *uuid.UUID `json:"staff_id,omitempty"`
	GroupID     *uuid.UUID `json:"group_id,omitempty"`
	Type        string     `json:"type"`
	Title       string     `json:"title"`
	Body        string     `json:"body"`
	IsRead      bool       `json:"is_read"`
	ReadAt      *time.Time `json:"read_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
}

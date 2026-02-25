package model

import (
	"time"

	"encoding/json"

	"github.com/google/uuid"
)

type Notification struct {
	ID         uuid.UUID       `json:"id" db:"id"`
	UserID     *uuid.UUID      `json:"user_id" db:"user_id"`
	Type       string          `json:"type" db:"type"`
	Title      string          `json:"title" db:"title"`
	Body       string          `json:"body" db:"body"`
	Meta       json.RawMessage `json:"meta" db:"meta"`
	IsRead     bool            `json:"is_read" db:"is_read"`
	PushSent   bool            `json:"push_sent" db:"push_sent"`
	ReadAt     time.Time       `json:"read_at" db:"read_at"`
	PushSentAt time.Time       `json:"push_sent_at" db:"push_sent_at"`
	CreatedAt  time.Time       `json:"created_at" db:"created_at"`
}

type CreateNotificationRequest struct {
	UserID  *uuid.UUID     `json:"user_id"`
	GroupID *uuid.UUID     `json:"group_id"`
	Title   string         `json:"title" binding:"required"`
	Body    string         `json:"body" binding:"required"`
	Type    string         `json:"type" binding:"required"`
	Meta    map[string]any `json:"meta"`
}

package model

import (
	"time"

	"encoding/json"

	"github.com/google/uuid"
)

type Notification struct {
	ID        uuid.UUID       `json:"id" db:"id"`
	UserID    *uuid.UUID      `json:"user_id" db:"user_id"`
	GroupID   *uuid.UUID      `json:"group_id" db:"group_id"`
	Type      string          `json:"type" db:"type"`
	Title     string          `json:"title" db:"title"`
	Body      string          `json:"body" db:"body"`
	Meta      json.RawMessage `json:"meta" db:"meta"`
	IsRead    bool            `json:"is_read" db:"is_read"`
	PushSent  bool            `json:"push_sent" db:"push_sent"`
	CreatedAt time.Time       `json:"created_at" db:"created_at"`
}

package model

import (
	"time"

	"github.com/google/uuid"
)

type AuditLogs struct {
	ID          uuid.UUID              `json:"id" db:"id"`
	ActorID     uuid.UUID              `json:"actor_id" db:"actor_id"`
	ActorType   string                 `json:"actor_type" db:"actor_type"`
	Action      string                 `json:"action" db:"action"`
	TargetTable *string                `json:"target_table" db:"target_table"`
	TargetID    *uuid.UUID             `json:"target_id" db:"target_id"`
	OldValue    map[string]interface{} `json:"old_value" db:"old_value"`
	NewValue    map[string]interface{} `json:"new_value" db:"new_value"`
	IPAddress   *string                `json:"ip_address" db:"ip_address"`
	CreatedAt   time.Time              `json:"created_at" db:"created_at"`
}

type UpdateAuditLogs struct {
	ActorID     *string                `json:"actor_id"`
	ActorType   *string                `json:"actor_type"`
	Action      *string                `json:"action"`
	TargetTable *string                `json:"target_table"`
	TargetID    *string                `json:"target_id"`
	OldValue    map[string]interface{} `json:"old_value"`
	NewValue    map[string]interface{} `json:"new_value"`
	IPAddress   *string                `json:"ip_address"`
}

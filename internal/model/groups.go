package model

import (
	"time"

	"github.com/google/uuid"
)

type Group struct {
	ID             uuid.UUID        `db:"id" json:"id"`
	Name           string           `db:"name" json:"name"`
	Type           string           `db:"type" json:"type"` // 'staff' or 'store'
	ParentID       *uuid.UUID       `db:"parent_id" json:"parent_id"`
	ParentName     *string          `json:"parent_name"`
	SupervisorID   *uuid.UUID       `db:"supervisor_id" json:"supervisor_id,omitempty"`
	SupervisorName *string          `db:"supervisor_name" json:"supervisor_name,omitempty"`
	Supervisor     *GroupSupervisor `json:"supervisor,omitempty"`
	Color          string           `db:"color" json:"color"`     // Hex color code for map visualization
	Layer          int              `db:"layer" json:"layer"`     // Layer order for map rendering
	Cluster        string           `db:"cluster" json:"cluster"` // Cluster name for grouping on map
	CreatedAt      time.Time        `db:"created_at" json:"created_at"`
	IsDeleted      bool             `json:"is_deleted" db:"is_deleted"`
	DeletedAt      *time.Time       `json:"deleted_at" db:"deleted_at"`

	Groups []Group `json:"groups"`
}

type GroupSupervisor struct {
	ID       uuid.UUID `json:"id"`
	FullName string    `json:"full_name"`
	Login    string    `json:"login"`
	Role     string    `json:"role"`
}

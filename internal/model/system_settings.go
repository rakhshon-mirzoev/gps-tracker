package model

import (
	"github.com/google/uuid"
)

type SystemSettings struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Key         string    `json:"key" db:"key"`
	Value       string    `json:"value" db:"value"`
	Description *string   `json:"description" db:"description"`
}

type UpdateSystemSettings struct {
	Key         *string `json:"key"`
	Value       *string `json:"value"`
	Description *string `json:"description"`
}

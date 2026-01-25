package model

import (
	"github.com/google/uuid"
)

type StaffStoresMapping struct {
	ID      uuid.UUID `json:"id" db:"id"`
	StaffID uuid.UUID `json:"staff_id" db:"staff_id"`
	StoreID uuid.UUID `json:"store_id" db:"store_id"`
}

type StaffStores struct {
	Staff  StaffUser `json:"staff_user"`
	Stores []*Store  `json:"stores"`
}

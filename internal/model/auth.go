package model

import (
	"time"

	"github.com/google/uuid"
)

type Auth struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type User struct {
	ID           uuid.UUID
	PasswordHash string     `db:"password_hash"`
	Role         string     `db:"role"`
	GroupID      *uuid.UUID `db:"group_id"`
	Type         string     `db:"type"`
}

type CreateUser struct {
	FullName        string     `json:"full_name" binding:"required"`
	Login           string     `json:"login" binding:"required,alphanum"`
	Password        string     `json:"password" binding:"required,min=8,max=20"`
	ConfirmPassword string     `json:"confirm_password" binding:"required,min=8,max=20"`
	Role            string     `json:"role,omitempty"`
	GroupID         *uuid.UUID `json:"group_id,omitempty"`
	PhoneNumber     string     `json:"phone_number,omitempty"`
	Inn             string     `json:"inn,omitempty"`
	Birthday        *time.Time `json:"birthday,omitempty"`
}

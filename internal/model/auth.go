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
	Type         string     `db:"user_type"`
}
type LoginUserInfo struct {
	ID          uuid.UUID `json:"id"`
	FullName    string    `json:"full_name"`
	Login       string    `json:"login"`
	Role        string    `json:"role"`
	PhoneNumber string    `json:"phone_number,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

type LoginResponse struct {
	AccessToken  string        `json:"access_token"`
	RefreshToken string        `json:"refresh_token"`
	User         LoginUserInfo `json:"user"`
}

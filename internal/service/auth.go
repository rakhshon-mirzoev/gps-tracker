package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/rakhshon-mirzoev/gps-tracker/pkg/token"
	"golang.org/x/crypto/bcrypt"
)

type Authorize interface {
	ComparePasswordWithHash(input, hash string) (err error)
	Authenticate(ctx context.Context, login, password, password_hash string, id uuid.UUID) (string, error)
}

type AuthorizeService struct {
	t token.Token
}

func InitAuthorizeService(t token.Token) Authorize {
	return &AuthorizeService{t: t}
}

func (s *AuthorizeService) ComparePasswordWithHash(input, hash string) (err error) {
	if len(hash) < 59 {
		return fmt.Errorf("invalid hash")
	}
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(input))
}

func (s *AuthorizeService) Authenticate(ctx context.Context, login, password, password_hash string, id uuid.UUID) (string, error) {

	err := bcrypt.CompareHashAndPassword([]byte(password_hash), []byte(password))
	if err != nil {
		return "", err
	}

	token, err := s.t.GenerateToken(id)
	if err != nil {
		return "", err
	}

	return token, nil
}

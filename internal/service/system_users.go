package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/model"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type SystemUsers struct {
	r repository.SystemUsersRepository
}

func InitSystemUsersService(r repository.SystemUsersRepository) *SystemUsers {
	return &SystemUsers{r: r}
}

func (s *SystemUsers) Create(ctx context.Context, sysUser *model.CreateSystemUser) error {
	if sysUser.FullName == "" || sysUser.Login == "" || sysUser.Password == "" {
		return errors.New("full_name, login and password are required")
	}
	hash, err := hashPassword(sysUser.Password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	u := &model.SystemUser{
		FullName:     sysUser.FullName,
		Login:        sysUser.Login,
		Role:         sysUser.Role,
		PasswordHash: hash,
	}

	err = s.r.Create(ctx, u)
	if err != nil {
		return err
	}

	return nil
}

func (s *SystemUsers) GetByID(ctx context.Context, id uuid.UUID) (*model.SystemUser, error) {
	return s.r.GetByID(ctx, id)
}

func (s *SystemUsers) GetByLogin(ctx context.Context, login string) (*model.SystemUser, error) {
	return s.r.GetByLogin(ctx, login)
}

func (s *SystemUsers) List(ctx context.Context) ([]model.SystemUser, error) {
	return s.r.List(ctx)
}

func (s *SystemUsers) Update(ctx context.Context, u *model.SystemUser) error {
	return s.r.Update(ctx, u)
}

func (s *SystemUsers) Delete(ctx context.Context, id uuid.UUID) error {
	return s.r.Delete(ctx, id)
}

func (s *SystemUsers) Authenticate(ctx context.Context, login, password string) (*model.SystemUser, error) {
	u, err := s.r.GetByLogin(ctx, login)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	if err != nil {
		return nil, err
	}

	return u, nil
}

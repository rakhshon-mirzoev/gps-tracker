package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/model"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/repository"
	errors "github.com/rakhshon-mirzoev/gps-tracker/pkg/error"
)

type SystemUsers interface {
	Create(ctx context.Context, sysUser *model.CreateSystemUser) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.SystemUser, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Update(ctx context.Context, u *model.SystemUser) error
	List(ctx context.Context) ([]model.SystemUser, error)
	GetByLogin(ctx context.Context, login string) (*model.SystemUser, error)
}

type SystemUserService struct {
	r repository.SystemUsersRepository
	a Authorize
}

func InitSystemUsersService(r repository.SystemUsersRepository) SystemUsers {
	return &SystemUserService{r: r}
}

func (s *SystemUserService) Create(ctx context.Context, sysUser *model.CreateSystemUser) error {

	if sysUser.Password != sysUser.ConfirmPassword {
		return errors.ErrorPasswordsDoNotMatch
	}

	hash, err := hashPassword(sysUser.Password)
	if err != nil {
		return errors.ErrorPasswordHashError
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

func (s *SystemUserService) GetByID(ctx context.Context, id uuid.UUID) (*model.SystemUser, error) {
	return s.r.GetByID(ctx, id)
}

func (s *SystemUserService) GetByLogin(ctx context.Context, login string) (*model.SystemUser, error) {
	return s.r.GetByLogin(ctx, login)
}

func (s *SystemUserService) List(ctx context.Context) ([]model.SystemUser, error) {
	return s.r.List(ctx)
}

func (s *SystemUserService) Update(ctx context.Context, u *model.SystemUser) error {
	return s.r.Update(ctx, u)
}

func (s *SystemUserService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.r.Delete(ctx, id)
}

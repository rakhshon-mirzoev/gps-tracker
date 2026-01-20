package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/model"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/repository"
	errors "github.com/rakhshon-mirzoev/gps-tracker/pkg/error"
	"github.com/rakhshon-mirzoev/gps-tracker/pkg/utils"
)

type SystemUsers interface {
	Create(ctx context.Context, sysUser *model.CreateSystemUser) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.SystemUser, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Update(ctx context.Context, u *model.SystemUser) error
	List(ctx context.Context) ([]model.SystemUser, error)
	GetByLogin(ctx context.Context, login string) (*model.SystemUser, error)
	Deactivate(ctx context.Context, id, dId uuid.UUID, dReason string) error
}

type systemUserService struct {
	r repository.SystemUsersRepository
}

func NewSystemUsersService(r repository.SystemUsersRepository) SystemUsers {
	return &systemUserService{r: r}
}

func (s *systemUserService) Create(ctx context.Context, sysUser *model.CreateSystemUser) error {

	if sysUser.Password != sysUser.ConfirmPassword {
		return errors.ErrorPasswordsDoNotMatch
	}

	hash, err := utils.HashPassword(sysUser.Password)
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

func (s *systemUserService) GetByID(ctx context.Context, id uuid.UUID) (*model.SystemUser, error) {
	return s.r.GetByID(ctx, id)
}

func (s *systemUserService) GetByLogin(ctx context.Context, login string) (*model.SystemUser, error) {
	return s.r.GetByLogin(ctx, login)
}

func (s *systemUserService) List(ctx context.Context) ([]model.SystemUser, error) {
	return s.r.List(ctx)
}

func (s *systemUserService) Update(ctx context.Context, u *model.SystemUser) error {
	return s.r.Update(ctx, u)
}

func (s *systemUserService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.r.Delete(ctx, id)
}

func (s *systemUserService) Deactivate(ctx context.Context, id, dId uuid.UUID, dReason string) error {
	return s.r.Deactivate(ctx, id, dId, dReason)
}

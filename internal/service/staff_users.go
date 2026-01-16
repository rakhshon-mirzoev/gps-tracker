package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/model"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/repository"
	errors "github.com/rakhshon-mirzoev/gps-tracker/pkg/error"
	"github.com/rakhshon-mirzoev/gps-tracker/pkg/utils.go"
)

type StaffUsers interface {
	Create(ctx context.Context, dto model.CreateStaffUser) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.StaffUser, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Update(ctx context.Context, id uuid.UUID, dto model.UpdateStaffUser) error
	List(ctx context.Context) ([]model.StaffUser, error)
	Deactivate(ctx context.Context, id, dId uuid.UUID, dReason string) error
}
type StaffUserService struct {
	r repository.StaffUsersRepository
}

func InitStaffUserService(r repository.StaffUsersRepository) StaffUsers {
	return &StaffUserService{r: r}
}

func (s *StaffUserService) Create(ctx context.Context, dto model.CreateStaffUser) error {

	hash, err := utils.HashPassword(dto.Password)
	if err != nil {
		return errors.ErrorPasswordHashError
	}

	user := &model.StaffUser{
		GroupID:             dto.GroupID,
		FullName:            dto.FullName,
		Login:               dto.Login,
		PhoneNumber:         dto.PhoneNumber,
		Inn:                 dto.Inn,
		Birthday:            dto.Birthday,
		PasswordHash:        hash,
		IsTemporaryPassword: true,
	}

	if err := s.r.Create(ctx, user); err != nil {
		return err
	}

	return nil
}

func (s *StaffUserService) GetByID(ctx context.Context, id uuid.UUID) (*model.StaffUser, error) {

	user, err := s.r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *StaffUserService) List(ctx context.Context) ([]model.StaffUser, error) {
	return s.r.List(ctx)
}

func (s *StaffUserService) Update(ctx context.Context, id uuid.UUID, dto model.UpdateStaffUser) error {

	user, err := s.r.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if user == nil {
		return errors.ErrorNotFound
	}

	if dto.FullName != nil {
		user.FullName = *dto.FullName
	}
	if dto.PhoneNumber != nil {
		user.PhoneNumber = *dto.PhoneNumber
	}
	if dto.IsActive != nil {
		user.IsActive = *dto.IsActive
	}

	return s.r.Update(ctx, user)
}

func (s *StaffUserService) Delete(ctx context.Context, id uuid.UUID) error {

	user, err := s.r.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.ErrorNotFound
	}

	return s.r.Delete(ctx, id)
}

func (s *StaffUserService) Deactivate(ctx context.Context, id, dId uuid.UUID, dReason string) error {
	return s.r.Deactivate(ctx, id, dId, dReason)
}

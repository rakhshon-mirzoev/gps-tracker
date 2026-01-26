package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/model"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/repository"
	errors "github.com/rakhshon-mirzoev/gps-tracker/pkg/error"
	errs "github.com/rakhshon-mirzoev/gps-tracker/pkg/error"
	"github.com/rakhshon-mirzoev/gps-tracker/pkg/utils"
)

type StaffUsers interface {
	Create(ctx context.Context, dto model.CreateStaffUser) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.StaffUser, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Update(ctx context.Context, id uuid.UUID, dto model.UpdateStaffUser) error
	List(ctx context.Context) ([]model.StaffUser, error)
	Deactivate(ctx context.Context, id, dId uuid.UUID, dReason string) error
	GetProfile(ctx context.Context, id uuid.UUID, lim, off int) (*model.StaffUserProfile, error)
}
type staffUserService struct {
	r         repository.StaffUsersRepository
	schedules Schedule
}

func NewStaffUserService(r repository.StaffUsersRepository, sch Schedule) StaffUsers {
	return &staffUserService{r: r, schedules: sch}
}

func (s *staffUserService) GetProfile(ctx context.Context, id uuid.UUID, lim, off int) (*model.StaffUserProfile, error) {

	user, err := s.r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	schedules, total, err := s.schedules.ListByStaffID(ctx, id, lim, off)
	if err != nil {
		return nil, err
	}

	return &model.StaffUserProfile{
		User:       user,
		Schedules:  schedules,
		TotalCount: total,
	}, nil
}

func (s *staffUserService) Create(ctx context.Context, dto model.CreateStaffUser) error {
	if dto.Password != dto.ConfirmPassword {
		return errs.ErrorPasswordsDoNotMatch
	}

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

func (s *staffUserService) GetByID(ctx context.Context, id uuid.UUID) (*model.StaffUser, error) {
	return s.r.GetByID(ctx, id)
}

func (s *staffUserService) List(ctx context.Context) ([]model.StaffUser, error) {
	return s.r.List(ctx)
}

func (s *staffUserService) Update(ctx context.Context, id uuid.UUID, dto model.UpdateStaffUser) error {
	return s.r.Update(ctx, id, dto)
}

func (s *staffUserService) Delete(ctx context.Context, id uuid.UUID) error {

	user, err := s.r.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.ErrorNotFound
	}

	return s.r.Delete(ctx, id)
}

func (s *staffUserService) Deactivate(ctx context.Context, id, dId uuid.UUID, dReason string) error {
	return s.r.Deactivate(ctx, id, dId, dReason)
}

package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/model"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type StaffUsers interface {
	Create(ctx context.Context, dto model.CreateStaffUser) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.StaffUser, error)
	Deactivate(ctx context.Context, id uuid.UUID) error
	Update(ctx context.Context, id uuid.UUID, dto model.UpdateStaffUser) error
	List(ctx context.Context) ([]*model.StaffUser, error)
}
type StaffUserService struct {
	r repository.StaffUsersRepository
}

func InitStaffUserService(r repository.StaffUsersRepository) StaffUsers {
	return &StaffUserService{r: r}
}

func (s *StaffUserService) Create(ctx context.Context, dto model.CreateStaffUser) error {

	if dto.FullName == "" || dto.Login == "" || dto.Password == "" {
		return errors.New("full_name, login and password are required")
	}

	hash, err := hashPassword(dto.Password)
	if err != nil {
		return err
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

	if user == nil {
		return nil, errors.New("staff user not found")
	}

	return user, nil
}

func (s *StaffUserService) List(ctx context.Context) ([]*model.StaffUser, error) {
	return s.r.List(ctx)
}

func (s *StaffUserService) Update(ctx context.Context, id uuid.UUID, dto model.UpdateStaffUser) error {

	user, err := s.r.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("staff user not found")
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

func (s *StaffUserService) Deactivate(ctx context.Context, id uuid.UUID) error {

	user, err := s.r.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("staff user not found")
	}

	return s.r.Delete(ctx, id)
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	return string(bytes), err
}

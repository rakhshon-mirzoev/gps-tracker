package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/model"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/repository"
	"github.com/rakhshon-mirzoev/gps-tracker/pkg/constants"
	errs "github.com/rakhshon-mirzoev/gps-tracker/pkg/error"
	"github.com/rakhshon-mirzoev/gps-tracker/pkg/token"
	"golang.org/x/crypto/bcrypt"
)

type Authorize interface {
	Authenticate(password string, user *model.User) (string, string, error)
	GetByLogin(ctx context.Context, login string) (*model.User, error)
	Create(ctx context.Context, u *model.CreateUser, utype string) error
	UpsertRefreshToken(ctx context.Context, id uuid.UUID, user_type, token string) error
	Refresh(ctx context.Context, oldRefreshToken string) (accessToken string, refreshToken string, err error)
	Logout(ctx context.Context, token string) error
}

type authorizeService struct {
	t     token.Token
	r     repository.Auth
	sys   SystemUsers
	staff StaffUsers
}

func NewAuthorizeService(t token.Token, r repository.Auth, sys SystemUsers, staff StaffUsers) Authorize {
	return &authorizeService{t: t, r: r, sys: sys, staff: staff}
}

func (s *authorizeService) Authenticate(password string, user *model.User) (string, string, error) {
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", "", err
	}

	token, err := s.t.GenerateToken(user.ID, user.GroupID, user.Role)
	if err != nil {
		return "", "", err
	}

	refresh_token, err := s.t.GenerateRefreshToken(user.ID, user.Type)
	if err != nil {
		return "", "", err
	}

	return token, refresh_token.Token, nil
}

func (s *authorizeService) GetByLogin(ctx context.Context, login string) (*model.User, error) {
	return s.r.GetByLoginUnified(ctx, login)
}

func (s *authorizeService) Create(ctx context.Context, u *model.CreateUser, utype string) error {

	switch utype {
	case constants.UserTypeSystem:
		user := &model.CreateSystemUser{
			FullName:        u.FullName,
			Login:           u.Login,
			Role:            u.Role,
			Password:        u.Password,
			ConfirmPassword: u.ConfirmPassword,
		}
		if err := s.sys.Create(ctx, user, uuid.Nil, "", ""); err != nil {
			return err
		}

		return nil
	case constants.UserTypeStaff:
		user := &model.CreateStaffUser{
			GroupID:         *u.GroupID,
			FullName:        u.FullName,
			Login:           u.Login,
			PhoneNumber:     u.PhoneNumber,
			Inn:             u.Inn,
			Birthday:        u.Birthday,
			Password:        u.Password,
			ConfirmPassword: u.ConfirmPassword,
		}

		if err := s.staff.Create(ctx, *user, uuid.Nil, "", "", "", ""); err != nil {
			return err
		}

		return nil
	default:
		return errs.ErrorInvalidLoginType
	}

}

func (s *authorizeService) UpsertRefreshToken(ctx context.Context, id uuid.UUID, user_type, token string) error {
	refresh := &model.Refresh{
		UserID:    id,
		UserType:  user_type,
		Token:     token,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 30),
	}

	return s.r.UpsertRefreshToken(ctx, refresh)

}

func (s *authorizeService) Refresh(ctx context.Context, oldRefreshToken string) (accessToken string, refreshToken string, err error) {
	var (
		gId     *uuid.UUID
		user_id uuid.UUID
		role    string
	)

	ref, err := s.r.GetRefreshByToken(ctx, oldRefreshToken)
	if err != nil {
		return "", "", err
	}

	if time.Now().After(ref.ExpiresAt) {
		_ = s.r.DeleteRefreshByID(ctx, ref.ID)
		return "", "", errs.ErrorInvalidToken
	}

	if err := s.r.DeleteRefreshByID(ctx, ref.ID); err != nil {
		return "", "", err
	}

	if ref.UserType == constants.UserTypeStaff {
		user, err := s.staff.GetByID(ctx, ref.UserID)
		if err != nil {
			return "", "", err
		}
		gId = &user.GroupID
		user_id = user.ID
		role = user.Role
	}

	if ref.UserType == constants.UserTypeSystem {
		user, err := s.sys.GetByID(ctx, ref.UserID)
		if err != nil {
			return "", "", err
		}
		user_id = user.ID
		role = user.Role
	}

	accessToken, err = s.t.GenerateToken(
		user_id,
		gId,
		role,
	)
	if err != nil {
		return "", "", err
	}

	refresh_token, err := s.t.GenerateRefreshToken(ref.UserID, ref.UserType)
	if err != nil {
		return "", "", err
	}

	if err := s.r.UpsertRefreshToken(ctx, refresh_token); err != nil {
		return "", "", err
	}

	return accessToken, refresh_token.Token, nil
}

func (s *authorizeService) Logout(ctx context.Context, token string) error {
	return s.r.DeleteRefreshByToken(ctx, token)
}

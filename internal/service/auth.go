package service

import (
	"context"

	"github.com/rakhshon-mirzoev/gps-tracker/internal/model"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/repository"
	errors "github.com/rakhshon-mirzoev/gps-tracker/pkg/error"
	errs "github.com/rakhshon-mirzoev/gps-tracker/pkg/error"
	"github.com/rakhshon-mirzoev/gps-tracker/pkg/token"
	"github.com/rakhshon-mirzoev/gps-tracker/pkg/utils.go"
	"golang.org/x/crypto/bcrypt"
)

type Authorize interface {
	ComparePasswordWithHash(input, hash string) (err error)
	Authenticate(password string, user *model.User) (string, error)
	GetByLogin(ctx context.Context, login string) (*model.User, error)
	Create(ctx context.Context, u *model.CreateUser, utype string) error
}

type authorizeService struct {
	t     token.Token
	r     repository.Auth
	sys   repository.SystemUsersRepository
	staff repository.StaffUsersRepository
}

func InitAuthorizeService(t token.Token) Authorize {
	return &authorizeService{t: t}
}

func (s *authorizeService) ComparePasswordWithHash(input, hash string) (err error) {
	if len(hash) < 59 {
		return errors.ErrorInvalidHash
	}
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(input))
}

func (s *authorizeService) Authenticate(password string, user *model.User) (string, error) {

	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", err
	}

	token, err := s.t.GenerateToken(user.ID, user.GroupID, user.Role)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *authorizeService) GetByLogin(ctx context.Context, login string) (*model.User, error) {
	return s.r.GetByLoginUnified(ctx, login)
}

func (s *authorizeService) Create(ctx context.Context, u *model.CreateUser, utype string) error {

	if u.Password != u.ConfirmPassword {
		return errors.ErrorPasswordsDoNotMatch
	}

	hash, err := utils.HashPassword(u.Password)
	if err != nil {
		return errors.ErrorPasswordHashError
	}

	switch utype {
	case "system":
		user := &model.SystemUser{
			FullName:     u.FullName,
			Login:        u.Login,
			Role:         u.Role,
			PasswordHash: hash,
		}
		err = s.sys.Create(ctx, user)
		if err != nil {
			return err
		}

		return nil
	case "staff":
		user := &model.StaffUser{
			GroupID:             u.GroupID,
			FullName:            u.FullName,
			Login:               u.Login,
			PhoneNumber:         u.PhoneNumber,
			Inn:                 u.Inn,
			Birthday:            u.Birthday,
			PasswordHash:        hash,
			IsTemporaryPassword: true,
		}

		if err := s.staff.Create(ctx, user); err != nil {
			return err
		}

		return nil
	default:
		return errs.ErrorInvalidLoginType
	}

}

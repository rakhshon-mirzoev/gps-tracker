package service

import (
	"context"
	"time"

	"github.com/IkbolNazarov/gps-tracker/internal/model"
	"github.com/IkbolNazarov/gps-tracker/internal/repository"
	"github.com/IkbolNazarov/gps-tracker/pkg/constants"
	errors "github.com/IkbolNazarov/gps-tracker/pkg/error"
	"github.com/IkbolNazarov/gps-tracker/pkg/timezone"
	"github.com/IkbolNazarov/gps-tracker/pkg/utils"
	"github.com/google/uuid"
)

type SystemUsers interface {
	Create(ctx context.Context, sysUser *model.CreateSystemUser, user_id uuid.UUID, ip, role string) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.SystemUser, error)
	Delete(ctx context.Context, id, user_id uuid.UUID, role, ip, tt, action string) error
	Update(ctx context.Context, u *model.UpdateSystemUser, id, user_id uuid.UUID, ip, role string) error
	List(ctx context.Context) ([]model.SystemUser, error)
	GetByLogin(ctx context.Context, login string) (*model.SystemUser, error)
	Deactivate(ctx context.Context, id, dId uuid.UUID, dReason, role, ip, tt, action string) error
}

type systemUserService struct {
	r     repository.SystemUsersRepository
	audit AuditLogs
}

func NewSystemUsersService(r repository.SystemUsersRepository, audit AuditLogs) SystemUsers {
	return &systemUserService{r: r, audit: audit}
}

func (s *systemUserService) Create(ctx context.Context, sysUser *model.CreateSystemUser, user_id uuid.UUID, ip, role string) error {

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

	if err := s.r.Create(ctx, u); err != nil {
		return err
	}

	if user_id != uuid.Nil {
		system := constants.UserTypeSystem

		if err := s.audit.Create(ctx, &model.AuditLogs{
			ActorID:     user_id,
			ActorType:   role,
			TargetTable: &system,
			Action:      "CreateSystemUser",
			IPAddress:   &ip,
			CreatedAt:   time.Now().In(timezone.Dushanbe),
		}); err != nil {
			return err
		}
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

func (s *systemUserService) Update(ctx context.Context, u *model.UpdateSystemUser, id, user_id uuid.UUID, ip, role string) error {
	if err := s.r.Update(ctx, u, id); err != nil {
		return err
	}

	system := constants.UserTypeSystem

	if err := s.audit.Create(ctx, &model.AuditLogs{
		ActorID:     user_id,
		ActorType:   role,
		Action:      "UpdateSystemUser",
		TargetTable: &system,
		TargetID:    &id,
		IPAddress:   &ip,
		CreatedAt:   time.Now().In(timezone.Dushanbe),
	}); err != nil {
		return err
	}

	return nil
}

func (s *systemUserService) Delete(ctx context.Context, id, user_id uuid.UUID, role, ip, tt, action string) error {
	if err := s.r.Delete(ctx, id, user_id); err != nil {
		return err
	}

	if err := s.audit.Create(ctx, &model.AuditLogs{
		ActorID:     user_id,
		ActorType:   role,
		IPAddress:   &ip,
		TargetID:    &id,
		TargetTable: &tt,
		Action:      action,
		CreatedAt:   time.Now().In(timezone.Dushanbe),
	}); err != nil {
		return err
	}

	return nil
}

func (s *systemUserService) Deactivate(ctx context.Context, id, dId uuid.UUID, dReason, role, ip, tt, action string) error {
	if err := s.r.Deactivate(ctx, id, dId, dReason); err != nil {
		return err
	}

	if err := s.audit.Create(ctx, &model.AuditLogs{
		ActorID:     dId,
		ActorType:   role,
		IPAddress:   &ip,
		TargetID:    &id,
		TargetTable: &tt,
		Action:      action,
		CreatedAt:   time.Now().In(timezone.Dushanbe),
	}); err != nil {
		return err
	}

	return nil
}

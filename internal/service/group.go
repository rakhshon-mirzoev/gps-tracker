package service

import (
	"context"
	"time"

	"github.com/IkbolNazarov/gps-tracker/internal/model"
	"github.com/IkbolNazarov/gps-tracker/internal/repository"
	"github.com/IkbolNazarov/gps-tracker/pkg/timezone"
	"github.com/google/uuid"
)

type Group interface {
	Create(ctx context.Context, group *model.Group, user_id uuid.UUID, ip, gr, role, action string) error
	Delete(ctx context.Context, id uuid.UUID, user_id uuid.UUID, ip, gr, role, action string) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Group, error)
	Update(ctx context.Context, group model.Group, id uuid.UUID, user_id uuid.UUID, ip, gr, role, action string) error
	List(ctx context.Context, groupType string, supervisorID *uuid.UUID, limit, offset int) ([]model.Group, int, error)
	GetGroupIDForSupervisor(ctx context.Context, userID uuid.UUID) (uuid.UUID, error)
	GetHierarchy(ctx context.Context, groupType string) ([]model.Group, error)
	GetDescendantIDs(ctx context.Context, id uuid.UUID) ([]uuid.UUID, error)
	GetDescendants(ctx context.Context, id uuid.UUID) ([]model.Group, error)
}

type GroupService struct {
	r     repository.GroupsRepository
	audit AuditLogs
}

func NewGroupService(r repository.GroupsRepository, audit AuditLogs) Group {
	return &GroupService{r: r, audit: audit}
}

func (g *GroupService) Create(ctx context.Context, group *model.Group, user_id uuid.UUID, ip, gr, role, action string) error {
	if err := g.r.Create(ctx, group); err != nil {
		return err
	}

	if err := g.audit.Create(ctx, &model.AuditLogs{
		ActorID:     user_id,
		ActorType:   role,
		Action:      action,
		TargetTable: &gr,
		IPAddress:   &ip,
		CreatedAt:   time.Now().In(timezone.Dushanbe),
	}); err != nil {
		return err
	}

	return nil
}

func (g *GroupService) List(ctx context.Context, groupType string, supervisorID *uuid.UUID, limit, offset int) ([]model.Group, int, error) {
	return g.r.List(ctx, groupType, supervisorID, limit, offset)
}

func (g *GroupService) Delete(ctx context.Context, id uuid.UUID, user_id uuid.UUID, ip, gr, role, action string) error {

	if err := g.r.Delete(ctx, id); err != nil {
		return err
	}

	if err := g.audit.Create(ctx, &model.AuditLogs{
		ActorID:     user_id,
		ActorType:   role,
		Action:      action,
		TargetTable: &gr,
		TargetID:    &id,
		IPAddress:   &ip,
		CreatedAt:   time.Now().In(timezone.Dushanbe),
	}); err != nil {
		return err
	}

	return nil
}

func (g *GroupService) GetByID(ctx context.Context, id uuid.UUID) (*model.Group, error) {
	return g.r.GetByID(ctx, id)
}

func (g *GroupService) GetGroupIDForSupervisor(ctx context.Context, userID uuid.UUID) (uuid.UUID, error) {
	return g.r.GetGroupIDForSupervisor(ctx, userID)
}

func (g *GroupService) Update(ctx context.Context, group model.Group, id uuid.UUID, user_id uuid.UUID, ip, gr, role, action string) error {
	if err := g.r.Update(ctx, group, id); err != nil {
		return err
	}

	if err := g.audit.Create(ctx, &model.AuditLogs{
		ActorID:     user_id,
		ActorType:   role,
		Action:      action,
		TargetTable: &gr,
		TargetID:    &id,
		IPAddress:   &ip,
		CreatedAt:   time.Now().In(timezone.Dushanbe),
	}); err != nil {
		return err
	}

	return nil
}

func (g *GroupService) GetHierarchy(ctx context.Context, groupType string) ([]model.Group, error) {
	return g.r.GetHierarchy(ctx, groupType)
}

func (g *GroupService) GetDescendantIDs(ctx context.Context, id uuid.UUID) ([]uuid.UUID, error) {
	return g.r.GetDescendantIDs(ctx, id)
}

func (g *GroupService) GetDescendants(ctx context.Context, id uuid.UUID) ([]model.Group, error) {
	return g.r.GetDescendants(ctx, id)
}

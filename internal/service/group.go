package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/model"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/repository"
)

type Group interface {
	Create(ctx context.Context, group *model.Group, user_id uuid.UUID, ip, gr, role, action string) error
	Delete(ctx context.Context, id uuid.UUID, user_id uuid.UUID, ip, gr, role, action string) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Group, error)
	Update(ctx context.Context, group model.Group, id uuid.UUID, user_id uuid.UUID, ip, gr, role, action string) error
	List(ctx context.Context, supervisorID *uuid.UUID, limit, offset int) ([]model.Group, int, error)
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
	}); err != nil {
		return err
	}

	return nil
}

func (g *GroupService) List(ctx context.Context, id *uuid.UUID, limit, offset int) ([]model.Group, int, error) {
	return g.r.List(ctx, id, limit, offset)
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
	}); err != nil {
		return err
	}

	return nil
}

func (g *GroupService) GetByID(ctx context.Context, id uuid.UUID) (*model.Group, error) {
	return g.r.GetByID(ctx, id)
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
	}); err != nil {
		return err
	}

	return nil
}

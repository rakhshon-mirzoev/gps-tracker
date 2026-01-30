package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/model"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/repository"
)

type Group interface {
	Create(ctx context.Context, group *model.Group) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Group, error)
	Update(ctx context.Context, gr model.Group, id uuid.UUID) error
	List(ctx context.Context, supervisorID *uuid.UUID) ([]model.Group, error)
}

type GroupService struct {
	r repository.GroupsRepository
}

func NewGroupService(r repository.GroupsRepository) Group {
	return &GroupService{r: r}
}

func (g *GroupService) Create(ctx context.Context, group *model.Group) error {
	return g.r.Create(ctx, group)
}

func (g *GroupService) List(ctx context.Context, id *uuid.UUID) ([]model.Group, error) {
	return g.r.List(ctx, id)
}

func (g *GroupService) Delete(ctx context.Context, id uuid.UUID) error {
	return g.r.Delete(ctx, id)
}

func (g *GroupService) GetByID(ctx context.Context, id uuid.UUID) (*model.Group, error) {
	return g.r.GetByID(ctx, id)
}

func (g *GroupService) Update(ctx context.Context, gr model.Group, id uuid.UUID) error {
	return g.r.Update(ctx, gr, id)
}

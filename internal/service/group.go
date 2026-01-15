package service

import (
	"context"

	"github.com/rakhshon-mirzoev/gps-tracker/internal/model"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/repository"
)

type Group interface {
	CreateGroup(ctx context.Context, group *model.Group) error
}

type GroupService struct {
	r repository.GroupsRepository
}

func InitGroupService(r repository.GroupsRepository) Group {
	return &GroupService{r: r}
}

func (g *GroupService) CreateGroup(ctx context.Context, group *model.Group) error {
	err := g.r.CreateGroup(ctx, group)
	if err != nil {
		return err
	}
	return nil
}

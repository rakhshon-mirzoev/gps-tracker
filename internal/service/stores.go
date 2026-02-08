package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/model"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/repository"
)

type Stores interface {
	Create(ctx context.Context, s *model.CreateStore, user_id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Store, error)
	List(ctx context.Context, filter model.StoreListFilter) ([]model.Store, int, error)
	Update(ctx context.Context, s model.UpdateStore, id uuid.UUID, user_id uuid.UUID, ip, role, st, action string) error
	Delete(ctx context.Context, id uuid.UUID, user_id uuid.UUID, ip, role, st, action string) error
}

type storesService struct {
	r     repository.StoreRepository
	audit AuditLogs
}

func NewStoresService(repo repository.StoreRepository, audit AuditLogs) Stores {
	return &storesService{r: repo, audit: audit}
}

func (s *storesService) Create(ctx context.Context, store *model.CreateStore, user_id uuid.UUID) error {
	staffID, err := uuid.Parse(store.CreatedByStaffID)
	if err != nil {
		return err
	}

	st := &model.Store{
		Name:             store.Name,
		Address:          store.Address,
		Latitude:         store.Latitude,
		Longitude:        store.Longitude,
		OwnerFullName:    store.OwnerFullName,
		OwnerINN:         store.OwnerINN,
		OwnerPassport:    store.OwnerPassport,
		CreatedByStaffID: staffID,
	}

	return s.r.Create(ctx, st)
}

func (s *storesService) GetByID(ctx context.Context, id uuid.UUID) (*model.Store, error) {
	if id == uuid.Nil {
		return nil, errors.New("id is required")
	}

	return s.r.GetByID(ctx, id)
}

func (s *storesService) List(ctx context.Context, filter model.StoreListFilter) ([]model.Store, int, error) {
	return s.r.List(ctx, filter)
}

func (s *storesService) Update(ctx context.Context, store model.UpdateStore, id uuid.UUID, user_id uuid.UUID, ip, role, st, action string) error {
	if err := s.r.Update(ctx, store, id); err != nil {
		return err
	}

	if err := s.audit.Create(ctx, &model.AuditLogs{
		ActorID:     user_id,
		ActorType:   role,
		IPAddress:   &ip,
		TargetID:    &id,
		TargetTable: &st,
		Action:      action,
	}); err != nil {
		return err
	}

	return nil
}

func (s *storesService) Delete(ctx context.Context, id uuid.UUID, user_id uuid.UUID, ip, role, st, action string) error {

	if err := s.r.Delete(ctx, id); err != nil {
		return err
	}

	if err := s.audit.Create(ctx, &model.AuditLogs{
		ActorID:     user_id,
		ActorType:   role,
		IPAddress:   &ip,
		TargetID:    &id,
		TargetTable: &st,
		Action:      action,
	}); err != nil {
		return err
	}

	return nil
}

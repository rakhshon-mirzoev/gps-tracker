package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/model"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/repository"
	errs "github.com/rakhshon-mirzoev/gps-tracker/pkg/error"
	minioservice "github.com/rakhshon-mirzoev/gps-tracker/pkg/storage/service"
	"github.com/rakhshon-mirzoev/gps-tracker/pkg/utils"
)

type Stores interface {
	Create(ctx context.Context, s *model.CreateStore, user_id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Store, error)
	List(ctx context.Context, isActive, status, sortBy, sortOrder string, lim, off int) ([]model.Store, int, error)
	Update(ctx context.Context, s model.UpdateStore, id uuid.UUID, user_id uuid.UUID, ip, role, st, action string) error
	Delete(ctx context.Context, id uuid.UUID, user_id uuid.UUID, ip, role, st, action string) error
	IfExists(ctx context.Context, name string) (bool, error)
}

type storesService struct {
	r     repository.StoreRepository
	audit AuditLogs
	minio minioservice.Minio
}

func NewStoresService(repo repository.StoreRepository, audit AuditLogs, minio minioservice.Minio) Stores {
	return &storesService{r: repo, audit: audit, minio: minio}
}

func IsTmpStoreKey(key string) bool {
	return strings.HasPrefix(key, "stores/tmp/")
}

func (s *storesService) Create(ctx context.Context, store *model.CreateStore, user_id uuid.UUID) error {

	st := &model.Store{
		Name:             store.Name,
		Address:          store.Address,
		Latitude:         store.Latitude,
		Longitude:        store.Longitude,
		OwnerFullName:    store.OwnerFullName,
		OwnerINN:         store.OwnerINN,
		OwnerPassport:    store.OwnerPassport,
		OwnerPhone:       store.OwnerPhone,
		AdditionalPhone:  store.AdditionalPhone,
		Comment:          store.Comment,
		CreatedByStaffID: user_id,
	}

	exists, err := s.IfExists(ctx, st.Name)
	if err != nil {
		return err
	}
	if exists {
		return errs.ErrorAlreadyExists
	}

	if err := s.r.Create(ctx, st); err != nil {
		return err
	}

	var newInterior, newExterior *string

	if store.InteriorPhoto != nil {
		if !strings.HasPrefix(*store.InteriorPhoto, "stores/tmp/") {
			return fmt.Errorf("invalid interior photo key")
		}

		newKey := strings.Replace(*store.InteriorPhoto, "stores/tmp/", "stores/proposals/", 1)

		if err := s.minio.MoveObject(ctx, *store.InteriorPhoto, newKey); err != nil {
			return err
		}

		newInterior = &newKey
	}

	if store.ExteriorPhoto != nil {
		if !strings.HasPrefix(*store.ExteriorPhoto, "stores/tmp/") {
			return fmt.Errorf("invalid exterior photo key")
		}

		newKey := strings.Replace(*store.ExteriorPhoto, "stores/tmp/", "stores/proposals/", 1)

		if err := s.minio.MoveObject(ctx, *store.ExteriorPhoto, newKey); err != nil {
			return err
		}

		newExterior = &newKey
	}

	if newInterior != nil || newExterior != nil {
		err := s.r.Update(ctx, model.UpdateStore{
			ExteriorPhoto: newExterior,
			InteriorPhoto: newInterior,
		}, st.ID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *storesService) GetByID(ctx context.Context, id uuid.UUID) (*model.Store, error) {
	if id == uuid.Nil {
		return nil, errors.New("id is required")
	}

	store, err := s.r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if store.InteriorPhoto != nil {
		url, err := s.minio.PresignedDownload(ctx, *store.InteriorPhoto, time.Minute*15)
		if err != nil {
			return nil, err
		}

		store.InteriorPhotoURL = &url
	}

	if store.ExteriorPhoto != nil {
		url, err := s.minio.PresignedDownload(ctx, *store.ExteriorPhoto, time.Minute*15)
		if err != nil {
			return nil, err
		}

		store.ExteriorPhotoURL = &url
	}

	return store, nil
}

func (s *storesService) List(ctx context.Context, isActive, status, sortBy, sortOrder string, lim, off int) ([]model.Store, int, error) {

	if err := utils.ValidateStoreStatus(status); err != nil {
		return nil, 0, err
	}

	stores, total, err := s.r.List(ctx, isActive, status, sortBy, sortOrder, lim, off)
	if err != nil {
		return nil, 0, err
	}

	for i, v := range stores {
		if v.InteriorPhoto != nil {
			url, err := s.minio.PresignedDownload(ctx, *v.InteriorPhoto, time.Minute*15)
			if err != nil {
				return nil, 0, err
			}

			v.InteriorPhotoURL = &url
		}

		if v.ExteriorPhoto != nil {
			url, err := s.minio.PresignedDownload(ctx, *v.ExteriorPhoto, time.Minute*15)
			if err != nil {
				return nil, 0, err
			}

			v.ExteriorPhotoURL = &url
		}

		stores[i] = v
	}

	return stores, total, nil
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

func (s *storesService) IfExists(ctx context.Context, name string) (bool, error) {
	return s.r.IfExists(ctx, name)
}

package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/model"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/repository"
	errs "github.com/rakhshon-mirzoev/gps-tracker/pkg/error"
	minioservice "github.com/rakhshon-mirzoev/gps-tracker/pkg/storage/service"
	"github.com/rakhshon-mirzoev/gps-tracker/pkg/utils"
	"golang.org/x/sync/errgroup"
)

type Stores interface {
	Create(ctx context.Context, s model.CreateStore, user_id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Store, error)
	GetByIDPhoto(ctx context.Context, id uuid.UUID) (*model.StoreWithPhoto, error)
	List(ctx context.Context, isActive, status, sortBy, sortOrder string, lim, off int) ([]model.StoreWithPhoto, int, error)
	Update(ctx context.Context, s model.UpdateStore, id uuid.UUID, user_id uuid.UUID, ip, role, st, action string) error
	Delete(ctx context.Context, id uuid.UUID, user_id uuid.UUID, ip, role, st, action string) error
	IfExists(ctx context.Context, name string) (bool, error)
	IfExistsTx(ctx context.Context, tx pgx.Tx, name string) (bool, error)
}

type storesService struct {
	r      repository.StoreRepository
	audit  AuditLogs
	minio  minioservice.Minio
	photos StorePhotos
}

func NewStoresService(repo repository.StoreRepository, audit AuditLogs, minio minioservice.Minio, photos StorePhotos) Stores {
	return &storesService{r: repo, audit: audit, minio: minio, photos: photos}
}

func (s *storesService) Create(ctx context.Context, store model.CreateStore, userID uuid.UUID) error {
	if err := utils.ValidateAllPhotos(store.InteriorPhotos, store.ExteriorPhotos); err != nil {
		return err
	}

	if len(store.ExteriorPhotos) > 10 || len(store.InteriorPhotos) > 10 {
		return errs.ErrorTooManyFiles
	}

	tx, err := s.r.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("start tx: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()

	st := &model.Store{
		Name:             store.Name,
		Address:          store.Address,
		Latitude:         store.Latitude,
		Longitude:        store.Longitude,
		OwnerFullName:    store.OwnerFullName,
		OwnerINN:         store.OwnerINN,
		OwnerPassport:    store.OwnerPassport,
		OwnerPhone:       &store.OwnerPhone,
		AdditionalPhone:  store.AdditionalPhone,
		Comment:          store.Comment,
		CreatedByStaffID: userID,
	}

	exists, err := s.IfExistsTx(ctx, tx, st.Name)
	if err != nil {
		return err
	}
	if exists {
		return errs.ErrorAlreadyExists
	}

	if err := s.r.CreateTx(ctx, tx, st); err != nil {
		return err
	}

	type photoItem struct {
		Key  string
		Type string
	}
	var allPhotos []photoItem
	for _, k := range store.InteriorPhotos {
		allPhotos = append(allPhotos, photoItem{Key: k, Type: "interior"})
	}
	for _, k := range store.ExteriorPhotos {
		allPhotos = append(allPhotos, photoItem{Key: k, Type: "exterior"})
	}

	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(5)

	for _, p := range allPhotos {
		p := p
		g.Go(func() error {
			if !strings.HasPrefix(p.Key, "stores/tmp/") {
				return fmt.Errorf("invalid photo key: %s", p.Key)
			}
			newKey := strings.Replace(p.Key, "stores/tmp/", "stores/proposals/", 1)

			if err := s.minio.MoveObject(ctx, p.Key, newKey); err != nil {
				return err
			}

			if err := s.photos.CreateTx(ctx, tx, &model.StorePhotos{
				StoreID:   st.ID,
				ObjectKey: newKey,
				Type:      p.Type,
			}); err != nil {
				return err
			}
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}

func (s *storesService) GetByID(ctx context.Context, id uuid.UUID) (*model.Store, error) {

	return s.r.GetByID(ctx, id)
}

func (s *storesService) GetByIDPhoto(ctx context.Context, id uuid.UUID) (*model.StoreWithPhoto, error) {

	store, err := s.r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	photos, err := s.photos.List(ctx, id)
	if err != nil {
		return nil, err
	}

	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(5)

	for i := range photos {
		i := i

		g.Go(func() error {
			url, err := s.minio.PresignedDownload(ctx, photos[i].ObjectKey, 15*time.Minute)
			if err != nil {
				return err
			}

			photos[i].URL = url
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return &model.StoreWithPhoto{
		Store:       *store,
		StorePhotos: photos,
	}, nil
}

func (s *storesService) List(ctx context.Context, isActive, status, sortBy, sortOrder string, lim, off int) ([]model.StoreWithPhoto, int, error) {

	if err := utils.ValidateStoreStatus(status); err != nil {
		return nil, 0, err
	}

	stores, total, err := s.r.List(ctx, isActive, status, sortBy, sortOrder, lim, off)
	if err != nil {
		return nil, 0, err
	}

	result := make([]model.StoreWithPhoto, 0, len(stores))

	for _, store := range stores {

		photos, err := s.photos.List(ctx, store.ID)
		if err != nil {
			return nil, 0, err
		}

		g, ctx := errgroup.WithContext(ctx)
		g.SetLimit(5)

		for i := range photos {
			i := i

			g.Go(func() error {
				url, err := s.minio.PresignedDownload(ctx, photos[i].ObjectKey, 15*time.Minute)
				if err != nil {
					return err
				}
				photos[i].URL = url
				return nil
			})
		}

		if err := g.Wait(); err != nil {
			return nil, 0, err
		}

		result = append(result, model.StoreWithPhoto{
			Store:       store,
			StorePhotos: photos,
		})
	}

	return result, total, nil
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

func (s *storesService) IfExistsTx(ctx context.Context, tx pgx.Tx, name string) (bool, error) {
	return s.r.IfExistsTx(ctx, tx, name)
}

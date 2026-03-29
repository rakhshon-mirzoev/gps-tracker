package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/IkbolNazarov/gps-tracker/internal/logger"
	"github.com/IkbolNazarov/gps-tracker/internal/model"
	"github.com/IkbolNazarov/gps-tracker/internal/repository"
	"github.com/IkbolNazarov/gps-tracker/pkg/fcm"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

type Notification interface {
	CreateAndPushToUser(ctx context.Context, userID uuid.UUID, title, body, ntype string, meta map[string]any, token *string) error
	CreateAndPushToGroup(ctx context.Context, groupID uuid.UUID, title, body, ntype string, meta map[string]any) error
	ListByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]model.Notification, int, error)
	MarkRead(ctx context.Context, notificationID uuid.UUID, userID uuid.UUID) error
	CreateTx(ctx context.Context, tx pgx.Tx, n *model.Notification) error
}

type notificationService struct {
	r         repository.NotificationsRepository
	staffRepo repository.StaffUsersRepository
	sender    fcm.PushSender
}

func NewNotificationService(r repository.NotificationsRepository, staff repository.StaffUsersRepository, sender fcm.PushSender) Notification {
	return &notificationService{r: r, staffRepo: staff, sender: sender}
}

func (s *notificationService) CreateTx(ctx context.Context, tx pgx.Tx, n *model.Notification) error {
	return s.r.CreateTx(ctx, tx, n)
}

func (s *notificationService) CreateAndPushToUser(ctx context.Context, userID uuid.UUID, title, body, ntype string, meta map[string]any, token *string) error {
	var metaRaw []byte
	if meta != nil {
		b, _ := json.Marshal(meta)
		metaRaw = b
	}

	n := &model.Notification{
		UserID:   &userID,
		Type:     ntype,
		Title:    title,
		Body:     body,
		Meta:     metaRaw,
		IsRead:   false,
		PushSent: false,
	}

	if err := s.r.Create(ctx, n); err != nil {
		return fmt.Errorf("create notification: %w", err)
	}

	if token == nil {
		u, err := s.staffRepo.GetByID(ctx, userID)
		if err != nil {
			return err
		}

		if u != nil && u.FcmToken != nil {
			token = u.FcmToken
		}

	}

	if token != nil {
		if err := s.sender.SendToToken(ctx, *token, title, body, map[string]string{"type": ntype}); err == nil {
			if err := s.r.UpdatePushSent(ctx, n.ID, true); err != nil {
				logger.Log.Error(fmt.Sprintf("failed to mark push_sent for %s", n.ID), zap.Error(err))
			}
		} else {
			logger.Log.Error(fmt.Sprintf("failed to send push to user %s", userID), zap.Error(err))
		}
	}

	return nil
}

func (s *notificationService) CreateAndPushToGroup(ctx context.Context, groupID uuid.UUID, title, body, ntype string, meta map[string]any) error {
	users, err := s.staffRepo.GetByGroupID(ctx, groupID)
	if err != nil {
		return err
	}

	for _, u := range users {
		var token *string
		if u.FcmToken != nil {
			token = u.FcmToken
		}
		err := s.CreateAndPushToUser(ctx, u.ID, title, body, ntype, meta, token)
		if err != nil {
			logger.Log.Error(fmt.Sprintf("failed to push to user: %s. body: %s", u.ID, body))
		}
	}

	return nil
}

func (s *notificationService) ListByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]model.Notification, int, error) {
	return s.r.ListByUser(ctx, userID, limit, offset)
}

func (s *notificationService) MarkRead(ctx context.Context, notificationID uuid.UUID, userID uuid.UUID) error {
	return s.r.MarkRead(ctx, notificationID, userID)
}

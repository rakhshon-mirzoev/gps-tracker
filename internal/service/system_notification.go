package service

import (
	"context"
	"fmt"

	"github.com/IkbolNazarov/gps-tracker/internal/events"
	"github.com/IkbolNazarov/gps-tracker/internal/logger"
	"github.com/IkbolNazarov/gps-tracker/internal/model"
	"github.com/IkbolNazarov/gps-tracker/internal/repository"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// SystemNotification — сервис уведомлений для admin/supervisor.
type SystemNotification interface {
	// Dispatch вызывается Kafka-консьюмером: сохраняет уведомление в БД
	// для всех подходящих получателей (admin — все, supervisor — своей группы).
	Dispatch(ctx context.Context, event events.NotificationEvent) ([]model.SystemNotification, error)
	ListByRecipient(ctx context.Context, recipientID uuid.UUID, limit, offset int) ([]model.SystemNotification, int, error)
	MarkRead(ctx context.Context, notificationID uuid.UUID, recipientID uuid.UUID) error
}

type systemNotificationService struct {
	repo     repository.SystemNotificationsRepository
	sysUsers repository.SystemUsersRepository
}

func NewSystemNotificationService(
	repo repository.SystemNotificationsRepository,
	sysUsers repository.SystemUsersRepository,
) SystemNotification {
	return &systemNotificationService{repo: repo, sysUsers: sysUsers}
}

func (s *systemNotificationService) Dispatch(ctx context.Context, event events.NotificationEvent) ([]model.SystemNotification, error) {
	// 1. Все активные admin-ы получают все уведомления
	admins, err := s.sysUsers.GetAllAdmins(ctx)
	if err != nil {
		return nil, fmt.Errorf("dispatch: get admins: %w", err)
	}

	// 2. Supervisor группы тоже получает уведомление
	supervisor, err := s.sysUsers.GetSupervisorByGroupID(ctx, event.GroupID)
	if err != nil {
		return nil, fmt.Errorf("dispatch: get supervisor: %w", err)
	}

	// Собираем уникальных получателей
	seen := make(map[uuid.UUID]bool)
	var recipients []model.SystemUser
	for _, a := range admins {
		if !seen[a.ID] {
			seen[a.ID] = true
			recipients = append(recipients, a)
		}
	}
	if supervisor != nil && !seen[supervisor.ID] {
		recipients = append(recipients, *supervisor)
	}

	groupID := event.GroupID
	staffID := event.StaffID

	var created []model.SystemNotification
	for _, r := range recipients {
		n := &model.SystemNotification{
			RecipientID: r.ID,
			StaffID:     &staffID,
			GroupID:     &groupID,
			Type:        event.Type,
			Title:       event.Title,
			Body:        event.Body,
		}
		if err := s.repo.Create(ctx, n); err != nil {
			logger.Log.Error("dispatch: create system notification",
				zap.Stringer("recipient_id", r.ID),
				zap.Error(err),
			)
			continue
		}
		created = append(created, *n)
	}

	return created, nil
}

func (s *systemNotificationService) ListByRecipient(ctx context.Context, recipientID uuid.UUID, limit, offset int) ([]model.SystemNotification, int, error) {
	return s.repo.ListByRecipient(ctx, recipientID, limit, offset)
}

func (s *systemNotificationService) MarkRead(ctx context.Context, notificationID uuid.UUID, recipientID uuid.UUID) error {
	return s.repo.MarkRead(ctx, notificationID, recipientID)
}

package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/IkbolNazarov/gps-tracker/internal/model"
	"github.com/IkbolNazarov/gps-tracker/pkg/db"
	errs "github.com/IkbolNazarov/gps-tracker/pkg/error"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type SystemNotificationsRepository interface {
	Create(ctx context.Context, n *model.SystemNotification) error
	ListByRecipient(ctx context.Context, recipientID uuid.UUID, limit, offset int) ([]model.SystemNotification, int, error)
	MarkRead(ctx context.Context, notificationID uuid.UUID, recipientID uuid.UUID) error
}

type systemNotificationsRepo struct {
	DB *db.DB
}

func NewSystemNotificationsRepository(db *db.DB) SystemNotificationsRepository {
	return &systemNotificationsRepo{DB: db}
}

func (r *systemNotificationsRepo) Create(ctx context.Context, n *model.SystemNotification) error {
	sql := `
		INSERT INTO system_notifications (
			recipient_id, staff_id, group_id, type, title, body, is_read
		) VALUES ($1, $2, $3, $4, $5, $6, FALSE)
		RETURNING id, created_at
	`
	return r.DB.Pool.QueryRow(ctx, sql,
		n.RecipientID,
		n.StaffID,
		n.GroupID,
		n.Type,
		n.Title,
		n.Body,
	).Scan(&n.ID, &n.CreatedAt)
}

func (r *systemNotificationsRepo) ListByRecipient(ctx context.Context, recipientID uuid.UUID, limit, offset int) ([]model.SystemNotification, int, error) {
	var total int
	if err := r.DB.Pool.QueryRow(ctx,
		`SELECT count(*) FROM system_notifications WHERE recipient_id = $1`,
		recipientID,
	).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count system notifications: %w", err)
	}

	sql := `
		SELECT id, recipient_id, staff_id, group_id, type, title, body,
		       is_read, read_at, created_at
		FROM system_notifications
		WHERE recipient_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.DB.Pool.Query(ctx, sql, recipientID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("list system notifications: %w", err)
	}
	defer rows.Close()

	var result []model.SystemNotification
	for rows.Next() {
		var n model.SystemNotification
		if err := rows.Scan(
			&n.ID, &n.RecipientID, &n.StaffID, &n.GroupID,
			&n.Type, &n.Title, &n.Body,
			&n.IsRead, &n.ReadAt, &n.CreatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("scan system notification: %w", err)
		}
		result = append(result, n)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("rows system notifications: %w", err)
	}

	return result, total, nil
}

func (r *systemNotificationsRepo) MarkRead(ctx context.Context, notificationID uuid.UUID, recipientID uuid.UUID) error {
	sql := `
		UPDATE system_notifications
		SET is_read = TRUE, read_at = now()
		WHERE id = $1 AND recipient_id = $2
	`
	ct, err := r.DB.Pool.Exec(ctx, sql, notificationID, recipientID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errs.ErrorNotFound
		}
		return fmt.Errorf("mark read system notification: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return errs.ErrorNotFound
	}
	return nil
}

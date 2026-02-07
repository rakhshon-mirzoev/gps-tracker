package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/model"
	"github.com/rakhshon-mirzoev/gps-tracker/pkg/db"
	errs "github.com/rakhshon-mirzoev/gps-tracker/pkg/error"
)

type NotificationsRepository interface {
	Create(ctx context.Context, n *model.Notification) error
	ListByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]model.Notification, int, error)
	MarkRead(ctx context.Context, notificationID uuid.UUID, userID uuid.UUID) error
	UpdatePushSent(ctx context.Context, notificationID uuid.UUID, sent bool) error
}

type notificationsRepo struct {
	DB *db.DB
}

func NewNotificationsRepository(db *db.DB) NotificationsRepository {
	return &notificationsRepo{DB: db}
}

func (r *notificationsRepo) Create(ctx context.Context, n *model.Notification) error {
	sql := `
		INSERT INTO notifications (
			user_id, 
			group_id, 
			type, 
			title, 
			body, 
			meta, 
			is_read, 
			push_sent
		) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8) 
		RETURNING id, created_at
	`

	err := r.DB.Pool.QueryRow(ctx, sql, n.UserID, n.GroupID, n.Type, n.Title, n.Body, n.Meta, n.IsRead, n.PushSent).Scan(&n.ID, &n.CreatedAt)
	if err != nil {
		return fmt.Errorf("create notification: %w", err)
	}
	return nil
}

func (r *notificationsRepo) ListByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]model.Notification, int, error) {
	total := 0
	sql_count := `SELECT count(*) FROM notifications WHERE user_id = $1`
	row := r.DB.Pool.QueryRow(ctx, sql_count, userID)
	if err := row.Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count notifications: %w", err)
	}

	sql := `SELECT id, user_id, group_id, type, title, body, meta, 
				is_read, push_sent, created_at FROM notifications 
			WHERE user_id = $1 
			ORDER BY created_at DESC LIMIT $2 OFFSET $3`

	rows, err := r.DB.Pool.Query(ctx, sql, userID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("list notifications: %w", err)
	}
	defer rows.Close()

	var res []model.Notification
	for rows.Next() {
		var n model.Notification
		if err := rows.Scan(&n.ID, &n.UserID, &n.GroupID, &n.Type, &n.Title, &n.Body, &n.Meta, &n.IsRead, &n.PushSent, &n.CreatedAt); err != nil {
			return nil, 0, fmt.Errorf("scan notification: %w", err)
		}
		res = append(res, n)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("rows notifications: %w", err)
	}

	return res, total, nil
}

func (r *notificationsRepo) MarkRead(ctx context.Context, notificationID uuid.UUID, userID uuid.UUID) error {
	sql := `UPDATE notifications SET is_read = true WHERE id = $1 AND user_id = $2`
	ct, err := r.DB.Pool.Exec(ctx, sql, notificationID, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errs.ErrorNotFound
		}
		return fmt.Errorf("mark read: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return errs.ErrorNotFound
	}
	return nil
}

func (r *notificationsRepo) UpdatePushSent(ctx context.Context, notificationID uuid.UUID, sent bool) error {
	sql := `UPDATE notifications SET push_sent = $1 WHERE id = $2`
	ct, err := r.DB.Pool.Exec(ctx, sql, sent, notificationID)
	if err != nil {
		return fmt.Errorf("update push_sent: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return errs.ErrorNotFound
	}
	return nil
}

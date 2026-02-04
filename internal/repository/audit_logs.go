package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/model"
	"github.com/rakhshon-mirzoev/gps-tracker/pkg/db"
)

type AuditLogsRepository interface {
	Create(ctx context.Context, al *model.AuditLogs) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.AuditLogs, error)
	List(ctx context.Context, actorID uuid.UUID) ([]model.AuditLogs, error)
	Update(ctx context.Context, id uuid.UUID, dto model.UpdateAuditLogs) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type auditLogs struct {
	DB *db.DB
}

func InitAuditLogsRepository(db *db.DB) AuditLogsRepository {
	return &auditLogs{DB: db}
}

func (r *auditLogs) Create(ctx context.Context, al *model.AuditLogs) error {
	sql := `
        INSERT INTO audit_logs (
            actor_id,
            actor_type,
            action,
            target_table,
            target_id,
            old_value,
            new_value,
            ip_address
        ) VALUES (
            $1, $2, $3, $4, $5, $6, $7, $8
        ) RETURNING id, created_at;
    `

	err := r.DB.Pool.QueryRow(ctx, sql,
		al.ActorID, al.ActorType, al.Action, al.TargetTable, al.TargetID,
		al.OldValue, al.NewValue, al.IPAddress,
	).Scan(&al.ID, &al.CreatedAt)
	if err != nil {
		return fmt.Errorf("insert audit logs: %w", err)
	}

	return nil
}

func (r *auditLogs) GetByID(ctx context.Context, id uuid.UUID) (*model.AuditLogs, error) {
	sql := `
        SELECT id, actor_id, actor_type, action, target_table, target_id,
               old_value, new_value, ip_address, created_at
        FROM audit_logs WHERE id = $1
    `

	var al model.AuditLogs
	err := r.DB.Pool.QueryRow(ctx, sql, id).Scan(
		&al.ID, &al.ActorID, &al.ActorType, &al.Action, &al.TargetTable, &al.TargetID,
		&al.OldValue, &al.NewValue, &al.IPAddress, &al.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("get by id audit logs: %w", err)
	}

	return &al, nil
}

func (r *auditLogs) List(ctx context.Context, actorID uuid.UUID) ([]model.AuditLogs, error) {
	sql := `
        SELECT id, actor_id, actor_type, action, target_table, target_id,
               old_value, new_value, ip_address, created_at
        FROM audit_logs WHERE actor_id = $1
    `

	rows, err := r.DB.Pool.Query(ctx, sql, actorID)
	if err != nil {
		return nil, fmt.Errorf("list audit logs: %w", err)
	}

	var logs = make([]model.AuditLogs, 0)

	for rows.Next() {
		var al model.AuditLogs

		if err := rows.Scan(
			&al.ID, &al.ActorID, &al.ActorType, &al.Action, &al.TargetTable, &al.TargetID,
			&al.OldValue, &al.NewValue, &al.IPAddress, &al.CreatedAt,
		); err != nil {
			return nil, err
		}

		logs = append(logs, al)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows audit logs: %w", err)
	}

	return logs, nil
}

func (r *auditLogs) Update(ctx context.Context, id uuid.UUID, dto model.UpdateAuditLogs) error {
	sql := `UPDATE audit_logs SET `
	var args []any
	n := 1

	if dto.ActorID != nil {
		sql += fmt.Sprintf("actor_id = $%d, ", n)
		args = append(args, *dto.ActorID)
		n++
	}

	if dto.ActorType != nil {
		sql += fmt.Sprintf("actor_type = $%d, ", n)
		args = append(args, *dto.ActorType)
		n++
	}

	if dto.Action != nil {
		sql += fmt.Sprintf("action = $%d, ", n)
		args = append(args, *dto.Action)
		n++
	}

	if dto.TargetTable != nil {
		sql += fmt.Sprintf("target_table = $%d, ", n)
		args = append(args, *dto.TargetTable)
		n++
	}

	if dto.TargetID != nil {
		sql += fmt.Sprintf("target_id = $%d, ", n)
		args = append(args, *dto.TargetID)
		n++
	}

	if dto.OldValue != nil {
		sql += fmt.Sprintf("old_value = $%d, ", n)
		args = append(args, dto.OldValue)
		n++
	}

	if dto.NewValue != nil {
		sql += fmt.Sprintf("new_value = $%d, ", n)
		args = append(args, dto.NewValue)
		n++
	}

	if dto.IPAddress != nil {
		sql += fmt.Sprintf("ip_address = $%d, ", n)
		args = append(args, *dto.IPAddress)
		n++
	}

	if len(args) == 0 {
		return nil
	}

	sql = strings.TrimSuffix(sql, ", ")
	sql += fmt.Sprintf(" WHERE id = $%d", n)
	args = append(args, id)

	_, err := r.DB.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("update audit logs %s: %w", id, err)
	}

	return nil
}

func (r *auditLogs) Delete(ctx context.Context, id uuid.UUID) error {
	sql := `
        DELETE FROM audit_logs WHERE id = $1
    `

	_, err := r.DB.Pool.Exec(ctx, sql, id)
	if err != nil {
		return fmt.Errorf("delete audit logs %s: %w", id, err)
	}

	return nil
}

package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/model"
	"github.com/rakhshon-mirzoev/gps-tracker/pkg/db"
)

type AuditLogsRepository interface {
	Create(ctx context.Context, al *model.AuditLogs) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.AuditLogs, error)
	List(ctx context.Context, actorID uuid.UUID) ([]model.AuditLogs, error)
	Update(ctx context.Context, al *model.AuditLogs) error
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
		return err
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
		return nil, err
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
		return nil, err
	}

	var logs []model.AuditLogs

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
		return nil, err
	}

	return logs, nil
}

func (r *auditLogs) Update(ctx context.Context, al *model.AuditLogs) error {
	sql := `
        UPDATE audit_logs
        SET actor_id = $1, actor_type = $2, action = $3, target_table = $4,
            target_id = $5, old_value = $6, new_value = $7, ip_address = $8
        WHERE id = $9
    `

	_, err := r.DB.Pool.Exec(ctx, sql,
		al.ActorID, al.ActorType, al.Action, al.TargetTable, al.TargetID,
		al.OldValue, al.NewValue, al.IPAddress, al.ID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *auditLogs) Delete(ctx context.Context, id uuid.UUID) error {
	sql := `
        DELETE FROM audit_logs WHERE id = $1
    `

	_, err := r.DB.Pool.Exec(ctx, sql, id)
	if err != nil {
		return err
	}

	return nil
}

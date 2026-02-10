package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/model"
	"github.com/rakhshon-mirzoev/gps-tracker/pkg/db"
	errs "github.com/rakhshon-mirzoev/gps-tracker/pkg/error"
)

type groupsRepository struct {
	DB *db.DB
}

type GroupsRepository interface {
	Create(context.Context, *model.Group) error
	List(ctx context.Context, supervisorID *uuid.UUID) ([]model.Group, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.Group, error)
	Update(ctx context.Context, gr model.Group, id uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID) error
}

func NewGroupsRepository(db *db.DB) GroupsRepository {
	return &groupsRepository{DB: db}
}

func (g *groupsRepository) Create(ctx context.Context, group *model.Group) error {
	sql := `
		INSERT INTO groups
		(
			name,
			supervisor_id
		)
		VALUES ($1, $2)
		RETURNING id;
	`
	err := g.DB.Pool.QueryRow(ctx, sql, group.Name, group.SupervisorID).Scan(&group.ID)
	if err != nil {
		return fmt.Errorf("create group: %w", err)
	}

	return nil
}

func (g *groupsRepository) List(ctx context.Context, supervisorID *uuid.UUID) ([]model.Group, error) {
	query := `
		SELECT 
			g.id, 
			g.name, 
			g.supervisor_id, 
			u.full_name as supervisor_name, 
			g.created_at 
		FROM groups g
		LEFT JOIN system_users u ON g.supervisor_id = u.id`

	var conditions []string
	var args []any

	if supervisorID != nil {
		conditions = append(conditions, fmt.Sprintf("g.supervisor_id = $%d", len(args)+1))
		args = append(args, *supervisorID)
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY g.id"

	rows, err := g.DB.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list groups: %w", err)
	}
	defer rows.Close()

	groups := make([]model.Group, 0)

	for rows.Next() {
		var g model.Group

		err := rows.Scan(
			&g.ID,
			&g.Name,
			&g.SupervisorID,
			&g.SupervisorName,
			&g.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan group: %w", err)
		}

		groups = append(groups, g)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows groups: %w", err)
	}

	return groups, nil
}

func (g *groupsRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Group, error) {
	sql := `
		SELECT 
			g.id,
			g.name,
			g.supervisor_id,
			s.full_name as supervisor_name,
			g.created_at
		FROM groups g
		LEFT JOIN system_users s ON g.supervisor_id = s.id
		WHERE g.id = $1
	`

	var gr model.Group

	err := g.DB.Pool.QueryRow(ctx, sql, id).Scan(
		&gr.ID,
		&gr.Name,
		&gr.SupervisorID,
		&gr.SupervisorName,
		&gr.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.ErrorNotFound
		}
		return nil, err
	}

	return &gr, nil
}

func (g *groupsRepository) Update(ctx context.Context, gr model.Group, id uuid.UUID) error {
	var setClauses []string
	var args []any

	if gr.Name != "" {
		args = append(args, gr.Name)
		setClauses = append(setClauses, fmt.Sprintf("name = $%d", len(args)))
	}

	if gr.SupervisorID != nil {
		args = append(args, gr.SupervisorID)
		setClauses = append(setClauses, fmt.Sprintf("supervisor_id = $%d", len(args)))
	}

	if len(setClauses) == 0 {
		return nil
	}

	args = append(args, id)
	query := fmt.Sprintf("UPDATE groups SET %s WHERE id = $%d", strings.Join(setClauses, ", "), len(args))

	ct, err := g.DB.Pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("update group: %w", err)
	}

	if ct.RowsAffected() == 0 {
		return errs.ErrorNotFound
	}

	return nil
}

func (g *groupsRepository) Delete(ctx context.Context, id uuid.UUID) error {
	sql := `DELETE FROM groups WHERE id = $1`

	ct, err := g.DB.Pool.Exec(ctx, sql, id)
	if err != nil {
		return fmt.Errorf("delete group: %w", err)
	}

	if ct.RowsAffected() == 0 {
		return errs.ErrorNotFound
	}

	return nil
}

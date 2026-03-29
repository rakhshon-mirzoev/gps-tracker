package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/IkbolNazarov/gps-tracker/internal/model"
	"github.com/IkbolNazarov/gps-tracker/pkg/db"
	errs "github.com/IkbolNazarov/gps-tracker/pkg/error"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type groupsRepository struct {
	DB *db.DB
}

type GroupsRepository interface {
	Create(context.Context, *model.Group) error
	List(ctx context.Context, supervisorID *uuid.UUID, limit, offset int) ([]model.Group, int, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.Group, error)
	Update(ctx context.Context, gr model.Group, id uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetGroupIDForSupervisor(ctx context.Context, userID uuid.UUID) (uuid.UUID, error)
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

func (g *groupsRepository) List(ctx context.Context, supervisorID *uuid.UUID, limit, offset int) ([]model.Group, int, error) {
	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	var conditions []string
	var args []any

	conditions = append(conditions, "g.is_deleted = false")

	if supervisorID != nil {
		conditions = append(conditions, fmt.Sprintf("g.supervisor_id = $%d", len(args)+1))
		args = append(args, *supervisorID)
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = " WHERE " + strings.Join(conditions, " AND ")
	}

	var total int
	countSQL := `SELECT COUNT(*) FROM groups g` + whereClause
	if err := g.DB.Pool.QueryRow(ctx, countSQL, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count groups: %w", err)
	}

	query := fmt.Sprintf(`
		SELECT 
			g.id, 
			g.name, 
			g.supervisor_id, 
			u.full_name as supervisor_name,
			u.id as sup_id,
			u.full_name as sup_full_name,
			u.login as sup_login,
			u.role as sup_role,
			g.created_at,
			g.is_deleted,
			g.deleted_at
		FROM groups g
		LEFT JOIN system_users u ON g.supervisor_id = u.id
		%s
		ORDER BY g.id
		LIMIT $%d OFFSET $%d`, whereClause, len(args)+1, len(args)+2)

	args = append(args, limit, offset)

	rows, err := g.DB.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("list groups: %w", err)
	}
	defer rows.Close()

	groups := make([]model.Group, 0)

	for rows.Next() {
		var g model.Group
		var supID *uuid.UUID
		var supFullName, supLogin, supRole *string

		err := rows.Scan(
			&g.ID,
			&g.Name,
			&g.SupervisorID,
			&g.SupervisorName,
			&supID,
			&supFullName,
			&supLogin,
			&supRole,
			&g.CreatedAt,
			&g.IsDeleted,
			&g.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("scan group: %w", err)
		}

		if supID != nil {
			g.Supervisor = &model.GroupSupervisor{
				ID:       *supID,
				FullName: *supFullName,
				Login:    *supLogin,
				Role:     *supRole,
			}
		}

		groups = append(groups, g)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("rows groups: %w", err)
	}

	return groups, total, nil
}

func (g *groupsRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Group, error) {
	sql := `
		SELECT 
			g.id,
			g.name,
			g.supervisor_id,
			s.full_name as supervisor_name,
			s.id as sup_id,
			s.full_name as sup_full_name,
			s.login as sup_login,
			s.role as sup_role,
			g.created_at,
			g.is_deleted,
			g.deleted_at
		FROM groups g
		LEFT JOIN system_users s ON g.supervisor_id = s.id
		WHERE g.id = $1
	`

	var gr model.Group
	var supID *uuid.UUID
	var supFullName, supLogin, supRole *string

	err := g.DB.Pool.QueryRow(ctx, sql, id).Scan(
		&gr.ID,
		&gr.Name,
		&gr.SupervisorID,
		&gr.SupervisorName,
		&supID,
		&supFullName,
		&supLogin,
		&supRole,
		&gr.CreatedAt,
		&gr.IsDeleted,
		&gr.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.ErrorNotFound
		}
		return nil, err
	}

	if supID != nil {
		gr.Supervisor = &model.GroupSupervisor{
			ID:       *supID,
			FullName: *supFullName,
			Login:    *supLogin,
			Role:     *supRole,
		}
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

func (g *groupsRepository) GetGroupIDForSupervisor(ctx context.Context, userID uuid.UUID) (uuid.UUID, error) {
	var groupID uuid.UUID
	err := g.DB.Pool.QueryRow(ctx, `
		SELECT id FROM groups WHERE supervisor_id = $1
	`, userID).Scan(&groupID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return uuid.Nil, fmt.Errorf("no group found for this supervisor")
		}
		return uuid.Nil, fmt.Errorf("get group id for supervisor: %w", err)
	}

	return groupID, nil
}

func (g *groupsRepository) Delete(ctx context.Context, id uuid.UUID) error {
	sql := `UPDATE groups 
		SET is_deleted = true,
			deleted_at = now()
		WHERE id = $1`

	ct, err := g.DB.Pool.Exec(ctx, sql, id)
	if err != nil {
		return fmt.Errorf("delete group: %w", err)
	}

	if ct.RowsAffected() == 0 {
		return errs.ErrorNotFound
	}

	return nil
}

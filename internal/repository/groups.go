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
	List(ctx context.Context, groupType string, supervisorID *uuid.UUID, limit, offset int) ([]model.Group, int, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.Group, error)
	Update(ctx context.Context, gr model.Group, id uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetGroupIDForSupervisor(ctx context.Context, userID uuid.UUID) (uuid.UUID, error)
	GetHierarchy(ctx context.Context, groupType string) ([]model.Group, error)
	GetDescendantIDs(ctx context.Context, id uuid.UUID) ([]uuid.UUID, error)
	GetDescendants(ctx context.Context, id uuid.UUID) ([]model.Group, error)
}

func NewGroupsRepository(db *db.DB) GroupsRepository {
	return &groupsRepository{DB: db}
}

func (g *groupsRepository) Create(ctx context.Context, group *model.Group) error {
	sql := `
		INSERT INTO groups
		(
			name,
			type,
			parent_id,
			supervisor_id,
			color,
			layer,
			cluster
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id;
	`
	err := g.DB.Pool.QueryRow(ctx, sql, group.Name, group.Type, group.ParentID, group.SupervisorID, group.Color, group.Layer, group.Cluster).Scan(&group.ID)
	if err != nil {
		return fmt.Errorf("create group: %w", err)
	}

	return nil
}

func (g *groupsRepository) List(ctx context.Context, groupType string, supervisorID *uuid.UUID, limit, offset int) ([]model.Group, int, error) {
	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	var conditions []string
	var args []any

	conditions = append(conditions, "g.is_deleted = false")

	if groupType != "" {
		conditions = append(conditions, fmt.Sprintf("g.type = $%d", len(args)+1))
		args = append(args, groupType)
	}

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
			g.type,
			g.parent_id,
			p.name as parent_name,
			g.supervisor_id, 
			u.full_name as supervisor_name,
			u.id as sup_id,
			u.full_name as sup_full_name,
			u.login as sup_login,
			u.role as sup_role,
			g.color,
			g.layer,
			g.cluster,
			g.created_at,
			g.is_deleted,
			g.deleted_at
		FROM groups g
		LEFT JOIN groups p ON g.parent_id = p.id
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
			&g.Type,
			&g.ParentID,
			&g.ParentName,
			&g.SupervisorID,
			&g.SupervisorName,
			&supID,
			&supFullName,
			&supLogin,
			&supRole,
			&g.Color,
			&g.Layer,
			&g.Cluster,
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
			g.type,
			g.supervisor_id,
			g.parent_id,
			COALESCE(s.full_name, '') as supervisor_name,
			COALESCE(s.id, '00000000-0000-0000-0000-000000000000'::uuid) as sup_id,
			COALESCE(s.full_name, '') as sup_full_name,
			COALESCE(s.login, '') as sup_login,
			COALESCE(s.role, '') as sup_role,
			g.color,
			g.layer,
			g.cluster,
			g.created_at,
			g.is_deleted,
			g.deleted_at
		FROM groups g
		LEFT JOIN system_users s ON g.supervisor_id = s.id
		WHERE g.id = $1
	`

	var gr model.Group
	var supFullName, supLogin, supRole string
	var supIDVal uuid.UUID

	err := g.DB.Pool.QueryRow(ctx, sql, id).Scan(
		&gr.ID,
		&gr.Name,
		&gr.Type,
		&gr.SupervisorID,
		&gr.ParentID,
		&gr.SupervisorName,
		&supIDVal,
		&supFullName,
		&supLogin,
		&supRole,
		&gr.Color,
		&gr.Layer,
		&gr.Cluster,
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

	// Check if supervisor exists (not the default zero UUID)
	zeroUUID := uuid.UUID{}
	if supIDVal != zeroUUID {
		gr.Supervisor = &model.GroupSupervisor{
			ID:       supIDVal,
			FullName: supFullName,
			Login:    supLogin,
			Role:     supRole,
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

	if gr.Type != "" {
		args = append(args, gr.Type)
		setClauses = append(setClauses, fmt.Sprintf("type = $%d", len(args)))
	}

	if gr.ParentID != nil {
		args = append(args, gr.ParentID)
		setClauses = append(setClauses, fmt.Sprintf("parent_id = $%d", len(args)))
	}

	if gr.SupervisorID != nil {
		args = append(args, gr.SupervisorID)
		setClauses = append(setClauses, fmt.Sprintf("supervisor_id = $%d", len(args)))
	}

	if gr.Color != "" {
		args = append(args, gr.Color)
		setClauses = append(setClauses, fmt.Sprintf("color = $%d", len(args)))
	}

	if gr.Layer != 0 {
		args = append(args, gr.Layer)
		setClauses = append(setClauses, fmt.Sprintf("layer = $%d", len(args)))
	}

	if gr.Cluster != "" {
		args = append(args, gr.Cluster)
		setClauses = append(setClauses, fmt.Sprintf("cluster = $%d", len(args)))
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

func (g *groupsRepository) GetHierarchy(ctx context.Context, groupType string) ([]model.Group, error) {
	sql := `
		WITH RECURSIVE group_tree AS (
			SELECT 
				g.id, g.name, g.type, g.parent_id, g.supervisor_id, g.color, g.layer, g.cluster, g.created_at, g.is_deleted, g.deleted_at,
				s.full_name as supervisor_name
			FROM groups g
			LEFT JOIN system_users s ON g.supervisor_id = s.id
			WHERE g.parent_id IS NULL AND g.type = $1 AND g.is_deleted = false
			
			UNION ALL
			
			SELECT 
				g.id, g.name, g.type, g.parent_id, g.supervisor_id, g.color, g.layer, g.cluster, g.created_at, g.is_deleted, g.deleted_at,
				s.full_name as supervisor_name
			FROM groups g
			LEFT JOIN system_users s ON g.supervisor_id = s.id
			INNER JOIN group_tree gt ON g.parent_id = gt.id
			WHERE g.is_deleted = false
		)
		SELECT * FROM group_tree ORDER BY parent_id NULLS FIRST, name
	`

	rows, err := g.DB.Pool.Query(ctx, sql, groupType)
	if err != nil {
		return nil, fmt.Errorf("get hierarchy: %w", err)
	}
	defer rows.Close()

	groupMap := make(map[uuid.UUID]*model.Group)

	for rows.Next() {
		gr := &model.Group{}
		var supervisorName *string
		err := rows.Scan(&gr.ID, &gr.Name, &gr.Type, &gr.ParentID, &gr.SupervisorID, &gr.Color, &gr.Layer, &gr.Cluster, &gr.CreatedAt, &gr.IsDeleted, &gr.DeletedAt, &supervisorName)
		if err != nil {
			return nil, fmt.Errorf("scan hierarchy: %w", err)
		}
		if supervisorName != nil {
			gr.SupervisorName = supervisorName
		}
		groupMap[gr.ID] = gr
	}

	for _, gr := range groupMap {
		if gr.ParentID != nil {
			if parent, exists := groupMap[*gr.ParentID]; exists {
				parent.Groups = append(parent.Groups, *gr)
			}
		}
	}

	var roots []model.Group
	for _, gr := range groupMap {
		if gr.ParentID == nil {
			roots = append(roots, *gr)
		}
	}

	return roots, nil
}

func (g *groupsRepository) GetDescendants(ctx context.Context, id uuid.UUID) ([]model.Group, error) {
	sql := `
		WITH RECURSIVE group_tree AS (
			SELECT id FROM groups WHERE id = $1 AND is_deleted = false
			UNION ALL
			SELECT g.id FROM groups g
			INNER JOIN group_tree gt ON g.parent_id = gt.id
			WHERE g.is_deleted = false
		)
		SELECT
			g.id,
			g.name,
			g.type,
			g.parent_id,
			p.name as parent_name,
			g.supervisor_id,
			u.full_name as supervisor_name,
			u.id as sup_id,
			u.full_name as sup_full_name,
			u.login as sup_login,
			u.role as sup_role,
			g.color,
			g.layer,
			g.cluster,
			g.created_at,
			g.is_deleted,
			g.deleted_at
		FROM group_tree gt
		JOIN groups g ON g.id = gt.id
		LEFT JOIN groups p ON g.parent_id = p.id
		LEFT JOIN system_users u ON g.supervisor_id = u.id
		ORDER BY g.id
	`

	rows, err := g.DB.Pool.Query(ctx, sql, id)
	if err != nil {
		return nil, fmt.Errorf("get descendants: %w", err)
	}
	defer rows.Close()

	var groups []model.Group
	for rows.Next() {
		var gr model.Group
		var supID *uuid.UUID
		var supFullName, supLogin, supRole *string

		if err := rows.Scan(
			&gr.ID,
			&gr.Name,
			&gr.Type,
			&gr.ParentID,
			&gr.ParentName,
			&gr.SupervisorID,
			&gr.SupervisorName,
			&supID,
			&supFullName,
			&supLogin,
			&supRole,
			&gr.Color,
			&gr.Layer,
			&gr.Cluster,
			&gr.CreatedAt,
			&gr.IsDeleted,
			&gr.DeletedAt,
		); err != nil {
			return nil, fmt.Errorf("scan descendant: %w", err)
		}

		if supID != nil {
			gr.Supervisor = &model.GroupSupervisor{
				ID:       *supID,
				FullName: *supFullName,
				Login:    *supLogin,
				Role:     *supRole,
			}
		}

		groups = append(groups, gr)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows descendants: %w", err)
	}

	return groups, nil
}

func (g *groupsRepository) GetDescendantIDs(ctx context.Context, id uuid.UUID) ([]uuid.UUID, error) {
	sql := `
		WITH RECURSIVE group_tree AS (
			SELECT id FROM groups WHERE id = $1 AND is_deleted = false
			UNION ALL
			SELECT g.id FROM groups g
			INNER JOIN group_tree gt ON g.parent_id = gt.id
			WHERE g.is_deleted = false
		)
		SELECT id FROM group_tree;
	`

	rows, err := g.DB.Pool.Query(ctx, sql, id)
	if err != nil {
		return nil, fmt.Errorf("get descendant ids: %w", err)
	}
	defer rows.Close()

	var ids []uuid.UUID
	for rows.Next() {
		var grID uuid.UUID
		if err := rows.Scan(&grID); err != nil {
			return nil, fmt.Errorf("scan descendant id: %w", err)
		}
		ids = append(ids, grID)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows descendant ids: %w", err)
	}

	return ids, nil
}

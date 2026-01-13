package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/rakhshon-mirzoev/gps-tracker/internal/model"
	"github.com/rakhshon-mirzoev/gps-tracker/pkg/db"
)

type groupsRepository struct {
	DB *db.DB
}

type GroupsRepository interface {
	CreateGroup(context.Context, *model.Group) error
}

func InitGroupsRepository(db *db.DB) GroupsRepository {
	return &groupsRepository{DB: db}
}

func (g *groupsRepository) CreateGroup(ctx context.Context, group *model.Group) error {
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
		return err
	}

	return nil
}

func (g *groupsRepository) GetGroups(ctx context.Context, supervisorID *int64) ([]model.Group, error) {
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
		return nil, err
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
			return nil, err
		}

		groups = append(groups, g)
	}

	return groups, nil
}

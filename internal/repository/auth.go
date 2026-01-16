package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/model"
	"github.com/rakhshon-mirzoev/gps-tracker/pkg/db"
	errs "github.com/rakhshon-mirzoev/gps-tracker/pkg/error"
)

type Auth interface {
	GetByLoginUnified(ctx context.Context, login string) (*model.User, error)
}

type authRepository struct {
	DB *db.DB
}

func InitAuthRepository(db *db.DB) Auth {
	return &authRepository{DB: db}
}

func (r *authRepository) GetByLoginUnified(ctx context.Context, login string) (*model.User, error) {

	sql := `SELECT *
	FROM (
		SELECT id, password_hash, NULL AS role, group_id, 'staff' AS user_type
		FROM staff_users
		WHERE login = $1

		UNION ALL

		SELECT id, password_hash, role, NULL AS group_id, 'system' AS user_type
		FROM system_users
		WHERE login = $1
	) u
	LIMIT 1;

	`

	var u model.User

	err := r.DB.Pool.QueryRow(ctx, sql, login).Scan(
		&u.ID,
		&u.PasswordHash,
		&u.Role,
		&u.GroupID,
		&u.Type,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errs.ErrorNotFound
		}
		return nil, err
	}

	return &u, nil
}

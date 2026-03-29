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

type Auth interface {
	GetByLoginUnified(ctx context.Context, login string) (*model.User, error)
	UpsertRefreshToken(ctx context.Context, ref *model.Refresh) error
	GetRefreshByToken(ctx context.Context, token string) (*model.Refresh, error)
	DeleteRefreshByID(ctx context.Context, id uuid.UUID) error
	DeleteRefreshByToken(ctx context.Context, token string) error
}

type authRepository struct {
	DB *db.DB
}

func NewAuthRepository(db *db.DB) Auth {
	return &authRepository{DB: db}
}

func (r *authRepository) GetByLoginUnified(ctx context.Context, login string) (*model.User, error) {

	sql := `SELECT *
	FROM (
		SELECT id, password_hash, role, group_id, 'staff' AS user_type
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
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.ErrorNotFound
		}
		return nil, err
	}

	return &u, nil
}

func (r *authRepository) UpsertRefreshToken(ctx context.Context, ref *model.Refresh) error {
	sql := `
		INSERT INTO refresh_token (
			user_id,
			user_type,
			role,
			token,
			expires_at
		)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (user_id, user_type)
		DO UPDATE SET
			token = EXCLUDED.token,
			expires_at = EXCLUDED.expires_at;
	`

	_, err := r.DB.Pool.Exec(
		ctx,
		sql,
		ref.UserID,
		ref.UserType,
		ref.Role,
		ref.Token,
		ref.ExpiresAt,
	)

	if err != nil {
		return fmt.Errorf("upsert refresh token: %w", err)
	}

	return nil
}

func (r *authRepository) GetRefreshByToken(ctx context.Context, token string) (*model.Refresh, error) {

	ref := &model.Refresh{}

	sql := `
		SELECT id, user_id, user_type, role, expires_at
		FROM refresh_token
		WHERE token = $1;
	`

	err := r.DB.Pool.QueryRow(ctx, sql, token).Scan(
		&ref.ID,
		&ref.UserID,
		&ref.UserType,
		&ref.Role,
		&ref.ExpiresAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.ErrorNotFound
		}
		return nil, err
	}

	return ref, nil
}

func (r *authRepository) DeleteRefreshByID(ctx context.Context, id uuid.UUID) error {
	sql := `DELETE FROM refresh_token WHERE id = $1`
	ct, err := r.DB.Pool.Exec(ctx, sql, id)
	if err != nil {
		return fmt.Errorf("delete refresh token by id: %w", err)
	}

	if ct.RowsAffected() == 0 {
		return errs.ErrorNotFound
	}

	return nil
}

func (r *authRepository) DeleteRefreshByToken(ctx context.Context, token string) error {
	sqlUpd := `UPDATE staff_users SET fcm_token = '' WHERE id = (select user_id from refresh_token where token = $1 and user_type = 'staff')`

	r.DB.Pool.Exec(ctx, sqlUpd, token)

	sql := `DELETE FROM refresh_token WHERE token = $1`

	ct, err := r.DB.Pool.Exec(ctx, sql, token)
	if err != nil {
		return fmt.Errorf("delete refresh token by token: %w", err)
	}

	if ct.RowsAffected() == 0 {
		return errs.ErrorNotFound
	}

	return nil
}

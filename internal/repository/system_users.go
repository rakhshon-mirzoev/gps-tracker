package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/model"
	"github.com/rakhshon-mirzoev/gps-tracker/pkg/db"
	errors "github.com/rakhshon-mirzoev/gps-tracker/pkg/error"
)

type systemUsers struct {
	DB *db.DB
}

type SystemUsersRepository interface {
	Create(ctx context.Context, sysUser *model.SystemUser) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.SystemUser, error)
	GetByLogin(ctx context.Context, login string) (*model.SystemUser, error) // Для логина
	Update(ctx context.Context, sysUser *model.SystemUser) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context) ([]model.SystemUser, error)
	Deactivate(ctx context.Context, id, dId uuid.UUID, dReason string) error
}

func InitSystemUsersRepository(db *db.DB) SystemUsersRepository {
	return &systemUsers{DB: db}
}

func (s *systemUsers) Create(ctx context.Context, sysUser *model.SystemUser) (err error) {
	sql := `
		INSERT INTO system_users 
		(
			full_name, 
			login, 
			password_hash, 
			role, 
			is_active
		) 
		VALUES($1, $2, $3, $4, $5) 
		RETURNING id, created_at;
	`

	err = s.DB.Pool.QueryRow(
		ctx,
		sql,
		sysUser.FullName,
		sysUser.Login,
		sysUser.PasswordHash,
		sysUser.Role,
		sysUser.IsActive).Scan(&sysUser.ID, &sysUser.CreatedAt)

	if err != nil {
		return err
	}

	return nil
}

func (s *systemUsers) GetByID(ctx context.Context, id uuid.UUID) (*model.SystemUser, error) {
	u := &model.SystemUser{}
	sql := `SELECT id, full_name, login, role, is_active, created_at FROM system_users WHERE id = $1`

	err := s.DB.Pool.QueryRow(ctx, sql, id).Scan(
		&u.ID, &u.FullName, &u.Login, &u.Role, &u.IsActive, &u.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return u, nil
}

func (s *systemUsers) GetByLogin(ctx context.Context, login string) (*model.SystemUser, error) {
	u := &model.SystemUser{}
	sql := `SELECT id, full_name, login, password_hash, role, is_active FROM system_users WHERE login = $1`

	err := s.DB.Pool.QueryRow(ctx, sql, login).Scan(
		&u.ID, &u.FullName, &u.Login, &u.PasswordHash, &u.Role, &u.IsActive,
	)
	if err != nil {
		if err == errors.ErrorNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("get user by login: %w", err)
	}
	return u, nil
}

func (s *systemUsers) Update(ctx context.Context, sysUser *model.SystemUser) error {
	sql := `
        UPDATE system_users 
        SET full_name = $1, role = $2, is_active = $3 
        WHERE id = $4`

	_, err := s.DB.Pool.Exec(ctx, sql,
		sysUser.FullName,
		sysUser.Role,
		sysUser.IsActive,
		sysUser.ID,
	)
	if err != nil {
		return fmt.Errorf("update user: %w", err)
	}
	return nil
}

func (s *systemUsers) Delete(ctx context.Context, id uuid.UUID) error {
	sql := `DELETE FROM system_users WHERE id = $1`
	_, err := s.DB.Pool.Exec(ctx, sql, id)
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}
	return nil
}

func (s *systemUsers) List(ctx context.Context) ([]model.SystemUser, error) {
	sql := `SELECT id, full_name, login, role, is_active, created_at FROM system_users ORDER BY created_at DESC`

	rows, err := s.DB.Pool.Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.SystemUser
	for rows.Next() {
		var u model.SystemUser
		if err := rows.Scan(&u.ID, &u.FullName, &u.Login, &u.Role, &u.IsActive, &u.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (s *systemUsers) Deactivate(ctx context.Context, id, dId uuid.UUID, dReason string) error {
	sql := `
		SET is_active = false,
		 	deactivated_by_id = $1
			deactivated_at = $2
			deactivation_reason = $3
		WHERE id = $4
	`

	_, err := s.DB.Pool.Exec(ctx, sql, dId, time.Now(), dReason, id)
	if err != nil {
		return err
	}

	return nil
}

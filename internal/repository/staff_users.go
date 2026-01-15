package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/model"
	"github.com/rakhshon-mirzoev/gps-tracker/pkg/db"
)

type StaffUsersRepository interface {
	Create(ctx context.Context, user *model.StaffUser) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.StaffUser, error)
	List(ctx context.Context) ([]*model.StaffUser, error)
	Update(ctx context.Context, user *model.StaffUser) error
	Delete(ctx context.Context, id uuid.UUID) error
	Deactivate(ctx context.Context, id, dId uuid.UUID, dReason string) error
}

type staffUsers struct {
	DB *db.DB
}

func InitStaffUsersRepository(db *db.DB) StaffUsersRepository {
	return &staffUsers{DB: db}
}

func (r *staffUsers) Create(ctx context.Context, u *model.StaffUser) error {
	query := `
		INSERT INTO staff_users (
			group_id, full_name, login, phone_number, inn,
			birthday, password_hash, is_temporary_password
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
		RETURNING id, created_at
	`

	return r.DB.Pool.QueryRow(
		ctx,
		query,
		u.GroupID,
		u.FullName,
		u.Login,
		u.PhoneNumber,
		u.Inn,
		u.Birthday,
		u.PasswordHash,
		u.IsTemporaryPassword,
	).Scan(&u.ID, &u.CreatedAt)
}

func (r *staffUsers) GetByID(ctx context.Context, id uuid.UUID) (*model.StaffUser, error) {
	query := `
		SELECT id, group_id, full_name, login, is_active, created_at
		FROM staff_users
		WHERE id = $1
	`

	var u model.StaffUser

	err := r.DB.Pool.QueryRow(ctx, query, id).Scan(
		&u.ID,
		&u.GroupID,
		&u.FullName,
		&u.Login,
		&u.IsActive,
		&u.CreatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &u, nil
}

func (r *staffUsers) List(ctx context.Context) ([]*model.StaffUser, error) {
	query := `
		SELECT id, full_name, login, is_active, created_at
		FROM staff_users
		ORDER BY created_at DESC
	`

	rows, err := r.DB.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*model.StaffUser

	for rows.Next() {
		var u model.StaffUser
		if err := rows.Scan(
			&u.ID,
			&u.FullName,
			&u.Login,
			&u.IsActive,
			&u.CreatedAt,
		); err != nil {
			return nil, err
		}
		users = append(users, &u)
	}

	return users, nil
}

func (r *staffUsers) Update(ctx context.Context, u *model.StaffUser) error {
	query := `
		UPDATE staff_users
		SET full_name = $1,
		    phone_number = $2,
		    is_active = $3
		WHERE id = $4
	`

	_, err := r.DB.Pool.Exec(
		ctx,
		query,
		u.FullName,
		u.PhoneNumber,
		u.IsActive,
		u.ID,
	)

	return err
}

func (r *staffUsers) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE staff_users
		SET is_active = false,
		    deactivated_at = now()
		WHERE id = $1
	`

	_, err := r.DB.Pool.Exec(ctx, query, id)
	return err
}

func (s *staffUsers) Deactivate(ctx context.Context, id, dId uuid.UUID, dReason string) error {
	sql := `
		UPDATE staff_users
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

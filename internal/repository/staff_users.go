package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/model"
	"github.com/rakhshon-mirzoev/gps-tracker/pkg/db"
	errs "github.com/rakhshon-mirzoev/gps-tracker/pkg/error"
)

type StaffUsersRepository interface {
	Create(ctx context.Context, user *model.StaffUser) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.StaffUser, error)
	List(ctx context.Context) ([]model.StaffUser, error)
	Update(ctx context.Context, id uuid.UUID, dto model.UpdateStaffUser) error
	Delete(ctx context.Context, id uuid.UUID) error
	Deactivate(ctx context.Context, id, dId uuid.UUID, dReason string) error
}

type staffUsers struct {
	DB *db.DB
}

func NewStaffUsersRepository(db *db.DB) StaffUsersRepository {
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
		SELECT id, group_id, full_name, role, login, is_active, created_at
		FROM staff_users
		WHERE id = $1
	`

	var u model.StaffUser

	err := r.DB.Pool.QueryRow(ctx, query, id).Scan(
		&u.ID,
		&u.GroupID,
		&u.FullName,
		&u.Role,
		&u.Login,
		&u.IsActive,
		&u.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.ErrorNotFound
		}
		return nil, err
	}

	return &u, nil
}

func (r *staffUsers) List(ctx context.Context) ([]model.StaffUser, error) {
	query := `
		SELECT id, full_name, role, login, is_active, created_at
		FROM staff_users
		ORDER BY created_at DESC
	`

	rows, err := r.DB.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.StaffUser

	for rows.Next() {
		var u model.StaffUser
		if err := rows.Scan(
			&u.ID,
			&u.FullName,
			&u.Role,
			&u.Login,
			&u.IsActive,
			&u.CreatedAt,
		); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *staffUsers) Update(ctx context.Context, id uuid.UUID, dto model.UpdateStaffUser) error {

	query := `UPDATE staff_users SET `
	args := []any{}
	n := 1

	if dto.FullName != nil {
		query += fmt.Sprintf("full_name = $%d, ", n)
		args = append(args, *dto.FullName)
		n++
	}

	if dto.PhoneNumber != nil {
		query += fmt.Sprintf("phone_number = $%d, ", n)
		args = append(args, *dto.PhoneNumber)
		n++
	}

	if dto.IsActive != nil {
		query += fmt.Sprintf("is_active = $%d, ", n)
		args = append(args, *dto.IsActive)
		n++
	}

	if dto.GroupID != nil {
		query += fmt.Sprintf("group_id = $%d, ", n)
		args = append(args, *dto.GroupID)
		n++
	}

	if dto.FcmToken != nil {
		query += fmt.Sprintf("fcm_token = $%d, ", n)
		args = append(args, *dto.FcmToken)
		n++
	}

	if dto.DeviceId != nil {
		query += fmt.Sprintf("device_id = $%d, ", n)
		args = append(args, *dto.DeviceId)
		n++
	}

	if dto.DeactivatedAt != nil {
		query += fmt.Sprintf("deactivated_at = $%d, ", n)
		args = append(args, *dto.DeactivatedAt)
		n++
	}

	if dto.DeactivatedByID != nil {
		query += fmt.Sprintf("deactivated_by_id = $%d, ", n)
		args = append(args, *dto.DeactivatedByID)
		n++
	}

	if dto.DeactivationReason != nil {
		query += fmt.Sprintf("deactivated_reason = $%d, ", n)
		args = append(args, *dto.DeactivationReason)
		n++
	}

	if len(args) == 0 {
		return nil
	}

	query = strings.TrimSuffix(query, ", ")
	query += fmt.Sprintf(" WHERE id = $%d", n)
	args = append(args, id)

	ct, err := r.DB.Pool.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	if ct.RowsAffected() == 0 {
		return errs.ErrorNotFound
	}

	return nil
}

func (r *staffUsers) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		DELETE FROM staff_users
		WHERE id = $1
	`

	ct, err := r.DB.Pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if ct.RowsAffected() == 0 {
		return errs.ErrorNotFound
	}

	return nil
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

	ct, err := s.DB.Pool.Exec(ctx, sql, dId, time.Now(), dReason, id)
	if err != nil {
		return err
	}

	if ct.RowsAffected() == 0 {
		return errs.ErrorNotFound
	}

	return nil
}

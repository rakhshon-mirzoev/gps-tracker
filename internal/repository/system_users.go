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

type systemUsers struct {
	DB *db.DB
}

type SystemUsersRepository interface {
	Create(ctx context.Context, sysUser *model.SystemUser) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.SystemUser, error)
	GetByLogin(ctx context.Context, login string) (*model.SystemUser, error)
	Update(ctx context.Context, sysUser *model.UpdateSystemUser, id uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context) ([]model.SystemUser, error)
	Deactivate(ctx context.Context, id, dId uuid.UUID, dReason string) error
	GetGroupList(ctx context.Context, groupId uuid.UUID) ([]model.StaffUser, error)
	GetAllAdmins(ctx context.Context) ([]model.SystemUser, error)
	GetSupervisorByGroupID(ctx context.Context, groupID uuid.UUID) (*model.SystemUser, error)
}

func NewSystemUsersRepository(db *db.DB) SystemUsersRepository {
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
		return fmt.Errorf("create system user: %w", err)
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
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.ErrorNotFound
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
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.ErrorNotFound
		}
		return nil, err
	}
	return u, nil
}

func (s *systemUsers) Update(ctx context.Context, dto *model.UpdateSystemUser, id uuid.UUID) error {
	sets := make([]string, 0)
	args := make([]any, 0)
	n := 1

	if dto.FullName != nil {
		sets = append(sets, fmt.Sprintf("full_name = $%d", n))
		args = append(args, *dto.FullName)
		n++
	}

	if dto.IsActive != nil {
		sets = append(sets, fmt.Sprintf("is_active = $%d", n))
		args = append(args, *dto.IsActive)
		n++
	}

	if dto.DeactivatedAt != nil {
		sets = append(sets, fmt.Sprintf("deactivated_at = $%d", n))
		args = append(args, *dto.DeactivatedAt)
		n++
	}

	if dto.DeactivatedByID != nil {
		sets = append(sets, fmt.Sprintf("deactivated_by_id = $%d", n))
		args = append(args, *dto.DeactivatedByID)
		n++
	}

	if dto.DeactivationReason != nil {
		sets = append(sets, fmt.Sprintf("deactivation_reason = $%d", n))
		args = append(args, *dto.DeactivationReason)
		n++
	}

	if len(sets) == 0 {
		return nil
	}

	query := fmt.Sprintf(
		"UPDATE system_users SET %s WHERE id = $%d",
		strings.Join(sets, ", "),
		n,
	)

	args = append(args, id)

	_, err := s.DB.Pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("update system user: %w", err)
	}

	return nil
}

func (s *systemUsers) Delete(ctx context.Context, id uuid.UUID) error {
	sql := `DELETE FROM system_users WHERE id = $1`
	ct, err := s.DB.Pool.Exec(ctx, sql, id)
	if err != nil {
		return fmt.Errorf("delete system user: %w", err)
	}

	if ct.RowsAffected() == 0 {
		return errs.ErrorNotFound
	}

	return nil
}

func (s *systemUsers) List(ctx context.Context) ([]model.SystemUser, error) {
	sql := `SELECT id, full_name, login, role, is_active, created_at FROM system_users ORDER BY created_at DESC`

	rows, err := s.DB.Pool.Query(ctx, sql)
	if err != nil {
		return nil, fmt.Errorf("list system users: %w", err)
	}
	defer rows.Close()

	var users = make([]model.SystemUser, 0)

	for rows.Next() {
		var u model.SystemUser
		if err := rows.Scan(&u.ID, &u.FullName, &u.Login, &u.Role, &u.IsActive, &u.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan system user: %w", err)
		}
		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows system users: %w", err)
	}

	return users, nil
}

func (s *systemUsers) Deactivate(ctx context.Context, id, dId uuid.UUID, dReason string) error {
	sql := `
		UPDATE system_users
		SET is_active = false,
		 	deactivated_by_id = $1,
			deactivated_at = $2,
			deactivation_reason = $3
		WHERE id = $4
	`

	ct, err := s.DB.Pool.Exec(ctx, sql, dId, time.Now(), dReason, id)
	if err != nil {
		return fmt.Errorf("deactivate system user: %w", err)
	}

	if ct.RowsAffected() == 0 {
		return errs.ErrorNotFound
	}

	return nil
}

func (s *systemUsers) GetGroupList(ctx context.Context, groupId uuid.UUID) ([]model.StaffUser, error) {
	sql := `
		SELECT	s.id, s.full_name, s.phone_number, s.inn, s.birthday, s.password_updated_at, s.device_id, s.fcm_token, s.app_version, 
				s.last_battery_level, s.last_seen_at, s.is_active, s.created_at
		FROM staff_users s 
		WHERE s.group_id = $1 AND s.deactivated_by_id IS NULL;
	`

	rows, err := s.DB.Pool.Query(ctx, sql, groupId)
	if err != nil {
		return nil, fmt.Errorf("get group list: %w", err)
	}
	defer rows.Close()

	var result = make([]model.StaffUser, 0)

	for rows.Next() {
		var su = model.StaffUser{}
		if err := rows.Scan(
			&su.ID,
			&su.FullName,
			&su.PhoneNumber,
			&su.Inn,
			&su.Birthday,
			&su.PasswordUpdatedAt,
			&su.DeviceID,
			&su.FcmToken,
			&su.AppVersion,
			&su.LastBatteryLevel,
			&su.LastSeenAt,
			&su.IsActive,
			&su.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan staff user: %w", err)
		}
		result = append(result, su)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows staff users group: %w", err)
	}

	return result, nil
}

// GetAllAdmins возвращает всех активных system_users с ролью admin.
func (s *systemUsers) GetAllAdmins(ctx context.Context) ([]model.SystemUser, error) {
	sql := `SELECT id, full_name, login, role, is_active, created_at
	        FROM system_users
	        WHERE role = 'admin' AND is_active = TRUE`

	rows, err := s.DB.Pool.Query(ctx, sql)
	if err != nil {
		return nil, fmt.Errorf("get all admins: %w", err)
	}
	defer rows.Close()

	var admins []model.SystemUser
	for rows.Next() {
		var u model.SystemUser
		if err := rows.Scan(&u.ID, &u.FullName, &u.Login, &u.Role, &u.IsActive, &u.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan admin: %w", err)
		}
		admins = append(admins, u)
	}
	return admins, rows.Err()
}

// GetSupervisorByGroupID возвращает supervisor группы (system_user).
func (s *systemUsers) GetSupervisorByGroupID(ctx context.Context, groupID uuid.UUID) (*model.SystemUser, error) {
	sql := `
		SELECT su.id, su.full_name, su.login, su.role, su.is_active, su.created_at
		FROM system_users su
		JOIN groups g ON g.supervisor_id = su.id
		WHERE g.id = $1 AND su.is_active = TRUE
		LIMIT 1
	`
	var u model.SystemUser
	err := s.DB.Pool.QueryRow(ctx, sql, groupID).Scan(
		&u.ID, &u.FullName, &u.Login, &u.Role, &u.IsActive, &u.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // группа без супервайзера — не ошибка
		}
		return nil, fmt.Errorf("get supervisor by group id: %w", err)
	}
	return &u, nil
}

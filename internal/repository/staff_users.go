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
	List(ctx context.Context, limit, offset int) ([]model.StaffUser, int, error)
	Update(ctx context.Context, id uuid.UUID, dto model.UpdateStaffUser) error
	Delete(ctx context.Context, id uuid.UUID) error
	Deactivate(ctx context.Context, id uuid.UUID, dReason model.DeactivateStaffUser) error
	ListByGroupID(ctx context.Context, user_id uuid.UUID) (staff []model.StaffUserLocation, err error)
	GetLocationByID(ctx context.Context, id uuid.UUID, begin, end time.Time) (*model.StaffUserLocationHistory, error)
	GetByGroupID(ctx context.Context, group_id uuid.UUID) ([]model.StaffUser, error)
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
		SELECT s.id, s.full_name, s.role, s.login, s.is_active, s.created_at, s.inn, s.birthday, s.is_temporary_password, s.group_id, s.phone_number,
			   s.password_updated_at, s.device_id, s.fcm_token, s.app_version, s.last_battery_level, s.last_seen_at, s.deactivated_at,
			   s.deactivated_by_id, s.deactivation_reason, s.avatar_key, s.start_at, s.end_at, s.day_started, g.name
		FROM staff_users s
		JOIN groups g ON g.id = s.group_id
		WHERE id = $1
	`

	var u model.StaffUser

	err := r.DB.Pool.QueryRow(ctx, query, id).Scan(
		&u.ID,
		&u.FullName,
		&u.Role,
		&u.Login,
		&u.IsActive,
		&u.CreatedAt,
		&u.Inn,
		&u.Birthday,
		&u.IsTemporaryPassword,
		&u.GroupID,
		&u.PhoneNumber,
		&u.PasswordUpdatedAt,
		&u.DeviceID,
		&u.FcmToken,
		&u.AppVersion,
		&u.LastBatteryLevel,
		&u.LastSeenAt,
		&u.DeactivatedAt,
		&u.DeactivatedByID,
		&u.DeactivationReason,
		&u.AvatarKey,
		&u.StartAt,
		&u.EndAt,
		&u.DayStarted,
		&u.GroupName,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.ErrorNotFound
		}
		return nil, err
	}

	return &u, nil
}

func (r *staffUsers) List(ctx context.Context, limit, offset int) ([]model.StaffUser, int, error) {
	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	var total int
	countSQL := `SELECT COUNT(*) FROM staff_users`
	if err := r.DB.Pool.QueryRow(ctx, countSQL).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count staff users: %w", err)
	}

	query := `
		SELECT s.id, s.full_name, s.role, s.login, s.is_active, s.created_at, s.inn, s.birthday, s.is_temporary_password, s.group_id, s.phone_number,
			   s.password_updated_at, s.device_id, s.fcm_token, s.app_version, s.last_battery_level, s.last_seen_at, s.deactivated_at,
			   s.deactivated_by_id, s.deactivation_reason, s.avatar_key, s.start_at, s.end_at, s.day_started, g.name
		FROM staff_users s
		JOIN groups g ON g.id = s.group_id
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.DB.Pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("list staff users: %w", err)
	}
	defer rows.Close()

	var users = make([]model.StaffUser, 0)

	for rows.Next() {
		var u model.StaffUser
		if err := rows.Scan(
			&u.ID,
			&u.FullName,
			&u.Role,
			&u.Login,
			&u.IsActive,
			&u.CreatedAt,
			&u.Inn,
			&u.Birthday,
			&u.IsTemporaryPassword,
			&u.GroupID,
			&u.PhoneNumber,
			&u.PasswordUpdatedAt,
			&u.DeviceID,
			&u.FcmToken,
			&u.AppVersion,
			&u.LastBatteryLevel,
			&u.LastSeenAt,
			&u.DeactivatedAt,
			&u.DeactivatedByID,
			&u.DeactivationReason,
			&u.StartAt,
			&u.EndAt,
			&u.DayStarted,
			&u.GroupName,
		); err != nil {
			return nil, 0, fmt.Errorf("scan staff user: %w", err)
		}
		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("rows staff users: %w", err)
	}

	return users, total, nil
}

func (r *staffUsers) Update(ctx context.Context, id uuid.UUID, dto model.UpdateStaffUser) error {

	query := `UPDATE staff_users SET `
	var args []any
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
	if dto.DeactivatedByID != nil {
		query += fmt.Sprintf("deactivated_by_id = $%d, ", n)
		args = append(args, *dto.DeactivatedByID)
		n++
	}

	if dto.DeactivationReason != nil {
		query += fmt.Sprintf("deactivation_reason = $%d, ", n)
		args = append(args, *dto.DeactivationReason)
		n++
	}

	if dto.DeactivatedAt != nil {
		query += fmt.Sprintf("deactivated_at = $%d, ", n)
		args = append(args, *dto.DeactivatedAt)
		n++
	}

	if dto.AvatarKey != nil {
		query += fmt.Sprintf("avatar_key = $%d, ", n)
		args = append(args, *dto.AvatarKey)
		n++
	}

	if dto.StartAt != nil {
		query += fmt.Sprintf("start_at = $%d, day_started = true, ", n)
		args = append(args, *dto.StartAt)
		n++
		if dto.EndAt == nil {
			query += "end_at = null, "
		}
	}

	if dto.EndAt != nil {
		query += fmt.Sprintf("end_at = $%d, day_started = false, ", n)
		args = append(args, *dto.EndAt)
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
		return fmt.Errorf("update staff user: %w", err)
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
		return fmt.Errorf("delete staff user: %w", err)
	}

	if ct.RowsAffected() == 0 {
		return errs.ErrorNotFound
	}

	return nil
}

func (s *staffUsers) Deactivate(ctx context.Context, id uuid.UUID, deactivate model.DeactivateStaffUser) error {
	sql := `
		UPDATE staff_users
		SET is_active = false,
		 	deactivated_by_id = $1,
			deactivated_at = $2,
			deactivation_reason = $3
		WHERE id = $4
	`

	ct, err := s.DB.Pool.Exec(ctx, sql, deactivate.DeactivatedByID, time.Now(), deactivate.DeactivationReason, id)
	if err != nil {
		return fmt.Errorf("deactivate staff user: %w", err)
	}

	if ct.RowsAffected() == 0 {
		return errs.ErrorNotFound
	}

	return nil
}

func (s *staffUsers) ListByGroupID(ctx context.Context, user_id uuid.UUID) ([]model.StaffUserLocation, error) {
	sql := `
		SELECT DISTINCT ON (s.id)
			s.id, s.full_name, s.login, s.inn, s.birthday, s.group_id, s.phone_number,
			s.device_id, s.fcm_token, s.app_version, s.last_battery_level, s.last_seen_at,
			s.start_at, s.end_at, s.day_started, g.name,
			l.id, l.latitude, l.longitude, l.battery_level, l.captured_at
		FROM staff_users s
		JOIN groups g ON g.id = s.group_id
		JOIN location_history l ON l.staff_id = s.id
		WHERE g.supervisor_id = $1 AND s.is_active = true
		ORDER BY s.id, l.captured_at DESC
	`

	rows, err := s.DB.Pool.Query(ctx, sql, user_id)
	if err != nil {
		return nil, fmt.Errorf("list staff by group id: %w", err)
	}
	defer rows.Close()

	var staff = make([]model.StaffUserLocation, 0)

	for rows.Next() {
		var st model.StaffUserLocation
		st.Location = &model.LocationHistory{}
		if err := rows.Scan(
			&st.ID,
			&st.FullName,
			&st.Login,
			&st.Inn,
			&st.Birthday,
			&st.GroupID,
			&st.PhoneNumber,
			&st.DeviceID,
			&st.FcmToken,
			&st.AppVersion,
			&st.LastBatteryLevel,
			&st.LastSeenAt,
			&st.StartAt,
			&st.EndAt,
			&st.DayStarted,
			&st.GroupName,
			&st.Location.ID,
			&st.Location.Latitude,
			&st.Location.Longitude,
			&st.Location.BatteryLevel,
			&st.Location.CapturedAt,
		); err != nil {
			return nil, fmt.Errorf("scan staff user location: %w", err)
		}
		staff = append(staff, st)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows staff users: %w", err)
	}

	return staff, nil
}

func (s *staffUsers) GetLocationByID(ctx context.Context, id uuid.UUID, begin, end time.Time) (*model.StaffUserLocationHistory, error) {
	sql := `
		SELECT
			s.id, s.full_name, s.login, s.inn, s.birthday, s.group_id, s.phone_number,
			s.device_id, s.fcm_token, s.app_version, s.last_battery_level, s.last_seen_at,
			s.start_at, s.end_at, s.day_started, g.name,
			l.id, l.latitude, l.longitude, l.battery_level, l.captured_at, l.status
		FROM staff_users s
		JOIN location_history l ON l.staff_id = s.id
		JOIN groups g ON g.id = s.group_id
		WHERE s.id = $1
		AND s.is_active = true
		AND (l.captured_at >= $2 AND l.captured_at < $3)
		ORDER BY l.captured_at
	`

	rows, err := s.DB.Pool.Query(ctx, sql, id, begin, end)
	if err != nil {
		return nil, fmt.Errorf("get location by id: %w", err)
	}
	defer rows.Close()

	var staff *model.StaffUserLocationHistory

	for rows.Next() {
		var loc model.LocationHistory

		if staff == nil {
			staff = &model.StaffUserLocationHistory{}
		}

		if err := rows.Scan(
			&staff.ID,
			&staff.FullName,
			&staff.Login,
			&staff.Inn,
			&staff.Birthday,
			&staff.GroupID,
			&staff.PhoneNumber,
			&staff.DeviceID,
			&staff.FcmToken,
			&staff.AppVersion,
			&staff.LastBatteryLevel,
			&staff.LastSeenAt,
			&staff.StartAt,
			&staff.EndAt,
			&staff.DayStarted,
			&staff.GroupName,
			&loc.ID,
			&loc.Latitude,
			&loc.Longitude,
			&loc.BatteryLevel,
			&loc.CapturedAt,
			&loc.Status,
		); err != nil {
			return nil, fmt.Errorf("scan location history: %w", err)
		}

		staff.Locations = append(staff.Locations, loc)
	}

	if staff == nil {
		return nil, errs.ErrorNotFound
	}

	return staff, nil
}

func (s *staffUsers) GetByGroupID(ctx context.Context, group_id uuid.UUID) ([]model.StaffUser, error) {
	sql := `id, full_name, role, login, is_active, created_at, inn, birthday, is_temporary_password, phone_number,
			   password_updated_at, device_id, fcm_token, app_version, last_battery_level, last_seen_at, deactivated_at,
			   deactivated_by_id, deactivation_reason, start_at, end_at, day_started
		FROM staff_users
		WHERE group_id = $1`

	rows, err := s.DB.Pool.Query(ctx, sql, group_id)
	if err != nil {
		return nil, fmt.Errorf("list staff users: %w", err)
	}
	defer rows.Close()

	var users = make([]model.StaffUser, 0)

	for rows.Next() {
		var u model.StaffUser
		if err := rows.Scan(
			&u.ID,
			&u.FullName,
			&u.Role,
			&u.Login,
			&u.IsActive,
			&u.CreatedAt,
			&u.Inn,
			&u.Birthday,
			&u.IsTemporaryPassword,
			&u.PhoneNumber,
			&u.PasswordUpdatedAt,
			&u.DeviceID,
			&u.FcmToken,
			&u.AppVersion,
			&u.LastBatteryLevel,
			&u.LastSeenAt,
			&u.DeactivatedAt,
			&u.DeactivatedByID,
			&u.DeactivationReason,
			&u.StartAt,
			&u.EndAt,
			&u.DayStarted,
		); err != nil {
			return nil, fmt.Errorf("scan staff user: %w", err)
		}
		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows staff users: %w", err)
	}

	return users, nil

}

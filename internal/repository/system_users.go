package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/model"
)

type systemUsersRepository struct {
	DB *pgxpool.Pool
}

type SystemUsersRepo interface {
}

func InitSystemUsersRepository(db *pgxpool.Pool) SystemUsersRepo {
	return &systemUsersRepository{DB: db}
}

func (s *systemUsersRepository) CreateSystemUser(ctx context.Context, sysUser *model.SystemUser) (err error) {
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

	err = s.DB.QueryRow(ctx, sql, sysUser.FullName, sysUser.Login, sysUser.PasswordHash, sysUser.Role, sysUser.IsActive).Scan(&sysUser.ID, &sysUser.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create sysUser: %v", err)
	}

	return nil
}

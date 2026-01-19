package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rakhshon-mirzoev/gps-tracker/config"
)

type DB struct {
	Pool   *pgxpool.Pool
	Ctx    context.Context
	Cancel context.CancelFunc
}

func ConnectToDb(cfg *config.DB) (*DB, error) {
	ctx, cancel := context.WithCancel(context.Background())

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.Username,
		cfg.Password,
		cfg.DbHost,
		cfg.DbPort,
		cfg.DbName,
	)

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		cancel()
		return nil, fmt.Errorf("error while parsing config into db: %v", err)
	}

	config.MaxConns = 10
	config.MinConns = 2
	config.MaxConnIdleTime = 5 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		cancel()
		return nil, fmt.Errorf("error while connecting to db: %v", err)
	}

	pingCtx, pingCancel := context.WithTimeout(ctx, 5*time.Second)
	defer pingCancel()

	if err := pool.Ping(pingCtx); err != nil {
		pool.Close()
		cancel()
		return nil, fmt.Errorf("failed while ping: %v", err)
	}

	return &DB{
		Pool:   pool,
		Ctx:    ctx,
		Cancel: cancel,
	}, nil

}

func (db *DB) Close() {
	db.Pool.Close()
	db.Cancel()
}

package websocket

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type Deps struct {
	Hub    *Hub
	DBPool *pgxpool.Pool
}

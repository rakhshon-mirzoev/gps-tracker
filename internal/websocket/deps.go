package websocket

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rakhshon-mirzoev/gps-tracker/pkg/token"
)

type Deps struct {
	Hub          *Hub
	DBPool       *pgxpool.Pool
	TokenService token.Token
}

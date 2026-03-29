package websocket

import (
	"context"

	"github.com/IkbolNazarov/gps-tracker/internal/events"
	"github.com/IkbolNazarov/gps-tracker/pkg/token"
	"github.com/google/uuid"
)

type LocationProvider interface {
	LatestByGroupID(ctx context.Context, groupID uuid.UUID) ([]events.LocationUpdate, error)
	LatestAll(ctx context.Context) ([]events.LocationUpdate, error)
}

type GroupResolver interface {
	GetGroupIDForSupervisor(ctx context.Context, userID uuid.UUID) (uuid.UUID, error)
}

type Deps struct {
	Hub              *Hub
	TokenService     token.Token
	LocationProvider LocationProvider
	GroupResolver    GroupResolver
}

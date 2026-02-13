package websocket

import (
	"context"

	"github.com/google/uuid"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/events"
	"github.com/rakhshon-mirzoev/gps-tracker/pkg/token"
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

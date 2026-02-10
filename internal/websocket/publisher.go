package websocket

import (
	"context"
	"encoding/json"

	"github.com/rakhshon-mirzoev/gps-tracker/internal/events"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/logger"
)

type wsLocationPublisher struct {
	hub *Hub
}

func NewLocationPublisher(hub *Hub) events.LocationPublisher {
	return &wsLocationPublisher{hub: hub}
}

func (p *wsLocationPublisher) PublishOnlineLocation(_ context.Context, update events.LocationUpdate) error {
	payload, err := json.Marshal(update)
	if err != nil {
		return err
	}
	select {
	case p.hub.Broadcast <- payload:
	default:
		logger.Log.Warn("p.hub ignore")
	}
	return nil
}

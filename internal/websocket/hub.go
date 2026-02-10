package websocket

import (
	"context"
	"encoding/json"

	"github.com/rakhshon-mirzoev/gps-tracker/internal/events"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/logger"
	"go.uber.org/zap"
)

type Hub struct {
	Clients    map[*Client]bool
	Broadcast  chan []byte
	Register   chan *Client
	Unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan []byte, 256),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

func (h *Hub) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case client := <-h.Register:
			h.Clients[client] = true

		case client := <-h.Unregister:
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.Send)
			}

		case message := <-h.Broadcast:
			var update events.LocationUpdate
			if err := json.Unmarshal(message, &update); err != nil {
				logger.Log.Error("invalid location update received", zap.Error(err))
				continue
			}

			for client := range h.Clients {
				if client.GroupID == update.GroupID {
					select {
					case client.Send <- message:
					default:
						close(client.Send)
						delete(h.Clients, client)
					}
				}
			}
		}
	}
}

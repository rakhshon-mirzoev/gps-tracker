package websocket

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/google/uuid"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/logger"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/model"
	"go.uber.org/zap"
)

// NotificationMessage — обёртка, которую кладут в канал Send хаба.
type NotificationMessage struct {
	RecipientID uuid.UUID
	Payload     []byte
}

// NotificationHub — адресный хаб: доставляет уведомления конкретным system_user-ам.
type NotificationHub struct {
	mu         sync.RWMutex
	clients    map[uuid.UUID][]*NotificationClient // recipientID -> список соединений
	Register   chan *NotificationClient
	Unregister chan *NotificationClient
	Send       chan NotificationMessage
}

func NewNotificationHub() *NotificationHub {
	return &NotificationHub{
		clients:    make(map[uuid.UUID][]*NotificationClient),
		Register:   make(chan *NotificationClient, 16),
		Unregister: make(chan *NotificationClient, 16),
		Send:       make(chan NotificationMessage, 256),
	}
}

func (h *NotificationHub) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return

		case client := <-h.Register:
			h.mu.Lock()
			h.clients[client.UserID] = append(h.clients[client.UserID], client)
			h.mu.Unlock()

		case client := <-h.Unregister:
			h.mu.Lock()
			if list, ok := h.clients[client.UserID]; ok {
				newList := list[:0]
				for _, c := range list {
					if c != client {
						newList = append(newList, c)
					}
				}
				if len(newList) == 0 {
					delete(h.clients, client.UserID)
				} else {
					h.clients[client.UserID] = newList
				}
				close(client.Send)
			}
			h.mu.Unlock()

		case msg := <-h.Send:
			h.mu.RLock()
			clients := h.clients[msg.RecipientID]
			h.mu.RUnlock()

			for _, c := range clients {
				select {
				case c.Send <- msg.Payload:
				default:
					logger.Log.Warn("notification client send buffer full, dropping message",
						zap.Stringer("user_id", c.UserID),
					)
				}
			}
		}
	}
}

// Push отправляет уведомление конкретному получателю всем его активным соединениям.
func (h *NotificationHub) Push(n model.SystemNotification) {
	payload, err := json.Marshal(n)
	if err != nil {
		logger.Log.Error("notification hub: marshal notification", zap.Error(err))
		return
	}
	select {
	case h.Send <- NotificationMessage{RecipientID: n.RecipientID, Payload: payload}:
	default:
		logger.Log.Warn("notification hub Send channel full, dropping",
			zap.Stringer("recipient_id", n.RecipientID),
		)
	}
}

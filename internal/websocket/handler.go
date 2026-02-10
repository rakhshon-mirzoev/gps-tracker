package websocket

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		return origin == "https://your-admin-domain.com"
	},
}

func getGroupIDForSystemUser(ctx context.Context, userID uuid.UUID, db *pgxpool.Pool) (uuid.UUID, error) {
	var groupID uuid.UUID
	err := db.QueryRow(ctx, `
		SELECT id 
		FROM groups 
		WHERE supervisor_id = $1
	`, userID).Scan(&groupID)

	if err == pgx.ErrNoRows {
		return uuid.Nil, fmt.Errorf("no group found for this supervisor")
	}
	if err != nil {
		return uuid.Nil, fmt.Errorf("query group id: %w", err)
	}

	return groupID, nil
}

func Handler(deps Deps) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := c.Get("user_id")
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		groupID, err := getGroupIDForSystemUser(c.Request.Context(), userID.(uuid.UUID), deps.DBPool)
		if err != nil {
			c.AbortWithStatusJSON(403, gin.H{"error": "no group assigned or access denied"})
			return
		}

		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			return
		}

		client := &Client{
			Hub:     deps.Hub,
			Conn:    conn,
			Send:    make(chan []byte, 256),
			UserID:  userID.(uuid.UUID),
			GroupID: groupID,
		}

		deps.Hub.Register <- client

		go client.WritePump()
		go client.ReadPump(c.Request.Context())
	}
}

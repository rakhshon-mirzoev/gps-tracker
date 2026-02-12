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
		return true
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
		tokenStr := c.Query("token")
		if tokenStr == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing token query parameter"})
			return
		}

		claims, err := deps.TokenService.ParseToken(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		if claims.Role != "admin" && claims.Role != "supervisor" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "access denied: admin or supervisor role required"})
			return
		}

		groupID, err := getGroupIDForSystemUser(c.Request.Context(), claims.UserID, deps.DBPool)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "no group assigned or access denied"})
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
			UserID:  claims.UserID,
			GroupID: groupID,
		}

		deps.Hub.Register <- client

		go client.WritePump()
		go client.ReadPump(c.Request.Context())
	}
}

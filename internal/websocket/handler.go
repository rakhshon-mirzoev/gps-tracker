package websocket

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/events"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/logger"
	"go.uber.org/zap"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
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

		ctx := c.Request.Context()

		var groupID uuid.UUID
		var initialLocations []events.LocationUpdate

		if claims.Role == "supervisor" {
			groupID, err = deps.GroupResolver.GetGroupIDForSupervisor(ctx, claims.UserID)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "no group assigned or access denied"})
				return
			}

			initialLocations, err = deps.LocationProvider.LatestByGroupID(ctx, groupID)
			if err != nil {
				logger.Log.Error("failed to get initial locations for supervisor", zap.Error(err))
			}
		} else {
			groupID = uuid.Nil

			initialLocations, err = deps.LocationProvider.LatestAll(ctx)
			if err != nil {
				logger.Log.Error("failed to get initial locations for admin", zap.Error(err))
			}
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
			Role:    claims.Role,
		}

		for _, loc := range initialLocations {
			payload, err := json.Marshal(loc)
			if err != nil {
				logger.Log.Error("failed to marshal initial location", zap.Error(err))
				continue
			}
			select {
			case client.Send <- payload:
			default:
				logger.Log.Warn("client send buffer full, skipping initial location")
			}
		}

		deps.Hub.Register <- client

		go client.WritePump()
		go client.ReadPump()
	}
}

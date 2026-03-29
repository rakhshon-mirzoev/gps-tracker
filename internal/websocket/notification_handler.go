package websocket

import (
	"net/http"

	"github.com/IkbolNazarov/gps-tracker/pkg/token"
	"github.com/gin-gonic/gin"
)

// NotificationHandler возвращает gin.HandlerFunc для WS-канала уведомлений.
// Подключение: GET /api/v1/admin/ws-notifications?token=<jwt>
func NotificationHandler(hub *NotificationHub, tokenSvc token.Token) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.Query("token")
		if tokenStr == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			return
		}

		claims, err := tokenSvc.ParseToken(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		if claims.Role != "admin" && claims.Role != "supervisor" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "admin or supervisor required"})
			return
		}

		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			return
		}

		client := &NotificationClient{
			Hub:    hub,
			Conn:   conn,
			Send:   make(chan []byte, 64),
			UserID: claims.UserID,
		}

		hub.Register <- client

		go client.WritePumpNotif()
		go client.ReadPumpNotif()
	}
}

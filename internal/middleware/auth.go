package middleware

import (
	"strings"

	"github.com/IkbolNazarov/gps-tracker/internal/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ctxKey string

const (
	CtxUserID  ctxKey = "user_id"
	CtxRole    ctxKey = "role"
	CtxGroupID ctxKey = "group_id"
)

func (s *MiddlewareService) Auth() gin.HandlerFunc {
	return func(c *gin.Context) {

		header := c.GetHeader("Authorization")
		if header == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "missing token"})
			return
		}

		const prefix = "Bearer "
		if !strings.HasPrefix(header, prefix) {
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid auth header"})
			return
		}

		tokenStr := strings.TrimPrefix(header, prefix)

		claims, err := s.t.ParseToken(tokenStr)
		if err != nil {
			logger.Log.Error("Auth: parse token error", zap.Error(err))
			c.AbortWithStatusJSON(401, gin.H{"error": err.Error()})
			return
		}

		c.Set(string(CtxUserID), claims.UserID)

		c.Set(string(CtxRole), claims.Role)

		if claims.GroupID != nil {
			c.Set(string(CtxGroupID), *claims.GroupID)
		}

		c.Next()
	}
}

package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/rakhshon-mirzoev/gps-tracker/pkg/token"
)

type Middleware interface {
	CorsMiddleWare() gin.HandlerFunc
	Auth() gin.HandlerFunc
}

type MiddlewareService struct {
	t token.Token
}

func InitMiddleware(t token.Token) Middleware {
	return &MiddlewareService{t: t}
}

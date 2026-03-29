package middleware

import (
	"github.com/IkbolNazarov/gps-tracker/pkg/token"
	"github.com/gin-gonic/gin"
)

type Middleware interface {
	CorsMiddleWare() gin.HandlerFunc
	Auth() gin.HandlerFunc
	Recovery() gin.HandlerFunc
	Logger() gin.HandlerFunc
	RequestID() gin.HandlerFunc
	Admin() gin.HandlerFunc
	Mobile() gin.HandlerFunc
	System() gin.HandlerFunc
	RateLimitMiddleware(limiter *IPRateLimiter) gin.HandlerFunc
}

type MiddlewareService struct {
	t token.Token
}

func New(t token.Token) Middleware {
	return &MiddlewareService{t: t}
}

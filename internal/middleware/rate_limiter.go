package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type item struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type IPRateLimiter struct {
	ips map[string]*item
	mu  sync.RWMutex
	r   rate.Limit
	b   int
}

func (i *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	v, exists := i.ips[ip]
	if !exists {
		limiter := rate.NewLimiter(i.r, i.b)
		i.ips[ip] = &item{limiter: limiter, lastSeen: time.Now()}
		return limiter
	}

	v.lastSeen = time.Now()
	return v.limiter
}

func (i *IPRateLimiter) cleanupVisitor() {
	for {
		time.Sleep(time.Minute * 10)
		i.mu.Lock()
		for ip, v := range i.ips {
			if time.Since(v.lastSeen) > time.Hour*1 {
				delete(i.ips, ip)
			}
		}
		i.mu.Unlock()
	}
}

func NewIPRateLimiter(r rate.Limit, b int) *IPRateLimiter {
	i := &IPRateLimiter{
		ips: make(map[string]*item),
		r:   r,
		b:   b,
	}
	go i.cleanupVisitor()
	return i
}

func (s *MiddlewareService) RateLimitMiddleware(limiter *IPRateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.ClientIP()

		if !limiter.GetLimiter(key).Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":       "Too many requests",
				"retry_after": "1s",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

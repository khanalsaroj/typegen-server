package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/khanalsaroj/typegen-server/internal/pkg/response"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type IPRateLimiter struct {
	ips map[string]*rate.Limiter
	mu  *sync.RWMutex
	r   rate.Limit
	b   int
}

func NewIPRateLimiter(r rate.Limit, b int) *IPRateLimiter {
	return &IPRateLimiter{
		ips: make(map[string]*rate.Limiter),
		mu:  &sync.RWMutex{},
		r:   r,
		b:   b,
	}
}

func (i *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	limiter, exists := i.ips[ip]
	if !exists {
		limiter = rate.NewLimiter(i.r, i.b)
		i.ips[ip] = limiter
	}

	return limiter
}

func RateLimit(rps int) gin.HandlerFunc {
	limiter := NewIPRateLimiter(rate.Limit(rps), rps*2)

	go func() {
		for {
			time.Sleep(time.Minute)
			limiter.mu.Lock()
			limiter.ips = make(map[string]*rate.Limiter)
			limiter.mu.Unlock()
		}
	}()

	return func(c *gin.Context) {
		ip := c.ClientIP()
		rateLimiter := limiter.GetLimiter(ip)

		if !rateLimiter.Allow() {
			response.Error(c, http.StatusTooManyRequests, "Rate limit exceeded", nil)
			c.Abort()
			return
		}

		c.Next()
	}
}

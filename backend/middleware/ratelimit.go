package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimiter implements a simple token bucket rate limiter
type RateLimiter struct {
	visitors map[string]*visitor
	mu       sync.Mutex
	rate     int           // requests per interval
	interval time.Duration // time interval
}

type visitor struct {
	tokens    int
	lastSeen  time.Time
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(rate int, interval time.Duration) *RateLimiter {
	rl := &RateLimiter{
		visitors: make(map[string]*visitor),
		rate:     rate,
		interval: interval,
	}

	// Cleanup old visitors every minute
	go rl.cleanupVisitors()

	return rl
}

func (rl *RateLimiter) cleanupVisitors() {
	for {
		time.Sleep(time.Minute)
		rl.mu.Lock()
		for ip, v := range rl.visitors {
			if time.Since(v.lastSeen) > 3*time.Minute {
				delete(rl.visitors, ip)
			}
		}
		rl.mu.Unlock()
	}
}

func (rl *RateLimiter) getVisitor(ip string) *visitor {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	v, exists := rl.visitors[ip]
	if !exists {
		rl.visitors[ip] = &visitor{tokens: rl.rate, lastSeen: time.Now()}
		return rl.visitors[ip]
	}

	// Refill tokens based on time elapsed
	elapsed := time.Since(v.lastSeen)
	tokensToAdd := int(elapsed / rl.interval) * rl.rate
	v.tokens += tokensToAdd
	if v.tokens > rl.rate {
		v.tokens = rl.rate
	}
	v.lastSeen = time.Now()

	return v
}

// RateLimitMiddleware returns a Gin middleware for rate limiting
func RateLimitMiddleware(rate int, interval time.Duration) gin.HandlerFunc {
	limiter := NewRateLimiter(rate, interval)

	return func(c *gin.Context) {
		ip := c.ClientIP()
		visitor := limiter.getVisitor(ip)

		if visitor.tokens <= 0 {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"success": false,
				"message": "Rate limit exceeded. Please try again later.",
			})
			c.Abort()
			return
		}

		visitor.tokens--
		c.Next()
	}
}

// StrictRateLimitMiddleware for sensitive endpoints like login
func StrictRateLimitMiddleware() gin.HandlerFunc {
	return RateLimitMiddleware(10, time.Minute) // 10 requests per minute
}

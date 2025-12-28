package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/SureshAmal/NimbusU-backend/shared/utils"
	"github.com/redis/go-redis/v9"
)

// RateLimiter implements rate limiting using Redis
type RateLimiter struct {
	client       *redis.Client
	maxRequests  int
	windowPeriod time.Duration
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(client *redis.Client, maxRequests int, windowPeriod time.Duration) *RateLimiter {
	return &RateLimiter{
		client:       client,
		maxRequests:  maxRequests,
		windowPeriod: windowPeriod,
	}
}

// RateLimitMiddleware creates a rate limiting middleware
func (rl *RateLimiter) RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Use IP address as identifier
		identifier := c.ClientIP()

		// Create Redis key
		key := fmt.Sprintf("rate_limit:%s", identifier)

		ctx := context.Background()

		// Increment request count
		count, err := rl.client.Incr(ctx, key).Result()
		if err != nil {
			// If Redis fails, allow the request but log the error
			c.Next()
			return
		}

		// Set expiry on first request
		if count == 1 {
			rl.client.Expire(ctx, key, rl.windowPeriod)
		}

		// Check if limit exceeded
		if count > int64(rl.maxRequests) {
			utils.ErrorResponse(c, http.StatusTooManyRequests, "Rate limit exceeded", nil)
			c.Abort()
			return
		}

		// Add rate limit headers
		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", rl.maxRequests))
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", rl.maxRequests-int(count)))

		c.Next()
	}
}

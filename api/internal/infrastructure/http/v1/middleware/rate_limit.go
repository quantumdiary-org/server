package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// RateLimiter implements a rate limiting mechanism
type RateLimiter struct {
	limits map[string]*rate.Limiter
	mutex  sync.RWMutex
	rate   rate.Limit
	burst  int
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(requestsPerSecond float64, burst int) *RateLimiter {
	return &RateLimiter{
		limits: make(map[string]*rate.Limiter),
		rate:   rate.Limit(requestsPerSecond),
		burst:  burst,
	}
}

// GetLimiter retrieves or creates a rate limiter for a specific key
func (rl *RateLimiter) GetLimiter(key string) *rate.Limiter {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	limiter, exists := rl.limits[key]
	if !exists {
		limiter = rate.NewLimiter(rl.rate, rl.burst)
		rl.limits[key] = limiter
	}

	return limiter
}

// RateLimitMiddleware applies rate limiting to requests
func (rl *RateLimiter) RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Use IP address as the key for rate limiting
		ip := c.ClientIP()
		limiter := rl.GetLimiter(ip)

		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded",
				"message": "Too many requests, please try again later",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RateLimitMiddlewareWithKey allows custom key for rate limiting
func (rl *RateLimiter) RateLimitMiddlewareWithKey(keyFunc func(*gin.Context) string) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := keyFunc(c)
		if key == "" {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal server error",
				"message": "Could not determine rate limit key",
			})
			c.Abort()
			return
		}

		limiter := rl.GetLimiter(key)

		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded",
				"message": "Too many requests, please try again later",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// TokenBucketRateLimiter implements a token bucket rate limiting algorithm
type TokenBucketRateLimiter struct {
	capacity  int
	refillRate time.Duration
	buckets   map[string]*TokenBucket
	mutex     sync.RWMutex
}

// TokenBucket represents a token bucket for a specific key
type TokenBucket struct {
	tokens    int
	lastRefill time.Time
	capacity  int
	refillRate time.Duration
}

// NewTokenBucket creates a new token bucket
func NewTokenBucket(capacity int, refillRate time.Duration) *TokenBucket {
	return &TokenBucket{
		tokens:     capacity,
		lastRefill: time.Now(),
		capacity:   capacity,
		refillRate: refillRate,
	}
}

// NewTokenBucketRateLimiter creates a new token bucket rate limiter
func NewTokenBucketRateLimiter(capacity int, refillRate time.Duration) *TokenBucketRateLimiter {
	return &TokenBucketRateLimiter{
		capacity:  capacity,
		refillRate: refillRate,
		buckets:   make(map[string]*TokenBucket),
	}
}

// Allow checks if a request is allowed based on token bucket algorithm
func (tb *TokenBucketRateLimiter) Allow(key string) bool {
	tb.mutex.Lock()
	defer tb.mutex.Unlock()

	bucket, exists := tb.buckets[key]
	if !exists {
		bucket = NewTokenBucket(tb.capacity, tb.refillRate)
		tb.buckets[key] = bucket
	}

	return bucket.consume(1)
}

// consume consumes tokens from the bucket
func (tb *TokenBucket) consume(tokens int) bool {
	tb.refill()
	
	if tb.tokens >= tokens {
		tb.tokens -= tokens
		return true
	}
	
	return false
}

// refill refills the bucket based on elapsed time
func (tb *TokenBucket) refill() {
	now := time.Now()
	elapsed := now.Sub(tb.lastRefill)
	
	if elapsed >= tb.refillRate {
		// Calculate how many tokens to add based on elapsed time
		tokensToAdd := int(elapsed / tb.refillRate)
		tb.tokens = min(tb.capacity, tb.tokens+tokensToAdd)
		tb.lastRefill = now
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// TokenBucketRateLimitMiddleware applies token bucket rate limiting
func (tb *TokenBucketRateLimiter) TokenBucketRateLimitMiddleware(keyFunc func(*gin.Context) string) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := keyFunc(c)
		if key == "" {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal server error",
				"message": "Could not determine rate limit key",
			})
			c.Abort()
			return
		}

		if !tb.Allow(key) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded",
				"message": "Too many requests, please try again later",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
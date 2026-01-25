package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)


type RateLimiter struct {
	limits map[string]*rate.Limiter
	mutex  sync.RWMutex
	rate   rate.Limit
	burst  int
}


func NewRateLimiter(requestsPerSecond float64, burst int) *RateLimiter {
	return &RateLimiter{
		limits: make(map[string]*rate.Limiter),
		rate:   rate.Limit(requestsPerSecond),
		burst:  burst,
	}
}


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


func (rl *RateLimiter) RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		
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


type TokenBucketRateLimiter struct {
	capacity  int
	refillRate time.Duration
	buckets   map[string]*TokenBucket
	mutex     sync.RWMutex
}


type TokenBucket struct {
	tokens    int
	lastRefill time.Time
	capacity  int
	refillRate time.Duration
}


func NewTokenBucket(capacity int, refillRate time.Duration) *TokenBucket {
	return &TokenBucket{
		tokens:     capacity,
		lastRefill: time.Now(),
		capacity:   capacity,
		refillRate: refillRate,
	}
}


func NewTokenBucketRateLimiter(capacity int, refillRate time.Duration) *TokenBucketRateLimiter {
	return &TokenBucketRateLimiter{
		capacity:  capacity,
		refillRate: refillRate,
		buckets:   make(map[string]*TokenBucket),
	}
}


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


func (tb *TokenBucket) consume(tokens int) bool {
	tb.refill()
	
	if tb.tokens >= tokens {
		tb.tokens -= tokens
		return true
	}
	
	return false
}


func (tb *TokenBucket) refill() {
	now := time.Now()
	elapsed := now.Sub(tb.lastRefill)
	
	if elapsed >= tb.refillRate {
		
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
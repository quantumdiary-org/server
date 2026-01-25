package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"netschool-proxy/api/api/internal/domain/cache"
)


type ResponseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r *ResponseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

type CacheMiddleware struct {
	cache cache.CacheStrategy
}

func NewCacheMiddleware(cache cache.CacheStrategy) *CacheMiddleware {
	return &CacheMiddleware{cache: cache}
}


func (cm *CacheMiddleware) CacheResponse(ttl time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		
		cacheKey := generateCacheKey(c)

		
		var cachedResponse map[string]interface{}
		found, err := cm.cache.Get(c.Request.Context(), cacheKey, &cachedResponse)
		if err == nil && found {
			
			c.JSON(200, cachedResponse)
			c.Abort()
			return
		}

		
		bodyWriter := &ResponseBodyWriter{
			ResponseWriter: c.Writer,
			body:          &bytes.Buffer{},
		}
		c.Writer = bodyWriter

		
		c.Next()

		
		if c.Writer.Status() == 200 && bodyWriter.body.Len() > 0 {
			
			var responseBody interface{}
			if err := json.Unmarshal(bodyWriter.body.Bytes(), &responseBody); err == nil {
				
				go func() {
					ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
					defer cancel()
					
					cm.cache.Set(ctx, cacheKey, responseBody, ttl)
				}()
			}
		}
	}
}


func generateCacheKey(c *gin.Context) string {
	return fmt.Sprintf("%s:%s?%s", c.Request.Method, c.Request.URL.Path, c.Request.URL.RawQuery)
}


func (cm *CacheMiddleware) CacheMiddlewareWithCondition(ttl time.Duration, conditionFunc func(*gin.Context) bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		
		if !conditionFunc(c) {
			c.Next()
			return
		}

		
		cacheKey := generateCacheKey(c)

		
		var cachedResponse map[string]interface{}
		found, err := cm.cache.Get(c.Request.Context(), cacheKey, &cachedResponse)
		if err == nil && found {
			
			c.JSON(200, cachedResponse)
			c.Abort()
			return
		}

		
		bodyWriter := &ResponseBodyWriter{
			ResponseWriter: c.Writer,
			body:          &bytes.Buffer{},
		}
		c.Writer = bodyWriter

		
		c.Next()

		
		if c.Writer.Status() == 200 && bodyWriter.body.Len() > 0 {
			
			var responseBody interface{}
			if err := json.Unmarshal(bodyWriter.body.Bytes(), &responseBody); err == nil {
				
				go func() {
					ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
					defer cancel()
					
					cm.cache.Set(ctx, cacheKey, responseBody, ttl)
				}()
			}
		}
	}
}
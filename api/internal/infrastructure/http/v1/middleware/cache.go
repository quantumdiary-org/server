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

// ResponseBodyWriter wraps gin.ResponseWriter to capture response body
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

// CacheResponse кэширует ответ обработчика
func (cm *CacheMiddleware) CacheResponse(ttl time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Генерируем ключ кэша на основе пути и параметров запроса
		cacheKey := generateCacheKey(c)

		// Пытаемся получить закэшированный ответ
		var cachedResponse map[string]interface{}
		found, err := cm.cache.Get(c.Request.Context(), cacheKey, &cachedResponse)
		if err == nil && found {
			// Если нашли закэшированный ответ, возвращаем его
			c.JSON(200, cachedResponse)
			c.Abort()
			return
		}

		// Если в кэше нет данных, оборачиваем ResponseWriter для захвата тела ответа
		bodyWriter := &ResponseBodyWriter{
			ResponseWriter: c.Writer,
			body:          &bytes.Buffer{},
		}
		c.Writer = bodyWriter

		// Продолжаем выполнение
		c.Next()

		// После выполнения обработчика сохраняем результат в кэш
		if c.Writer.Status() == 200 && bodyWriter.body.Len() > 0 {
			// Парсим тело ответа
			var responseBody interface{}
			if err := json.Unmarshal(bodyWriter.body.Bytes(), &responseBody); err == nil {
				// Сохраняем в кэш
				go func() {
					ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
					defer cancel()
					
					cm.cache.Set(ctx, cacheKey, responseBody, ttl)
				}()
			}
		}
	}
}

// generateCacheKey генерирует уникальный ключ для кэширования
func generateCacheKey(c *gin.Context) string {
	return fmt.Sprintf("%s:%s?%s", c.Request.Method, c.Request.URL.Path, c.Request.URL.RawQuery)
}

// CacheMiddlewareWithCondition позволяет кэшировать ответы с условием
func (cm *CacheMiddleware) CacheMiddlewareWithCondition(ttl time.Duration, conditionFunc func(*gin.Context) bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Проверяем условие кэширования
		if !conditionFunc(c) {
			c.Next()
			return
		}

		// Генерируем ключ кэша на основе пути и параметров запроса
		cacheKey := generateCacheKey(c)

		// Пытаемся получить закэшированный ответ
		var cachedResponse map[string]interface{}
		found, err := cm.cache.Get(c.Request.Context(), cacheKey, &cachedResponse)
		if err == nil && found {
			// Если нашли закэшированный ответ, возвращаем его
			c.JSON(200, cachedResponse)
			c.Abort()
			return
		}

		// Если в кэше нет данных, оборачиваем ResponseWriter для захвата тела ответа
		bodyWriter := &ResponseBodyWriter{
			ResponseWriter: c.Writer,
			body:          &bytes.Buffer{},
		}
		c.Writer = bodyWriter

		// Продолжаем выполнение
		c.Next()

		// После выполнения обработчика сохраняем результат в кэш
		if c.Writer.Status() == 200 && bodyWriter.body.Len() > 0 {
			// Парсим тело ответа
			var responseBody interface{}
			if err := json.Unmarshal(bodyWriter.body.Bytes(), &responseBody); err == nil {
				// Сохраняем в кэш
				go func() {
					ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
					defer cancel()
					
					cm.cache.Set(ctx, cacheKey, responseBody, ttl)
				}()
			}
		}
	}
}
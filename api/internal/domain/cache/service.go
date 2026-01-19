package cache

import (
	"context"
	"encoding/json"
	"time"
)

// Repository определяет интерфейс для работы с кэшем в базе данных
type Repository interface {
	Set(ctx context.Context, key string, value string, ttl time.Duration) error
	Get(ctx context.Context, key string) (string, bool, error)
	Delete(ctx context.Context, key string) error
	CleanupExpired(ctx context.Context) error
}

// CacheService предоставляет методы для работы с кэшем
type CacheService struct {
	repo Repository
}

// NewCacheService создает новый сервис кэширования
func NewCacheService(repo Repository) *CacheService {
	return &CacheService{repo: repo}
}

// Set устанавливает значение в кэш
func (c *CacheService) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	// Сериализуем значение в JSON
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return c.repo.Set(ctx, key, string(jsonValue), ttl)
}

// Get получает значение из кэша
func (c *CacheService) Get(ctx context.Context, key string, target interface{}) (bool, error) {
	jsonValue, exists, err := c.repo.Get(ctx, key)
	if err != nil || !exists {
		return exists, err
	}

	// Десериализуем JSON в целевую структуру
	if err := json.Unmarshal([]byte(jsonValue), target); err != nil {
		return false, err
	}

	return true, nil
}

// GetOrDefault получает значение из кэша или возвращает значение по умолчанию
func (c *CacheService) GetOrDefault(ctx context.Context, key string, defaultValue interface{}) (interface{}, error) {
	var result interface{}
	exists, err := c.Get(ctx, key, &result)
	if err != nil {
		return defaultValue, err
	}
	if !exists {
		return defaultValue, nil
	}
	return result, nil
}

// Delete удаляет значение из кэша
func (c *CacheService) Delete(ctx context.Context, key string) error {
	return c.repo.Delete(ctx, key)
}

// Cleanup удаляет просроченные записи из кэша
func (c *CacheService) Cleanup(ctx context.Context) error {
	return c.repo.CleanupExpired(ctx)
}
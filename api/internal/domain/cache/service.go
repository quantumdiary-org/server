package cache

import (
	"context"
	"encoding/json"
	"time"
)


type Repository interface {
	Set(ctx context.Context, key string, value string, ttl time.Duration) error
	Get(ctx context.Context, key string) (string, bool, error)
	Delete(ctx context.Context, key string) error
	CleanupExpired(ctx context.Context) error
}


type CacheService struct {
	repo Repository
}


func NewCacheService(repo Repository) *CacheService {
	return &CacheService{repo: repo}
}


func (c *CacheService) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return c.repo.Set(ctx, key, string(jsonValue), ttl)
}


func (c *CacheService) Get(ctx context.Context, key string, target interface{}) (bool, error) {
	jsonValue, exists, err := c.repo.Get(ctx, key)
	if err != nil || !exists {
		return exists, err
	}

	
	if err := json.Unmarshal([]byte(jsonValue), target); err != nil {
		return false, err
	}

	return true, nil
}


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


func (c *CacheService) Delete(ctx context.Context, key string) error {
	return c.repo.Delete(ctx, key)
}


func (c *CacheService) Cleanup(ctx context.Context) error {
	return c.repo.CleanupExpired(ctx)
}
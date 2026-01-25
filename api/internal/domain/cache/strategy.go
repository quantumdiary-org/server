package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
)

type CacheStrategy interface {
	Get(ctx context.Context, key string, target interface{}) (bool, error)
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
	Clear(ctx context.Context) error
}

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(client *redis.Client) *RedisCache {
	return &RedisCache{client: client}
}

func (r *RedisCache) Get(ctx context.Context, key string, target interface{}) (bool, error) {
	data, err := r.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, json.Unmarshal(data, target)
}

func (r *RedisCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, key, data, ttl).Err()
}

func (r *RedisCache) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

func (r *RedisCache) Clear(ctx context.Context) error {
	return r.client.FlushDB(ctx).Err()
}


type CacheItem struct {
	Value     interface{}
	ExpiresAt time.Time
}


func (c *CacheItem) IsExpired() bool {
	return c.ExpiresAt.Before(time.Now())
}


type MemoryCache struct {
	data map[string]*CacheItem
}

func NewMemoryCache() *MemoryCache {
	return &MemoryCache{
		data: make(map[string]*CacheItem),
	}
}

func (m *MemoryCache) Get(ctx context.Context, key string, target interface{}) (bool, error) {
	item, exists := m.data[key]
	if !exists || item.IsExpired() {
		delete(m.data, key) 
		return false, nil
	}

	data, ok := item.Value.([]byte)
	if !ok {
		return false, nil
	}

	return true, json.Unmarshal(data, target)
}

func (m *MemoryCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	m.data[key] = &CacheItem{
		Value:     data,
		ExpiresAt: time.Now().Add(ttl),
	}

	return nil
}

func (m *MemoryCache) Delete(ctx context.Context, key string) error {
	delete(m.data, key)
	return nil
}

func (m *MemoryCache) Clear(ctx context.Context) error {
	m.data = make(map[string]*CacheItem)
	return nil
}


type CacheWithFallback struct {
	primary  CacheStrategy
	fallback CacheStrategy
}

func NewCacheWithFallback(primary, fallback CacheStrategy) *CacheWithFallback {
	return &CacheWithFallback{
		primary:  primary,
		fallback: fallback,
	}
}

func (c *CacheWithFallback) Get(ctx context.Context, key string, target interface{}) (bool, error) {
	found, err := c.primary.Get(ctx, key, target)
	if err == nil && found {
		return true, nil
	}

	
	return c.fallback.Get(ctx, key, target)
}

func (c *CacheWithFallback) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	
	err1 := c.primary.Set(ctx, key, value, ttl)
	err2 := c.fallback.Set(ctx, key, value, ttl)
	
	if err1 != nil {
		return err1
	}
	return err2
}

func (c *CacheWithFallback) Delete(ctx context.Context, key string) error {
	
	err1 := c.primary.Delete(ctx, key)
	err2 := c.fallback.Delete(ctx, key)
	
	if err1 != nil {
		return err1
	}
	return err2
}

func (c *CacheWithFallback) Clear(ctx context.Context) error {
	
	err1 := c.primary.Clear(ctx)
	err2 := c.fallback.Clear(ctx)
	
	if err1 != nil {
		return err1
	}
	return err2
}
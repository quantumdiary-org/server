package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
)


type RedisCacheService struct {
	client *redis.Client
}


func NewRedisCacheService(redisAddr string) (*RedisCacheService, error) {
	client := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &RedisCacheService{client: client}, nil
}


func (r *RedisCacheService) Get(ctx context.Context, key string, target interface{}) (bool, error) {
	data, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, json.Unmarshal([]byte(data), target)
}


func (r *RedisCacheService) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, key, data, ttl).Err()
}


func (r *RedisCacheService) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}


func (r *RedisCacheService) Clear(ctx context.Context) error {
	return r.client.FlushDB(ctx).Err()
}


func (r *RedisCacheService) Close() error {
	return r.client.Close()
}


func (r *RedisCacheService) Exists(ctx context.Context, key string) (bool, error) {
	count, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return count > 0, nil
}


func (r *RedisCacheService) Expire(ctx context.Context, key string, ttl time.Duration) error {
	return r.client.Expire(ctx, key, ttl).Err()
}
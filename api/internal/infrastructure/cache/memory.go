package cache

import (
	"context"
	"encoding/json"
	"sync"
	"time"
)


type MemoryCacheService struct {
	data map[string]*CacheItem
	mu   sync.RWMutex
	maxSize int
}


type CacheItem struct {
	Value     interface{}
	ExpiresAt time.Time
}


func (c *CacheItem) IsExpired() bool {
	return time.Now().After(c.ExpiresAt)
}


func NewMemoryCacheService(maxSize int) *MemoryCacheService {
	return &MemoryCacheService{
		data:    make(map[string]*CacheItem),
		maxSize: maxSize,
	}
}


func (m *MemoryCacheService) Get(ctx context.Context, key string, target interface{}) (bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	item, exists := m.data[key]
	if !exists || item.IsExpired() {
		if exists {
			
			delete(m.data, key)
		}
		return false, nil
	}

	data, ok := item.Value.([]byte)
	if !ok {
		return false, nil
	}

	return true, json.Unmarshal(data, target)
}


func (m *MemoryCacheService) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	
	if len(m.data) >= m.maxSize {
		
		
		for k, v := range m.data {
			if v.IsExpired() {
				delete(m.data, k)
			}
		}
	}

	m.data[key] = &CacheItem{
		Value:     data,
		ExpiresAt: time.Now().Add(ttl),
	}

	return nil
}


func (m *MemoryCacheService) Delete(ctx context.Context, key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.data, key)
	return nil
}


func (m *MemoryCacheService) Clear(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.data = make(map[string]*CacheItem)
	return nil
}


func (m *MemoryCacheService) Exists(ctx context.Context, key string) (bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	item, exists := m.data[key]
	if exists && item.IsExpired() {
		delete(m.data, key)
		return false, nil
	}

	return exists, nil
}


func (m *MemoryCacheService) Expire(ctx context.Context, key string, ttl time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	item, exists := m.data[key]
	if !exists || item.IsExpired() {
		if exists {
			delete(m.data, key)
		}
		return nil
	}

	item.ExpiresAt = time.Now().Add(ttl)
	return nil
}


func (m *MemoryCacheService) Size() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	count := 0
	for _, item := range m.data {
		if !item.IsExpired() {
			count++
		} else {
			
			delete(m.data, "") 
		}
	}
	return count
}


func (m *MemoryCacheService) Cleanup() {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	for key, item := range m.data {
		if item.IsExpired() {
			delete(m.data, key)
		}
	}
}
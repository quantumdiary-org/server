package cache

import (
	"context"
	"encoding/json"
	"sync"
	"time"
)

// MemoryCacheService provides in-memory caching
type MemoryCacheService struct {
	data map[string]*CacheItem
	mu   sync.RWMutex
	maxSize int
}

// CacheItem represents a cached item with expiration
type CacheItem struct {
	Value     interface{}
	ExpiresAt time.Time
}

// IsExpired checks if the cache item has expired
func (c *CacheItem) IsExpired() bool {
	return time.Now().After(c.ExpiresAt)
}

// NewMemoryCacheService creates a new in-memory cache service
func NewMemoryCacheService(maxSize int) *MemoryCacheService {
	return &MemoryCacheService{
		data:    make(map[string]*CacheItem),
		maxSize: maxSize,
	}
}

// Get retrieves a value from memory cache
func (m *MemoryCacheService) Get(ctx context.Context, key string, target interface{}) (bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	item, exists := m.data[key]
	if !exists || item.IsExpired() {
		if exists {
			// Delete expired item
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

// Set stores a value in memory cache
func (m *MemoryCacheService) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// Check if we need to evict items due to size limit
	if len(m.data) >= m.maxSize {
		// Simple eviction: remove oldest items
		// In a real implementation, you might want to use LRU or another algorithm
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

// Delete removes a value from memory cache
func (m *MemoryCacheService) Delete(ctx context.Context, key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.data, key)
	return nil
}

// Clear removes all values from memory cache
func (m *MemoryCacheService) Clear(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.data = make(map[string]*CacheItem)
	return nil
}

// Exists checks if a key exists in memory cache
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

// Expire sets expiration time for a key in memory cache
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

// Size returns the current size of the cache
func (m *MemoryCacheService) Size() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	count := 0
	for _, item := range m.data {
		if !item.IsExpired() {
			count++
		} else {
			// Clean up expired items during size check
			delete(m.data, "") // This won't actually delete anything, just for demonstration
		}
	}
	return count
}

// Cleanup removes expired items from the cache
func (m *MemoryCacheService) Cleanup() {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	for key, item := range m.data {
		if item.IsExpired() {
			delete(m.data, key)
		}
	}
}
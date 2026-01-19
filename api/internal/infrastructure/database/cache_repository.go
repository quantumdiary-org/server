package database

import (
	"context"
	"time"

	"gorm.io/gorm"
	"netschool-proxy/api/api/internal/domain/cache"
)

type CacheRepository struct {
	db *gorm.DB
}

func NewCacheRepository(db *gorm.DB) *CacheRepository {
	return &CacheRepository{db: db}
}

func (r *CacheRepository) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	entry := &cache.CacheEntry{
		Key:       key,
		Value:     value,
		ExpiresAt: time.Now().Add(ttl),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Проверяем, существует ли уже запись с таким ключом
	var existingEntry cache.CacheEntry
	result := r.db.Where("key = ?", key).First(&existingEntry)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			// Запись не существует, создаем новую
			return r.db.WithContext(ctx).Create(entry).Error
		}
		return result.Error
	}

	// Запись существует, обновляем её
	entry.ID = existingEntry.ID // Сохраняем ID
	return r.db.WithContext(ctx).Save(entry).Error
}

func (r *CacheRepository) Get(ctx context.Context, key string) (string, bool, error) {
	var entry cache.CacheEntry
	result := r.db.WithContext(ctx).
		Where("key = ? AND expires_at > ?", key, time.Now()).
		First(&entry)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return "", false, nil
		}
		return "", false, result.Error
	}

	return entry.Value, true, nil
}

func (r *CacheRepository) Delete(ctx context.Context, key string) error {
	return r.db.WithContext(ctx).
		Where("key = ?", key).
		Delete(&cache.CacheEntry{}).Error
}

func (r *CacheRepository) CleanupExpired(ctx context.Context) error {
	return r.db.WithContext(ctx).
		Where("expires_at < ?", time.Now()).
		Delete(&cache.CacheEntry{}).Error
}
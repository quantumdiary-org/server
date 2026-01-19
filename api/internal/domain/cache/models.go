package cache

import (
	"time"
)

// CacheEntry представляет запись в кэше
type CacheEntry struct {
	ID        int       `json:"id" gorm:"column:id"`
	Key       string    `json:"key" gorm:"column:key"`
	Value     string    `json:"-" gorm:"column:value"` // Значение не передается в JSON, только хранится в БД
	ExpiresAt time.Time `json:"expires_at" gorm:"column:expires_at"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at"`
}

// TableName указывает имя таблицы в базе данных
func (CacheEntry) TableName() string {
	return "cache"
}


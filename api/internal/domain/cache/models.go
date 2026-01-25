package cache

import (
	"time"
)


type CacheEntry struct {
	ID        int       `json:"id" gorm:"column:id"`
	Key       string    `json:"key" gorm:"column:key"`
	Value     string    `json:"-" gorm:"column:value"` 
	ExpiresAt time.Time `json:"expires_at" gorm:"column:expires_at"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at"`
}


func (CacheEntry) TableName() string {
	return "cache"
}


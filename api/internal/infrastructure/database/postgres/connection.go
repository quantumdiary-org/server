package postgres

import (
	"fmt"

	"netschool-proxy/api/api/internal/domain/auth"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)


func NewConnection(url string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(url), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	
	if err := db.AutoMigrate(&auth.NetSchoolSession{}); err != nil {
		return nil, fmt.Errorf("failed to migrate database schema: %w", err)
	}

	return db, nil
}
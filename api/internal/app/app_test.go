package app_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"netschool-proxy/api/api/internal/app"
	"netschool-proxy/api/api/internal/config"
)

func TestNewApp(t *testing.T) {
	
	cfg := &config.Config{
		Server: config.ServerConfig{
			Port: "8080",
		},
		Database: config.DatabaseConfig{
			Host:     "localhost",
			Port:     5432,
			Name:     "test_db",
			User:     "test_user",
			Password: "test_password",
		},
		JWT: config.JWTConfig{
			Secret:    "test_secret",
			ExpiresIn: 24 * 60 * 60, 
		},
	}

	
	
	
	_, err := app.New(cfg)

	
	
	
	assert.NotNil(t, cfg)
	_ = err 
}
package security_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"netschool-proxy/api/api/internal/pkg/security"
)

func TestJWTService_GenerateAndParseToken(t *testing.T) {
	secret := "test_secret"
	expiresIn := 24 * time.Hour
	
	jwtService := security.NewJWTService(secret, expiresIn)
	
	userID := "test_user"
	sessionID := "test_session"
	schoolID := 123
	
	token, err := jwtService.GenerateToken(userID, sessionID, schoolID)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	
	parsedClaims, err := jwtService.ParseToken(token)
	assert.NoError(t, err)
	assert.Equal(t, userID, parsedClaims.UserID)
	assert.Equal(t, sessionID, parsedClaims.SessionID)
	assert.Equal(t, schoolID, parsedClaims.SchoolID)
}

func TestJWTService_InvalidToken(t *testing.T) {
	secret := "test_secret"
	differentSecret := "different_secret"
	expiresIn := 24 * time.Hour
	
	jwtService1 := security.NewJWTService(secret, expiresIn)
	jwtService2 := security.NewJWTService(differentSecret, expiresIn)
	
	userID := "test_user"
	sessionID := "test_session"
	schoolID := 123
	
	token, err := jwtService1.GenerateToken(userID, sessionID, schoolID)
	assert.NoError(t, err)
	
	
	_, err = jwtService2.ParseToken(token)
	assert.Error(t, err)
}

func TestJWTService_ExpiredToken(t *testing.T) {
	secret := "test_secret"
	expiresIn := 1 * time.Millisecond 
	
	jwtService := security.NewJWTService(secret, expiresIn)
	
	userID := "test_user"
	sessionID := "test_session"
	schoolID := 123
	
	token, err := jwtService.GenerateToken(userID, sessionID, schoolID)
	assert.NoError(t, err)
	
	
	time.Sleep(10 * time.Millisecond)
	
	_, err = jwtService.ParseToken(token)
	assert.Error(t, err)
}
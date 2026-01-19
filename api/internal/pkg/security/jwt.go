package security

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService struct {
	secretKey []byte
	expiresIn time.Duration
}

type Claims struct {
	UserID    string `json:"user_id"`
	SessionID string `json:"session_id"`
	SchoolID  int    `json:"school_id"`
	jwt.RegisteredClaims
}

func NewJWTService(secret string, expiresIn time.Duration) *JWTService {
	return &JWTService{
		secretKey: []byte(secret),
		expiresIn: expiresIn,
	}
}

func (s *JWTService) GenerateToken(userID, sessionID string, schoolID int) (string, error) {
	expiresAt := time.Now().Add(s.expiresIn)

	claims := &Claims{
		UserID:    userID,
		SessionID: sessionID,
		SchoolID:  schoolID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "netschool-proxy/api",
			Subject:   userID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secretKey)
}

func (s *JWTService) ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token claims")
}
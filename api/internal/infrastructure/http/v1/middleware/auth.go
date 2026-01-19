package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"netschool-proxy/api/api/internal/domain/auth"
	"netschool-proxy/api/api/internal/pkg/security"
)

type AuthMiddleware struct {
	authService *auth.Service
	jwtService  *security.JWTService
}

func NewAuthMiddleware(authService *auth.Service, jwtService *security.JWTService) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
		jwtService:  jwtService,
	}
}

func (m *AuthMiddleware) AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(401, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Проверяем формат заголовка Authorization
		authParts := strings.SplitN(authHeader, " ", 2)
		if len(authParts) != 2 || authParts[0] != "Bearer" {
			c.JSON(401, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		tokenString := authParts[1]

		// Валидируем токен
		claims, err := m.authService.ValidateToken(c.Request.Context(), tokenString)
		if err != nil {
			c.JSON(401, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Получаем сессию отдельно для дальнейшего использования
		session, err := m.authService.GetSessionByUserID(c.Request.Context(), claims.UserID)
		if err != nil {
			c.JSON(401, gin.H{"error": "Session not found"})
			c.Abort()
			return
		}

		// Добавляем информацию о пользователе в контекст
		c.Set("userID", claims.UserID)
		c.Set("sessionID", claims.SessionID)
		c.Set("schoolID", claims.SchoolID)
		c.Set("session", session) // Добавляем сессию в контекст для дальнейшего использования

		c.Next()
	}
}
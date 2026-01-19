package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"netschool-proxy/api/api/internal/api_types"
	"netschool-proxy/api/api/internal/pkg/security"
)

type Service struct {
	sessionRepo    SessionRepository
	apiClientFactory *api_types.APIClientFactory
	config         api_types.APIConfig
	jwtService     *security.JWTService
}

type SessionRepository interface {
	Create(ctx context.Context, session *NetSchoolSession) error
	GetByUserID(ctx context.Context, userID string) (*NetSchoolSession, error)
	Delete(ctx context.Context, userID string) error
	CleanupExpired(ctx context.Context) error
}

func NewService(sessionRepo SessionRepository, apiClientFactory *api_types.APIClientFactory, config api_types.APIConfig, jwtService *security.JWTService) *Service {
	return &Service{
		sessionRepo:      sessionRepo,
		apiClientFactory: apiClientFactory,
		config:           config,
		jwtService:       jwtService,
	}
}

func (s *Service) Login(ctx context.Context, username, password string, schoolID int, instanceURL string) (string, error) {
	// Используем конфигурацию по умолчанию
	return s.LoginWithAPIType(ctx, username, password, schoolID, instanceURL, string(s.config.Mode))
}

func (s *Service) LoginWithAPIType(ctx context.Context, username, password string, schoolID int, instanceURL string, apiType string) (string, error) {
	// 1. Создаем клиент API соответствующего типа
	apiMode := api_types.APIMode(apiType)
	clientConfig := s.config
	clientConfig.Mode = apiMode

	apiClient, err := s.apiClientFactory.NewAPIClient(apiMode, clientConfig)
	if err != nil {
		return "", fmt.Errorf("failed to create API client: %w", err)
	}

	// 2. Получаем loginData от API
	loginData, err := apiClient.GetLoginData(ctx, instanceURL)
	if err != nil {
		return "", fmt.Errorf("failed to get login data from API: %w", err)
	}

	// 3. Выполняем аутентификацию через API
	accessToken, err := apiClient.Login(ctx, username, password, schoolID, instanceURL, loginData)
	if err != nil {
		return "", fmt.Errorf("failed to authenticate with API: %w", err)
	}

	// 4. Генерируем уникальный ID пользователя (например, на основе логина и школы)
	userID := fmt.Sprintf("%s_%d_%s", username, schoolID, apiType)

	// 5. Получаем дополнительную информацию о пользователе из API
	// Попробуем получить информацию о пользователе для извлечения studentID
	var studentID string
	var yearID string

	// Получаем информацию о пользователе
	userInfo, err := apiClient.GetInfo(ctx, accessToken, instanceURL)
	if err != nil {
		// Если не удается получить информацию, используем дефолтные значения
		studentID = fmt.Sprintf("user_%s_%d", username, schoolID)
	} else {
		// Извлекаем studentID из информации о пользователе
		if userInfoMap, ok := userInfo.(map[string]interface{}); ok {
			if id, exists := userInfoMap["id"]; exists {
				studentID = fmt.Sprintf("%v", id)
			} else {
				studentID = fmt.Sprintf("user_%s_%d", username, schoolID)
			}
		} else {
			studentID = fmt.Sprintf("user_%s_%d", username, schoolID)
		}
	}

	// Получаем информацию о контексте для извлечения года
	// Используем GetSchoolInfo для получения информации о школе и году
	schoolInfo, err := apiClient.GetSchoolInfo(ctx, accessToken, instanceURL)
	if err != nil {
		// Если не удается получить информацию о школе, используем дефолтное значение
		yearID = "current_year"
	} else {
		// Извлекаем год из информации о школе
		if schoolInfoMap, ok := schoolInfo.(map[string]interface{}); ok {
			if year, exists := schoolInfoMap["yearId"]; exists {
				yearID = fmt.Sprintf("%v", year)
			} else if id, exists := schoolInfoMap["id"]; exists {
				yearID = fmt.Sprintf("%v", id)
			} else {
				yearID = "current_year"
			}
		} else {
			yearID = "current_year"
		}
	}

	// 6. Сохраняем сессию в базе данных (с токеном Сетевого Города)
	session := &NetSchoolSession{
		UserID:               userID,
		NetSchoolAccessToken: accessToken, // Токен Сетевого Города хранится только в БД
		RefreshToken:         "", // В реальной системе может быть refresh token
		ExpiresAt:            time.Now().Add(24 * time.Hour), // Используем фиксированный срок действия
		NetSchoolURL:         instanceURL,                     // Используем переданный URL
		SchoolID:             schoolID,
		StudentID:            studentID,
		YearID:               yearID,
		APIType:              apiType,                         // Сохраняем тип API
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}

	if err := s.sessionRepo.Create(ctx, session); err != nil {
		return "", fmt.Errorf("failed to save session: %w", err)
	}

	// 7. Генерируем JWT токен для клиента (токен прокси, не токен Сетевого Города)
	proxyToken, err := s.jwtService.GenerateToken(userID, fmt.Sprintf("%d", session.ID), schoolID)
	if err != nil {
		return "", fmt.Errorf("failed to generate proxy token: %w", err)
	}

	return proxyToken, nil
}

func (s *Service) ValidateToken(ctx context.Context, token string) (*security.Claims, error) {
	claims, err := s.jwtService.ParseToken(token)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	// Проверяем, что сессия все еще действительна в базе данных
	session, err := s.sessionRepo.GetByUserID(ctx, claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("session not found: %w", err)
	}

	// Проверяем, что сессия не истекла
	if session.ExpiresAt.Before(time.Now()) {
		// Удаляем просроченную сессию
		s.sessionRepo.Delete(ctx, claims.UserID)
		return nil, errors.New("session expired")
	}

	return claims, nil
}

// ValidateTokenWithSession валидирует токен и возвращает сессию
func (s *Service) ValidateTokenWithSession(ctx context.Context, token string) (*security.Claims, *NetSchoolSession, error) {
	claims, err := s.jwtService.ParseToken(token)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid token: %w", err)
	}

	// Проверяем, что сессия все еще действительна в базе данных
	session, err := s.sessionRepo.GetByUserID(ctx, claims.UserID)
	if err != nil {
		return nil, nil, fmt.Errorf("session not found: %w", err)
	}

	// Проверяем, что сессия не истекла
	if session.ExpiresAt.Before(time.Now()) {
		// Удаляем просроченную сессию
		s.sessionRepo.Delete(ctx, claims.UserID)
		return nil, nil, errors.New("session expired")
	}

	return claims, session, nil
}

func (s *Service) GetSessionByUserID(ctx context.Context, userID string) (*NetSchoolSession, error) {
	session, err := s.sessionRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}
	return session, nil
}

func (s *Service) Logout(ctx context.Context, userID string) error {
	return s.sessionRepo.Delete(ctx, userID)
}
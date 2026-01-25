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
	
	return s.LoginWithAPIType(ctx, username, password, schoolID, instanceURL, string(s.config.Mode))
}

func (s *Service) LoginWithAPIType(ctx context.Context, username, password string, schoolID int, instanceURL string, apiType string) (string, error) {
	
	apiMode := api_types.APIMode(apiType)
	clientConfig := s.config
	clientConfig.Mode = apiMode

	apiClient, err := s.apiClientFactory.NewAPIClient(apiMode, clientConfig)
	if err != nil {
		return "", fmt.Errorf("failed to create API client: %w", err)
	}

	
	loginData, err := apiClient.GetLoginData(ctx, instanceURL)
	if err != nil {
		return "", fmt.Errorf("failed to get login data from API: %w", err)
	}

	
	accessToken, err := apiClient.Login(ctx, username, password, schoolID, instanceURL, loginData)
	if err != nil {
		return "", fmt.Errorf("failed to authenticate with API: %w", err)
	}

	
	userID := fmt.Sprintf("%s_%d_%s", username, schoolID, apiType)

	
	
	var studentID string
	var yearID string

	
	userInfo, err := apiClient.GetInfo(ctx, accessToken, instanceURL)
	if err != nil {
		
		studentID = fmt.Sprintf("user_%s_%d", username, schoolID)
	} else {
		
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

	
	
	schoolInfo, err := apiClient.GetSchoolInfo(ctx, accessToken, instanceURL)
	if err != nil {
		
		yearID = "current_year"
	} else {
		
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

	
	session := &NetSchoolSession{
		UserID:               userID,
		NetSchoolAccessToken: accessToken, 
		RefreshToken:         "", 
		ExpiresAt:            time.Now().Add(24 * time.Hour), 
		NetSchoolURL:         instanceURL,                     
		SchoolID:             schoolID,
		StudentID:            studentID,
		YearID:               yearID,
		APIType:              apiType,                         
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}

	if err := s.sessionRepo.Create(ctx, session); err != nil {
		return "", fmt.Errorf("failed to save session: %w", err)
	}

	
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

	
	session, err := s.sessionRepo.GetByUserID(ctx, claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("session not found: %w", err)
	}

	
	if session.ExpiresAt.Before(time.Now()) {
		
		s.sessionRepo.Delete(ctx, claims.UserID)
		return nil, errors.New("session expired")
	}

	return claims, nil
}


func (s *Service) ValidateTokenWithSession(ctx context.Context, token string) (*security.Claims, *NetSchoolSession, error) {
	claims, err := s.jwtService.ParseToken(token)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid token: %w", err)
	}

	
	session, err := s.sessionRepo.GetByUserID(ctx, claims.UserID)
	if err != nil {
		return nil, nil, fmt.Errorf("session not found: %w", err)
	}

	
	if session.ExpiresAt.Before(time.Now()) {
		
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
package schedule

import (
	"context"
	"fmt"
	"time"

	"netschool-proxy/api/api/internal/api_types"
	"netschool-proxy/api/api/internal/domain/auth"
	"netschool-proxy/api/api/internal/domain/cache"
)

type Service struct {
	apiClientFactory *api_types.APIClientFactory
	sessionRepo      auth.SessionRepository
	cacheService     cache.CacheStrategy
	config           api_types.APIConfig
}

func NewService(apiClientFactory *api_types.APIClientFactory, sessionRepo auth.SessionRepository, cacheService cache.CacheStrategy, config api_types.APIConfig) *Service {
	return &Service{
		apiClientFactory: apiClientFactory,
		sessionRepo:      sessionRepo,
		cacheService:     cacheService,
		config:           config,
	}
}

func (s *Service) GetWeeklySchedule(ctx context.Context, userID, instanceURL string, weekStart time.Time) (interface{}, error) {
	
	cacheKey := fmt.Sprintf("schedule_weekly_%s_%s_%s", userID, instanceURL, weekStart.Format("2006-01-02"))
	var cachedSchedule interface{}

	if s.cacheService != nil {
		found, err := s.cacheService.Get(ctx, cacheKey, &cachedSchedule)
		if err == nil && found {
			return cachedSchedule, nil
		}
	}

	
	session, err := s.sessionRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user session: %w", err)
	}

	
	apiMode := api_types.APIMode(session.APIType)
	clientConfig := s.config
	clientConfig.Mode = apiMode

	apiClient, err := s.apiClientFactory.NewAPIClient(apiMode, clientConfig)
	if err != nil {
		
		if s.cacheService != nil {
			var backupSchedule interface{}
			_, err := s.cacheService.Get(ctx, cacheKey+"_backup", &backupSchedule)
			if err == nil {
				return backupSchedule, nil
			}
		}
		return nil, fmt.Errorf("failed to create API client: %w", err)
	}

	
	scheduleData, err := apiClient.GetSchedule(ctx, session.NetSchoolAccessToken, instanceURL, weekStart)
	if err != nil {
		
		if s.cacheService != nil {
			var backupSchedule interface{}
			_, cacheErr := s.cacheService.Get(ctx, cacheKey+"_backup", &backupSchedule)
			if cacheErr == nil {
				return backupSchedule, nil
			}
		}
		return nil, fmt.Errorf("failed to get schedule from API: %w", err)
	}

	
	if s.cacheService != nil {
		
		s.cacheService.Set(ctx, cacheKey, scheduleData, 30*time.Minute)

		
		s.cacheService.Set(ctx, cacheKey+"_backup", scheduleData, 24*time.Hour)
	}

	return scheduleData, nil
}

func (s *Service) GetDailySchedule(ctx context.Context, userID, instanceURL string, date time.Time) (interface{}, error) {
	
	cacheKey := fmt.Sprintf("schedule_daily_%s_%s_%s", userID, instanceURL, date.Format("2006-01-02"))
	var cachedSchedule interface{}

	if s.cacheService != nil {
		found, err := s.cacheService.Get(ctx, cacheKey, &cachedSchedule)
		if err == nil && found {
			return cachedSchedule, nil
		}
	}

	
	weeklySchedule, err := s.GetWeeklySchedule(ctx, userID, instanceURL, date.AddDate(0, 0, -int(date.Weekday())+1))
	if err != nil {
		return nil, fmt.Errorf("failed to get daily schedule: %w", err)
	}

	
	return weeklySchedule, nil
}
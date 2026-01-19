package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"netschool-proxy/api/api/internal/api_types"
	"netschool-proxy/api/api/internal/domain/auth"
)

// AutoCacheService автоматически обновляет кэш для всех аккаунтов
type AutoCacheService struct {
	cacheRepo      Repository
	sessionRepo    auth.SessionRepository
	apiClientFactory *api_types.APIClientFactory
	config         api_types.APIConfig
	interval       time.Duration
	stopChan       chan struct{}
	wg             sync.WaitGroup
}

// NewAutoCacheService создает новый сервис автоматического обновления кэша
func NewAutoCacheService(
	cacheRepo Repository,
	sessionRepo auth.SessionRepository,
	apiClientFactory *api_types.APIClientFactory,
	config api_types.APIConfig,
	interval time.Duration,
) *AutoCacheService {
	return &AutoCacheService{
		cacheRepo:      cacheRepo,
		sessionRepo:    sessionRepo,
		apiClientFactory: apiClientFactory,
		config:         config,
		interval:       interval,
		stopChan:       make(chan struct{}),
	}
}

// Start запускает автоматическое обновление кэша
func (s *AutoCacheService) Start(ctx context.Context) {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		
		ticker := time.NewTicker(s.interval)
		defer ticker.Stop()
		
		// Выполняем первое обновление сразу
		s.updateCache(ctx)
		
		for {
			select {
			case <-ticker.C:
				s.updateCache(ctx)
			case <-s.stopChan:
				return
			case <-ctx.Done():
				return
			}
		}
	}()
}

// Stop останавливает автоматическое обновление кэша
func (s *AutoCacheService) Stop() {
	close(s.stopChan)
	s.wg.Wait()
}

// updateCache обновляет кэш для всех активных аккаунтов
func (s *AutoCacheService) updateCache(ctx context.Context) {
	// Получаем все активные сессии
	sessions, err := s.getAllActiveSessions(ctx)
	if err != nil {
		log.Printf("Failed to get active sessions for cache update: %v", err)
		return
	}

	for _, session := range sessions {
		// Обновляем кэш для каждой сессии
		if err := s.updateSessionCache(ctx, session); err != nil {
			log.Printf("Failed to update cache for session %s: %v", session.UserID, err)
		}
	}
}

// getAllActiveSessions возвращает все активные сессии
func (s *AutoCacheService) getAllActiveSessions(ctx context.Context) ([]*auth.NetSchoolSession, error) {
	// В реальной реализации здесь будет запрос к базе данных
	// для получения всех активных сессий
	// Пока возвращаем пустой массив
	return []*auth.NetSchoolSession{}, nil
}

// updateSessionCache обновляет кэш для конкретной сессии
func (s *AutoCacheService) updateSessionCache(ctx context.Context, session *auth.NetSchoolSession) error {
	// Создаем клиент API соответствующего типа
	apiMode := api_types.APIMode(session.APIType)
	clientConfig := s.config
	clientConfig.Mode = apiMode

	apiClient, err := s.apiClientFactory.NewAPIClient(apiMode, clientConfig)
	if err != nil {
		return fmt.Errorf("failed to create API client: %w", err)
	}

	// Обновляем кэш для текущей недели
	if err := s.updateWeeklyCache(ctx, apiClient, session); err != nil {
		return fmt.Errorf("failed to update weekly cache: %w", err)
	}

	// Обновляем кэш итоговых оценок
	if err := s.updateFinalGradesCache(ctx, apiClient, session); err != nil {
		return fmt.Errorf("failed to update final grades cache: %w", err)
	}

	return nil
}

// updateWeeklyCache обновляет кэш расписания на неделю
func (s *AutoCacheService) updateWeeklyCache(ctx context.Context, apiClient api_types.APIClientInterface, session *auth.NetSchoolSession) error {
	// Получаем расписание на неделю
	weekStart := time.Now()
	scheduleData, err := apiClient.GetSchedule(ctx, session.NetSchoolAccessToken, session.NetSchoolURL, weekStart)
	if err != nil {
		// Если API недоступен, сохраняем текущие данные из кэша как резервные
		return s.saveBackupCache(ctx, session, "schedule_"+weekStart.Format("2006-01-02"))
	}

	// Сохраняем в кэш
	cacheKey := fmt.Sprintf("schedule_week_%s_%s", session.UserID, weekStart.Format("2006-01-02"))
	jsonData, err := json.Marshal(scheduleData)
	if err != nil {
		return fmt.Errorf("failed to marshal schedule data: %w", err)
	}

	if err := s.cacheRepo.Set(ctx, cacheKey, string(jsonData), 24*time.Hour); err != nil {
		return fmt.Errorf("failed to save schedule to cache: %w", err)
	}

	return nil
}

// updateFinalGradesCache обновляет кэш итоговых оценок
func (s *AutoCacheService) updateFinalGradesCache(ctx context.Context, apiClient api_types.APIClientInterface, session *auth.NetSchoolSession) error {
	// Получаем итоговые оценки
	gradesData, err := apiClient.GetGrades(ctx, session.NetSchoolAccessToken, session.StudentID, session.NetSchoolURL)
	if err != nil {
		// Если API недоступен, сохраняем текущие данные из кэша как резервные
		return s.saveBackupCache(ctx, session, "final_grades")
	}

	// Сохраняем в кэш
	cacheKey := fmt.Sprintf("final_grades_%s", session.UserID)
	jsonData, err := json.Marshal(gradesData)
	if err != nil {
		return fmt.Errorf("failed to marshal grades data: %w", err)
	}

	if err := s.cacheRepo.Set(ctx, cacheKey, string(jsonData), 24*time.Hour); err != nil {
		return fmt.Errorf("failed to save grades to cache: %w", err)
	}

	return nil
}

// saveBackupCache сохраняет резервную копию кэша
func (s *AutoCacheService) saveBackupCache(ctx context.Context, session *auth.NetSchoolSession, prefix string) error {
	// Получаем текущие данные из кэша
	backupKey := fmt.Sprintf("%s_backup_%s", prefix, session.UserID)
	currentKey := fmt.Sprintf("%s_%s", prefix, session.UserID)
	
	currentData, exists, err := s.cacheRepo.Get(ctx, currentKey)
	if err != nil || !exists {
		// Если нет текущих данных, ничего не делаем
		return nil
	}

	// Сохраняем как резервную копию
	if err := s.cacheRepo.Set(ctx, backupKey, currentData, 48*time.Hour); err != nil {
		return fmt.Errorf("failed to save backup cache: %w", err)
	}

	return nil
}

// GetCachedData получает данные из кэша или резервной копии
func (s *AutoCacheService) GetCachedData(ctx context.Context, userID, dataType string, weekStart time.Time) (interface{}, bool, error) {
	var cacheKey string
	if weekStart.IsZero() {
		cacheKey = fmt.Sprintf("%s_%s", dataType, userID)
	} else {
		cacheKey = fmt.Sprintf("%s_%s_%s", dataType, userID, weekStart.Format("2006-01-02"))
	}

	// Пробуем получить из основного кэша
	data, exists, err := s.cacheRepo.Get(ctx, cacheKey)
	if err != nil || !exists {
		// Пробуем получить из резервного кэша
		backupKey := fmt.Sprintf("%s_backup_%s", dataType, userID)
		data, exists, err = s.cacheRepo.Get(ctx, backupKey)
		if err != nil || !exists {
			return nil, false, err
		}
	}

	var result interface{}
	if err := json.Unmarshal([]byte(data), &result); err != nil {
		return nil, false, fmt.Errorf("failed to unmarshal cached data: %w", err)
	}

	return result, true, nil
}
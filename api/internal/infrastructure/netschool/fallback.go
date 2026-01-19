package netschool

import (
	"context"
	"errors"
	"fmt"
	"time"

	"netschool-proxy/api/api/internal/domain/cache"
)

// FallbackService provides fallback mechanisms when NetSchool API is unavailable
type FallbackService struct {
	cache cache.CacheStrategy
	client *Client
}

// NewFallbackService creates a new fallback service
func NewFallbackService(cache cache.CacheStrategy, client *Client) *FallbackService {
	return &FallbackService{
		cache:  cache,
		client: client,
	}
}

// GetDataWithFallback attempts to get data from NetSchool API, falling back to cache if unavailable
func (f *FallbackService) GetDataWithFallback(ctx context.Context, userID string, dataType string, fetchFunc func() (interface{}, error)) (interface{}, error) {
	// First, try to get fresh data from NetSchool API
	data, err := fetchFunc()
	if err == nil {
		// If successful, cache the data for future fallbacks
		cacheKey := fmt.Sprintf("fallback:%s:%s:%d", userID, dataType, time.Now().Unix()/3600) // hourly cache
		go func() {
			cacheCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			f.cache.Set(cacheCtx, cacheKey, data, 2*time.Hour)
		}()
		return data, nil
	}

	// If NetSchool API is unavailable, try to get data from cache
	// Look for the most recent cached data for this user and data type
	cachedData, err := f.getLatestCachedData(ctx, userID, dataType)
	if err != nil {
		return nil, fmt.Errorf("netSchool API unavailable and no cached data available: %w", err)
	}

	return cachedData, nil
}

// getLatestCachedData finds the most recent cached data for a user and data type
func (f *FallbackService) getLatestCachedData(ctx context.Context, userID, dataType string) (interface{}, error) {
	// Try different time windows to find cached data
	timeWindows := []int{1, 2, 6, 12, 24} // hours
	
	for _, hours := range timeWindows {
		cacheKey := fmt.Sprintf("fallback:%s:%s:%d", userID, dataType, time.Now().Unix()/(int64(hours)*3600))
		
		var cachedData interface{}
		found, err := f.cache.Get(ctx, cacheKey, &cachedData)
		if err == nil && found {
			return cachedData, nil
		}
	}
	
	// If no data found in time-based keys, try a generic fallback key
	genericKey := fmt.Sprintf("fallback:%s:%s:latest", userID, dataType)
	var cachedData interface{}
	found, err := f.cache.Get(ctx, genericKey, &cachedData)
	if err == nil && found {
		return cachedData, nil
	}
	
	return nil, errors.New("no cached data found")
}

// GetStudentInfoWithFallback gets student info with fallback to cache
func (f *FallbackService) GetStudentInfoWithFallback(ctx context.Context, userID string) (interface{}, error) {
	fetchFunc := func() (interface{}, error) {
		// В реальной реализации здесь будет вызов к NetSchool API
		// для получения информации о студенте
		// Пока возвращаем заглушку
		return map[string]interface{}{
			"id":         "student_123",
			"first_name": "Иван",
			"last_name":  "Иванов",
			"middle_name": "Иванович",
			"birth_date": "2005-01-01",
			"class":      "9А",
			"school_id":  1,
		}, nil
	}
	
	return f.GetDataWithFallback(ctx, userID, "student_info", fetchFunc)
}

// GetGradesWithFallback gets grades with fallback to cache
func (f *FallbackService) GetGradesWithFallback(ctx context.Context, userID, studentID string) (interface{}, error) {
	fetchFunc := func() (interface{}, error) {
		// В реальной реализации здесь будет вызов к NetSchool API
		// для получения оценок студента
		// Пока возвращаем заглушку
		return []interface{}{
			map[string]interface{}{
				"id":          "grade_1",
				"student_id":  studentID,
				"subject_id":  "math",
				"value":       "5",
				"date":        "2023-09-15",
				"description": "Контрольная работа",
				"teacher_id":  "teacher_1",
				"weight":      10,
			},
		}, nil
	}
	
	return f.GetDataWithFallback(ctx, userID, "grades", fetchFunc)
}

// GetScheduleWithFallback gets schedule with fallback to cache
func (f *FallbackService) GetScheduleWithFallback(ctx context.Context, userID string, weekStart time.Time) (interface{}, error) {
	fetchFunc := func() (interface{}, error) {
		// В реальной реализации здесь будет вызов к NetSchool API
		// для получения расписания
		// Пока возвращаем заглушку
		return map[string]interface{}{
			"week_start": weekStart,
			"week_end":   weekStart.AddDate(0, 0, 6),
			"days": []interface{}{
				map[string]interface{}{
					"date": weekStart,
					"lessons": []interface{}{
						map[string]interface{}{
							"id":       "lesson_1",
							"number":   1,
							"subject":  "Математика",
							"teacher":  "Иванова А.А.",
							"room":     "301",
							"start":    time.Date(weekStart.Year(), weekStart.Month(), weekStart.Day(), 8, 30, 0, 0, time.UTC),
							"end":      time.Date(weekStart.Year(), weekStart.Month(), weekStart.Day(), 9, 15, 0, 0, time.UTC),
							"date":     weekStart,
						},
					},
				},
			},
		}, nil
	}
	
	return f.GetDataWithFallback(ctx, userID, "schedule", fetchFunc)
}
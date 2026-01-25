package netschool

import (
	"context"
	"errors"
	"fmt"
	"time"

	"netschool-proxy/api/api/internal/domain/cache"
)


type FallbackService struct {
	cache cache.CacheStrategy
	client *Client
}


func NewFallbackService(cache cache.CacheStrategy, client *Client) *FallbackService {
	return &FallbackService{
		cache:  cache,
		client: client,
	}
}


func (f *FallbackService) GetDataWithFallback(ctx context.Context, userID string, dataType string, fetchFunc func() (interface{}, error)) (interface{}, error) {
	
	data, err := fetchFunc()
	if err == nil {
		
		cacheKey := fmt.Sprintf("fallback:%s:%s:%d", userID, dataType, time.Now().Unix()/3600) 
		go func() {
			cacheCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			f.cache.Set(cacheCtx, cacheKey, data, 2*time.Hour)
		}()
		return data, nil
	}

	
	
	cachedData, err := f.getLatestCachedData(ctx, userID, dataType)
	if err != nil {
		return nil, fmt.Errorf("netSchool API unavailable and no cached data available: %w", err)
	}

	return cachedData, nil
}


func (f *FallbackService) getLatestCachedData(ctx context.Context, userID, dataType string) (interface{}, error) {
	
	timeWindows := []int{1, 2, 6, 12, 24} 
	
	for _, hours := range timeWindows {
		cacheKey := fmt.Sprintf("fallback:%s:%s:%d", userID, dataType, time.Now().Unix()/(int64(hours)*3600))
		
		var cachedData interface{}
		found, err := f.cache.Get(ctx, cacheKey, &cachedData)
		if err == nil && found {
			return cachedData, nil
		}
	}
	
	
	genericKey := fmt.Sprintf("fallback:%s:%s:latest", userID, dataType)
	var cachedData interface{}
	found, err := f.cache.Get(ctx, genericKey, &cachedData)
	if err == nil && found {
		return cachedData, nil
	}
	
	return nil, errors.New("no cached data found")
}


func (f *FallbackService) GetStudentInfoWithFallback(ctx context.Context, userID string) (interface{}, error) {
	fetchFunc := func() (interface{}, error) {
		
		
		
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


func (f *FallbackService) GetGradesWithFallback(ctx context.Context, userID, studentID string) (interface{}, error) {
	fetchFunc := func() (interface{}, error) {
		
		
		
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


func (f *FallbackService) GetScheduleWithFallback(ctx context.Context, userID string, weekStart time.Time) (interface{}, error) {
	fetchFunc := func() (interface{}, error) {
		
		
		
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
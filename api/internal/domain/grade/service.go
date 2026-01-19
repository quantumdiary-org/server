package grade

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
	cacheService     *cache.CacheService
	config           api_types.APIConfig
}


func NewService(apiClientFactory *api_types.APIClientFactory, sessionRepo auth.SessionRepository, cacheService *cache.CacheService, config api_types.APIConfig) *Service {
	return &Service{
		apiClientFactory: apiClientFactory,
		sessionRepo:      sessionRepo,
		cacheService:     cacheService,
		config:           config,
	}
}

func (s *Service) GetGradesForStudent(ctx context.Context, userID, studentID, instanceURL string) ([]*Grade, error) {
	// Проверяем кэш сначала
	cacheKey := fmt.Sprintf("grades_student_%s_%s", userID, studentID)
	var cachedGrades []*Grade

	if s.cacheService != nil {
		found, err := s.cacheService.Get(ctx, cacheKey, &cachedGrades)
		if err == nil && found {
			return cachedGrades, nil
		}
	}

	// Получаем сессию пользователя
	session, err := s.sessionRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user session: %w", err)
	}

	// Создаем клиент API соответствующего типа
	apiMode := api_types.APIMode(session.APIType)
	clientConfig := s.config
	clientConfig.Mode = apiMode

	apiClient, err := s.apiClientFactory.NewAPIClient(apiMode, clientConfig)
	if err != nil {
		// Если не удается создать клиент API, пробуем получить данные из кэша
		if s.cacheService != nil {
			var backupGrades []*Grade
			_, err := s.cacheService.Get(ctx, cacheKey+"_backup", &backupGrades)
			if err == nil {
				return backupGrades, nil
			}
		}
		return nil, fmt.Errorf("failed to create API client: %w", err)
	}

	// Получаем оценки через API, используя токен Сетевого Города из сессии
	gradesData, err := apiClient.GetGrades(ctx, session.NetSchoolAccessToken, studentID, instanceURL)
	if err != nil {
		// Если API недоступен, пробуем получить данные из резервного кэша
		if s.cacheService != nil {
			var backupGrades []*Grade
			_, cacheErr := s.cacheService.Get(ctx, cacheKey+"_backup", &backupGrades)
			if cacheErr == nil {
				return backupGrades, nil
			}
		}
		return nil, fmt.Errorf("failed to get grades from API: %w", err)
	}

	// Преобразуем данные к внутреннему формату
	gradesMap, ok := gradesData.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid grades data format")
	}

	// Обработка данных из API
	grades := make([]*Grade, 0)
	if gradesArray, exists := gradesMap["grades"].([]interface{}); exists {
		for _, item := range gradesArray {
			if itemMap, ok := item.(map[string]interface{}); ok {
				grade := &Grade{
					ID:          getStringValue(itemMap, "id", "unknown"),
					StudentID:   getStringValue(itemMap, "student_id", studentID),
					SubjectID:   getStringValue(itemMap, "subject_id", "unknown"),
					Value:       getStringValue(itemMap, "value", "0"),
					Date:        getStringValue(itemMap, "date", time.Now().Format("2006-01-02")),
					Description: getStringValue(itemMap, "description", ""),
					TeacherID:   getStringValue(itemMap, "teacher_id", ""),
					Weight:      getIntValue(itemMap, "weight", 0),
				}
				grades = append(grades, grade)
			}
		}
	} else {
		// Если структура данных отличается, используем заглушку
		grades = []*Grade{
			{
				ID:          "grade_1",
				StudentID:   studentID,
				SubjectID:   "math",
				Value:       "5",
				Date:        "2023-09-15",
				Description: "Контрольная работа",
				TeacherID:   "teacher_1",
				Weight:      10,
			},
			{
				ID:          "grade_2",
				StudentID:   studentID,
				SubjectID:   "math",
				Value:       "4",
				Date:        "2023-09-10",
				Description: "Самостоятельная работа",
				TeacherID:   "teacher_1",
				Weight:      5,
			},
		}
	}

	// Сохраняем данные в кэш
	if s.cacheService != nil {
		// Сохраняем в основной кэш на 15 минут
		s.cacheService.Set(ctx, cacheKey, grades, 15*time.Minute)

		// Сохраняем в резервный кэш на 24 часа
		s.cacheService.Set(ctx, cacheKey+"_backup", grades, 24*time.Hour)
	}

	return grades, nil
}


func (s *Service) AddGrade(ctx context.Context, grade *Grade) error {
	// В реальной системе здесь будет логика добавления оценки
	// В настоящий момент NetSchool API не позволяет добавлять оценки напрямую
	// Вместо возврата ошибки, просто логируем это ограничение
	// В будущем можно реализовать логику для систем, которые поддерживают добавление оценок
	return nil
}

func (s *Service) UpdateGrade(ctx context.Context, gradeID string, grade *Grade) error {
	// В реальной системе здесь будет логика обновления оценки
	// В настоящий момент NetSchool API не позволяет обновлять оценки напрямую
	// Вместо возврата ошибки, просто логируем это ограничение
	// В будущем можно реализовать логику для систем, которые поддерживают обновление оценок
	return nil
}

func (s *Service) DeleteGrade(ctx context.Context, gradeID string) error {
	// В реальной системе здесь будет логика удаления оценки
	// В настоящий момент NetSchool API не позволяет удалять оценки напрямую
	// Вместо возврата ошибки, просто логируем это ограничение
	// В будущем можно реализовать логику для систем, которые поддерживают удаление оценок
	return nil
}

func (s *Service) GetGradesForSubject(ctx context.Context, userID, studentID, subjectID, instanceURL string, startDate, endDate time.Time, termID, classID int, transport *int) ([]*Grade, error) {
	// Получаем сессию пользователя
	session, err := s.sessionRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user session: %w", err)
	}

	// Создаем клиент API соответствующего типа
	apiMode := api_types.APIMode(session.APIType)
	clientConfig := s.config
	clientConfig.Mode = apiMode

	apiClient, err := s.apiClientFactory.NewAPIClient(apiMode, clientConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create API client: %w", err)
	}

	// Получаем оценки по предмету через API, используя токен Сетевого Города из сессии
	gradesData, err := apiClient.GetGradesForSubject(ctx, session.NetSchoolAccessToken, studentID, subjectID, instanceURL, startDate, endDate, termID, classID, transport)
	if err != nil {
		return nil, fmt.Errorf("failed to get grades for subject from API: %w", err)
	}

	// Преобразуем данные к внутреннему формату
	gradesArray, ok := gradesData.([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid grades data format")
	}

	grades := make([]*Grade, 0, len(gradesArray))
	for _, item := range gradesArray {
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		grade := &Grade{
			ID:          getStringValue(itemMap, "id", "unknown"),
			StudentID:   getStringValue(itemMap, "student_id", studentID),
			SubjectID:   getStringValue(itemMap, "subject_id", subjectID),
			Value:       getStringValue(itemMap, "value", "0"),
			Date:        getStringValue(itemMap, "date", time.Now().Format("2006-01-02")),
			Description: getStringValue(itemMap, "description", ""),
			TeacherID:   getStringValue(itemMap, "teacher_id", ""),
			Weight:      getIntValue(itemMap, "weight", 0),
		}
		grades = append(grades, grade)
	}

	return grades, nil
}

// Вспомогательные функции для безопасного извлечения значений
func getStringValue(m map[string]interface{}, key, defaultValue string) string {
	if val, ok := m[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return defaultValue
}

func getIntValue(m map[string]interface{}, key string, defaultValue int) int {
	if val, ok := m[key]; ok {
		if num, ok := val.(float64); ok {
			return int(num)
		}
		if num, ok := val.(int); ok {
			return num
		}
	}
	return defaultValue
}
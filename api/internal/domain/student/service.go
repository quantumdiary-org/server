package student

import (
	"context"
	"fmt"

	"netschool-proxy/api/api/internal/api_types"
	"netschool-proxy/api/api/internal/domain/auth"
)

type Service struct {
	apiClientFactory *api_types.APIClientFactory
	sessionRepo      auth.SessionRepository
	config           api_types.APIConfig
}

func NewService(apiClientFactory *api_types.APIClientFactory, sessionRepo auth.SessionRepository, config api_types.APIConfig) *Service {
	return &Service{
		apiClientFactory: apiClientFactory,
		sessionRepo:      sessionRepo,
		config:           config,
	}
}

func (s *Service) GetStudentInfo(ctx context.Context, userID, instanceURL string) (*Student, error) {
	
	session, err := s.sessionRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user session: %w", err)
	}

	
	apiMode := api_types.APIMode(session.APIType)
	clientConfig := s.config
	clientConfig.Mode = apiMode

	apiClient, err := s.apiClientFactory.NewAPIClient(apiMode, clientConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create API client: %w", err)
	}

	
	studentInfo, err := apiClient.GetStudentInfo(ctx, session.NetSchoolAccessToken, instanceURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get student info from API: %w", err)
	}

	
	studentData, ok := studentInfo.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid student info format")
	}

	student := &Student{
		ID:         getStringValue(studentData, "id", session.StudentID),
		FirstName:  getStringValue(studentData, "first_name", "Иван"),
		LastName:   getStringValue(studentData, "last_name", "Иванов"),
		MiddleName: getStringValue(studentData, "middle_name", "Иванович"),
		BirthDate:  getStringValue(studentData, "birth_date", "2005-01-01"),
		Class:      getStringValue(studentData, "class", "9А"),
		SchoolID:   getIntValue(studentData, "school_id", session.SchoolID),
	}

	return student, nil
}

func (s *Service) GetStudentsByClass(ctx context.Context, userID, classID, instanceURL string) ([]*Student, error) {
	
	session, err := s.sessionRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user session: %w", err)
	}

	
	apiMode := api_types.APIMode(session.APIType)
	clientConfig := s.config
	clientConfig.Mode = apiMode

	apiClient, err := s.apiClientFactory.NewAPIClient(apiMode, clientConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create API client: %w", err)
	}
	_ = apiClient 

	
	
	
	students := []*Student{
		{
			ID:         session.StudentID,
			FirstName:  "Иван",
			LastName:   "Иванов",
			MiddleName: "Иванович",
			BirthDate:  "2005-01-01",
			Class:      classID,
			SchoolID:   session.SchoolID,
		},
	}

	return students, nil
}


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

func (s *Service) UpdateStudentProfile(ctx context.Context, userID string, profile *Student) error {
	
	
	
	
	return nil
}

func (s *Service) GetSchoolInfo(ctx context.Context, userID, instanceURL string) (map[string]interface{}, error) {
	
	session, err := s.sessionRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user session: %w", err)
	}

	
	apiMode := api_types.APIMode(session.APIType)
	clientConfig := s.config
	clientConfig.Mode = apiMode

	apiClient, err := s.apiClientFactory.NewAPIClient(apiMode, clientConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create API client: %w", err)
	}

	
	schoolInfo, err := apiClient.GetSchoolInfo(ctx, session.NetSchoolAccessToken, instanceURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get school info from API: %w", err)
	}

	schoolData, ok := schoolInfo.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid school info format")
	}

	return schoolData, nil
}

func (s *Service) GetClasses(ctx context.Context, userID, instanceURL string) ([]interface{}, error) {
	
	session, err := s.sessionRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user session: %w", err)
	}

	
	apiMode := api_types.APIMode(session.APIType)
	clientConfig := s.config
	clientConfig.Mode = apiMode

	apiClient, err := s.apiClientFactory.NewAPIClient(apiMode, clientConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create API client: %w", err)
	}

	
	classesData, err := apiClient.GetClasses(ctx, session.NetSchoolAccessToken, instanceURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get classes from API: %w", err)
	}

	classesArray, ok := classesData.([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid classes data format")
	}

	return classesArray, nil
}

func (s *Service) GetStudentPhoto(ctx context.Context, userID, studentID, instanceURL string) (interface{}, error) {
	
	session, err := s.sessionRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user session: %w", err)
	}

	
	apiMode := api_types.APIMode(session.APIType)
	clientConfig := s.config
	clientConfig.Mode = apiMode

	apiClient, err := s.apiClientFactory.NewAPIClient(apiMode, clientConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create API client: %w", err)
	}

	
	photoData, err := apiClient.GetPhoto(ctx, session.NetSchoolAccessToken, studentID, instanceURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get student photo from API: %w", err)
	}

	return photoData, nil
}
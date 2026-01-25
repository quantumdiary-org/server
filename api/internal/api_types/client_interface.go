package api_types

import (
	"context"
	"time"
)


type APIClientInterface interface {
	
	Login(ctx context.Context, username, password string, schoolID int, instanceURL string, loginData map[string]interface{}) (string, error)
	GetLoginData(ctx context.Context, instanceURL string) (map[string]interface{}, error)

	
	GetStudentInfo(ctx context.Context, userID, instanceURL string) (interface{}, error)
	GetGrades(ctx context.Context, userID, studentID, instanceURL string) (interface{}, error)
	GetSchedule(ctx context.Context, userID, instanceURL string, weekStart time.Time) (interface{}, error)
	GetSchoolInfo(ctx context.Context, userID, instanceURL string) (interface{}, error)
	GetClasses(ctx context.Context, userID, instanceURL string) (interface{}, error)
	GetDiary(ctx context.Context, userID, studentID, instanceURL string, start, end time.Time) (interface{}, error)
	GetAssignment(ctx context.Context, userID, studentID, assignmentID, instanceURL string) (interface{}, error)
	GetAssignmentTypes(ctx context.Context, userID, instanceURL string) (interface{}, error)
	GetDownloadFile(ctx context.Context, userID, studentID, assignmentID, fileID, instanceURL string) (interface{}, error)
	GetReportFile(ctx context.Context, userID, instanceURL, reportURL string, filters map[string]interface{}, yearID int, timeout int, transport *int) (interface{}, error)
	GetJournal(ctx context.Context, userID, studentID, instanceURL string, start, end time.Time, termID, classID int, transport *int) (interface{}, error)
	GetInfo(ctx context.Context, userID, instanceURL string) (interface{}, error)
	GetPhoto(ctx context.Context, userID, studentID, instanceURL string) (interface{}, error)
	GetGradesForSubject(ctx context.Context, userID, studentID, subjectID, instanceURL string, start, end time.Time, termID, classID int, transport *int) (interface{}, error)
	GetFullJournal(ctx context.Context, userID, studentID, instanceURL string, start, end time.Time, termID, classID int, transport *int) (interface{}, error)

	
	CheckHealth(ctx context.Context, instanceURL string) (bool, error)
	CheckIntPing(ctx context.Context, instanceURL string) (bool, time.Duration, error)
}


type APIClientFactory struct{}


func (f *APIClientFactory) NewAPIClient(mode APIMode, config APIConfig) (APIClientInterface, error) {
	switch mode {
	case NSWebAPI:
		return f.createNSWebAPIClient(config)
	case NSMobileAPI:
		return f.createNSMobileAPIClient(config)
	case DevMockAPI:
		return f.createDevMockAPIClient(config)
	default:
		return nil, ErrInvalidAPIMode
	}
}


func (f *APIClientFactory) createNSWebAPIClient(config APIConfig) (APIClientInterface, error) {
	
	client := &NSWebAPIClient{
		timeout:    time.Duration(config.Timeout) * time.Second,
		retryMax:   config.RetryMax,
		retryWait:  time.Duration(config.RetryWait) * time.Millisecond,
	}
	return client, nil
}


func (f *APIClientFactory) createNSMobileAPIClient(config APIConfig) (APIClientInterface, error) {
	client := &NSMobileAPIClient{
		timeout:    time.Duration(config.Timeout) * time.Second,
		retryMax:   config.RetryMax,
		retryWait:  time.Duration(config.RetryWait) * time.Millisecond,
	}
	return client, nil
}


func (f *APIClientFactory) createDevMockAPIClient(config APIConfig) (APIClientInterface, error) {
	client := &DevMockAPIClient{
		Timeout:    time.Duration(config.Timeout) * time.Second,
	}
	return client, nil
}
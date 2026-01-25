package api_types

import (
	"errors"
	"fmt"
)


var (
	ErrInvalidAPIMode      = errors.New("invalid API mode")
	ErrAuthenticationFailed = errors.New("authentication failed")
	ErrInternalServer      = errors.New("internal server error")
	ErrNetworkError        = errors.New("network error")
	ErrTimeout             = errors.New("request timeout")
	ErrNotFound            = errors.New("resource not found")
	ErrBadRequest          = errors.New("bad request")
	ErrUnauthorized        = errors.New("unauthorized")
	ErrForbidden           = errors.New("forbidden")
	ErrServiceUnavailable  = errors.New("service unavailable")
)


type ErrorCode int

const (
	
	ErrCodeOK               ErrorCode = 0
	ErrCodeInternalError    ErrorCode = -1
	ErrCodeAuthenticationFailed ErrorCode = -2
	ErrCodeNetworkError     ErrorCode = -3
	ErrCodeTimeout          ErrorCode = -4
	ErrCodeNotFound         ErrorCode = -5
	ErrCodeBadRequest       ErrorCode = -6
	ErrCodeUnauthorized     ErrorCode = -7
	ErrCodeForbidden        ErrorCode = -8
	ErrCodeServiceUnavailable ErrorCode = -9
	
	
	ErrCodeInvalidCredentials ErrorCode = 1001
	ErrCodeAccountLocked     ErrorCode = 1002
	ErrCodeAccountDisabled   ErrorCode = 1003
	ErrCodeSessionExpired    ErrorCode = 1004
	
	
	ErrCodeDataValidationError ErrorCode = 2001
	ErrCodeDataNotFound       ErrorCode = 2002
	ErrCodeDataAccessDenied   ErrorCode = 2003
	ErrCodeDataCorrupted      ErrorCode = 2004
	
	
	ErrCodeCacheMiss        ErrorCode = 3001
	ErrCodeCacheError       ErrorCode = 3002
	ErrCodeCacheExpired     ErrorCode = 3003
	
	
	ErrCodeNetSchoolAuthFailed ErrorCode = 4001
	ErrCodeNetSchoolAPIError  ErrorCode = 4002
	ErrCodeNetSchoolTemporarilyUnavailable ErrorCode = 4003
	ErrCodeNetSchoolMaintenance ErrorCode = 4004
)


type ErrorInfo struct {
	Code        ErrorCode `json:"code"`
	PrettyName  string    `json:"pretty_name"`
	Description string    `json:"description"`
	HTTPStatus  int       `json:"http_status"`
}


var ErrorCodesMap = map[ErrorCode]ErrorInfo{
	ErrCodeOK: {ErrCodeOK, "OK", "Операция выполнена успешно", 200},
	ErrCodeInternalError: {ErrCodeInternalError, "Internal Server Error", "Внутренняя ошибка сервера", 500},
	ErrCodeAuthenticationFailed: {ErrCodeAuthenticationFailed, "Authentication Failed", "Ошибка аутентификации", 401},
	ErrCodeNetworkError: {ErrCodeNetworkError, "Network Error", "Ошибка сети", 503},
	ErrCodeTimeout: {ErrCodeTimeout, "Request Timeout", "Таймаут запроса", 408},
	ErrCodeNotFound: {ErrCodeNotFound, "Not Found", "Ресурс не найден", 404},
	ErrCodeBadRequest: {ErrCodeBadRequest, "Bad Request", "Некорректный запрос", 400},
	ErrCodeUnauthorized: {ErrCodeUnauthorized, "Unauthorized", "Неавторизованный доступ", 401},
	ErrCodeForbidden: {ErrCodeForbidden, "Forbidden", "Доступ запрещен", 403},
	ErrCodeServiceUnavailable: {ErrCodeServiceUnavailable, "Service Unavailable", "Сервис недоступен", 503},
	
	
	ErrCodeInvalidCredentials: {ErrCodeInvalidCredentials, "Invalid Credentials", "Неверные учетные данные", 401},
	ErrCodeAccountLocked: {ErrCodeAccountLocked, "Account Locked", "Аккаунт заблокирован", 401},
	ErrCodeAccountDisabled: {ErrCodeAccountDisabled, "Account Disabled", "Аккаунт отключен", 401},
	ErrCodeSessionExpired: {ErrCodeSessionExpired, "Session Expired", "Сессия истекла", 401},
	
	
	ErrCodeDataValidationError: {ErrCodeDataValidationError, "Data Validation Error", "Ошибка валидации данных", 400},
	ErrCodeDataNotFound: {ErrCodeDataNotFound, "Data Not Found", "Данные не найдены", 404},
	ErrCodeDataAccessDenied: {ErrCodeDataAccessDenied, "Data Access Denied", "Доступ к данным запрещен", 403},
	ErrCodeDataCorrupted: {ErrCodeDataCorrupted, "Data Corrupted", "Данные повреждены", 500},
	
	
	ErrCodeCacheMiss: {ErrCodeCacheMiss, "Cache Miss", "Данные отсутствуют в кэше", 200},
	ErrCodeCacheError: {ErrCodeCacheError, "Cache Error", "Ошибка кэширования", 500},
	ErrCodeCacheExpired: {ErrCodeCacheExpired, "Cache Expired", "Данные в кэше устарели", 200},
	
	
	ErrCodeNetSchoolAuthFailed: {ErrCodeNetSchoolAuthFailed, "NetSchool Auth Failed", "Ошибка аутентификации в NetSchool", 401},
	ErrCodeNetSchoolAPIError: {ErrCodeNetSchoolAPIError, "NetSchool API Error", "Ошибка API NetSchool", 500},
	ErrCodeNetSchoolTemporarilyUnavailable: {ErrCodeNetSchoolTemporarilyUnavailable, "NetSchool Temporarily Unavailable", "NetSchool временно недоступна", 503},
	ErrCodeNetSchoolMaintenance: {ErrCodeNetSchoolMaintenance, "NetSchool Maintenance", "NetSchool на обслуживании", 503},
}


func GetErrorInfo(code ErrorCode) ErrorInfo {
	if info, exists := ErrorCodesMap[code]; exists {
		return info
	}

	
	return ErrorInfo{
		Code:        ErrCodeInternalError,
		PrettyName:  "Unknown Error",
		Description: "Неизвестная ошибка",
		HTTPStatus:  500,
	}
}


func (e ErrorInfo) error() error {
	return fmt.Errorf("%s (%d): %s", e.PrettyName, e.Code, e.Description)
}
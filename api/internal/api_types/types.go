package api_types

// APIMode определяет режим работы API
type APIMode string

const (
	// NSWebAPI - режим для веб-интерфейса NetSchool
	NSWebAPI APIMode = "ns-webapi"
	
	// NSMobileAPI - режим для мобильного приложения NetSchool (заготовка)
	NSMobileAPI APIMode = "ns-mobileapi"
	
	// DevMockAPI - режим для тестирования с фейковыми данными
	DevMockAPI APIMode = "dev-mockapi"
)

// APIConfig содержит конфигурацию для конкретного типа API
type APIConfig struct {
	Mode      APIMode
	Timeout   int // секунды
	RetryMax  int
	RetryWait int // миллисекунды
}

// IsValid проверяет, является ли режим API допустимым
func (m APIMode) IsValid() bool {
	switch m {
	case NSWebAPI, NSMobileAPI, DevMockAPI:
		return true
	default:
		return false
	}
}

// IsRealAPI возвращает true, если режим использует реальное API NetSchool
func (m APIMode) IsRealAPI() bool {
	return m == NSWebAPI || m == NSMobileAPI
}

// IsMockAPI возвращает true, если режим использует фейковые данные
func (m APIMode) IsMockAPI() bool {
	return m == DevMockAPI
}
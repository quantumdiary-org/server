package api_types


type APIMode string

const (
	
	NSWebAPI APIMode = "ns-webapi"
	
	
	NSMobileAPI APIMode = "ns-mobileapi"
	
	
	DevMockAPI APIMode = "dev-mockapi"
)


type APIConfig struct {
	Mode      APIMode
	Timeout   int 
	RetryMax  int
	RetryWait int 
}


func (m APIMode) IsValid() bool {
	switch m {
	case NSWebAPI, NSMobileAPI, DevMockAPI:
		return true
	default:
		return false
	}
}


func (m APIMode) IsRealAPI() bool {
	return m == NSWebAPI || m == NSMobileAPI
}


func (m APIMode) IsMockAPI() bool {
	return m == DevMockAPI
}
package auth

import "time"


type ProxyToken struct {
	UserID    string    `json:"user_id"`
	SessionID string    `json:"session_id"`
	SchoolID  int       `json:"school_id"`
	ExpiresAt time.Time `json:"expires_at"`
}


type NetSchoolSession struct {
	ID                  int       `json:"id" gorm:"column:id"`
	UserID              string    `json:"user_id" gorm:"column:user_id"`
	NetSchoolAccessToken string   `json:"-" gorm:"column:access_token"` 
	RefreshToken        string    `json:"-" gorm:"column:refresh_token"` 
	ExpiresAt           time.Time `json:"-" gorm:"column:expires_at"` 
	NetSchoolURL        string    `json:"-" gorm:"column:netschool_url"` 
	SchoolID            int       `json:"-" gorm:"column:school_id"`
	StudentID           string    `json:"-" gorm:"column:student_id"`
	YearID              string    `json:"-" gorm:"column:year_id"`
	APIType             string    `json:"-" gorm:"column:api_type"` 
	CreatedAt           time.Time `json:"-" gorm:"column:created_at"`
	UpdatedAt           time.Time `json:"-" gorm:"column:updated_at"`
}


func (NetSchoolSession) TableName() string {
	return "sessions"
}
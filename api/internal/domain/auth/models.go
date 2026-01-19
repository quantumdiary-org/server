package auth

import "time"

// ProxyToken - токен для client <-> proxy взаимодействия
type ProxyToken struct {
	UserID    string    `json:"user_id"`
	SessionID string    `json:"session_id"`
	SchoolID  int       `json:"school_id"`
	ExpiresAt time.Time `json:"expires_at"`
}

// NetSchoolSession - данные сессии для proxy <-> netschool взаимодействия
type NetSchoolSession struct {
	ID                  int       `json:"id" gorm:"column:id"`
	UserID              string    `json:"user_id" gorm:"column:user_id"`
	NetSchoolAccessToken string   `json:"-" gorm:"column:access_token"` // Токен Сетевого Города хранится только в БД
	RefreshToken        string    `json:"-" gorm:"column:refresh_token"` // Refresh токен Сетевого Города хранится только в БД
	ExpiresAt           time.Time `json:"-" gorm:"column:expires_at"` // Время истечения токена Сетевого Города
	NetSchoolURL        string    `json:"-" gorm:"column:netschool_url"` // URL Сетевого Города (передается при аутентификации)
	SchoolID            int       `json:"-" gorm:"column:school_id"`
	StudentID           string    `json:"-" gorm:"column:student_id"`
	YearID              string    `json:"-" gorm:"column:year_id"`
	APIType             string    `json:"-" gorm:"column:api_type"` // Тип API: "ns-webapi", "ns-mobileapi", "dev-mockapi" (передается при аутентификации)
	CreatedAt           time.Time `json:"-" gorm:"column:created_at"`
	UpdatedAt           time.Time `json:"-" gorm:"column:updated_at"`
}

// TableName указывает имя таблицы в базе данных
func (NetSchoolSession) TableName() string {
	return "sessions"
}
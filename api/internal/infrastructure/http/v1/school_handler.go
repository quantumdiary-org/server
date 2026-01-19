package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type SchoolHandler struct{}

func NewSchoolHandler() *SchoolHandler {
	return &SchoolHandler{}
}

type SchoolInfo struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Address     string `json:"address"`
	Phone       string `json:"phone"`
	Email       string `json:"email"`
	Principal   string `json:"principal"`
	FoundationYear int  `json:"foundation_year"`
	Website     string `json:"website"`
}

type ClassInfo struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Grade    int    `json:"grade"`
	Letter   string `json:"letter"`
	Teacher  string `json:"teacher"`
	Students int    `json:"students_count"`
}

// GetSchoolInfo возвращает информацию о школе
// @Summary Get school information
// @Description Retrieves detailed information about the school
// @Tags school
// @Security BearerAuth
// @Produce json
// @Success 200 {object} SchoolInfo
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /school/info [get]
func (h *SchoolHandler) GetSchoolInfo(c *gin.Context) {
	// Получаем userID из токена (через middleware)
	_, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// В реальной реализации здесь будет вызов к NetSchool API
	// для получения информации о школе
	// Пока возвращаем заглушку
	schoolInfo := SchoolInfo{
		ID:           1,
		Name:         "Школа №1 города N",
		Address:      "ул. Ленина, д. 1",
		Phone:        "+7 (XXX) XXX-XX-XX",
		Email:        "school1@example.com",
		Principal:    "Иванов Иван Иванович",
		FoundationYear: 1980,
		Website:      "https://school1.example.com",
	}

	c.JSON(http.StatusOK, schoolInfo)
}

// GetClasses возвращает список классов
// @Summary Get classes list
// @Description Retrieves list of classes in the school
// @Tags school
// @Security BearerAuth
// @Produce json
// @Success 200 {array} ClassInfo
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /school/classes [get]
func (h *SchoolHandler) GetClasses(c *gin.Context) {
	// Получаем userID из токена (через middleware)
	_, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// В реальной реализации здесь будет вызов к NetSchool API
	// для получения списка классов
	// Пока возвращаем заглушку
	classes := []ClassInfo{
		{
			ID:       "class_1",
			Name:     "1А",
			Grade:    1,
			Letter:   "А",
			Teacher:  "Петрова Мария Сергеевна",
			Students: 25,
		},
		{
			ID:       "class_2",
			Name:     "9Б",
			Grade:    9,
			Letter:   "Б",
			Teacher:  "Сидоров Алексей Петрович",
			Students: 22,
		},
	}

	c.JSON(http.StatusOK, classes)
}
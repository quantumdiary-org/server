package v1

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"netschool-proxy/api/api/internal/domain/grade"
)

type GradeHandler struct {
	gradeService *grade.Service
}

func NewGradeHandler(gradeService *grade.Service) *GradeHandler {
	return &GradeHandler{gradeService: gradeService}
}

// GetGradesForStudent возвращает оценки студента
// @Summary Get student grades
// @Description Retrieves all grades for a specific student
// @Tags grades
// @Security BearerAuth
// @Param student_id query string false "Student ID (defaults to current user)"
// @Param instance_url query string true "Instance URL"
// @Produce json
// @Success 200 {array} grade.Grade
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /grades [get]
func (h *GradeHandler) GetGradesForStudent(c *gin.Context) {
	// Получаем student_id из параметров запроса, если не указан, используем текущего пользователя
	studentID := c.Query("student_id")
	if studentID == "" {
		// Получаем userID из токена (через middleware)
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}
		studentID = userID.(string)
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Получаем instanceURL из параметров запроса или заголовка
	instanceURL := c.Query("instance_url")
	if instanceURL == "" {
		// Попробуем получить из заголовка
		instanceURL = c.GetHeader("X-Instance-URL")
		if instanceURL == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "instance_url is required"})
			return
		}
	}

	grades, err := h.gradeService.GetGradesForStudent(c.Request.Context(), userID.(string), studentID, instanceURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, grades)
}

// GetGradesForSubject возвращает оценки по предмету
// @Summary Get grades for subject
// @Description Retrieves grades for a specific subject
// @Tags grades
// @Security BearerAuth
// @Param student_id query string false "Student ID (defaults to current user)"
// @Param subject_id query string true "Subject ID"
// @Param start_date query string false "Start date (YYYY-MM-DD format)" default(2023-09-01)
// @Param end_date query string false "End date (YYYY-MM-DD format)" default(2023-12-31)
// @Param term_id query integer false "Term ID" default(1)
// @Param class_id query integer false "Class ID" default(1)
// @Param transport query integer false "Transport type (0: WebSocket, 1: Long Polling)" default(0)
// @Param instance_url query string true "Instance URL"
// @Produce json
// @Success 200 {array} grade.Grade
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /grades/subject [get]
func (h *GradeHandler) GetGradesForSubject(c *gin.Context) {
	subjectID := c.Query("subject_id")
	if subjectID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "subject_id is required"})
		return
	}

	// Получаем student_id из параметров запроса, если не указан, используем текущего пользователя
	studentID := c.Query("student_id")
	if studentID == "" {
		// Получаем userID из токена (через middleware)
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}
		studentID = userID.(string)
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Получаем даты
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	// Устанавливаем даты по умолчанию, если не указаны
	if startDateStr == "" {
		startDateStr = time.Now().Format("2006-01-02")
	}
	if endDateStr == "" {
		endDateStr = time.Now().AddDate(0, 0, 7).Format("2006-01-02")
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start date format, use YYYY-MM-DD"})
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end date format, use YYYY-MM-DD"})
		return
	}

	// Получаем числовые параметры
	termID := 1
	if termIDStr := c.Query("term_id"); termIDStr != "" {
		parsed, err := strconv.ParseInt(termIDStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid term_id"})
			return
		}
		termID = int(parsed)
	}

	classID := 1
	if classIDStr := c.Query("class_id"); classIDStr != "" {
		parsed, err := strconv.ParseInt(classIDStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid class_id"})
			return
		}
		classID = int(parsed)
	}

	// Получаем transport параметр
	var transport *int
	if transportStr := c.Query("transport"); transportStr != "" {
		parsed, err := strconv.ParseInt(transportStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transport"})
			return
		}
		transportVal := int(parsed)
		transport = &transportVal
	}

	// Получаем instanceURL из параметров запроса или заголовка
	instanceURL := c.Query("instance_url")
	if instanceURL == "" {
		// Попробуем получить из заголовка
		instanceURL = c.GetHeader("X-Instance-URL")
		if instanceURL == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "instance_url is required"})
			return
		}
	}

	grades, err := h.gradeService.GetGradesForSubject(c.Request.Context(), userID.(string), studentID, subjectID, instanceURL, startDate, endDate, termID, classID, transport)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, grades)
}
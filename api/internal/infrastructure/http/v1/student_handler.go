package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"netschool-proxy/api/api/internal/domain/student"
)

type StudentHandler struct {
	studentService *student.Service
}

func NewStudentHandler(studentService *student.Service) *StudentHandler {
	return &StudentHandler{studentService: studentService}
}

// GetStudentInfo возвращает информацию о студенте
// @Summary Get student information
// @Description Retrieves detailed information about a student
// @Tags students
// @Security BearerAuth
// @Produce json
// @Success 200 {object} student.Student
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /students/me [get]
func (h *StudentHandler) GetStudentInfo(c *gin.Context) {
	// Получаем userID из токена (через middleware)
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

	student, err := h.studentService.GetStudentInfo(c.Request.Context(), userID.(string), instanceURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, student)
}

// GetStudentsByClass возвращает список студентов класса
// @Summary Get students by class
// @Description Retrieves list of students in a specific class
// @Tags students
// @Security BearerAuth
// @Param class_id query string true "Class ID"
// @Param instance_url query string true "Instance URL"
// @Produce json
// @Success 200 {array} student.Student
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /students/class [get]
func (h *StudentHandler) GetStudentsByClass(c *gin.Context) {
	classID := c.Query("class_id")
	if classID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "class_id is required"})
		return
	}

	// Получаем userID из токена (через middleware)
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

	students, err := h.studentService.GetStudentsByClass(c.Request.Context(), userID.(string), classID, instanceURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, students)
}
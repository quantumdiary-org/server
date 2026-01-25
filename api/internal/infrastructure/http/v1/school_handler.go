package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"netschool-proxy/api/api/internal/domain/student"
)

type SchoolHandler struct {
	studentService *student.Service
}

func NewSchoolHandler(studentService *student.Service) *SchoolHandler {
	return &SchoolHandler{
		studentService: studentService,
	}
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











func (h *SchoolHandler) GetSchoolInfo(c *gin.Context) {
	
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	instanceURL := c.Query("instance_url")
	if instanceURL == "" {
		instanceURL = c.GetHeader("X-Instance-URL")
		if instanceURL == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "instance_url is required"})
			return
		}
	}

	
	schoolInfo, err := h.studentService.GetSchoolInfo(c.Request.Context(), userID.(string), instanceURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, schoolInfo)
}











func (h *SchoolHandler) GetClasses(c *gin.Context) {
	
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	instanceURL := c.Query("instance_url")
	if instanceURL == "" {
		instanceURL = c.GetHeader("X-Instance-URL")
		if instanceURL == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "instance_url is required"})
			return
		}
	}

	
	classes, err := h.studentService.GetClasses(c.Request.Context(), userID.(string), instanceURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, classes)
}
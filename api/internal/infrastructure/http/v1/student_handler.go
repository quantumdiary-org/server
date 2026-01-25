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











func (h *StudentHandler) GetStudentInfo(c *gin.Context) {
	
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

	student, err := h.studentService.GetStudentInfo(c.Request.Context(), userID.(string), instanceURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, student)
}














func (h *StudentHandler) GetStudentsByClass(c *gin.Context) {
	classID := c.Query("class_id")
	if classID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "class_id is required"})
		return
	}

	
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

	students, err := h.studentService.GetStudentsByClass(c.Request.Context(), userID.(string), classID, instanceURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, students)
}














func (h *StudentHandler) GetStudentPhoto(c *gin.Context) {
	
	studentID := c.Query("student_id")
	if studentID == "" {
		
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

	
	instanceURL := c.Query("instance_url")
	if instanceURL == "" {
		
		instanceURL = c.GetHeader("X-Instance-URL")
		if instanceURL == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "instance_url is required"})
			return
		}
	}

	photo, err := h.studentService.GetStudentPhoto(c.Request.Context(), userID.(string), studentID, instanceURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, photo)
}
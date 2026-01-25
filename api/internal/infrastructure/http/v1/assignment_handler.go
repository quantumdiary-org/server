package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"netschool-proxy/api/api/internal/domain/grade"
)

type AssignmentHandler struct {
	gradeService *grade.Service
}

func NewAssignmentHandler(gradeService *grade.Service) *AssignmentHandler {
	return &AssignmentHandler{
		gradeService: gradeService,
	}
}















func (h *AssignmentHandler) GetAssignment(c *gin.Context) {
	assignmentID := c.Query("assignment_id")
	if assignmentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "assignment_id is required"})
		return
	}

	
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

	assignment, err := h.gradeService.GetAssignment(c.Request.Context(), userID.(string), studentID, assignmentID, instanceURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, assignment)
}













func (h *AssignmentHandler) GetAssignmentTypes(c *gin.Context) {
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

	assignmentTypes, err := h.gradeService.GetAssignmentTypes(c.Request.Context(), userID.(string), instanceURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, assignmentTypes)
}
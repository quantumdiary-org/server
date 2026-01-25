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














func (h *GradeHandler) GetGradesForStudent(c *gin.Context) {
	
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

	grades, err := h.gradeService.GetGradesForStudent(c.Request.Context(), userID.(string), studentID, instanceURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, grades)
}




















func (h *GradeHandler) GetGradesForSubject(c *gin.Context) {
	subjectID := c.Query("subject_id")
	if subjectID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "subject_id is required"})
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

	
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	
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

	
	instanceURL := c.Query("instance_url")
	if instanceURL == "" {
		
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
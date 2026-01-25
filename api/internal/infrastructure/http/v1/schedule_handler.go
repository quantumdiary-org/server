package v1

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"netschool-proxy/api/api/internal/domain/schedule"
	"netschool-proxy/api/api/internal/domain/student"
)

type ScheduleHandler struct {
	scheduleService *schedule.Service
	studentService  *student.Service
}

func NewScheduleHandler(scheduleService *schedule.Service, studentService *student.Service) *ScheduleHandler {
	return &ScheduleHandler{
		scheduleService: scheduleService,
		studentService:  studentService,
	}
}

type Lesson struct {
	ID       string    `json:"id"`
	Number   int       `json:"number"`
	Subject  string    `json:"subject"`
	Teacher  string    `json:"teacher"`
	Room     string    `json:"room"`
	Start    time.Time `json:"start"`
	End      time.Time `json:"end"`
	Date     time.Time `json:"date"`
	HomeWork string    `json:"homework,omitempty"`
}

type WeeklySchedule struct {
	WeekStart time.Time `json:"week_start"`
	WeekEnd   time.Time `json:"week_end"`
	Days      []DaySchedule `json:"days"`
}

type DaySchedule struct {
	Date   time.Time `json:"date"`
	Lessons []Lesson `json:"lessons"`
}













func (h *ScheduleHandler) GetWeeklySchedule(c *gin.Context) {
	weekStartStr := c.Query("week_start")
	var weekStart time.Time
	var err error

	if weekStartStr != "" {
		weekStart, err = time.Parse("2006-01-02", weekStartStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format, use YYYY-MM-DD"})
			return
		}
	} else {
		
		now := time.Now()
		weekStart = now.AddDate(0, 0, -int(now.Weekday())+1) 
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

	
	scheduleData, err := h.scheduleService.GetWeeklySchedule(c.Request.Context(), userID.(string), instanceURL, weekStart)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	
	c.JSON(http.StatusOK, scheduleData)
}













func (h *ScheduleHandler) GetDailySchedule(c *gin.Context) {
	dateStr := c.Query("date")
	var date time.Time
	var err error

	if dateStr != "" {
		date, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format, use YYYY-MM-DD"})
			return
		}
	} else {
		
		date = time.Now()
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

	
	scheduleData, err := h.scheduleService.GetDailySchedule(c.Request.Context(), userID.(string), instanceURL, date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	
	c.JSON(http.StatusOK, scheduleData)
}
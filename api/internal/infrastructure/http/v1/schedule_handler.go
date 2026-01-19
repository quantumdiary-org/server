package v1

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type ScheduleHandler struct{}

func NewScheduleHandler() *ScheduleHandler {
	return &ScheduleHandler{}
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

// GetWeeklySchedule возвращает расписание на неделю
// @Summary Get weekly schedule
// @Description Retrieves the weekly schedule for a student
// @Tags schedule
// @Security BearerAuth
// @Param week_start query string false "Start of week (YYYY-MM-DD format)" default(2023-09-01)
// @Produce json
// @Success 200 {object} WeeklySchedule
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /schedule/weekly [get]
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
		// По умолчанию берем начало текущей недели
		now := time.Now()
		weekStart = now.AddDate(0, 0, -int(now.Weekday())+1) // Понедельник недели
	}

	// Получаем userID из токена (через middleware)
	_, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// В реальной реализации здесь будет вызов к NetSchool API
	// для получения расписания
	// Пока возвращаем заглушку
	weekEnd := weekStart.AddDate(0, 0, 6) // Конец недели (воскресенье)

	mondayLessons := []Lesson{
		{
			ID:      "lesson_1",
			Number:  1,
			Subject: "Математика",
			Teacher: "Иванова А.А.",
			Room:    "301",
			Start:   time.Date(weekStart.Year(), weekStart.Month(), weekStart.Day(), 8, 30, 0, 0, time.UTC),
			End:     time.Date(weekStart.Year(), weekStart.Month(), weekStart.Day(), 9, 15, 0, 0, time.UTC),
			Date:    weekStart,
		},
		{
			ID:      "lesson_2",
			Number:  2,
			Subject: "Русский язык",
			Teacher: "Петрова Б.Б.",
			Room:    "302",
			Start:   time.Date(weekStart.Year(), weekStart.Month(), weekStart.Day(), 9, 30, 0, 0, time.UTC),
			End:     time.Date(weekStart.Year(), weekStart.Month(), weekStart.Day(), 10, 15, 0, 0, time.UTC),
			Date:    weekStart,
		},
	}

	tuesdayLessons := []Lesson{
		{
			ID:      "lesson_3",
			Number:  1,
			Subject: "Физика",
			Teacher: "Сидоров В.В.",
			Room:    "401",
			Start:   time.Date(weekStart.Year(), weekStart.Month(), weekStart.Day()+1, 8, 30, 0, 0, time.UTC),
			End:     time.Date(weekStart.Year(), weekStart.Month(), weekStart.Day()+1, 9, 15, 0, 0, time.UTC),
			Date:    weekStart.AddDate(0, 0, 1),
		},
	}

	days := []DaySchedule{
		{Date: weekStart, Lessons: mondayLessons},
		{Date: weekStart.AddDate(0, 0, 1), Lessons: tuesdayLessons},
		// Добавим остальные дни недели с пустыми расписаниями
		{Date: weekStart.AddDate(0, 0, 2), Lessons: []Lesson{}},
		{Date: weekStart.AddDate(0, 0, 3), Lessons: []Lesson{}},
		{Date: weekStart.AddDate(0, 0, 4), Lessons: []Lesson{}},
		{Date: weekStart.AddDate(0, 0, 5), Lessons: []Lesson{}},
		{Date: weekStart.AddDate(0, 0, 6), Lessons: []Lesson{}},
	}

	schedule := WeeklySchedule{
		WeekStart: weekStart,
		WeekEnd:   weekEnd,
		Days:      days,
	}

	c.JSON(http.StatusOK, schedule)
}

// GetDailySchedule возвращает расписание на день
// @Summary Get daily schedule
// @Description Retrieves the daily schedule for a student
// @Tags schedule
// @Security BearerAuth
// @Param date query string false "Date (YYYY-MM-DD format)" default(2023-09-01)
// @Produce json
// @Success 200 {object} DaySchedule
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /schedule/daily [get]
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
		// По умолчанию берем сегодняшнюю дату
		date = time.Now()
	}

	// Получаем userID из токена (через middleware)
	_, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// В реальной реализации здесь будет вызов к NetSchool API
	// для получения расписания на день
	// Пока возвращаем заглушку
	lessons := []Lesson{
		{
			ID:      "lesson_1",
			Number:  1,
			Subject: "Математика",
			Teacher: "Иванова А.А.",
			Room:    "301",
			Start:   time.Date(date.Year(), date.Month(), date.Day(), 8, 30, 0, 0, time.UTC),
			End:     time.Date(date.Year(), date.Month(), date.Day(), 9, 15, 0, 0, time.UTC),
			Date:    date,
		},
		{
			ID:      "lesson_2",
			Number:  2,
			Subject: "Русский язык",
			Teacher: "Петрова Б.Б.",
			Room:    "302",
			Start:   time.Date(date.Year(), date.Month(), date.Day(), 9, 30, 0, 0, time.UTC),
			End:     time.Date(date.Year(), date.Month(), date.Day(), 10, 15, 0, 0, time.UTC),
			Date:    date,
		},
	}

	daySchedule := DaySchedule{
		Date:    date,
		Lessons: lessons,
	}

	c.JSON(http.StatusOK, daySchedule)
}
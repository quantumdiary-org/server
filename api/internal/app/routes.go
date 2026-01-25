package app

import (
	"time"

	"github.com/gin-gonic/gin"
	"netschool-proxy/api/api/internal/api_types"
	"netschool-proxy/api/api/internal/domain/auth"
	"netschool-proxy/api/api/internal/domain/cache"
	"netschool-proxy/api/api/internal/domain/grade"
	"netschool-proxy/api/api/internal/domain/schedule"
	"netschool-proxy/api/api/internal/domain/student"
	"netschool-proxy/api/api/internal/infrastructure/http/v1"
	"netschool-proxy/api/api/internal/infrastructure/http/v1/middleware"
	"netschool-proxy/api/api/internal/pkg/security"
)


func setupRoutes(
	router *gin.Engine,
	authService *auth.Service,
	studentService *student.Service,
	gradeService *grade.Service,
	scheduleService *schedule.Service,
	cacheService cache.CacheStrategy,
	jwtService *security.JWTService,
	sessionRepo auth.SessionRepository,
	defaultAPIClient api_types.APIClientInterface,
) {
	
	authHandler := v1.NewAuthHandler(authService)
	healthHandler := v1.NewHealthHandler(defaultAPIClient, sessionRepo)
	studentHandler := v1.NewStudentHandler(studentService)
	gradeHandler := v1.NewGradeHandler(gradeService)
	scheduleHandler := v1.NewScheduleHandler(scheduleService, studentService)
	schoolHandler := v1.NewSchoolHandler(studentService)
	assignmentHandler := v1.NewAssignmentHandler(gradeService)

	
	authMiddleware := middleware.NewAuthMiddleware(authService, jwtService)
	cacheMiddleware := middleware.NewCacheMiddleware(cacheService)
	rateLimiter := middleware.NewRateLimiter(10, 20) 

	
	public := router.Group("/")
	{
		public.GET("/health/ping", healthHandler.Ping)
		public.GET("/health/intping", healthHandler.IntPing)
		public.GET("/health/full", healthHandler.FullHealth)
		public.POST("/auth/login", rateLimiter.RateLimitMiddleware(), authHandler.Login)
	}

	
	protected := router.Group("/api/v1")
	protected.Use(authMiddleware.AuthRequired())
	{
		
		protected.POST("/auth/logout", authHandler.Logout)

		
		protected.GET("/students/me", cacheMiddleware.CacheResponse(5*time.Minute), studentHandler.GetStudentInfo)
		protected.GET("/students/class", cacheMiddleware.CacheResponse(10*time.Minute), studentHandler.GetStudentsByClass)

		
		protected.GET("/grades", cacheMiddleware.CacheResponse(5*time.Minute), gradeHandler.GetGradesForStudent)
		protected.GET("/grades/subject", cacheMiddleware.CacheResponse(5*time.Minute), gradeHandler.GetGradesForSubject)

		
		protected.GET("/schedule/weekly", cacheMiddleware.CacheResponse(30*time.Minute), scheduleHandler.GetWeeklySchedule)
		protected.GET("/schedule/daily", cacheMiddleware.CacheResponse(30*time.Minute), scheduleHandler.GetDailySchedule)

		
		protected.GET("/school/info", cacheMiddleware.CacheResponse(1*time.Hour), schoolHandler.GetSchoolInfo)
		protected.GET("/school/classes", cacheMiddleware.CacheResponse(1*time.Hour), schoolHandler.GetClasses)

		
		protected.GET("/diary", cacheMiddleware.CacheResponse(1*time.Hour), scheduleHandler.GetWeeklySchedule)

		
		protected.GET("/assignments/detail", cacheMiddleware.CacheResponse(1*time.Hour), assignmentHandler.GetAssignment)
		protected.GET("/assignments/types", cacheMiddleware.CacheResponse(1*time.Hour), assignmentHandler.GetAssignmentTypes)

		
		protected.GET("/journal", cacheMiddleware.CacheResponse(10*time.Minute), gradeHandler.GetGradesForStudent)
		protected.GET("/journal/full", cacheMiddleware.CacheResponse(10*time.Minute), gradeHandler.GetGradesForSubject) 

		
		protected.GET("/info", cacheMiddleware.CacheResponse(1*time.Hour), studentHandler.GetStudentInfo)

		
		protected.GET("/student/photo", cacheMiddleware.CacheResponse(1*time.Hour), studentHandler.GetStudentPhoto)
	}

	

	
	admin := router.Group("/admin")
	admin.Use(authMiddleware.AuthRequired())
	{
		
		
		
	}
}
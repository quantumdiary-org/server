package app

import (
	"time"

	"github.com/gin-gonic/gin"
	"netschool-proxy/api/api/internal/api_types"
	"netschool-proxy/api/api/internal/domain/auth"
	"netschool-proxy/api/api/internal/domain/cache"
	"netschool-proxy/api/api/internal/domain/student"
	"netschool-proxy/api/api/internal/infrastructure/http/v1"
	"netschool-proxy/api/api/internal/infrastructure/http/v1/middleware"
	"netschool-proxy/api/api/internal/pkg/security"
)

// setupRoutes configures all application routes
func setupRoutes(
	router *gin.Engine,
	authService *auth.Service,
	studentService *student.Service,
	cacheService cache.CacheStrategy,
	jwtService *security.JWTService,
	sessionRepo auth.SessionRepository,
	defaultAPIClient api_types.APIClientInterface,
) {
	// Initialize handlers
	authHandler := v1.NewAuthHandler(authService)
	healthHandler := v1.NewHealthHandler(defaultAPIClient, sessionRepo)
	studentHandler := v1.NewStudentHandler(studentService)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(authService, jwtService)
	cacheMiddleware := middleware.NewCacheMiddleware(cacheService)
	rateLimiter := middleware.NewRateLimiter(10, 20) // 10 requests per second, burst of 20

	// Public routes
	public := router.Group("/")
	{
		public.GET("/health/ping", healthHandler.Ping)
		public.GET("/health/intping", healthHandler.IntPing)
		public.GET("/health/full", healthHandler.FullHealth)
		public.POST("/auth/login", rateLimiter.RateLimitMiddleware(), authHandler.Login)
	}

	// Protected routes
	protected := router.Group("/api/v1")
	protected.Use(authMiddleware.AuthRequired())
	{
		// Authentication routes
		protected.POST("/auth/logout", authHandler.Logout)

		// Student routes
		protected.GET("/students/me", cacheMiddleware.CacheResponse(5*time.Minute), studentHandler.GetStudentInfo)
		protected.GET("/students/class", cacheMiddleware.CacheResponse(10*time.Minute), studentHandler.GetStudentsByClass)

		// Add other protected routes here
		// For example:
		// protected.GET("/grades", gradeHandler.GetGradesForStudent)
		// protected.GET("/schedule/weekly", scheduleHandler.GetWeeklySchedule)
		// protected.GET("/school/info", schoolHandler.GetSchoolInfo)
	}

	// Add grade, schedule, and school handlers
	// gradeHandler := v1.NewGradeHandler(gradeService)
	// scheduleHandler := v1.NewScheduleHandler()
	// schoolHandler := v1.NewSchoolHandler()

	// Additional routes would be added here

	// Admin routes (optional)
	admin := router.Group("/admin")
	admin.Use(authMiddleware.AuthRequired())
	{
		// Only accessible to admin users
		// admin.GET("/stats", statsHandler.GetStats)
		// admin.GET("/users", usersHandler.GetAllUsers)
	}
}
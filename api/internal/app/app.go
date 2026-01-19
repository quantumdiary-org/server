package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"netschool-proxy/api/api/internal/api_types"
	"netschool-proxy/api/api/internal/config"
	"netschool-proxy/api/api/internal/domain/auth"
	"netschool-proxy/api/api/internal/domain/cache"
	"netschool-proxy/api/api/internal/domain/student"
	"netschool-proxy/api/api/internal/infrastructure/database"
	"netschool-proxy/api/api/internal/pkg/logger"
	"netschool-proxy/api/api/internal/pkg/security"
	infraCache "netschool-proxy/api/api/internal/infrastructure/cache"
)

// App represents the main application
type App struct {
	server *http.Server
	config *config.Config
	sessionRepo auth.SessionRepository
}

// New creates a new application instance
func New(cfg *config.Config) (*App, error) {
	// Initialize logger
	if err := logger.Init(cfg.Logging.Level, cfg.Logging.File); err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %w", err)
	}

	// Setup database connection
	dbConfig := database.DatabaseConfig{
		Type:     cfg.Database.Type,
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		Name:     cfg.Database.Name,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		SSLMode:  cfg.Database.SSLMode,
		URL:      cfg.Database.URL,
		SQLitePath: cfg.Database.SQLitePath,
	}

	dbManager := database.NewConnectionManager(dbConfig)
	db, err := dbManager.Connect()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Initialize API client factory
	apiFactory := &api_types.APIClientFactory{}
	apiConfig := api_types.APIConfig{
		Mode:       api_types.APIMode(cfg.NetSchool.Mode),
		Timeout:    int(cfg.NetSchool.Timeout.Seconds()),
		RetryMax:   cfg.NetSchool.RetryMax,
		RetryWait:  int(cfg.NetSchool.RetryWait.Milliseconds()),
	}

	// Initialize JWT service
	jwtService := security.NewJWTService(cfg.JWT.Secret, cfg.JWT.ExpiresIn)

	// Initialize session repository
	sessionRepo := database.NewSessionRepository(db)

	// Initialize auth service
	authService := auth.NewService(sessionRepo, apiFactory, apiConfig, jwtService)

	// Initialize cache
	var cacheService cache.CacheStrategy
	if cfg.Cache.Type == "redis" {
		redisCache, err := infraCache.NewRedisCacheService(cfg.Cache.RedisAddr)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to redis: %w", err)
		}
		cacheService = redisCache
	} else {
		cacheService = infraCache.NewMemoryCacheService(cfg.Cache.MemorySize)
	}

	// Initialize default API client for health checks
	defaultAPIClient, err := apiFactory.NewAPIClient(api_types.APIMode(cfg.NetSchool.Mode), apiConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create default API client: %w", err)
	}

	// Initialize services
	studentService := student.NewService(apiFactory, sessionRepo, apiConfig)

	// Create router
	router := gin.New()
	logger.Init(cfg.Logging.Level, cfg.Logging.File)
	router.Use(gin.Logger())

	// Setup routes - используем defaultAPIClient для health проверок
	setupRoutes(router, authService, studentService, cacheService, jwtService, sessionRepo, defaultAPIClient)

	// Create HTTP server
	server := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	return &App{
		server: server,
		config: cfg,
		sessionRepo: sessionRepo,
	}, nil
}

// Start starts the application
func (a *App) Start() error {
	// Start cleanup service in background
	cleanupService := auth.NewCleanupService(a.sessionRepo, 1*time.Hour)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go cleanupService.StartCleanup(ctx)

	// Start server
	logger.Info("Starting server", "port", a.config.Server.Port)
	
	go func() {
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", "error", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down server...")

	// Shutdown server gracefully
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	var err error
	err = a.server.Shutdown(shutdownCtx)
	if err != nil {
		logger.Error("Server forced to shutdown", "error", err)
		return err
	}

	logger.Info("Server exited")
	return nil
}
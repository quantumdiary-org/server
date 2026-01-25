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
	"netschool-proxy/api/api/internal/domain/grade"
	"netschool-proxy/api/api/internal/domain/schedule"
	"netschool-proxy/api/api/internal/domain/student"
	"netschool-proxy/api/api/internal/infrastructure/database"
	"netschool-proxy/api/api/internal/pkg/logger"
	"netschool-proxy/api/api/internal/pkg/security"
	infraCache "netschool-proxy/api/api/internal/infrastructure/cache"
)


type App struct {
	server *http.Server
	config *config.Config
	sessionRepo auth.SessionRepository
}


func New(cfg *config.Config) (*App, error) {
	
	if err := logger.Init(cfg.Logging.Level, cfg.Logging.File); err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %w", err)
	}

	
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

	
	apiFactory := &api_types.APIClientFactory{}
	apiConfig := api_types.APIConfig{
		Mode:       api_types.APIMode(cfg.NetSchool.Mode),
		Timeout:    int(cfg.NetSchool.Timeout.Seconds()),
		RetryMax:   cfg.NetSchool.RetryMax,
		RetryWait:  int(cfg.NetSchool.RetryWait.Milliseconds()),
	}

	
	jwtService := security.NewJWTService(cfg.JWT.Secret, cfg.JWT.ExpiresIn)

	
	sessionRepo := database.NewSessionRepository(db)

	
	authService := auth.NewService(sessionRepo, apiFactory, apiConfig, jwtService)

	
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

	
	defaultAPIClient, err := apiFactory.NewAPIClient(api_types.APIMode(cfg.NetSchool.Mode), apiConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create default API client: %w", err)
	}

	
	studentService := student.NewService(apiFactory, sessionRepo, apiConfig)
	gradeService := grade.NewService(apiFactory, sessionRepo, cacheService, apiConfig)
	scheduleService := schedule.NewService(apiFactory, sessionRepo, cacheService, apiConfig)

	
	router := gin.New()
	logger.Init(cfg.Logging.Level, cfg.Logging.File)
	router.Use(gin.Logger())

	
	setupRoutes(router, authService, studentService, gradeService, scheduleService, cacheService, jwtService, sessionRepo, defaultAPIClient)

	
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


func (a *App) Start() error {
	
	cleanupService := auth.NewCleanupService(a.sessionRepo, 1*time.Hour)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go cleanupService.StartCleanup(ctx)

	
	logger.Info("Starting server", "port", a.config.Server.Port)
	
	go func() {
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", "error", err)
		}
	}()

	
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down server...")

	
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
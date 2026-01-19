package v1

import (
	"context"
	"net/http"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"netschool-proxy/api/api/internal/api_types"
	"netschool-proxy/api/api/internal/domain/auth"
)

type HealthHandler struct {
	apiClient   api_types.APIClientInterface
	sessionRepo auth.SessionRepository
}

func NewHealthHandler(apiClient api_types.APIClientInterface, sessionRepo auth.SessionRepository) *HealthHandler {
	return &HealthHandler{
		apiClient:   apiClient,
		sessionRepo: sessionRepo,
	}
}

// Ping проверяет доступность прокси-сервера
// @Summary Check proxy server health
// @Description Checks if the proxy server is running
// @Tags health
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /health/ping [get]
func (h *HealthHandler) Ping(c *gin.Context) {
	// Сбор метрик производительности
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	c.JSON(http.StatusOK, gin.H{
		"status":    "ok",
		"version":   "1.0.0",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"service":   "netschool-proxy/api",
		"uptime":    time.Since(startTime).String(),
		"metrics": gin.H{
			"goroutines": runtime.NumGoroutine(),
			"alloc":      m.Alloc,
			"total_alloc": m.TotalAlloc,
			"sys":        m.Sys,
			"num_gc":     m.NumGC,
		},
	})
}

// startTime - время запуска сервиса
var startTime = time.Now()

// IntPing проверяет соединение с NetSchool API
// @Summary Check NetSchool API connectivity
// @Description Checks if the proxy can connect to NetSchool API
// @Tags health
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 503 {object} map[string]interface{}
// @Router /health/intping [get]
func (h *HealthHandler) IntPing(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	// Получаем instanceURL из параметров запроса или заголовка
	instanceURL := c.Query("instance_url")
	if instanceURL == "" {
		// Попробуем получить из заголовка
		instanceURL = c.GetHeader("X-Instance-URL")
		if instanceURL == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "instance_url is required"})
			return
		}
	}

	start := time.Now()
	// Проверяем доступность API
	_, err := h.apiClient.GetLoginData(ctx, instanceURL)
	duration := time.Since(start).Milliseconds()

	status := "available"
	cacheStatus := "fresh"
	lastCheck := time.Now().UTC()

	if err != nil {
		status = "unavailable"
		cacheStatus = "stale"
		// Здесь можно добавить логику проверки кэша
	}

	c.JSON(http.StatusOK, gin.H{
		"api_status":           status,
		"response_time_ms":     duration,
		"cache_status":         cacheStatus,
		"last_successful_call": lastCheck.Format(time.RFC3339),
		"can_use_cache":        status == "unavailable",
		"instance_url":         instanceURL,
	})
}

// FullHealth проверяет полное состояние системы
// @Summary Check full system health
// @Description Checks the health of all system components
// @Tags health
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /health/full [get]
func (h *HealthHandler) FullHealth(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	// Получаем instanceURL из параметров запроса или заголовка
	instanceURL := c.Query("instance_url")
	if instanceURL == "" {
		// Попробуем получить из заголовка
		instanceURL = c.GetHeader("X-Instance-URL")
		if instanceURL == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "instance_url is required"})
			return
		}
	}

	// Проверяем доступность API
	apiStart := time.Now()
	_, apiErr := h.apiClient.GetLoginData(ctx, instanceURL)
	apiDuration := time.Since(apiStart).Milliseconds()

	// Проверяем доступность базы данных
	dbStart := time.Now()
	_, dbErr := h.sessionRepo.GetByUserID(ctx, "health_check") // используем фиктивный ID для проверки
	dbDuration := time.Since(dbStart).Milliseconds()

	// Сбор метрик
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	// Формируем ответ
	response := gin.H{
		"status":    "ok",
		"version":   "1.0.0",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"service":   "netschool-proxy/api",
		"uptime":    time.Since(startTime).String(),
		"components": gin.H{
			"api": gin.H{
				"status":         getStatus(apiErr),
				"response_time":  apiDuration,
				"last_checked":   time.Now().UTC().Format(time.RFC3339),
			},
			"database": gin.H{
				"status":         getStatus(dbErr),
				"response_time":  dbDuration,
				"last_checked":   time.Now().UTC().Format(time.RFC3339),
			},
		},
		"metrics": gin.H{
			"goroutines": runtime.NumGoroutine(),
			"alloc":      m.Alloc,
			"total_alloc": m.TotalAlloc,
			"sys":        m.Sys,
			"num_gc":     m.NumGC,
		},
	}

	// Если есть ошибки, меняем статус на warning или error
	if apiErr != nil || dbErr != nil {
		response["status"] = "warning"
		if apiErr != nil {
			response["api_error"] = apiErr.Error()
		}
		if dbErr != nil && dbErr.Error() != "pg: no rows in result set" { // игнорируем ошибку отсутствия строки
			response["database_error"] = dbErr.Error()
		}
	}

	c.JSON(http.StatusOK, response)
}

// getStatus вспомогательная функция для определения статуса компонента
func getStatus(err error) string {
	if err != nil {
		return "unavailable"
	}
	return "available"
}
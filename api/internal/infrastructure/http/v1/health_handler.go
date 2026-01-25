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








func (h *HealthHandler) Ping(c *gin.Context) {
	
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	c.JSON(http.StatusOK, gin.H{
		"status":    "ok",
		"version":   "1.1.0",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"service":   "quantum-diary/api",
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


var startTime = time.Now()









func (h *HealthHandler) IntPing(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	
	instanceURL := c.Query("instance_url")
	if instanceURL == "" {
		
		instanceURL = c.GetHeader("X-Instance-URL")
		if instanceURL == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "instance_url is required"})
			return
		}
	}

	start := time.Now()
	
	_, err := h.apiClient.GetLoginData(ctx, instanceURL)
	duration := time.Since(start).Milliseconds()

	status := "available"
	cacheStatus := "fresh"
	lastCheck := time.Now().UTC()

	if err != nil {
		status = "unavailable"
		cacheStatus = "stale"
		
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








func (h *HealthHandler) FullHealth(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	
	instanceURL := c.Query("instance_url")
	if instanceURL == "" {
		
		instanceURL = c.GetHeader("X-Instance-URL")
		if instanceURL == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "instance_url is required"})
			return
		}
	}

	
	apiStart := time.Now()
	_, apiErr := h.apiClient.GetLoginData(ctx, instanceURL)
	apiDuration := time.Since(apiStart).Milliseconds()

	
	dbStart := time.Now()
	_, dbErr := h.sessionRepo.GetByUserID(ctx, "health_check") 
	dbDuration := time.Since(dbStart).Milliseconds()

	
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	
	response := gin.H{
		"status":    "ok",
		"version":   "1.1.0",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"service":   "quantum-diary/api",
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

	
	if apiErr != nil || dbErr != nil {
		response["status"] = "warning"
		if apiErr != nil {
			response["api_error"] = apiErr.Error()
		}
		if dbErr != nil && dbErr.Error() != "pg: no rows in result set" { 
			response["database_error"] = dbErr.Error()
		}
	}

	c.JSON(http.StatusOK, response)
}


func getStatus(err error) string {
	if err != nil {
		return "unavailable"
	}
	return "available"
}
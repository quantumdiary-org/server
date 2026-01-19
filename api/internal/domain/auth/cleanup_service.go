package auth

import (
	"context"
	"log"
	"time"

	"netschool-proxy/api/api/internal/pkg/logger"
)

// CleanupService manages automatic cleanup of expired sessions
type CleanupService struct {
	sessionRepo SessionRepository
	interval    time.Duration
}

// NewCleanupService creates a new cleanup service
func NewCleanupService(sessionRepo SessionRepository, interval time.Duration) *CleanupService {
	return &CleanupService{
		sessionRepo: sessionRepo,
		interval:    interval,
	}
}

// StartCleanup starts the periodic cleanup process
func (cs *CleanupService) StartCleanup(ctx context.Context) {
	ticker := time.NewTicker(cs.interval)
	defer ticker.Stop()

	// Perform initial cleanup
	cs.cleanupExpiredSessions()

	for {
		select {
		case <-ticker.C:
			cs.cleanupExpiredSessions()
		case <-ctx.Done():
			log.Println("Cleanup service shutting down...")
			return
		}
	}
}

// cleanupExpiredSessions removes all expired sessions from the database
func (cs *CleanupService) cleanupExpiredSessions() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := cs.sessionRepo.CleanupExpired(ctx); err != nil {
		logger.Error("Failed to cleanup expired sessions", "error", err)
	} else {
		logger.Info("Successfully cleaned up expired sessions")
	}
}

// ManualCleanup performs a one-time cleanup of expired sessions
func (cs *CleanupService) ManualCleanup() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	return cs.sessionRepo.CleanupExpired(ctx)
}

// GetCleanupStats returns statistics about cleanup operations
func (cs *CleanupService) GetCleanupStats() map[string]interface{} {
	// В реальной реализации здесь будет логика получения статистики
	// о количестве удаленных сессий и т.д.
	return map[string]interface{}{
		"cleanup_interval": cs.interval.String(),
		"last_run":        time.Now(),
		"status":          "active",
	}
}
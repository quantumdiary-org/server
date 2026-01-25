package auth

import (
	"context"
	"log"
	"time"

	"netschool-proxy/api/api/internal/pkg/logger"
)


type CleanupService struct {
	sessionRepo SessionRepository
	interval    time.Duration
}


func NewCleanupService(sessionRepo SessionRepository, interval time.Duration) *CleanupService {
	return &CleanupService{
		sessionRepo: sessionRepo,
		interval:    interval,
	}
}


func (cs *CleanupService) StartCleanup(ctx context.Context) {
	ticker := time.NewTicker(cs.interval)
	defer ticker.Stop()

	
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


func (cs *CleanupService) cleanupExpiredSessions() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := cs.sessionRepo.CleanupExpired(ctx); err != nil {
		logger.Error("Failed to cleanup expired sessions", "error", err)
	} else {
		logger.Info("Successfully cleaned up expired sessions")
	}
}


func (cs *CleanupService) ManualCleanup() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	return cs.sessionRepo.CleanupExpired(ctx)
}


func (cs *CleanupService) GetCleanupStats() map[string]interface{} {
	
	
	return map[string]interface{}{
		"cleanup_interval": cs.interval.String(),
		"last_run":        time.Now(),
		"status":          "active",
	}
}
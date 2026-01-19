package database

import (
	"context"
	"time"

	"gorm.io/gorm"
	"netschool-proxy/api/api/internal/domain/auth"
)

type SessionRepository struct {
	db *gorm.DB
}

func NewSessionRepository(db *gorm.DB) *SessionRepository {
	return &SessionRepository{db: db}
}

func (r *SessionRepository) Create(ctx context.Context, session *auth.NetSchoolSession) error {
	// Check if session already exists
	var existingSession auth.NetSchoolSession
	result := r.db.Where("user_id = ?", session.UserID).First(&existingSession)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			// Record doesn't exist, create new one
			return r.db.WithContext(ctx).Create(session).Error
		}
		return result.Error
	}

	// Record exists, update it
	session.ID = existingSession.ID // Preserve the ID
	return r.db.WithContext(ctx).Save(session).Error
}

func (r *SessionRepository) GetByUserID(ctx context.Context, userID string) (*auth.NetSchoolSession, error) {
	var session auth.NetSchoolSession
	result := r.db.WithContext(ctx).
		Where("user_id = ? AND expires_at > ?", userID, time.Now()).
		First(&session)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}

	return &session, nil
}

func (r *SessionRepository) Delete(ctx context.Context, userID string) error {
	return r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Delete(&auth.NetSchoolSession{}).Error
}

func (r *SessionRepository) CleanupExpired(ctx context.Context) error {
	return r.db.WithContext(ctx).
		Where("expires_at < ?", time.Now()).
		Delete(&auth.NetSchoolSession{}).Error
}
package services

import (
	"fmt"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

// AuditService handles security audit logging
type AuditService struct {
	logger *zap.Logger
	db     *gorm.DB
}

// NewAuditService creates a new audit service
func NewAuditService(logger *zap.Logger, db *gorm.DB) *AuditService {
	return &AuditService{
		logger: logger,
		db:     db,
	}
}

// AuditLog represents a security audit log entry
type AuditLog struct {
	ID            uint `gorm:"primarykey"`
	CreatedAt     time.Time
	UserID        uint   `json:"user_id"`
	UserType      string `json:"user_type"` // admin, distributor
	Action        string `json:"action"`
	Resource      string `json:"resource"`
	ResourceID    uint   `json:"resource_id,omitempty"`
	IP            string `json:"ip"`
	UserAgent     string `json:"user_agent"`
	Success       bool   `json:"success"`
	FailureReason string `json:"failure_reason,omitempty"`
	Metadata      string `json:"metadata,omitempty"` // JSON string for additional data
}

// LogCardAccess logs card access attempts
func (s *AuditService) LogCardAccess(userID uint, userType string, cardID uint, action string, ip string, userAgent string, success bool, failureReason string) {
	log := &AuditLog{
		UserID:        userID,
		UserType:      userType,
		Action:        action,
		Resource:      "card",
		ResourceID:    cardID,
		IP:            ip,
		UserAgent:     userAgent,
		Success:       success,
		FailureReason: failureReason,
		CreatedAt:     time.Now(),
	}

	if err := s.db.Create(log).Error; err != nil {
		s.logger.Error("failed to create audit log", zap.Error(err))
	}
}

// LogLoginAttempt logs login attempts
func (s *AuditService) LogLoginAttempt(username string, userType string, ip string, userAgent string, success bool) {
	log := &AuditLog{
		UserType:  userType,
		Action:    "login",
		Resource:  "auth",
		IP:        ip,
		UserAgent: userAgent,
		Success:   success,
		Metadata:  fmt.Sprintf(`{"username":"%s"}`, username),
		CreatedAt: time.Now(),
	}

	if !success {
		log.FailureReason = "invalid credentials"
	}

	if err := s.db.Create(log).Error; err != nil {
		s.logger.Error("failed to create audit log", zap.Error(err))
	}
}

// GetFailedLoginAttempts gets the number of failed login attempts in the last duration
func (s *AuditService) GetFailedLoginAttempts(username string, userType string, duration time.Duration) (int64, error) {
	var count int64
	since := time.Now().Add(-duration)

	err := s.db.Model(&AuditLog{}).
		Where("user_type = ? AND action = ? AND success = ? AND created_at > ? AND metadata LIKE ?",
			userType, "login", false, since, fmt.Sprintf(`%%"username":"%s"%%`, username)).
		Count(&count).Error

	return count, err
}

// LogAPIAccess logs API access
func (s *AuditService) LogAPIAccess(userID uint, userType string, endpoint string, method string, ip string, statusCode int) {
	log := &AuditLog{
		UserID:    userID,
		UserType:  userType,
		Action:    fmt.Sprintf("%s %s", method, endpoint),
		Resource:  "api",
		IP:        ip,
		Success:   statusCode >= 200 && statusCode < 300,
		Metadata:  fmt.Sprintf(`{"status_code":%d}`, statusCode),
		CreatedAt: time.Now(),
	}

	if err := s.db.Create(log).Error; err != nil {
		s.logger.Error("failed to create audit log", zap.Error(err))
	}
}

// GetSuspiciousActivity gets suspicious activity patterns
func (s *AuditService) GetSuspiciousActivity() ([]*AuditLog, error) {
	var logs []*AuditLog

	// Look for patterns like:
	// - Multiple failed login attempts
	// - Unusual access patterns
	// - Access from multiple IPs in short time
	err := s.db.
		Where("success = ? AND created_at > ?", false, time.Now().Add(-24*time.Hour)).
		Order("created_at desc").
		Limit(100).
		Find(&logs).Error

	return logs, err
}

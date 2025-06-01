package config

// SecurityConfig holds security-related configuration
type SecurityConfig struct {
	// MaskCardCodes whether to mask card codes in responses
	MaskCardCodes bool
	// LogCardCodes whether to log card codes (should be false in production)
	LogCardCodes bool
	// MaxLoginAttempts maximum login attempts before account lockout
	MaxLoginAttempts int
	// LockoutDuration duration in minutes for account lockout
	LockoutDuration int
	// RateLimitPerMinute API rate limit per minute
	RateLimitPerMinute int
	// EnableAuditLog whether to enable audit logging
	EnableAuditLog bool
}

// GetSecurityConfig returns the security configuration
func GetSecurityConfig() *SecurityConfig {
	return &SecurityConfig{
		MaskCardCodes:      true,  // Always mask in production
		LogCardCodes:       false, // Never log card codes
		MaxLoginAttempts:   5,
		LockoutDuration:    30, // 30 minutes
		RateLimitPerMinute: 60,
		EnableAuditLog:     true,
	}
}
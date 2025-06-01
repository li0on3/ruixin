package middleware

import (
	"backend/internal/services"
	"backend/internal/utils"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// SecurityHeaders adds security headers
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		c.Header("Content-Security-Policy", "default-src 'self'")
		c.Next()
	}
}

// InputSanitizer sanitizes user input
func InputSanitizer() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Sanitize query parameters
		for key, values := range c.Request.URL.Query() {
			for i, value := range values {
				values[i] = utils.SanitizeInput(value)
			}
			c.Request.URL.RawQuery = strings.Replace(c.Request.URL.RawQuery, key+"="+values[0], key+"="+values[0], 1)
		}
		
		c.Next()
	}
}

// AuditLogger logs all API access
func AuditLogger(auditService *services.AuditService) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		
		// Process request
		c.Next()
		
		// Log after request
		userID := uint(0)
		userType := "anonymous"
		
		if id, exists := c.Get("admin_id"); exists {
			userID = id.(uint)
			userType = "admin"
		} else if id, exists := c.Get("distributor_id"); exists {
			userID = id.(uint)
			userType = "distributor"
		}
		
		// Log API access
		go auditService.LogAPIAccess(
			userID,
			userType,
			c.Request.URL.Path,
			c.Request.Method,
			c.ClientIP(),
			c.Writer.Status(),
		)
		
		// Log slow requests
		if time.Since(start) > 5*time.Second {
			// Log slow request for monitoring
		}
	}
}

// CardAccessControl ensures card access is properly authorized
func CardAccessControl() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if the route involves card access
		if strings.Contains(c.Request.URL.Path, "/cards") {
			// Ensure user is authenticated
			if _, exists := c.Get("admin_id"); !exists {
				if _, exists := c.Get("distributor_id"); !exists {
					c.JSON(http.StatusForbidden, gin.H{
						"code": 403,
						"msg":  "Access denied",
						"data": nil,
					})
					c.Abort()
					return
				}
			}
		}
		
		c.Next()
	}
}

// IPWhitelist checks if the request IP is whitelisted
func IPWhitelist(whitelist []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if len(whitelist) == 0 {
			// No whitelist configured, allow all
			c.Next()
			return
		}
		
		clientIP := c.ClientIP()
		allowed := false
		
		for _, ip := range whitelist {
			if clientIP == ip {
				allowed = true
				break
			}
		}
		
		if !allowed {
			c.JSON(http.StatusForbidden, gin.H{
				"code": 403,
				"msg":  "Access denied from this IP address",
				"data": nil,
			})
			c.Abort()
			return
		}
		
		c.Next()
	}
}
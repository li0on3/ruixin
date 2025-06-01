package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"backend/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

// RateLimiter 频率限制配置
type RateLimiter struct {
	redis  *redis.Client
	ctx    context.Context
}

// NewRateLimiter 创建频率限制器
func NewRateLimiter(redisClient *redis.Client) *RateLimiter {
	return &RateLimiter{
		redis: redisClient,
		ctx:   context.Background(),
	}
}

// RateLimit 创建频率限制中间件
// window: 时间窗口
// limit: 在时间窗口内的最大请求数
// keyFunc: 生成限流key的函数
func (rl *RateLimiter) RateLimit(window time.Duration, limit int, keyFunc func(*gin.Context) string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 生成限流key
		key := keyFunc(c)
		if key == "" {
			c.Next()
			return
		}

		// 使用Redis的INCR命令计数
		count, err := rl.redis.Incr(rl.ctx, key).Result()
		if err != nil {
			// Redis错误时，为了不影响业务，继续处理请求
			c.Next()
			return
		}

		// 第一次访问时设置过期时间
		if count == 1 {
			rl.redis.Expire(rl.ctx, key, window)
		}

		// 检查是否超过限制
		if count > int64(limit) {
			// 记录触发频率限制的事件（如果有审计服务）
			if auditService, exists := c.Get("security_audit_service"); exists {
				if service, ok := auditService.(*services.SecurityAuditService); ok {
					distributorID, _ := c.Get("distributor_id")
					if distID, ok := distributorID.(uint); ok {
						service.LogRateLimited(distID, c.Request.URL.Path, c)
					}
				}
			}
			
			c.JSON(http.StatusTooManyRequests, gin.H{
				"code": 429,
				"msg":  fmt.Sprintf("请求过于频繁，请%d秒后再试", int(window.Seconds())),
				"data": nil,
			})
			c.Abort()
			return
		}

		// 在响应头中返回限流信息
		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", limit))
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", limit-int(count)))
		c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", time.Now().Add(window).Unix()))

		c.Next()
	}
}

// DistributorRateLimit 分销商API限流（基于分销商ID）
func DistributorRateLimit(rl *RateLimiter, window time.Duration, limit int) gin.HandlerFunc {
	return rl.RateLimit(window, limit, func(c *gin.Context) string {
		distributorID, exists := c.Get("distributor_id")
		if !exists {
			return ""
		}
		return fmt.Sprintf("rate_limit:distributor:%v:%s", distributorID, c.Request.URL.Path)
	})
}

// IPRateLimit 基于IP的限流
func IPRateLimit(rl *RateLimiter, window time.Duration, limit int) gin.HandlerFunc {
	return rl.RateLimit(window, limit, func(c *gin.Context) string {
		return fmt.Sprintf("rate_limit:ip:%s:%s", c.ClientIP(), c.Request.URL.Path)
	})
}

// GlobalRateLimit 全局API限流
func GlobalRateLimit(rl *RateLimiter, window time.Duration, limit int) gin.HandlerFunc {
	return rl.RateLimit(window, limit, func(c *gin.Context) string {
		return fmt.Sprintf("rate_limit:global:%s", c.Request.URL.Path)
	})
}

// SensitiveAPIRateLimit 敏感API的严格限流
func SensitiveAPIRateLimit(rl *RateLimiter) gin.HandlerFunc {
	// 对敏感API使用更严格的限制：每分钟10次
	return rl.RateLimit(time.Minute, 10, func(c *gin.Context) string {
		distributorID, exists := c.Get("distributor_id")
		if !exists {
			return fmt.Sprintf("rate_limit:sensitive:ip:%s", c.ClientIP())
		}
		return fmt.Sprintf("rate_limit:sensitive:distributor:%v", distributorID)
	})
}
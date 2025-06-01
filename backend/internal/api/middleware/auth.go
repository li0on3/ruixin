package middleware

import (
	"net/http"
	"strings"

	"backend/internal/models"
	"backend/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// DistributorAuth 分销商API认证中间件
func DistributorAuth(distributorRepo *repository.DistributorRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从Header获取API Key
		apiKey := c.GetHeader("X-API-Key")
		apiSecret := c.GetHeader("X-API-Secret")

		if apiKey == "" || apiSecret == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "Missing API credentials",
				"data": nil,
			})
			c.Abort()
			return
		}

		// 验证API Key
		distributor, err := distributorRepo.GetByAPIKey(apiKey)
		if err != nil || distributor.APISecret != apiSecret {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "Invalid API credentials",
				"data": nil,
			})
			c.Abort()
			return
		}

		// 检查分销商状态
		if distributor.Status != 1 {
			c.JSON(http.StatusForbidden, gin.H{
				"code": 403,
				"msg":  "Distributor account is disabled",
				"data": nil,
			})
			c.Abort()
			return
		}

		// 记录API调用日志
		go distributorRepo.LogAPICall(&models.DistributorAPILog{
			DistributorID: distributor.ID,
			APIEndpoint:   c.Request.URL.Path,
			Method:        c.Request.Method,
			IPAddress:     c.ClientIP(),
			UserAgent:     c.GetHeader("User-Agent"),
		})

		// 设置分销商信息到上下文
		c.Set("distributor_id", distributor.ID)
		c.Set("distributor", distributor)

		c.Next()
	}
}

// AdminAuth 管理员认证中间件
func AdminAuth(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从Header获取Token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "Missing authorization header",
				"data": nil,
			})
			c.Abort()
			return
		}

		// 解析Token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "无效的令牌",
				"data": nil,
			})
			c.Abort()
			return
		}

		// 获取Claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "Invalid token claims",
				"data": nil,
			})
			c.Abort()
			return
		}

		// 设置管理员信息到上下文
		c.Set("admin_id", uint(claims["admin_id"].(float64)))
		c.Set("admin_role", claims["role"].(string))

		c.Next()
	}
}

// RoleRequired 角色权限中间件
func RoleRequired(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		adminRole, exists := c.Get("admin_role")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{
				"code": 403,
				"msg":  "Access denied",
				"data": nil,
			})
			c.Abort()
			return
		}

		// 检查角色权限
		roleStr := adminRole.(string)
		authorized := false
		for _, role := range roles {
			if roleStr == role {
				authorized = true
				break
			}
		}

		if !authorized {
			c.JSON(http.StatusForbidden, gin.H{
				"code": 403,
				"msg":  "Insufficient permissions",
				"data": nil,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

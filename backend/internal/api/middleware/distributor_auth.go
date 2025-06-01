package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// DistributorJWTAuth 分销商JWT认证中间件
func DistributorJWTAuth(jwtSecret string) gin.HandlerFunc {
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

		// 验证token类型
		if claims["type"] != "distributor" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "Invalid token type",
				"data": nil,
			})
			c.Abort()
			return
		}

		// 设置分销商信息到上下文
		c.Set("distributor_id", uint(claims["distributor_id"].(float64)))
		c.Set("distributor_email", claims["email"].(string))

		c.Next()
	}
}
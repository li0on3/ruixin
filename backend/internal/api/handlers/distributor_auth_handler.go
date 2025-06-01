package handlers

import (
	"net/http"
	"time"

	"backend/internal/config"
	"backend/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type DistributorAuthHandler struct {
	distributorService *services.DistributorService
	jwtConfig          *config.JWTConfig
}

func NewDistributorAuthHandler(distributorService *services.DistributorService, jwtConfig *config.JWTConfig) *DistributorAuthHandler {
	return &DistributorAuthHandler{
		distributorService: distributorService,
		jwtConfig:          jwtConfig,
	}
}

type DistributorLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// Login 分销商登录
func (h *DistributorAuthHandler) Login(c *gin.Context) {
	var req DistributorLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "请求参数无效",
			"data": nil,
		})
		return
	}

	// 查找分销商
	distributor, err := h.distributorService.GetByEmail(req.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": 401,
			"msg":  "Invalid email or password",
			"data": nil,
		})
		return
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(distributor.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": 401,
			"msg":  "Invalid email or password",
			"data": nil,
		})
		return
	}

	// 检查状态
	if distributor.Status != 1 {
		c.JSON(http.StatusForbidden, gin.H{
			"code": 403,
			"msg":  "Account is disabled",
			"data": nil,
		})
		return
	}

	// 生成JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"distributor_id": distributor.ID,
		"email":          distributor.Email,
		"type":           "distributor",
		"exp":            time.Now().Add(time.Duration(h.jwtConfig.Expiration) * time.Second).Unix(),
	})

	tokenString, err := token.SignedString([]byte(h.jwtConfig.Secret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "Failed to generate token",
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "登录成功",
		"data": gin.H{
			"token": tokenString,
			"distributor": gin.H{
				"id":           distributor.ID,
				"name":         distributor.Name,
				"email":        distributor.Email,
				"company_name": distributor.CompanyName,
				"api_key":      distributor.APIKey,
				"api_secret":   distributor.APISecret,
			},
		},
	})
}

// Logout 分销商登出
func (h *DistributorAuthHandler) Logout(c *gin.Context) {
	// 客户端清除token即可
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "登出成功",
		"data": nil,
	})
}

// GetProfile 获取当前分销商信息
func (h *DistributorAuthHandler) GetProfile(c *gin.Context) {
	distributorID, exists := c.Get("distributor_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": 401,
			"msg":  "未授权",
			"data": nil,
		})
		return
	}

	distributor, err := h.distributorService.GetByID(distributorID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "分销商不存在",
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "成功",
		"data": gin.H{
			"id":                 distributor.ID,
			"name":               distributor.Name,
			"company_name":       distributor.CompanyName,
			"contact_name":       distributor.ContactName,
			"phone":              distributor.Phone,
			"email":              distributor.Email,
			"api_key":            distributor.APIKey,
			"api_secret":         distributor.APISecret,
			"status":             distributor.Status,
			"balance":            distributor.Balance,
			"credit_limit":       distributor.CreditLimit,
			"callback_url":       distributor.CallbackURL,
			"daily_order_limit":  distributor.DailyOrderLimit,
			"monthly_order_limit": distributor.MonthlyOrderLimit,
		},
	})
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

// ChangePassword 修改密码
func (h *DistributorAuthHandler) ChangePassword(c *gin.Context) {
	distributorID, exists := c.Get("distributor_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": 401,
			"msg":  "未授权",
			"data": nil,
		})
		return
	}

	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "请求参数无效",
			"data": nil,
		})
		return
	}

	distributor, err := h.distributorService.GetByID(distributorID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "分销商不存在",
			"data": nil,
		})
		return
	}

	// 验证旧密码
	if err := bcrypt.CompareHashAndPassword([]byte(distributor.Password), []byte(req.OldPassword)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "Old password is incorrect",
			"data": nil,
		})
		return
	}

	// 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "Failed to encrypt password",
			"data": nil,
		})
		return
	}

	// 更新密码
	distributor.Password = string(hashedPassword)
	if err := h.distributorService.Update(distributor); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "Failed to update password",
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "Password changed successfully",
		"data": nil,
	})
}
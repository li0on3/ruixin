package handlers

import (
	"net/http"
	"time"

	"backend/internal/config"
	"backend/internal/models"
	"backend/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AdminAuthHandler struct {
	adminService  *services.AdminService
	jwtSecret     string
	jwtExpiration int
}

func NewAdminAuthHandler(adminService *services.AdminService, cfg *config.JWTConfig) *AdminAuthHandler {
	return &AdminAuthHandler{
		adminService:  adminService,
		jwtSecret:     cfg.Secret,
		jwtExpiration: cfg.Expiration,
	}
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Login 管理员登录
func (h *AdminAuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "请求参数无效",
			"data": nil,
		})
		return
	}

	// 查找管理员
	admin, err := h.adminService.GetByUsername(req.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": 401,
			"msg":  "Invalid username or password",
			"data": nil,
		})
		return
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": 401,
			"msg":  "Invalid username or password",
			"data": nil,
		})
		return
	}

	// 检查账号状态
	if admin.Status != 1 {
		c.JSON(http.StatusForbidden, gin.H{
			"code": 403,
			"msg":  "Account is disabled",
			"data": nil,
		})
		return
	}

	// 生成JWT token
	token, err := h.generateToken(admin)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "Failed to generate token",
			"data": nil,
		})
		return
	}

	// 更新登录信息
	h.adminService.UpdateLoginInfo(admin.ID, c.ClientIP())

	// 记录操作日志
	h.adminService.LogOperation(&models.AdminOperationLog{
		AdminID:   admin.ID,
		AdminName: admin.Username,
		Operation: "login",
		Module:    "auth",
		Details:   "Admin login successful",
		IPAddress: c.ClientIP(),
		UserAgent: c.GetHeader("User-Agent"),
	})

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "登录成功",
		"data": gin.H{
			"token": token,
			"user": gin.H{
				"id":         admin.ID,
				"username":   admin.Username,
				"email":      admin.Email,
				"real_name":  admin.RealName,
				"role":       admin.Role,
				"created_at": admin.CreatedAt,
			},
		},
	})
}

// Logout 管理员登出
func (h *AdminAuthHandler) Logout(c *gin.Context) {
	adminID, _ := c.Get("admin_id")

	// 记录操作日志
	h.adminService.LogOperation(&models.AdminOperationLog{
		AdminID:   adminID.(uint),
		Operation: "logout",
		Module:    "auth",
		Details:   "Admin logout",
		IPAddress: c.ClientIP(),
		UserAgent: c.GetHeader("User-Agent"),
	})

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "登出成功",
		"data": nil,
	})
}

// GetUserInfo 获取当前用户信息
func (h *AdminAuthHandler) GetUserInfo(c *gin.Context) {
	adminID, _ := c.Get("admin_id")

	admin, err := h.adminService.GetByID(adminID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "User not found",
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "成功",
		"data": gin.H{
			"id":            admin.ID,
			"username":      admin.Username,
			"email":         admin.Email,
			"phone":         admin.Phone,
			"real_name":     admin.RealName,
			"role":          admin.Role,
			"status":        admin.Status,
			"last_login_at": admin.LastLoginAt,
			"created_at":    admin.CreatedAt,
		},
	})
}

// ChangePassword 修改密码
func (h *AdminAuthHandler) ChangePassword(c *gin.Context) {
	adminID, _ := c.Get("admin_id")

	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "请求参数无效",
			"data": nil,
		})
		return
	}

	admin, err := h.adminService.GetByID(adminID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "Failed to get user info",
			"data": nil,
		})
		return
	}

	// 验证旧密码
	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(req.OldPassword)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "Invalid old password",
			"data": nil,
		})
		return
	}

	// 更新密码
	if err := h.adminService.ChangePassword(admin.ID, req.NewPassword); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "Failed to change password",
			"data": nil,
		})
		return
	}

	// 记录操作日志
	h.adminService.LogOperation(&models.AdminOperationLog{
		AdminID:   admin.ID,
		AdminName: admin.Username,
		Operation: "change_password",
		Module:    "auth",
		Details:   "Password changed",
		IPAddress: c.ClientIP(),
		UserAgent: c.GetHeader("User-Agent"),
	})

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "Password changed successfully",
		"data": nil,
	})
}

func (h *AdminAuthHandler) generateToken(admin *models.Admin) (string, error) {
	claims := jwt.MapClaims{
		"admin_id": admin.ID,
		"username": admin.Username,
		"role":     admin.Role,
		"exp":      time.Now().Add(time.Duration(h.jwtExpiration) * time.Second).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(h.jwtSecret))
}

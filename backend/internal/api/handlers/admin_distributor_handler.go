package handlers

import (
	"net/http"
	"strconv"

	"backend/internal/models"
	"backend/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AdminDistributorHandler struct {
	distributorService *services.DistributorService
	adminService       *services.AdminService
}

func NewAdminDistributorHandler(distributorService *services.DistributorService, adminService *services.AdminService) *AdminDistributorHandler {
	return &AdminDistributorHandler{
		distributorService: distributorService,
		adminService:       adminService,
	}
}

// DistributorWithAPIKeys 包含API密钥的分销商信息
type DistributorWithAPIKeys struct {
	*models.Distributor
	APIKey    string `json:"api_key"`
	APISecret string `json:"api_secret"`
}

// ListDistributors 获取分销商列表
func (h *AdminDistributorHandler) ListDistributors(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	status, _ := strconv.Atoi(c.Query("status"))
	name := c.Query("name")

	offset := (page - 1) * pageSize
	filters := make(map[string]interface{})

	if status > 0 {
		filters["status"] = status - 1
	}
	if name != "" {
		filters["name"] = name
	}

	distributors, total, err := h.distributorService.List(offset, pageSize, filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "Failed to get distributor list",
			"data": nil,
		})
		return
	}

	// 转换为包含API密钥的响应结构
	var distributorsWithKeys []DistributorWithAPIKeys
	for _, d := range distributors {
		distributorsWithKeys = append(distributorsWithKeys, DistributorWithAPIKeys{
			Distributor: d,
			APIKey:      d.APIKey,
			APISecret:   d.APISecret,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "成功",
		"data": gin.H{
			"list":      distributorsWithKeys,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetDistributor 获取分销商详情
func (h *AdminDistributorHandler) GetDistributor(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的分销商ID",
			"data": nil,
		})
		return
	}

	distributor, err := h.distributorService.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "分销商不存在",
			"data": nil,
		})
		return
	}

	// 包含API密钥的响应
	distributorWithKeys := DistributorWithAPIKeys{
		Distributor: distributor,
		APIKey:      distributor.APIKey,
		APISecret:   distributor.APISecret,
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "成功",
		"data": distributorWithKeys,
	})
}

type CreateDistributorRequest struct {
	Name              string  `json:"name" binding:"required"`
	CompanyName       string  `json:"company_name"`
	ContactName       string  `json:"contact_name"`
	Phone             string  `json:"phone"`
	Email             string  `json:"email" binding:"required,email"`
	Password          string  `json:"password"`
	CallbackURL       string  `json:"callback_url"`
	CreditLimit       float64 `json:"credit_limit"`
	DailyOrderLimit   int     `json:"daily_order_limit"`
	MonthlyOrderLimit int     `json:"monthly_order_limit"`
}

// CreateDistributor 创建分销商
func (h *AdminDistributorHandler) CreateDistributor(c *gin.Context) {
	var req CreateDistributorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "请求参数无效",
			"data": nil,
		})
		return
	}

	// 生成API密钥
	apiKey := uuid.New().String()
	apiSecret := uuid.New().String()

	// 处理密码
	password := req.Password
	if password == "" {
		password = "123456" // 默认密码
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "Failed to encrypt password",
			"data": nil,
		})
		return
	}

	distributor := &models.Distributor{
		Name:              req.Name,
		CompanyName:       req.CompanyName,
		ContactName:       req.ContactName,
		Phone:             req.Phone,
		Email:             req.Email,
		Password:          string(hashedPassword),
		APIKey:            apiKey,
		APISecret:         apiSecret,
		Status:            1, // 默认正常
		Balance:           0,
		CreditLimit:       req.CreditLimit,
		CallbackURL:       req.CallbackURL,
		DailyOrderLimit:   req.DailyOrderLimit,
		MonthlyOrderLimit: req.MonthlyOrderLimit,
	}

	if err := h.distributorService.Create(distributor); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	// 记录操作日志
	adminID, _ := c.Get("admin_id")
	h.adminService.LogOperation(&models.AdminOperationLog{
		AdminID:   adminID.(uint),
		Operation: "create_distributor",
		Module:    "distributor",
		Details:   "Created distributor: " + distributor.Name,
		IPAddress: c.ClientIP(),
		UserAgent: c.GetHeader("User-Agent"),
	})

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "分销商创建成功",
		"data": gin.H{
			"id":              distributor.ID,
			"api_key":         apiKey,
			"api_secret":      apiSecret,
			"default_password": password, // 返回初始密码供管理员告知分销商
		},
	})
}

type UpdateDistributorRequest struct {
	Name              *string  `json:"name"`
	CompanyName       *string  `json:"company_name"`
	ContactName       *string  `json:"contact_name"`
	Phone             *string  `json:"phone"`
	Email             *string  `json:"email"`
	Status            *int     `json:"status"`
	CreditLimit       *float64 `json:"credit_limit"`
	CallbackURL       *string  `json:"callback_url"`
	DailyOrderLimit   *int     `json:"daily_order_limit"`
	MonthlyOrderLimit *int     `json:"monthly_order_limit"`
}

// UpdateDistributor 更新分销商
func (h *AdminDistributorHandler) UpdateDistributor(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的分销商ID",
			"data": nil,
		})
		return
	}

	var req UpdateDistributorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "请求参数无效",
			"data": nil,
		})
		return
	}

	distributor, err := h.distributorService.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "分销商不存在",
			"data": nil,
		})
		return
	}

	// 更新字段
	if req.Name != nil {
		distributor.Name = *req.Name
	}
	if req.CompanyName != nil {
		distributor.CompanyName = *req.CompanyName
	}
	if req.ContactName != nil {
		distributor.ContactName = *req.ContactName
	}
	if req.Phone != nil {
		distributor.Phone = *req.Phone
	}
	if req.Email != nil {
		distributor.Email = *req.Email
	}
	if req.Status != nil {
		distributor.Status = *req.Status
	}
	if req.CreditLimit != nil {
		distributor.CreditLimit = *req.CreditLimit
	}
	if req.CallbackURL != nil {
		distributor.CallbackURL = *req.CallbackURL
	}
	if req.DailyOrderLimit != nil {
		distributor.DailyOrderLimit = *req.DailyOrderLimit
	}
	if req.MonthlyOrderLimit != nil {
		distributor.MonthlyOrderLimit = *req.MonthlyOrderLimit
	}

	if err := h.distributorService.Update(distributor); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "Failed to update distributor",
			"data": nil,
		})
		return
	}

	// 记录操作日志
	adminID, _ := c.Get("admin_id")
	h.adminService.LogOperation(&models.AdminOperationLog{
		AdminID:   adminID.(uint),
		Operation: "update_distributor",
		Module:    "distributor",
		Details:   "Updated distributor: " + distributor.Name,
		IPAddress: c.ClientIP(),
		UserAgent: c.GetHeader("User-Agent"),
	})

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "分销商更新成功",
		"data": distributor,
	})
}

// DeleteDistributor 删除分销商
func (h *AdminDistributorHandler) DeleteDistributor(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的分销商ID",
			"data": nil,
		})
		return
	}

	distributor, err := h.distributorService.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "分销商不存在",
			"data": nil,
		})
		return
	}

	if err := h.distributorService.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	// 记录操作日志
	adminID, _ := c.Get("admin_id")
	h.adminService.LogOperation(&models.AdminOperationLog{
		AdminID:   adminID.(uint),
		Operation: "delete_distributor",
		Module:    "distributor",
		Details:   "Deleted distributor: " + distributor.Name,
		IPAddress: c.ClientIP(),
		UserAgent: c.GetHeader("User-Agent"),
	})

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "分销商删除成功",
		"data": nil,
	})
}

// ResetAPIKey 重置API密钥
func (h *AdminDistributorHandler) ResetAPIKey(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的分销商ID",
			"data": nil,
		})
		return
	}

	distributor, err := h.distributorService.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "分销商不存在",
			"data": nil,
		})
		return
	}

	// 生成新的API密钥
	apiKey := uuid.New().String()
	apiSecret := uuid.New().String()

	distributor.APIKey = apiKey
	distributor.APISecret = apiSecret

	if err := h.distributorService.Update(distributor); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "Failed to reset API key",
			"data": nil,
		})
		return
	}

	// 记录操作日志
	adminID, _ := c.Get("admin_id")
	h.adminService.LogOperation(&models.AdminOperationLog{
		AdminID:   adminID.(uint),
		Operation: "reset_api_key",
		Module:    "distributor",
		Details:   "Reset API key for distributor: " + distributor.Name,
		IPAddress: c.ClientIP(),
		UserAgent: c.GetHeader("User-Agent"),
	})

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "API密钥重置成功",
		"data": gin.H{
			"api_key":    apiKey,
			"api_secret": apiSecret,
		},
	})
}

// ResetPassword 重置分销商密码
func (h *AdminDistributorHandler) ResetPassword(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的分销商ID",
			"data": nil,
		})
		return
	}

	var req struct {
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "请求参数无效",
			"data": nil,
		})
		return
	}

	distributor, err := h.distributorService.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "分销商不存在",
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
			"msg":  "Failed to reset password",
			"data": nil,
		})
		return
	}

	// 记录操作日志
	adminID, _ := c.Get("admin_id")
	h.adminService.LogOperation(&models.AdminOperationLog{
		AdminID:   adminID.(uint),
		Operation: "reset_distributor_password",
		Module:    "distributor",
		Details:   "Reset password for distributor: " + distributor.Name,
		IPAddress: c.ClientIP(),
		UserAgent: c.GetHeader("User-Agent"),
	})

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "密码重置成功",
		"data": nil,
	})
}

// GetAPILogs 获取API调用日志
func (h *AdminDistributorHandler) GetAPILogs(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的分销商ID",
			"data": nil,
		})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	offset := (page - 1) * pageSize

	logs, total, err := h.distributorService.GetAPILogs(uint(id), offset, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "Failed to get API logs",
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "成功",
		"data": gin.H{
			"list":      logs,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

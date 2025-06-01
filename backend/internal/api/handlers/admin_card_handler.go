package handlers

import (
	"net/http"
	"strconv"
	"time"

	"backend/internal/models"
	"backend/internal/services"
	"github.com/gin-gonic/gin"
)

type AdminCardHandler struct {
	cardService         *services.CardService
	adminService        *services.AdminService
	productService      *services.ProductService
	systemConfigService *services.SystemConfigService
	securityService     *services.SecurityService
}

func NewAdminCardHandler(
	cardService *services.CardService,
	adminService *services.AdminService,
	productService *services.ProductService,
	systemConfigService *services.SystemConfigService,
) *AdminCardHandler {
	return &AdminCardHandler{
		cardService:         cardService,
		adminService:        adminService,
		productService:      productService,
		systemConfigService: systemConfigService,
		securityService:     services.NewSecurityService(),
	}
}

// ListCards 获取卡片列表
func (h *AdminCardHandler) ListCards(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	
	filters := make(map[string]interface{})
	
	// 状态筛选
	if status := c.Query("status"); status != "" {
		statusInt, _ := strconv.Atoi(status)
		filters["status"] = statusInt
	}
	
	// 价格ID筛选
	if priceID := c.Query("price_id"); priceID != "" {
		priceIDInt, _ := strconv.ParseInt(priceID, 10, 64)
		filters["price_id"] = priceIDInt
	}
	
	// 搜索
	if search := c.Query("search"); search != "" {
		filters["search"] = search
	}
	
	offset := (page - 1) * pageSize
	cards, total, err := h.cardService.List(offset, pageSize, filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "获取卡片列表失败",
			"data": nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "成功",
		"data": gin.H{
			"list":  cards,
			"total": total,
		},
	})
}

// GetCard 获取卡片详情
func (h *AdminCardHandler) GetCard(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的卡片ID",
			"data": nil,
		})
		return
	}
	
	card, err := h.cardService.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "卡片不存在",
			"data": nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "成功",
		"data": card,
	})
}

type CreateCardRequest struct {
	CardCode    string     `json:"card_code" binding:"required"`
	PriceID     int64      `json:"price_id" binding:"required"`
	CostPrice   float64    `json:"cost_price" binding:"required,min=0"`
	SellPrice   float64    `json:"sell_price" binding:"required,min=0"`
	ExpiredAt   *time.Time `json:"expired_at"`
	Description string     `json:"description"`
}

// CreateCard 创建卡片
func (h *AdminCardHandler) CreateCard(c *gin.Context) {
	var req CreateCardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "请求参数无效",
			"data": nil,
		})
		return
	}

	card := &models.Card{
		CardCode:    req.CardCode,
		Status:      0, // 默认未使用
		PriceID:     req.PriceID,
		CostPrice:   req.CostPrice,
		SellPrice:   req.SellPrice,
		Description: req.Description,
	}
	
	// 设置过期时间，如果没有提供则默认1年后
	if req.ExpiredAt != nil {
		card.ExpiredAt = req.ExpiredAt
	} else {
		expiredAt := time.Now().AddDate(1, 0, 0)
		card.ExpiredAt = &expiredAt
	}

	if err := h.cardService.Create(card); err != nil {
		// 判断是否是卡片代码重复错误
		if err.Error() == "card code already exists" {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 400,
				"msg":  "卡片代码已存在",
				"data": nil,
			})
			return
		}
		
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "创建卡片失败: " + err.Error(),
			"data": nil,
		})
		return
	}

	// 记录操作日志
	adminID, _ := c.Get("admin_id")
	h.adminService.LogOperation(&models.AdminOperationLog{
		AdminID:   adminID.(uint),
		Operation: "create_card",
		Module:    "card",
		Details:   "Created card: " + card.CardCode,
		IPAddress: c.ClientIP(),
		UserAgent: c.GetHeader("User-Agent"),
	})

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "卡片创建成功",
		"data": card,
	})
}

// CreateCardWithoutValidation 创建卡片（不验证）- 仅用于测试
func (h *AdminCardHandler) CreateCardWithoutValidation(c *gin.Context) {
	var req CreateCardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "请求参数无效",
			"data": nil,
		})
		return
	}

	card := &models.Card{
		CardCode:    req.CardCode,
		Status:      0, // 默认未使用
		PriceID:     req.PriceID,
		CostPrice:   req.CostPrice,
		SellPrice:   req.SellPrice,
		Description: req.Description,
	}
	
	// 设置过期时间
	if req.ExpiredAt != nil {
		card.ExpiredAt = req.ExpiredAt
	} else {
		expiredAt := time.Now().AddDate(1, 0, 0)
		card.ExpiredAt = &expiredAt
	}

	// 直接创建，跳过验证
	if err := h.cardService.CreateWithoutValidation(card); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "创建卡片失败: " + err.Error(),
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "卡片创建成功（未验证）",
		"data": card,
	})
}

type UpdateCardRequest struct {
	Status      *int       `json:"status"`
	CostPrice   *float64   `json:"cost_price"`
	SellPrice   *float64   `json:"sell_price"`
	ExpiredAt   *time.Time `json:"expired_at"`
	Description *string    `json:"description"`
}

// UpdateCard 更新卡片
func (h *AdminCardHandler) UpdateCard(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的卡片ID",
			"data": nil,
		})
		return
	}
	
	var req UpdateCardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "请求参数无效",
			"data": nil,
		})
		return
	}
	
	card, err := h.cardService.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "卡片不存在",
			"data": nil,
		})
		return
	}
	
	// 更新字段
	if req.Status != nil {
		card.Status = *req.Status
	}
	if req.CostPrice != nil {
		card.CostPrice = *req.CostPrice
	}
	if req.SellPrice != nil {
		card.SellPrice = *req.SellPrice
	}
	if req.ExpiredAt != nil {
		card.ExpiredAt = req.ExpiredAt
	}
	if req.Description != nil {
		card.Description = *req.Description
	}
	
	if err := h.cardService.Update(card); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "更新卡片失败",
			"data": nil,
		})
		return
	}
	
	// 记录操作日志
	adminID, _ := c.Get("admin_id")
	h.adminService.LogOperation(&models.AdminOperationLog{
		AdminID:   adminID.(uint),
		Operation: "update_card",
		Module:    "card",
		Details:   "Updated card: " + card.CardCode,
		IPAddress: c.ClientIP(),
		UserAgent: c.GetHeader("User-Agent"),
	})
	
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "卡片更新成功",
		"data": card,
	})
}

// DeleteCard 删除卡片
func (h *AdminCardHandler) DeleteCard(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的卡片ID",
			"data": nil,
		})
		return
	}
	
	if err := h.cardService.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "删除卡片失败: " + err.Error(),
			"data": nil,
		})
		return
	}
	
	// 记录操作日志
	adminID, _ := c.Get("admin_id")
	h.adminService.LogOperation(&models.AdminOperationLog{
		AdminID:   adminID.(uint),
		Operation: "delete_card",
		Module:    "card",
		Details:   "Deleted card ID: " + c.Param("id"),
		IPAddress: c.ClientIP(),
		UserAgent: c.GetHeader("User-Agent"),
	})
	
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "卡片删除成功",
		"data": nil,
	})
}

// GetCardUsageLogs 获取卡片使用记录
func (h *AdminCardHandler) GetCardUsageLogs(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的卡片ID",
			"data": nil,
		})
		return
	}
	
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	
	offset := (page - 1) * pageSize
	logs, total, err := h.cardService.GetUsageLogs(uint(id), offset, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "获取使用记录失败",
			"data": nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "成功",
		"data": gin.H{
			"list":  logs,
			"total": total,
		},
	})
}

// BatchImportCards 批量导入卡片
func (h *AdminCardHandler) BatchImportCards(c *gin.Context) {
	var req services.BatchImportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "请求参数无效: " + err.Error(),
			"data": nil,
		})
		return
	}
	
	// 设置管理员ID
	adminID, _ := c.Get("admin_id")
	req.AdminID = adminID.(uint)
	
	// 批量导入
	batch, err := h.cardService.BatchImport(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "批量导入失败: " + err.Error(),
			"data": nil,
		})
		return
	}
	
	// 记录操作日志
	h.adminService.LogOperation(&models.AdminOperationLog{
		AdminID:   req.AdminID,
		Operation: "batch_import_cards",
		Module:    "card",
		Details:   "Batch imported " + strconv.Itoa(batch.TotalCount) + " cards",
		IPAddress: c.ClientIP(),
		UserAgent: c.GetHeader("User-Agent"),
	})
	
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "批量导入成功",
		"data": batch,
	})
}

// GetCardStats 获取卡片统计
func (h *AdminCardHandler) GetCardStats(c *gin.Context) {
	// 检查是否指定了价格ID
	priceIDStr := c.Query("price_id")
	if priceIDStr != "" {
		priceID, err := strconv.ParseInt(priceIDStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 400,
				"msg":  "无效的价格ID",
				"data": nil,
			})
			return
		}
		
		// 获取指定价格的卡片统计
		stats, err := h.cardService.GetPriceStats(priceID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": 500,
				"msg":  "获取价格统计数据失败",
				"data": nil,
			})
			return
		}
		
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "成功",
			"data": stats,
		})
		return
	}
	
	// 获取全部卡片统计
	stats, err := h.cardService.GetCardStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "获取统计数据失败",
			"data": nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "成功",
		"data": stats,
	})
}

// ListBatches 获取批次列表
func (h *AdminCardHandler) ListBatches(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	
	offset := (page - 1) * pageSize
	batches, total, err := h.cardService.ListBatches(offset, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "获取批次列表失败",
			"data": nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "成功",
		"data": gin.H{
			"list":  batches,
			"total": total,
		},
	})
}

// GetBatch 获取批次详情
func (h *AdminCardHandler) GetBatch(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的批次ID",
			"data": nil,
		})
		return
	}
	
	batch, err := h.cardService.GetBatch(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "批次不存在",
			"data": nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "成功",
		"data": batch,
	})
}

// GetBatchCards 获取批次下的卡片
func (h *AdminCardHandler) GetBatchCards(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的批次ID",
			"data": nil,
		})
		return
	}
	
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	
	offset := (page - 1) * pageSize
	cards, total, err := h.cardService.GetBatchCards(uint(id), offset, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "获取批次卡片失败",
			"data": nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "成功",
		"data": gin.H{
			"list":  cards,
			"total": total,
		},
	})
}

// SyncCardProducts 同步卡片商品
func (h *AdminCardHandler) SyncCardProducts(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的卡片ID",
			"data": nil,
		})
		return
	}
	
	// 获取卡片信息
	card, err := h.cardService.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "卡片不存在",
			"data": nil,
		})
		return
	}
	
	// 获取同步店铺代码
	storeCode, err := h.systemConfigService.GetSyncStoreCode()
	if err != nil || storeCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "请先配置同步店铺代码",
			"data": nil,
		})
		return
	}
	
	// 同步商品
	result, err := h.productService.SyncProductsFromCard(card.CardCode, storeCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "同步商品失败: " + err.Error(),
			"data": nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "同步商品成功",
		"data": result,
	})
}

// ValidateCard 验证单个卡片
func (h *AdminCardHandler) ValidateCard(c *gin.Context) {
	var req struct {
		CardCode string `json:"card_code" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "请求参数无效",
			"data": nil,
		})
		return
	}
	
	// 使用新的验证并更新方法
	isValid, message, updated, err := h.cardService.ValidateAndUpdateCard(req.CardCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "验证失败: " + err.Error(),
			"data": nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "验证完成",
		"data": gin.H{
			"is_valid": isValid,
			"message":  message,
			"updated":  updated,
		},
	})
}

// BatchValidateCards 批量验证卡片
func (h *AdminCardHandler) BatchValidateCards(c *gin.Context) {
	result, err := h.cardService.BatchValidateAndUpdateCards()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "批量验证失败: " + err.Error(),
			"data": nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "批量验证完成",
		"data": result,
	})
}

// StartBatchValidation 启动批量验证（支持双模式）
func (h *AdminCardHandler) StartBatchValidation(c *gin.Context) {
	var req struct {
		Mode string `json:"mode" binding:"required"` // "all" 或 "smart"
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误: " + err.Error(),
			"data": nil,
		})
		return
	}
	
	// 验证模式
	var mode services.ValidationMode
	switch req.Mode {
	case "all":
		mode = services.ValidationModeAll
	case "smart":
		mode = services.ValidationModeSmart
	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的验证模式，支持: all, smart",
			"data": nil,
		})
		return
	}
	
	// 启动验证任务
	task, err := h.cardService.StartBatchValidation(mode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "启动验证任务失败: " + err.Error(),
			"data": nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "验证任务已启动",
		"data": task,
	})
}

// GetValidationProgress 获取验证进度
func (h *AdminCardHandler) GetValidationProgress(c *gin.Context) {
	taskID := c.Param("taskId")
	if taskID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "任务ID不能为空",
			"data": nil,
		})
		return
	}
	
	task, err := h.cardService.GetValidationTask(taskID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "任务不存在",
			"data": nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": task,
	})
}

// CancelValidation 取消验证任务
func (h *AdminCardHandler) CancelValidation(c *gin.Context) {
	taskID := c.Param("taskId")
	if taskID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "任务ID不能为空",
			"data": nil,
		})
		return
	}
	
	err := h.cardService.CancelValidationTask(taskID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "取消任务失败: " + err.Error(),
			"data": nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "任务已取消",
		"data": nil,
	})
}

// GetValidationStatistics 获取验证统计信息
func (h *AdminCardHandler) GetValidationStatistics(c *gin.Context) {
	stats := h.cardService.GetValidationStatistics()
	
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": stats,
	})
}
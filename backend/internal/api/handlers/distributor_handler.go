package handlers

import (
	"net/http"
	"strconv"

	"backend/internal/services"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type DistributorHandler struct {
	orderService         *services.OrderService
	storeService         *services.StoreService
	cityService          *services.CityService
	cardService          *services.CardService
	securityAuditService *services.SecurityAuditService
	logger               *zap.Logger
}

func NewDistributorHandler(
	orderService *services.OrderService,
	storeService *services.StoreService,
	cityService *services.CityService,
	cardService *services.CardService,
	securityAuditService *services.SecurityAuditService,
	logger *zap.Logger,
) *DistributorHandler {
	return &DistributorHandler{
		orderService:         orderService,
		storeService:         storeService,
		cityService:          cityService,
		cardService:          cardService,
		securityAuditService: securityAuditService,
		logger:               logger,
	}
}

// CreateOrder 创建订单
func (h *DistributorHandler) CreateOrder(c *gin.Context) {
	var req services.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "请求参数无效",
			"data": nil,
		})
		return
	}

	// 从上下文获取分销商ID
	distributorID, _ := c.Get("distributor_id")
	req.DistributorID = distributorID.(uint)

	order, err := h.orderService.CreateOrder(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "订单创建成功",
		"data": gin.H{
			"order_no":      order.OrderNo,
			"status":        order.Status,
			"total_amount":  order.TotalAmount,
			"store_name":    order.StoreName,
			"store_address": order.StoreAddress,
		},
	})
}

// BatchCreateOrders 批量创建订单
func (h *DistributorHandler) BatchCreateOrders(c *gin.Context) {
	var req struct {
		Orders []services.CreateOrderRequest `json:"orders" binding:"required,min=1,max=10"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "请求参数无效",
			"data": nil,
		})
		return
	}

	// 从上下文获取分销商ID
	distributorID, _ := c.Get("distributor_id")

	// 批量处理订单
	results := make([]map[string]interface{}, len(req.Orders))
	successCount := 0

	for i, orderReq := range req.Orders {
		orderReq.DistributorID = distributorID.(uint)

		order, err := h.orderService.CreateOrder(&orderReq)
		if err != nil {
			results[i] = map[string]interface{}{
				"success": false,
				"error":   err.Error(),
				"index":   i,
			}
		} else {
			successCount++
			orderData := BuildOrderResponse(order, h.orderService)
			orderData["success"] = true
			orderData["index"] = i
			results[i] = orderData
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "Batch order creation completed",
		"data": gin.H{
			"total":   len(req.Orders),
			"success": successCount,
			"failed":  len(req.Orders) - successCount,
			"results": results,
		},
	})
}

// QueryOrder 查询订单状态
func (h *DistributorHandler) QueryOrder(c *gin.Context) {
	orderNo := c.Param("orderNo")
	if orderNo == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "Order number is required",
			"data": nil,
		})
		return
	}

	order, err := h.orderService.QueryOrderStatus(orderNo)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "订单不存在",
			"data": nil,
		})
		return
	}

	// 验证订单属于当前分销商
	distributorID, _ := c.Get("distributor_id")
	if order.DistributorID != distributorID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{
			"code": 403,
			"msg":  "Access denied",
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "成功",
		"data": BuildOrderResponse(order, h.orderService),
	})
}

// SearchStores 搜索门店
func (h *DistributorHandler) SearchStores(c *gin.Context) {
	card := c.Query("card")
	cityIDStr := c.Query("city_id")
	cityName := c.Query("city_name") // 新增：支持城市名称
	keywords := c.Query("keywords")

	if card == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "Card is required",
			"data": nil,
		})
		return
	}

	// 获取分销商ID进行权限检查（软模式：仅记录，不阻止）
	distributorID := c.GetUint("distributor_id")
	hasAccess, err := h.cardService.CheckDistributorCardAccess(distributorID, card)
	if err != nil {
		h.logger.Error("检查卡片权限失败",
			zap.Uint("distributor_id", distributorID),
			zap.String("card", card),
			zap.Error(err))
	} else {
		// 记录审计日志
		h.securityAuditService.LogCardAccess(distributorID, card, "SearchStores", hasAccess, c)
		
		if !hasAccess {
			// 软模式：仅记录潜在的越权访问，不阻止
			h.logger.Warn("分销商访问未授权卡片",
				zap.Uint("distributor_id", distributorID),
				zap.String("card", card),
				zap.String("api", "SearchStores"))
		}
	}

	var cityID int

	// 如果提供了城市名称，转换为城市ID
	if cityName != "" {
		cityID, err = h.storeService.GetCityIDByName(cityName)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 400,
				"msg":  "Invalid city name",
				"data": nil,
			})
			return
		}
	} else {
		cityID, _ = strconv.Atoi(cityIDStr)
	}

	stores, err := h.storeService.SearchStores(card, cityID, keywords)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "成功",
		"data": stores,
	})
}

// GetMenu 获取门店菜单
func (h *DistributorHandler) GetMenu(c *gin.Context) {
	card := c.Query("card")
	storeCode := c.Query("store_code")

	if card == "" || storeCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "Card and store_code are required",
			"data": nil,
		})
		return
	}

	// 权限检查（软模式）
	distributorID := c.GetUint("distributor_id")
	hasAccess, err := h.cardService.CheckDistributorCardAccess(distributorID, card)
	if err != nil {
		h.logger.Error("检查卡片权限失败",
			zap.Uint("distributor_id", distributorID),
			zap.String("card", card),
			zap.Error(err))
	} else {
		// 记录审计日志
		h.securityAuditService.LogCardAccess(distributorID, card, "GetMenu", hasAccess, c)
		
		if !hasAccess {
			h.logger.Warn("分销商访问未授权卡片",
				zap.Uint("distributor_id", distributorID),
				zap.String("card", card),
				zap.String("api", "GetMenu"))
		}
	}

	menu, err := h.storeService.GetMenu(card, storeCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "成功",
		"data": menu,
	})
}

// GetGoodsDetail 获取商品详情
func (h *DistributorHandler) GetGoodsDetail(c *gin.Context) {
	card := c.Query("card")
	storeCode := c.Query("store_code")
	goodsCode := c.Query("goods_code")

	if card == "" || storeCode == "" || goodsCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "Card, store_code and goods_code are required",
			"data": nil,
		})
		return
	}

	// 权限检查（软模式）
	distributorID := c.GetUint("distributor_id")
	hasAccess, err := h.cardService.CheckDistributorCardAccess(distributorID, card)
	if err != nil {
		h.logger.Error("检查卡片权限失败",
			zap.Uint("distributor_id", distributorID),
			zap.String("card", card),
			zap.Error(err))
	} else {
		// 记录审计日志
		h.securityAuditService.LogCardAccess(distributorID, card, "GetGoodsDetail", hasAccess, c)
		
		if !hasAccess {
			h.logger.Warn("分销商访问未授权卡片",
				zap.Uint("distributor_id", distributorID),
				zap.String("card", card),
				zap.String("api", "GetGoodsDetail"))
		}
	}

	goods, err := h.storeService.GetGoodsDetail(card, storeCode, goodsCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "成功",
		"data": goods,
	})
}

// GetCities 获取城市列表
func (h *DistributorHandler) GetCities(c *gin.Context) {
	cities, err := h.cityService.GetAllCities()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "Failed to get cities",
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "成功",
		"data": cities,
	})
}

// SyncCities 同步城市数据
func (h *DistributorHandler) SyncCities(c *gin.Context) {
	card := c.Query("card")
	if card == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "Card is required",
			"data": nil,
		})
		return
	}

	if err := h.cityService.SyncCities(card); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "同步城市数据失败",
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "城市数据同步成功",
		"data": nil,
	})
}

// ParseSpecsCode 解析规格代码（分销商接口）
func (h *DistributorHandler) ParseSpecsCode(c *gin.Context) {
	var req struct {
		GoodsCode string `json:"goods_code" binding:"required"`
		SKUCode   string `json:"sku_code" binding:"required"`
		SpecsCode string `json:"specs_code" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "请求参数错误: " + err.Error(),
			"data": nil,
		})
		return
	}

	// 暂时返回501，建议使用管理员接口
	c.JSON(http.StatusNotImplemented, gin.H{
		"code": 501,
		"msg":  "功能暂未实现，请使用管理员接口 POST /api/v1/admin/products/parse-specs",
		"data": nil,
	})
}

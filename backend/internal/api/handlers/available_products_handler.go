package handlers

import (
	"backend/internal/models"
	"backend/internal/services"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// AvailableProductsHandler 可用商品处理器
type AvailableProductsHandler struct {
	productService *services.ProductService
	cacheService   *services.CacheService
}

// NewAvailableProductsHandler 创建可用商品处理器
func NewAvailableProductsHandler(productService *services.ProductService, cacheService *services.CacheService) *AvailableProductsHandler {
	return &AvailableProductsHandler{
		productService: productService,
		cacheService:   cacheService,
	}
}

// GetDistributorAvailableProducts 获取分销商可用商品列表
// @Summary 获取可用商品列表
// @Description 获取当前系统中可下单的商品信息，包含中文名称和编码，不暴露卡片信息
// @Tags Distributor
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} models.AvailableProductsResponse
// @Failure 401 {object} gin.H
// @Failure 429 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/distributor/available-products [get]
func (h *AvailableProductsHandler) GetDistributorAvailableProducts(c *gin.Context) {
	// 获取分销商信息
	distributor, exists := c.Get("distributor")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": 401,
			"msg":  "未授权的访问",
			"data": nil,
		})
		return
	}

	distributorModel := distributor.(*models.Distributor)

	// 检查缓存
	cacheKey := "available_products:" + string(distributorModel.ID)
	if h.cacheService != nil {
		cachedData, err := h.cacheService.Get(cacheKey)
		if err == nil && cachedData != nil {
			// 返回缓存数据
			c.JSON(http.StatusOK, gin.H{
				"code": 200,
				"msg":  "成功",
				"data": cachedData,
			})
			return
		}
	}

	// 获取可用商品数据
	response, err := h.productService.GetAvailableProducts(distributorModel.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "获取商品信息失败",
			"data": nil,
		})
		return
	}

	// 缓存数据（5分钟）
	if h.cacheService != nil {
		h.cacheService.Set(cacheKey, response, 5*time.Minute)
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "成功",
		"data": response,
	})
}

// GetAdminAvailableProducts 获取管理员可用商品列表
// @Summary 获取所有可用商品列表（管理员）
// @Description 获取系统中所有可下单的商品信息，供管理员查看
// @Tags Admin
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} models.AvailableProductsResponse
// @Failure 401 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/admin/available-products [get]
func (h *AvailableProductsHandler) GetAdminAvailableProducts(c *gin.Context) {
	// 管理员可以看到所有商品
	response, err := h.productService.GetAvailableProducts(0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "获取商品信息失败",
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "成功",
		"data": response,
	})
}
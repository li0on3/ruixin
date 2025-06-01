package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"backend/internal/models"
	"backend/internal/services"
	"github.com/gin-gonic/gin"
)

type AdminProductHandler struct {
	productService       *services.ProductService
	systemConfigService  *services.SystemConfigService
}

func NewAdminProductHandler(productService *services.ProductService, systemConfigService *services.SystemConfigService) *AdminProductHandler {
	return &AdminProductHandler{
		productService:      productService,
		systemConfigService: systemConfigService,
	}
}

// GetProductList 获取商品列表
func (h *AdminProductHandler) GetProductList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	search := c.Query("search")
	cardID, _ := strconv.Atoi(c.Query("card_id"))

	// 如果指定了卡片ID，获取该卡片绑定的产品
	if cardID > 0 {
		products, total, err := h.productService.GetCardBoundProducts(uint(cardID), page, pageSize, search)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "获取成功",
			"data": gin.H{
				"list":  products,
				"total": total,
			},
		})
		return
	}

	// 否则获取所有产品
	products, total, err := h.productService.GetProductList(page, pageSize, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": gin.H{
			"list":  products,
			"total": total,
		},
	})
}

// SearchProducts 搜索商品（用于价格配置时选择商品）
func (h *AdminProductHandler) SearchProducts(c *gin.Context) {
	keyword := c.Query("keyword")
	if keyword == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "搜索关键词不能为空"})
		return
	}

	products, err := h.productService.SearchProducts(keyword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "搜索成功",
		"data": products,
	})
}

// GetProductsByCodes 根据商品代码批量获取商品信息
func (h *AdminProductHandler) GetProductsByCodes(c *gin.Context) {
	var req struct {
		Codes string `json:"codes" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 分割商品代码
	codes := strings.Split(req.Codes, ",")
	products, err := h.productService.GetProductsByCodes(codes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": products,
	})
}

// SyncProducts 同步商品信息
func (h *AdminProductHandler) SyncProducts(c *gin.Context) {
	var req struct {
		StoreCode string `json:"store_code"` // 可选，不传则使用系统配置
		CardCode  string `json:"card_code"`  // 必需，用于获取商品信息
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "请求参数错误",
			"data": nil,
		})
		return
	}

	// 如果没有提供卡片代码，返回错误
	if req.CardCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "卡片代码不能为空",
			"data": nil,
		})
		return
	}

	// 如果没有提供店铺代码，从系统配置获取
	storeCode := req.StoreCode
	if storeCode == "" {
		var err error
		storeCode, err = h.systemConfigService.GetSyncStoreCode()
		if err != nil || storeCode == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 400,
				"msg":  "请配置同步店铺代码或在请求中提供",
				"data": nil,
			})
			return
		}
	}

	// 执行同步
	result, err := h.productService.SyncProductsFromCard(req.CardCode, storeCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "同步失败: " + err.Error(),
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "同步成功",
		"data": gin.H{
			"synced_count":  result.SyncedCount,
			"new_count":     result.NewCount,
			"updated_count": result.UpdatedCount,
			"failed_count":  result.FailedCount,
			"alias_count":   result.AliasCount,
			"duration":      result.EndTime.Sub(result.StartTime).Seconds(),
		},
	})
}

// GetMatchLogs 获取匹配失败日志
func (h *AdminProductHandler) GetMatchLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	distributorID, _ := strconv.ParseUint(c.Query("distributor_id"), 10, 32)
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	logs, total, err := h.productService.GetMatchLogs(page, pageSize, uint(distributorID), startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "获取日志失败: " + err.Error(),
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "查询成功",
		"data": gin.H{
			"list":      logs,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// ParseSpecsCode 解析规格代码
func (h *AdminProductHandler) ParseSpecsCode(c *gin.Context) {
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

	specItems, err := h.productService.ParseSpecsCode(req.GoodsCode, req.SKUCode, req.SpecsCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "解析规格代码失败: " + err.Error(),
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "解析成功",
		"data": specItems,
	})
}

// GetProductCards 获取产品的卡片绑定详情
func (h *AdminProductHandler) GetProductCards(c *gin.Context) {
	goodsCode := c.Query("goods_code")
	if goodsCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "商品代码不能为空",
		})
		return
	}

	// 查询产品
	var product models.Product
	if err := h.productService.GetDB().Where("goods_code = ?", goodsCode).First(&product).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "商品不存在",
		})
		return
	}

	// 查询该产品的所有卡片绑定
	var bindings []models.CardProductBinding
	err := h.productService.GetDB().
		Where("product_id = ? AND is_active = ?", product.ID, true).
		Preload("Card").
		Preload("Card.Price").
		Find(&bindings).Error
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "查询失败: " + err.Error(),
		})
		return
	}

	// 构建响应数据
	var cards []map[string]interface{}
	for _, binding := range bindings {
		if binding.Card != nil {
			cardData := map[string]interface{}{
				"card_id":    binding.Card.ID,
				"card_code":  binding.Card.CardCode,
				"status":     binding.Card.Status,
				"price":      0.0,
				"price_name": "",
			}
			
			// 添加价格信息
			if binding.Card.Price != nil {
				cardData["price"] = binding.Card.Price.PriceValue
				cardData["price_name"] = binding.Card.Price.PriceCode // 使用价格代码作为名称
			}
			
			cards = append(cards, cardData)
		}
	}

	// 构建完整的产品信息
	productData := map[string]interface{}{
		"id":         product.ID,
		"goods_code": product.GoodsCode,
		"goods_name": product.GoodsName,
		"cards":      cards,
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "查询成功",
		"data": gin.H{
			"list":  []interface{}{productData},
			"total": 1,
		},
	})
}
package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"backend/internal/models"
	"backend/internal/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AdminLuckinConfigHandler struct {
	configService *services.LuckinConfigService
}

func NewAdminLuckinConfigHandler(db *gorm.DB) *AdminLuckinConfigHandler {
	return &AdminLuckinConfigHandler{
		configService: services.NewLuckinConfigService(db),
	}
}

// Price Management Handlers

// GetPriceList 获取价格列表
func (h *AdminLuckinConfigHandler) GetPriceList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	prices, total, err := h.configService.GetPriceList(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "成功",
		"data": gin.H{
			"list": prices,
			"pagination": gin.H{
				"total":      total,
				"page":       page,
				"page_size":  pageSize,
				"total_page": (total + int64(pageSize) - 1) / int64(pageSize),
			},
		},
	})
}

// CreatePrice 创建价格
func (h *AdminLuckinConfigHandler) CreatePrice(c *gin.Context) {
	var req struct {
		PriceID      string  `json:"price_id" binding:"required"`
		PriceValue   float64 `json:"price_value" binding:"required,gt=0"`
		Status       int     `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	adminID := c.GetUint("admin_id")
	price := &models.LuckinPrice{
		PriceCode:    req.PriceID,
		PriceValue:   req.PriceValue,
		Status:       req.Status,
		CreatedBy:    int64(adminID),
	}

	if err := h.configService.CreatePrice(price); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "价格创建成功", "data": price})
}

// UpdatePrice 更新价格
func (h *AdminLuckinConfigHandler) UpdatePrice(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的价格ID"})
		return
	}

	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.configService.UpdatePrice(id, req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "价格更新成功", "data": nil})
}

// DeletePrice 删除价格
func (h *AdminLuckinConfigHandler) DeletePrice(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的价格ID"})
		return
	}

	if err := h.configService.DeletePrice(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "价格删除成功", "data": nil})
}

// Product Management Handlers

// GetProductList 获取产品列表
func (h *AdminLuckinConfigHandler) GetProductList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	filters := make(map[string]interface{})
	if name := c.Query("name"); name != "" {
		filters["name"] = name
	}
	if category := c.Query("category"); category != "" {
		filters["category"] = category
	}
	if status := c.Query("status"); status != "" {
		if s, err := strconv.Atoi(status); err == nil {
			filters["status"] = s
		}
	}

	products, total, err := h.configService.GetProductList(filters, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "成功",
		"data": gin.H{
			"data": products,
			"pagination": gin.H{
				"total":      total,
				"page":       page,
				"page_size":  pageSize,
				"total_page": (total + int64(pageSize) - 1) / int64(pageSize),
			},
		},
	})
}

// CreateProduct 创建产品
func (h *AdminLuckinConfigHandler) CreateProduct(c *gin.Context) {
	var req struct {
		ProductID   string `json:"product_id" binding:"required"`
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		Category    string `json:"category"`
		ImageURL    string `json:"image_url"`
		Status      int    `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	adminID := int64(c.GetUint("admin_id"))
	product := &models.LuckinProduct{
		ProductID:   req.ProductID,
		Name:        req.Name,
		Description: req.Description,
		Category:    req.Category,
		ImageURL:    req.ImageURL,
		Status:      req.Status,
		CreatedBy:   adminID,
	}

	if err := h.configService.CreateProduct(product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "产品创建成功", "data": product})
}

// UpdateProduct 更新产品
func (h *AdminLuckinConfigHandler) UpdateProduct(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的产品ID"})
		return
	}

	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.configService.UpdateProduct(id, req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "产品更新成功", "data": nil})
}

// DeleteProduct 删除产品
func (h *AdminLuckinConfigHandler) DeleteProduct(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的产品ID"})
		return
	}

	if err := h.configService.DeleteProduct(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "产品删除成功", "data": nil})
}

// GetProductCategories 获取产品类别列表
func (h *AdminLuckinConfigHandler) GetProductCategories(c *gin.Context) {
	categories, err := h.configService.GetProductCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "Success", "data": categories})
}

// Category Binding Handlers

// GetCategoryBindings 获取种类绑定列表
func (h *AdminLuckinConfigHandler) GetCategoryBindings(c *gin.Context) {
	categoryID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的种类ID"})
		return
	}

	bindings, err := h.configService.GetCategoryBindings(categoryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "Success", "data": bindings})
}

// CreateBinding 创建绑定
func (h *AdminLuckinConfigHandler) CreateBinding(c *gin.Context) {
	categoryID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的种类ID"})
		return
	}

	var req struct {
		TargetType string `json:"target_type" binding:"required,oneof=price product"`
		TargetID   string `json:"target_id" binding:"required"`
		Priority   int    `json:"priority"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	adminID := int64(c.GetUint("admin_id"))

	if err := h.configService.CreateBinding(categoryID, req.TargetType, req.TargetID, req.Priority, adminID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "绑定创建成功", "data": nil})
}

// DeleteBinding 删除绑定
func (h *AdminLuckinConfigHandler) DeleteBinding(c *gin.Context) {
	bindingID, err := strconv.ParseInt(c.Param("bindingId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的绑定ID"})
		return
	}

	if err := h.configService.DeleteBinding(bindingID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "绑定删除成功", "data": nil})
}

// UpdateBindingPriority 更新绑定优先级
func (h *AdminLuckinConfigHandler) UpdateBindingPriority(c *gin.Context) {
	bindingID, err := strconv.ParseInt(c.Param("bindingId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的绑定ID"})
		return
	}

	var req struct {
		Priority int `json:"priority" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.configService.UpdateBindingPriority(bindingID, req.Priority); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "优先级更新成功", "data": nil})
}

// BatchImportProducts 批量导入产品
func (h *AdminLuckinConfigHandler) BatchImportProducts(c *gin.Context) {
	var req struct {
		Products []struct {
			ProductID   string `json:"product_id" binding:"required"`
			Name        string `json:"name" binding:"required"`
			Description string `json:"description"`
			Category    string `json:"category"`
			ImageURL    string `json:"image_url"`
		} `json:"products" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	adminID := int64(c.GetUint("admin_id"))
	products := make([]*models.LuckinProduct, len(req.Products))
	for i, p := range req.Products {
		products[i] = &models.LuckinProduct{
			ProductID:   p.ProductID,
			Name:        p.Name,
			Description: p.Description,
			Category:    p.Category,
			ImageURL:    p.ImageURL,
			Status:      1,
			CreatedBy:   adminID,
		}
	}

	if err := h.configService.BatchImportProducts(products); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": fmt.Sprintf("成功导入%d个产品", len(products)), "data": nil})
}

// GetActiveOptions 获取有效的价格和产品选项（用于下拉框）
func (h *AdminLuckinConfigHandler) GetActiveOptions(c *gin.Context) {
	prices, _ := h.configService.GetActivePrices()
	products, _ := h.configService.GetActiveProducts()

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "成功",
		"data": gin.H{
			"prices":   prices,
			"products": products,
		},
	})
}


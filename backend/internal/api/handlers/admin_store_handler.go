package handlers

import (
	"backend/internal/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type AdminStoreHandler struct {
	storeService *services.StoreService
	cityService  *services.CityService
	cardService  *services.CardService
}

func NewAdminStoreHandler(storeService *services.StoreService, cityService *services.CityService, cardService *services.CardService) *AdminStoreHandler {
	return &AdminStoreHandler{
		storeService: storeService,
		cityService:  cityService,
		cardService:  cardService,
	}
}

// SearchStores 搜索店铺
func (h *AdminStoreHandler) SearchStores(c *gin.Context) {
	// 获取查询参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	cityID := c.Query("city_id")
	keyword := c.Query("keyword")
	longitude := c.Query("longitude")
	latitude := c.Query("latitude")
	
	// 参数验证
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	
	// 构建搜索参数
	params := &services.StoreSearchParams{
		CityID:    cityID,
		Keyword:   keyword,
		Longitude: longitude,
		Latitude:  latitude,
		Page:      page,
		PageSize:  pageSize,
	}
	
	// 执行搜索
	stores, total, err := h.storeService.SearchStoresAdmin(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "搜索店铺失败: " + err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "成功",
		"data": gin.H{
			"list":  stores,
			"total": total,
		},
	})
}

// GetCities 获取城市列表
func (h *AdminStoreHandler) GetCities(c *gin.Context) {
	cities, err := h.cityService.GetAllCities()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "获取城市列表失败",
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "成功",
		"data": cities,
	})
}

// GetHotStores 获取热门店铺
func (h *AdminStoreHandler) GetHotStores(c *gin.Context) {
	// 返回一些常用的店铺，这里可以根据实际需求调整
	hotStores := []map[string]interface{}{
		{
			"store_code": "390840",
			"store_name": "北京王府井店",
			"city_name":  "北京",
		},
		{
			"store_code": "123456", 
			"store_name": "上海南京路店",
			"city_name":  "上海",
		},
		{
			"store_code": "234567",
			"store_name": "广州天河店", 
			"city_name":  "广州",
		},
		{
			"store_code": "345678",
			"store_name": "深圳福田店",
			"city_name":  "深圳",
		},
		{
			"store_code": "456789",
			"store_name": "成都春熙路店",
			"city_name":  "成都",
		},
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "成功",
		"data": hotStores,
	})
}

// SyncCities 同步城市列表
func (h *AdminStoreHandler) SyncCities(c *gin.Context) {
	// 先尝试从内部获取可用的卡片
	var cardCode string
	
	// 检查是否有可用的卡片
	cards, err := h.cardService.GetAvailableCards(1) // 只需要一张卡片
	if err == nil && len(cards) > 0 {
		cardCode = cards[0].CardCode
		// 自动同步
		err = h.cityService.SyncCities(cardCode)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": 500,
				"msg":  "同步城市失败: " + err.Error(),
			})
			return
		}
	} else {
		// 如果没有可用卡片，要求用户提供
		type SyncRequest struct {
			CardCode string `json:"card_code"`
		}
		
		var req SyncRequest
		if err := c.ShouldBindJSON(&req); err != nil || req.CardCode == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 400,
				"msg":  "系统中没有可用的卡片，请提供一个有效的卡片代码",
			})
			return
		}
		
		// 使用用户提供的卡片同步
		err = h.cityService.SyncCities(req.CardCode)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": 500,
				"msg":  "同步城市失败: " + err.Error(),
			})
			return
		}
	}

	// 返回最新的城市列表
	cities, err := h.cityService.GetAllCities()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "获取城市列表失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "同步成功",
		"data": cities,
	})
}
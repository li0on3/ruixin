package handlers

import (
	"backend/internal/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AdminSystemConfigHandler struct {
	systemConfigService *services.SystemConfigService
}

func NewAdminSystemConfigHandler(systemConfigService *services.SystemConfigService) *AdminSystemConfigHandler {
	return &AdminSystemConfigHandler{
		systemConfigService: systemConfigService,
	}
}

// GetConfigs 获取系统配置
func (h *AdminSystemConfigHandler) GetConfigs(c *gin.Context) {
	configs, err := h.systemConfigService.GetAllConfigs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "获取配置失败",
		})
		return
	}

	// 返回配置，如果某些配置不存在则使用默认值
	result := map[string]interface{}{
		"sync_store_code": "390840", // 默认店铺代码
		"sync_enabled":    "false",   // 默认不启用自动同步
	}
	
	// 添加安全配置的默认值
	securityConfigs := h.systemConfigService.GetSecurityConfig()
	for key, value := range securityConfigs {
		result[key] = value
	}

	// 覆盖数据库中的配置
	for key, value := range configs {
		result[key] = value
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "成功",
		"data": result,
	})
}

// UpdateConfigs 更新系统配置
func (h *AdminSystemConfigHandler) UpdateConfigs(c *gin.Context) {
	var req map[string]string
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}

	// 更新配置
	for key, value := range req {
		var description string
		switch key {
		case "sync_store_code":
			description = "商品同步和卡片验证使用的店铺代码"
		case "sync_enabled":
			description = "是否启用自动商品同步"
		default:
			continue // 跳过未知的配置项
		}

		if err := h.systemConfigService.SetConfig(key, value, description); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": 500,
				"msg":  "更新配置失败: " + err.Error(),
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "配置更新成功",
	})
}

// GetStoreCode 获取店铺代码配置
func (h *AdminSystemConfigHandler) GetStoreCode(c *gin.Context) {
	storeCode, err := h.systemConfigService.GetSyncStoreCode()
	if err != nil || storeCode == "" {
		storeCode = "390840" // 默认值
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "成功",
		"data": gin.H{
			"store_code": storeCode,
		},
	})
}

// UpdateStoreCode 更新店铺代码配置
func (h *AdminSystemConfigHandler) UpdateStoreCode(c *gin.Context) {
	var req struct {
		StoreCode string `json:"store_code" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}

	if err := h.systemConfigService.SetSyncStoreCode(req.StoreCode); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "更新店铺代码失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "店铺代码更新成功",
	})
}
package handlers

import (
	"net/http"

	"backend/internal/services"
	"github.com/gin-gonic/gin"
)

type AdminDashboardHandler struct {
	dashboardService *services.DashboardService
}

func NewAdminDashboardHandler(dashboardService *services.DashboardService) *AdminDashboardHandler {
	return &AdminDashboardHandler{
		dashboardService: dashboardService,
	}
}

// GetStatistics 获取统计数据
func (h *AdminDashboardHandler) GetStatistics(c *gin.Context) {
	// 获取日期参数
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	
	stats, err := h.dashboardService.GetStatisticsWithDateRange(startDate, endDate)
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

// GetOrderTrend 获取订单趋势
func (h *AdminDashboardHandler) GetOrderTrend(c *gin.Context) {
	dateRange := c.Query("dateRange")
	days := 7

	// 支持直接传递天数或者预设的时间范围
	switch dateRange {
	case "7":
		days = 7
	case "30": 
		days = 30
	case "90":
		days = 90
	case "week":
		days = 7
	case "month":
		days = 30
	default:
		days = 7
	}

	trend, err := h.dashboardService.GetOrderTrend(days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "获取订单趋势失败",
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "成功",
		"data": trend,
	})
}

// GetHotGoods 获取热门商品
func (h *AdminDashboardHandler) GetHotGoods(c *gin.Context) {
	// 获取日期参数
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	
	hotGoods, err := h.dashboardService.GetHotGoodsWithDateRange(startDate, endDate, 10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "获取热门商品失败",
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "成功",
		"data": hotGoods,
	})
}

// GetRecentOrders 获取最新订单
func (h *AdminDashboardHandler) GetRecentOrders(c *gin.Context) {
	// 获取日期参数
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	
	orders, err := h.dashboardService.GetRecentOrdersWithDateRange(startDate, endDate, 10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "获取最新订单失败",
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "成功",
		"data": orders,
	})
}

// GetHourDistribution 获取小时分布
func (h *AdminDashboardHandler) GetHourDistribution(c *gin.Context) {
	// 获取24小时的订单分布数据
	distribution, err := h.dashboardService.GetHourDistribution()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "获取小时分布失败",
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "成功",
		"data": distribution,
	})
}

// GetDistributorRank 获取分销商排行
func (h *AdminDashboardHandler) GetDistributorRank(c *gin.Context) {
	period := c.Query("period")
	if period == "" {
		period = "today"
	}
	
	rank, err := h.dashboardService.GetDistributorRank(period)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "获取分销商排行失败",
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "成功",
		"data": rank,
	})
}

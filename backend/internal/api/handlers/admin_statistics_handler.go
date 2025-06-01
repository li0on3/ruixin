package handlers

import (
	"net/http"
	"strconv"
	"time"

	"backend/internal/services"
	"github.com/gin-gonic/gin"
)

type AdminStatisticsHandler struct {
	statisticsService *services.StatisticsService
}

func NewAdminStatisticsHandler(statisticsService *services.StatisticsService) *AdminStatisticsHandler {
	return &AdminStatisticsHandler{
		statisticsService: statisticsService,
	}
}

// GetMetrics 获取核心指标
func (h *AdminStatisticsHandler) GetMetrics(c *gin.Context) {
	startDate, endDate, err := h.parseDateRange(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "日期参数错误", "data": nil})
		return
	}

	metrics, err := h.statisticsService.GetMetrics(startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "获取指标失败", "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "Success", "data": metrics})
}

// GetSalesTrend 获取销售趋势
func (h *AdminStatisticsHandler) GetSalesTrend(c *gin.Context) {
	startDate, endDate, err := h.parseDateRange(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "日期参数错误", "data": nil})
		return
	}

	trendType := c.DefaultQuery("type", "both")
	
	trend, err := h.statisticsService.GetSalesTrend(startDate, endDate, trendType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "获取趋势失败", "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "Success", "data": trend})
}

// GetDistributorRank 获取分销商排行
func (h *AdminStatisticsHandler) GetDistributorRank(c *gin.Context) {
	startDate, endDate, err := h.parseDateRange(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "日期参数错误", "data": nil})
		return
	}

	rankType := c.DefaultQuery("type", "revenue")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if limit <= 0 {
		limit = 10
	}

	rank, err := h.statisticsService.GetDistributorRank(startDate, endDate, rankType, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "获取排行失败", "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "Success", "data": rank})
}

// GetProductAnalysis 获取商品分析
func (h *AdminStatisticsHandler) GetProductAnalysis(c *gin.Context) {
	startDate, endDate, err := h.parseDateRange(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "日期参数错误", "data": nil})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if limit <= 0 {
		limit = 10
	}

	analysis, err := h.statisticsService.GetProductAnalysis(startDate, endDate, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "获取分析失败", "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "Success", "data": analysis})
}

// GetHourDistribution 获取时段分布
func (h *AdminStatisticsHandler) GetHourDistribution(c *gin.Context) {
	startDate, endDate, err := h.parseDateRange(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "日期参数错误", "data": nil})
		return
	}

	distribution, err := h.statisticsService.GetHourDistribution(startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "获取分布失败", "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "Success", "data": distribution})
}

// GetRegionDistribution 获取地区分布
func (h *AdminStatisticsHandler) GetRegionDistribution(c *gin.Context) {
	startDate, endDate, err := h.parseDateRange(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "日期参数错误", "data": nil})
		return
	}

	distribution, err := h.statisticsService.GetRegionDistribution(startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "获取分布失败", "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "Success", "data": distribution})
}

// GetDetailData 获取详细数据
func (h *AdminStatisticsHandler) GetDetailData(c *gin.Context) {
	startDate, endDate, err := h.parseDateRange(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "日期参数错误", "data": nil})
		return
	}

	detailType := c.DefaultQuery("type", "daily")

	details, err := h.statisticsService.GetDetailData(startDate, endDate, detailType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "获取详情失败", "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "Success", "data": details})
}

// ExportData 导出数据
func (h *AdminStatisticsHandler) ExportData(c *gin.Context) {
	startDate, endDate, err := h.parseDateRange(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "日期参数错误", "data": nil})
		return
	}

	format := c.DefaultQuery("format", "excel")

	// TODO: 实现数据导出功能
	// 这里需要集成Excel导出库（如excelize）来生成Excel文件

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "导出功能开发中",
		"data": map[string]interface{}{
			"start_date": startDate.Format("2006-01-02"),
			"end_date":   endDate.Format("2006-01-02"),
			"format":     format,
		},
	})
}

// parseDateRange 解析日期范围参数
func (h *AdminStatisticsHandler) parseDateRange(c *gin.Context) (time.Time, time.Time, error) {
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	if startDateStr == "" || endDateStr == "" {
		// 默认最近7天
		endDate := time.Now()
		startDate := endDate.AddDate(0, 0, -7)
		return startDate, endDate, nil
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	// 将结束日期调整到当天的23:59:59
	endDate = endDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	return startDate, endDate, nil
}
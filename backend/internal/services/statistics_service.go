package services

import (
	"time"

	"backend/internal/repository"
)

type StatisticsService struct {
	orderRepo       *repository.OrderRepository
	distributorRepo *repository.DistributorRepository
}

func NewStatisticsService(orderRepo *repository.OrderRepository, distributorRepo *repository.DistributorRepository) *StatisticsService {
	return &StatisticsService{
		orderRepo:       orderRepo,
		distributorRepo: distributorRepo,
	}
}

// MetricsData 核心指标数据
type MetricsData struct {
	TotalRevenue     float64 `json:"totalRevenue"`
	RevenueGrowth    float64 `json:"revenueGrowth"`
	TotalOrders      int64   `json:"totalOrders"`
	OrderGrowth      float64 `json:"orderGrowth"`
	AvgOrderValue    float64 `json:"avgOrderValue"`
	AvgValueGrowth   float64 `json:"avgValueGrowth"`
	ConversionRate   float64 `json:"conversionRate"`
	ConversionGrowth float64 `json:"conversionGrowth"`
}

// GetMetrics 获取核心指标
func (s *StatisticsService) GetMetrics(startDate, endDate time.Time) (*MetricsData, error) {
	// 当前期间数据
	currentStats, err := s.orderRepo.GetStatsByDateRange(startDate, endDate)
	if err != nil {
		return nil, err
	}

	// 计算环比期间（相同天数的前一个时间段）
	duration := endDate.Sub(startDate)
	prevEndDate := startDate.Add(-time.Hour * 24)
	prevStartDate := prevEndDate.Add(-duration)

	// 上一期间数据
	prevStats, err := s.orderRepo.GetStatsByDateRange(prevStartDate, prevEndDate)
	if err != nil {
		return nil, err
	}

	metrics := &MetricsData{
		TotalRevenue: currentStats.TotalRevenue,
		TotalOrders:  currentStats.TotalOrders,
	}

	// 计算平均客单价
	if currentStats.TotalOrders > 0 {
		metrics.AvgOrderValue = currentStats.TotalRevenue / float64(currentStats.TotalOrders)
	}

	// 计算环比增长
	if prevStats.TotalRevenue > 0 {
		metrics.RevenueGrowth = (currentStats.TotalRevenue - prevStats.TotalRevenue) / prevStats.TotalRevenue
	}

	if prevStats.TotalOrders > 0 {
		metrics.OrderGrowth = float64(currentStats.TotalOrders-prevStats.TotalOrders) / float64(prevStats.TotalOrders)

		prevAvgValue := prevStats.TotalRevenue / float64(prevStats.TotalOrders)
		if prevAvgValue > 0 {
			metrics.AvgValueGrowth = (metrics.AvgOrderValue - prevAvgValue) / prevAvgValue
		}
	}

	// TODO: 计算转化率需要访问量数据，暂时返回模拟数据
	metrics.ConversionRate = 0.15
	metrics.ConversionGrowth = 0.05

	return metrics, nil
}

// SalesTrendData 销售趋势数据
type SalesTrendData struct {
	Date     string  `json:"date"`
	Revenue  float64 `json:"revenue"`
	Orders   int64   `json:"orders"`
	AvgValue float64 `json:"avgValue"`
}

// GetSalesTrend 获取销售趋势
func (s *StatisticsService) GetSalesTrend(startDate, endDate time.Time, trendType string) ([]SalesTrendData, error) {
	dailyStats, err := s.orderRepo.GetDailyStatsByDateRange(startDate, endDate)
	if err != nil {
		return nil, err
	}

	trend := make([]SalesTrendData, 0, len(dailyStats))
	for _, stat := range dailyStats {
		data := SalesTrendData{
			Date:    stat.Date.Format("2006-01-02"),
			Revenue: stat.Revenue,
			Orders:  stat.Orders,
		}
		if stat.Orders > 0 {
			data.AvgValue = stat.Revenue / float64(stat.Orders)
		}
		trend = append(trend, data)
	}

	return trend, nil
}

// DistributorRankData 分销商排行数据
type DistributorRankData struct {
	DistributorID   uint    `json:"distributorId"`
	DistributorName string  `json:"distributorName"`
	Revenue         float64 `json:"revenue"`
	Orders          int64   `json:"orders"`
	AvgValue        float64 `json:"avgValue"`
}

// GetDistributorRank 获取分销商排行
func (s *StatisticsService) GetDistributorRank(startDate, endDate time.Time, rankType string, limit int) ([]DistributorRankData, error) {
	stats, err := s.orderRepo.GetDistributorStatsByDateRange(startDate, endDate, limit)
	if err != nil {
		return nil, err
	}

	rank := make([]DistributorRankData, 0, len(stats))
	for _, stat := range stats {
		data := DistributorRankData{
			DistributorID:   stat.DistributorID,
			DistributorName: stat.DistributorName,
			Revenue:         stat.Revenue,
			Orders:          stat.Orders,
		}
		if stat.Orders > 0 {
			data.AvgValue = stat.Revenue / float64(stat.Orders)
		}
		rank = append(rank, data)
	}

	return rank, nil
}

// ProductAnalysisData 商品分析数据
type ProductAnalysisData struct {
	ProductName string  `json:"productName"`
	Quantity    int64   `json:"quantity"`
	Revenue     float64 `json:"revenue"`
	Percentage  float64 `json:"percentage"`
}

// GetProductAnalysis 获取商品分析
func (s *StatisticsService) GetProductAnalysis(startDate, endDate time.Time, limit int) ([]ProductAnalysisData, error) {
	stats, err := s.orderRepo.GetProductStatsByDateRange(startDate, endDate, limit)
	if err != nil {
		return nil, err
	}

	// 计算总销售额
	var totalRevenue float64
	for _, stat := range stats {
		totalRevenue += stat.Revenue
	}

	analysis := make([]ProductAnalysisData, 0, len(stats))
	for _, stat := range stats {
		data := ProductAnalysisData{
			ProductName: stat.ProductName,
			Quantity:    stat.Quantity,
			Revenue:     stat.Revenue,
		}
		if totalRevenue > 0 {
			data.Percentage = stat.Revenue / totalRevenue
		}
		analysis = append(analysis, data)
	}

	return analysis, nil
}

// HourDistributionData 时段分布数据
type HourDistributionData struct {
	Hour   int   `json:"hour"`
	Orders int64 `json:"orders"`
}

// GetHourDistribution 获取时段分布
func (s *StatisticsService) GetHourDistribution(startDate, endDate time.Time) ([]HourDistributionData, error) {
	stats, err := s.orderRepo.GetHourDistributionByDateRange(startDate, endDate)
	if err != nil {
		return nil, err
	}

	distribution := make([]HourDistributionData, 24)
	for i := 0; i < 24; i++ {
		distribution[i] = HourDistributionData{
			Hour:   i,
			Orders: 0,
		}
	}

	for _, stat := range stats {
		if stat.Hour >= 0 && stat.Hour < 24 {
			distribution[stat.Hour].Orders = stat.Orders
		}
	}

	return distribution, nil
}

// RegionDistributionData 地区分布数据
type RegionDistributionData struct {
	Region  string  `json:"region"`
	Orders  int64   `json:"orders"`
	Revenue float64 `json:"revenue"`
}

// GetRegionDistribution 获取地区分布
func (s *StatisticsService) GetRegionDistribution(startDate, endDate time.Time) ([]RegionDistributionData, error) {
	// TODO: 实现地区分布统计，需要从门店地址中提取地区信息
	// 暂时返回模拟数据
	return []RegionDistributionData{
		{Region: "深圳", Orders: 1000, Revenue: 15000},
		{Region: "广州", Orders: 800, Revenue: 12000},
		{Region: "上海", Orders: 600, Revenue: 9000},
		{Region: "北京", Orders: 500, Revenue: 7500},
		{Region: "杭州", Orders: 400, Revenue: 6000},
	}, nil
}

// DetailData 详细数据 - 已移除，使用 map[string]interface{} 代替

// GetDetailData 获取详细数据
func (s *StatisticsService) GetDetailData(startDate, endDate time.Time, detailType string) ([]map[string]interface{}, error) {
	switch detailType {
	case "daily":
		stats, err := s.orderRepo.GetDailyStatsByDateRange(startDate, endDate)
		if err != nil {
			return nil, err
		}
		details := make([]map[string]interface{}, 0, len(stats))
		for _, stat := range stats {
			avgOrderValue := float64(0)
			if stat.Orders > 0 {
				avgOrderValue = stat.Revenue / float64(stat.Orders)
			}
			
			// 计算模拟数据
			newCustomers := stat.Orders * 30 / 100 // 假设30%是新客户
			conversionRate := 0.15 + float64(stat.Orders%10)*0.01 // 模拟转化率
			profit := stat.Revenue * 0.3 // 假设30%利润率
			profitRate := 0.3
			
			details = append(details, map[string]interface{}{
				"date":           stat.Date.Format("2006-01-02"),
				"orderCount":     stat.Orders,
				"revenue":        stat.Revenue,
				"avgOrderValue":  avgOrderValue,
				"newCustomers":   newCustomers,
				"conversionRate": conversionRate,
				"profit":         profit,
				"profitRate":     profitRate,
			})
		}
		return details, nil
		
	case "distributor":
		stats, err := s.orderRepo.GetDistributorStatsByDateRange(startDate, endDate, 0)
		if err != nil {
			return nil, err
		}
		
		// 计算总销售额
		var totalRevenue float64
		for _, stat := range stats {
			totalRevenue += stat.Revenue
		}
		
		details := make([]map[string]interface{}, 0, len(stats))
		for _, stat := range stats {
			avgOrderValue := float64(0)
			if stat.Orders > 0 {
				avgOrderValue = stat.Revenue / float64(stat.Orders)
			}
			
			percentage := float64(0)
			if totalRevenue > 0 {
				percentage = (stat.Revenue / totalRevenue) * 100
			}
			
			commission := stat.Revenue * 0.1 // 假设10%佣金
			conversionRate := 0.15 + float64(stat.Orders%10)*0.01
			status := 1 // 默认活跃
			if stat.Orders < 5 {
				status = 0 // 休眠
			}
			
			details = append(details, map[string]interface{}{
				"distributorName": stat.DistributorName,
				"orderCount":      stat.Orders,
				"revenue":         stat.Revenue,
				"commission":      commission,
				"avgOrderValue":   avgOrderValue,
				"conversionRate":  conversionRate,
				"percentage":      percentage,
				"status":          status,
			})
		}
		return details, nil
		
	case "product":
		stats, err := s.orderRepo.GetProductStatsByDateRange(startDate, endDate, 0)
		if err != nil {
			return nil, err
		}
		
		// 计算总销售额
		var totalRevenue float64
		for _, stat := range stats {
			totalRevenue += stat.Revenue
		}
		
		details := make([]map[string]interface{}, 0, len(stats))
		for _, stat := range stats {
			avgPrice := float64(0)
			if stat.Quantity > 0 {
				avgPrice = stat.Revenue / float64(stat.Quantity)
			}
			
			percentage := float64(0)
			if totalRevenue > 0 {
				percentage = (stat.Revenue / totalRevenue) * 100
			}
			
			cost := stat.Revenue * 0.6 // 假设60%成本
			profit := stat.Revenue - cost
			profitRate := 0.4
			
			details = append(details, map[string]interface{}{
				"productName": stat.ProductName,
				"category":    "咖啡", // 需要从产品信息中获取
				"quantity":    stat.Quantity,
				"revenue":     stat.Revenue,
				"cost":        cost,
				"profit":      profit,
				"avgPrice":    avgPrice,
				"profitRate":  profitRate,
				"percentage":  percentage,
			})
		}
		return details, nil
		
	case "store":
		// TODO: 实现门店统计，暂时返回模拟数据
		stores := []map[string]interface{}{
			{
				"storeName": "深圳福田店",
				"storeCode": "SZ001",
				"city": "深圳",
				"district": "福田区",
				"orderCount": 150,
				"revenue": 3500.00,
				"avgOrderValue": 23.33,
				"activeRate": 0.85,
				"percentage": 25.5,
			},
			{
				"storeName": "深圳南山店",
				"storeCode": "SZ002",
				"city": "深圳",
				"district": "南山区",
				"orderCount": 120,
				"revenue": 2800.00,
				"avgOrderValue": 23.33,
				"activeRate": 0.75,
				"percentage": 20.4,
			},
		}
		return stores, nil
		
	default:
		return []map[string]interface{}{}, nil
	}
}

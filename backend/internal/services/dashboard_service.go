package services

import (
	"backend/internal/repository"
	"fmt"
	"time"
)

type DashboardService struct {
	orderRepo       *repository.OrderRepository
	cardRepo        *repository.CardRepository
	distributorRepo *repository.DistributorRepository
}

func NewDashboardService(
	orderRepo *repository.OrderRepository,
	cardRepo *repository.CardRepository,
	distributorRepo *repository.DistributorRepository,
) *DashboardService {
	return &DashboardService{
		orderRepo:       orderRepo,
		cardRepo:        cardRepo,
		distributorRepo: distributorRepo,
	}
}

// GetStatistics 获取仪表盘统计数据
func (s *DashboardService) GetStatistics() (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// 今日统计
	today := time.Now().Truncate(24 * time.Hour)
	tomorrow := today.Add(24 * time.Hour)
	todayStats, err := s.orderRepo.GetStatistics(0, today, tomorrow)
	if err != nil {
		return nil, err
	}

	// 昨日统计
	yesterday := today.Add(-24 * time.Hour)
	yesterdayStats, err := s.orderRepo.GetStatistics(0, yesterday, today)
	if err != nil {
		return nil, err
	}

	// 计算增长率
	orderGrowth := 0.0
	if yesterdayOrders, ok := yesterdayStats["total_orders"].(int64); ok && yesterdayOrders > 0 {
		if todayOrders, ok := todayStats["total_orders"].(int64); ok {
			orderGrowth = float64(todayOrders-yesterdayOrders) / float64(yesterdayOrders) * 100
		}
	}

	amountGrowth := 0.0
	if yesterdayAmount, ok := yesterdayStats["total_amount"].(float64); ok && yesterdayAmount > 0 {
		if todayAmount, ok := todayStats["total_amount"].(float64); ok {
			amountGrowth = (todayAmount - yesterdayAmount) / yesterdayAmount * 100
		}
	}

	// 活跃分销商统计
	activeDistributors, _, err := s.distributorRepo.List(0, 1000, map[string]interface{}{"status": 1})
	_, totalDistributors, _ := s.distributorRepo.List(0, 1, nil)

	// 卡片统计 - status=0 表示未使用（活跃）
	activeCards, _, err := s.cardRepo.List(0, 1000, map[string]interface{}{"status": 0})
	if err != nil {
		return nil, err
	}
	// 获取总卡片数，这里应该使用返回的第二个值（total）
	_, totalCardsCount, err := s.cardRepo.List(0, 1, nil)
	if err != nil {
		return nil, err
	}

	// 确保返回的是数字，不是其他数据
	var todayOrdersCount int64 = 0
	var todayAmountValue float64 = 0

	if v, ok := todayStats["total_orders"].(int64); ok {
		todayOrdersCount = v
	}
	if v, ok := todayStats["total_amount"].(float64); ok {
		todayAmountValue = v
	}

	stats["todayOrders"] = todayOrdersCount
	stats["todayAmount"] = todayAmountValue
	stats["orderGrowth"] = orderGrowth
	stats["amountGrowth"] = amountGrowth
	stats["activeDistributors"] = len(activeDistributors)
	stats["totalDistributors"] = totalDistributors
	stats["activeCards"] = len(activeCards)
	stats["totalCards"] = totalCardsCount

	return stats, nil
}

// GetStatisticsWithDateRange 获取指定日期范围的统计数据
func (s *DashboardService) GetStatisticsWithDateRange(startDateStr, endDateStr string) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// 解析日期参数
	var startDate, endDate time.Time
	var err error

	if startDateStr != "" && endDateStr != "" {
		startDate, err = time.Parse("2006-01-02", startDateStr)
		if err != nil {
			return nil, fmt.Errorf("解析开始日期失败: %v", err)
		}
		endDate, err = time.Parse("2006-01-02", endDateStr)
		if err != nil {
			return nil, fmt.Errorf("解析结束日期失败: %v", err)
		}
		// 将结束日期调整到当天的23:59:59
		endDate = endDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
	} else {
		// 默认返回今日数据
		return s.GetStatistics()
	}

	// 获取日期范围内的统计
	periodStats, err := s.orderRepo.GetStatistics(0, startDate, endDate)
	if err != nil {
		return nil, err
	}

	// 获取前一个相同时长的统计数据用于计算增长率
	duration := endDate.Sub(startDate)
	prevStartDate := startDate.Add(-duration)
	prevEndDate := startDate
	prevStats, err := s.orderRepo.GetStatistics(0, prevStartDate, prevEndDate)
	if err != nil {
		// 如果获取前期数据失败，忽略错误，增长率设为0
		prevStats = map[string]interface{}{
			"total_orders": int64(0),
			"total_amount": float64(0),
		}
	}

	// 计算增长率
	orderGrowth := 0.0
	if prevOrders, ok := prevStats["total_orders"].(int64); ok && prevOrders > 0 {
		if currOrders, ok := periodStats["total_orders"].(int64); ok {
			orderGrowth = float64(currOrders-prevOrders) / float64(prevOrders) * 100
		}
	}

	amountGrowth := 0.0
	if prevAmount, ok := prevStats["total_amount"].(float64); ok && prevAmount > 0 {
		if currAmount, ok := periodStats["total_amount"].(float64); ok {
			amountGrowth = (currAmount - prevAmount) / prevAmount * 100
		}
	}

	// 活跃分销商统计（这个不受日期范围影响）
	activeDistributors, _, _ := s.distributorRepo.List(0, 1000, map[string]interface{}{"status": 1})
	_, totalDistributors, _ := s.distributorRepo.List(0, 1, nil)

	// 卡片统计（这个不受日期范围影响）
	activeCards, _, _ := s.cardRepo.List(0, 1000, map[string]interface{}{"status": 0})
	_, totalCardsCount, _ := s.cardRepo.List(0, 1, nil)

	// 确保返回的是数字，不是其他数据
	var periodOrdersCount int64 = 0
	var periodAmountValue float64 = 0

	if v, ok := periodStats["total_orders"].(int64); ok {
		periodOrdersCount = v
	}
	if v, ok := periodStats["total_amount"].(float64); ok {
		periodAmountValue = v
	}

	// 根据日期范围调整标签
	isToday := startDateStr == time.Now().Format("2006-01-02") && endDateStr == time.Now().Format("2006-01-02")
	
	stats["todayOrders"] = periodOrdersCount  // 保持字段名不变以兼容前端
	stats["todayAmount"] = periodAmountValue  // 保持字段名不变以兼容前端
	stats["orderGrowth"] = orderGrowth
	stats["amountGrowth"] = amountGrowth
	stats["activeDistributors"] = len(activeDistributors)
	stats["totalDistributors"] = totalDistributors
	stats["activeCards"] = len(activeCards)
	stats["totalCards"] = totalCardsCount
	stats["dateRange"] = map[string]string{
		"start": startDateStr,
		"end": endDateStr,
		"label": formatDateRangeLabel(startDate, endDate, isToday),
	}

	return stats, nil
}

// formatDateRangeLabel 格式化日期范围标签
func formatDateRangeLabel(startDate, endDate time.Time, isToday bool) string {
	if isToday {
		return "今日"
	}
	if startDate.Format("2006-01-02") == endDate.Format("2006-01-02") {
		return startDate.Format("01月02日")
	}
	return fmt.Sprintf("%s 至 %s", startDate.Format("01月02日"), endDate.Format("01月02日"))
}

// GetOrderTrend 获取订单趋势
func (s *DashboardService) GetOrderTrend(days int) ([]map[string]interface{}, error) {
	result, err := s.orderRepo.GetDailyStatistics(0, days)
	if err != nil {
		fmt.Printf("GetOrderTrend error: %v\n", err)
		return nil, err
	}
	fmt.Printf("GetOrderTrend result count: %d\n", len(result))
	return result, nil
}

// GetHotGoods 获取热门商品
func (s *DashboardService) GetHotGoods(limit int) ([]map[string]interface{}, error) {
	// 获取最近30天的商品统计数据
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -30) // 最近30天
	
	productStats, err := s.orderRepo.GetProductStatsByDateRange(startDate, endDate, 0)
	if err != nil {
		fmt.Printf("GetHotGoods error: %v\n", err)
		return nil, err
	}
	
	// 如果没有真实数据，返回空数组而不是假数据
	if len(productStats) == 0 {
		fmt.Printf("GetHotGoods: no product stats found\n")
		return []map[string]interface{}{}, nil
	}
	
	// 计算总数量用于计算百分比
	var totalCount int64 = 0
	for _, stat := range productStats {
		totalCount += stat.Quantity
	}
	
	// 转换为返回格式
	result := make([]map[string]interface{}, 0, len(productStats))
	for _, stat := range productStats {
		percentage := 0.0
		if totalCount > 0 {
			percentage = float64(stat.Quantity) / float64(totalCount) * 100
		}
		
		result = append(result, map[string]interface{}{
			"name":       stat.ProductName,
			"count":      stat.Quantity,
			"percentage": fmt.Sprintf("%.1f", percentage),
		})
	}
	
	// 应用限制
	if limit > 0 && limit < len(result) {
		return result[:limit], nil
	}
	
	fmt.Printf("GetHotGoods: returning %d products\n", len(result))
	return result, nil
}

// GetHotGoodsWithDateRange 获取指定日期范围的热门商品
func (s *DashboardService) GetHotGoodsWithDateRange(startDateStr, endDateStr string, limit int) ([]map[string]interface{}, error) {
	// 解析日期参数
	var startDate, endDate time.Time
	var err error

	if startDateStr != "" && endDateStr != "" {
		startDate, err = time.Parse("2006-01-02", startDateStr)
		if err != nil {
			return nil, fmt.Errorf("解析开始日期失败: %v", err)
		}
		endDate, err = time.Parse("2006-01-02", endDateStr)
		if err != nil {
			return nil, fmt.Errorf("解析结束日期失败: %v", err)
		}
		// 将结束日期调整到当天的23:59:59
		endDate = endDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
	} else {
		// 默认最近30天
		return s.GetHotGoods(limit)
	}
	
	productStats, err := s.orderRepo.GetProductStatsByDateRange(startDate, endDate, 0)
	if err != nil {
		fmt.Printf("GetHotGoodsWithDateRange error: %v\n", err)
		return nil, err
	}
	
	// 如果没有真实数据，返回空数组
	if len(productStats) == 0 {
		fmt.Printf("GetHotGoodsWithDateRange: no product stats found for date range %s to %s\n", startDateStr, endDateStr)
		return []map[string]interface{}{}, nil
	}
	
	// 计算总数量用于计算百分比
	var totalCount int64 = 0
	for _, stat := range productStats {
		totalCount += stat.Quantity
	}
	
	// 转换为返回格式
	result := make([]map[string]interface{}, 0, len(productStats))
	for _, stat := range productStats {
		percentage := 0.0
		if totalCount > 0 {
			percentage = float64(stat.Quantity) / float64(totalCount) * 100
		}
		
		result = append(result, map[string]interface{}{
			"name":       stat.ProductName,
			"count":      stat.Quantity,
			"percentage": fmt.Sprintf("%.1f", percentage),
		})
	}
	
	// 应用限制
	if limit > 0 && limit < len(result) {
		return result[:limit], nil
	}
	
	fmt.Printf("GetHotGoodsWithDateRange: returning %d products for date range %s to %s\n", len(result), startDateStr, endDateStr)
	return result, nil
}

// GetRecentOrders 获取最新订单
func (s *DashboardService) GetRecentOrders(limit int) ([]map[string]interface{}, error) {
	orders, total, err := s.orderRepo.List(0, limit, nil)
	if err != nil {
		fmt.Printf("GetRecentOrders error: %v\n", err)
		return nil, err
	}
	fmt.Printf("GetRecentOrders: found %d orders (total: %d)\n", len(orders), total)

	result := make([]map[string]interface{}, 0, len(orders))
	for _, order := range orders {
		// 获取分销商信息
		distributor, _ := s.distributorRepo.GetByID(order.DistributorID)
		distributorName := ""
		if distributor != nil {
			distributorName = distributor.Name
		}

		result = append(result, map[string]interface{}{
			"id":              order.ID,
			"orderNo":         order.OrderNo,
			"distributorName": distributorName,
			"storeName":       order.StoreName,
			"totalAmount":     order.TotalAmount,
			"status":          order.Status,
			"createdAt":       order.CreatedAt,
		})
	}

	return result, nil
}

// GetRecentOrdersWithDateRange 获取指定日期范围的最新订单
func (s *DashboardService) GetRecentOrdersWithDateRange(startDateStr, endDateStr string, limit int) ([]map[string]interface{}, error) {
	var filters map[string]interface{}
	
	// 解析日期参数
	if startDateStr != "" && endDateStr != "" {
		startDate, err := time.Parse("2006-01-02", startDateStr)
		if err != nil {
			return nil, fmt.Errorf("解析开始日期失败: %v", err)
		}
		endDate, err := time.Parse("2006-01-02", endDateStr)
		if err != nil {
			return nil, fmt.Errorf("解析结束日期失败: %v", err)
		}
		// 将结束日期调整到当天的23:59:59
		endDate = endDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
		
		filters = map[string]interface{}{
			"start_date": startDate,
			"end_date":   endDate,
		}
	}
	
	orders, total, err := s.orderRepo.List(0, limit, filters)
	if err != nil {
		fmt.Printf("GetRecentOrdersWithDateRange error: %v\n", err)
		return nil, err
	}
	fmt.Printf("GetRecentOrdersWithDateRange: found %d orders (total: %d) for date range %s to %s\n", len(orders), total, startDateStr, endDateStr)

	result := make([]map[string]interface{}, 0, len(orders))
	for _, order := range orders {
		// 获取分销商信息
		distributor, _ := s.distributorRepo.GetByID(order.DistributorID)
		distributorName := ""
		if distributor != nil {
			distributorName = distributor.Name
		}

		result = append(result, map[string]interface{}{
			"id":              order.ID,
			"orderNo":         order.OrderNo,
			"distributorName": distributorName,
			"storeName":       order.StoreName,
			"totalAmount":     order.TotalAmount,
			"status":          order.Status,
			"createdAt":       order.CreatedAt,
		})
	}

	return result, nil
}

// GetHourDistribution 获取小时分布
func (s *DashboardService) GetHourDistribution() ([]int, error) {
	// 获取今日的小时分布数据
	distribution, err := s.orderRepo.GetHourDistribution()
	if err != nil {
		fmt.Printf("GetHourDistribution error: %v\n", err)
		return nil, err
	}
	
	// 确保返回24小时的数据（0-23时）
	result := make([]int, 24)
	for _, item := range distribution {
		if hour, ok := item["hour"].(int); ok && hour >= 0 && hour < 24 {
			if count, ok := item["count"].(int64); ok {
				result[hour] = int(count)
			}
		}
	}
	
	fmt.Printf("GetHourDistribution: returning 24-hour distribution data\n")
	return result, nil
}

// GetDistributorRank 获取分销商排行
func (s *DashboardService) GetDistributorRank(period string) ([]map[string]interface{}, error) {
	// 根据时间周期计算时间范围
	var startTime, endTime time.Time
	now := time.Now()
	
	switch period {
	case "today":
		startTime = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		endTime = startTime.AddDate(0, 0, 1)
	case "week":
		// 本周一开始
		weekday := int(now.Weekday())
		if weekday == 0 { // 周日
			weekday = 7
		}
		startTime = now.AddDate(0, 0, -(weekday-1)).Truncate(24 * time.Hour)
		endTime = now
	case "month":
		// 本月开始
		startTime = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		endTime = now
	default:
		// 默认今日
		startTime = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		endTime = startTime.AddDate(0, 0, 1)
	}
	
	// 获取分销商排行数据
	rank, err := s.orderRepo.GetDistributorRank(startTime, endTime)
	if err != nil {
		fmt.Printf("GetDistributorRank error: %v\n", err)
		return nil, err
	}
	
	result := make([]map[string]interface{}, 0, len(rank))
	for _, item := range rank {
		// 获取分销商信息
		distributorID := item["distributor_id"].(uint)
		distributor, _ := s.distributorRepo.GetByID(distributorID)
		distributorName := fmt.Sprintf("分销商%d", distributorID)
		if distributor != nil {
			distributorName = distributor.Name
		}
		
		result = append(result, map[string]interface{}{
			"id":      distributorID,
			"name":    distributorName,
			"orders":  item["order_count"],
			"revenue": item["total_amount"],
		})
	}
	
	fmt.Printf("GetDistributorRank: returning %d distributors for period %s\n", len(result), period)
	return result, nil
}

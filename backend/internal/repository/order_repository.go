package repository

import (
	"backend/internal/models"
	"database/sql"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type OrderRepository struct {
	*BaseRepository
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

func (r *OrderRepository) Create(order *models.Order) error {
	return r.db.Create(order).Error
}

func (r *OrderRepository) Update(order *models.Order) error {
	return r.db.Save(order).Error
}

func (r *OrderRepository) GetByID(id uint) (*models.Order, error) {
	var order models.Order
	err := r.db.First(&order, id).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *OrderRepository) GetByOrderNo(orderNo string) (*models.Order, error) {
	var order models.Order
	err := r.db.Where("order_no = ?", orderNo).First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *OrderRepository) List(offset, limit int, filters map[string]interface{}) ([]*models.Order, int64, error) {
	var orders []*models.Order
	var total int64

	query := r.db.Model(&models.Order{})

	// Apply filters
	if distributorID, ok := filters["distributor_id"]; ok {
		query = query.Where("distributor_id = ?", distributorID)
	}
	if status, ok := filters["status"]; ok {
		query = query.Where("status = ?", status)
	}
	if storeCode, ok := filters["store_code"]; ok {
		query = query.Where("store_code = ?", storeCode)
	}
	if cardCode, ok := filters["card_code"]; ok {
		query = query.Where("card_code = ?", cardCode)
	}
	if startDate, ok := filters["start_date"]; ok {
		query = query.Where("created_at >= ?", startDate)
	}
	if endDate, ok := filters["end_date"]; ok {
		query = query.Where("created_at <= ?", endDate)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get data
	err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&orders).Error
	if err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

func (r *OrderRepository) GetStatistics(distributorID uint, startDate, endDate time.Time) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	query := r.db.Model(&models.Order{})
	if distributorID > 0 {
		query = query.Where("distributor_id = ?", distributorID)
	}
	query = query.Where("created_at BETWEEN ? AND ?", startDate, endDate)

	// Total orders
	var totalOrders int64
	query.Count(&totalOrders)
	stats["total_orders"] = totalOrders

	// Success orders
	var successOrders int64
	query.Where("status = ?", models.OrderStatusSuccess).Count(&successOrders)
	stats["success_orders"] = successOrders

	// Total amount
	var totalAmount sql.NullFloat64
	query.Select("COALESCE(SUM(total_amount), 0)").Scan(&totalAmount)
	stats["total_amount"] = totalAmount.Float64

	// Total profit
	var totalProfit sql.NullFloat64
	query.Select("COALESCE(SUM(profit_amount), 0)").Scan(&totalProfit)
	stats["total_profit"] = totalProfit.Float64

	// Success rate
	if totalOrders > 0 {
		stats["success_rate"] = float64(successOrders) / float64(totalOrders) * 100
	} else {
		stats["success_rate"] = 0
	}

	return stats, nil
}

func (r *OrderRepository) GetDailyStatistics(distributorID uint, days int) ([]map[string]interface{}, error) {
	results := make([]map[string]interface{}, 0)

	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -days)

	fmt.Printf("GetDailyStatistics: startDate=%v, endDate=%v\n", startDate, endDate)

	query := r.db.Model(&models.Order{})
	if distributorID > 0 {
		query = query.Where("distributor_id = ?", distributorID)
	}

	// 先检查是否有数据
	var count int64
	query.Where("created_at BETWEEN ? AND ?", startDate, endDate).Count(&count)
	fmt.Printf("GetDailyStatistics: total orders in date range = %d\n", count)

	rows, err := query.
		Select("DATE(created_at) as date, COUNT(*) as total_orders, "+
			"SUM(CASE WHEN status = ? THEN 1 ELSE 0 END) as success_orders, "+
			"COALESCE(SUM(total_amount), 0) as total_amount, COALESCE(SUM(profit_amount), 0) as total_profit",
			models.OrderStatusSuccess).
		Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Group("DATE(created_at)").
		Order("date DESC").
		Rows()

	if err != nil {
		fmt.Printf("GetDailyStatistics query error: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var date string
		var totalOrders, successOrders int64
		var totalAmount, totalProfit sql.NullFloat64

		if err := rows.Scan(&date, &totalOrders, &successOrders, &totalAmount, &totalProfit); err != nil {
			continue
		}

		successRate := 0.0
		if totalOrders > 0 {
			successRate = float64(successOrders) / float64(totalOrders) * 100
		}

		results = append(results, map[string]interface{}{
			"date":           date,
			"total_orders":   totalOrders,
			"success_orders": successOrders,
			"total_amount":   totalAmount.Float64,
			"total_profit":   totalProfit.Float64,
			"success_rate":   successRate,
		})
	}

	return results, nil
}

// 统计相关结构体
type OrderStats struct {
	TotalRevenue float64
	TotalOrders  int64
}

type DailyStats struct {
	Date    time.Time
	Revenue float64
	Orders  int64
}

type DistributorStats struct {
	DistributorID   uint
	DistributorName string
	Revenue         float64
	Orders          int64
}

type ProductStats struct {
	ProductName string
	Quantity    int64
	Revenue     float64
}

type HourStats struct {
	Hour   int
	Orders int64
}

// GetStatsByDateRange 获取指定日期范围的统计数据
func (r *OrderRepository) GetStatsByDateRange(startDate, endDate time.Time) (*OrderStats, error) {
	var stats OrderStats
	
	err := r.db.Model(&models.Order{}).
		Select("COALESCE(SUM(total_amount), 0) as total_revenue, COUNT(*) as total_orders").
		Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Where("status = ?", models.OrderStatusSuccess).
		Scan(&stats).Error
		
	return &stats, err
}

// GetDailyStatsByDateRange 获取按日统计数据
func (r *OrderRepository) GetDailyStatsByDateRange(startDate, endDate time.Time) ([]DailyStats, error) {
	var stats []DailyStats
	
	err := r.db.Model(&models.Order{}).
		Select("DATE(created_at) as date, SUM(total_amount) as revenue, COUNT(*) as orders").
		Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Where("status = ?", models.OrderStatusSuccess).
		Group("DATE(created_at)").
		Order("date ASC").
		Find(&stats).Error
		
	return stats, err
}

// GetDistributorStatsByDateRange 获取分销商统计数据
func (r *OrderRepository) GetDistributorStatsByDateRange(startDate, endDate time.Time, limit int) ([]DistributorStats, error) {
	var stats []DistributorStats
	
	query := r.db.Table("orders o").
		Select("o.distributor_id, d.name as distributor_name, SUM(o.total_amount) as revenue, COUNT(*) as orders").
		Joins("LEFT JOIN distributors d ON o.distributor_id = d.id").
		Where("o.created_at BETWEEN ? AND ?", startDate, endDate).
		Where("o.status = ?", models.OrderStatusSuccess).
		Group("o.distributor_id, d.name").
		Order("revenue DESC")
		
	if limit > 0 {
		query = query.Limit(limit)
	}
	
	err := query.Find(&stats).Error
	return stats, err
}

// GetProductStatsByDateRange 获取产品统计数据
func (r *OrderRepository) GetProductStatsByDateRange(startDate, endDate time.Time, limit int) ([]ProductStats, error) {
	var stats []ProductStats
	
	// 由于商品信息存储在JSON字段中，需要使用JSON_EXTRACT提取商品名称
	// 注意：这里假设每个订单只有一个商品，如果有多个商品需要更复杂的处理
	query := r.db.Model(&models.Order{}).
		Select(`JSON_UNQUOTE(JSON_EXTRACT(goods, '$[0].goods_name')) as product_name, 
				COUNT(*) as quantity, 
				SUM(total_amount) as revenue`).
		Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Where("status = ?", models.OrderStatusSuccess).
		Where("JSON_EXTRACT(goods, '$[0].goods_name') IS NOT NULL").
		Group("JSON_UNQUOTE(JSON_EXTRACT(goods, '$[0].goods_name'))").
		Order("revenue DESC")
		
	if limit > 0 {
		query = query.Limit(limit)
	}
	
	err := query.Find(&stats).Error
	return stats, err
}

// GetHourDistributionByDateRange 获取时段分布数据
func (r *OrderRepository) GetHourDistributionByDateRange(startDate, endDate time.Time) ([]HourStats, error) {
	var stats []HourStats
	
	err := r.db.Model(&models.Order{}).
		Select("HOUR(created_at) as hour, COUNT(*) as orders").
		Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Where("status = ?", models.OrderStatusSuccess).
		Group("HOUR(created_at)").
		Order("hour ASC").
		Find(&stats).Error
		
	return stats, err
}

// GetHourDistribution 获取今日小时分布数据
func (r *OrderRepository) GetHourDistribution() ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	
	// 获取今日的时间范围
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)
	
	rows, err := r.db.Model(&models.Order{}).
		Select("HOUR(created_at) as hour, COUNT(*) as count").
		Where("created_at BETWEEN ? AND ?", startOfDay, endOfDay).
		Group("HOUR(created_at)").
		Order("hour ASC").
		Rows()
		
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	for rows.Next() {
		var hour int
		var count int64
		
		if err := rows.Scan(&hour, &count); err != nil {
			continue
		}
		
		results = append(results, map[string]interface{}{
			"hour":  hour,
			"count": count,
		})
	}
	
	return results, nil
}

// GetDistributorRank 获取分销商排行
func (r *OrderRepository) GetDistributorRank(startTime, endTime time.Time) ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	
	rows, err := r.db.Table("orders").
		Select("distributor_id, COUNT(*) as order_count, COALESCE(SUM(total_amount), 0) as total_amount").
		Where("created_at BETWEEN ? AND ?", startTime, endTime).
		Where("status = ?", models.OrderStatusSuccess).
		Group("distributor_id").
		Order("total_amount DESC").
		Limit(10).
		Rows()
		
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	for rows.Next() {
		var distributorID uint
		var orderCount int64
		var totalAmount sql.NullFloat64
		
		if err := rows.Scan(&distributorID, &orderCount, &totalAmount); err != nil {
			continue
		}
		
		results = append(results, map[string]interface{}{
			"distributor_id": distributorID,
			"order_count":    orderCount,
			"total_amount":   totalAmount.Float64,
		})
	}
	
	return results, nil
}

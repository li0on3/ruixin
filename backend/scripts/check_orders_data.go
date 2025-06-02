package main

import (
	"backend/internal/config"
	"backend/internal/models"
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// 加载配置
	if err := config.Init("../configs/config.yaml"); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 连接数据库
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.AppConfig.Database.Username,
		config.AppConfig.Database.Password,
		config.AppConfig.Database.Host,
		config.AppConfig.Database.Port,
		config.AppConfig.Database.Database,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 1. 检查订单总数
	var totalCount int64
	db.Model(&models.Order{}).Count(&totalCount)
	fmt.Printf("订单总数: %d\n\n", totalCount)

	// 2. 按状态统计订单
	fmt.Println("按状态统计:")
	var statusStats []struct {
		Status models.OrderStatus
		Count  int64
	}
	db.Model(&models.Order{}).
		Select("status, COUNT(*) as count").
		Group("status").
		Order("status").
		Find(&statusStats)
	
	statusNames := map[models.OrderStatus]string{
		0: "待处理",
		1: "处理中",
		2: "已完成",
		3: "失败",
		4: "已退款",
		5: "已取消",
	}
	
	for _, stat := range statusStats {
		fmt.Printf("状态 %d (%s): %d 个订单\n", stat.Status, statusNames[stat.Status], stat.Count)
	}

	// 3. 检查最近30天的订单
	fmt.Println("\n最近30天的订单统计:")
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -30)
	
	var recentCount int64
	db.Model(&models.Order{}).
		Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Count(&recentCount)
	fmt.Printf("最近30天订单数: %d\n", recentCount)

	// 4. 按日统计最近7天的订单
	fmt.Println("\n最近7天每日订单数:")
	for i := 6; i >= 0; i-- {
		date := time.Now().AddDate(0, 0, -i)
		start := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
		end := start.Add(24 * time.Hour)
		
		var dayCount int64
		var successCount int64
		
		db.Model(&models.Order{}).
			Where("created_at >= ? AND created_at < ?", start, end).
			Count(&dayCount)
		
		db.Model(&models.Order{}).
			Where("created_at >= ? AND created_at < ?", start, end).
			Where("status = ?", models.OrderStatusSuccess).
			Count(&successCount)
		
		fmt.Printf("%s: 总计 %d 个订单, 成功 %d 个\n", start.Format("2006-01-02"), dayCount, successCount)
	}

	// 5. 查看最新的10个订单
	fmt.Println("\n最新的10个订单:")
	var orders []models.Order
	db.Order("created_at DESC").Limit(10).Find(&orders)
	
	for _, order := range orders {
		fmt.Printf("订单号: %s, 状态: %d, 金额: %.2f, 创建时间: %s\n",
			order.OrderNo, order.Status, order.TotalAmount, order.CreatedAt.Format("2006-01-02 15:04:05"))
	}
}
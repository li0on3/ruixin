package main

import (
	"backend/internal/config"
	"backend/internal/models"
	"fmt"
	"log"
	"math/rand"
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

	// 商品列表
	products := []struct {
		Name  string
		Price float64
		Cost  float64
	}{
		{"拿铁", 15.00, 8.00},
		{"美式咖啡", 12.00, 6.00},
		{"卡布奇诺", 16.00, 8.50},
		{"摩卡", 18.00, 10.00},
		{"焦糖玛奇朵", 20.00, 11.00},
		{"香草拿铁", 17.00, 9.00},
		{"抹茶拿铁", 19.00, 10.50},
		{"红茶拿铁", 16.00, 8.50},
	}

	// 门店列表
	stores := []struct {
		Code    string
		Name    string
		Address string
	}{
		{"SZ001", "深圳福田店", "深圳市福田区深南大道1234号"},
		{"SZ002", "深圳南山店", "深圳市南山区科技园路5678号"},
		{"SZ003", "深圳罗湖店", "深圳市罗湖区人民南路910号"},
		{"GZ001", "广州天河店", "广州市天河区天河路111号"},
		{"GZ002", "广州越秀店", "广州市越秀区北京路222号"},
	}

	// 创建最近30天的订单
	rand.Seed(time.Now().UnixNano())
	createdCount := 0

	for i := 29; i >= 0; i-- {
		date := time.Now().AddDate(0, 0, -i)
		// 每天创建2-8个订单
		ordersPerDay := rand.Intn(7) + 2

		for j := 0; j < ordersPerDay; j++ {
			// 随机选择商品
			product := products[rand.Intn(len(products))]
			
			// 随机选择门店
			store := stores[rand.Intn(len(stores))]
			
			// 创建订单商品信息
			quantity := rand.Intn(3) + 1 // 1-3杯
			goods := models.OrderGoods{
				{
					GoodsID:       fmt.Sprintf("G%03d", rand.Intn(100)+1),
					GoodsName:     product.Name,
					SKUCode:       fmt.Sprintf("SKU%03d", rand.Intn(100)+1),
					SKUName:       product.Name,
					Quantity:      quantity,
					OriginalPrice: product.Price,
					SalePrice:     product.Price,
				},
			}

			// 计算金额
			totalAmount := product.Price * float64(quantity)
			costAmount := product.Cost * float64(quantity)
			profitAmount := totalAmount - costAmount

			// 设置具体的创建时间
			createdAt := time.Date(date.Year(), date.Month(), date.Day(), 
				rand.Intn(14)+8, // 8-21点
				rand.Intn(60),   // 随机分钟
				rand.Intn(60),   // 随机秒
				0, date.Location())

			// 创建订单
			order := &models.Order{
				OrderNo:         fmt.Sprintf("TEST%d%02d%02d%04d", date.Year(), date.Month(), date.Day(), j+1),
				PutOrderID:      fmt.Sprintf("LUCKIN%d", time.Now().UnixNano()),
				DistributorID:   uint(rand.Intn(3) + 1), // 假设有3个分销商
				CardID:          uint(rand.Intn(10) + 1), // 假设有10张卡
				CardCode:        fmt.Sprintf("CARD%04d", rand.Intn(10)+1),
				Status:          models.OrderStatusSuccess, // 设置为成功状态
				StoreCode:       store.Code,
				StoreName:       store.Name,
				StoreAddress:    store.Address,
				Goods:           goods,
				TotalAmount:     totalAmount,
				CostAmount:      costAmount,
				ProfitAmount:    profitAmount,
				LuckinPrice:     product.Price,
				LuckinCostPrice: product.Cost,
			}

			// 创建订单并设置创建时间
			err := db.Create(order).Error
			if err != nil {
				log.Printf("创建订单失败: %v", err)
			} else {
				// 更新创建时间
				db.Model(order).Update("created_at", createdAt)
				createdCount++
			}
		}
	}

	fmt.Printf("成功创建 %d 个测试订单\n", createdCount)

	// 显示统计
	var stats []struct {
		Date  time.Time
		Count int64
		Total float64
	}
	
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -29)
	
	db.Model(&models.Order{}).
		Select("DATE(created_at) as date, COUNT(*) as count, SUM(total_amount) as total").
		Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Where("status = ?", models.OrderStatusSuccess).
		Group("DATE(created_at)").
		Order("date DESC").
		Limit(10).
		Find(&stats)
	
	fmt.Println("\n最近10天的订单统计:")
	for _, stat := range stats {
		fmt.Printf("%s: %d 个订单, 总金额: %.2f\n", 
			stat.Date.Format("2006-01-02"), stat.Count, stat.Total)
	}
}
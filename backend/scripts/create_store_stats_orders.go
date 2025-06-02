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

	// 门店列表 - 扩展更多门店
	stores := []struct {
		Code    string
		Name    string
		Address string
		City    string
		District string
	}{
		{"SZ001", "深圳福田店", "深圳市福田区深南大道1234号", "深圳", "福田区"},
		{"SZ002", "深圳南山店", "深圳市南山区科技园路5678号", "深圳", "南山区"},
		{"SZ003", "深圳罗湖店", "深圳市罗湖区人民南路910号", "深圳", "罗湖区"},
		{"SZ004", "深圳龙岗店", "深圳市龙岗区龙城广场111号", "深圳", "龙岗区"},
		{"SZ005", "深圳宝安店", "深圳市宝安区新安路222号", "深圳", "宝安区"},
		{"GZ001", "广州天河店", "广州市天河区天河路111号", "广州", "天河区"},
		{"GZ002", "广州越秀店", "广州市越秀区北京路222号", "广州", "越秀区"},
		{"GZ003", "广州海珠店", "广州市海珠区江南大道333号", "广州", "海珠区"},
		{"GZ004", "广州番禺店", "广州市番禺区市桥路444号", "广州", "番禺区"},
		{"BJ001", "北京朝阳店", "北京市朝阳区三里屯路555号", "北京", "朝阳区"},
		{"BJ002", "北京海淀店", "北京市海淀区中关村路666号", "北京", "海淀区"},
		{"SH001", "上海浦东店", "上海市浦东新区陆家嘴路777号", "上海", "浦东新区"},
		{"SH002", "上海静安店", "上海市静安区南京西路888号", "上海", "静安区"},
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
		{"澳白咖啡", 18.00, 9.50},
		{"生椰拿铁", 19.00, 10.00},
	}

	// 获取现有的分销商
	var distributors []models.Distributor
	db.Find(&distributors)
	if len(distributors) == 0 {
		log.Fatal("No distributors found, please create distributors first")
	}

	// 创建最近30天的订单，每个门店有不同的销售量
	rand.Seed(time.Now().UnixNano())
	createdCount := 0

	for i := 29; i >= 0; i-- {
		date := time.Now().AddDate(0, 0, -i)
		
		// 为每个门店创建订单
		for _, store := range stores {
			// 根据门店位置设置不同的订单量
			baseOrders := 3
			if store.City == "深圳" || store.City == "上海" {
				baseOrders = 5 // 一线城市订单更多
			} else if store.City == "广州" {
				baseOrders = 4
			}
			
			// 每天每个门店创建随机数量的订单
			ordersPerStore := rand.Intn(baseOrders) + baseOrders/2
			
			for j := 0; j < ordersPerStore; j++ {
				// 随机选择商品
				product := products[rand.Intn(len(products))]
				
				// 随机选择分销商
				distributor := distributors[rand.Intn(len(distributors))]
				
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
				hour := 8 + rand.Intn(14) // 8-21点
				createdAt := time.Date(date.Year(), date.Month(), date.Day(), 
					hour, rand.Intn(60), rand.Intn(60), 0, date.Location())

				// 创建订单
				order := &models.Order{
					OrderNo:         fmt.Sprintf("STORE%s%d%02d%02d%04d", store.Code, date.Year(), date.Month(), date.Day(), j+1),
					PutOrderID:      fmt.Sprintf("LUCKIN%d", time.Now().UnixNano()),
					DistributorID:   distributor.ID,
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
	}

	fmt.Printf("成功创建 %d 个测试订单\n", createdCount)

	// 显示每个门店的统计
	var storeStats []struct {
		StoreCode string
		StoreName string
		Count     int64
		Total     float64
	}
	
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -29)
	
	db.Model(&models.Order{}).
		Select("store_code, store_name, COUNT(*) as count, SUM(total_amount) as total").
		Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Where("status = ?", models.OrderStatusSuccess).
		Group("store_code, store_name").
		Order("total DESC").
		Find(&storeStats)

	fmt.Println("\n各门店订单统计:")
	for _, stat := range storeStats {
		fmt.Printf("%s (%s): %d 个订单, 总金额: %.2f\n", 
			stat.StoreName, stat.StoreCode, stat.Count, stat.Total)
	}
}
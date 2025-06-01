package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"backend/internal/config"
	"backend/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// 加载配置
	if err := config.Init("./configs/config.yaml"); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 连接数据库
	dsn := config.GetDSN()
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 创建测试卡片
	cards := []models.Card{
		{
			CardCode:    "TEST001",
			Status:      1,
			ProductID:   6,
			DailyLimit:  20,
			TotalLimit:  500,
			UsedCount:   45,
			Description: "测试卡片1",
			ExpiredAt:   time.Now().AddDate(0, 6, 0),
		},
		{
			CardCode:    "TEST002",
			Status:      1,
			ProductID:   10,
			DailyLimit:  10,
			TotalLimit:  200,
			UsedCount:   23,
			Description: "测试卡片2",
			ExpiredAt:   time.Now().AddDate(0, 3, 0),
		},
		{
			CardCode:    "8HVXKL76",
			Status:      1,
			ProductID:   6,
			DailyLimit:  10,
			TotalLimit:  100,
			UsedCount:   5,
			Description: "演示卡片",
			ExpiredAt:   time.Now().AddDate(1, 0, 0),
		},
	}

	for _, card := range cards {
		var existing models.Card
		if err := db.Where("card_code = ?", card.CardCode).First(&existing).Error; err != nil {
			db.Create(&card)
			fmt.Printf("Created card: %s\n", card.CardCode)
		}
	}

	// 创建测试分销商
	var distributor models.Distributor
	if err := db.Where("email = ?", "test@example.com").First(&distributor).Error; err != nil {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("test123"), bcrypt.DefaultCost)
		distributor = models.Distributor{
			Name:              "测试分销商",
			CompanyName:       "测试科技有限公司",
			ContactName:       "李测试",
			Phone:             "13900000001",
			Email:             "test@example.com",
			Password:          string(hashedPassword),
			APIKey:            "test-api-key-001",
			APISecret:         "test-api-secret-001",
			Status:            1,
			Balance:           5000.00,
			FrozenAmount:      0,
			CreditLimit:       20000.00,
			DailyOrderLimit:   200,
			MonthlyOrderLimit: 5000,
		}
		db.Create(&distributor)
		fmt.Println("Created test distributor")
	}

	// 创建测试订单（最近7天的数据）
	rand.Seed(time.Now().UnixNano())
	storeNames := []string{"瑞幸咖啡（世纪大道店）", "瑞幸咖啡（人民广场店）", "瑞幸咖啡（静安寺店）", "瑞幸咖啡（徐家汇店）"}
	goodsNames := []string{"拿铁", "美式咖啡", "燕麦拿铁", "橙C冰茶", "生椰拿铁"}
	
	for i := 0; i < 7; i++ {
		date := time.Now().AddDate(0, 0, -i)
		// 每天创建5-15个订单
		orderCount := rand.Intn(10) + 5
		
		for j := 0; j < orderCount; j++ {
			// 随机状态，大部分成功
			status := models.OrderStatusSuccess
			if rand.Float32() < 0.1 { // 10%失败率
				status = models.OrderStatusFailed
			}
			
			// 随机金额 15-50元
			amount := float64(rand.Intn(35)+15) + float64(rand.Intn(100))/100
			
			order := models.Order{
				OrderNo:          fmt.Sprintf("DD%d%02d%02d%04d", date.Year(), date.Month(), date.Day(), j+1),
				DistributorID:    distributor.ID,
				CardCode:         cards[rand.Intn(len(cards))].CardCode,
				StoreCode:        fmt.Sprintf("3877%02d", rand.Intn(20)+1),
				StoreName:        storeNames[rand.Intn(len(storeNames))],
				GoodsDetails:     fmt.Sprintf("[{\"name\":\"%s\",\"quantity\":1,\"price\":%.2f}]", goodsNames[rand.Intn(len(goodsNames))], amount),
				TotalAmount:      amount,
				CostAmount:       amount * 0.85, // 85%成本
				ProfitAmount:     amount * 0.15, // 15%利润
				PhoneNumber:      fmt.Sprintf("139%08d", rand.Intn(100000000)),
				Status:           status,
				TakeCode:         fmt.Sprintf("%04d", rand.Intn(10000)),
				QRData:           fmt.Sprintf("QR_%s", time.Now().Format("20060102150405")),
				CallbackTime:     &date,
				CallbackURL:      "https://callback.example.com",
				CreatedAt:        date,
				UpdatedAt:        date,
			}
			
			if status == models.OrderStatusSuccess {
				completedAt := date.Add(time.Minute * time.Duration(rand.Intn(30)+5))
				order.CompletedAt = &completedAt
			}
			
			db.Create(&order)
		}
		
		fmt.Printf("Created %d orders for %s\n", orderCount, date.Format("2006-01-02"))
	}

	// 创建财务记录
	// 为分销商创建一些充值记录
	transactions := []models.Transaction{
		{
			DistributorID:   distributor.ID,
			Type:            "recharge",
			Amount:          1000.00,
			BalanceBefore:   4000.00,
			BalanceAfter:    5000.00,
			Description:     "系统充值",
			RelatedOrderNo:  "",
			CreatedAt:       time.Now().AddDate(0, 0, -7),
		},
		{
			DistributorID:   distributor.ID,
			Type:            "recharge",
			Amount:          2000.00,
			BalanceBefore:   2000.00,
			BalanceAfter:    4000.00,
			Description:     "系统充值",
			RelatedOrderNo:  "",
			CreatedAt:       time.Now().AddDate(0, 0, -14),
		},
	}

	for _, tx := range transactions {
		db.Create(&tx)
	}

	fmt.Println("Test data created successfully!")
}
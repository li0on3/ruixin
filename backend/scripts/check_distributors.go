package main

import (
	"fmt"
	"log"

	"backend/internal/config"
	"backend/internal/models"
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

	// 查询所有分销商
	var distributors []models.Distributor
	if err := db.Find(&distributors).Error; err != nil {
		log.Fatalf("Failed to query distributors: %v", err)
	}

	fmt.Printf("找到 %d 个分销商:\n", len(distributors))
	for _, d := range distributors {
		fmt.Printf("ID: %d, 名称: %s, 邮箱: %s, 余额: %.2f, 状态: %d\n", 
			d.ID, d.Name, d.Email, d.Balance, d.Status)
	}

	if len(distributors) == 0 {
		fmt.Println("数据库中没有分销商数据，创建测试分销商...")
		
		// 创建测试分销商
		testDistributor := models.Distributor{
			Name:              "手动创建的测试分销商",
			CompanyName:       "手动测试公司",
			ContactName:       "测试联系人",
			Phone:             "13900000000",
			Email:             "manual-test@example.com",
			Password:          "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", // demo123
			APIKey:            "manual-test-api-key",
			APISecret:         "manual-test-api-secret",
			Status:            1,
			Balance:           1500.00,
			FrozenAmount:      0,
			CreditLimit:       15000.00,
			DailyOrderLimit:   150,
			MonthlyOrderLimit: 4000,
		}
		
		if err := db.Create(&testDistributor).Error; err != nil {
			log.Fatalf("Failed to create test distributor: %v", err)
		}
		
		fmt.Printf("成功创建测试分销商 ID: %d\n", testDistributor.ID)
	}
}
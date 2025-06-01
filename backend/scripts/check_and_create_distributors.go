package main

import (
	"fmt"
	"log"

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

	// 检查现有分销商数量
	var count int64
	db.Model(&models.Distributor{}).Count(&count)
	fmt.Printf("当前分销商数量: %d\n", count)

	// 列出现有分销商
	var distributors []models.Distributor
	db.Find(&distributors)
	fmt.Println("现有分销商列表:")
	for _, d := range distributors {
		fmt.Printf("ID: %d, Name: %s, Email: %s, Status: %d, Balance: %.2f\n", 
			d.ID, d.Name, d.Email, d.Status, d.Balance)
	}

	// 如果没有分销商，创建测试分销商
	if count == 0 {
		fmt.Println("没有分销商，创建测试分销商...")
		createTestDistributors(db)
	} else {
		fmt.Println("已有分销商存在")
	}
}

func createTestDistributors(db *gorm.DB) {
	// 创建默认密码
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("demo123"), bcrypt.DefaultCost)
	
	distributors := []models.Distributor{
		{
			Name:              "示例分销商A",
			CompanyName:       "示例公司A",
			ContactName:       "张三",
			Phone:             "13800138001",
			Email:             "demo1@distributor.com",
			Password:          string(hashedPassword),
			APIKey:            "demo-api-key-123456",
			APISecret:         "demo-api-secret-654321",
			Status:            1,
			Balance:           1000.00,
			FrozenAmount:      0,
			CreditLimit:       10000.00,
			DailyOrderLimit:   100,
			MonthlyOrderLimit: 3000,
		},
		{
			Name:              "示例分销商B",
			CompanyName:       "示例公司B",
			ContactName:       "李四",
			Phone:             "13800138002",
			Email:             "demo2@distributor.com",
			Password:          string(hashedPassword),
			APIKey:            "demo-api-key-223456",
			APISecret:         "demo-api-secret-754321",
			Status:            1,
			Balance:           2000.00,
			FrozenAmount:      0,
			CreditLimit:       15000.00,
			DailyOrderLimit:   150,
			MonthlyOrderLimit: 4500,
		},
		{
			Name:              "示例分销商C",
			CompanyName:       "示例公司C",
			ContactName:       "王五",
			Phone:             "13800138003",
			Email:             "demo3@distributor.com",
			Password:          string(hashedPassword),
			APIKey:            "demo-api-key-323456",
			APISecret:         "demo-api-secret-854321",
			Status:            1,
			Balance:           500.00,
			FrozenAmount:      0,
			CreditLimit:       5000.00,
			DailyOrderLimit:   50,
			MonthlyOrderLimit: 1500,
		},
	}

	for i, distributor := range distributors {
		if err := db.Create(&distributor).Error; err != nil {
			log.Printf("Failed to create distributor %d: %v", i+1, err)
		} else {
			fmt.Printf("创建分销商: %s (Email: %s)\n", distributor.Name, distributor.Email)
		}
	}
}
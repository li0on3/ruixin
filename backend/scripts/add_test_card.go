package main

import (
	"log"
	"time"

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

	// 添加测试卡片
	expiredAt := time.Now().AddDate(1, 0, 0)
	card := models.Card{
		CardCode:    "VTEX97M4",
		Status:      1,
		ProductID:   6,
		DailyLimit:  100,
		TotalLimit:  1000,
		UsedCount:   0,
		Description: "测试优惠卡",
		ExpiredAt:   &expiredAt,
	}

	if err := db.Create(&card).Error; err != nil {
		log.Printf("Failed to create card: %v", err)
	} else {
		log.Println("Test card VTEX97M4 created successfully")
	}
}
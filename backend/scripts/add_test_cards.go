package main

import (
	"fmt"
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

	// 获取第一个价格配置
	var price models.LuckinPrice
	if err := db.First(&price).Error; err != nil {
		log.Fatalf("No price found in database. Please create prices first.")
	}

	// 先创建一个批次
	batch := models.CardBatch{
		BatchNo:     fmt.Sprintf("TEST-%s", time.Now().Format("20060102150405")),
		PriceID:     price.ID,
		CostPrice:   5.00,
		SellPrice:   price.PriceValue,
		TotalCount:  5,
		UsedCount:   0,
		ImportedAt:  time.Now(),
		ImportedBy:  1, // 假设管理员ID为1
		Description: "测试批次",
	}
	if err := db.Create(&batch).Error; err != nil {
		log.Fatalf("Failed to create batch: %v", err)
	}
	fmt.Printf("Created batch: %s (ID: %d)\n\n", batch.BatchNo, batch.ID)

	// 创建测试卡片数据
	testCards := []models.Card{
		{
			CardCode:    "TEST001",
			BatchID:     batch.ID,
			PriceID:     price.ID,
			CostPrice:   5.00,
			SellPrice:   price.PriceValue,
			Status:      0, // 未使用
			ExpiredAt:   time.Now().AddDate(0, 6, 0), // 6个月后过期
			Description: "测试卡片001",
		},
		{
			CardCode:    "TEST002",
			BatchID:     batch.ID,
			PriceID:     price.ID,
			CostPrice:   5.00,
			SellPrice:   price.PriceValue,
			Status:      0, // 未使用
			ExpiredAt:   time.Now().AddDate(0, 6, 0),
			Description: "测试卡片002",
		},
		{
			CardCode:    "TEST003",
			BatchID:     batch.ID,
			PriceID:     price.ID,
			CostPrice:   5.00,
			SellPrice:   price.PriceValue,
			Status:      0, // 未使用
			ExpiredAt:   time.Now().AddDate(0, 6, 0),
			Description: "测试卡片003",
		},
		{
			CardCode:    "TEST004",
			BatchID:     batch.ID,
			PriceID:     price.ID,
			CostPrice:   5.00,
			SellPrice:   price.PriceValue,
			Status:      1, // 已使用
			ExpiredAt:   time.Now().AddDate(0, 6, 0),
			Description: "测试卡片004 - 已使用",
			UsedAt:      &[]time.Time{time.Now().Add(-24 * time.Hour)}[0],
		},
		{
			CardCode:    "TEST005",
			BatchID:     batch.ID,
			PriceID:     price.ID,
			CostPrice:   5.00,
			SellPrice:   price.PriceValue,
			Status:      2, // 预占中
			ExpiredAt:   time.Now().AddDate(0, 6, 0),
			Description: "测试卡片005 - 预占中",
			ReservedAt:  &[]time.Time{time.Now().Add(-1 * time.Hour)}[0],
		},
	}

	// 批量创建卡片
	successCount := 0
	for _, card := range testCards {
		// 先检查是否已存在
		var existing models.Card
		if err := db.Where("card_code = ?", card.CardCode).First(&existing).Error; err == nil {
			fmt.Printf("Card %s already exists, skipping...\n", card.CardCode)
			continue
		}

		if err := db.Create(&card).Error; err != nil {
			fmt.Printf("Failed to create card %s: %v\n", card.CardCode, err)
		} else {
			fmt.Printf("Created card %s (Status: %d)\n", card.CardCode, card.Status)
			successCount++
		}
	}

	fmt.Printf("\nTotal cards created: %d\n", successCount)

	// 显示统计信息
	var stats []struct {
		Status int
		Count  int64
	}
	db.Model(&models.Card{}).Select("status, count(*) as count").Group("status").Scan(&stats)
	
	fmt.Println("\nCard statistics:")
	for _, stat := range stats {
		statusName := "Unknown"
		switch stat.Status {
		case 0:
			statusName = "Unused"
		case 1:
			statusName = "Used"
		case 2:
			statusName = "Reserved"
		}
		fmt.Printf("Status %d (%s): %d cards\n", stat.Status, statusName, stat.Count)
	}
}
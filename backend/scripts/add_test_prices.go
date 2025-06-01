package main

import (
	"backend/internal/models"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"time"
)

func main() {
	// 连接数据库
	dsn := "dev:!QAZzse4@tcp(127.0.0.1:3306)/ruixin_platform?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// 创建测试价格数据
	prices := []models.LuckinPrice{
		{
			PriceCode:  "P15",
			PriceValue: 15.0,
			Status:     1,
			CreatedBy:  1,
		},
		{
			PriceCode:  "P20",
			PriceValue: 20.0,
			Status:     1,
			CreatedBy:  1,
		},
		{
			PriceCode:  "P25",
			PriceValue: 25.0,
			Status:     1,
			CreatedBy:  1,
		},
		{
			PriceCode:  "P30",
			PriceValue: 30.0,
			Status:     1,
			CreatedBy:  1,
		},
		{
			PriceCode:  "P35",
			PriceValue: 35.0,
			Status:     1,
			CreatedBy:  1,
		},
	}

	for _, price := range prices {
		var existing models.LuckinPrice
		if err := db.Where("price_id = ?", price.PriceCode).First(&existing).Error; err == gorm.ErrRecordNotFound {
			price.CreatedAt = time.Now()
			price.UpdatedAt = time.Now()
			if err := db.Create(&price).Error; err != nil {
				log.Printf("Failed to create price %s: %v", price.PriceCode, err)
			} else {
				fmt.Printf("Created price: %s - ¥%.2f\n", price.PriceCode, price.PriceValue)
			}
		} else {
			fmt.Printf("Price %s already exists\n", price.PriceCode)
		}
	}

	fmt.Println("Test prices added successfully!")
}
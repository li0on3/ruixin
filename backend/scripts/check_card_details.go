package main

import (
	"backend/internal/config"
	"backend/internal/models"
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// 加载配置
	cfg := &config.Config{
		Database: config.DatabaseConfig{
			Host:     "localhost",
			Port:     3306,
			Username: "test",
			Password: "!QAZzse4",
			Database: "ruixin_platform",
			Charset:  "utf8mb4",
		},
	}

	// 连接数据库
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		cfg.Database.Username,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Database,
		cfg.Database.Charset,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// 查询可用的卡片
	var cards []models.Card
	err = db.Where("status = 0").Limit(5).Find(&cards).Error
	if err != nil {
		log.Fatal("Failed to query cards:", err)
	}

	fmt.Println("可用卡片详情：")
	fmt.Println("=================")
	for _, card := range cards {
		fmt.Printf("卡片代码: %s\n", card.CardCode)
		fmt.Printf("ID: %d\n", card.ID)
		fmt.Printf("PriceID: %d\n", card.PriceID)
		fmt.Printf("成本价: %.2f\n", card.CostPrice)
		fmt.Printf("销售价: %.2f\n", card.SellPrice)
		fmt.Printf("状态: %d\n", card.Status)
		fmt.Println("-----------------")
	}

	// 查询商品4929绑定的卡片
	fmt.Println("\n商品4929绑定的卡片：")
	fmt.Println("===================")
	
	var product models.Product
	err = db.Where("goods_code = ?", "4929").First(&product).Error
	if err == nil {
		var bindings []models.CardProductBinding
		err = db.Where("product_id = ? AND is_active = ?", product.ID, true).
			Preload("Card").
			Find(&bindings).Error
		
		if err == nil {
			for _, binding := range bindings {
				if binding.Card != nil && binding.Card.Status == 0 {
					fmt.Printf("卡片代码: %s, PriceID: %d\n", binding.Card.CardCode, binding.Card.PriceID)
				}
			}
		}
	}

	// 查询luckin_prices表
	fmt.Println("\n价格配置表(luckin_prices)：")
	fmt.Println("==========================")
	type LuckinPrice struct {
		ID        int64   `gorm:"primarykey"`
		PriceID   int64   `gorm:"column:price_id"`
		PriceName string  `gorm:"column:price_name"`
		CostPrice float64 `gorm:"column:cost_price"`
		SellPrice float64 `gorm:"column:sell_price"`
	}
	
	var prices []LuckinPrice
	err = db.Table("luckin_prices").Limit(10).Find(&prices).Error
	if err == nil && len(prices) > 0 {
		for _, price := range prices {
			fmt.Printf("ID: %d, PriceID: %d, 名称: %s, 成本: %.2f, 销售: %.2f\n", 
				price.ID, price.PriceID, price.PriceName, price.CostPrice, price.SellPrice)
		}
	} else {
		fmt.Println("没有找到价格配置")
	}
}
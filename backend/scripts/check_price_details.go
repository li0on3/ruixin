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

	// 查询最近的订单
	fmt.Println("最近的订单信息：")
	fmt.Println("=================")
	var orders []models.Order
	err = db.Order("created_at DESC").Limit(5).Find(&orders).Error
	if err == nil {
		for _, order := range orders {
			fmt.Printf("订单号: %s\n", order.OrderNo)
			fmt.Printf("总金额: %.2f\n", order.TotalAmount)
			fmt.Printf("成本金额: %.2f\n", order.CostAmount)
			fmt.Printf("利润金额: %.2f\n", order.ProfitAmount)
			fmt.Printf("卡片代码: %s\n", order.CardCode)
			fmt.Printf("商品: %v\n", order.Goods)
			fmt.Printf("创建时间: %s\n", order.CreatedAt.Format("2006-01-02 15:04:05"))
			fmt.Println("-----------------")
		}
	}

	// 查询卡片的价格配置
	fmt.Println("\n卡片价格配置：")
	fmt.Println("===============")
	var cards []models.Card
	err = db.Where("card_code IN ?", []string{"KSGHGK4C", "8HVXKL76", "KK8NY7YK"}).Find(&cards).Error
	if err == nil {
		for _, card := range cards {
			fmt.Printf("卡片代码: %s\n", card.CardCode)
			fmt.Printf("成本价: %.2f\n", card.CostPrice)
			fmt.Printf("销售价: %.2f\n", card.SellPrice)
			fmt.Printf("PriceID: %d\n", card.PriceID)
			
			// 查询对应的价格配置
			var priceConfig struct {
				ID        int64   `gorm:"column:id"`
				PriceID   int64   `gorm:"column:price_id"`
				PriceName string  `gorm:"column:price_name"`
				CostPrice float64 `gorm:"column:cost_price"`
				SellPrice float64 `gorm:"column:sell_price"`
			}
			db.Table("luckin_prices").Where("id = ?", card.PriceID).First(&priceConfig)
			fmt.Printf("价格配置: ID=%d, PriceID=%d, 名称=%s\n", priceConfig.ID, priceConfig.PriceID, priceConfig.PriceName)
			fmt.Println("-----------------")
		}
	}

	// 查询商品4929的原价
	fmt.Println("\n商品4929的信息：")
	fmt.Println("=================")
	var product models.Product
	err = db.Where("goods_code = ?", "4929").First(&product).Error
	if err == nil {
		fmt.Printf("商品代码: %s\n", product.GoodsCode)
		fmt.Printf("商品名称: %s\n", product.GoodsName)
		
		// 查询SKU价格
		var skus []models.ProductSKU
		db.Where("product_id = ?", product.ID).Find(&skus)
		fmt.Println("SKU信息：")
		for _, sku := range skus {
			fmt.Printf("- SKU: %s, 名称: %s\n", sku.SKUCode, sku.SKUName)
		}
	}
}
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

	// 查询商品4929的信息
	var product models.Product
	err = db.Where("goods_code = ?", "4929").First(&product).Error
	if err != nil {
		fmt.Printf("商品代码 4929 不存在: %v\n", err)
		
		// 列出所有商品
		var products []models.Product
		db.Limit(10).Find(&products)
		fmt.Println("\n数据库中的商品列表（前10个）：")
		for _, p := range products {
			fmt.Printf("- %s: %s\n", p.GoodsCode, p.GoodsName)
		}
		return
	}

	fmt.Printf("找到商品: %s (%s)\n", product.GoodsCode, product.GoodsName)

	// 查询绑定的卡片
	var bindings []models.CardProductBinding
	err = db.Where("product_id = ? AND is_active = ?", product.ID, true).
		Preload("Card").
		Find(&bindings).Error
	
	if err != nil {
		fmt.Printf("查询绑定关系失败: %v\n", err)
		return
	}

	fmt.Printf("\n绑定的卡片数量: %d\n", len(bindings))
	
	// 统计可用卡片
	availableCount := 0
	for _, binding := range bindings {
		if binding.Card != nil && binding.Card.Status == 0 {
			availableCount++
			fmt.Printf("- 卡片 %s (ID: %d) - 状态: 可用\n", binding.Card.CardCode, binding.Card.ID)
		}
	}
	
	fmt.Printf("\n可用卡片数量: %d\n", availableCount)

	// 如果没有绑定，检查是否有通过价格映射的卡片
	if len(bindings) == 0 {
		fmt.Println("\n没有直接绑定的卡片，检查价格映射...")
		
		var mappings []models.ProductPriceMapping
		err = db.Where("product_code = ? AND status = 1", product.GoodsCode).Find(&mappings).Error
		if err == nil && len(mappings) > 0 {
			fmt.Printf("找到 %d 个价格映射\n", len(mappings))
			
			for _, mapping := range mappings {
				var cards []models.Card
				db.Where("price_id = ? AND status = 0", mapping.PriceID).Find(&cards)
				fmt.Printf("- 价格ID %d: %d 张可用卡片\n", mapping.PriceID, len(cards))
			}
		} else {
			fmt.Println("没有找到价格映射")
		}
	}
}
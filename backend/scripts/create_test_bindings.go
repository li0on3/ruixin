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

	fmt.Println("=== 创建测试绑定关系 ===")
	
	// 获取所有未使用的卡片
	var cards []models.Card
	db.Where("status = 0").Find(&cards)
	fmt.Printf("找到 %d 张未使用的卡片\n", len(cards))
	
	// 获取所有商品
	var products []models.Product
	db.Find(&products)
	fmt.Printf("找到 %d 个商品\n", len(products))
	
	if len(cards) == 0 {
		fmt.Println("错误：没有未使用的卡片，无法创建绑定关系")
		return
	}
	
	if len(products) == 0 {
		fmt.Println("错误：没有商品，无法创建绑定关系")
		return
	}
	
	// 为每个商品绑定第一张未使用的卡片
	createdCount := 0
	for _, product := range products {
		// 检查是否已存在绑定
		var existing models.CardProductBinding
		err := db.Where("card_id = ? AND product_id = ?", cards[0].ID, product.ID).First(&existing).Error
		if err == gorm.ErrRecordNotFound {
			// 创建绑定
			binding := models.CardProductBinding{
				CardID:    cards[0].ID,
				ProductID: product.ID,
				Priority:  0,
				IsActive:  true,
			}
			if err := db.Create(&binding).Error; err != nil {
				fmt.Printf("创建绑定失败: %v\n", err)
			} else {
				fmt.Printf("成功创建绑定: 卡片 %s -> 商品 %s\n", cards[0].CardCode, product.GoodsName)
				createdCount++
			}
		} else {
			fmt.Printf("绑定已存在: 卡片 %s -> 商品 %s\n", cards[0].CardCode, product.GoodsName)
		}
	}
	
	fmt.Printf("\n创建了 %d 个新的绑定关系\n", createdCount)
	
	// 验证可用商品查询
	var availableProducts []struct {
		GoodsName      string
		AvailableCount int
	}
	
	db.Raw(`
		SELECT p.goods_name, COUNT(DISTINCT c.id) as available_count
		FROM products p
		JOIN card_product_bindings cpb ON cpb.product_id = p.id
		JOIN cards c ON c.id = cpb.card_id
		WHERE c.status = 0 AND cpb.is_active = true AND c.deleted_at IS NULL
		GROUP BY p.id
	`).Scan(&availableProducts)
	
	fmt.Println("\n现在的可用商品:")
	for i, p := range availableProducts {
		fmt.Printf("%d. %s - 可用卡片数: %d\n", i+1, p.GoodsName, p.AvailableCount)
	}
}
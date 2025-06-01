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

	fmt.Println("=== 检查 card_product_bindings 表 ===")
	
	// 1. 统计总记录数
	var totalBindings int64
	db.Model(&models.CardProductBinding{}).Count(&totalBindings)
	fmt.Printf("\n1. 总绑定关系数: %d\n", totalBindings)
	
	// 2. 统计激活的绑定关系
	var activeBindings int64
	db.Model(&models.CardProductBinding{}).Where("is_active = ?", true).Count(&activeBindings)
	fmt.Printf("2. 激活的绑定关系数: %d\n", activeBindings)
	
	// 3. 检查绑定关系对应的卡片状态
	type CardStatus struct {
		Status int
		Count  int64
	}
	var cardStatuses []CardStatus
	db.Raw(`
		SELECT c.status, COUNT(DISTINCT c.id) as count
		FROM card_product_bindings cpb
		JOIN cards c ON c.id = cpb.card_id
		WHERE cpb.is_active = true AND c.deleted_at IS NULL
		GROUP BY c.status
	`).Scan(&cardStatuses)
	
	fmt.Println("\n3. 激活绑定关系对应的卡片状态分布:")
	for _, cs := range cardStatuses {
		statusName := "未知"
		switch cs.Status {
		case 0:
			statusName = "未使用"
		case 1:
			statusName = "已使用"
		case 2:
			statusName = "预占中"
		}
		fmt.Printf("   状态 %d (%s): %d 张卡片\n", cs.Status, statusName, cs.Count)
	}
	
	// 4. 查看可用商品查询的实际结果
	var availableProducts []struct {
		ProductID      uint
		GoodsCode      string
		GoodsName      string
		AvailableCount int
	}
	
	db.Raw(`
		SELECT p.id as product_id, p.goods_code, p.goods_name, COUNT(DISTINCT c.id) as available_count
		FROM products p
		JOIN card_product_bindings cpb ON cpb.product_id = p.id
		JOIN cards c ON c.id = cpb.card_id
		WHERE c.status = 0 AND cpb.is_active = true AND c.deleted_at IS NULL
		GROUP BY p.id
		ORDER BY available_count DESC
		LIMIT 10
	`).Scan(&availableProducts)
	
	fmt.Printf("\n4. 可用商品（前10个）:\n")
	if len(availableProducts) == 0 {
		fmt.Println("   没有找到可用商品！")
		fmt.Println("   原因：没有商品绑定了状态为0（未使用）的卡片")
	} else {
		for i, p := range availableProducts {
			fmt.Printf("   %d. %s (%s) - 可用卡片数: %d\n", i+1, p.GoodsName, p.GoodsCode, p.AvailableCount)
		}
	}
	
	// 5. 列出前5个绑定关系的详细信息
	fmt.Println("\n5. 前5个绑定关系详情:")
	var bindings []models.CardProductBinding
	db.Preload("Card").Preload("Product").Limit(5).Find(&bindings)
	
	for i, binding := range bindings {
		fmt.Printf("\n   绑定 %d:\n", i+1)
		fmt.Printf("   - 卡片: %s (状态: %d)\n", binding.Card.CardCode, binding.Card.Status)
		fmt.Printf("   - 商品: %s (%s)\n", binding.Product.GoodsName, binding.Product.GoodsCode)
		fmt.Printf("   - 激活状态: %v\n", binding.IsActive)
		fmt.Printf("   - 优先级: %d\n", binding.Priority)
	}
	
	// 6. 检查是否有商品但没有绑定关系
	var productsWithoutBindings int64
	db.Raw(`
		SELECT COUNT(DISTINCT p.id)
		FROM products p
		LEFT JOIN card_product_bindings cpb ON cpb.product_id = p.id
		WHERE cpb.id IS NULL
	`).Scan(&productsWithoutBindings)
	
	fmt.Printf("\n6. 没有绑定关系的商品数: %d\n", productsWithoutBindings)
}
package main

import (
	"fmt"
	"log"
	"time"

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
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 自动迁移 - 按照依赖关系顺序
	err = db.AutoMigrate(
		// 基础表（无外键依赖）
		&models.Admin{},
		&models.Distributor{},
		&models.LuckinPrice{},
		&models.LuckinProduct{},
		&models.City{},
		&models.SystemConfig{},   // 系统配置表
		
		// 依赖基础表的表
		&models.AdminOperationLog{},
		&models.CardBatch{},      // 依赖 LuckinPrice 和 Admin
		&models.Card{},           // 依赖 CardBatch 和 LuckinPrice
		&models.Order{},          // 依赖 Distributor 和 Card
		&models.CardUsageLog{},   // 依赖 Card 和 Distributor
		&models.DistributorAPILog{},
		&models.Transaction{},
		&models.Withdrawal{},
		&models.CategoryBinding{},
		
		// 商品相关表
		&models.Product{},
		&models.ProductSKU{},
		&models.ProductSpec{},
		&models.ProductSpecOption{},
		&models.ProductPriceMapping{},
		&models.ProductAlias{},      // 商品别名表
		&models.ProductMatchLog{},   // 匹配失败日志表
		&models.ProductSKUMapping{}, // SKU中文映射表
		&models.ProductSpecConfig{}, // 商品规格配置表
		&models.SpecInputAlias{},    // 规格输入别名表
		&models.CardProductBinding{}, // 卡片商品绑定表
	)
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	fmt.Println("Database migration completed")

	// 创建默认管理员账号
	var admin models.Admin
	if err := db.Where("username = ?", "admin").First(&admin).Error; err != nil {
		// 管理员不存在，创建新的
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		if err != nil {
			log.Fatalf("Failed to hash password: %v", err)
		}

		admin = models.Admin{
			Username:    "admin",
			Password:    string(hashedPassword),
			Email:       "admin@ruixin.com",
			RealName:    "系统管理员",
			Role:        "super_admin",
			Status:      1,
			LastLoginAt: time.Now(),
		}

		if err := db.Create(&admin).Error; err != nil {
			log.Fatalf("Failed to create admin user: %v", err)
		}

		fmt.Println("Default admin user created successfully")
		fmt.Println("Username: admin")
		fmt.Println("Password: admin123")
	} else {
		fmt.Println("Admin user already exists")
	}

	// 创建示例数据（可选）
	createSampleData(db)
}

func createSampleData(db *gorm.DB) {
	// 创建示例卡片批次和卡片
	var batch models.CardBatch
	if err := db.Where("batch_no = ?", "BATCH-2025-001").First(&batch).Error; err != nil {
		// 获取价格ID为6的价格
		var price models.LuckinPrice
		db.Where("price_id = ?", "6").First(&price)
		
		batch = models.CardBatch{
			BatchNo:     "BATCH-2025-001",
			PriceID:     price.ID,
			CostPrice:   7.50,
			SellPrice:   8.50,
			TotalCount:  5,
			UsedCount:   0,
			ImportedAt:  time.Now(),
			ImportedBy:  1,
			Description: "示例批次",
		}
		db.Create(&batch)
		fmt.Println("Sample card batch created")
		
		// 创建批次下的卡片
		expiredAt := time.Now().AddDate(1, 0, 0)
		cards := []models.Card{
			{CardCode: "DEMO123456", BatchID: &batch.ID, PriceID: price.ID, CostPrice: 7.50, SellPrice: 8.50, Status: 0, ExpiredAt: &expiredAt},
			{CardCode: "DEMO123457", BatchID: &batch.ID, PriceID: price.ID, CostPrice: 7.50, SellPrice: 8.50, Status: 0, ExpiredAt: &expiredAt},
			{CardCode: "DEMO123458", BatchID: &batch.ID, PriceID: price.ID, CostPrice: 7.50, SellPrice: 8.50, Status: 0, ExpiredAt: &expiredAt},
			{CardCode: "DEMO123459", BatchID: &batch.ID, PriceID: price.ID, CostPrice: 7.50, SellPrice: 8.50, Status: 0, ExpiredAt: &expiredAt},
			{CardCode: "DEMO123460", BatchID: &batch.ID, PriceID: price.ID, CostPrice: 7.50, SellPrice: 8.50, Status: 0, ExpiredAt: &expiredAt},
		}
		for _, card := range cards {
			db.Create(&card)
		}
		fmt.Println("Sample cards created")
	}

	// 创建示例分销商
	var distributor models.Distributor
	if err := db.Where("email = ?", "demo@distributor.com").First(&distributor).Error; err != nil {
		// 创建默认密码
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("demo123"), bcrypt.DefaultCost)
		
		distributor = models.Distributor{
			Name:              "示例分销商",
			CompanyName:       "示例公司",
			ContactName:       "张三",
			Phone:             "13800138000",
			Email:             "demo@distributor.com",
			Password:          string(hashedPassword),
			APIKey:            "demo-api-key-123456",
			APISecret:         "demo-api-secret-654321",
			Status:            1,
			Balance:           1000.00,
			FrozenAmount:      0,
			CreditLimit:       10000.00,
			DailyOrderLimit:   100,
			MonthlyOrderLimit: 3000,
		}
		db.Create(&distributor)
		fmt.Println("Sample distributor created")
		fmt.Println("Email: demo@distributor.com")
		fmt.Println("Password: demo123")
		fmt.Println("API Key: demo-api-key-123456")
		fmt.Println("API Secret: demo-api-secret-654321")
	}

	// 创建示例瑞幸价格
	prices := []models.LuckinPrice{
		{PriceCode: "6", PriceValue: 9.10, Status: 1, CreatedBy: 1},
		{PriceCode: "10", PriceValue: 15.50, Status: 1, CreatedBy: 1},
		{PriceCode: "15", PriceValue: 20.00, Status: 1, CreatedBy: 1},
		{PriceCode: "20", PriceValue: 25.00, Status: 1, CreatedBy: 1},
	}
	for _, price := range prices {
		var existing models.LuckinPrice
		if err := db.Where("price_id = ?", price.PriceCode).First(&existing).Error; err != nil {
			db.Create(&price)
			fmt.Printf("Sample price created: %s (¥%.2f)\n", price.PriceCode, price.PriceValue)
		}
	}

	// 创建示例瑞幸产品
	products := []models.LuckinProduct{
		{ProductID: "2500", Name: "标准美式", Description: "经典美式咖啡", Category: "美式家族", Status: 1, CreatedBy: 1},
		{ProductID: "4500", Name: "燕麦拿铁", Description: "燕麦奶拿铁咖啡", Category: "拿铁系列", Status: 1, CreatedBy: 1},
		{ProductID: "4805", Name: "拿铁", Description: "经典拿铁咖啡", Category: "拿铁系列", Status: 1, CreatedBy: 1},
		{ProductID: "4929", Name: "橙C冰茶", Description: "橙汁茶饮", Category: "茶饮系列", Status: 1, CreatedBy: 1},
	}
	for _, product := range products {
		var existing models.LuckinProduct
		if err := db.Where("product_id = ?", product.ProductID).First(&existing).Error; err != nil {
			db.Create(&product)
			fmt.Printf("Sample product created: %s (%s)\n", product.ProductID, product.Name)
		}
	}
	
	// 创建默认系统配置
	var config models.SystemConfig
	if err := db.Where("config_key = ?", "sync_store_code").First(&config).Error; err != nil {
		config = models.SystemConfig{
			ConfigKey:   "sync_store_code",
			ConfigValue: "390840", // 使用提供的店铺代码
			Description: "商品同步使用的店铺代码",
		}
		db.Create(&config)
		fmt.Println("System config created: sync_store_code = 390840")
	}
	
	// 创建默认规格别名
	createDefaultSpecAliases(db)
}

// createDefaultSpecAliases 创建默认的规格别名
func createDefaultSpecAliases(db *gorm.DB) {
	// 杯型别名
	sizeAliases := []models.SpecInputAlias{
		{SpecType: "size", StandardValue: "大杯", AliasValue: "大"},
		{SpecType: "size", StandardValue: "大杯", AliasValue: "16oz"},
		{SpecType: "size", StandardValue: "大杯", AliasValue: "大杯 16oz"},
		{SpecType: "size", StandardValue: "中杯", AliasValue: "中"},
		{SpecType: "size", StandardValue: "中杯", AliasValue: "12oz"},
		{SpecType: "size", StandardValue: "小杯", AliasValue: "小"},
		{SpecType: "size", StandardValue: "小杯", AliasValue: "8oz"},
	}
	
	// 温度别名
	tempAliases := []models.SpecInputAlias{
		{SpecType: "temperature", StandardValue: "冰", AliasValue: "冰饮"},
		{SpecType: "temperature", StandardValue: "冰", AliasValue: "加冰"},
		{SpecType: "temperature", StandardValue: "热", AliasValue: "热饮"},
		{SpecType: "temperature", StandardValue: "热", AliasValue: "加热"},
	}
	
	// 甜度别名
	sweetAliases := []models.SpecInputAlias{
		{SpecType: "sweetness", StandardValue: "标准甜", AliasValue: "标准"},
		{SpecType: "sweetness", StandardValue: "标准甜", AliasValue: "正常甜"},
		{SpecType: "sweetness", StandardValue: "微甜", AliasValue: "微微甜"},
		{SpecType: "sweetness", StandardValue: "微甜", AliasValue: "少糖"},
		{SpecType: "sweetness", StandardValue: "少甜", AliasValue: "半糖"},
		{SpecType: "sweetness", StandardValue: "不另外加糖", AliasValue: "无糖"},
		{SpecType: "sweetness", StandardValue: "不另外加糖", AliasValue: "不加糖"},
	}
	
	// 批量创建别名
	for _, alias := range sizeAliases {
		var existing models.SpecInputAlias
		if err := db.Where("spec_type = ? AND standard_value = ? AND alias_value = ?", 
			alias.SpecType, alias.StandardValue, alias.AliasValue).First(&existing).Error; err != nil {
			db.Create(&alias)
		}
	}
	
	for _, alias := range tempAliases {
		var existing models.SpecInputAlias
		if err := db.Where("spec_type = ? AND standard_value = ? AND alias_value = ?", 
			alias.SpecType, alias.StandardValue, alias.AliasValue).First(&existing).Error; err != nil {
			db.Create(&alias)
		}
	}
	
	for _, alias := range sweetAliases {
		var existing models.SpecInputAlias
		if err := db.Where("spec_type = ? AND standard_value = ? AND alias_value = ?", 
			alias.SpecType, alias.StandardValue, alias.AliasValue).First(&existing).Error; err != nil {
			db.Create(&alias)
		}
	}
	
	fmt.Println("Default spec aliases created")
}

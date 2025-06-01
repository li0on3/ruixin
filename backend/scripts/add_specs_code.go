package main

import (
	"log"
	"backend/internal/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// 加载配置
	if err := config.Init("configs/config.yaml"); err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// 使用配置中的DSN
	dsn := config.GetDSN()

	// 连接数据库
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// 检查字段是否已存在
	var columnExists bool
	err = db.Raw(`
		SELECT COUNT(*) > 0 
		FROM INFORMATION_SCHEMA.COLUMNS 
		WHERE TABLE_SCHEMA = ? 
		AND TABLE_NAME = 'product_sku_mappings' 
		AND COLUMN_NAME = 'specs_code'
	`, config.AppConfig.Database.Database).Scan(&columnExists).Error
	
	if err != nil {
		log.Fatal("Failed to check column existence:", err)
	}

	if columnExists {
		log.Println("Column specs_code already exists in product_sku_mappings table")
	} else {
		// 添加specs_code字段到product_sku_mappings表
		if err := db.Exec("ALTER TABLE product_sku_mappings ADD COLUMN specs_code VARCHAR(100) COMMENT '完整的规格代码组合'").Error; err != nil {
			log.Fatal("Failed to add column:", err)
		} else {
			log.Println("Successfully added specs_code column to product_sku_mappings table")
		}
	}

	log.Println("Migration completed")
}
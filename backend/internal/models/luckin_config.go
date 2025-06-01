package models

import (
	"gorm.io/gorm"
	"time"
)

// LuckinPrice 瑞幸面价表
type LuckinPrice struct {
	ID         int64     `gorm:"primaryKey" json:"id"`
	PriceCode  string    `gorm:"uniqueIndex;size:50;not null;column:price_id" json:"price_id"`  // 使用 column tag 保持数据库字段名不变
	PriceValue float64   `gorm:"not null" json:"price_value"`
	Status     int       `gorm:"default:1" json:"status"` // 0:禁用 1:启用
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	CreatedBy  int64     `gorm:"not null" json:"created_by"`

	// Associations
	Creator Admin `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
}

func (LuckinPrice) TableName() string {
	return "luckin_prices"
}

// LuckinProduct 瑞幸产品表
type LuckinProduct struct {
	ID          int64     `gorm:"primaryKey" json:"id"`
	ProductID   string    `gorm:"uniqueIndex;size:50;not null" json:"product_id"`
	Name        string    `gorm:"size:200;not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	Category    string    `gorm:"size:100" json:"category"`
	ImageURL    string    `gorm:"size:500" json:"image_url"`
	Status      int       `gorm:"default:1" json:"status"` // 0:禁用 1:启用
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	CreatedBy   int64     `gorm:"not null" json:"created_by"`

	// Associations
	Creator Admin `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
}

func (LuckinProduct) TableName() string {
	return "luckin_products"
}

// CategoryBinding 种类绑定表
type CategoryBinding struct {
	ID         int64     `gorm:"primaryKey" json:"id"`
	CategoryID int64     `gorm:"not null;index" json:"category_id"`
	TargetType string    `gorm:"size:20;not null" json:"target_type"` // price, product
	TargetID   string    `gorm:"size:50;not null" json:"target_id"`
	Priority   int       `gorm:"default:0" json:"priority"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
	CreatedBy  int64     `gorm:"not null" json:"created_by"`

	// Associations
	Category Card  `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Creator  Admin `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
}

func (CategoryBinding) TableName() string {
	return "category_bindings"
}

// 添加唯一索引
func (c *CategoryBinding) BeforeCreate(tx *gorm.DB) error {
	// 确保同一个种类下，同一类型的目标ID只能绑定一次
	var count int64
	tx.Model(&CategoryBinding{}).
		Where("category_id = ? AND target_type = ? AND target_id = ?",
			c.CategoryID, c.TargetType, c.TargetID).
		Count(&count)
	if count > 0 {
		return gorm.ErrDuplicatedKey
	}
	return nil
}

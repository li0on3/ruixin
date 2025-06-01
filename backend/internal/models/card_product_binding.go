package models

import (
	"time"
	"gorm.io/gorm"
)

// CardProductBinding 卡片商品绑定关系
type CardProductBinding struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	CardID      uint   `gorm:"index;not null" json:"card_id"`      // 卡片ID
	ProductID   uint   `gorm:"index;not null" json:"product_id"`   // 商品ID
	Priority    int    `gorm:"default:0" json:"priority"`          // 优先级
	IsActive    bool   `gorm:"default:true" json:"is_active"`      // 是否启用
	
	// 关联
	Card    *Card    `gorm:"foreignKey:CardID" json:"card,omitempty"`
	Product *Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
}

func (CardProductBinding) TableName() string {
	return "card_product_bindings"
}
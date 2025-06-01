package models

import (
	"gorm.io/gorm"
	"time"
)

// ProductAlias 商品别名表
type ProductAlias struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	ProductID uint   `gorm:"index;not null" json:"product_id"`           // 商品ID
	AliasName string `gorm:"type:varchar(200);index;not null" json:"alias_name"` // 别名

	// 关联
	Product *Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
}

func (ProductAlias) TableName() string {
	return "product_aliases"
}
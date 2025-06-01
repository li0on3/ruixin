package models

import (
	"gorm.io/gorm"
	"time"
)

// ProductSKUMapping 商品SKU中文映射表
type ProductSKUMapping struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	ProductID   uint   `gorm:"index;not null" json:"product_id"`                      // 商品ID
	SKUCode     string `gorm:"type:varchar(50);not null" json:"sku_code"`            // SKU代码
	ChineseDesc string `gorm:"type:varchar(200);index;not null" json:"chinese_desc"` // 中文描述，如"大杯冰微甜茉莉花香"
	SpecsCode   string `gorm:"type:varchar(100)" json:"specs_code"`                  // 完整的规格代码组合，如"0_0_0_1"

	// 关联
	Product *Product    `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	SKU     *ProductSKU `gorm:"foreignKey:SKUCode;references:SKUCode" json:"sku,omitempty"`
}

// ProductSpecConfig 商品规格配置表
type ProductSpecConfig struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	ProductID    uint   `gorm:"index;not null" json:"product_id"`                     // 商品ID
	SpecType     string `gorm:"type:varchar(50);not null" json:"spec_type"`          // 规格类型：size, temperature, sweetness, flavor
	ChineseName  string `gorm:"type:varchar(50);not null" json:"chinese_name"`       // 中文名称：杯型, 温度, 甜度, 口味
	IsRequired   bool   `gorm:"default:true" json:"is_required"`                      // 是否必填
	DefaultValue string `gorm:"type:varchar(50)" json:"default_value"`               // 默认值
	DisplayOrder int    `gorm:"default:0" json:"display_order"`                       // 显示顺序

	// 关联
	Product *Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
}

// SpecInputAlias 规格输入别名表
type SpecInputAlias struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	SpecType       string `gorm:"type:varchar(50);index;not null" json:"spec_type"`       // 规格类型
	StandardValue  string `gorm:"type:varchar(50);not null" json:"standard_value"`        // 标准值，如"大杯"
	AliasValue     string `gorm:"type:varchar(50);index;not null" json:"alias_value"`     // 别名，如"大"、"16oz"
}

func (ProductSKUMapping) TableName() string {
	return "product_sku_mappings"
}

func (ProductSpecConfig) TableName() string {
	return "product_spec_configs"
}

func (SpecInputAlias) TableName() string {
	return "spec_input_aliases"
}
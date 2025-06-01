package models

import (
	"gorm.io/gorm"
	"time"
)

// Product 商品信息表 - 仅保存下单必需的信息
type Product struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	GoodsCode      string    `gorm:"type:varchar(50);uniqueIndex;not null" json:"goods_code"`  // 商品代码
	GoodsName      string    `gorm:"type:varchar(200);not null" json:"goods_name"`             // 商品名称（仅用于显示）
	AvailableSpecs string    `gorm:"type:text" json:"available_specs"`                         // 可用规格信息JSON
	LastSyncAt     time.Time `json:"last_sync_at"`                                             // 最后同步时间
	
	// 关联
	SKUs    []ProductSKU    `gorm:"foreignKey:ProductID" json:"skus,omitempty"`
	Aliases []ProductAlias  `gorm:"foreignKey:ProductID" json:"aliases,omitempty"`
}

// ProductSKU 商品SKU信息 - 仅保存下单必需的信息
type ProductSKU struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	ProductID   uint   `gorm:"index" json:"product_id"`                                  // 商品ID
	SKUCode     string `gorm:"type:varchar(50);index;not null" json:"sku_code"`         // SKU代码
	SKUName     string `gorm:"type:varchar(200)" json:"sku_name"`                       // SKU名称（仅用于显示）
	
	// 关联
	Product *Product      `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	Specs   []ProductSpec `gorm:"foreignKey:SKUID" json:"specs,omitempty"`
	Mapping *ProductSKUMapping `gorm:"foreignKey:SKUCode;references:SKUCode" json:"mapping,omitempty"`
}

// ProductSpec 商品规格信息
type ProductSpec struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	SKUID       uint   `gorm:"index" json:"sku_id"`                                      // SKU ID
	SpecsCode   string `gorm:"type:varchar(50);not null" json:"specs_code"`             // 规格代码
	SpecsName   string `gorm:"type:varchar(100)" json:"specs_name"`                     // 规格名称
	SpecsType   int    `json:"specs_type"`                                               // 规格类型
	IsRequired  bool   `gorm:"default:false" json:"is_required"`                        // 是否必选
	
	// 关联
	SKU     *ProductSKU          `gorm:"foreignKey:SKUID" json:"sku,omitempty"`
	Options []ProductSpecOption  `gorm:"foreignKey:SpecID" json:"options,omitempty"`
}

// ProductSpecOption 规格选项
type ProductSpecOption struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	SpecID      uint   `gorm:"index" json:"spec_id"`                                     // 规格ID
	Code        string `gorm:"type:varchar(50);not null" json:"code"`                   // 选项代码
	Name        string `gorm:"type:varchar(100)" json:"name"`                           // 选项名称
	IsDefault   bool   `gorm:"default:false" json:"is_default"`                         // 是否默认选项
	
	// 关联
	Spec *ProductSpec `gorm:"foreignKey:SpecID" json:"spec,omitempty"`
}

// ProductPriceMapping 商品价格映射表 - 记录哪些商品可以使用哪个价格
type ProductPriceMapping struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	PriceID     int64  `gorm:"index" json:"price_id"`                                    // 价格ID
	ProductCode string `gorm:"type:varchar(50);index" json:"product_code"`              // 商品代码
	SKUCode     string `gorm:"type:varchar(50)" json:"sku_code"`                        // SKU代码（可选，为空表示所有SKU）
	Priority    int    `gorm:"default:0" json:"priority"`                                // 优先级（数字越大优先级越高）
	Status      int    `gorm:"default:1" json:"status"`                                  // 状态
	
	// 关联
	Price *LuckinPrice `gorm:"foreignKey:PriceID;references:ID" json:"price,omitempty"`
}

func (Product) TableName() string {
	return "products"
}

func (ProductSKU) TableName() string {
	return "product_skus"
}

func (ProductSpec) TableName() string {
	return "product_specs"
}

func (ProductSpecOption) TableName() string {
	return "product_spec_options"
}

func (ProductPriceMapping) TableName() string {
	return "product_price_mappings"
}
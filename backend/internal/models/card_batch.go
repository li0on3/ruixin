package models

import (
	"gorm.io/gorm"
	"time"
)

// CardBatch 卡片批次 - 用于管理批量导入的卡片
type CardBatch struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	BatchNo         string    `gorm:"type:varchar(50);uniqueIndex;not null" json:"batch_no"`     // 批次号
	LuckinProductID int       `gorm:"default:6" json:"luckin_product_id"`                        // 瑞幸产品ID，默认6
	PriceID         int64     `gorm:"type:bigint;index" json:"price_id,omitempty"`               // 关联价格ID（已废弃，保留用于兼容）
	CostPrice       float64   `gorm:"type:decimal(10,2)" json:"cost_price"`                     // 成本价
	SellPrice       float64   `gorm:"type:decimal(10,2)" json:"sell_price"`                     // 销售价
	TotalCount  int       `json:"total_count"`                                               // 总数量
	UsedCount   int       `json:"used_count"`                                                // 已使用数量
	ImportedAt  time.Time `json:"imported_at"`                                               // 导入时间
	ImportedBy  uint      `json:"imported_by"`                                               // 导入人ID
	Description string    `gorm:"type:text" json:"description"`                              // 备注

	// 关联
	Price    *LuckinPrice `gorm:"foreignKey:PriceID;references:ID" json:"price,omitempty"`
	Admin    *Admin       `gorm:"foreignKey:ImportedBy" json:"admin,omitempty"`
	Cards    []Card       `gorm:"foreignKey:BatchID" json:"cards,omitempty"`
}

func (CardBatch) TableName() string {
	return "card_batches"
}

// 计算批次的使用率
func (b *CardBatch) GetUsageRate() float64 {
	if b.TotalCount == 0 {
		return 0
	}
	return float64(b.UsedCount) / float64(b.TotalCount) * 100
}

// 计算平台利润率
func (b *CardBatch) GetPlatformProfitRate() float64 {
	if b.CostPrice == 0 {
		return 0
	}
	return (b.SellPrice - b.CostPrice) / b.CostPrice * 100
}

// 计算分销商利润率（需要商品价格）
func (b *CardBatch) GetDistributorProfitRate(productPrice float64) float64 {
	if b.SellPrice == 0 {
		return 0
	}
	return (productPrice - b.SellPrice) / b.SellPrice * 100
}
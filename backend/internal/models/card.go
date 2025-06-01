package models

import (
	"errors"
	"gorm.io/gorm"
	"time"
)

// Card status constants
const (
	CardStatusNormal   = 0 // 未使用
	CardStatusUsed     = 1 // 已使用
	CardStatusReserved = 2 // 预占中
)

var (
	// ErrInvalidPrice is returned when cost price is greater than sell price
	ErrInvalidPrice = errors.New("invalid price: cost price cannot be greater than sell price")
	// ErrCardExpired is returned when card is expired
	ErrCardExpired = errors.New("card is expired")
	// ErrInvalidCardCode is returned when card code format is invalid
	ErrInvalidCardCode = errors.New("invalid card code format")
	// ErrInvalidPhoneNumber is returned when phone number format is invalid
	ErrInvalidPhoneNumber = errors.New("invalid phone number format")
	// ErrInvalidOrderNo is returned when order number format is invalid
	ErrInvalidOrderNo = errors.New("invalid order number format")
)

type Card struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	CardCode         string     `gorm:"type:varchar(100);not null" json:"card_code"`
	LuckinProductID  int        `gorm:"default:6" json:"luckin_product_id"`           // 瑞幸产品ID，默认6
	BatchID          *uint      `json:"batch_id"`                                      // 批次ID，单个创建的卡片可为空
	PriceID          int64      `gorm:"type:bigint" json:"price_id,omitempty"`        // 价格ID（已废弃，保留用于兼容）
	CostPrice        float64    `gorm:"type:decimal(10,2)" json:"cost_price"`        // 成本价
	SellPrice        float64    `gorm:"type:decimal(10,2)" json:"sell_price"`        // 销售价
	Status      int        `gorm:"default:0" json:"status"`                      // 0: 未使用, 1: 已使用, 2: 预占中
	SyncStatus  string     `gorm:"type:varchar(20);default:'pending'" json:"sync_status"` // pending, syncing, synced, failed
	SyncedAt    *time.Time `json:"synced_at"`                                     // 最后同步时间
	UsedAt      *time.Time `json:"used_at"`                                       // 使用时间
	ReservedAt  *time.Time `json:"reserved_at"`                                   // 预占时间
	OrderID     *uint      `json:"order_id"`                                      // 订单ID
	ExpiredAt   *time.Time `json:"expired_at"`                                   // 过期时间
	Description string     `gorm:"type:text" json:"description"`                  // 描述

	// 关联
	Batch *CardBatch   `gorm:"foreignKey:BatchID;references:ID" json:"batch,omitempty"`
	Price *LuckinPrice `gorm:"foreignKey:PriceID;references:ID" json:"price,omitempty"`
	Order *Order       `gorm:"foreignKey:OrderID;references:ID" json:"order,omitempty"`
}

// CardUsageLog 卡片使用日志 - 记录所有使用尝试，用于审计和分析
type CardUsageLog struct {
	ID            uint      `gorm:"primarykey" json:"id"`
	CreatedAt     time.Time `json:"created_at"`
	CardID        uint      `json:"card_id"`
	CardCode      string    `gorm:"type:varchar(100)" json:"card_code"`      // 冗余存储，方便查询
	DistributorID uint      `json:"distributor_id"`
	OrderNo       string    `gorm:"type:varchar(50)" json:"order_no"`        // 订单号
	Success       bool      `json:"success"`                                  // 是否成功
	FailReason    string    `gorm:"type:varchar(50)" json:"fail_reason"`     // 失败原因类型
	ErrorMessage  string    `gorm:"type:text" json:"error_message,omitempty"` // 详细错误信息
	RequestIP     string    `gorm:"type:varchar(45)" json:"request_ip"`      // 请求IP
	UserAgent     string    `gorm:"type:varchar(255)" json:"user_agent"`     // 用户代理
	
	// 关联
	Card        *Card        `gorm:"foreignKey:CardID" json:"card,omitempty"`
	Distributor *Distributor `gorm:"foreignKey:DistributorID" json:"distributor,omitempty"`
}

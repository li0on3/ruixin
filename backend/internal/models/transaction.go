package models

import (
	"time"
)

type TransactionType int

const (
	TransactionTypeRecharge TransactionType = 1 // 充值
	TransactionTypeConsume  TransactionType = 2 // 消费
	TransactionTypeWithdraw TransactionType = 3 // 提现
	TransactionTypeRefund   TransactionType = 4 // 退款
	TransactionTypeAdjust   TransactionType = 5 // 调整
)

type Transaction struct {
	ID            int64           `gorm:"primaryKey" json:"id"`
	DistributorID int64           `gorm:"not null;index" json:"distributor_id"`
	Type          TransactionType `gorm:"not null;index" json:"type"`
	Amount        float64         `gorm:"not null" json:"amount"`
	BalanceBefore float64         `gorm:"not null" json:"balance_before"`
	BalanceAfter  float64         `gorm:"not null" json:"balance_after"`
	RelatedID     string          `gorm:"size:100" json:"related_id"`
	Remark        string          `gorm:"size:500" json:"remark"`
	CreatedBy     int64           `gorm:"not null" json:"created_by"`
	CreatedAt     time.Time       `gorm:"autoCreateTime" json:"created_at"`

	// Associations
	Distributor Distributor `gorm:"foreignKey:DistributorID" json:"distributor,omitempty"`
	Creator     Admin       `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
}

func (Transaction) TableName() string {
	return "transactions"
}

func (t TransactionType) String() string {
	switch t {
	case TransactionTypeRecharge:
		return "充值"
	case TransactionTypeConsume:
		return "消费"
	case TransactionTypeWithdraw:
		return "提现"
	case TransactionTypeRefund:
		return "退款"
	case TransactionTypeAdjust:
		return "调整"
	default:
		return "未知"
	}
}
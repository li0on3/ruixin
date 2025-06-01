package models

import (
	"time"
)

type WithdrawalStatus int

const (
	WithdrawalStatusPending   WithdrawalStatus = 0 // 待处理
	WithdrawalStatusProcessed WithdrawalStatus = 1 // 已处理
	WithdrawalStatusRejected  WithdrawalStatus = 2 // 已拒绝
)

type Withdrawal struct {
	ID            int64            `gorm:"primaryKey" json:"id"`
	DistributorID int64            `gorm:"not null;index" json:"distributor_id"`
	Amount        float64          `gorm:"not null" json:"amount"`
	Status        WithdrawalStatus `gorm:"default:0;index" json:"status"`
	AccountInfo   string           `gorm:"type:text;not null" json:"account_info"`
	Remark        string           `gorm:"size:500" json:"remark"`
	CreatedAt     time.Time        `gorm:"autoCreateTime" json:"created_at"`
	ProcessedAt   *time.Time       `json:"processed_at"`
	ProcessedBy   *int64           `json:"processed_by"`
	RejectReason  string           `gorm:"size:500" json:"reject_reason"`

	// Associations
	Distributor Distributor `gorm:"foreignKey:DistributorID" json:"distributor,omitempty"`
	Processor   *Admin      `gorm:"foreignKey:ProcessedBy" json:"processor,omitempty"`
}

func (Withdrawal) TableName() string {
	return "withdrawals"
}

func (s WithdrawalStatus) String() string {
	switch s {
	case WithdrawalStatusPending:
		return "待处理"
	case WithdrawalStatusProcessed:
		return "已处理"
	case WithdrawalStatusRejected:
		return "已拒绝"
	default:
		return "未知"
	}
}

type AccountInfo struct {
	Type        string `json:"type"`         // 账户类型：alipay/wechat/bank
	AccountName string `json:"account_name"` // 账户名称
	AccountNo   string `json:"account_no"`   // 账号
	BankName    string `json:"bank_name"`    // 银行名称（银行卡时需要）
	BankBranch  string `json:"bank_branch"`  // 开户行（银行卡时需要）
}
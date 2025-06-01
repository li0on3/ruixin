package models

import (
	"gorm.io/gorm"
	"time"
)

type Distributor struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Name           string  `gorm:"type:varchar(100);not null" json:"name"`
	CompanyName    string  `gorm:"type:varchar(200)" json:"company_name"`
	ContactName    string  `gorm:"type:varchar(100)" json:"contact_name"`
	Phone          string  `gorm:"type:varchar(20)" json:"phone"`
	Email          string  `gorm:"type:varchar(100);uniqueIndex" json:"email"`
	Password       string  `gorm:"type:varchar(255)" json:"-"`
	APIKey         string  `gorm:"type:varchar(64);uniqueIndex;not null" json:"-"`
	APISecret      string  `gorm:"type:varchar(64)" json:"-"`
	Status         int     `json:"status"` // 0: 待审核, 1: 正常, 2: 禁用
	Balance        float64 `json:"balance"`
	FrozenAmount   float64 `json:"frozen_amount"` // 冻结金额
	CreditLimit    float64 `json:"credit_limit"`
	CallbackURL    string  `gorm:"type:varchar(500)" json:"callback_url"`
	WarningBalance float64 `json:"warning_balance"`                          // 余额预警阈值
	WarningEnabled bool    `json:"warning_enabled"`                          // 是否启用预警
	WarningEmail   string  `gorm:"type:varchar(100)" json:"warning_email"`   // 预警邮箱
	WarningWebhook string  `gorm:"type:varchar(500)" json:"warning_webhook"` // 预警Webhook

	// 配额限制
	DailyOrderLimit   int `json:"daily_order_limit"`
	MonthlyOrderLimit int `json:"monthly_order_limit"`

	// 统计信息
	TotalOrders int     `json:"total_orders"`
	TotalAmount float64 `json:"total_amount"`
}

type DistributorAPILog struct {
	ID            uint      `gorm:"primarykey" json:"id"`
	CreatedAt     time.Time `json:"created_at"`
	DistributorID uint      `json:"distributor_id"`
	APIEndpoint   string    `json:"api_endpoint"`
	Method        string    `json:"method"`
	RequestBody   string    `json:"request_body"`
	ResponseCode  int       `json:"response_code"`
	ResponseBody  string    `json:"response_body"`
	IPAddress     string    `json:"ip_address"`
	UserAgent     string    `json:"user_agent"`
}

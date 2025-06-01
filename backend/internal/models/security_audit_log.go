package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// SecurityAuditLog 安全审计日志
type SecurityAuditLog struct {
	ID            uint            `gorm:"primaryKey" json:"id"`
	DistributorID uint            `gorm:"index:idx_distributor_action" json:"distributor_id"`
	Action        string          `gorm:"size:50;index:idx_distributor_action" json:"action"`
	Resource      string          `gorm:"size:100" json:"resource"`
	Details       JSONMap         `gorm:"type:json" json:"details"`
	IPAddress     string          `gorm:"size:45" json:"ip_address"`
	UserAgent     string          `gorm:"type:text" json:"user_agent"`
	Status        string          `gorm:"size:20" json:"status"`
	ErrorMsg      string          `gorm:"type:text" json:"error_msg"`
	CreatedAt     time.Time       `gorm:"index:idx_created_at" json:"created_at"`
}

// JSONMap 用于存储JSON数据
type JSONMap map[string]interface{}

// Value 实现 driver.Valuer 接口
func (j JSONMap) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// Scan 实现 sql.Scanner 接口
func (j *JSONMap) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, j)
}

// 审计日志操作类型常量
const (
	AuditActionCardAccess       = "CARD_ACCESS"           // 访问卡片
	AuditActionUnauthorizedCard = "UNAUTHORIZED_CARD"     // 未授权卡片访问
	AuditActionRateLimited      = "RATE_LIMITED"          // 触发频率限制
	AuditActionOrderCreate      = "ORDER_CREATE"          // 创建订单
	AuditActionWithdrawal       = "WITHDRAWAL"            // 提现操作
	AuditActionAPIError         = "API_ERROR"             // API错误
)

// 审计日志状态常量
const (
	AuditStatusSuccess = "success"
	AuditStatusWarning = "warning"
	AuditStatusFailed  = "failed"
)
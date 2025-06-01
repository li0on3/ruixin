package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// ProductMatchLog 商品匹配失败日志表
type ProductMatchLog struct {
	ID            uint      `gorm:"primarykey" json:"id"`
	CreatedAt     time.Time `json:"created_at"`
	
	DistributorID uint      `gorm:"index;not null" json:"distributor_id"`             // 分销商ID
	RequestTime   time.Time `gorm:"index;not null" json:"request_time"`               // 请求时间
	InputProduct  string    `gorm:"type:varchar(200)" json:"input_product"`           // 输入的商品名
	InputSpecs    JSON      `gorm:"type:json" json:"input_specs"`                     // 输入的规格
	ErrorReason   string    `gorm:"type:varchar(500)" json:"error_reason"`            // 错误原因
	Suggestions   JSON      `gorm:"type:json" json:"suggestions"`                     // 建议选项
	
	// 关联
	Distributor *Distributor `gorm:"foreignKey:DistributorID" json:"distributor,omitempty"`
}

// JSON 自定义JSON类型，用于存储JSON数据
type JSON json.RawMessage

// Scan 从数据库读取
func (j *JSON) Scan(value interface{}) error {
	if value == nil {
		*j = JSON("null")
		return nil
	}
	s, ok := value.([]byte)
	if !ok {
		return nil
	}
	*j = JSON(s)
	return nil
}

// Value 写入数据库
func (j JSON) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	return json.RawMessage(j), nil
}

// MarshalJSON 序列化
func (j JSON) MarshalJSON() ([]byte, error) {
	if j == nil {
		return []byte("null"), nil
	}
	return json.RawMessage(j), nil
}

// UnmarshalJSON 反序列化
func (j *JSON) UnmarshalJSON(data []byte) error {
	if j == nil {
		return nil
	}
	*j = JSON(data)
	return nil
}

func (ProductMatchLog) TableName() string {
	return "product_match_logs"
}
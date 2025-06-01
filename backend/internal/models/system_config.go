package models

import (
	"time"
)

// SystemConfig 系统配置表
type SystemConfig struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	
	ConfigKey   string `gorm:"type:varchar(100);uniqueIndex;not null" json:"config_key"`  // 配置键
	ConfigValue string `gorm:"type:text" json:"config_value"`                            // 配置值
	Description string `gorm:"type:varchar(500)" json:"description"`                     // 配置说明
}

func (SystemConfig) TableName() string {
	return "system_configs"
}
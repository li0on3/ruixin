package repository

import (
	"backend/internal/models"
	"gorm.io/gorm"
)

type SystemConfigRepository struct {
	db *gorm.DB
}

func NewSystemConfigRepository(db *gorm.DB) *SystemConfigRepository {
	return &SystemConfigRepository{db: db}
}

// GetByKey 根据配置键获取配置
func (r *SystemConfigRepository) GetByKey(key string) (*models.SystemConfig, error) {
	var config models.SystemConfig
	err := r.db.Where("config_key = ?", key).First(&config).Error
	return &config, err
}

// Set 设置配置值
func (r *SystemConfigRepository) Set(key, value, description string) error {
	var config models.SystemConfig
	err := r.db.Where("config_key = ?", key).First(&config).Error
	
	if err == gorm.ErrRecordNotFound {
		// 创建新配置
		config = models.SystemConfig{
			ConfigKey:   key,
			ConfigValue: value,
			Description: description,
		}
		return r.db.Create(&config).Error
	} else if err != nil {
		return err
	}
	
	// 更新现有配置
	config.ConfigValue = value
	if description != "" {
		config.Description = description
	}
	return r.db.Save(&config).Error
}

// GetAll 获取所有配置
func (r *SystemConfigRepository) GetAll() ([]*models.SystemConfig, error) {
	var configs []*models.SystemConfig
	err := r.db.Find(&configs).Error
	return configs, err
}
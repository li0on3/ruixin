package services

import (
	"backend/internal/repository"
	"gorm.io/gorm"
)

type SystemConfigService struct {
	db         *gorm.DB
	configRepo *repository.SystemConfigRepository
}

func NewSystemConfigService(db *gorm.DB) *SystemConfigService {
	return &SystemConfigService{
		db:         db,
		configRepo: repository.NewSystemConfigRepository(db),
	}
}

// GetSyncStoreCode 获取同步店铺代码
func (s *SystemConfigService) GetSyncStoreCode() (string, error) {
	config, err := s.configRepo.GetByKey("sync_store_code")
	if err != nil {
		return "", err
	}
	return config.ConfigValue, nil
}

// SetSyncStoreCode 设置同步店铺代码
func (s *SystemConfigService) SetSyncStoreCode(storeCode string) error {
	return s.configRepo.Set("sync_store_code", storeCode, "商品同步使用的店铺代码")
}

// GetConfig 获取配置值
func (s *SystemConfigService) GetConfig(key string) (string, error) {
	config, err := s.configRepo.GetByKey(key)
	if err != nil {
		return "", err
	}
	return config.ConfigValue, nil
}

// SetConfig 设置配置值
func (s *SystemConfigService) SetConfig(key, value, description string) error {
	return s.configRepo.Set(key, value, description)
}

// GetAllConfigs 获取所有配置
func (s *SystemConfigService) GetAllConfigs() (map[string]string, error) {
	configs, err := s.configRepo.GetAll()
	if err != nil {
		return nil, err
	}
	
	result := make(map[string]string)
	for _, config := range configs {
		result[config.ConfigKey] = config.ConfigValue
	}
	return result, nil
}

// GetSecurityConfig 获取安全配置
func (s *SystemConfigService) GetSecurityConfig() map[string]string {
	// 默认安全配置
	defaultConfigs := map[string]string{
		"security_card_access_mode":   "soft",  // soft: 仅记录, strict: 阻止访问
		"security_rate_limit_enabled": "true",  // 是否启用频率限制
		"security_audit_log_enabled":  "true",  // 是否启用审计日志
		"security_rate_limit_window":  "60",    // 频率限制时间窗口（秒）
		"security_rate_limit_count":   "300",   // 频率限制请求数
	}
	
	// 从数据库获取配置
	for key := range defaultConfigs {
		if value, err := s.GetConfig(key); err == nil && value != "" {
			defaultConfigs[key] = value
		}
	}
	
	return defaultConfigs
}

// IsCardAccessStrict 检查是否启用严格的卡片访问控制
func (s *SystemConfigService) IsCardAccessStrict() bool {
	config, err := s.GetConfig("security_card_access_mode")
	if err != nil {
		return false // 默认软模式
	}
	return config == "strict"
}

// IsRateLimitEnabled 检查是否启用频率限制
func (s *SystemConfigService) IsRateLimitEnabled() bool {
	config, err := s.GetConfig("security_rate_limit_enabled")
	if err != nil {
		return true // 默认启用
	}
	return config == "true"
}

// IsAuditLogEnabled 检查是否启用审计日志
func (s *SystemConfigService) IsAuditLogEnabled() bool {
	config, err := s.GetConfig("security_audit_log_enabled")
	if err != nil {
		return true // 默认启用
	}
	return config == "true"
}

// InitializeSecurityConfigs 初始化安全配置
func (s *SystemConfigService) InitializeSecurityConfigs() error {
	configs := []struct {
		key         string
		value       string
		description string
	}{
		{"security_card_access_mode", "soft", "卡片访问控制模式：soft(仅记录)/strict(阻止访问)"},
		{"security_rate_limit_enabled", "true", "是否启用API频率限制"},
		{"security_audit_log_enabled", "true", "是否启用安全审计日志"},
		{"security_rate_limit_window", "60", "频率限制时间窗口（秒）"},
		{"security_rate_limit_count", "300", "频率限制最大请求数"},
	}
	
	for _, config := range configs {
		// 只有在配置不存在时才创建
		if _, err := s.GetConfig(config.key); err != nil {
			if err := s.SetConfig(config.key, config.value, config.description); err != nil {
				return err
			}
		}
	}
	
	return nil
}
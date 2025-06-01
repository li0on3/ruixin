package services

import (
	"backend/internal/models"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SecurityAuditService 安全审计服务
type SecurityAuditService struct {
	db *gorm.DB
}

// NewSecurityAuditService 创建安全审计服务
func NewSecurityAuditService(db *gorm.DB) *SecurityAuditService {
	return &SecurityAuditService{db: db}
}

// LogSecurityEvent 记录安全事件
func (s *SecurityAuditService) LogSecurityEvent(distributorID uint, action string, resource string, details map[string]interface{}, status string, errorMsg string, c *gin.Context) error {
	log := &models.SecurityAuditLog{
		DistributorID: distributorID,
		Action:        action,
		Resource:      resource,
		Details:       models.JSONMap(details),
		Status:        status,
		ErrorMsg:      errorMsg,
		CreatedAt:     time.Now(),
	}

	if c != nil {
		log.IPAddress = c.ClientIP()
		log.UserAgent = c.GetHeader("User-Agent")
	}

	return s.db.Create(log).Error
}

// LogCardAccess 记录卡片访问
func (s *SecurityAuditService) LogCardAccess(distributorID uint, cardCode string, api string, authorized bool, c *gin.Context) error {
	action := models.AuditActionCardAccess
	status := models.AuditStatusSuccess
	if !authorized {
		action = models.AuditActionUnauthorizedCard
		status = models.AuditStatusWarning
	}

	details := map[string]interface{}{
		"card_code": cardCode,
		"api":       api,
		"authorized": authorized,
	}

	return s.LogSecurityEvent(distributorID, action, cardCode, details, status, "", c)
}

// LogRateLimited 记录频率限制事件
func (s *SecurityAuditService) LogRateLimited(distributorID uint, api string, c *gin.Context) error {
	details := map[string]interface{}{
		"api": api,
		"timestamp": time.Now(),
	}

	return s.LogSecurityEvent(distributorID, models.AuditActionRateLimited, api, details, models.AuditStatusWarning, "请求频率超限", c)
}

// LogOrderCreate 记录订单创建
func (s *SecurityAuditService) LogOrderCreate(distributorID uint, orderNo string, amount float64, success bool, errorMsg string, c *gin.Context) error {
	status := models.AuditStatusSuccess
	if !success {
		status = models.AuditStatusFailed
	}

	details := map[string]interface{}{
		"order_no": orderNo,
		"amount":   amount,
		"success":  success,
	}

	return s.LogSecurityEvent(distributorID, models.AuditActionOrderCreate, orderNo, details, status, errorMsg, c)
}

// GetAuditLogs 获取审计日志
func (s *SecurityAuditService) GetAuditLogs(filters map[string]interface{}, page, pageSize int) ([]models.SecurityAuditLog, int64, error) {
	var logs []models.SecurityAuditLog
	var total int64

	query := s.db.Model(&models.SecurityAuditLog{})

	// 应用过滤条件
	if distributorID, ok := filters["distributor_id"].(uint); ok && distributorID > 0 {
		query = query.Where("distributor_id = ?", distributorID)
	}
	if action, ok := filters["action"].(string); ok && action != "" {
		query = query.Where("action = ?", action)
	}
	if status, ok := filters["status"].(string); ok && status != "" {
		query = query.Where("status = ?", status)
	}
	if startDate, ok := filters["start_date"].(time.Time); ok {
		query = query.Where("created_at >= ?", startDate)
	}
	if endDate, ok := filters["end_date"].(time.Time); ok {
		query = query.Where("created_at <= ?", endDate)
	}

	// 计数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

// GetSecurityMetrics 获取安全指标
func (s *SecurityAuditService) GetSecurityMetrics(distributorID uint, days int) (map[string]interface{}, error) {
	startDate := time.Now().AddDate(0, 0, -days)
	
	metrics := make(map[string]interface{})
	
	// 统计各类事件数量
	var stats []struct {
		Action string
		Count  int64
	}
	
	query := s.db.Model(&models.SecurityAuditLog{}).
		Select("action, COUNT(*) as count").
		Where("created_at >= ?", startDate).
		Group("action")
		
	if distributorID > 0 {
		query = query.Where("distributor_id = ?", distributorID)
	}
	
	if err := query.Find(&stats).Error; err != nil {
		return nil, err
	}
	
	// 转换为map
	actionCounts := make(map[string]int64)
	for _, stat := range stats {
		actionCounts[stat.Action] = stat.Count
	}
	
	metrics["action_counts"] = actionCounts
	metrics["start_date"] = startDate
	metrics["end_date"] = time.Now()
	
	// 计算未授权访问率
	unauthorizedCount := actionCounts[models.AuditActionUnauthorizedCard]
	totalCardAccess := actionCounts[models.AuditActionCardAccess] + unauthorizedCount
	
	if totalCardAccess > 0 {
		metrics["unauthorized_rate"] = float64(unauthorizedCount) / float64(totalCardAccess) * 100
	} else {
		metrics["unauthorized_rate"] = 0
	}
	
	return metrics, nil
}
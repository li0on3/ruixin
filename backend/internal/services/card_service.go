package services

import (
	"backend/internal/models"
	"backend/internal/repository"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"net/url"
	"strings"
	"sync"
	"time"
)

// BatchImportRequest 批量导入请求
type BatchImportRequest struct {
	LuckinProductID int        `json:"luckin_product_id" binding:"required,min=1,max=100"`
	CostPrice       float64    `json:"cost_price" binding:"required,min=0"`
	SellPrice       float64    `json:"sell_price" binding:"required,min=0"`
	CardCodes       []string   `json:"card_codes,omitempty"` // 直接提供卡片代码
	CardURLs        []string   `json:"card_urls,omitempty"`  // 提供卡片URL
	BatchName       string     `json:"batch_name"`           // 批次名称
	ExpiredAt       *time.Time `json:"expired_at"`
	Description     string     `json:"description"`
	AdminID         uint       `json:"-"`
}

type CardService struct {
	db                  *gorm.DB
	cardRepo            *repository.CardRepository
	batchRepo           *repository.CardBatchRepository
	priceService        *LuckinConfigService
	productService      *ProductService
	systemConfigService *SystemConfigService
}

func NewCardService(
	db *gorm.DB,
	cardRepo *repository.CardRepository,
	batchRepo *repository.CardBatchRepository,
	priceService *LuckinConfigService,
	productService *ProductService,
	systemConfigService *SystemConfigService,
) *CardService {
	return &CardService{
		db:                  db,
		cardRepo:            cardRepo,
		batchRepo:           batchRepo,
		priceService:        priceService,
		productService:      productService,
		systemConfigService: systemConfigService,
	}
}

func (s *CardService) GetByID(id uint) (*models.Card, error) {
	return s.cardRepo.GetByID(id)
}

func (s *CardService) GetByCode(code string) (*models.Card, error) {
	return s.cardRepo.GetByCode(code)
}

func (s *CardService) Create(card *models.Card) error {
	// 先验证卡片是否有效
	isValid, message, err := s.ValidateCard(card.CardCode)
	if err != nil {
		return fmt.Errorf("验证卡片失败: %v", err)
	}
	if !isValid {
		return fmt.Errorf("卡片无效: %s", message)
	}

	// 检查卡片代码是否已存在（包括已删除的）
	existingCard, err := s.cardRepo.GetByCode(card.CardCode)
	if err == nil && existingCard != nil {
		return errors.New("card code already exists")
	}

	// 如果是因为软删除导致的重复，先物理删除旧记录
	var deletedCard models.Card
	result := s.db.Unscoped().Where("card_code = ? AND deleted_at IS NOT NULL", card.CardCode).First(&deletedCard)
	if result.Error == nil {
		// 物理删除已软删除的记录
		s.db.Unscoped().Delete(&deletedCard)
	}

	// 创建卡片
	if err := s.cardRepo.Create(card); err != nil {
		return err
	}

	// 异步同步商品（如果启用了自动同步）
	go s.syncProductsIfEnabled(card.CardCode)

	return nil
}

func (s *CardService) Update(card *models.Card) error {
	return s.cardRepo.Update(card)
}

// CreateWithoutValidation 创建卡片（不验证）
func (s *CardService) CreateWithoutValidation(card *models.Card) error {
	// 检查卡片代码是否已存在（包括已删除的）
	existingCard, err := s.cardRepo.GetByCode(card.CardCode)
	if err == nil && existingCard != nil {
		return errors.New("card code already exists")
	}

	// 如果是因为软删除导致的重复，先物理删除旧记录
	var deletedCard models.Card
	result := s.db.Unscoped().Where("card_code = ? AND deleted_at IS NOT NULL", card.CardCode).First(&deletedCard)
	if result.Error == nil {
		// 物理删除已软删除的记录
		s.db.Unscoped().Delete(&deletedCard)
	}

	// 创建卡片
	if err := s.cardRepo.Create(card); err != nil {
		return err
	}

	// 异步同步商品（如果启用了自动同步）
	go s.syncProductsIfEnabled(card.CardCode)

	return nil
}

func (s *CardService) Delete(id uint) error {
	// 检查卡片是否已被使用
	card, err := s.cardRepo.GetByID(id)
	if err != nil {
		return err
	}

	if card.Status == 1 {
		return errors.New("无法删除已使用的卡片")
	}

	return s.cardRepo.Delete(id)
}

func (s *CardService) List(offset, limit int, filters map[string]interface{}) ([]*models.Card, int64, error) {
	return s.cardRepo.List(offset, limit, filters)
}

func (s *CardService) GetUsageLogs(cardID uint, offset, limit int) ([]*models.CardUsageLog, int64, error) {
	return s.cardRepo.GetUsageLogs(cardID, offset, limit)
}

// GetAvailableCards 获取可用的卡片（未使用且未过期）
func (s *CardService) GetAvailableCards(limit int) ([]*models.Card, error) {
	var cards []*models.Card

	err := s.db.Where("status = ? AND (expired_at IS NULL OR expired_at > ?)", 0, time.Now()).
		Order("created_at DESC").
		Limit(limit).
		Find(&cards).Error

	if err != nil {
		return nil, err
	}

	return cards, nil
}

// GetAvailableCard 获取单个可用的卡片
func (s *CardService) GetAvailableCard() (*models.Card, error) {
	cards, err := s.GetAvailableCards(1)
	if err != nil {
		return nil, err
	}

	if len(cards) == 0 {
		return nil, errors.New("no available cards")
	}

	return cards[0], nil
}

// GetCardStats 获取卡片统计信息
func (s *CardService) GetCardStats() (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// 总卡片数
	var totalCards int64
	s.db.Model(&models.Card{}).Count(&totalCards)
	stats["total_cards"] = totalCards

	// 未使用卡片数
	var unusedCards int64
	s.db.Model(&models.Card{}).Where("status = ?", 0).Count(&unusedCards)
	stats["unused_cards"] = unusedCards

	// 已使用卡片数
	var usedCards int64
	s.db.Model(&models.Card{}).Where("status = ?", 1).Count(&usedCards)
	stats["used_cards"] = usedCards

	// 预占中卡片数
	var reservedCards int64
	s.db.Model(&models.Card{}).Where("status = ?", 2).Count(&reservedCards)
	stats["reserved_cards"] = reservedCards

	// 按价格分组统计
	type PriceStats struct {
		PriceID     int64   `json:"price_id"`
		PriceValue  float64 `json:"price_value"`
		TotalCount  int64   `json:"total_count"`
		UnusedCount int64   `json:"unused_count"`
		UsedCount   int64   `json:"used_count"`
	}

	var priceStats []PriceStats
	s.db.Table("cards").
		Select("cards.price_id, luckin_prices.price_value, COUNT(*) as total_count, " +
			"SUM(CASE WHEN cards.status = 0 THEN 1 ELSE 0 END) as unused_count, " +
			"SUM(CASE WHEN cards.status = 1 THEN 1 ELSE 0 END) as used_count").
		Joins("LEFT JOIN luckin_prices ON cards.price_id = luckin_prices.id").
		Group("cards.price_id, luckin_prices.price_value").
		Order("luckin_prices.price_value").
		Scan(&priceStats)

	stats["price_stats"] = priceStats

	return stats, nil
}

// ListBatches 获取批次列表
func (s *CardService) ListBatches(offset, limit int) ([]*models.CardBatch, int64, error) {
	filters := make(map[string]interface{})
	return s.batchRepo.List(offset, limit, filters)
}

// GetBatch 获取批次详情
func (s *CardService) GetBatch(id uint) (*models.CardBatch, error) {
	return s.batchRepo.GetByID(id)
}

// GetBatchCards 获取批次下的卡片
func (s *CardService) GetBatchCards(batchID uint, offset, limit int) ([]*models.Card, int64, error) {
	var cards []*models.Card
	var total int64

	// 计算总数
	err := s.db.Model(&models.Card{}).
		Where("batch_id = ?", batchID).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err = s.db.Preload("Price").
		Where("batch_id = ?", batchID).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&cards).Error
	if err != nil {
		return nil, 0, err
	}

	return cards, total, nil
}

// BatchImport 批量导入卡片
func (s *CardService) BatchImport(batchReq *BatchImportRequest) (*models.CardBatch, error) {
	// 验证瑞幸产品ID
	if batchReq.LuckinProductID < 1 || batchReq.LuckinProductID > 100 {
		return nil, errors.New("瑞幸产品ID必须在1-100之间")
	}

	// 处理卡片代码：从URL中提取或直接使用提供的代码
	allCardCodes := make([]string, 0)

	// 添加直接提供的卡片代码
	allCardCodes = append(allCardCodes, batchReq.CardCodes...)

	// 从URL中提取卡片代码
	for _, cardURL := range batchReq.CardURLs {
		cardCode, err := s.extractCardCodeFromURL(cardURL)
		if err != nil {
			// 跳过无效的URL
			continue
		}
		allCardCodes = append(allCardCodes, cardCode)
	}

	// 去重
	uniqueCardCodes := s.removeDuplicates(allCardCodes)

	if len(uniqueCardCodes) == 0 {
		return nil, errors.New("没有有效的卡片代码")
	}

	// 创建批次
	batch := &models.CardBatch{
		BatchNo:         s.batchRepo.GenerateBatchNo(),
		LuckinProductID: batchReq.LuckinProductID,
		CostPrice:       batchReq.CostPrice,
		SellPrice:       batchReq.SellPrice,
		TotalCount:      len(uniqueCardCodes),
		UsedCount:       0,
		ImportedAt:      time.Now(),
		ImportedBy:      batchReq.AdminID,
		Description:     batchReq.Description,
	}

	if err := s.batchRepo.Create(batch); err != nil {
		return nil, err
	}

	// 批量创建卡片
	for _, code := range uniqueCardCodes {
		var expiredAt *time.Time
		if batchReq.ExpiredAt != nil {
			expiredAt = batchReq.ExpiredAt
		} else {
			defaultExpired := time.Now().AddDate(1, 0, 0)
			expiredAt = &defaultExpired
		}

		card := &models.Card{
			CardCode:        code,
			BatchID:         &batch.ID,
			LuckinProductID: batchReq.LuckinProductID,
			CostPrice:       batchReq.CostPrice,
			SellPrice:       batchReq.SellPrice,
			Status:          0, // 未使用
			ExpiredAt:       expiredAt,
			Description:     batchReq.Description,
		}

		if err := s.cardRepo.Create(card); err != nil {
			// 如果是重复的卡片代码，跳过
			if err.Error() == "card code already exists" {
				continue
			}
			return nil, err
		}
	}

	// 异步同步商品
	go func() {
		// 对每个唯一的卡片代码执行同步
		for _, code := range uniqueCardCodes {
			s.syncProductsIfEnabled(code)
			time.Sleep(1 * time.Second) // 避免请求频率过高
		}
	}()

	return batch, nil
}

// GetAvailableByPriceID 根据价格ID获取可用卡片
func (s *CardService) GetAvailableByPriceID(priceID int64) (*models.Card, error) {
	return s.cardRepo.GetAvailableByPriceID(priceID)
}

// MarkAsUsed 标记卡片为已使用
func (s *CardService) MarkAsUsed(cardID uint, orderID uint) error {
	card, err := s.cardRepo.GetByID(cardID)
	if err != nil {
		return err
	}

	if card.Status != 0 {
		return errors.New("卡片状态不正确")
	}

	// 标记为已使用
	if err := s.cardRepo.MarkAsUsed(cardID, orderID); err != nil {
		return err
	}

	// 更新批次使用数量
	if card.BatchID != nil {
		return s.batchRepo.UpdateUsedCount(*card.BatchID)
	}

	return nil
}

// GetBatchList 获取批次列表
func (s *CardService) GetBatchList(offset, limit int, filters map[string]interface{}) ([]*models.CardBatch, int64, error) {
	return s.batchRepo.List(offset, limit, filters)
}

// GetBatchByID 获取批次详情
func (s *CardService) GetBatchByID(id uint) (*models.CardBatch, error) {
	return s.batchRepo.GetByID(id)
}

// GetCardsByBatchID 获取批次下的卡片
func (s *CardService) GetCardsByBatchID(batchID uint, offset, limit int) ([]*models.Card, int64, error) {
	return s.cardRepo.GetCardsByBatchID(batchID, offset, limit)
}

// GetAvailableCountByPriceID 获取某个价格下可用卡片数量
func (s *CardService) GetAvailableCountByPriceID(priceID int64) (int64, error) {
	return s.cardRepo.GetAvailableCountByPriceID(priceID)
}

// GetPriceStats 获取价格相关的卡片统计
func (s *CardService) GetPriceStats(priceID int64) (map[string]interface{}, error) {
	return s.batchRepo.GetStatsByPriceID(priceID)
}

// ReserveCard 预占卡片（标记为使用中，但不关联订单）
func (s *CardService) ReserveCard(cardID uint) error {
	now := time.Now()
	return s.db.Model(&models.Card{}).Where("id = ?", cardID).Updates(map[string]interface{}{
		"status":      2,
		"reserved_at": &now,
	}).Error
}

// ReleaseCard 释放预占的卡片
func (s *CardService) ReleaseCard(cardID uint) error {
	// 只释放预占状态的卡片，已使用的不能释放
	card, err := s.cardRepo.GetByID(cardID)
	if err != nil {
		return err
	}

	if card.Status == 2 { // 只有预占状态才能释放
		return s.db.Model(&models.Card{}).Where("id = ?", cardID).Updates(map[string]interface{}{
			"status":      0,
			"reserved_at": nil,
		}).Error
	}

	return nil
}

// GetAnyAvailableCard 获取任意一张可用的卡片（用于店铺搜索等场景）
func (s *CardService) GetAnyAvailableCard() (*models.Card, error) {
	return s.cardRepo.GetAnyAvailable()
}

// GetBoundProducts 获取卡片绑定的商品
func (s *CardService) GetBoundProducts(cardID uint) ([]*models.Product, error) {
	var bindings []models.CardProductBinding
	err := s.db.Preload("Product").Where("card_id = ? AND is_active = ?", cardID, true).Order("priority DESC").Find(&bindings).Error
	if err != nil {
		return nil, err
	}

	products := make([]*models.Product, len(bindings))
	for i, binding := range bindings {
		products[i] = binding.Product
	}
	return products, nil
}

// GetBoundProductCount 获取卡片绑定的商品数量
func (s *CardService) GetBoundProductCount(cardID uint) (int64, error) {
	var count int64
	err := s.db.Model(&models.CardProductBinding{}).Where("card_id = ? AND is_active = ?", cardID, true).Count(&count).Error
	return count, err
}

// syncProductsIfEnabled 如果启用了自动同步，则同步商品
func (s *CardService) syncProductsIfEnabled(cardCode string) {
	// 检查是否启用了自动同步
	syncEnabled, err := s.systemConfigService.GetConfig("sync_enabled")
	if err != nil || syncEnabled != "true" {
		return
	}

	// 获取同步使用的门店代码
	storeCode, err := s.systemConfigService.GetSyncStoreCode()
	if err != nil || storeCode == "" {
		storeCode = "390840" // 默认门店代码
	}

	// 调用商品同步服务
	_, err = s.productService.SyncProductsFromCard(cardCode, storeCode)
	if err != nil {
		// 记录错误日志，但不影响卡片创建
		// 可以考虑添加日志记录
	}
}

// ValidateCard 验证卡片是否可用
func (s *CardService) ValidateCard(cardCode string) (bool, string, error) {
	// 获取瑞幸客户端
	if s.productService == nil {
		return false, "产品服务未初始化", errors.New("product service not initialized")
	}

	luckinClient := s.productService.GetLuckinClient()
	if luckinClient == nil {
		return false, "瑞幸客户端未初始化", errors.New("luckin client not initialized")
	}

	// 记录验证开始
	fmt.Printf("开始验证卡片: %s\n", cardCode)

	// 调用瑞幸API验证卡片
	resp, err := luckinClient.CheckGoodsCard(cardCode)
	if err != nil {
		fmt.Printf("调用瑞幸API失败: %v\n", err)
		return false, "验证失败", fmt.Errorf("API调用失败: %v", err)
	}

	// 记录响应
	fmt.Printf("瑞幸API响应: code=%d, msg=%s, status=%d\n", resp.Code, resp.Msg, resp.Data.Status)

	// 解析响应
	if resp.Code != 200 {
		return false, resp.Msg, nil
	}

	// status: 1 表示可用
	if resp.Data.Status == 1 {
		return true, "卡片有效", nil
	}

	// 其他状态都视为不可用
	statusMsg := "卡片无效或已被使用"
	if resp.Data.Status == 0 {
		statusMsg = "卡片不存在"
	} else if resp.Data.Status == 2 {
		statusMsg = "卡片已被使用"
	} else if resp.Data.Status == 3 {
		statusMsg = "卡片已过期"
	}

	return false, statusMsg, nil
}

// ValidateAndUpdateCard 验证卡片并更新状态
func (s *CardService) ValidateAndUpdateCard(cardCode string) (bool, string, bool, error) {
	// 先验证卡片
	isValid, message, err := s.ValidateCard(cardCode)
	if err != nil {
		return false, message, false, err
	}

	// 查找数据库中的卡片
	card, err := s.cardRepo.GetByCode(cardCode)
	if err != nil {
		// 卡片不在我们的数据库中
		return isValid, message, false, nil
	}

	// 检查是否需要更新状态
	needUpdate := false
	now := time.Now()
	updateData := map[string]interface{}{
		"sync_status": "synced",
		"synced_at":   &now,
	}

	if isValid && card.Status == 1 {
		// 瑞幸系统显示可用，但本地标记为已使用，恢复为未使用
		updateData["status"] = 0
		updateData["used_at"] = gorm.Expr("NULL")
		updateData["order_id"] = gorm.Expr("NULL")
		needUpdate = true
		message = "卡片有效，状态已从已使用恢复为未使用"
	} else if !isValid && card.Status == 0 {
		// 瑞幸系统显示不可用，但本地标记为未使用，更新为已使用
		updateData["status"] = 1
		updateData["used_at"] = &now
		needUpdate = true
		message = message + "，状态已从未使用更新为已使用"
	} else if card.Status == 2 {
		// 预占中的卡片不更新状态
		message = message + "（卡片预占中，未更新状态）"
	}

	// 执行更新
	if needUpdate {
		if err := s.db.Model(card).Updates(updateData).Error; err != nil {
			return isValid, message, false, fmt.Errorf("更新卡片状态失败: %v", err)
		}
	}

	return isValid, message, needUpdate, nil
}

// BatchValidateResult 批量验证结果
type BatchValidateResult struct {
	TotalCards    int `json:"total_cards"`
	ValidCards    int `json:"valid_cards"`
	InvalidCards  int `json:"invalid_cards"`
	FailedCards   int `json:"failed_cards"`
	UpdatedCards  int `json:"updated_cards"`
	RestoredCards int `json:"restored_cards"` // 从已使用恢复为未使用的卡片数
	MarkedUsed    int `json:"marked_used"`    // 从未使用标记为已使用的卡片数
}

// ValidationMode 验证模式枚举
type ValidationMode string

const (
	ValidationModeAll   ValidationMode = "all"   // 全量验证
	ValidationModeSmart ValidationMode = "smart" // 智能验证
)

// ValidationTask 验证任务
type ValidationTask struct {
	ID          string              `json:"id"`
	Mode        ValidationMode      `json:"mode"`
	Status      string              `json:"status"` // queued/running/completed/failed/cancelled
	Progress    *ValidationProgress `json:"progress"`
	CreatedAt   time.Time           `json:"created_at"`
	StartedAt   *time.Time          `json:"started_at,omitempty"`
	CompletedAt *time.Time          `json:"completed_at,omitempty"`
	Error       string              `json:"error,omitempty"`
}

// ValidationProgress 验证进度
type ValidationProgress struct {
	Total     int `json:"total"`
	Processed int `json:"processed"`
	Valid     int `json:"valid"`
	Invalid   int `json:"invalid"`
	Failed    int `json:"failed"`
	Skipped   int `json:"skipped"`
}

// StartBatchValidation 启动批量验证（支持两种模式）
func (s *CardService) StartBatchValidation(mode ValidationMode) (*ValidationTask, error) {
	taskID := fmt.Sprintf("validation_%d", time.Now().UnixNano())
	
	task := &ValidationTask{
		ID:        taskID,
		Mode:      mode,
		Status:    "queued",
		Progress:  &ValidationProgress{},
		CreatedAt: time.Now(),
	}
	
	// 保存任务到Redis（如果有的话）或内存
	s.saveValidationTask(task)
	
	// 异步执行验证
	go s.executeValidationTask(task)
	
	return task, nil
}

// executeValidationTask 执行验证任务
func (s *CardService) executeValidationTask(task *ValidationTask) {
	defer func() {
		if r := recover(); r != nil {
			task.Status = "failed"
			task.Error = fmt.Sprintf("验证任务异常: %v", r)
			now := time.Now()
			task.CompletedAt = &now
			s.saveValidationTask(task)
		}
	}()
	
	// 更新任务状态
	task.Status = "running"
	now := time.Now()
	task.StartedAt = &now
	s.saveValidationTask(task)
	
	switch task.Mode {
	case ValidationModeAll:
		s.executeFullValidation(task)
	case ValidationModeSmart:
		s.executeSmartValidation(task)
	default:
		task.Status = "failed"
		task.Error = "未知验证模式"
	}
	
	// 完成任务
	task.Status = "completed"
	completedAt := time.Now()
	task.CompletedAt = &completedAt
	s.saveValidationTask(task)
}

// executeFullValidation 执行全量验证（安全模式）
func (s *CardService) executeFullValidation(task *ValidationTask) {
	// 获取所有需要验证的卡片（除了预占中的）
	var cards []models.Card
	err := s.db.Where("status IN ?", []int{0, 1}).Find(&cards).Error
	if err != nil {
		task.Status = "failed"
		task.Error = fmt.Sprintf("获取卡片列表失败: %v", err)
		return
	}
	
	task.Progress.Total = len(cards)
	s.saveValidationTask(task)
	
	// 安全验证参数
	config := struct {
		batchSize     int
		batchInterval time.Duration
		cardInterval  time.Duration
	}{
		batchSize:     5,                    // 每批5张
		batchInterval: 2 * time.Minute,      // 批间间隔2分钟
		cardInterval:  10 * time.Second,     // 卡片间隔10秒
	}
	
	// 分批处理
	for i := 0; i < len(cards); i += config.batchSize {
		// 检查任务是否被取消
		if s.isTaskCancelled(task.ID) {
			task.Status = "cancelled"
			return
		}
		
		end := i + config.batchSize
		if end > len(cards) {
			end = len(cards)
		}
		batch := cards[i:end]
		
		// 处理这批卡片
		for _, card := range batch {
			if s.isTaskCancelled(task.ID) {
				task.Status = "cancelled"
				return
			}
			
			s.validateSingleCard(card, task)
			time.Sleep(config.cardInterval)
		}
		
		// 批间休息
		if end < len(cards) {
			time.Sleep(config.batchInterval)
		}
	}
}

// executeSmartValidation 执行智能验证（只验证重要卡片）
func (s *CardService) executeSmartValidation(task *ValidationTask) {
	// 分层获取需要验证的卡片
	priorityCards := s.getPriorityCards()
	
	task.Progress.Total = len(priorityCards)
	s.saveValidationTask(task)
	
	// 快速验证参数
	config := struct {
		cardInterval time.Duration
		maxDuration  time.Duration
	}{
		cardInterval: 3 * time.Second,     // 卡片间隔3秒
		maxDuration:  10 * time.Minute,   // 最大10分钟
	}
	
	startTime := time.Now()
	
	for _, card := range priorityCards {
		// 检查超时
		if time.Since(startTime) > config.maxDuration {
			break
		}
		
		// 检查任务是否被取消
		if s.isTaskCancelled(task.ID) {
			task.Status = "cancelled"
			return
		}
		
		s.validateSingleCard(card, task)
		time.Sleep(config.cardInterval)
	}
}

// getPriorityCards 获取需要优先验证的卡片
func (s *CardService) getPriorityCards() []models.Card {
	var cards []models.Card
	
	// 分层查询策略
	now24h := time.Now().Add(-24 * time.Hour)
	now7d := time.Now().Add(-7 * 24 * time.Hour)
	
	query := s.db.Where(`
		(
			(status = 0 AND (synced_at IS NULL OR synced_at < ?)) OR
			(status = 1 AND used_at > ? AND (synced_at IS NULL OR synced_at < used_at)) OR
			(order_id IS NOT NULL AND status = 0) OR
			(created_at > ? AND (synced_at IS NULL OR sync_status = 'failed'))
		) AND status != 2
	`, now24h, now7d, now24h).Order(`
		CASE 
			WHEN order_id IS NOT NULL AND status = 0 THEN 1
			WHEN status = 0 AND (synced_at IS NULL OR synced_at < '` + now24h.Format("2006-01-02 15:04:05") + `') THEN 2
			WHEN status = 1 AND used_at > '` + now7d.Format("2006-01-02 15:04:05") + `' THEN 3
			ELSE 4
		END, created_at DESC
	`)
	
	err := query.Limit(200).Find(&cards).Error // 最多验证200张
	if err != nil {
		return []models.Card{}
	}
	
	return cards
}

// validateSingleCard 验证单张卡片
func (s *CardService) validateSingleCard(card models.Card, task *ValidationTask) {
	isValid, _, err := s.ValidateCard(card.CardCode)
	
	// 更新进度
	task.Progress.Processed++
	
	if err != nil {
		task.Progress.Failed++
		// 更新卡片同步状态为失败
		now := time.Now()
		s.db.Model(&card).Updates(map[string]interface{}{
			"sync_status": "failed",
			"synced_at":   &now,
		})
	} else {
		// 根据验证结果更新卡片状态
		now := time.Now()
		updateData := map[string]interface{}{
			"sync_status": "synced",
			"synced_at":   &now,
		}
		
		if isValid {
			task.Progress.Valid++
			// 卡片在瑞幸系统中可用
			if card.Status == 1 {
				// 本地显示已使用，但瑞幸显示可用，恢复为未使用
				updateData["status"] = 0
				updateData["used_at"] = nil
				updateData["order_id"] = nil
			}
		} else {
			task.Progress.Invalid++
			// 卡片在瑞幸系统中不可用
			if card.Status == 0 {
				// 本地显示未使用，但瑞幸显示已使用，更新为已使用
				updateData["status"] = 1
				updateData["used_at"] = &now
			}
		}
		
		s.db.Model(&card).Updates(updateData)
	}
	
	// 保存进度
	s.saveValidationTask(task)
}

// 任务管理相关方法
var validationTasks = make(map[string]*ValidationTask)
var taskMutex = sync.RWMutex{}

// saveValidationTask 保存验证任务
func (s *CardService) saveValidationTask(task *ValidationTask) {
	taskMutex.Lock()
	defer taskMutex.Unlock()
	validationTasks[task.ID] = task
}

// getValidationTask 获取验证任务
func (s *CardService) GetValidationTask(taskID string) (*ValidationTask, error) {
	taskMutex.RLock()
	defer taskMutex.RUnlock()
	
	task, exists := validationTasks[taskID]
	if !exists {
		return nil, fmt.Errorf("任务不存在: %s", taskID)
	}
	
	return task, nil
}

// isTaskCancelled 检查任务是否被取消
func (s *CardService) isTaskCancelled(taskID string) bool {
	taskMutex.RLock()
	defer taskMutex.RUnlock()
	
	task, exists := validationTasks[taskID]
	if !exists {
		return true
	}
	
	return task.Status == "cancelled"
}

// CancelValidationTask 取消验证任务
func (s *CardService) CancelValidationTask(taskID string) error {
	taskMutex.Lock()
	defer taskMutex.Unlock()
	
	task, exists := validationTasks[taskID]
	if !exists {
		return fmt.Errorf("任务不存在: %s", taskID)
	}
	
	if task.Status == "running" {
		task.Status = "cancelled"
		now := time.Now()
		task.CompletedAt = &now
	}
	
	return nil
}

// GetValidationStatistics 获取验证统计信息
func (s *CardService) GetValidationStatistics() map[string]interface{} {
	var stats struct {
		TotalCards         int64 `json:"total_cards"`
		UnusedCards        int64 `json:"unused_cards"`
		UsedCards          int64 `json:"used_cards"`
		ReservedCards      int64 `json:"reserved_cards"`
		NeedValidation     int64 `json:"need_validation"`
		PriorityValidation int64 `json:"priority_validation"`
	}
	
	// 总卡片数
	s.db.Model(&models.Card{}).Count(&stats.TotalCards)
	
	// 各状态卡片数
	s.db.Model(&models.Card{}).Where("status = ?", 0).Count(&stats.UnusedCards)
	s.db.Model(&models.Card{}).Where("status = ?", 1).Count(&stats.UsedCards)
	s.db.Model(&models.Card{}).Where("status = ?", 2).Count(&stats.ReservedCards)
	
	// 需要验证的卡片数（24小时内未验证的）
	s.db.Model(&models.Card{}).Where("synced_at IS NULL OR synced_at < ?", time.Now().Add(-24*time.Hour)).Count(&stats.NeedValidation)
	
	// 优先验证的卡片数
	priorityCards := s.getPriorityCards()
	stats.PriorityValidation = int64(len(priorityCards))
	
	return map[string]interface{}{
		"total_cards":         stats.TotalCards,
		"unused_cards":        stats.UnusedCards,
		"used_cards":          stats.UsedCards,
		"reserved_cards":      stats.ReservedCards,
		"need_validation":     stats.NeedValidation,
		"priority_validation": stats.PriorityValidation,
	}
}

// BatchValidateAndUpdateCards 批量验证并更新卡片状态（保持兼容性）
func (s *CardService) BatchValidateAndUpdateCards() (*BatchValidateResult, error) {
	result := &BatchValidateResult{}

	// 先统计所有未删除的卡片数量
	var totalCount int64
	s.db.Model(&models.Card{}).Count(&totalCount)
	fmt.Printf("数据库中总卡片数: %d\n", totalCount)

	// 统计各种状态的卡片数量
	var statusCounts []struct {
		Status int
		Count  int64
	}
	s.db.Model(&models.Card{}).Select("status, count(*) as count").Group("status").Scan(&statusCounts)
	for _, sc := range statusCounts {
		fmt.Printf("状态 %d 的卡片数: %d\n", sc.Status, sc.Count)
	}

	// 获取所有卡片（包括已使用的），但排除预占中的
	var cards []models.Card
	err := s.db.Where("status IN ?", []int{0, 1}).Find(&cards).Error
	if err != nil {
		return nil, fmt.Errorf("获取卡片列表失败: %v", err)
	}

	result.TotalCards = len(cards)
	fmt.Printf("开始批量验证 %d 张卡片（包括未使用和已使用的）\n", len(cards))

	// 如果没有卡片，直接返回
	if len(cards) == 0 {
		return result, nil
	}

	// 批量更新同步状态为 syncing
	cardIDs := make([]uint, len(cards))
	for i, card := range cards {
		cardIDs[i] = card.ID
	}

	now := time.Now()
	s.db.Model(&models.Card{}).Where("id IN ?", cardIDs).Updates(map[string]interface{}{
		"sync_status": "syncing",
		"synced_at":   &now,
	})

	// 逐个验证并更新状态
	for i, card := range cards {
		// 避免请求过于频繁
		if i > 0 && i%10 == 0 {
			time.Sleep(1 * time.Second)
		}

		isValid, _, err := s.ValidateCard(card.CardCode)
		if err != nil {
			fmt.Printf("验证卡片 %s 失败: %v\n", card.CardCode, err)
			result.FailedCards++

			// 更新同步状态为失败
			s.db.Model(&card).Updates(map[string]interface{}{
				"sync_status": "failed",
				"synced_at":   &now,
			})
			continue
		}

		// 根据验证结果更新状态
		updateData := map[string]interface{}{
			"sync_status": "synced",
			"synced_at":   &now,
		}

		// 处理验证结果
		if isValid {
			// 卡片在瑞幸系统中是可用的
			if card.Status == 1 {
				// 如果本地标记为已使用，但瑞幸系统显示可用，则恢复为未使用
				updateData["status"] = 0
				updateData["used_at"] = nil
				updateData["order_id"] = nil
				result.UpdatedCards++
				fmt.Printf("卡片 %s 状态从已使用恢复为未使用\n", card.CardCode)
			}
			result.ValidCards++
		} else {
			// 卡片在瑞幸系统中不可用
			if card.Status == 0 {
				// 如果本地标记为未使用，但瑞幸系统显示已使用，则更新为已使用
				updateData["status"] = 1
				updateData["used_at"] = &now
				result.UpdatedCards++
				fmt.Printf("卡片 %s 状态从未使用更新为已使用\n", card.CardCode)
			}
			result.InvalidCards++
		}

		err = s.db.Model(&card).Updates(updateData).Error
		if err != nil {
			fmt.Printf("更新卡片 %s 状态失败: %v\n", card.CardCode, err)
			result.FailedCards++
		}
	}

	fmt.Printf("批量验证完成: 总计 %d 张，有效 %d 张，无效 %d 张，失败 %d 张，状态更新 %d 张\n",
		result.TotalCards, result.ValidCards, result.InvalidCards, result.FailedCards, result.UpdatedCards)

	return result, nil
}

// extractCardCodeFromURL 从URL中提取卡片代码
func (s *CardService) extractCardCodeFromURL(cardURL string) (string, error) {
	// 解析URL
	parsedURL, err := url.Parse(cardURL)
	if err != nil {
		return "", fmt.Errorf("无效的URL: %v", err)
	}

	// 获取查询参数
	queryParams := parsedURL.Query()
	cardCode := queryParams.Get("card")

	if cardCode == "" {
		return "", fmt.Errorf("URL中没有找到card参数")
	}

	// 验证卡片代码格式（可根据实际情况调整）
	cardCode = strings.TrimSpace(cardCode)
	if len(cardCode) < 4 || len(cardCode) > 20 {
		return "", fmt.Errorf("卡片代码格式无效")
	}

	return cardCode, nil
}

// removeDuplicates 去除重复的卡片代码
func (s *CardService) removeDuplicates(cardCodes []string) []string {
	seen := make(map[string]bool)
	unique := make([]string, 0)

	for _, code := range cardCodes {
		code = strings.TrimSpace(code)
		if code != "" && !seen[code] {
			seen[code] = true
			unique = append(unique, code)
		}
	}

	return unique
}

// CheckDistributorCardAccess 检查分销商是否有权限使用指定卡片
// 权限判断逻辑：
// 1. 如果分销商之前使用过这个卡片（有订单记录），则认为有权限
// 2. 如果系统配置允许，可以访问未使用的卡片（可选）
func (s *CardService) CheckDistributorCardAccess(distributorID uint, cardCode string) (bool, error) {
	// 首先检查卡片是否存在且有效
	var card models.Card
	err := s.db.Where("card_code = ? AND status = ?", cardCode, models.CardStatusNormal).First(&card).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil // 卡片不存在或无效
		}
		return false, err
	}

	// 检查是否有使用记录
	var count int64
	err = s.db.Table("orders").
		Where("distributor_id = ? AND card_id = ?", distributorID, card.ID).
		Count(&count).Error
	if err != nil {
		return false, err
	}

	// 如果有订单记录，则有权限
	if count > 0 {
		return true, nil
	}

	// 检查是否有特殊授权（可以通过系统配置控制）
	// 这里可以添加更多的权限判断逻辑，比如：
	// - 检查分销商的特殊权限标记
	// - 检查卡片是否属于公共池
	// - 检查是否在白名单中等

	// 默认情况下，如果没有使用过，则没有权限
	return false, nil
}

// GetDistributorCards 获取分销商有权访问的卡片列表
func (s *CardService) GetDistributorCards(distributorID uint) ([]string, error) {
	var cardCodes []string
	
	// 查询该分销商使用过的所有卡片
	err := s.db.Table("orders").
		Joins("JOIN cards ON cards.id = orders.card_id").
		Where("orders.distributor_id = ? AND cards.status = ?", distributorID, models.CardStatusNormal).
		Distinct("cards.card_code").
		Pluck("cards.card_code", &cardCodes).Error
		
	if err != nil {
		return nil, err
	}
	
	return cardCodes, nil
}

// ValidateCardForDistributor 验证卡片对分销商的可用性
func (s *CardService) ValidateCardForDistributor(distributorID uint, cardCode string) (*models.Card, error) {
	// 检查权限
	hasAccess, err := s.CheckDistributorCardAccess(distributorID, cardCode)
	if err != nil {
		return nil, fmt.Errorf("检查权限失败: %v", err)
	}
	
	if !hasAccess {
		return nil, fmt.Errorf("无权使用该卡片")
	}
	
	// 获取卡片信息
	var card models.Card
	err = s.db.Where("card_code = ? AND status = ?", cardCode, models.CardStatusNormal).First(&card).Error
	if err != nil {
		return nil, fmt.Errorf("卡片不存在或已失效")
	}
	
	return &card, nil
}

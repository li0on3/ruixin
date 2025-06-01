package repository

import (
	"backend/internal/models"
	"fmt"
	"time"
	"gorm.io/gorm"
)

type CardBatchRepository struct {
	*BaseRepository
}

func NewCardBatchRepository(db *gorm.DB) *CardBatchRepository {
	return &CardBatchRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

func (r *CardBatchRepository) Create(batch *models.CardBatch) error {
	return r.db.Create(batch).Error
}

func (r *CardBatchRepository) Update(batch *models.CardBatch) error {
	return r.db.Save(batch).Error
}

func (r *CardBatchRepository) GetByID(id uint) (*models.CardBatch, error) {
	var batch models.CardBatch
	err := r.db.Preload("Price").Preload("Admin").First(&batch, id).Error
	if err != nil {
		return nil, err
	}
	return &batch, nil
}

func (r *CardBatchRepository) GetByBatchNo(batchNo string) (*models.CardBatch, error) {
	var batch models.CardBatch
	err := r.db.Where("batch_no = ?", batchNo).First(&batch).Error
	if err != nil {
		return nil, err
	}
	return &batch, nil
}

func (r *CardBatchRepository) List(offset, limit int, filters map[string]interface{}) ([]*models.CardBatch, int64, error) {
	var batches []*models.CardBatch
	var total int64

	query := r.db.Model(&models.CardBatch{}).Preload("Price").Preload("Admin")

	// Apply filters
	if priceID, ok := filters["price_id"]; ok {
		query = query.Where("price_id = ?", priceID)
	}
	if importedBy, ok := filters["imported_by"]; ok {
		query = query.Where("imported_by = ?", importedBy)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get data
	err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&batches).Error
	if err != nil {
		return nil, 0, err
	}

	return batches, total, nil
}

func (r *CardBatchRepository) Delete(id uint) error {
	return r.db.Delete(&models.CardBatch{}, id).Error
}

// GenerateBatchNo 生成批次号
func (r *CardBatchRepository) GenerateBatchNo() string {
	return fmt.Sprintf("BATCH-%s-%d", time.Now().Format("20060102"), time.Now().Unix()%10000)
}

// UpdateUsedCount 更新已使用数量
func (r *CardBatchRepository) UpdateUsedCount(batchID uint) error {
	var usedCount int64
	err := r.db.Model(&models.Card{}).
		Where("batch_id = ? AND status = 1", batchID).
		Count(&usedCount).Error
	if err != nil {
		return err
	}

	return r.db.Model(&models.CardBatch{}).
		Where("id = ?", batchID).
		Update("used_count", usedCount).Error
}

// GetStatsByPriceID 获取某个价格的批次统计
func (r *CardBatchRepository) GetStatsByPriceID(priceID int64) (map[string]interface{}, error) {
	var totalCount, usedCount int64

	// 获取总数量和已使用数量
	err := r.db.Model(&models.Card{}).
		Where("price_id = ?", priceID).
		Count(&totalCount).Error
	if err != nil {
		return nil, err
	}

	err = r.db.Model(&models.Card{}).
		Where("price_id = ? AND status = 1", priceID).
		Count(&usedCount).Error
	if err != nil {
		return nil, err
	}

	// 计算总成本和总收入
	type Result struct {
		TotalCost    float64
		TotalRevenue float64
	}
	var result Result
	err = r.db.Model(&models.Card{}).
		Select("SUM(cost_price) as total_cost, SUM(CASE WHEN status = 1 THEN sell_price ELSE 0 END) as total_revenue").
		Where("price_id = ?", priceID).
		Scan(&result).Error
	if err != nil {
		return nil, err
	}

	// 计算使用率，避免除零错误
	usageRate := float64(0)
	if totalCount > 0 {
		usageRate = float64(usedCount) / float64(totalCount) * 100
	}

	// 计算利润率，避免除零错误
	profitMargin := float64(0)
	totalProfit := result.TotalRevenue - result.TotalCost
	if result.TotalRevenue > 0 {
		profitMargin = (totalProfit / result.TotalRevenue) * 100
	}

	return map[string]interface{}{
		"total_count":     totalCount,
		"used_count":      usedCount,
		"available_count": totalCount - usedCount,
		"usage_rate":      usageRate,
		"total_cost":      result.TotalCost,
		"total_revenue":   result.TotalRevenue,
		"total_profit":    totalProfit,
		"profit_margin":   profitMargin,
	}, nil
}
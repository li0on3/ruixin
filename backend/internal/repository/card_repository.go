package repository

import (
	"backend/internal/models"
	"gorm.io/gorm"
)

type CardRepository struct {
	*BaseRepository
}

func NewCardRepository(db *gorm.DB) *CardRepository {
	return &CardRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

func (r *CardRepository) Create(card *models.Card) error {
	return r.db.Create(card).Error
}

func (r *CardRepository) Update(card *models.Card) error {
	// 使用 Updates 而不是 Save，避免零值问题
	return r.db.Model(card).Where("id = ?", card.ID).Updates(card).Error
}

func (r *CardRepository) GetByID(id uint) (*models.Card, error) {
	var card models.Card
	err := r.db.First(&card, id).Error
	if err != nil {
		return nil, err
	}
	return &card, nil
}

func (r *CardRepository) GetByCode(code string) (*models.Card, error) {
	var card models.Card
	err := r.db.Where("card_code = ?", code).First(&card).Error
	if err != nil {
		return nil, err
	}
	return &card, nil
}

func (r *CardRepository) List(offset, limit int, filters map[string]interface{}) ([]*models.Card, int64, error) {
	var cards []*models.Card
	var total int64

	query := r.db.Model(&models.Card{})

	// Apply filters
	if status, ok := filters["status"]; ok {
		query = query.Where("status = ?", status)
	}
	if priceID, ok := filters["price_id"]; ok {
		query = query.Where("price_id = ?", priceID)
	}
	if batchID, ok := filters["batch_id"]; ok {
		if batchID == 0 {
			query = query.Where("batch_id IS NULL")
		} else {
			query = query.Where("batch_id = ?", batchID)
		}
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get data
	// 优先显示未使用的卡片，然后按创建时间倒序
	err := query.Offset(offset).Limit(limit).Order("status ASC, created_at DESC").Find(&cards).Error
	if err != nil {
		return nil, 0, err
	}

	return cards, total, nil
}

func (r *CardRepository) Delete(id uint) error {
	return r.db.Delete(&models.Card{}, id).Error
}

func (r *CardRepository) LogUsage(log *models.CardUsageLog) error {
	return r.db.Create(log).Error
}

func (r *CardRepository) GetUsageLogs(cardID uint, offset, limit int) ([]*models.CardUsageLog, int64, error) {
	var logs []*models.CardUsageLog
	var total int64

	query := r.db.Model(&models.CardUsageLog{}).Where("card_id = ?", cardID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&logs).Error
	if err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

// GetAvailableByPriceID 根据价格ID获取可用的卡片
func (r *CardRepository) GetAvailableByPriceID(priceID int64) (*models.Card, error) {
	var card models.Card
	err := r.db.Where("price_id = ? AND status = 0 AND expired_at > NOW()", priceID).
		First(&card).Error
	if err != nil {
		return nil, err
	}
	return &card, nil
}

// MarkAsUsed 标记卡片为已使用
func (r *CardRepository) MarkAsUsed(id uint, orderID uint) error {
	return r.db.Model(&models.Card{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":   1,
			"used_at":  gorm.Expr("NOW()"),
			"order_id": orderID,
		}).Error
}

// GetAvailableCountByPriceID 获取某个价格下可用卡片数量
func (r *CardRepository) GetAvailableCountByPriceID(priceID int64) (int64, error) {
	var count int64
	err := r.db.Model(&models.Card{}).
		Where("price_id = ? AND status = 0 AND expired_at > NOW()", priceID).
		Count(&count).Error
	return count, err
}

// GetCardsByBatchID 获取批次下的所有卡片
func (r *CardRepository) GetCardsByBatchID(batchID uint, offset, limit int) ([]*models.Card, int64, error) {
	var cards []*models.Card
	var total int64

	query := r.db.Model(&models.Card{}).Where("batch_id = ?", batchID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&cards).Error
	if err != nil {
		return nil, 0, err
	}

	return cards, total, nil
}

// GetCardsByPriceID 获取价格下的所有卡片
func (r *CardRepository) GetCardsByPriceID(priceID int64) ([]*models.Card, error) {
	var cards []*models.Card
	err := r.db.Where("price_id = ?", priceID).Find(&cards).Error
	if err != nil {
		return nil, err
	}
	return cards, nil
}


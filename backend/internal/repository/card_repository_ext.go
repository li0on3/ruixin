package repository

import (
	"backend/internal/models"
)

// GetAnyAvailable 获取任意一张可用的卡片
func (r *CardRepository) GetAnyAvailable() (*models.Card, error) {
	var card models.Card
	err := r.db.Where("status = ?", 0).First(&card).Error
	if err != nil {
		return nil, err
	}
	return &card, nil
}
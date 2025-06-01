package repository

import (
	"backend/internal/models"
	"gorm.io/gorm"
	"time"
)

type TransactionRepository struct {
	BaseRepository
}

func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{
		BaseRepository: BaseRepository{db: db},
	}
}

func (r *TransactionRepository) Create(tx *models.Transaction) error {
	return r.db.Create(tx).Error
}

func (r *TransactionRepository) GetByID(id int64) (*models.Transaction, error) {
	var tx models.Transaction
	err := r.db.Preload("Distributor").Preload("Creator").First(&tx, id).Error
	if err != nil {
		return nil, err
	}
	return &tx, nil
}

func (r *TransactionRepository) ListByDistributor(distributorID int64, offset, limit int) ([]*models.Transaction, int64, error) {
	var transactions []*models.Transaction
	var total int64

	query := r.db.Model(&models.Transaction{}).Where("distributor_id = ?", distributorID)

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Preload("Creator").
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&transactions).Error

	return transactions, total, err
}

func (r *TransactionRepository) ListByFilters(filters map[string]interface{}, offset, limit int) ([]*models.Transaction, int64, error) {
	var transactions []*models.Transaction
	var total int64

	query := r.db.Model(&models.Transaction{})

	// Apply filters
	if distributorID, ok := filters["distributor_id"].(int64); ok && distributorID > 0 {
		query = query.Where("distributor_id = ?", distributorID)
	}
	if txType, ok := filters["type"].(models.TransactionType); ok && txType > 0 {
		query = query.Where("type = ?", txType)
	}
	if startTime, ok := filters["start_time"].(time.Time); ok {
		query = query.Where("created_at >= ?", startTime)
	}
	if endTime, ok := filters["end_time"].(time.Time); ok {
		query = query.Where("created_at <= ?", endTime)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Preload("Distributor").Preload("Creator").
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&transactions).Error

	return transactions, total, err
}

func (r *TransactionRepository) GetDistributorBalance(distributorID int64) (float64, float64, error) {
	var distributor models.Distributor
	err := r.db.Select("balance, frozen_amount").First(&distributor, distributorID).Error
	if err != nil {
		return 0, 0, err
	}
	return distributor.Balance, distributor.FrozenAmount, nil
}

func (r *TransactionRepository) UpdateDistributorBalance(tx *gorm.DB, distributorID int64, balance, frozenAmount float64) error {
	return tx.Model(&models.Distributor{}).
		Where("id = ?", distributorID).
		Updates(map[string]interface{}{
			"balance":       balance,
			"frozen_amount": frozenAmount,
		}).Error
}

func (r *TransactionRepository) GetMonthlyStatistics(distributorID int64, year, month int) (map[string]interface{}, error) {
	var result struct {
		TotalRecharge float64
		TotalConsume  float64
		TotalWithdraw float64
		TotalRefund   float64
		TotalAdjust   float64
	}

	startTime := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)
	endTime := startTime.AddDate(0, 1, 0)

	err := r.db.Model(&models.Transaction{}).
		Where("distributor_id = ? AND created_at >= ? AND created_at < ?", distributorID, startTime, endTime).
		Select("SUM(CASE WHEN type = ? THEN amount ELSE 0 END) as total_recharge", models.TransactionTypeRecharge).
		Select("SUM(CASE WHEN type = ? THEN amount ELSE 0 END) as total_consume", models.TransactionTypeConsume).
		Select("SUM(CASE WHEN type = ? THEN amount ELSE 0 END) as total_withdraw", models.TransactionTypeWithdraw).
		Select("SUM(CASE WHEN type = ? THEN amount ELSE 0 END) as total_refund", models.TransactionTypeRefund).
		Select("SUM(CASE WHEN type = ? THEN amount ELSE 0 END) as total_adjust", models.TransactionTypeAdjust).
		Scan(&result).Error

	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"total_recharge": result.TotalRecharge,
		"total_consume":  result.TotalConsume,
		"total_withdraw": result.TotalWithdraw,
		"total_refund":   result.TotalRefund,
		"total_adjust":   result.TotalAdjust,
		"net_change":     result.TotalRecharge - result.TotalConsume - result.TotalWithdraw + result.TotalRefund + result.TotalAdjust,
	}, nil
}

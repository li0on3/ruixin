package repository

import (
	"backend/internal/models"
	"gorm.io/gorm"
	"time"
)

type WithdrawalRepository struct {
	BaseRepository
}

func NewWithdrawalRepository(db *gorm.DB) *WithdrawalRepository {
	return &WithdrawalRepository{
		BaseRepository: BaseRepository{db: db},
	}
}

func (r *WithdrawalRepository) Create(withdrawal *models.Withdrawal) error {
	return r.db.Create(withdrawal).Error
}

func (r *WithdrawalRepository) GetByID(id int64) (*models.Withdrawal, error) {
	var withdrawal models.Withdrawal
	err := r.db.Preload("Distributor").Preload("Processor").First(&withdrawal, id).Error
	if err != nil {
		return nil, err
	}
	return &withdrawal, nil
}

func (r *WithdrawalRepository) Update(withdrawal *models.Withdrawal) error {
	return r.db.Save(withdrawal).Error
}

func (r *WithdrawalRepository) ListByDistributor(distributorID int64, offset, limit int) ([]*models.Withdrawal, int64, error) {
	var withdrawals []*models.Withdrawal
	var total int64

	query := r.db.Model(&models.Withdrawal{}).Where("distributor_id = ?", distributorID)

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Preload("Processor").
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&withdrawals).Error

	return withdrawals, total, err
}

func (r *WithdrawalRepository) ListPending(offset, limit int) ([]*models.Withdrawal, int64, error) {
	var withdrawals []*models.Withdrawal
	var total int64

	query := r.db.Model(&models.Withdrawal{}).Where("status = ?", models.WithdrawalStatusPending)

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Preload("Distributor").
		Order("created_at ASC").
		Offset(offset).
		Limit(limit).
		Find(&withdrawals).Error

	return withdrawals, total, err
}

func (r *WithdrawalRepository) ListByFilters(filters map[string]interface{}, offset, limit int) ([]*models.Withdrawal, int64, error) {
	var withdrawals []*models.Withdrawal
	var total int64

	query := r.db.Model(&models.Withdrawal{})

	// Apply filters
	if distributorID, ok := filters["distributor_id"].(int64); ok && distributorID > 0 {
		query = query.Where("distributor_id = ?", distributorID)
	}
	if status, ok := filters["status"].(models.WithdrawalStatus); ok {
		query = query.Where("status = ?", status)
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

	err = query.Preload("Distributor").Preload("Processor").
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&withdrawals).Error

	return withdrawals, total, err
}

func (r *WithdrawalRepository) GetPendingCount() (int64, error) {
	var count int64
	err := r.db.Model(&models.Withdrawal{}).Where("status = ?", models.WithdrawalStatusPending).Count(&count).Error
	return count, err
}

func (r *WithdrawalRepository) GetPendingAmount() (float64, error) {
	var total float64
	err := r.db.Model(&models.Withdrawal{}).
		Where("status = ?", models.WithdrawalStatusPending).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&total).Error
	return total, err
}

func (r *WithdrawalRepository) ProcessWithdrawal(tx *gorm.DB, withdrawalID int64, processorID int64, status models.WithdrawalStatus, rejectReason string) error {
	now := time.Now()
	updates := map[string]interface{}{
		"status":       status,
		"processed_at": now,
		"processed_by": processorID,
	}

	if status == models.WithdrawalStatusRejected && rejectReason != "" {
		updates["reject_reason"] = rejectReason
	}

	return tx.Model(&models.Withdrawal{}).
		Where("id = ? AND status = ?", withdrawalID, models.WithdrawalStatusPending).
		Updates(updates).Error
}

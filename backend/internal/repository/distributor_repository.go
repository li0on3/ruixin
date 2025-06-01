package repository

import (
	"backend/internal/models"
	"gorm.io/gorm"
)

type DistributorRepository struct {
	*BaseRepository
}

func NewDistributorRepository(db *gorm.DB) *DistributorRepository {
	return &DistributorRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

func (r *DistributorRepository) Create(distributor *models.Distributor) error {
	return r.db.Create(distributor).Error
}

func (r *DistributorRepository) Update(distributor *models.Distributor) error {
	return r.db.Save(distributor).Error
}

func (r *DistributorRepository) GetByID(id uint) (*models.Distributor, error) {
	var distributor models.Distributor
	err := r.db.First(&distributor, id).Error
	if err != nil {
		return nil, err
	}
	return &distributor, nil
}

func (r *DistributorRepository) GetByAPIKey(apiKey string) (*models.Distributor, error) {
	var distributor models.Distributor
	err := r.db.Where("api_key = ?", apiKey).First(&distributor).Error
	if err != nil {
		return nil, err
	}
	return &distributor, nil
}

func (r *DistributorRepository) GetByEmail(email string) (*models.Distributor, error) {
	var distributor models.Distributor
	err := r.db.Where("email = ?", email).First(&distributor).Error
	if err != nil {
		return nil, err
	}
	return &distributor, nil
}

func (r *DistributorRepository) List(offset, limit int, filters map[string]interface{}) ([]*models.Distributor, int64, error) {
	var distributors []*models.Distributor
	var total int64

	query := r.db.Model(&models.Distributor{})

	// Apply filters
	if status, ok := filters["status"]; ok {
		query = query.Where("status = ?", status)
	}
	if name, ok := filters["name"]; ok {
		query = query.Where("name LIKE ?", "%"+name.(string)+"%")
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get data
	err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&distributors).Error
	if err != nil {
		return nil, 0, err
	}

	return distributors, total, nil
}

func (r *DistributorRepository) Delete(id uint) error {
	return r.db.Delete(&models.Distributor{}, id).Error
}

func (r *DistributorRepository) UpdateBalance(id uint, amount float64) error {
	return r.db.Model(&models.Distributor{}).Where("id = ?", id).
		Update("balance", gorm.Expr("balance + ?", amount)).Error
}

func (r *DistributorRepository) IncrementOrderCount(id uint) error {
	return r.db.Model(&models.Distributor{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"total_orders": gorm.Expr("total_orders + ?", 1),
		}).Error
}

func (r *DistributorRepository) LogAPICall(log *models.DistributorAPILog) error {
	return r.db.Create(log).Error
}

func (r *DistributorRepository) GetAPILogs(distributorID uint, offset, limit int) ([]*models.DistributorAPILog, int64, error) {
	var logs []*models.DistributorAPILog
	var total int64

	query := r.db.Model(&models.DistributorAPILog{}).Where("distributor_id = ?", distributorID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&logs).Error
	if err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

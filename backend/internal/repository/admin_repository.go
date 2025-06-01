package repository

import (
	"backend/internal/models"
	"gorm.io/gorm"
)

type AdminRepository struct {
	*BaseRepository
}

func NewAdminRepository(db *gorm.DB) *AdminRepository {
	return &AdminRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

func (r *AdminRepository) Create(admin *models.Admin) error {
	return r.db.Create(admin).Error
}

func (r *AdminRepository) Update(admin *models.Admin) error {
	return r.db.Save(admin).Error
}

func (r *AdminRepository) GetByID(id uint) (*models.Admin, error) {
	var admin models.Admin
	err := r.db.First(&admin, id).Error
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

func (r *AdminRepository) GetByUsername(username string) (*models.Admin, error) {
	var admin models.Admin
	err := r.db.Where("username = ?", username).First(&admin).Error
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

func (r *AdminRepository) GetByEmail(email string) (*models.Admin, error) {
	var admin models.Admin
	err := r.db.Where("email = ?", email).First(&admin).Error
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

func (r *AdminRepository) List(offset, limit int, filters map[string]interface{}) ([]*models.Admin, int64, error) {
	var admins []*models.Admin
	var total int64

	query := r.db.Model(&models.Admin{})

	// Apply filters
	if role, ok := filters["role"]; ok {
		query = query.Where("role = ?", role)
	}
	if status, ok := filters["status"]; ok {
		query = query.Where("status = ?", status)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get data
	err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&admins).Error
	if err != nil {
		return nil, 0, err
	}

	return admins, total, nil
}

func (r *AdminRepository) Delete(id uint) error {
	return r.db.Delete(&models.Admin{}, id).Error
}

func (r *AdminRepository) UpdateLoginInfo(id uint, ip string) error {
	return r.db.Model(&models.Admin{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"last_login_at": gorm.Expr("NOW()"),
			"last_login_ip": ip,
		}).Error
}

func (r *AdminRepository) LogOperation(log *models.AdminOperationLog) error {
	return r.db.Create(log).Error
}

func (r *AdminRepository) GetOperationLogs(adminID uint, offset, limit int) ([]*models.AdminOperationLog, int64, error) {
	var logs []*models.AdminOperationLog
	var total int64

	query := r.db.Model(&models.AdminOperationLog{})
	if adminID > 0 {
		query = query.Where("admin_id = ?", adminID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&logs).Error
	if err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

// DistributorAPILog represents API call log entry
type DistributorAPILog struct {
	models.DistributorAPILog
}

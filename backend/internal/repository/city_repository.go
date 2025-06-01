package repository

import (
	"backend/internal/models"
	"gorm.io/gorm"
)

type CityRepository struct {
	db *gorm.DB
}

func NewCityRepository(db *gorm.DB) *CityRepository {
	return &CityRepository{db: db}
}

// GetByCityID 根据城市ID获取城市信息
func (r *CityRepository) GetByCityID(cityID int) (*models.City, error) {
	var city models.City
	err := r.db.Where("city_id = ?", cityID).First(&city).Error
	return &city, err
}

// GetByCityName 根据城市名称获取城市信息
func (r *CityRepository) GetByCityName(cityName string) (*models.City, error) {
	var city models.City
	err := r.db.Where("city_name = ?", cityName).First(&city).Error
	return &city, err
}

// List 获取城市列表
func (r *CityRepository) List() ([]*models.City, error) {
	var cities []*models.City
	err := r.db.Order("city_id").Find(&cities).Error
	return cities, err
}

// CreateOrUpdate 创建或更新城市
func (r *CityRepository) CreateOrUpdate(city *models.City) error {
	return r.db.Save(city).Error
}

// BatchCreateOrUpdate 批量创建或更新城市
func (r *CityRepository) BatchCreateOrUpdate(cities []*models.City) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, city := range cities {
			if err := tx.Save(city).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// SearchByName 根据名称模糊搜索城市
func (r *CityRepository) SearchByName(keyword string) ([]*models.City, error) {
	var cities []*models.City
	err := r.db.Where("city_name LIKE ?", "%"+keyword+"%").Find(&cities).Error
	return cities, err
}
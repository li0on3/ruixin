package repository

import (
	"backend/internal/models"
	"gorm.io/gorm"
)

type LuckinConfigRepository struct {
	BaseRepository
}

func NewLuckinConfigRepository(db *gorm.DB) *LuckinConfigRepository {
	return &LuckinConfigRepository{
		BaseRepository: BaseRepository{db: db},
	}
}

// LuckinPrice methods
func (r *LuckinConfigRepository) CreatePrice(price *models.LuckinPrice) error {
	return r.db.Create(price).Error
}

func (r *LuckinConfigRepository) GetPriceByID(id int64) (*models.LuckinPrice, error) {
	var price models.LuckinPrice
	err := r.db.Preload("Creator").First(&price, id).Error
	if err != nil {
		return nil, err
	}
	return &price, nil
}

func (r *LuckinConfigRepository) GetPriceByPriceID(priceID string) (*models.LuckinPrice, error) {
	var price models.LuckinPrice
	err := r.db.Where("price_id = ?", priceID).First(&price).Error
	if err != nil {
		return nil, err
	}
	return &price, nil
}

func (r *LuckinConfigRepository) UpdatePrice(price *models.LuckinPrice) error {
	return r.db.Save(price).Error
}

func (r *LuckinConfigRepository) DeletePrice(id int64) error {
	return r.db.Delete(&models.LuckinPrice{}, id).Error
}

func (r *LuckinConfigRepository) ListPrices(offset, limit int) ([]*models.LuckinPrice, int64, error) {
	var prices []*models.LuckinPrice
	var total int64

	err := r.db.Model(&models.LuckinPrice{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Preload("Creator").
		Order("price_value ASC").
		Offset(offset).
		Limit(limit).
		Find(&prices).Error

	return prices, total, err
}

func (r *LuckinConfigRepository) GetActivePrices() ([]*models.LuckinPrice, error) {
	var prices []*models.LuckinPrice
	err := r.db.Where("status = ?", 1).Order("price_value ASC").Find(&prices).Error
	return prices, err
}

// LuckinProduct methods
func (r *LuckinConfigRepository) CreateProduct(product *models.LuckinProduct) error {
	return r.db.Create(product).Error
}

func (r *LuckinConfigRepository) GetProductByID(id int64) (*models.LuckinProduct, error) {
	var product models.LuckinProduct
	err := r.db.Preload("Creator").First(&product, id).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *LuckinConfigRepository) GetProductByProductID(productID string) (*models.LuckinProduct, error) {
	var product models.LuckinProduct
	err := r.db.Where("product_id = ?", productID).First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *LuckinConfigRepository) UpdateProduct(product *models.LuckinProduct) error {
	return r.db.Save(product).Error
}

func (r *LuckinConfigRepository) DeleteProduct(id int64) error {
	return r.db.Delete(&models.LuckinProduct{}, id).Error
}

func (r *LuckinConfigRepository) ListProducts(filters map[string]interface{}, offset, limit int) ([]*models.LuckinProduct, int64, error) {
	var products []*models.LuckinProduct
	var total int64

	query := r.db.Model(&models.LuckinProduct{})

	// Apply filters
	if name, ok := filters["name"].(string); ok && name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if category, ok := filters["category"].(string); ok && category != "" {
		query = query.Where("category = ?", category)
	}
	if status, ok := filters["status"].(int); ok && status >= 0 {
		query = query.Where("status = ?", status)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Preload("Creator").
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&products).Error

	return products, total, err
}

func (r *LuckinConfigRepository) GetActiveProducts() ([]*models.LuckinProduct, error) {
	var products []*models.LuckinProduct
	err := r.db.Where("status = ?", 1).Order("name ASC").Find(&products).Error
	return products, err
}

func (r *LuckinConfigRepository) GetProductCategories() ([]string, error) {
	var categories []string
	err := r.db.Model(&models.LuckinProduct{}).
		Where("category IS NOT NULL AND category != ''").
		Distinct("category").
		Pluck("category", &categories).Error
	return categories, err
}

// CategoryBinding methods
func (r *LuckinConfigRepository) CreateBinding(binding *models.CategoryBinding) error {
	return r.db.Create(binding).Error
}

func (r *LuckinConfigRepository) DeleteBinding(id int64) error {
	return r.db.Delete(&models.CategoryBinding{}, id).Error
}

func (r *LuckinConfigRepository) ListBindingsByCategory(categoryID int64) ([]*models.CategoryBinding, error) {
	var bindings []*models.CategoryBinding
	err := r.db.Where("category_id = ?", categoryID).
		Preload("Creator").
		Order("target_type ASC, priority DESC").
		Find(&bindings).Error
	return bindings, err
}

func (r *LuckinConfigRepository) ListBindingsByTarget(targetType, targetID string) ([]*models.CategoryBinding, error) {
	var bindings []*models.CategoryBinding
	err := r.db.Where("target_type = ? AND target_id = ?", targetType, targetID).
		Preload("Category").
		Preload("Creator").
		Find(&bindings).Error
	return bindings, err
}

func (r *LuckinConfigRepository) UpdateBindingPriority(id int64, priority int) error {
	return r.db.Model(&models.CategoryBinding{}).Where("id = ?", id).Update("priority", priority).Error
}

// 批量导入产品
func (r *LuckinConfigRepository) BatchImportProducts(products []*models.LuckinProduct) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, product := range products {
			// 尝试更新，如果不存在则创建
			var existing models.LuckinProduct
			err := tx.Where("product_id = ?", product.ProductID).First(&existing).Error
			if err == gorm.ErrRecordNotFound {
				// 创建新产品
				if err := tx.Create(product).Error; err != nil {
					return err
				}
			} else if err == nil {
				// 更新现有产品
				existing.Name = product.Name
				existing.Description = product.Description
				existing.Category = product.Category
				existing.ImageURL = product.ImageURL
				existing.Status = product.Status
				if err := tx.Save(&existing).Error; err != nil {
					return err
				}
			} else {
				return err
			}
		}
		return nil
	})
}

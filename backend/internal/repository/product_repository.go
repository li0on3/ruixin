package repository

import (
	"backend/internal/models"
	"gorm.io/gorm"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

// GetByCode 根据商品代码获取商品
func (r *ProductRepository) GetByCode(code string) (*models.Product, error) {
	var product models.Product
	err := r.db.Where("goods_code = ?", code).
		Preload("SKUs").
		Preload("SKUs.Mapping").
		Preload("SKUs.Specs").
		Preload("SKUs.Specs.Options").
		First(&product).Error
	return &product, err
}

// GetPriceMappings 获取商品的价格映射
func (r *ProductRepository) GetPriceMappings(productCode string, skuCode string) ([]*models.ProductPriceMapping, error) {
	var mappings []*models.ProductPriceMapping
	
	query := r.db.Where("product_code = ? AND status = 1", productCode)
	
	// 如果指定了SKU，优先查找SKU特定的映射
	if skuCode != "" {
		// 先查找SKU特定的映射
		var skuMappings []*models.ProductPriceMapping
		err := query.Where("sku_code = ?", skuCode).
			Order("priority DESC").
			Preload("Price").
			Find(&skuMappings).Error
		if err == nil && len(skuMappings) > 0 {
			return skuMappings, nil
		}
	}
	
	// 查找通用映射（sku_code为空）
	err := query.Where("sku_code = '' OR sku_code IS NULL").
		Order("priority DESC").
		Preload("Price").
		Find(&mappings).Error
	
	return mappings, err
}

// CreateProduct 创建商品
func (r *ProductRepository) CreateProduct(product *models.Product) error {
	return r.db.Create(product).Error
}

// UpdateProduct 更新商品
func (r *ProductRepository) UpdateProduct(product *models.Product) error {
	return r.db.Save(product).Error
}

// CreatePriceMapping 创建价格映射
func (r *ProductRepository) CreatePriceMapping(mapping *models.ProductPriceMapping) error {
	return r.db.Create(mapping).Error
}

// DeletePriceMapping 删除价格映射
func (r *ProductRepository) DeletePriceMapping(id uint) error {
	return r.db.Delete(&models.ProductPriceMapping{}, id).Error
}

// GetProductList 获取商品列表
func (r *ProductRepository) GetProductList(offset, limit int, search string) ([]*models.Product, int64, error) {
	var products []*models.Product
	var total int64
	
	query := r.db.Model(&models.Product{})
	
	if search != "" {
		query = query.Where("goods_code LIKE ? OR goods_name LIKE ?", "%"+search+"%", "%"+search+"%")
	}
	
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	
	err = query.Offset(offset).Limit(limit).
		Order("id DESC").
		Preload("SKUs").
		Preload("SKUs.Mapping").
		Preload("SKUs.Specs").
		Preload("SKUs.Specs.Options").
		Find(&products).Error
	
	return products, total, err
}

// BatchCreateProducts 批量创建商品
func (r *ProductRepository) BatchCreateProducts(products []*models.Product) error {
	return r.db.CreateInBatches(products, 100).Error
}

// GetSKUByCode 根据SKU代码获取SKU信息
func (r *ProductRepository) GetSKUByCode(skuCode string) (*models.ProductSKU, error) {
	var sku models.ProductSKU
	err := r.db.Where("sku_code = ?", skuCode).
		Preload("Product").
		Preload("Specs").
		Preload("Specs.Options").
		First(&sku).Error
	return &sku, err
}

// SearchProducts 搜索商品
func (r *ProductRepository) SearchProducts(keyword string) ([]*models.Product, error) {
	var products []*models.Product
	err := r.db.Where("goods_code LIKE ? OR goods_name LIKE ?", "%"+keyword+"%", "%"+keyword+"%").
		Limit(20).
		Find(&products).Error
	return products, err
}

// GetProductsByCodes 根据商品代码批量获取商品
func (r *ProductRepository) GetProductsByCodes(codes []string) ([]*models.Product, error) {
	var products []*models.Product
	if len(codes) == 0 {
		return products, nil
	}
	
	err := r.db.Where("goods_code IN ?", codes).Find(&products).Error
	return products, err
}

// GetProductsByCardIDs 根据卡片ID获取绑定的商品（去重）
func (r *ProductRepository) GetProductsByCardIDs(cardIDs []uint) ([]*models.Product, error) {
	var products []*models.Product
	if len(cardIDs) == 0 {
		return products, nil
	}
	
	err := r.db.Distinct("products.*").
		Joins("JOIN card_product_bindings ON card_product_bindings.product_id = products.id").
		Where("card_product_bindings.card_id IN ? AND card_product_bindings.is_active = ?", cardIDs, true).
		Order("products.goods_name").
		Find(&products).Error
		
	return products, err
}
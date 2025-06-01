package services

import (
	"errors"
	"fmt"

	"backend/internal/models"
	"backend/internal/repository"
	"gorm.io/gorm"
)

type LuckinConfigService struct {
	db         *gorm.DB
	configRepo *repository.LuckinConfigRepository
	cardRepo   *repository.CardRepository
}

func NewLuckinConfigService(db *gorm.DB) *LuckinConfigService {
	return &LuckinConfigService{
		db:         db,
		configRepo: repository.NewLuckinConfigRepository(db),
		cardRepo:   repository.NewCardRepository(db),
	}
}

// Price Management
func (s *LuckinConfigService) CreatePrice(price *models.LuckinPrice) error {
	// 检查price_id是否已存在
	existing, _ := s.configRepo.GetPriceByPriceID(price.PriceCode)
	if existing != nil {
		return errors.New("价格ID已存在")
	}
	return s.configRepo.CreatePrice(price)
}

func (s *LuckinConfigService) UpdatePrice(id int64, updates map[string]interface{}) error {
	price, err := s.configRepo.GetPriceByID(id)
	if err != nil {
		return err
	}

	// 如果要更新price_id，检查是否重复
	if newPriceID, ok := updates["price_id"].(string); ok && newPriceID != price.PriceCode {
		existing, _ := s.configRepo.GetPriceByPriceID(newPriceID)
		if existing != nil {
			return errors.New("价格ID已存在")
		}
		price.PriceCode = newPriceID
	}

	if priceValue, ok := updates["price_value"].(float64); ok {
		price.PriceValue = priceValue
	}
	if status, ok := updates["status"].(int); ok {
		price.Status = status
	}
	// status might come as float64 from JSON
	if status, ok := updates["status"].(float64); ok {
		price.Status = int(status)
	}

	return s.configRepo.UpdatePrice(price)
}

func (s *LuckinConfigService) DeletePrice(id int64) error {
	// 检查是否有卡片绑定了这个价格
	var bindings []*models.CategoryBinding
	bindings, _ = s.configRepo.ListBindingsByTarget("price", fmt.Sprintf("%d", id))
	if len(bindings) > 0 {
		return errors.New("该价格已被卡片种类绑定，无法删除")
	}
	return s.configRepo.DeletePrice(id)
}

func (s *LuckinConfigService) GetPriceList(page, pageSize int) ([]*models.LuckinPrice, int64, error) {
	offset := (page - 1) * pageSize
	return s.configRepo.ListPrices(offset, pageSize)
}

func (s *LuckinConfigService) GetActivePrices() ([]*models.LuckinPrice, error) {
	return s.configRepo.GetActivePrices()
}

func (s *LuckinConfigService) GetPriceByID(id int64) (*models.LuckinPrice, error) {
	return s.configRepo.GetPriceByID(id)
}

// Product Management
func (s *LuckinConfigService) CreateProduct(product *models.LuckinProduct) error {
	// 检查product_id是否已存在
	existing, _ := s.configRepo.GetProductByProductID(product.ProductID)
	if existing != nil {
		return errors.New("产品ID已存在")
	}
	return s.configRepo.CreateProduct(product)
}

func (s *LuckinConfigService) UpdateProduct(id int64, updates map[string]interface{}) error {
	product, err := s.configRepo.GetProductByID(id)
	if err != nil {
		return err
	}

	// 如果要更新product_id，检查是否重复
	if newProductID, ok := updates["product_id"].(string); ok && newProductID != product.ProductID {
		existing, _ := s.configRepo.GetProductByProductID(newProductID)
		if existing != nil {
			return errors.New("产品ID已存在")
		}
		product.ProductID = newProductID
	}

	if name, ok := updates["name"].(string); ok {
		product.Name = name
	}
	if description, ok := updates["description"].(string); ok {
		product.Description = description
	}
	if category, ok := updates["category"].(string); ok {
		product.Category = category
	}
	if imageURL, ok := updates["image_url"].(string); ok {
		product.ImageURL = imageURL
	}
	if status, ok := updates["status"].(int); ok {
		product.Status = status
	}

	return s.configRepo.UpdateProduct(product)
}

func (s *LuckinConfigService) DeleteProduct(id int64) error {
	// 检查是否有卡片绑定了这个产品
	var bindings []*models.CategoryBinding
	bindings, _ = s.configRepo.ListBindingsByTarget("product", fmt.Sprintf("%d", id))
	if len(bindings) > 0 {
		return errors.New("该产品已被卡片种类绑定，无法删除")
	}
	return s.configRepo.DeleteProduct(id)
}

func (s *LuckinConfigService) GetProductList(filters map[string]interface{}, page, pageSize int) ([]*models.LuckinProduct, int64, error) {
	offset := (page - 1) * pageSize
	return s.configRepo.ListProducts(filters, offset, pageSize)
}

func (s *LuckinConfigService) GetActiveProducts() ([]*models.LuckinProduct, error) {
	return s.configRepo.GetActiveProducts()
}

func (s *LuckinConfigService) GetProductCategories() ([]string, error) {
	return s.configRepo.GetProductCategories()
}

// Category Binding Management
func (s *LuckinConfigService) CreateBinding(categoryID int64, targetType, targetID string, priority int, createdBy int64) error {
	// 验证目标是否存在
	switch targetType {
	case "price":
		price, err := s.configRepo.GetPriceByPriceID(targetID)
		if err != nil {
			return errors.New("价格不存在")
		}
		if price.Status != 1 {
			return errors.New("价格已禁用")
		}
	case "product":
		product, err := s.configRepo.GetProductByProductID(targetID)
		if err != nil {
			return errors.New("产品不存在")
		}
		if product.Status != 1 {
			return errors.New("产品已禁用")
		}
	default:
		return errors.New("无效的目标类型")
	}

	// 验证卡片种类是否存在
	_, err := s.cardRepo.GetByID(uint(categoryID))
	if err != nil {
		return errors.New("卡片种类不存在")
	}

	binding := &models.CategoryBinding{
		CategoryID: categoryID,
		TargetType: targetType,
		TargetID:   targetID,
		Priority:   priority,
		CreatedBy:  createdBy,
	}

	return s.configRepo.CreateBinding(binding)
}

func (s *LuckinConfigService) DeleteBinding(id int64) error {
	return s.configRepo.DeleteBinding(id)
}

func (s *LuckinConfigService) GetCategoryBindings(categoryID int64) ([]*models.CategoryBinding, error) {
	bindings, err := s.configRepo.ListBindingsByCategory(categoryID)
	if err != nil {
		return nil, err
	}

	// 填充关联的价格和产品信息
	for _, binding := range bindings {
		switch binding.TargetType {
		case "price":
			if price, err := s.configRepo.GetPriceByPriceID(binding.TargetID); err == nil {
				binding.TargetID = fmt.Sprintf("%s (¥%.2f)", price.PriceCode, price.PriceValue)
			}
		case "product":
			if product, err := s.configRepo.GetProductByProductID(binding.TargetID); err == nil {
				binding.TargetID = fmt.Sprintf("%s (%s)", product.ProductID, product.Name)
			}
		}
	}

	return bindings, nil
}

func (s *LuckinConfigService) UpdateBindingPriority(id int64, priority int) error {
	return s.configRepo.UpdateBindingPriority(id, priority)
}

// Batch Import
func (s *LuckinConfigService) BatchImportProducts(products []*models.LuckinProduct) error {
	if len(products) == 0 {
		return errors.New("没有要导入的产品")
	}
	return s.configRepo.BatchImportProducts(products)
}

// 同步瑞幸产品数据（从API获取并更新本地数据）
func (s *LuckinConfigService) SyncProductsFromAPI() error {
	// TODO: 调用瑞幸API获取最新产品列表
	// 这里需要实现从瑞幸API获取产品数据的逻辑
	// 暂时返回未实现错误
	return errors.New("产品同步功能尚未实现")
}


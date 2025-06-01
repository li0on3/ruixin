package services

import (
	"backend/internal/models"
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"strings"
	"time"
)

type ProductMappingService struct {
	db *gorm.DB
}

func NewProductMappingService(db *gorm.DB) *ProductMappingService {
	return &ProductMappingService{db: db}
}

// MapProductRequest 映射商品请求
type MapProductRequest struct {
	ProductName string                 `json:"product"`    // 商品名称
	Specs       map[string]string      `json:"specs"`      // 规格映射
	CardID      uint                   `json:"card_id,omitempty"` // 卡片ID（用于限定商品范围）
}

// MapProductResult 映射结果
type MapProductResult struct {
	ProductCode string                    `json:"product_code"`
	SKUCode     string                    `json:"sku_code"`
	GoodsJSON   []map[string]interface{}  `json:"goods_json"`
}

// MapProduct 将中文商品描述映射到具体的商品代码和SKU
func (s *ProductMappingService) MapProduct(req MapProductRequest) (*MapProductResult, error) {
	// 1. 查找商品（支持别名）
	product, err := s.findProductByName(req.ProductName, req.CardID)
	if err != nil {
		return nil, fmt.Errorf("商品'%s'未找到", req.ProductName)
	}
	
	// 2. 标准化规格输入
	normalizedSpecs, err := s.normalizeSpecs(req.Specs)
	if err != nil {
		return nil, err
	}
	
	// 3. 检查必填规格
	if err := s.checkRequiredSpecs(product.ID, normalizedSpecs); err != nil {
		return nil, err
	}
	
	// 4. 构建中文描述查找SKU
	sku, err := s.findSKUBySpecs(product.ID, normalizedSpecs)
	if err != nil {
		return nil, err
	}
	
	// 5. 构建下单所需的goods JSON
	goodsJSON, err := s.buildGoodsJSON(product.GoodsCode, sku, normalizedSpecs)
	if err != nil {
		return nil, err
	}
	
	return &MapProductResult{
		ProductCode: product.GoodsCode,
		SKUCode:     sku.SKUCode,
		GoodsJSON:   goodsJSON,
	}, nil
}

// findProductByName 通过名称查找商品（支持别名，可限定卡片范围）
func (s *ProductMappingService) findProductByName(name string, cardID uint) (*models.Product, error) {
	var product models.Product
	
	// 如果指定了卡片ID，则只在该卡片绑定的商品中查找
	if cardID > 0 {
		// 先尝试直接匹配商品名（在卡片绑定的商品中）
		query := s.db.Table("products").
			Joins("JOIN card_product_bindings ON products.id = card_product_bindings.product_id").
			Where("card_product_bindings.card_id = ? AND card_product_bindings.is_active = ?", cardID, true).
			Where("products.goods_name = ?", name)
		
		err := query.First(&product).Error
		if err == nil {
			return &product, nil
		}
		
		// 尝试通过别名查找（在卡片绑定的商品中）
		var alias models.ProductAlias
		err = s.db.Where("alias_name = ?", name).First(&alias).Error
		if err == nil {
			// 检查该商品是否绑定到指定卡片
			query = s.db.Table("products").
				Joins("JOIN card_product_bindings ON products.id = card_product_bindings.product_id").
				Where("card_product_bindings.card_id = ? AND card_product_bindings.is_active = ?", cardID, true).
				Where("products.id = ?", alias.ProductID)
			
			err = query.First(&product).Error
			if err == nil {
				return &product, nil
			}
		}
		
		return nil, fmt.Errorf("商品'%s'未在该卡片的可用商品中找到", name)
	}
	
	// 如果没有指定卡片ID，则在所有商品中查找（原有逻辑）
	// 先尝试直接匹配商品名
	err := s.db.Where("goods_name = ?", name).First(&product).Error
	if err == nil {
		return &product, nil
	}
	
	// 尝试通过别名查找
	var alias models.ProductAlias
	err = s.db.Where("alias_name = ?", name).First(&alias).Error
	if err == nil {
		err = s.db.First(&product, alias.ProductID).Error
		if err == nil {
			return &product, nil
		}
	}
	
	return nil, gorm.ErrRecordNotFound
}

// normalizeSpecs 标准化规格输入
func (s *ProductMappingService) normalizeSpecs(specs map[string]string) (map[string]string, error) {
	normalized := make(map[string]string)
	
	for specType, value := range specs {
		// 标准化规格类型
		standardType := s.standardizeSpecType(specType)
		
		// 标准化规格值（通过别名）
		standardValue, err := s.standardizeSpecValue(standardType, value)
		if err != nil {
			return nil, fmt.Errorf("规格'%s'的值'%s'无效", specType, value)
		}
		
		normalized[standardType] = standardValue
	}
	
	return normalized, nil
}

// standardizeSpecType 标准化规格类型
func (s *ProductMappingService) standardizeSpecType(input string) string {
	// 规格类型映射
	typeMap := map[string]string{
		"杯型": "size",
		"大小": "size",
		"规格": "size",
		"温度": "temperature",
		"冷热": "temperature",
		"甜度": "sweetness",
		"糖度": "sweetness",
		"奶": "milk",
		"奶量": "milk",
		"口味": "flavor",
		"风味": "flavor",
		"咖啡豆": "bean",
	}
	
	if standardType, ok := typeMap[input]; ok {
		return standardType
	}
	return input
}

// standardizeSpecValue 标准化规格值
func (s *ProductMappingService) standardizeSpecValue(specType, value string) (string, error) {
	// 先尝试直接返回（已经是标准值）
	if s.isStandardValue(specType, value) {
		return value, nil
	}
	
	// 查找别名映射
	var alias models.SpecInputAlias
	err := s.db.Where("spec_type = ? AND alias_value = ?", specType, value).First(&alias).Error
	if err == nil {
		return alias.StandardValue, nil
	}
	
	// 没找到映射，返回错误
	return "", errors.New("未找到对应的标准值")
}

// isStandardValue 检查是否已经是标准值
func (s *ProductMappingService) isStandardValue(specType, value string) bool {
	// 这里可以维护一个标准值列表，或者从数据库查询
	standardValues := map[string][]string{
		"size": {"大杯", "中杯", "小杯"},
		"temperature": {"冰", "热", "常温"},
		"sweetness": {"标准甜", "少甜", "少少甜", "微甜", "不另外加糖"},
		"milk": {"双份奶", "单份奶", "无奶"},
	}
	
	if values, ok := standardValues[specType]; ok {
		for _, v := range values {
			if v == value {
				return true
			}
		}
	}
	return false
}

// checkRequiredSpecs 检查必填规格
func (s *ProductMappingService) checkRequiredSpecs(productID uint, specs map[string]string) error {
	var configs []models.ProductSpecConfig
	err := s.db.Where("product_id = ? AND is_required = ?", productID, true).Find(&configs).Error
	if err != nil {
		return err
	}
	
	var missingSpecs []string
	for _, config := range configs {
		if _, ok := specs[config.SpecType]; !ok {
			// 如果有默认值，自动填充
			if config.DefaultValue != "" {
				specs[config.SpecType] = config.DefaultValue
			} else {
				missingSpecs = append(missingSpecs, config.ChineseName)
			}
		}
	}
	
	if len(missingSpecs) > 0 {
		return fmt.Errorf("缺少必填规格: %s", strings.Join(missingSpecs, ", "))
	}
	
	return nil
}

// findSKUBySpecs 通过规格查找SKU
func (s *ProductMappingService) findSKUBySpecs(productID uint, specs map[string]string) (*models.ProductSKU, error) {
	// 获取商品的所有规格配置，按显示顺序排序
	var configs []models.ProductSpecConfig
	err := s.db.Where("product_id = ?", productID).Order("display_order").Find(&configs).Error
	if err != nil {
		return nil, err
	}
	
	// 构建中文描述
	var descParts []string
	for _, config := range configs {
		if value, ok := specs[config.SpecType]; ok {
			descParts = append(descParts, value)
		}
	}
	chineseDesc := strings.Join(descParts, "")
	
	// 查找对应的SKU映射
	var mapping models.ProductSKUMapping
	err = s.db.Where("product_id = ? AND chinese_desc = ?", productID, chineseDesc).First(&mapping).Error
	if err != nil {
		// SKU mapping not found - productID: %d, desc: %s
		return nil, fmt.Errorf("未找到对应的商品规格组合: %s", chineseDesc)
	}
	
	// 获取SKU信息
	var sku models.ProductSKU
	err = s.db.Where("sku_code = ?", mapping.SKUCode).First(&sku).Error
	if err != nil {
		return nil, err
	}
	
	return &sku, nil
}

// buildGoodsJSON 构建下单所需的goods JSON
func (s *ProductMappingService) buildGoodsJSON(goodsCode string, sku *models.ProductSKU, specs map[string]string) ([]map[string]interface{}, error) {
	// 获取SKU的所有规格信息
	var productSpecs []models.ProductSpec
	err := s.db.Where("sku_id = ?", sku.ID).
		Preload("Options").
		Order("id").Find(&productSpecs).Error
	if err != nil {
		return nil, err
	}
	
	// 构建规格数组
	var specsArray []map[string]interface{}
	for _, spec := range productSpecs {
		specType := s.getSpecType(spec.SpecsCode)
		selectedValue, ok := specs[specType]
		if !ok {
			continue
		}
		
		// 找到对应的选项
		for _, option := range spec.Options {
			if option.Name == selectedValue {
				specsArray = append(specsArray, map[string]interface{}{
					"specsCode": spec.SpecsCode,
					"code":      option.Code,
					"name":      option.Name,
				})
				break
			}
		}
	}
	
	// 构建最终的goods JSON（不设置数量，由调用方决定）
	goods := []map[string]interface{}{
		{
			"goodsCode": goodsCode,
			"skuCode":   sku.SKUCode,
			"specs":     specsArray,
		},
	}
	
	return goods, nil
}

// RecordMatchFailure 记录匹配失败
func (s *ProductMappingService) RecordMatchFailure(distributorID uint, productName string, specs map[string]string, errorReason string) error {
	specsJSON, _ := json.Marshal(specs)
	
	log := models.ProductMatchLog{
		DistributorID: distributorID,
		InputProduct:  productName,
		InputSpecs:    models.JSON(specsJSON),
		ErrorReason:   errorReason,
		RequestTime:   time.Now(),
	}
	
	// 尝试生成建议
	suggestions := s.generateSuggestions(productName)
	if len(suggestions) > 0 {
		suggestionsJSON, _ := json.Marshal(suggestions)
		log.Suggestions = models.JSON(suggestionsJSON)
	}
	
	return s.db.Create(&log).Error
}

// generateSuggestions 生成建议
func (s *ProductMappingService) generateSuggestions(productName string) []string {
	var suggestions []string
	
	// 模糊搜索相似商品
	var products []models.Product
	s.db.Where("goods_name LIKE ?", "%"+productName+"%").Limit(3).Find(&products)
	
	for _, p := range products {
		suggestions = append(suggestions, p.GoodsName)
	}
	
	return suggestions
}

// getSpecType 获取规格类型标识
func (s *ProductMappingService) getSpecType(specsCode string) string {
	// 根据specsCode映射到类型标识
	typeMap := map[string]string{
		"64": "size",
		"50": "bean",
		"17": "temperature",
		"14": "sweetness",
		"18": "sweetness",
		"15": "milk",
		"100": "flavor",
	}
	
	if specType, ok := typeMap[specsCode]; ok {
		return specType
	}
	return "other"
}
package services

import (
	"backend/internal/models"
	"backend/internal/repository"
	"backend/pkg/httpclient"
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"strings"
	"time"
)

type ProductService struct {
	db           *gorm.DB
	productRepo  *repository.ProductRepository
	luckinClient *httpclient.LuckinClient
}

func NewProductService(db *gorm.DB, luckinClient *httpclient.LuckinClient) *ProductService {
	return &ProductService{
		db:           db,
		productRepo:  repository.NewProductRepository(db),
		luckinClient: luckinClient,
	}
}

// GetProductByCode 获取商品信息
func (s *ProductService) GetProductByCode(code string) (*models.Product, error) {
	return s.productRepo.GetByCode(code)
}

// GetLuckinClient 获取瑞幸客户端
func (s *ProductService) GetLuckinClient() *httpclient.LuckinClient {
	return s.luckinClient
}

// FindBestCard 为商品找到最佳的卡片
func (s *ProductService) FindBestCard(productCode string, cardService *CardService, priceRepo *repository.LuckinConfigRepository) (*models.Card, error) {
	// 1. 首先查找绑定了该商品的卡片
	var cards []models.Card
	err := s.db.Table("cards").
		Joins("JOIN card_product_bindings ON cards.id = card_product_bindings.card_id").
		Joins("JOIN products ON card_product_bindings.product_id = products.id").
		Where("products.goods_code = ? AND cards.status = 0 AND card_product_bindings.is_active = ?", productCode, true).
		Find(&cards).Error

	if err == nil && len(cards) > 0 {
		// 返回第一张可用的卡片
		return &cards[0], nil
	}

	// 2. 如果还是没有找到，返回错误
	return nil, fmt.Errorf("没有找到适用于商品 %s 的可用卡片", productCode)
}

// SyncProductsFromCard 从指定卡片同步商品信息
func (s *ProductService) SyncProductsFromCard(cardCode string, storeCode string) (*SyncResult, error) {
	// Starting product sync for card and store

	result := &SyncResult{
		StartTime: time.Now(),
	}

	// 1. 获取菜单列表
	menuResp, err := s.fetchMenuByCard(cardCode, storeCode)
	if err != nil {
		return result, fmt.Errorf("获取菜单失败: %v", err)
	}

	// 2. 遍历所有商品获取详情
	for _, category := range menuResp.Data.Menu {
		if category.GoodsList == nil {
			continue
		}

		for _, goods := range category.GoodsList {
			// 获取商品详情
			goodsResp, err := s.fetchGoodsByCard(cardCode, storeCode, goods.GoodsCode)
			if err != nil {
				// Failed to fetch goods detail
				continue
			}

			// 同步单个商品
			if err := s.syncSingleProduct(goodsResp.Data, cardCode); err != nil {
				// Failed to sync product
				result.FailedCount++
			} else {
				result.SyncedCount++
			}
		}
	}

	result.EndTime = time.Now()
	// Product sync completed
	return result, nil
}

// fetchMenuByCard 获取菜单
func (s *ProductService) fetchMenuByCard(cardCode string, storeCode string) (*httpclient.MenuByCardResponse, error) {
	req := httpclient.MenuByCardRequest{
		ProductID:      6, // 固定值
		OrderType:      1,
		StoreCode:      storeCode,
		UpDiscountRate: "0",
		Card:           cardCode,
	}

	return s.luckinClient.MenuByCard(&req)
}

// fetchGoodsByCard 获取商品详情
func (s *ProductService) fetchGoodsByCard(cardCode string, storeCode string, goodsCode string) (*httpclient.GoodsByCardResponse, error) {
	req := httpclient.GoodsByCardRequest{
		ProductID:      6,
		OrderType:      1,
		StoreCode:      storeCode,
		LinkID:         goodsCode,
		UpDiscountRate: "0",
		Card:           cardCode,
	}

	return s.luckinClient.GoodsByCard(&req)
}

// syncSingleProduct 同步单个商品
func (s *ProductService) syncSingleProduct(goodsData struct {
	GoodsCode  string                 `json:"goodsCode"`
	GoodsName  string                 `json:"goodsName"`
	GoodsImage string                 `json:"goodsImage"`
	CostPrice  string                 `json:"costPrice"`
	LinePrice  string                 `json:"linePrice"`
	GoodsDesc  string                 `json:"goodsDesc"`
	GoodsSkus  []httpclient.GoodsSku  `json:"goodsSkus"`
	GoodsSpecs []httpclient.GoodsSpec `json:"goodsSpecs"`
}, cardCode string) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 1. 查找或创建商品
		var product models.Product
		err := tx.Where("goods_code = ?", goodsData.GoodsCode).First(&product).Error
		if err == gorm.ErrRecordNotFound {
			// 创建新商品
			product = models.Product{
				GoodsCode:  goodsData.GoodsCode,
				GoodsName:  goodsData.GoodsName,
				LastSyncAt: time.Now(),
			}
			if err := tx.Create(&product).Error; err != nil {
				return err
			}
		} else if err != nil {
			return err
		} else {
			// 检查是否需要添加别名
			if product.GoodsName != goodsData.GoodsName {
				// 添加旧名称为别名
				var alias models.ProductAlias
				err := tx.Where("product_id = ? AND alias_name = ?", product.ID, product.GoodsName).First(&alias).Error
				if err == gorm.ErrRecordNotFound {
					alias = models.ProductAlias{
						ProductID: product.ID,
						AliasName: product.GoodsName,
					}
					tx.Create(&alias)
				}
			}

			// 更新商品信息
			product.GoodsName = goodsData.GoodsName
			product.LastSyncAt = time.Now()
			tx.Save(&product)
		}

		// 2. 解析并保存SKU和规格信息
		for _, sku := range goodsData.GoodsSkus {
			if err := s.syncProductSKU(tx, product.ID, sku, goodsData.GoodsSpecs); err != nil {
				// Failed to sync SKU
				fmt.Printf("Failed to sync SKU %s: %v\n", sku.SkuCode, err)
			}
		}

		// 2.5 保存商品的可用规格信息到 available_specs 字段
		availableSpecs := make(map[string]interface{})
		for _, spec := range goodsData.GoodsSpecs {
			specItems := make([]map[string]interface{}, 0)
			for _, item := range spec.SpecsItems {
				specItems = append(specItems, map[string]interface{}{
					"code":       item.Code,
					"name":       item.Name,
					"is_default": item.IsDefault,
				})
			}
			availableSpecs[spec.SpecsCode] = map[string]interface{}{
				"name":     spec.SpecsName,
				"type":     spec.SpecsType,
				"required": spec.SpecsMinChoices > 0,
				"items":    specItems,
			}
		}

		// 将规格信息保存到产品的 available_specs 字段
		if len(availableSpecs) > 0 {
			specsJSON, _ := json.Marshal(availableSpecs)
			product.AvailableSpecs = string(specsJSON)
			tx.Save(&product)
		}

		// 3. 建立商品与卡片的直接绑定关系
		var card models.Card
		if err := tx.Where("card_code = ?", cardCode).First(&card).Error; err == nil {
			// 创建卡片与商品的直接绑定关系
			var binding models.CardProductBinding
			err = tx.Where("card_id = ? AND product_id = ?", card.ID, product.ID).First(&binding).Error
			if err == gorm.ErrRecordNotFound {
				binding = models.CardProductBinding{
					CardID:    card.ID,
					ProductID: product.ID,
					Priority:  0,
					IsActive:  true,
				}
				tx.Create(&binding)
			}
		}

		return nil
	})
}

// syncProductSKU 同步SKU信息
func (s *ProductService) syncProductSKU(tx *gorm.DB, productID uint, skuData httpclient.GoodsSku, specsData []httpclient.GoodsSpec) error {
	// 查找或创建SKU
	var sku models.ProductSKU
	err := tx.Where("product_id = ? AND sku_code = ?", productID, skuData.SkuCode).First(&sku).Error
	if err == gorm.ErrRecordNotFound {
		sku = models.ProductSKU{
			ProductID: productID,
			SKUCode:   skuData.SkuCode,
			SKUName:   skuData.SkuShowName,
		}
		if err := tx.Create(&sku).Error; err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	// 解析SKU名称中的规格选项（如 "0_0_0_4_0"）
	specValues := strings.Split(skuData.SkuName, "_")

	// 构建中文描述
	var chineseDesc []string
	var specConfigs []*models.ProductSpecConfig

	// 保存规格信息
	for i, spec := range specsData {
		if i >= len(specValues) {
			break
		}

		// 查找或创建规格
		var productSpec models.ProductSpec
		err := tx.Where("sku_id = ? AND specs_code = ?", sku.ID, spec.SpecsCode).First(&productSpec).Error
		if err == gorm.ErrRecordNotFound {
			productSpec = models.ProductSpec{
				SKUID:      sku.ID,
				SpecsCode:  spec.SpecsCode,
				SpecsName:  spec.SpecsName,
				SpecsType:  spec.SpecsType,
				IsRequired: spec.SpecsMinChoices > 0,
			}
			if err := tx.Create(&productSpec).Error; err != nil {
				return err
			}
		}

		// 确定规格类型的中文名称
		chineseName := s.getSpecChineseName(spec.SpecsCode, spec.SpecsName)
		specType := s.getSpecType(spec.SpecsCode)

		// 查找或创建商品规格配置
		var specConfig models.ProductSpecConfig
		err = tx.Where("product_id = ? AND spec_type = ?", productID, specType).First(&specConfig).Error
		if err == gorm.ErrRecordNotFound {
			specConfig = models.ProductSpecConfig{
				ProductID:    productID,
				SpecType:     specType,
				ChineseName:  chineseName,
				IsRequired:   spec.SpecsMinChoices > 0,
				DisplayOrder: i,
			}
			tx.Create(&specConfig)
		}
		specConfigs = append(specConfigs, &specConfig)

		// 保存选项信息并找到当前SKU对应的选项
		selectedValue := specValues[i]
		for _, item := range spec.SpecsItems {
			var option models.ProductSpecOption
			err := tx.Where("spec_id = ? AND code = ?", productSpec.ID, item.Code).First(&option).Error
			if err == gorm.ErrRecordNotFound {
				option = models.ProductSpecOption{
					SpecID:    productSpec.ID,
					Code:      item.Code,
					Name:      item.Name,
					IsDefault: item.IsDefault == 1,
				}
				tx.Create(&option)
			}

			// 如果是当前SKU选中的选项，添加到中文描述
			if item.Code == selectedValue {
				chineseDesc = append(chineseDesc, item.Name)

				// 如果是默认值，更新规格配置的默认值
				if item.IsDefault == 1 && specConfig.DefaultValue == "" {
					specConfig.DefaultValue = item.Name
					tx.Save(&specConfig)
				}
			}
		}
	}

	// 创建或更新SKU中文映射
	chineseDescStr := strings.Join(chineseDesc, "")
	specsCodeStr := strings.Join(specValues[:len(specsData)], "_") // SKU的specs_code就是specValues的组合

	var skuMapping models.ProductSKUMapping
	err = tx.Where("product_id = ? AND sku_code = ?", productID, skuData.SkuCode).First(&skuMapping).Error
	if err == gorm.ErrRecordNotFound {
		skuMapping = models.ProductSKUMapping{
			ProductID:   productID,
			SKUCode:     skuData.SkuCode,
			ChineseDesc: chineseDescStr,
			SpecsCode:   specsCodeStr,
		}
		tx.Create(&skuMapping)
	} else if skuMapping.ChineseDesc != chineseDescStr || skuMapping.SpecsCode != specsCodeStr {
		skuMapping.ChineseDesc = chineseDescStr
		skuMapping.SpecsCode = specsCodeStr
		tx.Save(&skuMapping)
	}

	return nil
}

// getSpecChineseName 获取规格的中文名称
func (s *ProductService) getSpecChineseName(specsCode string, specsName string) string {
	// 根据specsCode映射到中文名称
	nameMap := map[string]string{
		"64":  "杯型",
		"50":  "咖啡豆",
		"17":  "温度",
		"14":  "甜度",
		"18":  "甜度", // 有些商品可能用不同的code
		"15":  "奶",
		"100": "口味",
	}

	if chineseName, ok := nameMap[specsCode]; ok {
		return chineseName
	}
	return specsName
}

// getSpecType 获取规格类型标识
func (s *ProductService) getSpecType(specsCode string) string {
	// 根据specsCode映射到类型标识
	typeMap := map[string]string{
		"64":  "size",
		"50":  "bean",
		"17":  "temperature",
		"14":  "sweetness",
		"18":  "sweetness",
		"15":  "milk",
		"100": "flavor",
	}

	if specType, ok := typeMap[specsCode]; ok {
		return specType
	}
	return "other"
}

// SyncResult 同步结果
type SyncResult struct {
	StartTime    time.Time
	EndTime      time.Time
	SyncedCount  int
	NewCount     int
	UpdatedCount int
	FailedCount  int
	AliasCount   int
}

// GetProductList 获取商品列表
func (s *ProductService) GetProductList(page, pageSize int, search string) ([]map[string]interface{}, int64, error) {
	offset := (page - 1) * pageSize

	// 获取商品列表
	products, total, err := s.productRepo.GetProductList(offset, pageSize, search)
	if err != nil {
		return nil, 0, err
	}

	// 构建响应数据，包含每个商品的卡片数量
	var result []map[string]interface{}
	for _, product := range products {
		// 统计该商品关联的卡片数量
		var cardCount int64
		err := s.db.Model(&models.CardProductBinding{}).
			Where("product_id = ? AND is_active = ?", product.ID, true).
			Count(&cardCount).Error
		if err != nil {
			// 记录错误但不影响主流程
			cardCount = 0
		}

		// 获取别名
		var aliases []models.ProductAlias
		s.db.Where("product_id = ?", product.ID).Find(&aliases)

		productData := map[string]interface{}{
			"id":              product.ID,
			"goods_code":      product.GoodsCode,
			"goods_name":      product.GoodsName,
			"available_specs": product.AvailableSpecs,
			"last_sync_at":    product.LastSyncAt,
			"created_at":      product.CreatedAt,
			"updated_at":      product.UpdatedAt,
			"card_count":      cardCount,
			"aliases":         aliases,
			"skus":            product.SKUs,
		}

		result = append(result, productData)
	}

	return result, total, nil
}

// GetCardBoundProducts 获取卡片绑定的产品
func (s *ProductService) GetCardBoundProducts(cardID uint, page, pageSize int, search string) ([]map[string]interface{}, int64, error) {
	offset := (page - 1) * pageSize

	query := s.db.Model(&models.Product{}).
		Joins("JOIN card_product_bindings ON card_product_bindings.product_id = products.id").
		Where("card_product_bindings.card_id = ? AND card_product_bindings.is_active = ?", cardID, true)

	if search != "" {
		query = query.Where("products.goods_name LIKE ? OR products.goods_code LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var products []*models.Product
	err := query.Offset(offset).Limit(pageSize).
		Order("card_product_bindings.priority DESC, products.id DESC").
		Preload("SKUs").
		Preload("SKUs.Mapping").
		Preload("SKUs.Specs").
		Preload("SKUs.Specs.Options").
		Find(&products).Error
	if err != nil {
		return nil, 0, err
	}

	// 构建响应数据
	var result []map[string]interface{}
	for _, product := range products {
		// 获取别名
		var aliases []models.ProductAlias
		s.db.Where("product_id = ?", product.ID).Find(&aliases)

		productData := map[string]interface{}{
			"id":              product.ID,
			"goods_code":      product.GoodsCode,
			"goods_name":      product.GoodsName,
			"available_specs": product.AvailableSpecs,
			"last_sync_at":    product.LastSyncAt,
			"created_at":      product.CreatedAt,
			"updated_at":      product.UpdatedAt,
			"card_count":      1, // 至少有当前卡片的绑定
			"aliases":         aliases,
			"skus":            product.SKUs,
		}

		result = append(result, productData)
	}

	return result, total, nil
}

// SearchProducts 搜索商品
func (s *ProductService) SearchProducts(keyword string) ([]*models.Product, error) {
	return s.productRepo.SearchProducts(keyword)
}

// GetProductsByCodes 根据商品代码批量获取商品
func (s *ProductService) GetProductsByCodes(codes []string) ([]*models.Product, error) {
	return s.productRepo.GetProductsByCodes(codes)
}

// GetMatchLogs 获取匹配失败日志
func (s *ProductService) GetMatchLogs(page, pageSize int, distributorID uint, startDate, endDate string) ([]*models.ProductMatchLog, int64, error) {
	query := s.db.Model(&models.ProductMatchLog{}).
		Preload("Distributor")

	// 筛选条件
	if distributorID > 0 {
		query = query.Where("distributor_id = ?", distributorID)
	}

	if startDate != "" {
		query = query.Where("request_time >= ?", startDate)
	}

	if endDate != "" {
		query = query.Where("request_time <= ?", endDate+" 23:59:59")
	}

	// 计算总数
	var total int64
	query.Count(&total)

	// 分页查询
	var logs []*models.ProductMatchLog
	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&logs).Error

	return logs, total, err
}

// GetAvailableProducts 获取所有可用商品信息（不暴露卡片信息）
func (s *ProductService) GetAvailableProducts(distributorID uint) (*models.AvailableProductsResponse, error) {
	// 查询所有有可用卡片的商品
	query := `
		SELECT DISTINCT p.*, COUNT(DISTINCT c.id) as available_count
		FROM products p
		JOIN card_product_bindings cpb ON cpb.product_id = p.id
		JOIN cards c ON c.id = cpb.card_id
		WHERE c.status = 0 AND cpb.is_active = true
		GROUP BY p.id
		ORDER BY available_count DESC, p.goods_name
	`

	var products []struct {
		models.Product
		AvailableCount int `gorm:"column:available_count"`
	}

	if err := s.db.Raw(query).Scan(&products).Error; err != nil {
		return nil, err
	}

	// 构建响应数据
	publicProducts := make([]models.PublicProductInfo, 0)

	for _, p := range products {
		// 获取商品的所有SKU和规格信息
		var product models.Product
		err := s.db.Preload("SKUs").
			Preload("SKUs.Mapping").
			Preload("SKUs.Specs").
			Preload("SKUs.Specs.Options").
			Preload("Aliases").
			First(&product, p.ID).Error
		if err != nil {
			continue
		}

		// 构建公开的SKU信息
		publicSKUs := make([]models.PublicSKUInfo, 0)
		for _, sku := range product.SKUs {
			specs := make(map[string]models.SpecOption)
			specsCode := ""

			// 构建规格信息
			for _, spec := range sku.Specs {
				for _, option := range spec.Options {
					if option.IsDefault {
						specs[s.getSpecType(spec.SpecsCode)] = models.SpecOption{
							Code: option.Code,
							Name: option.Name,
						}
						if specsCode != "" {
							specsCode += "_"
						}
						specsCode += option.Code
					}
				}
			}

			// 使用数据库中存储的specs_code和中文描述
			chineseDesc := ""
			if sku.Mapping != nil {
				if sku.Mapping.SpecsCode != "" {
					specsCode = sku.Mapping.SpecsCode
				}
				chineseDesc = sku.Mapping.ChineseDesc
			}

			publicSKU := models.PublicSKUInfo{
				SKUCode:     sku.SKUCode,
				ChineseDesc: chineseDesc,
				Specs:       specs,
				SpecsCode:   specsCode,
				IsDefault:   len(publicSKUs) == 0, // 第一个SKU作为默认
				IsAvailable: true,
			}
			publicSKUs = append(publicSKUs, publicSKU)
		}

		// 获取商品别名
		aliases := make([]string, 0)
		for _, alias := range product.Aliases {
			aliases = append(aliases, alias.AliasName)
		}

		// 收集所有可用的规格选项
		specsOptions := make(map[string][]models.SpecOption)
		specTypes := make(map[string]bool)

		// 遍历所有SKU的规格来收集选项
		for _, sku := range product.SKUs {
			for _, spec := range sku.Specs {
				specType := s.getSpecType(spec.SpecsCode)
				if !specTypes[specType] {
					specTypes[specType] = true
					for _, option := range spec.Options {
						specsOptions[specType] = append(specsOptions[specType], models.SpecOption{
							Code: option.Code,
							Name: option.Name,
						})
					}
				}
			}
		}

		// 构建公开的商品信息
		publicProduct := models.PublicProductInfo{
			GoodsCode:    product.GoodsCode,
			GoodsName:    product.GoodsName,
			Category:     "咖啡", // TODO: 从数据库获取分类
			IsAvailable:  p.AvailableCount > 0,
			StockStatus:  models.CalculateStockStatus(p.AvailableCount),
			PriceRange:   "15-30元", // TODO: 从价格映射获取
			SKUs:         publicSKUs,
			Aliases:      aliases,
			SpecsOptions: specsOptions,
			OrderingTips: fmt.Sprintf("建议使用商品名'%s'下单", product.GoodsName),
		}

		publicProducts = append(publicProducts, publicProduct)
	}

	// 构建响应
	response := &models.AvailableProductsResponse{
		UpdatedAt: time.Now(),
		Products:  publicProducts,
		Notice:    "商品可用性每5分钟更新一次",
	}

	return response, nil
}

// ParseSpecsCode 解析规格代码
type SpecItem struct {
	SpecsCode string `json:"specs_code"`
	Code      string `json:"code"`
	Name      string `json:"name"`
}

// ParseSpecsCode 将specs_code（如"0_0_0_0"）解析成具体的规格项数组
func (s *ProductService) ParseSpecsCode(goodsCode, skuCode, specsCode string) ([]SpecItem, error) {
	// 使用固定的规格类型映射，不依赖数据库
	// 这样可以处理不同商品有不同规格数量的情况
	specValues := strings.Split(specsCode, "_")

	// 常见的规格类型顺序（基于瑞幸API的常见模式）
	commonSpecTypes := []struct {
		SpecsCode string
		Name      string
		Options   map[string]string
	}{
		{
			SpecsCode: "64", // 杯型
			Name:      "杯型",
			Options: map[string]string{
				"0": "中杯", "1": "大杯", "2": "超大杯",
			},
		},
		{
			SpecsCode: "17", // 温度
			Name:      "温度",
			Options: map[string]string{
				"0": "热", "1": "冰",
			},
		},
		{
			SpecsCode: "18", // 甜度
			Name:      "甜度",
			Options: map[string]string{
				"0": "不另外加糖", "1": "微甜", "2": "半糖", "3": "正常糖",
			},
		},
		{
			SpecsCode: "15", // 奶类型
			Name:      "奶类型",
			Options: map[string]string{
				"0": "全脂牛奶", "1": "脱脂牛奶", "2": "豆奶", "3": "椰奶", "4": "燕麦奶",
			},
		},
		{
			SpecsCode: "100", // 口味/其他
			Name:      "口味",
			Options: map[string]string{
				"0": "原味", "1": "香草", "2": "焦糖", "3": "榛果",
			},
		},
	}

	var specItems []SpecItem
	for i, value := range specValues {
		if i >= len(commonSpecTypes) {
			// 如果超出了已知的规格类型，停止解析
			break
		}

		specType := commonSpecTypes[i]
		optionName, exists := specType.Options[value]
		if !exists {
			// 如果找不到对应的选项名称，使用默认值
			optionName = "选项" + value
		}

		specItems = append(specItems, SpecItem{
			SpecsCode: specType.SpecsCode,
			Code:      value,
			Name:      optionName,
		})
	}

	if len(specItems) == 0 {
		return nil, fmt.Errorf("无效的规格代码: %s", specsCode)
	}

	return specItems, nil
}

// GetDB 获取数据库实例
func (s *ProductService) GetDB() *gorm.DB {
	return s.db
}

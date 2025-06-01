package services

import (
	"backend/internal/models"
	"errors"
	"fmt"
	"go.uber.org/zap"
)

// SimplifiedOrderRequest 简化下单请求
type SimplifiedOrderRequest struct {
	DistributorID uint                     `json:"distributor_id"`
	CardCode      string                   `json:"card_code" binding:"required"`
	StoreName     string                   `json:"store_name" binding:"required"`
	PhoneNumber   string                   `json:"phone_number" binding:"required"`
	CallbackURL   string                   `json:"callback_url" binding:"required"`
	Items         []SimplifiedOrderItem    `json:"items" binding:"required,min=1"`
}

// SimplifiedOrderItem 简化订单项
type SimplifiedOrderItem struct {
	Product  string            `json:"product" binding:"required"`
	Specs    map[string]string `json:"specs" binding:"required"`
	Quantity int               `json:"quantity" binding:"required,min=1"`
}

// CreateSimplifiedOrder 创建简化订单
func (s *OrderService) CreateSimplifiedOrder(req *SimplifiedOrderRequest) (*models.Order, error) {
	s.logger.Info("Creating simplified order", 
		zap.Uint("distributorID", req.DistributorID),
		zap.String("storeName", req.StoreName))
	
	// 1. 查找店铺
	storeCode, storeAddress, err := s.findStoreByName(req.StoreName)
	if err != nil {
		// 记录查找失败
		s.mappingService.RecordMatchFailure(req.DistributorID, 
			"店铺:"+req.StoreName, nil, err.Error())
		return nil, fmt.Errorf("店铺'%s'未找到", req.StoreName)
	}
	
	// 2. 获取卡片信息（用于限定商品范围）
	card, err := s.cardService.GetByCode(req.CardCode)
	if err != nil {
		return nil, fmt.Errorf("卡片不存在: %w", err)
	}
	
	// 3. 映射商品和规格
	var goodsList []struct {
		GoodsCode string `json:"goods_code"`
		SKUCode   string `json:"sku_code"`
		Quantity  int    `json:"quantity"`
		SpecsCode string `json:"specs_code"` // 支持直接传入规格编码如"0_0_0_4_0"
		Specs     []struct {
			SpecsCode string `json:"specs_code"`
			Code      string `json:"code"`
			Name      string `json:"name"`
		} `json:"specs"` // 可选，如果提供SpecsCode则忽略此字段
	}
	
	for _, item := range req.Items {
		// 映射商品（限定在卡片绑定的商品范围内）
		mapReq := MapProductRequest{
			ProductName: item.Product,
			Specs:       item.Specs,
			CardID:      card.ID,
		}
		
		mapResult, err := s.mappingService.MapProduct(mapReq)
		if err != nil {
			// 记录映射失败
			s.mappingService.RecordMatchFailure(req.DistributorID, 
				item.Product, item.Specs, err.Error())
			return nil, err
		}
		
		// 构建商品数据
		if len(mapResult.GoodsJSON) > 0 {
			goodsData := mapResult.GoodsJSON[0]
			goodsItem := struct {
				GoodsCode string `json:"goods_code"`
				SKUCode   string `json:"sku_code"`
				Quantity  int    `json:"quantity"`
				SpecsCode string `json:"specs_code"` // 支持直接传入规格编码如"0_0_0_4_0"
				Specs     []struct {
					SpecsCode string `json:"specs_code"`
					Code      string `json:"code"`
					Name      string `json:"name"`
				} `json:"specs"` // 可选，如果提供SpecsCode则忽略此字段
			}{
				GoodsCode: goodsData["goodsCode"].(string),
				SKUCode:   goodsData["skuCode"].(string),
				Quantity:  item.Quantity,
			}
			
			// 转换规格数据
			if specsData, ok := goodsData["specs"].([]map[string]interface{}); ok {
				for _, spec := range specsData {
					goodsItem.Specs = append(goodsItem.Specs, struct {
						SpecsCode string `json:"specs_code"`
						Code      string `json:"code"`
						Name      string `json:"name"`
					}{
						SpecsCode: spec["specsCode"].(string),
						Code:      spec["code"].(string),
						Name:      spec["name"].(string),
					})
				}
			}
			
			goodsList = append(goodsList, goodsItem)
		}
	}
	
	// 4. 创建标准订单请求
	standardReq := &CreateOrderRequest{
		DistributorID: req.DistributorID,
		CardCode:      req.CardCode,
		StoreCode:     storeCode,
		PhoneNumber:   req.PhoneNumber,
		CallbackURL:   req.CallbackURL,
		Goods:         goodsList,
	}
	
	// 5. 调用标准下单流程
	order, err := s.CreateOrder(standardReq)
	if err != nil {
		return nil, err
	}
	
	// 6. 更新订单的店铺信息（如果标准流程没有设置）
	if order.StoreAddress == "" {
		order.StoreAddress = storeAddress
		s.orderRepo.Update(order)
	}
	
	return order, nil
}

// findStoreByName 通过店铺名称查找店铺代码
func (s *OrderService) findStoreByName(storeName string) (string, string, error) {
	// 先尝试从数据库查找缓存的店铺信息
	store, err := s.storeService.FindByName(storeName)
	if err == nil && store != nil {
		return store.StoreCode, store.Address, nil
	}
	
	// 如果数据库没有，尝试通过API搜索
	// 这里需要一个默认的卡片来搜索，可以从配置或者可用卡片中获取
	card, err := s.cardService.GetAnyAvailableCard()
	if err != nil {
		return "", "", errors.New("没有可用的卡片进行店铺搜索")
	}
	
	// 调用瑞幸API搜索店铺
	stores, err := s.storeService.SearchStores(card.CardCode, 0, storeName)
	if err != nil {
		return "", "", err
	}
	
	if len(stores) == 0 {
		return "", "", errors.New("未找到匹配的店铺")
	}
	
	// 返回第一个匹配的店铺
	return stores[0].StoreCode, stores[0].Address, nil
}
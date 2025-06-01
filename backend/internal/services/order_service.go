package services

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"backend/internal/models"
	"backend/internal/repository"
	"backend/pkg/httpclient"
	"github.com/google/uuid"
	"github.com/skip2/go-qrcode"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type OrderService struct {
	db              *gorm.DB
	orderRepo       *repository.OrderRepository
	cardRepo        *repository.CardRepository
	distributorRepo *repository.DistributorRepository
	luckinClient    *httpclient.LuckinClient
	financeService  *FinanceService
	cardService     *CardService
	mappingService  *ProductMappingService
	storeService    *StoreService
	logger          *zap.Logger
}

func NewOrderService(
	db *gorm.DB,
	orderRepo *repository.OrderRepository,
	cardRepo *repository.CardRepository,
	distributorRepo *repository.DistributorRepository,
	luckinClient *httpclient.LuckinClient,
	cardService *CardService,
	logger *zap.Logger,
) *OrderService {
	return &OrderService{
		db:              db,
		orderRepo:       orderRepo,
		cardRepo:        cardRepo,
		distributorRepo: distributorRepo,
		luckinClient:    luckinClient,
		financeService:  NewFinanceService(db),
		cardService:     cardService,
		mappingService:  NewProductMappingService(db),
		storeService:    NewStoreService(cardRepo, repository.NewCityRepository(db), luckinClient, logger),
		logger:          logger,
	}
}

type CreateOrderRequest struct {
	DistributorID uint   `json:"distributor_id"`
	CardCode      string `json:"card_code"`
	StoreCode     string `json:"store_code"`
	PhoneNumber   string `json:"phone_number"`
	CallbackURL   string `json:"callback_url"`
	Goods         []struct {
		GoodsCode string `json:"goods_code"`
		SKUCode   string `json:"sku_code"`
		Quantity  int    `json:"quantity"`
		SpecsCode string `json:"specs_code"` // 支持直接传入规格编码如"0_0_0_4_0"
		Specs     []struct {
			SpecsCode string `json:"specs_code"`
			Code      string `json:"code"`
			Name      string `json:"name"`
		} `json:"specs"` // 可选，如果提供SpecsCode则忽略此字段
	} `json:"goods"`
}

func (s *OrderService) CreateOrder(req *CreateOrderRequest) (*models.Order, error) {
	// 1. 验证卡片状态
	var card *models.Card
	var err error
	var needReleaseCard bool // 标记是否需要释放卡片
	
	if req.CardCode != "" {
		// 兼容旧逻辑：如果提供了卡片代码，直接使用
		card, err = s.cardRepo.GetByCode(req.CardCode)
		if err != nil {
			return nil, fmt.Errorf("卡片不存在: %w", err)
		}
		if card.Status != 0 {
			return nil, errors.New("卡片状态不正确")
		}
	} else {
		// 新逻辑：根据商品自动选择卡片
		if len(req.Goods) == 0 {
			return nil, errors.New("商品信息不能为空")
		}
		
		// 使用第一个商品来匹配价格（通常一个订单只有一个主商品）
		productService := NewProductService(s.db, s.luckinClient)
		luckinConfigRepo := repository.NewLuckinConfigRepository(s.db)
		card, err = productService.FindBestCard(req.Goods[0].GoodsCode, s.cardService, luckinConfigRepo)
		if err != nil {
			return nil, fmt.Errorf("自动选卡失败: %w", err)
		}
		req.CardCode = card.CardCode
	}
	
	// 预占卡片库存（标记为使用中）
	if err := s.cardService.ReserveCard(card.ID); err != nil {
		return nil, fmt.Errorf("预占卡片失败: %w", err)
	}
	needReleaseCard = true // 如果后续失败需要释放
	
	// 使用defer确保失败时释放卡片
	defer func() {
		if needReleaseCard {
			s.cardService.ReleaseCard(card.ID)
		}
	}()

	// 2. 验证分销商状态
	distributor, err := s.distributorRepo.GetByID(req.DistributorID)
	if err != nil {
		return nil, fmt.Errorf("distributor not found: %w", err)
	}

	if distributor.Status != 1 {
		return nil, errors.New("distributor is not active")
	}

	// 3. 生成订单号
	orderNo := s.generateOrderNo()

	// 4. 构建商品JSON字符串
	goodsJSON, err := s.buildGoodsJSON(req.Goods)
	if err != nil {
		return nil, fmt.Errorf("failed to build goods JSON: %w", err)
	}

	// 5. 检查订单（调用瑞幸API）
	// 使用卡片的瑞幸产品ID
	productID := card.LuckinProductID
	if productID == 0 {
		productID = 6 // 默认值
	}
	
	// 添加调试日志
	s.logger.Info("CheckByCard request",
		zap.String("card_code", req.CardCode),
		zap.Int("luckin_product_id", card.LuckinProductID),
		zap.Int("product_id", productID),
		zap.String("store_code", req.StoreCode),
		zap.String("goods_json", goodsJSON))
	
	checkReq := &httpclient.CheckByCardRequest{
		StoreCode:      req.StoreCode,
		ProductID:      productID,
		Goods:          goodsJSON,
		OrderNo:        orderNo,
		TakeMode:       1,
		UpDiscountRate: "0",
		Card:           req.CardCode,
	}

	checkResp, err := s.luckinClient.CheckByCard(checkReq)
	if err != nil {
		return nil, fmt.Errorf("failed to check order: %w", err)
	}

	if checkResp.Code != 200 {
		return nil, fmt.Errorf("check order failed: %s", checkResp.Msg)
	}

	// 6. 解析金额
	// 保存瑞幸的原始价格（仅用于记录和对账）
	luckinPrice, _ := strconv.ParseFloat(checkResp.Data.TotalSalePrice, 64)
	luckinCostPrice, _ := strconv.ParseFloat(checkResp.Data.TotalPlatformCostPrice, 64)
	
	// 使用卡片配置的自定义价格
	totalAmount := card.SellPrice   // 使用卡片的销售价
	costAmount := card.CostPrice    // 使用卡片的成本价
	
	// 记录价格对比信息
	s.logger.Info("使用自定义价格",
		zap.String("order_no", orderNo),
		zap.Float64("luckin_price", luckinPrice),
		zap.Float64("luckin_cost", luckinCostPrice),
		zap.Float64("custom_price", totalAmount),
		zap.Float64("custom_cost", costAmount))

	// 7. 检查余额并扣款
	err = s.financeService.Consume(int64(req.DistributorID), totalAmount, orderNo, 1) // 系统自动扣款
	if err != nil {
		return nil, fmt.Errorf("balance insufficient: %w", err)
	}

	// 8. 创建订单记录
	order := &models.Order{
		OrderNo:         orderNo,
		DistributorID:   req.DistributorID,
		CardID:          card.ID,
		CardCode:        req.CardCode,
		Status:          models.OrderStatusPending,
		StoreCode:       req.StoreCode,
		StoreName:       checkResp.Data.Name,
		StoreAddress:    checkResp.Data.Address,
		TotalAmount:     totalAmount,
		CostAmount:      costAmount,
		ProfitAmount:    totalAmount - costAmount,
		LuckinPrice:     luckinPrice,     // 保存瑞幸原始价格
		LuckinCostPrice: luckinCostPrice, // 保存瑞幸原始成本
		TakeMode:        1,
		PhoneNumber:     req.PhoneNumber,
		CallbackURL:     req.CallbackURL,
		Goods:           s.parseOrderGoods(checkResp.Data.Goods),
	}

	// 9. 提交订单到瑞幸
	orderReq := &httpclient.OrderByCardRequest{
		ProductID:      productID, // 使用之前获取的正确ProductID
		OrderNo:        orderNo,
		PhoneNo:        req.PhoneNumber,
		TakeMode:       "1",
		CallbackURL:    req.CallbackURL,
		UpDiscountRate: "0",
		Card:           req.CardCode,
	}

	orderResp, err := s.luckinClient.OrderByCard(orderReq)
	if err != nil {
		order.Status = models.OrderStatusFailed
		s.orderRepo.Create(order)
		// 退款
		s.financeService.Refund(int64(req.DistributorID), totalAmount, orderNo, "订单提交失败，自动退款", 1)
		return nil, fmt.Errorf("failed to submit order: %w", err)
	}

	if orderResp.Code != 200 {
		order.Status = models.OrderStatusFailed
		s.orderRepo.Create(order)
		// 退款
		s.financeService.Refund(int64(req.DistributorID), totalAmount, orderNo, fmt.Sprintf("订单提交失败：%s，自动退款", orderResp.Msg), 1)
		return nil, fmt.Errorf("submit order failed: %s", orderResp.Msg)
	}

	// 9. 更新订单信息
	order.PutOrderID = orderResp.Data.PutOrderID
	order.Status = models.OrderStatusDoing
	respJSON, _ := json.Marshal(orderResp)
	order.LuckinResponse = string(respJSON)

	// 10. 保存订单
	if err := s.orderRepo.Create(order); err != nil {
		return nil, fmt.Errorf("failed to save order: %w", err)
	}

	// 11. 标记卡片为已使用
	if err := s.cardService.MarkAsUsed(card.ID, order.ID); err != nil {
		s.logger.Error("failed to mark card as used", zap.Error(err))
	}
	needReleaseCard = false // 成功了，不需要释放卡片

	// 12. 记录卡片使用日志
	s.cardRepo.LogUsage(&models.CardUsageLog{
		CardID:        card.ID,
		CardCode:      card.CardCode,
		DistributorID: req.DistributorID,
		OrderNo:       orderNo,
		Success:       true,
	})

	// 13. 启动异步查询任务
	go s.queryOrderStatus(order)

	return order, nil
}

func (s *OrderService) generateOrderNo() string {
	timestamp := time.Now().UnixNano() / 1000000
	uuid := uuid.New().String()[:8]
	return fmt.Sprintf("DD%d%s", timestamp, uuid)
}

func (s *OrderService) buildGoodsJSON(goods []struct {
	GoodsCode string `json:"goods_code"`
	SKUCode   string `json:"sku_code"`
	Quantity  int    `json:"quantity"`
	SpecsCode string `json:"specs_code"` // 支持直接传入规格编码如"0_0_0_4_0"
	Specs     []struct {
		SpecsCode string `json:"specs_code"`
		Code      string `json:"code"`
		Name      string `json:"name"`
	} `json:"specs"` // 可选，如果提供SpecsCode则忽略此字段
}) (string, error) {
	var items []map[string]interface{}

	for _, g := range goods {
		var specs []map[string]interface{}
		
		// 优先使用SpecsCode字段进行自动解析
		if g.SpecsCode != "" {
			// 使用产品服务解析specs_code
			productService := NewProductService(s.db, s.luckinClient)
			specItems, err := productService.ParseSpecsCode(g.GoodsCode, g.SKUCode, g.SpecsCode)
			if err != nil {
				return "", fmt.Errorf("解析规格代码失败: %w", err)
			}
			
			// 转换为API需要的格式
			for _, spec := range specItems {
				specs = append(specs, map[string]interface{}{
					"specsType": 2,
					"specsCode": spec.SpecsCode,
					"code":      spec.Code,
					"name":      spec.Name,
					"num":       g.Quantity,
					"itemList":  []interface{}{},
				})
			}
		} else {
			// 兼容原有的详细specs数组模式
			for _, spec := range g.Specs {
				specs = append(specs, map[string]interface{}{
					"specsType": 2,
					"specsCode": spec.SpecsCode,
					"code":      spec.Code,
					"name":      spec.Name,
					"num":       g.Quantity,
					"itemList":  []interface{}{},
				})
			}
		}

		item := map[string]interface{}{
			"code":       g.GoodsCode,
			"num":        g.Quantity,
			"skuCode":    g.SKUCode,
			"goodsSpecs": specs,
		}
		items = append(items, item)
	}

	jsonBytes, err := json.Marshal(items)
	if err != nil {
		return "", err
	}

	return string(jsonBytes), nil
}

func (s *OrderService) parseOrderGoods(checkGoods []httpclient.CheckGoods) models.OrderGoods {
	var orderGoods models.OrderGoods

	for _, g := range checkGoods {
		totalPrice, _ := strconv.ParseFloat(g.TotalOriginalPrice, 64)
		salePrice, _ := strconv.ParseFloat(g.SalePrice, 64)

		item := models.OrderGoodsItem{
			GoodsID:       g.GoodsID,
			GoodsName:     g.GoodsName,
			GoodsImage:    g.GoodsImg,
			Quantity:      g.Num,
			OriginalPrice: totalPrice,
			SalePrice:     salePrice,
		}

		// 解析规格信息
		if g.PutDetailJSON != "" {
			var detail map[string]interface{}
			if err := json.Unmarshal([]byte(g.PutDetailJSON), &detail); err == nil {
				if skuName, ok := detail["skuName"].(string); ok {
					item.SKUName = skuName
				}
			}
		}

		orderGoods = append(orderGoods, item)
	}

	return orderGoods
}

func (s *OrderService) queryOrderStatus(order *models.Order) {
	// 多次重试查询，等待订单完成
	maxRetries := 10
	retryInterval := 3 * time.Second
	
	for i := 0; i < maxRetries; i++ {
		time.Sleep(retryInterval)
		
		queryReq := &httpclient.QueryByCardRequest{
			Brand:   "lk",
			Card:    order.CardCode,
			OrderNo: order.OrderNo,
		}

		queryResp, err := s.luckinClient.QueryByCard(queryReq)
		if err != nil {
			s.logger.Error("failed to query order status",
				zap.String("order_no", order.OrderNo),
				zap.Int("retry", i+1),
				zap.Error(err))
			continue
		}

		s.logger.Info("order status query result",
			zap.String("order_no", order.OrderNo),
			zap.Int("retry", i+1),
			zap.String("status", queryResp.Data.Status),
			zap.String("take_code", queryResp.Data.TakeCode))

		if queryResp.Code == 200 && queryResp.Data.Status == "success" {
			// 更新订单状态
			order.Status = models.OrderStatusSuccess
			order.TakeCode = queryResp.Data.TakeCode
			order.QRData = queryResp.Data.QRData

			// 如果主字段为空，尝试从takeInfoList中获取
			if (order.TakeCode == "" || order.QRData == "") && len(queryResp.Data.TakeInfoList) > 0 {
				takeInfo := queryResp.Data.TakeInfoList[0] // 使用第一个取餐信息
				if order.TakeCode == "" {
					order.TakeCode = takeInfo.TakeInfo
				}
				if order.QRData == "" {
					order.QRData = takeInfo.QRData
				}
			}

			respJSON, _ := json.Marshal(queryResp)
			order.LuckinResponse = string(respJSON)

			if err := s.orderRepo.Update(order); err != nil {
				s.logger.Error("failed to update order",
					zap.String("order_no", order.OrderNo),
					zap.Error(err))
			}

			// 触发回调
			if order.CallbackURL != "" {
				go s.sendCallback(order)
			}
			
			s.logger.Info("order completed successfully",
				zap.String("order_no", order.OrderNo),
				zap.String("take_code", order.TakeCode))
			return
		}
		
		// 如果状态不是success，继续重试
		s.logger.Info("order not ready yet, retrying",
			zap.String("order_no", order.OrderNo),
			zap.String("status", queryResp.Data.Status))
	}
	
	// 重试次数用完，记录警告
	s.logger.Warn("order status query timed out",
		zap.String("order_no", order.OrderNo),
		zap.Int("max_retries", maxRetries))
}

func (s *OrderService) sendCallback(order *models.Order) {
	// 使用独立的回调服务
	callbackService := NewCallbackService(s.orderRepo, s.logger)
	if err := callbackService.SendCallback(order); err != nil {
		s.logger.Error("failed to send callback",
			zap.String("order_no", order.OrderNo),
			zap.Error(err))
	}
}

func (s *OrderService) GetOrderByNo(orderNo string) (*models.Order, error) {
	return s.orderRepo.GetByOrderNo(orderNo)
}

func (s *OrderService) List(offset, limit int, filters map[string]interface{}) ([]*models.Order, int64, error) {
	return s.orderRepo.List(offset, limit, filters)
}

func (s *OrderService) GetStatistics(distributorID uint, startDate, endDate time.Time) (map[string]interface{}, error) {
	return s.orderRepo.GetStatistics(distributorID, startDate, endDate)
}

func (s *OrderService) QueryOrderStatus(orderNo string) (*models.Order, error) {
	order, err := s.orderRepo.GetByOrderNo(orderNo)
	if err != nil {
		return nil, err
	}

	// 如果订单已完成，直接返回
	if order.Status == models.OrderStatusSuccess {
		return order, nil
	}

	return order, nil
}

// RefreshOrderStatus 手动刷新订单状态
func (s *OrderService) RefreshOrderStatus(orderNo string) (*models.Order, error) {
	order, err := s.orderRepo.GetByOrderNo(orderNo)
	if err != nil {
		return nil, err
	}

	// 如果订单已完成，直接返回
	if order.Status == models.OrderStatusSuccess {
		return order, nil
	}

	// 立即查询瑞幸API获取最新状态
	queryReq := &httpclient.QueryByCardRequest{
		Brand:   "lk",
		Card:    order.CardCode,
		OrderNo: order.OrderNo,
	}

	queryResp, err := s.luckinClient.QueryByCard(queryReq)
	if err != nil {
		return nil, fmt.Errorf("查询订单状态失败: %w", err)
	}

	// 添加调试日志
	s.logger.Info("QueryByCard response details",
		zap.String("order_no", orderNo),
		zap.Int("response_code", queryResp.Code),
		zap.String("response_msg", queryResp.Msg),
		zap.String("data_status", queryResp.Data.Status),
		zap.String("take_code", queryResp.Data.TakeCode),
		zap.String("qr_data", queryResp.Data.QRData),
		zap.Int("take_info_list_count", len(queryResp.Data.TakeInfoList)))

	if queryResp.Code != 200 {
		return nil, fmt.Errorf("查询订单失败: %s", queryResp.Msg)
	}

	// 无论状态如何，都更新订单信息
	if queryResp.Data.Status == "success" {
		order.Status = models.OrderStatusSuccess
		order.TakeCode = queryResp.Data.TakeCode
		order.QRData = queryResp.Data.QRData

		// 如果主字段为空，尝试从takeInfoList中获取
		if (order.TakeCode == "" || order.QRData == "") && len(queryResp.Data.TakeInfoList) > 0 {
			takeInfo := queryResp.Data.TakeInfoList[0]
			if order.TakeCode == "" {
				order.TakeCode = takeInfo.TakeInfo
			}
			if order.QRData == "" {
				order.QRData = takeInfo.QRData
			}
		}
	}

	// 更新瑞幸响应
	respJSON, _ := json.Marshal(queryResp)
	order.LuckinResponse = string(respJSON)

	// 保存到数据库
	if err := s.orderRepo.Update(order); err != nil {
		return nil, fmt.Errorf("更新订单失败: %w", err)
	}

	return order, nil
}

// GenerateQRCode 生成二维码图片（Base64格式）
func (s *OrderService) GenerateQRCode(qrData string) (string, error) {
	if qrData == "" {
		return "", fmt.Errorf("二维码数据不能为空")
	}
	
	// 生成二维码
	png, err := qrcode.Encode(qrData, qrcode.Medium, 256)
	if err != nil {
		return "", fmt.Errorf("生成二维码失败: %w", err)
	}
	
	// 转换为Base64格式
	encodedString := base64.StdEncoding.EncodeToString(png)
	
	// 返回Data URL格式
	return fmt.Sprintf("data:image/png;base64,%s", encodedString), nil
}

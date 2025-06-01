package services

import (
	"fmt"
	"strconv"

	"backend/internal/repository"
	"backend/pkg/httpclient"
	"go.uber.org/zap"
)

type StoreService struct {
	cardRepo     *repository.CardRepository
	cityRepo     *repository.CityRepository
	luckinClient *httpclient.LuckinClient
	logger       *zap.Logger
}

func NewStoreService(
	cardRepo *repository.CardRepository,
	cityRepo *repository.CityRepository,
	luckinClient *httpclient.LuckinClient,
	logger *zap.Logger,
) *StoreService {
	return &StoreService{
		cardRepo:     cardRepo,
		cityRepo:     cityRepo,
		luckinClient: luckinClient,
		logger:       logger,
	}
}

// SearchStores 搜索门店
func (s *StoreService) SearchStores(cardCode string, cityID int, keywords string) ([]httpclient.Store, error) {
	// 验证卡片
	card, err := s.cardRepo.GetByCode(cardCode)
	if err != nil {
		return nil, fmt.Errorf("invalid card: %w", err)
	}

	req := &httpclient.StoreByCardRequest{
		Card:      cardCode,
		ProductID: card.LuckinProductID,
		CityID:    cityID,
		Keywords:  keywords,
	}

	resp, err := s.luckinClient.StoreByCard(req)
	if err != nil {
		return nil, fmt.Errorf("failed to search stores: %w", err)
	}

	if resp.Code != 200 {
		return nil, fmt.Errorf("search stores failed: %s", resp.Msg)
	}

	return resp.Data, nil
}

// GetMenu 获取门店菜单
func (s *StoreService) GetMenu(cardCode string, storeCode string) (interface{}, error) {
	// 验证卡片
	card, err := s.cardRepo.GetByCode(cardCode)
	if err != nil {
		return nil, fmt.Errorf("invalid card: %w", err)
	}

	req := &httpclient.MenuByCardRequest{
		ProductID:      card.LuckinProductID,
		OrderType:      1,
		StoreCode:      storeCode,
		UpDiscountRate: "0",
		Card:           cardCode,
	}

	resp, err := s.luckinClient.MenuByCard(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get menu: %w", err)
	}

	if resp.Code != 200 {
		return nil, fmt.Errorf("get menu failed: %s", resp.Msg)
	}

	return resp.Data, nil
}

// GetGoodsDetail 获取商品详情
func (s *StoreService) GetGoodsDetail(cardCode string, storeCode string, goodsCode string) (interface{}, error) {
	// 验证卡片
	card, err := s.cardRepo.GetByCode(cardCode)
	if err != nil {
		return nil, fmt.Errorf("invalid card: %w", err)
	}

	req := &httpclient.GoodsByCardRequest{
		ProductID:      card.LuckinProductID,
		OrderType:      1,
		StoreCode:      storeCode,
		LinkID:         goodsCode,
		UpDiscountRate: "0",
		Card:           cardCode,
	}

	resp, err := s.luckinClient.GoodsByCard(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get goods detail: %w", err)
	}

	if resp.Code != 200 {
		return nil, fmt.Errorf("get goods detail failed: %s", resp.Msg)
	}

	return resp.Data, nil
}

// GetCityIDByName 根据城市名称获取城市ID
func (s *StoreService) GetCityIDByName(cityName string) (int, error) {
	city, err := s.cityRepo.GetByCityName(cityName)
	if err != nil {
		return 0, fmt.Errorf("city not found: %w", err)
	}
	return city.CityID, nil
}

// FindByName 通过名称查找门店（在实际实现中，可能需要缓存门店信息到数据库）
func (s *StoreService) FindByName(name string) (*httpclient.Store, error) {
	// 这是一个简化的实现，实际中可能需要缓存门店信息到数据库
	// 现在返回nil，让调用方使用API搜索
	return nil, nil
}

// StoreSearchParams 店铺搜索参数
type StoreSearchParams struct {
	CityID    string
	Keyword   string
	Longitude string
	Latitude  string
	Page      int
	PageSize  int
}

// SearchStoresAdmin 管理员搜索店铺（使用任意可用卡片）
func (s *StoreService) SearchStoresAdmin(params *StoreSearchParams) ([]map[string]interface{}, int, error) {
	// 获取任意一张可用的卡片
	card, err := s.cardRepo.GetAnyAvailable()
	if err != nil {
		return nil, 0, fmt.Errorf("没有可用的卡片用于查询: %w", err)
	}

	// 构建请求参数
	req := &httpclient.StoreByCardRequest{
		Card:      card.CardCode,
		ProductID: 6, // 固定值
	}
	
	// 设置城市ID
	if params.CityID != "" {
		cityID, _ := strconv.Atoi(params.CityID)
		req.CityID = cityID
	}
	
	// 设置关键词
	if params.Keyword != "" {
		req.Keywords = params.Keyword
	}
	
	// 注意：瑞幸API不支持直接传经纬度搜索，需要通过其他方式实现

	// 调用瑞幸API
	resp, err := s.luckinClient.StoreByCard(req)
	if err != nil {
		return nil, 0, fmt.Errorf("调用门店搜索API失败: %w", err)
	}

	if resp.Code != 200 {
		return nil, 0, fmt.Errorf("门店搜索失败: %s", resp.Msg)
	}

	// 转换为通用格式
	stores := make([]map[string]interface{}, 0)
	start := (params.Page - 1) * params.PageSize
	end := start + params.PageSize
	
	// 处理分页
	totalStores := len(resp.Data)
	if start >= totalStores {
		return stores, totalStores, nil
	}
	if end > totalStores {
		end = totalStores
	}
	
	// 转换数据格式
	for _, store := range resp.Data[start:end] {
		// 计算距离（如果提供了经纬度）
		var distance float64
		if params.Longitude != "" && params.Latitude != "" {
			// 简单的距离计算，实际可能需要更精确的算法
			distKm, _ := strconv.ParseFloat(store.DistanceKM, 64)
			distance = distKm * 1000 // 转换为米
		}
		
		storeMap := map[string]interface{}{
			"store_code":     store.StoreCode,
			"store_name":     store.Name,
			"address":        store.Address,
			"phone":          "", // API没有返回电话
			"business_hours": fmt.Sprintf("%s - %s", store.StartTime, store.EndTime),
			"is_open":        store.OfficialStatus,
			"latitude":       store.Lat,
			"longitude":      store.Lng,
			"distance":       distance,
			"features":       []string{}, // API没有返回特性
			"city_id":        store.CityID,
			"city":           store.City,
			"close_note":     store.CloseNote,
		}
		stores = append(stores, storeMap)
	}

	return stores, totalStores, nil
}

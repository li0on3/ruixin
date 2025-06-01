package services

import (
	"backend/internal/models"
	"backend/internal/repository"
	"backend/pkg/httpclient"
	"fmt"
	"go.uber.org/zap"
)

type CityService struct {
	cityRepo     *repository.CityRepository
	cardRepo     *repository.CardRepository
	luckinClient *httpclient.LuckinClient
	logger       *zap.Logger
}

func NewCityService(
	cityRepo *repository.CityRepository,
	cardRepo *repository.CardRepository,
	luckinClient *httpclient.LuckinClient,
	logger *zap.Logger,
) *CityService {
	return &CityService{
		cityRepo:     cityRepo,
		cardRepo:     cardRepo,
		luckinClient: luckinClient,
		logger:       logger,
	}
}

// SyncCities 从瑞幸API同步城市数据
func (s *CityService) SyncCities(cardCode string) error {
	// 验证卡片
	card, err := s.cardRepo.GetByCode(cardCode)
	if err != nil {
		return err
	}

	// 调用瑞幸API获取城市列表
	// 使用卡片的瑞幸产品ID
	productID := 6 // 默认值
	if card.LuckinProductID > 0 {
		productID = card.LuckinProductID
	}
	
	req := &httpclient.CityByCardRequest{
		ProductID: productID,
		Card:      cardCode,
	}

	resp, err := s.luckinClient.CityByCard(req)
	if err != nil {
		s.logger.Error("failed to get cities from luckin", zap.Error(err), zap.Int("productId", productID))
		return err
	}

	if resp.Code != 200 {
		s.logger.Error("luckin api returned error", zap.Int("code", resp.Code), zap.String("msg", resp.Msg), zap.String("status", resp.Status))
		return fmt.Errorf("luckin api error: %s", resp.Msg)
	}

	// 批量更新城市数据
	var cities []*models.City
	for _, luckinCity := range resp.Data {
		city := &models.City{
			CityID:   luckinCity.CityID,
			CityName: luckinCity.CityName,
			Pinyin:   luckinCity.CityPinyin,
		}
		cities = append(cities, city)
	}

	if err := s.cityRepo.BatchCreateOrUpdate(cities); err != nil {
		s.logger.Error("failed to update cities", zap.Error(err))
		return err
	}

	s.logger.Info("cities synced successfully", zap.Int("count", len(cities)))
	return nil
}

// GetAllCities 获取所有城市
func (s *CityService) GetAllCities() ([]*models.City, error) {
	return s.cityRepo.List()
}

// GetCityByID 根据ID获取城市
func (s *CityService) GetCityByID(cityID int) (*models.City, error) {
	return s.cityRepo.GetByCityID(cityID)
}

// GetCityByName 根据名称获取城市
func (s *CityService) GetCityByName(cityName string) (*models.City, error) {
	return s.cityRepo.GetByCityName(cityName)
}

// SearchCities 搜索城市
func (s *CityService) SearchCities(keyword string) ([]*models.City, error) {
	return s.cityRepo.SearchByName(keyword)
}

// ConvertCityNameToID 将城市名称转换为ID
func (s *CityService) ConvertCityNameToID(cityName string) (int, error) {
	city, err := s.cityRepo.GetByCityName(cityName)
	if err != nil {
		return 0, err
	}
	return city.CityID, nil
}
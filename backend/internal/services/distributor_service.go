package services

import (
	"backend/internal/models"
	"backend/internal/repository"
	"errors"
)

type DistributorService struct {
	distributorRepo *repository.DistributorRepository
}

func NewDistributorService(distributorRepo *repository.DistributorRepository) *DistributorService {
	return &DistributorService{
		distributorRepo: distributorRepo,
	}
}

func (s *DistributorService) GetByID(id uint) (*models.Distributor, error) {
	return s.distributorRepo.GetByID(id)
}

func (s *DistributorService) GetByAPIKey(apiKey string) (*models.Distributor, error) {
	return s.distributorRepo.GetByAPIKey(apiKey)
}

func (s *DistributorService) GetByEmail(email string) (*models.Distributor, error) {
	return s.distributorRepo.GetByEmail(email)
}

func (s *DistributorService) Create(distributor *models.Distributor) error {
	// 检查邮箱是否已存在
	if _, err := s.distributorRepo.GetByEmail(distributor.Email); err == nil {
		return errors.New("email already exists")
	}

	return s.distributorRepo.Create(distributor)
}

func (s *DistributorService) Update(distributor *models.Distributor) error {
	// 如果更新邮箱，检查是否已存在
	if distributor.Email != "" {
		existing, err := s.distributorRepo.GetByEmail(distributor.Email)
		if err == nil && existing.ID != distributor.ID {
			return errors.New("email already exists")
		}
	}

	return s.distributorRepo.Update(distributor)
}

func (s *DistributorService) Delete(id uint) error {
	// 检查是否有未完成的订单
	distributor, err := s.distributorRepo.GetByID(id)
	if err != nil {
		return err
	}

	if distributor.TotalOrders > 0 {
		return errors.New("cannot delete distributor with order history")
	}

	return s.distributorRepo.Delete(id)
}

func (s *DistributorService) List(offset, limit int, filters map[string]interface{}) ([]*models.Distributor, int64, error) {
	return s.distributorRepo.List(offset, limit, filters)
}

func (s *DistributorService) UpdateBalance(id uint, amount float64) error {
	return s.distributorRepo.UpdateBalance(id, amount)
}

func (s *DistributorService) IncrementOrderCount(id uint) error {
	return s.distributorRepo.IncrementOrderCount(id)
}

func (s *DistributorService) GetAPILogs(distributorID uint, offset, limit int) ([]*models.DistributorAPILog, int64, error) {
	return s.distributorRepo.GetAPILogs(distributorID, offset, limit)
}

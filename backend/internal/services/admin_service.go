package services

import (
	"backend/internal/models"
	"backend/internal/repository"
	"errors"
	"golang.org/x/crypto/bcrypt"
)

type AdminService struct {
	adminRepo *repository.AdminRepository
}

func NewAdminService(adminRepo *repository.AdminRepository) *AdminService {
	return &AdminService{
		adminRepo: adminRepo,
	}
}

func (s *AdminService) GetByID(id uint) (*models.Admin, error) {
	return s.adminRepo.GetByID(id)
}

func (s *AdminService) GetByUsername(username string) (*models.Admin, error) {
	return s.adminRepo.GetByUsername(username)
}

func (s *AdminService) Create(admin *models.Admin) error {
	// 检查用户名是否已存在
	if _, err := s.adminRepo.GetByUsername(admin.Username); err == nil {
		return errors.New("username already exists")
	}

	// 检查邮箱是否已存在
	if admin.Email != "" {
		if _, err := s.adminRepo.GetByEmail(admin.Email); err == nil {
			return errors.New("email already exists")
		}
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(admin.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	admin.Password = string(hashedPassword)

	return s.adminRepo.Create(admin)
}

func (s *AdminService) Update(admin *models.Admin) error {
	// 如果更新密码，需要重新加密
	if admin.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(admin.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		admin.Password = string(hashedPassword)
	}

	return s.adminRepo.Update(admin)
}

func (s *AdminService) Delete(id uint) error {
	// 不能删除最后一个超级管理员
	admin, err := s.adminRepo.GetByID(id)
	if err != nil {
		return err
	}

	if admin.Role == "super_admin" {
		// 检查是否还有其他超级管理员
		admins, _, err := s.adminRepo.List(0, 100, map[string]interface{}{"role": "super_admin"})
		if err != nil {
			return err
		}
		if len(admins) <= 1 {
			return errors.New("cannot delete the last super admin")
		}
	}

	return s.adminRepo.Delete(id)
}

func (s *AdminService) List(offset, limit int, filters map[string]interface{}) ([]*models.Admin, int64, error) {
	return s.adminRepo.List(offset, limit, filters)
}

func (s *AdminService) UpdateLoginInfo(id uint, ip string) error {
	return s.adminRepo.UpdateLoginInfo(id, ip)
}

func (s *AdminService) ChangePassword(id uint, newPassword string) error {
	admin, err := s.adminRepo.GetByID(id)
	if err != nil {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	admin.Password = string(hashedPassword)
	return s.adminRepo.Update(admin)
}

func (s *AdminService) LogOperation(log *models.AdminOperationLog) error {
	return s.adminRepo.LogOperation(log)
}

func (s *AdminService) GetOperationLogs(adminID uint, offset, limit int) ([]*models.AdminOperationLog, int64, error) {
	return s.adminRepo.GetOperationLogs(adminID, offset, limit)
}

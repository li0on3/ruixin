package services

import (
	"backend/internal/models"
	"backend/internal/repository"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"go.uber.org/zap"
)

type BatchOrderService struct {
	db             *gorm.DB
	orderService   *OrderService
	cardService    *CardService
	financeService *FinanceService
	logger         *zap.Logger
}

func NewBatchOrderService(
	db *gorm.DB,
	orderService *OrderService,
	cardService *CardService,
	financeService *FinanceService,
	logger *zap.Logger,
) *BatchOrderService {
	return &BatchOrderService{
		db:             db,
		orderService:   orderService,
		cardService:    cardService,
		financeService: financeService,
		logger:         logger,
	}
}

type BatchOrderRequest struct {
	Orders []CreateOrderRequest `json:"orders"`
}

type BatchOrderResult struct {
	Success      bool                `json:"success"`
	TotalOrders  int                 `json:"total_orders"`
	SuccessCount int                 `json:"success_count"`
	FailedCount  int                 `json:"failed_count"`
	Orders       []*models.Order     `json:"orders,omitempty"`
	Errors       []string            `json:"errors,omitempty"`
}

// CreateBatchOrders 批量创建订单 - 全部成功或全部失败
func (s *BatchOrderService) CreateBatchOrders(distributorID uint, requests []CreateOrderRequest) (*BatchOrderResult, error) {
	result := &BatchOrderResult{
		TotalOrders: len(requests),
		Orders:      make([]*models.Order, 0),
		Errors:      make([]string, 0),
	}

	if len(requests) == 0 {
		return result, errors.New("批量订单不能为空")
	}

	// 1. 预检查阶段 - 检查所有订单是否都能创建
	reservedCards := make([]*models.Card, 0)
	totalAmount := 0.0
	
	// 开启事务
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			// 释放所有预占的卡片
			for _, card := range reservedCards {
				s.cardService.ReleaseCard(card.ID)
			}
		}
	}()

	// 2. 预占所有需要的卡片
	productService := NewProductService(tx, nil)
	priceRepo := repository.NewLuckinConfigRepository(tx)
	
	for i, req := range requests {
		req.DistributorID = distributorID
		
		// 如果没有指定卡片，自动选择
		if req.CardCode == "" {
			if len(req.Goods) == 0 {
				tx.Rollback()
				// 释放已预占的卡片
				for _, card := range reservedCards {
					s.cardService.ReleaseCard(card.ID)
				}
				return nil, fmt.Errorf("订单%d: 商品信息不能为空", i+1)
			}
			
			// 自动选卡
			card, err := productService.FindBestCard(
				req.Goods[0].GoodsCode, 
				s.cardService,
				priceRepo,
			)
			if err != nil {
				tx.Rollback()
				// 释放已预占的卡片
				for _, card := range reservedCards {
					s.cardService.ReleaseCard(card.ID)
				}
				return nil, fmt.Errorf("订单%d: %v", i+1, err)
			}
			
			req.CardCode = card.CardCode
			reservedCards = append(reservedCards, card)
			
			// 预占卡片
			if err := s.cardService.ReserveCard(card.ID); err != nil {
				tx.Rollback()
				// 释放已预占的卡片
				for _, c := range reservedCards {
					if c.ID != card.ID {
						s.cardService.ReleaseCard(c.ID)
					}
				}
				return nil, fmt.Errorf("订单%d: 预占卡片失败", i+1)
			}
		}
		
		// 预估订单金额（这里简化处理，实际应该调用检查接口）
		// TODO: 调用瑞幸API获取准确金额
		totalAmount += 20.0 // 假设每单20元
	}

	// 3. 检查余额是否充足
	distributor, err := repository.NewDistributorRepository(tx).GetByID(distributorID)
	if err != nil {
		tx.Rollback()
		// 释放所有预占的卡片
		for _, card := range reservedCards {
			s.cardService.ReleaseCard(card.ID)
		}
		return nil, errors.New("分销商不存在")
	}

	if distributor.Balance < totalAmount {
		tx.Rollback()
		// 释放所有预占的卡片
		for _, card := range reservedCards {
			s.cardService.ReleaseCard(card.ID)
		}
		return nil, errors.New("余额不足")
	}

	// 4. 创建所有订单
	successOrders := make([]*models.Order, 0)
	for i, req := range requests {
		order, err := s.orderService.CreateOrder(&req)
		if err != nil {
			// 如果任何一个订单失败，回滚所有操作
			tx.Rollback()
			
			// 释放所有预占的卡片
			for _, card := range reservedCards {
				s.cardService.ReleaseCard(card.ID)
			}
			
			// 退还已扣的款项
			for _, o := range successOrders {
				s.financeService.Refund(
					int64(distributorID), 
					o.TotalAmount, 
					o.OrderNo,
					"批量订单失败，自动退款",
					1,
				)
			}
			
			return nil, fmt.Errorf("订单%d创建失败: %v", i+1, err)
		}
		
		successOrders = append(successOrders, order)
	}

	// 5. 提交事务
	if err := tx.Commit().Error; err != nil {
		// 释放所有预占的卡片
		for _, card := range reservedCards {
			s.cardService.ReleaseCard(card.ID)
		}
		
		// 退还已扣的款项
		for _, o := range successOrders {
			s.financeService.Refund(
				int64(distributorID), 
				o.TotalAmount, 
				o.OrderNo,
				"批量订单提交失败，自动退款",
				1,
			)
		}
		
		return nil, errors.New("提交批量订单失败")
	}

	// 6. 返回成功结果
	result.Success = true
	result.SuccessCount = len(successOrders)
	result.Orders = successOrders
	
	return result, nil
}
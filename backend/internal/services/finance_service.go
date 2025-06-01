package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"backend/internal/models"
	"backend/internal/repository"
	"gorm.io/gorm"
)

type FinanceService struct {
	db              *gorm.DB
	transactionRepo *repository.TransactionRepository
	withdrawalRepo  *repository.WithdrawalRepository
	distributorRepo *repository.DistributorRepository
}

func NewFinanceService(db *gorm.DB) *FinanceService {
	return &FinanceService{
		db:              db,
		transactionRepo: repository.NewTransactionRepository(db),
		withdrawalRepo:  repository.NewWithdrawalRepository(db),
		distributorRepo: repository.NewDistributorRepository(db),
	}
}

// Recharge 充值
func (s *FinanceService) Recharge(distributorID int64, amount float64, remark string, createdBy int64) error {
	if amount <= 0 {
		return errors.New("充值金额必须大于0")
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		// 获取当前余额
		balance, frozenAmount, err := s.transactionRepo.GetDistributorBalance(distributorID)
		if err != nil {
			return err
		}

		// 创建交易记录
		transaction := &models.Transaction{
			DistributorID: distributorID,
			Type:          models.TransactionTypeRecharge,
			Amount:        amount,
			BalanceBefore: balance,
			BalanceAfter:  balance + amount,
			Remark:        remark,
			CreatedBy:     createdBy,
		}

		if err := tx.Create(transaction).Error; err != nil {
			return err
		}

		// 更新余额
		return s.transactionRepo.UpdateDistributorBalance(tx, distributorID, balance+amount, frozenAmount)
	})
}

// Consume 消费（下单时扣款）
func (s *FinanceService) Consume(distributorID int64, amount float64, orderID string, createdBy int64) error {
	if amount <= 0 {
		return errors.New("消费金额必须大于0")
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		// 获取当前余额
		balance, frozenAmount, err := s.transactionRepo.GetDistributorBalance(distributorID)
		if err != nil {
			return err
		}

		// 检查余额是否充足
		availableBalance := balance - frozenAmount
		if availableBalance < amount {
			return fmt.Errorf("余额不足，可用余额：%.2f，所需金额：%.2f", availableBalance, amount)
		}

		// 创建交易记录
		transaction := &models.Transaction{
			DistributorID: distributorID,
			Type:          models.TransactionTypeConsume,
			Amount:        amount,
			BalanceBefore: balance,
			BalanceAfter:  balance - amount,
			RelatedID:     orderID,
			Remark:        fmt.Sprintf("订单消费：%s", orderID),
			CreatedBy:     createdBy,
		}

		if err := tx.Create(transaction).Error; err != nil {
			return err
		}

		// 更新余额
		err = s.transactionRepo.UpdateDistributorBalance(tx, distributorID, balance-amount, frozenAmount)
		if err != nil {
			return err
		}

		// 检查余额预警
		go s.checkBalanceWarning(distributorID, balance-amount)

		return nil
	})
}

// Refund 退款
func (s *FinanceService) Refund(distributorID int64, amount float64, orderID string, remark string, createdBy int64) error {
	if amount <= 0 {
		return errors.New("退款金额必须大于0")
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		// 获取当前余额
		balance, frozenAmount, err := s.transactionRepo.GetDistributorBalance(distributorID)
		if err != nil {
			return err
		}

		// 创建交易记录
		transaction := &models.Transaction{
			DistributorID: distributorID,
			Type:          models.TransactionTypeRefund,
			Amount:        amount,
			BalanceBefore: balance,
			BalanceAfter:  balance + amount,
			RelatedID:     orderID,
			Remark:        remark,
			CreatedBy:     createdBy,
		}

		if err := tx.Create(transaction).Error; err != nil {
			return err
		}

		// 更新余额
		return s.transactionRepo.UpdateDistributorBalance(tx, distributorID, balance+amount, frozenAmount)
	})
}

// AdjustBalance 余额调整
func (s *FinanceService) AdjustBalance(distributorID int64, amount float64, remark string, createdBy int64) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 获取当前余额
		balance, frozenAmount, err := s.transactionRepo.GetDistributorBalance(distributorID)
		if err != nil {
			return err
		}

		// 检查调整后余额是否为负
		newBalance := balance + amount
		if newBalance < 0 {
			return fmt.Errorf("调整后余额不能为负数，当前余额：%.2f，调整金额：%.2f", balance, amount)
		}

		// 创建交易记录
		transaction := &models.Transaction{
			DistributorID: distributorID,
			Type:          models.TransactionTypeAdjust,
			Amount:        amount,
			BalanceBefore: balance,
			BalanceAfter:  newBalance,
			Remark:        remark,
			CreatedBy:     createdBy,
		}

		if err := tx.Create(transaction).Error; err != nil {
			return err
		}

		// 更新余额
		return s.transactionRepo.UpdateDistributorBalance(tx, distributorID, newBalance, frozenAmount)
	})
}

// CreateWithdrawal 创建提现申请
func (s *FinanceService) CreateWithdrawal(distributorID int64, amount float64, accountInfo models.AccountInfo, remark string) (*models.Withdrawal, error) {
	if amount <= 0 {
		return nil, errors.New("提现金额必须大于0")
	}

	// 序列化账户信息
	accountInfoJSON, err := json.Marshal(accountInfo)
	if err != nil {
		return nil, err
	}

	withdrawal := &models.Withdrawal{
		DistributorID: distributorID,
		Amount:        amount,
		Status:        models.WithdrawalStatusPending,
		AccountInfo:   string(accountInfoJSON),
		Remark:        remark,
	}

	err = s.db.Transaction(func(tx *gorm.DB) error {
		// 获取当前余额
		balance, frozenAmount, err := s.transactionRepo.GetDistributorBalance(distributorID)
		if err != nil {
			return err
		}

		// 检查余额是否充足
		availableBalance := balance - frozenAmount
		if availableBalance < amount {
			return fmt.Errorf("余额不足，可用余额：%.2f，提现金额：%.2f", availableBalance, amount)
		}

		// 创建提现申请
		if err := tx.Create(withdrawal).Error; err != nil {
			return err
		}

		// 冻结金额
		return s.transactionRepo.UpdateDistributorBalance(tx, distributorID, balance, frozenAmount+amount)
	})

	return withdrawal, err
}

// ProcessWithdrawal 处理提现申请
func (s *FinanceService) ProcessWithdrawal(withdrawalID int64, approved bool, rejectReason string, processorID int64) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 获取提现申请
		withdrawal, err := s.withdrawalRepo.GetByID(withdrawalID)
		if err != nil {
			return err
		}

		if withdrawal.Status != models.WithdrawalStatusPending {
			return errors.New("该提现申请已处理")
		}

		// 获取当前余额
		balance, frozenAmount, err := s.transactionRepo.GetDistributorBalance(withdrawal.DistributorID)
		if err != nil {
			return err
		}

		if approved {
			// 批准提现
			withdrawal.Status = models.WithdrawalStatusProcessed

			// 创建提现交易记录
			transaction := &models.Transaction{
				DistributorID: withdrawal.DistributorID,
				Type:          models.TransactionTypeWithdraw,
				Amount:        withdrawal.Amount,
				BalanceBefore: balance,
				BalanceAfter:  balance - withdrawal.Amount,
				RelatedID:     fmt.Sprintf("W%d", withdrawalID),
				Remark:        fmt.Sprintf("提现申请：%s", withdrawal.Remark),
				CreatedBy:     processorID,
			}

			if err := tx.Create(transaction).Error; err != nil {
				return err
			}

			// 更新余额（扣除余额，解冻金额）
			if err := s.transactionRepo.UpdateDistributorBalance(tx, withdrawal.DistributorID, balance-withdrawal.Amount, frozenAmount-withdrawal.Amount); err != nil {
				return err
			}
		} else {
			// 拒绝提现
			withdrawal.Status = models.WithdrawalStatusRejected
			withdrawal.RejectReason = rejectReason

			// 解冻金额
			if err := s.transactionRepo.UpdateDistributorBalance(tx, withdrawal.DistributorID, balance, frozenAmount-withdrawal.Amount); err != nil {
				return err
			}
		}

		// 更新提现申请
		now := time.Now()
		withdrawal.ProcessedAt = &now
		withdrawal.ProcessedBy = &processorID

		return tx.Save(withdrawal).Error
	})
}

// GetBalance 获取余额信息
func (s *FinanceService) GetBalance(distributorID int64) (map[string]interface{}, error) {
	balance, frozenAmount, err := s.transactionRepo.GetDistributorBalance(distributorID)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"balance":           balance,
		"frozen_amount":     frozenAmount,
		"available_balance": balance - frozenAmount,
	}, nil
}

// GetTransactionList 获取交易记录列表
func (s *FinanceService) GetTransactionList(filters map[string]interface{}, page, pageSize int) ([]*models.Transaction, int64, error) {
	offset := (page - 1) * pageSize
	return s.transactionRepo.ListByFilters(filters, offset, pageSize)
}

// GetWithdrawalList 获取提现申请列表
func (s *FinanceService) GetWithdrawalList(filters map[string]interface{}, page, pageSize int) ([]*models.Withdrawal, int64, error) {
	offset := (page - 1) * pageSize
	return s.withdrawalRepo.ListByFilters(filters, offset, pageSize)
}

// GetPendingWithdrawals 获取待处理的提现申请
func (s *FinanceService) GetPendingWithdrawals(page, pageSize int) ([]*models.Withdrawal, int64, error) {
	offset := (page - 1) * pageSize
	return s.withdrawalRepo.ListPending(offset, pageSize)
}

// GetFinanceStatistics 获取财务统计
func (s *FinanceService) GetFinanceStatistics() (map[string]interface{}, error) {
	pendingCount, _ := s.withdrawalRepo.GetPendingCount()
	pendingAmount, _ := s.withdrawalRepo.GetPendingAmount()

	// 获取今日充值总额
	var todayRecharge float64
	today := time.Now().Format("2006-01-02")
	s.db.Model(&models.Transaction{}).
		Where("type = ? AND DATE(created_at) = ?", models.TransactionTypeRecharge, today).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&todayRecharge)

	// 获取今日消费总额
	var todayConsume float64
	s.db.Model(&models.Transaction{}).
		Where("type = ? AND DATE(created_at) = ?", models.TransactionTypeConsume, today).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&todayConsume)

	// 获取今日提现总额
	var todayWithdraw float64
	s.db.Model(&models.Transaction{}).
		Where("type = ? AND DATE(created_at) = ?", models.TransactionTypeWithdraw, today).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&todayWithdraw)

	// 获取平台总余额
	var totalBalance float64
	s.db.Model(&models.Distributor{}).
		Select("COALESCE(SUM(balance), 0)").
		Scan(&totalBalance)

	// 获取平台总冻结金额
	var totalFrozen float64
	s.db.Model(&models.Distributor{}).
		Select("COALESCE(SUM(frozen_amount), 0)").
		Scan(&totalFrozen)

	return map[string]interface{}{
		"pending_withdrawals": map[string]interface{}{
			"count":  pendingCount,
			"amount": pendingAmount,
		},
		"today_summary": map[string]interface{}{
			"recharge": todayRecharge,
			"consume":  todayConsume,
			"withdraw": todayWithdraw,
		},
		"platform_balance": map[string]interface{}{
			"total":     totalBalance,
			"frozen":    totalFrozen,
			"available": totalBalance - totalFrozen,
		},
	}, nil
}

// checkBalanceWarning 检查余额预警
func (s *FinanceService) checkBalanceWarning(distributorID int64, currentBalance float64) {
	// 获取分销商信息
	distributor, err := s.distributorRepo.GetByID(uint(distributorID))
	if err != nil {
		return
	}

	// 检查是否启用预警
	if !distributor.WarningEnabled || distributor.WarningBalance <= 0 {
		return
	}

	// 检查余额是否低于预警阈值
	if currentBalance <= distributor.WarningBalance {
		// 发送预警通知
		s.sendBalanceWarning(distributor, currentBalance)
	}
}

// sendBalanceWarning 发送余额预警通知
func (s *FinanceService) sendBalanceWarning(distributor *models.Distributor, currentBalance float64) {
	// TODO: 实现邮件通知
	if distributor.WarningEmail != "" {
		// 发送邮件通知
	}

	// TODO: 实现Webhook通知
	if distributor.WarningWebhook != "" {
		// 发送Webhook通知
	}
}

// UpdateWarningSettings 更新预警设置
func (s *FinanceService) UpdateWarningSettings(distributorID int64, warningBalance float64, warningEnabled bool, warningEmail, warningWebhook string) error {
	updates := map[string]interface{}{
		"warning_balance": warningBalance,
		"warning_enabled": warningEnabled,
		"warning_email":   warningEmail,
		"warning_webhook": warningWebhook,
	}

	return s.db.Model(&models.Distributor{}).Where("id = ?", distributorID).Updates(updates).Error
}

// GetProfitReport 获取利润报表
func (s *FinanceService) GetProfitReport(startDate, endDate time.Time) (map[string]interface{}, error) {
	// 计算订单利润
	type OrderProfit struct {
		Count            int64   `json:"count"`
		TotalSalePrice   float64 `json:"total_sale_price"`
		TotalCostPrice   float64 `json:"total_cost_price"`
		TotalPlatformProfit float64 `json:"total_platform_profit"`
		TotalDistributorProfit float64 `json:"total_distributor_profit"`
	}

	var orderProfit OrderProfit
	err := s.db.Model(&models.Order{}).
		Select(`
			COUNT(*) as count,
			COALESCE(SUM(total_amount), 0) as total_sale_price,
			COALESCE(SUM(cost_amount), 0) as total_cost_price,
			COALESCE(SUM(profit_amount), 0) as total_platform_profit,
			COALESCE(SUM(total_amount - profit_amount), 0) as total_distributor_profit
		`).
		Where("status = ? AND created_at BETWEEN ? AND ?", models.OrderStatusSuccess, startDate, endDate).
		Scan(&orderProfit).Error
	if err != nil {
		return nil, err
	}

	// 计算卡片成本和利润
	type CardProfit struct {
		TotalCards       int64   `json:"total_cards"`
		UsedCards        int64   `json:"used_cards"`
		TotalCostPrice   float64 `json:"total_cost_price"`
		TotalSellPrice   float64 `json:"total_sell_price"`
		TotalProfit      float64 `json:"total_profit"`
	}

	var cardProfit CardProfit
	err = s.db.Model(&models.Card{}).
		Select(`
			COUNT(*) as total_cards,
			SUM(CASE WHEN status = 1 THEN 1 ELSE 0 END) as used_cards,
			COALESCE(SUM(cost_price), 0) as total_cost_price,
			COALESCE(SUM(CASE WHEN status = 1 THEN sell_price ELSE 0 END), 0) as total_sell_price,
			COALESCE(SUM(CASE WHEN status = 1 THEN sell_price - cost_price ELSE 0 END), 0) as total_profit
		`).
		Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Scan(&cardProfit).Error
	if err != nil {
		return nil, err
	}

	// 计算日均利润
	days := int(endDate.Sub(startDate).Hours() / 24)
	if days == 0 {
		days = 1
	}

	return map[string]interface{}{
		"period": map[string]interface{}{
			"start_date": startDate.Format("2006-01-02"),
			"end_date":   endDate.Format("2006-01-02"),
			"days":       days,
		},
		"order_profit": map[string]interface{}{
			"count":                  orderProfit.Count,
			"total_sale_price":       orderProfit.TotalSalePrice,
			"total_cost_price":       orderProfit.TotalCostPrice,
			"platform_profit":        orderProfit.TotalPlatformProfit,
			"distributor_profit":     orderProfit.TotalDistributorProfit,
			"avg_daily_profit":       orderProfit.TotalPlatformProfit / float64(days),
		},
		"card_profit": map[string]interface{}{
			"total_cards":      cardProfit.TotalCards,
			"used_cards":       cardProfit.UsedCards,
			"usage_rate":       float64(cardProfit.UsedCards) / float64(cardProfit.TotalCards) * 100,
			"total_cost":       cardProfit.TotalCostPrice,
			"total_revenue":    cardProfit.TotalSellPrice,
			"total_profit":     cardProfit.TotalProfit,
			"profit_margin":    cardProfit.TotalProfit / cardProfit.TotalCostPrice * 100,
		},
		"summary": map[string]interface{}{
			"total_profit":      orderProfit.TotalPlatformProfit,
			"avg_daily_profit":  orderProfit.TotalPlatformProfit / float64(days),
			"avg_order_profit":  orderProfit.TotalPlatformProfit / float64(orderProfit.Count),
		},
	}, nil
}

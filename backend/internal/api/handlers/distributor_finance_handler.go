package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"backend/internal/models"
	"backend/internal/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DistributorFinanceHandler struct {
	financeService *services.FinanceService
}

func NewDistributorFinanceHandler(db *gorm.DB) *DistributorFinanceHandler {
	return &DistributorFinanceHandler{
		financeService: services.NewFinanceService(db),
	}
}

// GetBalance 获取余额信息
func (h *DistributorFinanceHandler) GetBalance(c *gin.Context) {
	distributorID := c.GetInt64("distributor_id")

	balance, err := h.financeService.GetBalance(distributorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, balance)
}

// GetTransactionList 获取交易记录
func (h *DistributorFinanceHandler) GetTransactionList(c *gin.Context) {
	distributorID := c.GetInt64("distributor_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	filters := map[string]interface{}{
		"distributor_id": distributorID,
	}

	transactions, total, err := h.financeService.GetTransactionList(filters, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": transactions,
		"pagination": gin.H{
			"total":      total,
			"page":       page,
			"page_size":  pageSize,
			"total_page": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

// CreateWithdrawal 创建提现申请
func (h *DistributorFinanceHandler) CreateWithdrawal(c *gin.Context) {
	distributorID := c.GetInt64("distributor_id")

	var req struct {
		Amount      float64            `json:"amount" binding:"required,gt=0"`
		AccountInfo models.AccountInfo `json:"account_info" binding:"required"`
		Remark      string             `json:"remark"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 验证账户信息
	if req.AccountInfo.Type == "" || req.AccountInfo.AccountName == "" || req.AccountInfo.AccountNo == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "账户信息不完整"})
		return
	}

	if req.AccountInfo.Type == "bank" && (req.AccountInfo.BankName == "" || req.AccountInfo.BankBranch == "") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "银行卡信息不完整"})
		return
	}

	withdrawal, err := h.financeService.CreateWithdrawal(distributorID, req.Amount, req.AccountInfo, req.Remark)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "提现申请已提交",
		"data":    withdrawal,
	})
}

// GetWithdrawalList 获取提现记录
func (h *DistributorFinanceHandler) GetWithdrawalList(c *gin.Context) {
	distributorID := c.GetInt64("distributor_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	filters := map[string]interface{}{
		"distributor_id": distributorID,
	}

	withdrawals, total, err := h.financeService.GetWithdrawalList(filters, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": withdrawals,
		"pagination": gin.H{
			"total":      total,
			"page":       page,
			"page_size":  pageSize,
			"total_page": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

// UpdateWarningSettings 更新预警设置
func (h *DistributorFinanceHandler) UpdateWarningSettings(c *gin.Context) {
	distributorID := c.GetInt64("distributor_id")

	var req struct {
		WarningBalance float64 `json:"warning_balance" binding:"min=0"`
		WarningEnabled bool    `json:"warning_enabled"`
		WarningEmail   string  `json:"warning_email"`
		WarningWebhook string  `json:"warning_webhook"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 验证邮箱格式
	if req.WarningEmail != "" && !isValidEmail(req.WarningEmail) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "邮箱格式不正确"})
		return
	}

	// 验证Webhook URL
	if req.WarningWebhook != "" && !isValidURL(req.WarningWebhook) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Webhook URL格式不正确"})
		return
	}

	err := h.financeService.UpdateWarningSettings(
		distributorID,
		req.WarningBalance,
		req.WarningEnabled,
		req.WarningEmail,
		req.WarningWebhook,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "预警设置更新成功"})
}

// 辅助函数：验证邮箱格式
func isValidEmail(email string) bool {
	// 简单的邮箱格式验证
	return len(email) > 3 && strings.Contains(email, "@") && strings.Contains(email, ".")
}

// 辅助函数：验证URL格式
func isValidURL(url string) bool {
	// 简单的URL格式验证
	return strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")
}

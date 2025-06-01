package handlers

import (
	"net/http"
	"strconv"
	"time"

	"backend/internal/models"
	"backend/internal/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AdminFinanceHandler struct {
	financeService *services.FinanceService
}

func NewAdminFinanceHandler(db *gorm.DB) *AdminFinanceHandler {
	return &AdminFinanceHandler{
		financeService: services.NewFinanceService(db),
	}
}

// Recharge 充值
func (h *AdminFinanceHandler) Recharge(c *gin.Context) {
	var req struct {
		DistributorID int64   `json:"distributor_id" binding:"required"`
		Amount        float64 `json:"amount" binding:"required,gt=0"`
		Remark        string  `json:"remark"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	adminID, exists := c.Get("admin_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "管理员身份验证失败"})
		return
	}
	
	adminIDUint := adminID.(uint)
	adminIDInt64 := int64(adminIDUint)

	err := h.financeService.Recharge(req.DistributorID, req.Amount, req.Remark, adminIDInt64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "充值成功", "data": nil})
}

// AdjustBalance 余额调整
func (h *AdminFinanceHandler) AdjustBalance(c *gin.Context) {
	var req struct {
		DistributorID int64   `json:"distributor_id" binding:"required"`
		Amount        float64 `json:"amount" binding:"required"`
		Remark        string  `json:"remark" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	adminID, exists := c.Get("admin_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "管理员身份验证失败"})
		return
	}
	
	adminIDUint := adminID.(uint)
	adminIDInt64 := int64(adminIDUint)

	err := h.financeService.AdjustBalance(req.DistributorID, req.Amount, req.Remark, adminIDInt64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "余额调整成功", "data": nil})
}

// GetTransactionList 获取交易记录列表
func (h *AdminFinanceHandler) GetTransactionList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	distributorID, _ := strconv.ParseInt(c.Query("distributor_id"), 10, 64)
	txType, _ := strconv.Atoi(c.Query("type"))
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	filters := make(map[string]interface{})
	if distributorID > 0 {
		filters["distributor_id"] = distributorID
	}
	if txType > 0 {
		filters["type"] = models.TransactionType(txType)
	}
	if startDate != "" {
		if t, err := time.Parse("2006-01-02", startDate); err == nil {
			filters["start_time"] = t
		}
	}
	if endDate != "" {
		if t, err := time.Parse("2006-01-02", endDate); err == nil {
			filters["end_time"] = t.Add(24 * time.Hour).Add(-1 * time.Second)
		}
	}

	transactions, total, err := h.financeService.GetTransactionList(filters, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "成功",
		"data": gin.H{
			"data": transactions,
			"pagination": gin.H{
				"total":      total,
				"page":       page,
				"page_size":  pageSize,
				"total_page": (total + int64(pageSize) - 1) / int64(pageSize),
			},
		},
	})
}

// GetWithdrawalList 获取提现申请列表
func (h *AdminFinanceHandler) GetWithdrawalList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	distributorID, _ := strconv.ParseInt(c.Query("distributor_id"), 10, 64)
	status, _ := strconv.Atoi(c.Query("status"))
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	filters := make(map[string]interface{})
	if distributorID > 0 {
		filters["distributor_id"] = distributorID
	}
	if status >= 0 {
		filters["status"] = models.WithdrawalStatus(status)
	}
	if startDate != "" {
		if t, err := time.Parse("2006-01-02", startDate); err == nil {
			filters["start_time"] = t
		}
	}
	if endDate != "" {
		if t, err := time.Parse("2006-01-02", endDate); err == nil {
			filters["end_time"] = t.Add(24 * time.Hour).Add(-1 * time.Second)
		}
	}

	withdrawals, total, err := h.financeService.GetWithdrawalList(filters, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "成功",
		"data": gin.H{
			"data": withdrawals,
			"pagination": gin.H{
				"total":      total,
				"page":       page,
				"page_size":  pageSize,
				"total_page": (total + int64(pageSize) - 1) / int64(pageSize),
			},
		},
	})
}

// GetPendingWithdrawals 获取待处理的提现申请
func (h *AdminFinanceHandler) GetPendingWithdrawals(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	withdrawals, total, err := h.financeService.GetPendingWithdrawals(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "成功",
		"data": gin.H{
			"data": withdrawals,
			"pagination": gin.H{
				"total":      total,
				"page":       page,
				"page_size":  pageSize,
				"total_page": (total + int64(pageSize) - 1) / int64(pageSize),
			},
		},
	})
}

// ProcessWithdrawal 处理提现申请
func (h *AdminFinanceHandler) ProcessWithdrawal(c *gin.Context) {
	withdrawalID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的提现ID"})
		return
	}

	var req struct {
		Approved     bool   `json:"approved" binding:"required"`
		RejectReason string `json:"reject_reason"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !req.Approved && req.RejectReason == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "拒绝提现必须填写理由"})
		return
	}

	adminID := c.GetInt64("admin_id")

	err = h.financeService.ProcessWithdrawal(withdrawalID, req.Approved, req.RejectReason, adminID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "提现申请处理成功", "data": nil})
}

// GetFinanceStatistics 获取财务统计
func (h *AdminFinanceHandler) GetFinanceStatistics(c *gin.Context) {
	stats, err := h.financeService.GetFinanceStatistics()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "Success", "data": stats})
}

// GetDistributorBalance 获取分销商余额
func (h *AdminFinanceHandler) GetDistributorBalance(c *gin.Context) {
	distributorID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的分销商ID"})
		return
	}

	balance, err := h.financeService.GetBalance(distributorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "Success", "data": balance})
}

// GetProfitReport 获取利润报表
func (h *AdminFinanceHandler) GetProfitReport(c *gin.Context) {
	startDateStr := c.DefaultQuery("start_date", time.Now().AddDate(0, -1, 0).Format("2006-01-02"))
	endDateStr := c.DefaultQuery("end_date", time.Now().Format("2006-01-02"))

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "无效的开始日期", "data": nil})
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "无效的结束日期", "data": nil})
		return
	}

	// 设置结束日期为当天的23:59:59
	endDate = endDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	report, err := h.financeService.GetProfitReport(startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "获取利润报表失败", "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "成功", "data": report})
}

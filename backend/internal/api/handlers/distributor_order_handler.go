package handlers

import (
	"net/http"
	"backend/internal/services"
	"github.com/gin-gonic/gin"
)

// CreateSimplifiedOrder 创建简化订单
func (h *DistributorHandler) CreateSimplifiedOrder(c *gin.Context) {
	var req services.SimplifiedOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "请求参数无效: " + err.Error(),
			"data": nil,
		})
		return
	}

	// 从上下文获取分销商ID
	distributorID, _ := c.Get("distributor_id")
	req.DistributorID = distributorID.(uint)

	order, err := h.orderService.CreateSimplifiedOrder(&req)
	if err != nil {
		// 根据错误类型返回不同的状态码
		statusCode := http.StatusInternalServerError
		if err.Error() == "余额不足" {
			statusCode = http.StatusBadRequest
		}
		
		c.JSON(statusCode, gin.H{
			"code": statusCode,
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "订单创建成功",
		"data": gin.H{
			"order_no":      order.OrderNo,
			"status":        order.Status,
			"total_amount":  order.TotalAmount,
			"store_name":    order.StoreName,
			"store_address": order.StoreAddress,
		},
	})
}
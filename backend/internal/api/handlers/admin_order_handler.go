package handlers

import (
	"net/http"
	"strconv"
	"time"

	"backend/internal/services"
	"github.com/gin-gonic/gin"
)

type AdminOrderHandler struct {
	orderService *services.OrderService
}

func NewAdminOrderHandler(orderService *services.OrderService) *AdminOrderHandler {
	return &AdminOrderHandler{
		orderService: orderService,
	}
}

// ListOrders 获取订单列表
func (h *AdminOrderHandler) ListOrders(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	// 构建过滤条件
	filters := make(map[string]interface{})

	if distributorID := c.Query("distributor_id"); distributorID != "" {
		if id, err := strconv.ParseUint(distributorID, 10, 32); err == nil {
			filters["distributor_id"] = uint(id)
		}
	}

	if status := c.Query("status"); status != "" {
		if s, err := strconv.Atoi(status); err == nil {
			filters["status"] = s
		}
	}

	if storeCode := c.Query("store_code"); storeCode != "" {
		filters["store_code"] = storeCode
	}

	if cardCode := c.Query("card_code"); cardCode != "" {
		filters["card_code"] = cardCode
	}

	if startDate := c.Query("start_date"); startDate != "" {
		if t, err := time.Parse("2006-01-02", startDate); err == nil {
			filters["start_date"] = t
		}
	}

	if endDate := c.Query("end_date"); endDate != "" {
		if t, err := time.Parse("2006-01-02", endDate); err == nil {
			filters["end_date"] = t.Add(24 * time.Hour)
		}
	}

	offset := (page - 1) * pageSize
	orders, total, err := h.orderService.List(offset, pageSize, filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "Failed to get order list",
			"data": nil,
		})
		return
	}

	// 处理订单列表，为每个订单生成二维码
	orderList := make([]map[string]interface{}, len(orders))
	for i, order := range orders {
		orderList[i] = BuildOrderResponse(order, h.orderService)
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "成功",
		"data": gin.H{
			"list":      orderList,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetOrder 获取订单详情
func (h *AdminOrderHandler) GetOrder(c *gin.Context) {
	orderNo := c.Param("orderNo")
	if orderNo == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "Order number is required",
			"data": nil,
		})
		return
	}

	order, err := h.orderService.GetOrderByNo(orderNo)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "订单不存在",
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "成功",
		"data": BuildOrderResponse(order, h.orderService),
	})
}

// RefreshOrderStatus 刷新订单状态
func (h *AdminOrderHandler) RefreshOrderStatus(c *gin.Context) {
	orderNo := c.Param("orderNo")
	if orderNo == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "Order number is required",
			"data": nil,
		})
		return
	}

	order, err := h.orderService.RefreshOrderStatus(orderNo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "刷新订单状态失败: " + err.Error(),
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "订单状态已刷新",
		"data": BuildOrderResponse(order, h.orderService),
	})
}

// GenerateQRCode 生成订单二维码
func (h *AdminOrderHandler) GenerateQRCode(c *gin.Context) {
	orderNo := c.Param("orderNo")
	if orderNo == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "订单号不能为空",
			"data": nil,
		})
		return
	}

	order, err := h.orderService.GetOrderByNo(orderNo)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "订单不存在",
			"data": nil,
		})
		return
	}

	if order.QRData == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "订单暂无二维码数据",
			"data": nil,
		})
		return
	}

	qrCodeImage, err := h.orderService.GenerateQRCode(order.QRData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "生成二维码失败: " + err.Error(),
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "二维码生成成功",
		"data": gin.H{
			"qr_code_image": qrCodeImage,
			"qr_data":       order.QRData,
			"take_code":     order.TakeCode,
		},
	})
}

// GetOrderStatistics 获取订单统计
func (h *AdminOrderHandler) GetOrderStatistics(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	distributorIDStr := c.Query("distributor_id")

	var start, end time.Time
	var err error

	if startDate != "" {
		start, err = time.Parse("2006-01-02", startDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 400,
				"msg":  "Invalid start date format",
				"data": nil,
			})
			return
		}
	} else {
		start = time.Now().AddDate(0, -1, 0) // 默认最近一个月
	}

	if endDate != "" {
		end, err = time.Parse("2006-01-02", endDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 400,
				"msg":  "Invalid end date format",
				"data": nil,
			})
			return
		}
		end = end.Add(24 * time.Hour)
	} else {
		end = time.Now()
	}

	var distributorID uint
	if distributorIDStr != "" {
		if id, err := strconv.ParseUint(distributorIDStr, 10, 32); err == nil {
			distributorID = uint(id)
		}
	}

	stats, err := h.orderService.GetStatistics(distributorID, start, end)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "Failed to get order statistics",
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "成功",
		"data": stats,
	})
}

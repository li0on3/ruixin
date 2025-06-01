package handlers

import (
	"backend/internal/models"
	"backend/internal/services"
)

// BuildOrderResponse 构建订单响应数据，包括生成二维码图片
func BuildOrderResponse(order *models.Order, orderService *services.OrderService) map[string]interface{} {
	response := map[string]interface{}{
		"order_no":         order.OrderNo,
		"status":           order.Status,
		"take_code":        order.TakeCode,
		"qr_data":          order.QRData,
		"total_amount":     order.TotalAmount,
		"cost_amount":      order.CostAmount,
		"profit_amount":    order.ProfitAmount,
		"luckin_price":     order.LuckinPrice,     // 瑞幸原始价格
		"luckin_cost_price": order.LuckinCostPrice, // 瑞幸原始成本
		"store_code":       order.StoreCode,
		"store_name":       order.StoreName,
		"store_address":    order.StoreAddress,
		"phone_number":     order.PhoneNumber,
		"goods":            order.Goods,
		"created_at":       order.CreatedAt,
		"updated_at":       order.UpdatedAt,
	}

	// 如果有二维码数据，生成二维码图片
	if order.QRData != "" && orderService != nil {
		if qrCodeImage, err := orderService.GenerateQRCode(order.QRData); err == nil {
			response["qr_code_image"] = qrCodeImage
		}
	}

	// 添加状态描述
	response["status_text"] = GetOrderStatusText(order.Status)

	return response
}

// GetOrderStatusText 获取订单状态的中文描述
func GetOrderStatusText(status models.OrderStatus) string {
	switch status {
	case models.OrderStatusPending:
		return "待处理"
	case models.OrderStatusDoing:
		return "处理中"
	case models.OrderStatusSuccess:
		return "已完成"
	case models.OrderStatusFailed:
		return "失败"
	case models.OrderStatusRefunded:
		return "已退款"
	case models.OrderStatusCancelled:
		return "已取消"
	default:
		return "未知"
	}
}
package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"backend/internal/models"
	"backend/internal/repository"
	"go.uber.org/zap"
)

type CallbackService struct {
	orderRepo *repository.OrderRepository
	logger    *zap.Logger
	client    *http.Client
}

func NewCallbackService(orderRepo *repository.OrderRepository, logger *zap.Logger) *CallbackService {
	return &CallbackService{
		orderRepo: orderRepo,
		logger:    logger,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

type CallbackRequest struct {
	OrderNo     string                 `json:"order_no"`
	Status      string                 `json:"status"`
	TakeCode    string                 `json:"take_code,omitempty"`
	QRData      string                 `json:"qr_data,omitempty"`
	TotalAmount float64                `json:"total_amount"`
	Timestamp   int64                  `json:"timestamp"`
	Extra       map[string]interface{} `json:"extra,omitempty"`
}

func (s *CallbackService) SendCallback(order *models.Order) error {
	if order.CallbackURL == "" {
		return nil
	}

	// 构建回调数据
	callbackReq := CallbackRequest{
		OrderNo:     order.OrderNo,
		Status:      s.getStatusString(order.Status),
		TakeCode:    order.TakeCode,
		QRData:      order.QRData,
		TotalAmount: order.TotalAmount,
		Timestamp:   time.Now().Unix(),
		Extra: map[string]interface{}{
			"store_name":    order.StoreName,
			"store_address": order.StoreAddress,
			"phone_number":  order.PhoneNumber,
		},
	}

	jsonData, err := json.Marshal(callbackReq)
	if err != nil {
		s.logger.Error("failed to marshal callback data",
			zap.String("order_no", order.OrderNo),
			zap.Error(err))
		return err
	}

	// 发送HTTP请求
	req, err := http.NewRequest("POST", order.CallbackURL, bytes.NewBuffer(jsonData))
	if err != nil {
		s.logger.Error("failed to create callback request",
			zap.String("order_no", order.OrderNo),
			zap.Error(err))
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Order-No", order.OrderNo)
	req.Header.Set("X-Timestamp", fmt.Sprintf("%d", callbackReq.Timestamp))

	resp, err := s.client.Do(req)
	if err != nil {
		s.logger.Error("failed to send callback",
			zap.String("order_no", order.OrderNo),
			zap.String("url", order.CallbackURL),
			zap.Error(err))

		// 更新回调状态为失败
		order.CallbackStatus = 2
		s.orderRepo.Update(order)
		return err
	}
	defer resp.Body.Close()

	// 更新回调状态
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		order.CallbackStatus = 1
		now := time.Now()
		order.CallbackTime = &now
		s.logger.Info("callback sent successfully",
			zap.String("order_no", order.OrderNo),
			zap.Int("status_code", resp.StatusCode))
	} else {
		order.CallbackStatus = 2
		s.logger.Warn("callback returned non-2xx status",
			zap.String("order_no", order.OrderNo),
			zap.Int("status_code", resp.StatusCode))
	}

	s.orderRepo.Update(order)
	return nil
}

func (s *CallbackService) getStatusString(status models.OrderStatus) string {
	switch status {
	case models.OrderStatusPending:
		return "pending"
	case models.OrderStatusDoing:
		return "processing"
	case models.OrderStatusSuccess:
		return "success"
	case models.OrderStatusFailed:
		return "failed"
	case models.OrderStatusRefunded:
		return "refunded"
	case models.OrderStatusCancelled:
		return "cancelled"
	default:
		return "unknown"
	}
}

// RetryFailedCallbacks 重试失败的回调
func (s *CallbackService) RetryFailedCallbacks() {
	// 查找需要重试的订单
	orders, _, err := s.orderRepo.List(0, 100, map[string]interface{}{
		"callback_status": 2,
		"status":          models.OrderStatusSuccess,
	})

	if err != nil {
		s.logger.Error("failed to get orders for callback retry", zap.Error(err))
		return
	}

	for _, order := range orders {
		// 限制重试次数
		if order.CallbackStatus >= 5 {
			continue
		}

		s.logger.Info("retrying callback",
			zap.String("order_no", order.OrderNo),
			zap.Int("attempt", order.CallbackStatus))

		go s.SendCallback(order)
	}
}

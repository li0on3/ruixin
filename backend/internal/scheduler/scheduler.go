package scheduler

import (
	"time"

	"backend/internal/models"
	"backend/internal/services"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Scheduler struct {
	db              *gorm.DB
	callbackService *services.CallbackService
	cardService     *services.CardService
	logger          *zap.Logger
	stopCh          chan struct{}
}

func NewScheduler(
	db *gorm.DB,
	callbackService *services.CallbackService,
	cardService *services.CardService,
	logger *zap.Logger,
) *Scheduler {
	return &Scheduler{
		db:              db,
		callbackService: callbackService,
		cardService:     cardService,
		logger:          logger,
		stopCh:          make(chan struct{}),
	}
}

func (s *Scheduler) Start() {
	s.logger.Info("Starting scheduler")

	// 回调重试任务 - 每5分钟执行一次
	go s.runPeriodic("callback-retry", 5*time.Minute, func() {
		s.logger.Info("Running callback retry task")
		s.callbackService.RetryFailedCallbacks()
	})

	// 卡片状态检查任务 - 每小时执行一次
	go s.runPeriodic("card-status-check", 1*time.Hour, func() {
		s.logger.Info("Running card status check task")
		s.checkCardStatus()
	})

	// 卡片验证任务 - 每5分钟执行一次
	go s.runPeriodic("card-validation", 5*time.Minute, func() {
		s.logger.Info("Running card validation task")
		s.validateCards()
	})

	// 数据清理任务 - 每天凌晨2点执行
	go s.runDaily("data-cleanup", 2, 0, func() {
		s.logger.Info("Running data cleanup task")
		s.cleanupOldData()
	})
}

func (s *Scheduler) Stop() {
	s.logger.Info("Stopping scheduler")
	close(s.stopCh)
}

func (s *Scheduler) runPeriodic(name string, interval time.Duration, task func()) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.logger.Debug("Running periodic task", zap.String("task", name))
			task()
		case <-s.stopCh:
			s.logger.Debug("Stopping periodic task", zap.String("task", name))
			return
		}
	}
}

func (s *Scheduler) runDaily(name string, hour, minute int, task func()) {
	for {
		now := time.Now()
		next := time.Date(now.Year(), now.Month(), now.Day(), hour, minute, 0, 0, now.Location())
		if next.Before(now) {
			next = next.Add(24 * time.Hour)
		}

		duration := next.Sub(now)
		s.logger.Info("Scheduling daily task",
			zap.String("task", name),
			zap.Time("next_run", next),
			zap.Duration("duration", duration))

		select {
		case <-time.After(duration):
			s.logger.Debug("Running daily task", zap.String("task", name))
			task()
		case <-s.stopCh:
			s.logger.Debug("Stopping daily task", zap.String("task", name))
			return
		}
	}
}

func (s *Scheduler) checkCardStatus() {
	// 清理预占超时的卡片（超过5分钟的预占自动释放）
	s.logger.Debug("Checking for expired card reservations")
	
	expiredTime := time.Now().Add(-5 * time.Minute)
	
	// 查询所有预占超时的卡片
	var expiredCards []models.Card
	err := s.db.Where("status = 2 AND reserved_at < ?", expiredTime).Find(&expiredCards).Error
	if err != nil {
		s.logger.Error("Failed to query expired reservations", zap.Error(err))
		return
	}
	
	// 释放超时的卡片
	for _, card := range expiredCards {
		if err := s.cardService.ReleaseCard(card.ID); err != nil {
			s.logger.Error("Failed to release expired card reservation",
				zap.Uint("card_id", card.ID),
				zap.String("card_code", card.CardCode),
				zap.Error(err))
		} else {
			s.logger.Info("Released expired card reservation",
				zap.Uint("card_id", card.ID),
				zap.String("card_code", card.CardCode))
		}
	}
	
	s.logger.Debug("Card reservation cleanup completed", 
		zap.Int("released_count", len(expiredCards)))
}

func (s *Scheduler) cleanupOldData() {
	// TODO: 实现数据清理逻辑
	// 1. 清理过期的API日志
	// 2. 清理过期的操作日志
	// 3. 归档历史订单数据
	s.logger.Info("Data cleanup completed")
}

func (s *Scheduler) validateCards() {
	// 批量验证并更新卡片状态
	result, err := s.cardService.BatchValidateAndUpdateCards()
	if err != nil {
		s.logger.Error("Failed to validate cards", zap.Error(err))
		return
	}
	
	s.logger.Info("Card validation completed",
		zap.Int("total", result.TotalCards),
		zap.Int("valid", result.ValidCards),
		zap.Int("invalid", result.InvalidCards),
		zap.Int("failed", result.FailedCards))
}

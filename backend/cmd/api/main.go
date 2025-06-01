package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"backend/internal/api/handlers"
	"backend/internal/api/middleware"
	"backend/internal/config"
	"backend/internal/models"
	"backend/internal/repository"
	"backend/internal/scheduler"
	"backend/internal/services"
	"backend/pkg/httpclient"
	"backend/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// 加载配置
	if err := config.Init("configs/config.yaml"); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化日志
	appLogger := logger.NewLogger(config.AppConfig.Log)

	// 初始化数据库
	db, err := initDB()
	if err != nil {
		appLogger.Fatal("Failed to connect to database", err)
	}

	// 自动迁移
	if err := migrateDB(db); err != nil {
		appLogger.Fatal("Failed to migrate database", err)
	}

	// 初始化Redis
	rdb := initRedis()
	
	// 初始化频率限制器
	rateLimiter := middleware.NewRateLimiter(rdb)

	// 初始化仓储层
	adminRepo := repository.NewAdminRepository(db)
	cardRepo := repository.NewCardRepository(db)
	cardBatchRepo := repository.NewCardBatchRepository(db)
	distributorRepo := repository.NewDistributorRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	cityRepo := repository.NewCityRepository(db)

	// 初始化HTTP客户端
	luckinClient := httpclient.NewLuckinClient(&config.AppConfig.Luckin)

	// 初始化服务层
	adminService := services.NewAdminService(adminRepo)
	luckinConfigService := services.NewLuckinConfigService(db)
	
	// 先创建产品和系统配置服务
	productService := services.NewProductService(db, luckinClient)
	systemConfigService := services.NewSystemConfigService(db)
	
	// 初始化安全配置
	if err := systemConfigService.InitializeSecurityConfigs(); err != nil {
		appLogger.Warn("Failed to initialize security configs", zap.Error(err))
	}
	
	// 然后创建卡片服务（依赖产品和系统配置服务）
	cardService := services.NewCardService(db, cardRepo, cardBatchRepo, luckinConfigService, productService, systemConfigService)
	
	distributorService := services.NewDistributorService(distributorRepo)
	cityService := services.NewCityService(cityRepo, cardRepo, luckinClient, appLogger.Logger)
	storeService := services.NewStoreService(cardRepo, cityRepo, luckinClient, appLogger.Logger)
	orderService := services.NewOrderService(db, orderRepo, cardRepo, distributorRepo, luckinClient, cardService, appLogger.Logger)
	dashboardService := services.NewDashboardService(orderRepo, cardRepo, distributorRepo)
	callbackService := services.NewCallbackService(orderRepo, appLogger.Logger)
	statisticsService := services.NewStatisticsService(orderRepo, distributorRepo)
	cacheService := services.NewCacheService()
	securityAuditService := services.NewSecurityAuditService(db)

	// 初始化处理器
	// 分销商处理器
	distributorHandler := handlers.NewDistributorHandler(orderService, storeService, cityService, cardService, securityAuditService, appLogger.Logger)
	distributorFinanceHandler := handlers.NewDistributorFinanceHandler(db)
	distributorAuthHandler := handlers.NewDistributorAuthHandler(distributorService, &config.AppConfig.JWT)
	
	// 管理员处理器
	adminAuthHandler := handlers.NewAdminAuthHandler(adminService, &config.AppConfig.JWT)
	adminCardHandler := handlers.NewAdminCardHandler(cardService, adminService, productService, systemConfigService)
	adminDistributorHandler := handlers.NewAdminDistributorHandler(distributorService, adminService)
	adminOrderHandler := handlers.NewAdminOrderHandler(orderService)
	adminDashboardHandler := handlers.NewAdminDashboardHandler(dashboardService)
	adminFinanceHandler := handlers.NewAdminFinanceHandler(db)
	adminLuckinConfigHandler := handlers.NewAdminLuckinConfigHandler(db)
	adminStatisticsHandler := handlers.NewAdminStatisticsHandler(statisticsService)
	adminSystemConfigHandler := handlers.NewAdminSystemConfigHandler(systemConfigService)
	
	// 可用商品处理器
	availableProductsHandler := handlers.NewAvailableProductsHandler(productService, cacheService)

	// 设置Gin模式
	gin.SetMode(config.AppConfig.Server.Mode)

	// 创建路由
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	
	// CORS中间件 - 动态允许端口3000的访问
	router.Use(func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		
		// 检查是否是端口3000的请求（不限制IP）
		if strings.Contains(origin, ":3000") || origin == "" {
			if origin != "" {
				c.Header("Access-Control-Allow-Origin", origin)
			} else {
				c.Header("Access-Control-Allow-Origin", "*")
			}
		}
		
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		c.Header("Access-Control-Allow-Credentials", "true")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	})

	// 健康检查
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// API路由组
	apiV1 := router.Group("/api/v1")

	// 分销商认证接口（公开）
	apiV1.POST("/distributor/login", distributorAuthHandler.Login)

	// 分销商API（需要认证）
	distributorAPI := apiV1.Group("/distributor")
	distributorAPI.Use(middleware.DistributorAuth(distributorRepo))
	
	// 应用基础频率限制：每分钟300次请求（宽松限制，不影响正常业务）
	distributorAPI.Use(middleware.DistributorRateLimit(rateLimiter, time.Minute, 300))
	{
		// 订单相关接口使用正常限制
		distributorAPI.POST("/order", distributorHandler.CreateOrder)
		distributorAPI.POST("/order/simplified", distributorHandler.CreateSimplifiedOrder)  // 简化下单接口
		distributorAPI.POST("/orders/batch", distributorHandler.BatchCreateOrders)
		distributorAPI.GET("/order/:orderNo", distributorHandler.QueryOrder)
		
		// 查询类接口使用更严格的限制，防止爬取
		sensitiveAPILimiter := middleware.SensitiveAPIRateLimit(rateLimiter)
		distributorAPI.GET("/stores", sensitiveAPILimiter, distributorHandler.SearchStores)
		distributorAPI.GET("/menu", sensitiveAPILimiter, distributorHandler.GetMenu)
		distributorAPI.GET("/goods", sensitiveAPILimiter, distributorHandler.GetGoodsDetail)
		
		// 城市相关
		distributorAPI.GET("/cities", distributorHandler.GetCities)
		distributorAPI.POST("/cities/sync", distributorHandler.SyncCities)

		// 财务相关
		distributorAPI.GET("/balance", distributorFinanceHandler.GetBalance)
		distributorAPI.GET("/transactions", distributorFinanceHandler.GetTransactionList)
		distributorAPI.POST("/withdrawal", distributorFinanceHandler.CreateWithdrawal)
		distributorAPI.GET("/withdrawals", distributorFinanceHandler.GetWithdrawalList)
		distributorAPI.PUT("/warning-settings", distributorFinanceHandler.UpdateWarningSettings)
		
		// 可用商品查询
		distributorAPI.GET("/available-products", availableProductsHandler.GetDistributorAvailableProducts)
	}

	// 分销商JWT认证的API
	distributorJWTAPI := apiV1.Group("/distributor")
	distributorJWTAPI.Use(middleware.DistributorJWTAuth(config.AppConfig.JWT.Secret))
	{
		distributorJWTAPI.POST("/logout", distributorAuthHandler.Logout)
		distributorJWTAPI.GET("/profile", distributorAuthHandler.GetProfile)
		distributorJWTAPI.PUT("/password", distributorAuthHandler.ChangePassword)
	}

	// 公开的管理员认证接口
	apiV1.POST("/admin/login", adminAuthHandler.Login)

	// 管理员API（需要认证）
	adminAPI := apiV1.Group("/admin")
	adminAPI.Use(middleware.AdminAuth(config.AppConfig.JWT.Secret))
	{
		// 用户相关
		adminAPI.POST("/logout", adminAuthHandler.Logout)
		adminAPI.GET("/user/info", adminAuthHandler.GetUserInfo)
		adminAPI.PUT("/user/password", adminAuthHandler.ChangePassword)

		// 仪表盘
		adminAPI.GET("/dashboard/statistics", adminDashboardHandler.GetStatistics)
		adminAPI.GET("/dashboard/order-trend", adminDashboardHandler.GetOrderTrend)
		adminAPI.GET("/dashboard/hot-goods", adminDashboardHandler.GetHotGoods)
		adminAPI.GET("/dashboard/recent-orders", adminDashboardHandler.GetRecentOrders)

		// 卡片管理
		adminAPI.GET("/cards", adminCardHandler.ListCards)
		adminAPI.GET("/cards/:id", adminCardHandler.GetCard)
		adminAPI.POST("/cards", adminCardHandler.CreateCard)
		adminAPI.PUT("/cards/:id", adminCardHandler.UpdateCard)
		adminAPI.DELETE("/cards/:id", adminCardHandler.DeleteCard)
		adminAPI.GET("/cards/:id/usage-logs", adminCardHandler.GetCardUsageLogs)
		adminAPI.POST("/cards/batch-import", adminCardHandler.BatchImportCards)
		adminAPI.GET("/cards/stats", adminCardHandler.GetCardStats)
		adminAPI.POST("/cards/:id/sync-products", adminCardHandler.SyncCardProducts)
		adminAPI.POST("/cards/validate", adminCardHandler.ValidateCard)
		adminAPI.POST("/cards/batch-validate", adminCardHandler.BatchValidateCards)
		
		// 新版双模式批量验证
		adminAPI.POST("/cards/batch-validation/start", adminCardHandler.StartBatchValidation)
		adminAPI.GET("/cards/batch-validation/:taskId", adminCardHandler.GetValidationProgress)
		adminAPI.DELETE("/cards/batch-validation/:taskId", adminCardHandler.CancelValidation)
		adminAPI.GET("/cards/validation-stats", adminCardHandler.GetValidationStatistics)
		
		// 卡片批次管理
		adminAPI.GET("/card-batches", adminCardHandler.ListBatches)
		adminAPI.GET("/card-batches/:id", adminCardHandler.GetBatch)
		adminAPI.GET("/card-batches/:id/cards", adminCardHandler.GetBatchCards)

		// 分销商管理
		adminAPI.GET("/distributors", adminDistributorHandler.ListDistributors)
		adminAPI.GET("/distributors/:id", adminDistributorHandler.GetDistributor)
		adminAPI.POST("/distributors", adminDistributorHandler.CreateDistributor)
		adminAPI.PUT("/distributors/:id", adminDistributorHandler.UpdateDistributor)
		adminAPI.DELETE("/distributors/:id", adminDistributorHandler.DeleteDistributor)
		adminAPI.POST("/distributors/:id/reset-api-key", adminDistributorHandler.ResetAPIKey)
		adminAPI.POST("/distributors/:id/reset-password", adminDistributorHandler.ResetPassword)
		adminAPI.GET("/distributors/:id/api-logs", adminDistributorHandler.GetAPILogs)

		// 订单管理
		adminAPI.GET("/orders", adminOrderHandler.ListOrders)
		adminAPI.GET("/orders/:orderNo", adminOrderHandler.GetOrder)
		adminAPI.POST("/orders/:orderNo/refresh", adminOrderHandler.RefreshOrderStatus)
		adminAPI.POST("/orders/:orderNo/qrcode", adminOrderHandler.GenerateQRCode)
		adminAPI.GET("/orders/statistics", adminOrderHandler.GetOrderStatistics)

		// 财务管理
		adminAPI.POST("/finance/recharge", adminFinanceHandler.Recharge)
		adminAPI.POST("/finance/adjust", adminFinanceHandler.AdjustBalance)
		adminAPI.GET("/finance/transactions", adminFinanceHandler.GetTransactionList)
		adminAPI.GET("/finance/withdrawals", adminFinanceHandler.GetWithdrawalList)
		adminAPI.GET("/finance/withdrawals/pending", adminFinanceHandler.GetPendingWithdrawals)
		adminAPI.POST("/finance/withdrawals/:id/process", adminFinanceHandler.ProcessWithdrawal)
		adminAPI.GET("/finance/statistics", adminFinanceHandler.GetFinanceStatistics)
		adminAPI.GET("/finance/distributors/:id/balance", adminFinanceHandler.GetDistributorBalance)
		adminAPI.GET("/finance/profit-report", adminFinanceHandler.GetProfitReport)

		// 瑞幸配置管理
		// 价格管理
		adminAPI.GET("/luckin/prices", adminLuckinConfigHandler.GetPriceList)
		adminAPI.POST("/luckin/prices", adminLuckinConfigHandler.CreatePrice)
		adminAPI.PUT("/luckin/prices/:id", adminLuckinConfigHandler.UpdatePrice)
		adminAPI.DELETE("/luckin/prices/:id", adminLuckinConfigHandler.DeletePrice)

		// 产品管理
		adminAPI.GET("/luckin/products", adminLuckinConfigHandler.GetProductList)
		adminAPI.POST("/luckin/products", adminLuckinConfigHandler.CreateProduct)
		adminAPI.PUT("/luckin/products/:id", adminLuckinConfigHandler.UpdateProduct)
		adminAPI.DELETE("/luckin/products/:id", adminLuckinConfigHandler.DeleteProduct)
		adminAPI.GET("/luckin/products/categories", adminLuckinConfigHandler.GetProductCategories)
		adminAPI.POST("/luckin/products/import", adminLuckinConfigHandler.BatchImportProducts)

		// 种类绑定管理
		adminAPI.GET("/luckin/categories/:id/bindings", adminLuckinConfigHandler.GetCategoryBindings)
		adminAPI.POST("/luckin/categories/:id/bindings", adminLuckinConfigHandler.CreateBinding)
		adminAPI.DELETE("/luckin/bindings/:bindingId", adminLuckinConfigHandler.DeleteBinding)
		adminAPI.PUT("/luckin/bindings/:bindingId/priority", adminLuckinConfigHandler.UpdateBindingPriority)

		// 获取有效选项
		adminAPI.GET("/luckin/active-options", adminLuckinConfigHandler.GetActiveOptions)
		
		// 商品管理（用于价格配置）
		adminProductHandler := handlers.NewAdminProductHandler(productService, systemConfigService)
		adminAPI.GET("/products", adminProductHandler.GetProductList)
		adminAPI.GET("/products/search", adminProductHandler.SearchProducts)
		adminAPI.POST("/products/by-codes", adminProductHandler.GetProductsByCodes)
		adminAPI.POST("/products/sync", adminProductHandler.SyncProducts)
		adminAPI.GET("/products/match-logs", adminProductHandler.GetMatchLogs)
		adminAPI.POST("/products/parse-specs", adminProductHandler.ParseSpecsCode)
		adminAPI.GET("/products/cards", adminProductHandler.GetProductCards)

		// 统计分析
		adminAPI.GET("/statistics/metrics", adminStatisticsHandler.GetMetrics)
		adminAPI.GET("/statistics/sales-trend", adminStatisticsHandler.GetSalesTrend)
		adminAPI.GET("/statistics/distributor-rank", adminStatisticsHandler.GetDistributorRank)
		adminAPI.GET("/statistics/product-analysis", adminStatisticsHandler.GetProductAnalysis)
		adminAPI.GET("/statistics/hour-distribution", adminStatisticsHandler.GetHourDistribution)
		adminAPI.GET("/statistics/region-distribution", adminStatisticsHandler.GetRegionDistribution)
		adminAPI.GET("/statistics/details", adminStatisticsHandler.GetDetailData)
		adminAPI.GET("/statistics/export", adminStatisticsHandler.ExportData)
		
		// 系统配置
		adminAPI.GET("/system/configs", adminSystemConfigHandler.GetConfigs)
		adminAPI.PUT("/system/configs", adminSystemConfigHandler.UpdateConfigs)
		adminAPI.GET("/system/store-code", adminSystemConfigHandler.GetStoreCode)
		adminAPI.PUT("/system/store-code", adminSystemConfigHandler.UpdateStoreCode)
		
		// 店铺查询
		adminStoreHandler := handlers.NewAdminStoreHandler(storeService, cityService, cardService)
		adminAPI.GET("/stores/search", adminStoreHandler.SearchStores)
		adminAPI.GET("/stores/cities", adminStoreHandler.GetCities)
		adminAPI.GET("/stores/hot", adminStoreHandler.GetHotStores)
		adminAPI.POST("/stores/sync-cities", adminStoreHandler.SyncCities)
		
		// 可用商品查询（管理员）
		adminAPI.GET("/available-products", availableProductsHandler.GetAdminAvailableProducts)
	}

	// 创建服务器
	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", config.AppConfig.Server.Host, config.AppConfig.Server.Port),
		Handler:      router,
		ReadTimeout:  time.Duration(config.AppConfig.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(config.AppConfig.Server.WriteTimeout) * time.Second,
	}

	// 启动调度器
	taskScheduler := scheduler.NewScheduler(db, callbackService, cardService, appLogger.Logger)
	taskScheduler.Start()

	// 启动服务器
	go func() {
		appLogger.Info("Server starting",
			appLogger.String("address", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			appLogger.Fatal("Failed to start server", err)
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	appLogger.Info("Shutting down server...")

	// 优雅关闭
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		appLogger.Fatal("Server forced to shutdown", err)
	}

	// 停止调度器
	taskScheduler.Stop()

	// 关闭数据库连接
	sqlDB, _ := db.DB()
	sqlDB.Close()

	// 关闭Redis连接
	rdb.Close()

	appLogger.Info("Server exited")
}

func initDB() (*gorm.DB, error) {
	dsn := config.GetDSN()
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(config.AppConfig.Database.MaxIdle)
	sqlDB.SetMaxOpenConns(config.AppConfig.Database.MaxOpen)

	return db, nil
}

func migrateDB(db *gorm.DB) error {
	// 禁用外键约束以避免迁移时的问题
	db.Exec("SET FOREIGN_KEY_CHECKS = 0")
	defer db.Exec("SET FOREIGN_KEY_CHECKS = 1")
	
	return db.AutoMigrate(
		// 基础表（无外键依赖）
		&models.Admin{},
		&models.Distributor{},
		&models.LuckinPrice{},
		&models.LuckinProduct{},
		&models.City{},
		&models.Product{},
		&models.ProductSKU{},
		&models.ProductSpec{},
		&models.ProductSpecOption{},
		&models.ProductPriceMapping{},
		&models.ProductAlias{},
		&models.ProductSKUMapping{},
		&models.ProductMatchLog{},
		&models.SystemConfig{},
		&models.ProductSpecConfig{},
		
		// 依赖基础表的表
		&models.AdminOperationLog{},
		&models.CardBatch{},      // 依赖 LuckinPrice 和 Admin
		&models.Card{},           // 依赖 CardBatch 和 LuckinPrice
		&models.CardProductBinding{}, // 依赖 Card 和 Product
		&models.Order{},          // 依赖 Distributor 和 Card
		&models.CardUsageLog{},   // 依赖 Card 和 Distributor
		&models.DistributorAPILog{},
		&models.Transaction{},
		&models.Withdrawal{},
		&models.CategoryBinding{},
		&models.SecurityAuditLog{},
	)
}

func initRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.AppConfig.Redis.Host, config.AppConfig.Redis.Port),
		Password: config.AppConfig.Redis.Password,
		DB:       config.AppConfig.Redis.DB,
	})
}

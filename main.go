package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/yourusername/cloud-eye/internal/api/handler"
	"github.com/yourusername/cloud-eye/internal/api/router"
	"github.com/yourusername/cloud-eye/internal/pkg/config"
	"github.com/yourusername/cloud-eye/internal/pkg/database"
	"github.com/yourusername/cloud-eye/internal/pkg/logger"
	"github.com/yourusername/cloud-eye/internal/repository"
	"github.com/yourusername/cloud-eye/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	cfg, err := config.LoadConfig("configs/config.yaml")
	if err != nil {
		fmt.Println("Failed to load config:", err)
		os.Exit(1)
	}

	// 初始化日志
	err = logger.InitLogger()
	if err != nil {
		fmt.Println("Failed to initialize logger:", err)
		os.Exit(1)
	}

	logger.Info("Starting CloudEye Server...")

	// 设置Gin模式
	gin.SetMode(cfg.Server.Mode)

	// 初始化数据库
	err = database.InitDB()
	if err != nil {
		logger.Fatal("Failed to initialize database", err)
	}

	// 创建仓库层
	providerRepo := repository.NewCloudProviderRepository(database.DBClient)
	productRepo := repository.NewCloudProductRepository(database.DBClient)
	configItemRepo := repository.NewConfigurationItemRepository(database.DBClient)

	// 创建服务层
	providerService := service.NewCloudProviderService(providerRepo)
	productService := service.NewCloudProductService(productRepo, providerRepo)
	configItemService := service.NewConfigurationItemService(configItemRepo, providerRepo, productRepo)

	// 创建处理器层
	providerHandler := handler.NewCloudProviderHandler(providerService)
	productHandler := handler.NewCloudProductHandler(productService)
	configItemHandler := handler.NewConfigurationItemHandler(configItemService)

	// 初始化路由
	r := router.InitRouter(providerHandler, productHandler, configItemHandler)

	// 创建HTTP服务器
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: r,
	}

	// 优雅关闭服务器
	serverShutdown := make(chan struct{})
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		logger.Info("Shutting down server...")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			logger.Fatal("Server forced to shutdown:", err)
		}

		close(serverShutdown)
	}()

	// 启动服务器
	logger.Info(fmt.Sprintf("Server is running on port %d in %s mode", cfg.Server.Port, cfg.Server.Mode))
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatal("Failed to start server", err)
	}

	<-serverShutdown
	logger.Info("Server stopped")
}
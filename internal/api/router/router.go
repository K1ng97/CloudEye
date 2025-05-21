package router

import (
	"github.com/gin-gonic/gin"
	"github.com/yourusername/cloud-eye/internal/api/handler"
	"github.com/yourusername/cloud-eye/internal/pkg/logger"
)

// InitRouter 初始化路由
func InitRouter(
	cloudProviderHandler *handler.CloudProviderHandler,
	cloudProductHandler *handler.CloudProductHandler,
	configItemHandler *handler.ConfigurationItemHandler,
) *gin.Engine {
	r := gin.New()

	// 使用自定义中间件
	r.Use(gin.Recovery())
	r.Use(LoggerMiddleware())
	r.Use(CORSMiddleware())

	// API路由组
	api := r.Group("/api/v1")
	{
		// 云服务商相关路由
		providers := api.Group("/cloud-providers")
		{
			providers.GET("", cloudProviderHandler.GetAll)
			providers.GET("/:id", cloudProviderHandler.GetByID)
			providers.POST("", cloudProviderHandler.Create)
			providers.PUT("/:id", cloudProviderHandler.Update)
			providers.DELETE("/:id", cloudProviderHandler.Delete)

			// 获取指定云服务商的产品列表
			providers.GET("/:provider_id/products", cloudProductHandler.GetByProviderID)
			
			// 根据云服务商代码获取产品列表
			providers.GET("/code/:provider_code/products", cloudProductHandler.GetByProviderCode)
			
			// 获取指定云服务商和产品的配置项列表
			providers.GET("/:provider_id/products/:product_id/config-items", configItemHandler.GetByProviderAndProduct)
		}

		// 云产品相关路由
		products := api.Group("/cloud-products")
		{
			products.GET("", cloudProductHandler.GetAll)
			products.GET("/:id", cloudProductHandler.GetByID)
			products.POST("", cloudProductHandler.Create)
			products.PUT("/:id", cloudProductHandler.Update)
			products.DELETE("/:id", cloudProductHandler.Delete)
		}

		// 配置项相关路由
		configItems := api.Group("/config-items")
		{
			configItems.GET("", configItemHandler.GetByFilter)
			configItems.GET("/:id", configItemHandler.GetByID)
			configItems.POST("", configItemHandler.Create)
			configItems.PUT("/:id", configItemHandler.Update)
			configItems.DELETE("/:id", configItemHandler.Delete)
			
			// Excel导入导出
			configItems.GET("/export", configItemHandler.ExportExcel)
			configItems.POST("/import", configItemHandler.ImportExcel)
		}
	}

	// 添加健康检查接口
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	return r
}

// LoggerMiddleware 日志中间件
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 请求开始前
		path := c.Request.URL.Path
		method := c.Request.Method
		
		// 处理请求
		c.Next()
		
		// 请求结束后
		statusCode := c.Writer.Status()
		logger.Info("API Request",
			map[string]interface{}{
				"path":       path,
				"method":     method,
				"status":     statusCode,
				"client_ip":  c.ClientIP(),
				"user_agent": c.Request.UserAgent(),
			},
		)
	}
}

// CORSMiddleware 跨域中间件
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	}
}
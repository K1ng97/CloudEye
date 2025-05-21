package database

import (
	"fmt"
	"time"

	"github.com/yourusername/cloud-eye/internal/pkg/config"
	"github.com/yourusername/cloud-eye/internal/pkg/logger"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// DBClient 全局数据库客户端
var DBClient *gorm.DB

// InitDB 初始化数据库连接
func InitDB() error {
	cfg := config.GetConfig().Database

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
		cfg.Charset,
	)

	var err error
	DBClient, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名
		},
		Logger: config.NewGormLogger(), // 使用自定义日志器
	})
	if err != nil {
		logger.Error("Failed to connect to database", err)
		return err
	}

	sqlDB, err := DBClient.DB()
	if err != nil {
		logger.Error("Failed to get SQL DB", err)
		return err
	}

	// 设置连接池配置
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Hour)

	logger.Info("Database connection established successfully")
	return nil
}
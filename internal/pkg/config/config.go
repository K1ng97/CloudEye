package config

import (
	"log"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config 应用配置结构
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Log      LogConfig
	Excel    ExcelConfig
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port int
	Mode string
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Driver      string
	Host        string
	Port        int
	Username    string
	Password    string
	DBName      string
	Charset     string
	MaxIdleConns int
	MaxOpenConns int
	LogLevel     string
}

// LogConfig 日志配置
type LogConfig struct {
	Level    string
	Format   string
	Output   string
	Filename string
}

// ExcelConfig Excel导入导出配置
type ExcelConfig struct {
	ImportPath string
	ExportPath string
}

var config *Config

// LoadConfig 加载配置文件
func LoadConfig(configPath string) (*Config, error) {
	v := viper.New()
	v.SetConfigFile(configPath)
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	err := v.ReadInConfig()
	if err != nil {
		log.Printf("Error reading config file: %s", err)
		return nil, err
	}

	config = &Config{}
	err = v.Unmarshal(config)
	if err != nil {
		log.Printf("Error unmarshaling config: %s", err)
		return nil, err
	}

	return config, nil
}

// GetConfig 获取配置
func GetConfig() *Config {
	if config == nil {
		log.Fatal("Config not initialized")
	}
	return config
}

// 自定义GORM日志器
func NewGormLogger() GormLogger {
	return GormLogger{
		LogLevel: GetConfig().Database.LogLevel,
	}
}

// GormLogger 是一个GORM日志适配器
type GormLogger struct {
	LogLevel string
}

// 实现GORM日志接口
// ...实际实现略（将在logger包完成后集成）
package logger

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/yourusername/cloud-eye/internal/pkg/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

// InitLogger 初始化日志器
func InitLogger() error {
	cfg := config.GetConfig().Log

	// 解析日志级别
	var level zapcore.Level
	switch strings.ToLower(cfg.Level) {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	case "dpanic":
		level = zapcore.DPanicLevel
	case "panic":
		level = zapcore.PanicLevel
	case "fatal":
		level = zapcore.FatalLevel
	default:
		level = zapcore.InfoLevel
	}

	// 配置编码器
	var encoderConfig zapcore.EncoderConfig
	encoderConfig = zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 配置输出
	var writer zapcore.WriteSyncer
	if cfg.Output == "stdout" {
		writer = zapcore.AddSync(os.Stdout)
	} else {
		// 确保日志目录存在
		logDir := filepath.Dir(cfg.Filename)
		if err := os.MkdirAll(logDir, 0755); err != nil {
			return err
		}
		file, err := os.OpenFile(cfg.Filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		writer = zapcore.AddSync(file)
	}

	// 配置核心
	var encoder zapcore.Encoder
	if cfg.Format == "json" {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	core := zapcore.NewCore(encoder, writer, level)

	// 创建日志器
	logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	zap.ReplaceGlobals(logger)

	Info("Logger initialized successfully")
	return nil
}

// Debug 记录调试级别日志
func Debug(msg string, fields ...zap.Field) {
	if logger == nil {
		return
	}
	logger.Debug(msg, fields...)
}

// Info 记录信息级别日志
func Info(msg string, fields ...zap.Field) {
	if logger == nil {
		return
	}
	logger.Info(msg, fields...)
}

// Warn 记录警告级别日志
func Warn(msg string, fields ...zap.Field) {
	if logger == nil {
		return
	}
	logger.Warn(msg, fields...)
}

// Error 记录错误级别日志
func Error(msg string, err error, fields ...zap.Field) {
	if logger == nil {
		return
	}
	if err != nil {
		fields = append(fields, zap.Error(err))
	}
	logger.Error(msg, fields...)
}

// Fatal 记录致命级别日志并终止程序
func Fatal(msg string, err error, fields ...zap.Field) {
	if logger == nil {
		os.Exit(1)
	}
	if err != nil {
		fields = append(fields, zap.Error(err))
	}
	logger.Fatal(msg, fields...)
}
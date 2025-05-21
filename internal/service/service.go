package service

import (
	"context"

	"github.com/yourusername/cloud-eye/internal/repository"
)

// Service 定义了所有服务的通用接口
type Service interface {
	// 可以添加一些通用的服务方法
}

// BaseService 提供基础服务功能
type BaseService struct {
	// 可以添加一些共享的依赖
}

// 服务错误码定义
const (
	ErrCodeNotFound    = 1001 // 资源未找到
	ErrCodeInvalidData = 1002 // 无效的数据
	ErrCodeDatabase    = 1003 // 数据库错误
	ErrCodeDuplicate   = 1004 // 重复数据
	ErrCodeInternal    = 1005 // 内部错误
)

// ServiceError 服务错误类型
type ServiceError struct {
	Code    int    // 错误码
	Message string // 错误信息
	Err     error  // 原始错误
}

// Error 实现error接口
func (e *ServiceError) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

// NewServiceError 创建服务错误
func NewServiceError(code int, message string, err error) *ServiceError {
	return &ServiceError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// WithContext 给服务方法添加上下文
func WithContext(ctx context.Context) context.Context {
	if ctx == nil {
		return context.Background()
	}
	return ctx
}
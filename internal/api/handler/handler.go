package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/cloud-eye/internal/pkg/logger"
	"github.com/yourusername/cloud-eye/internal/service"
)

// Response 通用响应结构
type Response struct {
	Code    int         `json:"code"`              // 状态码，0表示成功
	Message string      `json:"message"`           // 消息
	Data    interface{} `json:"data,omitempty"`    // 响应数据
}

// BaseHandler 基础处理器
type BaseHandler struct {
	// 可以添加一些共享的依赖
}

// Success 成功响应
func (h *BaseHandler) Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

// Error 错误响应
func (h *BaseHandler) Error(c *gin.Context, httpStatus int, code int, message string) {
	c.JSON(httpStatus, Response{
		Code:    code,
		Message: message,
	})
}

// HandleServiceError 处理服务错误
func (h *BaseHandler) HandleServiceError(c *gin.Context, err error) {
	if err == nil {
		return
	}

	serviceErr, ok := err.(*service.ServiceError)
	if !ok {
		h.Error(c, http.StatusInternalServerError, 5000, "服务器内部错误")
		return
	}

	switch serviceErr.Code {
	case service.ErrCodeNotFound:
		h.Error(c, http.StatusNotFound, 4004, serviceErr.Message)
	case service.ErrCodeInvalidData:
		h.Error(c, http.StatusBadRequest, 4000, serviceErr.Message)
	case service.ErrCodeDuplicate:
		h.Error(c, http.StatusConflict, 4009, serviceErr.Message)
	default:
		h.Error(c, http.StatusInternalServerError, 5000, serviceErr.Message)
	}
}

// GetIDFromPath 从路径参数中获取ID
func (h *BaseHandler) GetIDFromPath(c *gin.Context, paramName string) (uint, bool) {
	idStr := c.Param(paramName)
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.Error(c, http.StatusBadRequest, 4000, "无效的ID参数")
		logger.Error("Invalid ID parameter", err)
		return 0, false
	}
	return uint(id), true
}

// GetQueryParam 获取查询参数
func (h *BaseHandler) GetQueryParam(c *gin.Context, paramName string) (string, bool) {
	value := c.Query(paramName)
	return value, value != ""
}

// GetIntQueryParam 获取整型查询参数
func (h *BaseHandler) GetIntQueryParam(c *gin.Context, paramName string, defaultValue int) int {
	valueStr := c.Query(paramName)
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

// GetUintQueryParam 获取无符号整型查询参数
func (h *BaseHandler) GetUintQueryParam(c *gin.Context, paramName string) (uint, bool) {
	valueStr := c.Query(paramName)
	if valueStr == "" {
		return 0, false
	}

	value, err := strconv.ParseUint(valueStr, 10, 32)
	if err != nil {
		return 0, false
	}
	return uint(value), true
}

// BindJSON 绑定JSON请求体
func (h *BaseHandler) BindJSON(c *gin.Context, obj interface{}) bool {
	if err := c.ShouldBindJSON(obj); err != nil {
		h.Error(c, http.StatusBadRequest, 4000, "无效的请求参数: "+err.Error())
		logger.Error("Invalid request body", err)
		return false
	}
	return true
}
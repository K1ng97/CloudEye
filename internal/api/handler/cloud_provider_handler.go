package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/cloud-eye/internal/models"
	"github.com/yourusername/cloud-eye/internal/pkg/logger"
	"github.com/yourusername/cloud-eye/internal/service"
	"go.uber.org/zap"
)

// CloudProviderHandler 云服务商API处理器
type CloudProviderHandler struct {
	BaseHandler
	service service.CloudProviderService
}

// NewCloudProviderHandler 创建云服务商处理器
func NewCloudProviderHandler(service service.CloudProviderService) *CloudProviderHandler {
	return &CloudProviderHandler{
		service: service,
	}
}

// GetAll 获取所有云服务商
// @Summary 获取所有云服务商
// @Description 获取系统中所有可用的云服务商列表
// @Tags 云服务商
// @Produce json
// @Success 200 {object} Response{data=[]models.CloudProvider} "成功"
// @Failure 500 {object} Response "服务器内部错误"
// @Router /api/v1/cloud-providers [get]
func (h *CloudProviderHandler) GetAll(c *gin.Context) {
	providers, err := h.service.GetAllProviders(c)
	if err != nil {
		logger.Error("Failed to get all cloud providers", err)
		h.HandleServiceError(c, err)
		return
	}

	h.Success(c, providers)
}

// GetByID 根据ID获取云服务商
// @Summary 获取云服务商详情
// @Description 根据ID获取云服务商详细信息
// @Tags 云服务商
// @Produce json
// @Param id path int true "云服务商ID"
// @Success 200 {object} Response{data=models.CloudProvider} "成功"
// @Failure 400 {object} Response "无效的ID参数"
// @Failure 404 {object} Response "云服务商不存在"
// @Failure 500 {object} Response "服务器内部错误"
// @Router /api/v1/cloud-providers/{id} [get]
func (h *CloudProviderHandler) GetByID(c *gin.Context) {
	id, ok := h.GetIDFromPath(c, "id")
	if !ok {
		return
	}

	provider, err := h.service.GetProviderByID(c, id)
	if err != nil {
		logger.Error("Failed to get cloud provider by ID", err, zap.Uint("id", id))
		h.HandleServiceError(c, err)
		return
	}

	h.Success(c, provider)
}

// Create 创建云服务商
// @Summary 创建云服务商
// @Description 创建新的云服务商
// @Tags 云服务商
// @Accept json
// @Produce json
// @Param provider body models.CloudProvider true "云服务商信息"
// @Success 200 {object} Response{data=models.CloudProvider} "成功"
// @Failure 400 {object} Response "无效的请求参数"
// @Failure 409 {object} Response "云服务商代码已存在"
// @Failure 500 {object} Response "服务器内部错误"
// @Router /api/v1/cloud-providers [post]
func (h *CloudProviderHandler) Create(c *gin.Context) {
	var provider models.CloudProvider
	if !h.BindJSON(c, &provider) {
		return
	}

	err := h.service.CreateProvider(c, &provider)
	if err != nil {
		logger.Error("Failed to create cloud provider", err)
		h.HandleServiceError(c, err)
		return
	}

	h.Success(c, provider)
}

// Update 更新云服务商
// @Summary 更新云服务商
// @Description 更新已有的云服务商信息
// @Tags 云服务商
// @Accept json
// @Produce json
// @Param id path int true "云服务商ID"
// @Param provider body models.CloudProvider true "云服务商信息"
// @Success 200 {object} Response "成功"
// @Failure 400 {object} Response "无效的请求参数"
// @Failure 404 {object} Response "云服务商不存在"
// @Failure 409 {object} Response "云服务商代码已存在"
// @Failure 500 {object} Response "服务器内部错误"
// @Router /api/v1/cloud-providers/{id} [put]
func (h *CloudProviderHandler) Update(c *gin.Context) {
	id, ok := h.GetIDFromPath(c, "id")
	if !ok {
		return
	}

	var provider models.CloudProvider
	if !h.BindJSON(c, &provider) {
		return
	}

	// 确保路径参数ID与请求体ID一致
	provider.ID = id

	err := h.service.UpdateProvider(c, &provider)
	if err != nil {
		logger.Error("Failed to update cloud provider", err, zap.Uint("id", id))
		h.HandleServiceError(c, err)
		return
	}

	h.Success(c, gin.H{"message": "云服务商更新成功"})
}

// Delete 删除云服务商
// @Summary 删除云服务商
// @Description 删除指定的云服务商
// @Tags 云服务商
// @Produce json
// @Param id path int true "云服务商ID"
// @Success 200 {object} Response "成功"
// @Failure 400 {object} Response "无效的ID参数"
// @Failure 404 {object} Response "云服务商不存在"
// @Failure 500 {object} Response "服务器内部错误"
// @Router /api/v1/cloud-providers/{id} [delete]
func (h *CloudProviderHandler) Delete(c *gin.Context) {
	id, ok := h.GetIDFromPath(c, "id")
	if !ok {
		return
	}

	err := h.service.DeleteProvider(c, id)
	if err != nil {
		logger.Error("Failed to delete cloud provider", err, zap.Uint("id", id))
		h.HandleServiceError(c, err)
		return
	}

	h.Success(c, gin.H{"message": "云服务商删除成功"})
}
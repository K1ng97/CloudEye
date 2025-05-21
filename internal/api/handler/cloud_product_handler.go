package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/cloud-eye/internal/models"
	"github.com/yourusername/cloud-eye/internal/pkg/logger"
	"github.com/yourusername/cloud-eye/internal/service"
	"go.uber.org/zap"
)

// CloudProductHandler 云产品API处理器
type CloudProductHandler struct {
	BaseHandler
	service service.CloudProductService
}

// NewCloudProductHandler 创建云产品处理器
func NewCloudProductHandler(service service.CloudProductService) *CloudProductHandler {
	return &CloudProductHandler{
		service: service,
	}
}

// GetAll 获取所有云产品
// @Summary 获取所有云产品
// @Description 获取系统中所有可用的云产品列表
// @Tags 云产品
// @Produce json
// @Success 200 {object} Response{data=[]models.CloudProduct} "成功"
// @Failure 500 {object} Response "服务器内部错误"
// @Router /api/v1/cloud-products [get]
func (h *CloudProductHandler) GetAll(c *gin.Context) {
	products, err := h.service.GetAllProducts(c)
	if err != nil {
		logger.Error("Failed to get all cloud products", err)
		h.HandleServiceError(c, err)
		return
	}

	h.Success(c, products)
}

// GetByID 根据ID获取云产品
// @Summary 获取云产品详情
// @Description 根据ID获取云产品详细信息
// @Tags 云产品
// @Produce json
// @Param id path int true "云产品ID"
// @Success 200 {object} Response{data=models.CloudProduct} "成功"
// @Failure 400 {object} Response "无效的ID参数"
// @Failure 404 {object} Response "云产品不存在"
// @Failure 500 {object} Response "服务器内部错误"
// @Router /api/v1/cloud-products/{id} [get]
func (h *CloudProductHandler) GetByID(c *gin.Context) {
	id, ok := h.GetIDFromPath(c, "id")
	if !ok {
		return
	}

	product, err := h.service.GetProductByID(c, id)
	if err != nil {
		logger.Error("Failed to get cloud product by ID", err, zap.Uint("id", id))
		h.HandleServiceError(c, err)
		return
	}

	h.Success(c, product)
}

// GetByProviderID 获取指定云服务商的产品列表
// @Summary 获取指定云服务商的产品列表
// @Description 根据云服务商ID获取其所有产品
// @Tags 云产品
// @Produce json
// @Param provider_id path int true "云服务商ID"
// @Success 200 {object} Response{data=[]models.CloudProduct} "成功"
// @Failure 400 {object} Response "无效的ID参数"
// @Failure 404 {object} Response "云服务商不存在"
// @Failure 500 {object} Response "服务器内部错误"
// @Router /api/v1/cloud-providers/{provider_id}/products [get]
func (h *CloudProductHandler) GetByProviderID(c *gin.Context) {
	providerID, ok := h.GetIDFromPath(c, "provider_id")
	if !ok {
		return
	}

	products, err := h.service.GetProductsByProviderID(c, providerID)
	if err != nil {
		logger.Error("Failed to get products by provider ID", err, zap.Uint("providerId", providerID))
		h.HandleServiceError(c, err)
		return
	}

	h.Success(c, products)
}

// GetByProviderCode 获取指定云服务商代码的产品列表
// @Summary 获取指定云服务商代码的产品列表
// @Description 根据云服务商代码获取其所有产品
// @Tags 云产品
// @Produce json
// @Param provider_code path string true "云服务商代码"
// @Success 200 {object} Response{data=[]models.CloudProduct} "成功"
// @Failure 404 {object} Response "云服务商不存在"
// @Failure 500 {object} Response "服务器内部错误"
// @Router /api/v1/cloud-providers/code/{provider_code}/products [get]
func (h *CloudProductHandler) GetByProviderCode(c *gin.Context) {
	providerCode := c.Param("provider_code")
	if providerCode == "" {
		h.Error(c, http.StatusBadRequest, 4000, "云服务商代码不能为空")
		return
	}

	products, err := h.service.GetProductsByProviderCode(c, providerCode)
	if err != nil {
		logger.Error("Failed to get products by provider code", err, zap.String("providerCode", providerCode))
		h.HandleServiceError(c, err)
		return
	}

	h.Success(c, products)
}

// Create 创建云产品
// @Summary 创建云产品
// @Description 创建新的云产品
// @Tags 云产品
// @Accept json
// @Produce json
// @Param product body models.CloudProduct true "云产品信息"
// @Success 200 {object} Response{data=models.CloudProduct} "成功"
// @Failure 400 {object} Response "无效的请求参数"
// @Failure 404 {object} Response "云服务商不存在"
// @Failure 409 {object} Response "云产品代码已存在"
// @Failure 500 {object} Response "服务器内部错误"
// @Router /api/v1/cloud-products [post]
func (h *CloudProductHandler) Create(c *gin.Context) {
	var product models.CloudProduct
	if !h.BindJSON(c, &product) {
		return
	}

	err := h.service.CreateProduct(c, &product)
	if err != nil {
		logger.Error("Failed to create cloud product", err)
		h.HandleServiceError(c, err)
		return
	}

	h.Success(c, product)
}

// Update 更新云产品
// @Summary 更新云产品
// @Description 更新已有的云产品信息
// @Tags 云产品
// @Accept json
// @Produce json
// @Param id path int true "云产品ID"
// @Param product body models.CloudProduct true "云产品信息"
// @Success 200 {object} Response "成功"
// @Failure 400 {object} Response "无效的请求参数"
// @Failure 404 {object} Response "云产品不存在或云服务商不存在"
// @Failure 409 {object} Response "云产品代码已存在"
// @Failure 500 {object} Response "服务器内部错误"
// @Router /api/v1/cloud-products/{id} [put]
func (h *CloudProductHandler) Update(c *gin.Context) {
	id, ok := h.GetIDFromPath(c, "id")
	if !ok {
		return
	}

	var product models.CloudProduct
	if !h.BindJSON(c, &product) {
		return
	}

	// 确保路径参数ID与请求体ID一致
	product.ID = id

	err := h.service.UpdateProduct(c, &product)
	if err != nil {
		logger.Error("Failed to update cloud product", err, zap.Uint("id", id))
		h.HandleServiceError(c, err)
		return
	}

	h.Success(c, gin.H{"message": "云产品更新成功"})
}

// Delete 删除云产品
// @Summary 删除云产品
// @Description 删除指定的云产品
// @Tags 云产品
// @Produce json
// @Param id path int true "云产品ID"
// @Success 200 {object} Response "成功"
// @Failure 400 {object} Response "无效的ID参数"
// @Failure 404 {object} Response "云产品不存在"
// @Failure 500 {object} Response "服务器内部错误"
// @Router /api/v1/cloud-products/{id} [delete]
func (h *CloudProductHandler) Delete(c *gin.Context) {
	id, ok := h.GetIDFromPath(c, "id")
	if !ok {
		return
	}

	err := h.service.DeleteProduct(c, id)
	if err != nil {
		logger.Error("Failed to delete cloud product", err, zap.Uint("id", id))
		h.HandleServiceError(c, err)
		return
	}

	h.Success(c, gin.H{"message": "云产品删除成功"})
}
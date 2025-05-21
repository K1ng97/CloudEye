package handler

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/cloud-eye/internal/models"
	"github.com/yourusername/cloud-eye/internal/pkg/excel"
	"github.com/yourusername/cloud-eye/internal/pkg/logger"
	"github.com/yourusername/cloud-eye/internal/repository"
	"github.com/yourusername/cloud-eye/internal/service"
	"go.uber.org/zap"
)

// ConfigurationItemHandler 配置项API处理器
type ConfigurationItemHandler struct {
	BaseHandler
	service service.ConfigurationItemService
	exporter *excel.ConfigItemExporter
	importer *excel.ConfigItemImporter
}

// NewConfigurationItemHandler 创建配置项处理器
func NewConfigurationItemHandler(service service.ConfigurationItemService) *ConfigurationItemHandler {
	return &ConfigurationItemHandler{
		service:  service,
		exporter: excel.NewConfigItemExporter(),
		importer: excel.NewConfigItemImporter(),
	}
}

// GetByID 根据ID获取配置项
// @Summary 获取配置项详情
// @Description 根据ID获取配置项详细信息
// @Tags 配置项
// @Produce json
// @Param id path int true "配置项ID"
// @Success 200 {object} Response{data=models.ConfigurationItem} "成功"
// @Failure 400 {object} Response "无效的ID参数"
// @Failure 404 {object} Response "配置项不存在"
// @Failure 500 {object} Response "服务器内部错误"
// @Router /api/v1/config-items/{id} [get]
func (h *ConfigurationItemHandler) GetByID(c *gin.Context) {
	id, ok := h.GetIDFromPath(c, "id")
	if !ok {
		return
	}

	item, err := h.service.GetConfigItemByID(c, id)
	if err != nil {
		logger.Error("Failed to get config item by ID", err, zap.Uint("id", id))
		h.HandleServiceError(c, err)
		return
	}

	h.Success(c, item)
}

// GetByFilter 根据过滤条件获取配置项列表（支持分页）
// @Summary 获取配置项列表
// @Description 根据过滤条件获取配置项列表，支持分页
// @Tags 配置项
// @Produce json
// @Param cloud_provider_id query int false "云服务商ID"
// @Param product_id query int false "产品ID"
// @Param keyword query string false "关键词搜索"
// @Param page query int false "页码，默认1"
// @Param page_size query int false "每页记录数，默认10"
// @Success 200 {object} Response{data=repository.PageResult} "成功"
// @Failure 500 {object} Response "服务器内部错误"
// @Router /api/v1/config-items [get]
func (h *ConfigurationItemHandler) GetByFilter(c *gin.Context) {
	filter := repository.ConfigItemFilter{
		Page:     h.GetIntQueryParam(c, "page", 1),
		PageSize: h.GetIntQueryParam(c, "page_size", 10),
	}

	// 获取可选过滤参数
	if providerID, ok := h.GetUintQueryParam(c, "cloud_provider_id"); ok {
		filter.CloudProviderID = &providerID
	}

	if productID, ok := h.GetUintQueryParam(c, "product_id"); ok {
		filter.ProductID = &productID
	}

	if keyword, ok := h.GetQueryParam(c, "keyword"); ok {
		filter.Keyword = &keyword
	}

	result, err := h.service.GetConfigItemsByFilter(c, filter)
	if err != nil {
		logger.Error("Failed to get config items by filter", err)
		h.HandleServiceError(c, err)
		return
	}

	h.Success(c, result)
}

// GetByProviderAndProduct 获取指定云服务商和产品的配置项列表
// @Summary 获取指定云服务商和产品的配置项列表
// @Description 根据云服务商ID和产品ID获取配置项列表
// @Tags 配置项
// @Produce json
// @Param provider_id path int true "云服务商ID"
// @Param product_id path int true "产品ID"
// @Success 200 {object} Response{data=[]models.ConfigurationItem} "成功"
// @Failure 400 {object} Response "无效的ID参数"
// @Failure 404 {object} Response "云服务商或产品不存在"
// @Failure 500 {object} Response "服务器内部错误"
// @Router /api/v1/cloud-providers/{provider_id}/products/{product_id}/config-items [get]
func (h *ConfigurationItemHandler) GetByProviderAndProduct(c *gin.Context) {
	providerID, ok := h.GetIDFromPath(c, "provider_id")
	if !ok {
		return
	}

	productID, ok := h.GetIDFromPath(c, "product_id")
	if !ok {
		return
	}

	items, err := h.service.GetConfigItemsByProviderAndProduct(c, providerID, productID)
	if err != nil {
		logger.Error("Failed to get config items by provider and product", err,
			zap.Uint("providerId", providerID),
			zap.Uint("productId", productID))
		h.HandleServiceError(c, err)
		return
	}

	h.Success(c, items)
}

// Create 创建配置项
// @Summary 创建配置项
// @Description 创建新的配置项
// @Tags 配置项
// @Accept json
// @Produce json
// @Param item body models.ConfigurationItem true "配置项信息"
// @Success 200 {object} Response{data=models.ConfigurationItem} "成功"
// @Failure 400 {object} Response "无效的请求参数"
// @Failure 404 {object} Response "云服务商或产品不存在"
// @Failure 500 {object} Response "服务器内部错误"
// @Router /api/v1/config-items [post]
func (h *ConfigurationItemHandler) Create(c *gin.Context) {
	var item models.ConfigurationItem
	if !h.BindJSON(c, &item) {
		return
	}

	err := h.service.CreateConfigItem(c, &item)
	if err != nil {
		logger.Error("Failed to create config item", err)
		h.HandleServiceError(c, err)
		return
	}

	h.Success(c, item)
}

// Update 更新配置项
// @Summary 更新配置项
// @Description 更新已有的配置项信息
// @Tags 配置项
// @Accept json
// @Produce json
// @Param id path int true "配置项ID"
// @Param item body models.ConfigurationItem true "配置项信息"
// @Success 200 {object} Response "成功"
// @Failure 400 {object} Response "无效的请求参数"
// @Failure 404 {object} Response "配置项不存在或云服务商或产品不存在"
// @Failure 500 {object} Response "服务器内部错误"
// @Router /api/v1/config-items/{id} [put]
func (h *ConfigurationItemHandler) Update(c *gin.Context) {
	id, ok := h.GetIDFromPath(c, "id")
	if !ok {
		return
	}

	var item models.ConfigurationItem
	if !h.BindJSON(c, &item) {
		return
	}

	// 确保路径参数ID与请求体ID一致
	item.ID = id

	err := h.service.UpdateConfigItem(c, &item)
	if err != nil {
		logger.Error("Failed to update config item", err, zap.Uint("id", id))
		h.HandleServiceError(c, err)
		return
	}

	h.Success(c, gin.H{"message": "配置项更新成功"})
}

// Delete 删除配置项
// @Summary 删除配置项
// @Description 删除指定的配置项
// @Tags 配置项
// @Produce json
// @Param id path int true "配置项ID"
// @Success 200 {object} Response "成功"
// @Failure 400 {object} Response "无效的ID参数"
// @Failure 404 {object} Response "配置项不存在"
// @Failure 500 {object} Response "服务器内部错误"
// @Router /api/v1/config-items/{id} [delete]
func (h *ConfigurationItemHandler) Delete(c *gin.Context) {
	id, ok := h.GetIDFromPath(c, "id")
	if !ok {
		return
	}

	err := h.service.DeleteConfigItem(c, id)
	if err != nil {
		logger.Error("Failed to delete config item", err, zap.Uint("id", id))
		h.HandleServiceError(c, err)
		return
	}

	h.Success(c, gin.H{"message": "配置项删除成功"})
}

// ExportExcel 导出配置项到Excel
// @Summary 导出配置项到Excel
// @Description 根据过滤条件导出配置项到Excel文件
// @Tags 配置项
// @Produce json
// @Param cloud_provider_id query int false "云服务商ID"
// @Param product_id query int false "产品ID"
// @Param keyword query string false "关键词搜索"
// @Success 200 {object} Response "成功"
// @Failure 500 {object} Response "服务器内部错误"
// @Router /api/v1/config-items/export [get]
func (h *ConfigurationItemHandler) ExportExcel(c *gin.Context) {
	filter := repository.ConfigItemFilter{
		// 导出时不分页，获取所有符合条件的数据
		Page:     1,
		PageSize: 1000, // 较大的页大小，实际会限制在100以内
	}

	// 获取可选过滤参数
	if providerID, ok := h.GetUintQueryParam(c, "cloud_provider_id"); ok {
		filter.CloudProviderID = &providerID
	}

	if productID, ok := h.GetUintQueryParam(c, "product_id"); ok {
		filter.ProductID = &productID
	}

	if keyword, ok := h.GetQueryParam(c, "keyword"); ok {
		filter.Keyword = &keyword
	}

	// 获取数据
	result, err := h.service.GetConfigItemsByFilter(c, filter)
	if err != nil {
		logger.Error("Failed to get config items for export", err)
		h.HandleServiceError(c, err)
		return
	}

	// 将数据转换为模型列表
	items, ok := result.Data.([]models.ConfigurationItem)
	if !ok {
		logger.Error("Failed to convert data to configuration items", nil)
		h.Error(c, http.StatusInternalServerError, 5000, "导出失败：数据类型错误")
		return
	}

	// 导出到Excel
	filePath, err := h.exporter.Export(c, items)
	if err != nil {
		logger.Error("Failed to export to Excel", err)
		h.Error(c, http.StatusInternalServerError, 5000, "导出Excel失败："+err.Error())
		return
	}

	// 返回下载链接
	h.Success(c, gin.H{
		"message":  "导出成功",
		"filePath": filePath,
		"fileName": filepath.Base(filePath),
	})
}

// ImportExcel 从Excel导入配置项
// @Summary 从Excel导入配置项
// @Description 从上传的Excel文件导入配置项
// @Tags 配置项
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Excel文件"
// @Success 200 {object} Response "成功"
// @Failure 400 {object} Response "无效的文件"
// @Failure 500 {object} Response "服务器内部错误"
// @Router /api/v1/config-items/import [post]
func (h *ConfigurationItemHandler) ImportExcel(c *gin.Context) {
	// 从表单获取文件
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		h.Error(c, http.StatusBadRequest, 4000, "请选择要导入的Excel文件")
		return
	}
	defer file.Close()

	// 验证文件类型
	if !strings.HasSuffix(strings.ToLower(header.Filename), ".xlsx") {
		h.Error(c, http.StatusBadRequest, 4000, "只支持.xlsx格式的Excel文件")
		return
	}

	// 生成唯一文件名
	filename := fmt.Sprintf("%d_%s", time.Now().Unix(), header.Filename)

	// 保存文件
	filePath, err := h.importer.SaveUploadedFile(file, filename)
	if err != nil {
		logger.Error("Failed to save uploaded file", err)
		h.Error(c, http.StatusInternalServerError, 5000, "保存文件失败："+err.Error())
		return
	}

	// 解析Excel
	items, err := h.importer.ImportConfigItems(c, filePath)
	if err != nil {
		logger.Error("Failed to parse Excel file", err)
		h.Error(c, http.StatusBadRequest, 4000, "解析Excel文件失败："+err.Error())
		return
	}

	// 批量导入数据
	err = h.service.BatchImportConfigItems(c, items)
	if err != nil {
		logger.Error("Failed to import config items", err)
		h.HandleServiceError(c, err)
		return
	}

	h.Success(c, gin.H{
		"message": "导入成功",
		"count":   len(items),
	})
}
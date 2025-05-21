package repository

import (
	"context"
	"errors"

	"github.com/yourusername/cloud-eye/internal/models"
	"github.com/yourusername/cloud-eye/internal/pkg/logger"
	"gorm.io/gorm"
)

// ConfigItemFilter 配置项查询过滤条件
type ConfigItemFilter struct {
	CloudProviderID *uint   `json:"cloud_provider_id,omitempty"`
	ProductID       *uint   `json:"product_id,omitempty"`
	Keyword         *string `json:"keyword,omitempty"`
	Page            int     `json:"page"`
	PageSize        int     `json:"page_size"`
}

// ConfigurationItemRepository 配置项仓库接口
type ConfigurationItemRepository interface {
	Repository
	GetByID(ctx context.Context, id uint) (*models.ConfigurationItem, error)
	GetByFilter(ctx context.Context, filter ConfigItemFilter) (*PageResult, error)
	GetByProviderAndProduct(ctx context.Context, providerID, productID uint) ([]models.ConfigurationItem, error)
	Create(ctx context.Context, item *models.ConfigurationItem) error
	Update(ctx context.Context, item *models.ConfigurationItem) error
	Delete(ctx context.Context, id uint) error
	BatchInsert(ctx context.Context, items []models.ConfigurationItem) error
}

// configurationItemRepository 配置项仓库实现
type configurationItemRepository struct {
	BaseRepository
}

// NewConfigurationItemRepository 创建配置项仓库
func NewConfigurationItemRepository(db *gorm.DB) ConfigurationItemRepository {
	return &configurationItemRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

// GetByID 根据ID获取配置项
func (r *configurationItemRepository) GetByID(ctx context.Context, id uint) (*models.ConfigurationItem, error) {
	var item models.ConfigurationItem
	err := r.DB.WithContext(ctx).First(&item, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		logger.Error("Failed to get configuration item by ID", err)
		return nil, err
	}
	return &item, nil
}

// GetByFilter 根据过滤条件获取配置项，支持分页
func (r *configurationItemRepository) GetByFilter(ctx context.Context, filter ConfigItemFilter) (*PageResult, error) {
	var items []models.ConfigurationItem
	var total int64

	query := r.DB.WithContext(ctx).Model(&models.ConfigurationItem{})

	// 应用过滤条件
	if filter.CloudProviderID != nil {
		query = query.Where("cloud_provider_id = ?", *filter.CloudProviderID)
	}

	if filter.ProductID != nil {
		query = query.Where("product_id = ?", *filter.ProductID)
	}

	if filter.Keyword != nil && *filter.Keyword != "" {
		query = query.Where("name LIKE ? OR recommended_value LIKE ? OR risk_description LIKE ?",
			"%"+*filter.Keyword+"%", "%"+*filter.Keyword+"%", "%"+*filter.Keyword+"%")
	}

	// 计算总数
	err := query.Count(&total).Error
	if err != nil {
		logger.Error("Failed to count configuration items", err)
		return nil, err
	}

	// 应用分页并查询数据
	err = query.Scopes(Paginate(filter.Page, filter.PageSize)).
		Preload("Provider").
		Preload("Product").
		Find(&items).Error
	if err != nil {
		logger.Error("Failed to get configuration items by filter", err)
		return nil, err
	}

	return &PageResult{
		Total:    total,
		Page:     filter.Page,
		PageSize: filter.PageSize,
		Data:     items,
	}, nil
}

// GetByProviderAndProduct 根据云服务商ID和产品ID获取配置项
func (r *configurationItemRepository) GetByProviderAndProduct(ctx context.Context, providerID, productID uint) ([]models.ConfigurationItem, error) {
	var items []models.ConfigurationItem
	err := r.DB.WithContext(ctx).
		Where("cloud_provider_id = ? AND product_id = ?", providerID, productID).
		Find(&items).Error
	if err != nil {
		logger.Error("Failed to get configuration items by provider and product", err)
		return nil, err
	}
	return items, nil
}

// Create 创建配置项
func (r *configurationItemRepository) Create(ctx context.Context, item *models.ConfigurationItem) error {
	err := r.DB.WithContext(ctx).Create(item).Error
	if err != nil {
		logger.Error("Failed to create configuration item", err)
		return err
	}
	return nil
}

// Update 更新配置项
func (r *configurationItemRepository) Update(ctx context.Context, item *models.ConfigurationItem) error {
	err := r.DB.WithContext(ctx).Save(item).Error
	if err != nil {
		logger.Error("Failed to update configuration item", err)
		return err
	}
	return nil
}

// Delete 删除配置项
func (r *configurationItemRepository) Delete(ctx context.Context, id uint) error {
	err := r.DB.WithContext(ctx).Delete(&models.ConfigurationItem{}, id).Error
	if err != nil {
		logger.Error("Failed to delete configuration item", err)
		return err
	}
	return nil
}

// BatchInsert 批量插入配置项（用于Excel导入）
func (r *configurationItemRepository) BatchInsert(ctx context.Context, items []models.ConfigurationItem) error {
	return r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, item := range items {
			if err := tx.Create(&item).Error; err != nil {
				logger.Error("Failed to batch insert configuration item", err)
				return err
			}
		}
		return nil
	})
}
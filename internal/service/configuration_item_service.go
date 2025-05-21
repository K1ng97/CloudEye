package service

import (
	"context"

	"github.com/yourusername/cloud-eye/internal/models"
	"github.com/yourusername/cloud-eye/internal/pkg/logger"
	"github.com/yourusername/cloud-eye/internal/repository"
	"go.uber.org/zap"
)

// ConfigurationItemService 配置项服务接口
type ConfigurationItemService interface {
	Service
	GetConfigItemByID(ctx context.Context, id uint) (*models.ConfigurationItem, error)
	GetConfigItemsByFilter(ctx context.Context, filter repository.ConfigItemFilter) (*repository.PageResult, error)
	GetConfigItemsByProviderAndProduct(ctx context.Context, providerID, productID uint) ([]models.ConfigurationItem, error)
	CreateConfigItem(ctx context.Context, item *models.ConfigurationItem) error
	UpdateConfigItem(ctx context.Context, item *models.ConfigurationItem) error
	DeleteConfigItem(ctx context.Context, id uint) error
	BatchImportConfigItems(ctx context.Context, items []models.ConfigurationItem) error
}

// configurationItemService 配置项服务实现
type configurationItemService struct {
	BaseService
	repo         repository.ConfigurationItemRepository
	providerRepo repository.CloudProviderRepository
	productRepo  repository.CloudProductRepository
}

// NewConfigurationItemService 创建配置项服务
func NewConfigurationItemService(
	repo repository.ConfigurationItemRepository,
	providerRepo repository.CloudProviderRepository,
	productRepo repository.CloudProductRepository,
) ConfigurationItemService {
	return &configurationItemService{
		repo:         repo,
		providerRepo: providerRepo,
		productRepo:  productRepo,
	}
}

// GetConfigItemByID 根据ID获取配置项
func (s *configurationItemService) GetConfigItemByID(ctx context.Context, id uint) (*models.ConfigurationItem, error) {
	ctx = WithContext(ctx)
	logger.Info("Getting configuration item by ID", zap.Uint("id", id))

	item, err := s.repo.GetByID(ctx, id)
	if err != nil {
		logger.Error("Failed to get configuration item by ID", err, zap.Uint("id", id))
		return nil, NewServiceError(ErrCodeDatabase, "获取配置项详情失败", err)
	}

	if item == nil {
		return nil, NewServiceError(ErrCodeNotFound, "配置项不存在", nil)
	}

	return item, nil
}

// GetConfigItemsByFilter 根据过滤条件获取配置项列表，支持分页
func (s *configurationItemService) GetConfigItemsByFilter(ctx context.Context, filter repository.ConfigItemFilter) (*repository.PageResult, error) {
	ctx = WithContext(ctx)
	logger.Info("Getting configuration items by filter",
		zap.Any("filter", filter))

	// 参数检查和验证
	if filter.CloudProviderID != nil {
		// 检查服务商是否存在
		provider, err := s.providerRepo.GetByID(ctx, *filter.CloudProviderID)
		if err != nil {
			logger.Error("Failed to check provider existence", err, zap.Uint("providerId", *filter.CloudProviderID))
			return nil, NewServiceError(ErrCodeDatabase, "获取配置项列表失败", err)
		}

		if provider == nil {
			return nil, NewServiceError(ErrCodeNotFound, "云服务商不存在", nil)
		}
	}

	if filter.ProductID != nil {
		// 检查产品是否存在
		product, err := s.productRepo.GetByID(ctx, *filter.ProductID)
		if err != nil {
			logger.Error("Failed to check product existence", err, zap.Uint("productId", *filter.ProductID))
			return nil, NewServiceError(ErrCodeDatabase, "获取配置项列表失败", err)
		}

		if product == nil {
			return nil, NewServiceError(ErrCodeNotFound, "云产品不存在", nil)
		}
	}

	// 设置默认分页参数
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 {
		filter.PageSize = 10
	} else if filter.PageSize > 100 {
		filter.PageSize = 100
	}

	result, err := s.repo.GetByFilter(ctx, filter)
	if err != nil {
		logger.Error("Failed to get configuration items by filter", err)
		return nil, NewServiceError(ErrCodeDatabase, "获取配置项列表失败", err)
	}

	return result, nil
}

// GetConfigItemsByProviderAndProduct 根据云服务商ID和产品ID获取配置项列表
func (s *configurationItemService) GetConfigItemsByProviderAndProduct(ctx context.Context, providerID, productID uint) ([]models.ConfigurationItem, error) {
	ctx = WithContext(ctx)
	logger.Info("Getting configuration items by provider and product",
		zap.Uint("providerId", providerID),
		zap.Uint("productId", productID))

	// 检查服务商是否存在
	provider, err := s.providerRepo.GetByID(ctx, providerID)
	if err != nil {
		logger.Error("Failed to check provider existence", err, zap.Uint("providerId", providerID))
		return nil, NewServiceError(ErrCodeDatabase, "获取配置项列表失败", err)
	}

	if provider == nil {
		return nil, NewServiceError(ErrCodeNotFound, "云服务商不存在", nil)
	}

	// 检查产品是否存在
	product, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		logger.Error("Failed to check product existence", err, zap.Uint("productId", productID))
		return nil, NewServiceError(ErrCodeDatabase, "获取配置项列表失败", err)
	}

	if product == nil {
		return nil, NewServiceError(ErrCodeNotFound, "云产品不存在", nil)
	}

	items, err := s.repo.GetByProviderAndProduct(ctx, providerID, productID)
	if err != nil {
		logger.Error("Failed to get configuration items by provider and product", err,
			zap.Uint("providerId", providerID),
			zap.Uint("productId", productID))
		return nil, NewServiceError(ErrCodeDatabase, "获取配置项列表失败", err)
	}

	return items, nil
}

// CreateConfigItem 创建配置项
func (s *configurationItemService) CreateConfigItem(ctx context.Context, item *models.ConfigurationItem) error {
	ctx = WithContext(ctx)
	logger.Info("Creating configuration item", zap.String("name", item.Name))

	// 检查服务商是否存在
	provider, err := s.providerRepo.GetByID(ctx, item.CloudProviderID)
	if err != nil {
		logger.Error("Failed to check provider existence", err, zap.Uint("providerId", item.CloudProviderID))
		return NewServiceError(ErrCodeDatabase, "创建配置项失败", err)
	}

	if provider == nil {
		return NewServiceError(ErrCodeNotFound, "云服务商不存在", nil)
	}

	// 检查产品是否存在
	product, err := s.productRepo.GetByID(ctx, item.ProductID)
	if err != nil {
		logger.Error("Failed to check product existence", err, zap.Uint("productId", item.ProductID))
		return NewServiceError(ErrCodeDatabase, "创建配置项失败", err)
	}

	if product == nil {
		return NewServiceError(ErrCodeNotFound, "云产品不存在", nil)
	}

	// 检查产品是否属于指定的服务商
	if product.CloudProviderID != item.CloudProviderID {
		return NewServiceError(ErrCodeInvalidData, "云产品不属于指定的云服务商", nil)
	}

	if err := s.repo.Create(ctx, item); err != nil {
		logger.Error("Failed to create configuration item", err)
		return NewServiceError(ErrCodeDatabase, "创建配置项失败", err)
	}

	return nil
}

// UpdateConfigItem 更新配置项
func (s *configurationItemService) UpdateConfigItem(ctx context.Context, item *models.ConfigurationItem) error {
	ctx = WithContext(ctx)
	logger.Info("Updating configuration item", zap.Uint("id", item.ID))

	// 检查配置项是否存在
	existingItem, err := s.repo.GetByID(ctx, item.ID)
	if err != nil {
		logger.Error("Failed to check configuration item existence", err, zap.Uint("id", item.ID))
		return NewServiceError(ErrCodeDatabase, "更新配置项失败", err)
	}

	if existingItem == nil {
		return NewServiceError(ErrCodeNotFound, "配置项不存在", nil)
	}

	// 如果云服务商ID有变更，检查新的服务商是否存在
	if item.CloudProviderID != existingItem.CloudProviderID {
		provider, err := s.providerRepo.GetByID(ctx, item.CloudProviderID)
		if err != nil {
			logger.Error("Failed to check provider existence", err, zap.Uint("providerId", item.CloudProviderID))
			return NewServiceError(ErrCodeDatabase, "更新配置项失败", err)
		}

		if provider == nil {
			return NewServiceError(ErrCodeNotFound, "云服务商不存在", nil)
		}
	}

	// 如果产品ID有变更，检查新的产品是否存在
	if item.ProductID != existingItem.ProductID {
		product, err := s.productRepo.GetByID(ctx, item.ProductID)
		if err != nil {
			logger.Error("Failed to check product existence", err, zap.Uint("productId", item.ProductID))
			return NewServiceError(ErrCodeDatabase, "更新配置项失败", err)
		}

		if product == nil {
			return NewServiceError(ErrCodeNotFound, "云产品不存在", nil)
		}

		// 检查产品是否属于指定的服务商
		if product.CloudProviderID != item.CloudProviderID {
			return NewServiceError(ErrCodeInvalidData, "云产品不属于指定的云服务商", nil)
		}
	}

	if err := s.repo.Update(ctx, item); err != nil {
		logger.Error("Failed to update configuration item", err)
		return NewServiceError(ErrCodeDatabase, "更新配置项失败", err)
	}

	return nil
}

// DeleteConfigItem 删除配置项
func (s *configurationItemService) DeleteConfigItem(ctx context.Context, id uint) error {
	ctx = WithContext(ctx)
	logger.Info("Deleting configuration item", zap.Uint("id", id))

	// 检查配置项是否存在
	existingItem, err := s.repo.GetByID(ctx, id)
	if err != nil {
		logger.Error("Failed to check configuration item existence", err, zap.Uint("id", id))
		return NewServiceError(ErrCodeDatabase, "删除配置项失败", err)
	}

	if existingItem == nil {
		return NewServiceError(ErrCodeNotFound, "配置项不存在", nil)
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		logger.Error("Failed to delete configuration item", err)
		return NewServiceError(ErrCodeDatabase, "删除配置项失败", err)
	}

	return nil
}

// BatchImportConfigItems 批量导入配置项（用于Excel导入）
func (s *configurationItemService) BatchImportConfigItems(ctx context.Context, items []models.ConfigurationItem) error {
	ctx = WithContext(ctx)
	logger.Info("Batch importing configuration items", zap.Int("count", len(items)))

	// 批量验证：确保所有服务商和产品的有效性
	for i, item := range items {
		// 检查服务商是否存在
		provider, err := s.providerRepo.GetByID(ctx, item.CloudProviderID)
		if err != nil {
			logger.Error("Failed to check provider existence", err, 
				zap.Uint("providerId", item.CloudProviderID),
				zap.Int("index", i))
			return NewServiceError(ErrCodeDatabase, "批量导入配置项失败：验证云服务商出错", err)
		}

		if provider == nil {
			return NewServiceError(ErrCodeNotFound, 
				"批量导入配置项失败：第" + string(i+1) + "条记录的云服务商不存在", nil)
		}

		// 检查产品是否存在
		product, err := s.productRepo.GetByID(ctx, item.ProductID)
		if err != nil {
			logger.Error("Failed to check product existence", err, 
				zap.Uint("productId", item.ProductID),
				zap.Int("index", i))
			return NewServiceError(ErrCodeDatabase, "批量导入配置项失败：验证云产品出错", err)
		}

		if product == nil {
			return NewServiceError(ErrCodeNotFound, 
				"批量导入配置项失败：第" + string(i+1) + "条记录的云产品不存在", nil)
		}

		// 检查产品是否属于指定的服务商
		if product.CloudProviderID != item.CloudProviderID {
			return NewServiceError(ErrCodeInvalidData, 
				"批量导入配置项失败：第" + string(i+1) + "条记录的云产品不属于指定的云服务商", nil)
		}
	}

	if err := s.repo.BatchInsert(ctx, items); err != nil {
		logger.Error("Failed to batch import configuration items", err)
		return NewServiceError(ErrCodeDatabase, "批量导入配置项失败", err)
	}

	return nil
}
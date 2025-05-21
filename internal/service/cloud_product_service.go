package service

import (
	"context"

	"github.com/yourusername/cloud-eye/internal/models"
	"github.com/yourusername/cloud-eye/internal/pkg/logger"
	"github.com/yourusername/cloud-eye/internal/repository"
	"go.uber.org/zap"
)

// CloudProductService 云产品服务接口
type CloudProductService interface {
	Service
	GetAllProducts(ctx context.Context) ([]models.CloudProduct, error)
	GetProductByID(ctx context.Context, id uint) (*models.CloudProduct, error)
	GetProductsByProviderID(ctx context.Context, providerID uint) ([]models.CloudProduct, error)
	GetProductsByProviderCode(ctx context.Context, providerCode string) ([]models.CloudProduct, error)
	CreateProduct(ctx context.Context, product *models.CloudProduct) error
	UpdateProduct(ctx context.Context, product *models.CloudProduct) error
	DeleteProduct(ctx context.Context, id uint) error
}

// cloudProductService 云产品服务实现
type cloudProductService struct {
	BaseService
	repo            repository.CloudProductRepository
	providerRepo    repository.CloudProviderRepository
}

// NewCloudProductService 创建云产品服务
func NewCloudProductService(
	repo repository.CloudProductRepository,
	providerRepo repository.CloudProviderRepository,
) CloudProductService {
	return &cloudProductService{
		repo:         repo,
		providerRepo: providerRepo,
	}
}

// GetAllProducts 获取所有云产品
func (s *cloudProductService) GetAllProducts(ctx context.Context) ([]models.CloudProduct, error) {
	ctx = WithContext(ctx)
	logger.Info("Getting all cloud products")

	products, err := s.repo.GetAll(ctx)
	if err != nil {
		logger.Error("Failed to get all cloud products", err)
		return nil, NewServiceError(ErrCodeDatabase, "获取云产品列表失败", err)
	}

	return products, nil
}

// GetProductByID 根据ID获取云产品
func (s *cloudProductService) GetProductByID(ctx context.Context, id uint) (*models.CloudProduct, error) {
	ctx = WithContext(ctx)
	logger.Info("Getting cloud product by ID", zap.Uint("id", id))

	product, err := s.repo.GetByID(ctx, id)
	if err != nil {
		logger.Error("Failed to get cloud product by ID", err, zap.Uint("id", id))
		return nil, NewServiceError(ErrCodeDatabase, "获取云产品详情失败", err)
	}

	if product == nil {
		return nil, NewServiceError(ErrCodeNotFound, "云产品不存在", nil)
	}

	return product, nil
}

// GetProductsByProviderID 根据云服务商ID获取云产品
func (s *cloudProductService) GetProductsByProviderID(ctx context.Context, providerID uint) ([]models.CloudProduct, error) {
	ctx = WithContext(ctx)
	logger.Info("Getting cloud products by provider ID", zap.Uint("providerId", providerID))

	// 检查服务商是否存在
	provider, err := s.providerRepo.GetByID(ctx, providerID)
	if err != nil {
		logger.Error("Failed to check provider existence", err, zap.Uint("providerId", providerID))
		return nil, NewServiceError(ErrCodeDatabase, "获取云产品列表失败", err)
	}

	if provider == nil {
		return nil, NewServiceError(ErrCodeNotFound, "云服务商不存在", nil)
	}

	products, err := s.repo.GetByProviderID(ctx, providerID)
	if err != nil {
		logger.Error("Failed to get cloud products by provider ID", err, zap.Uint("providerId", providerID))
		return nil, NewServiceError(ErrCodeDatabase, "获取云产品列表失败", err)
	}

	return products, nil
}

// GetProductsByProviderCode 根据云服务商代码获取云产品
func (s *cloudProductService) GetProductsByProviderCode(ctx context.Context, providerCode string) ([]models.CloudProduct, error) {
	ctx = WithContext(ctx)
	logger.Info("Getting cloud products by provider code", zap.String("providerCode", providerCode))

	// 检查服务商是否存在
	provider, err := s.providerRepo.GetByCode(ctx, providerCode)
	if err != nil {
		logger.Error("Failed to check provider existence", err, zap.String("providerCode", providerCode))
		return nil, NewServiceError(ErrCodeDatabase, "获取云产品列表失败", err)
	}

	if provider == nil {
		return nil, NewServiceError(ErrCodeNotFound, "云服务商不存在", nil)
	}

	products, err := s.repo.GetByProviderCode(ctx, providerCode)
	if err != nil {
		logger.Error("Failed to get cloud products by provider code", err, zap.String("providerCode", providerCode))
		return nil, NewServiceError(ErrCodeDatabase, "获取云产品列表失败", err)
	}

	return products, nil
}

// CreateProduct 创建云产品
func (s *cloudProductService) CreateProduct(ctx context.Context, product *models.CloudProduct) error {
	ctx = WithContext(ctx)
	logger.Info("Creating cloud product", 
		zap.String("name", product.Name), 
		zap.String("code", product.Code),
		zap.Uint("providerId", product.CloudProviderID))

	// 检查服务商是否存在
	provider, err := s.providerRepo.GetByID(ctx, product.CloudProviderID)
	if err != nil {
		logger.Error("Failed to check provider existence", err, zap.Uint("providerId", product.CloudProviderID))
		return NewServiceError(ErrCodeDatabase, "创建云产品失败", err)
	}

	if provider == nil {
		return NewServiceError(ErrCodeNotFound, "云服务商不存在", nil)
	}

	// 检查代码是否已存在于该服务商下
	existingProduct, err := s.repo.GetByCode(ctx, product.CloudProviderID, product.Code)
	if err != nil {
		logger.Error("Failed to check product code", err, zap.String("code", product.Code))
		return NewServiceError(ErrCodeDatabase, "创建云产品失败", err)
	}

	if existingProduct != nil {
		return NewServiceError(ErrCodeDuplicate, "该服务商下云产品代码已存在", nil)
	}

	if err := s.repo.Create(ctx, product); err != nil {
		logger.Error("Failed to create cloud product", err)
		return NewServiceError(ErrCodeDatabase, "创建云产品失败", err)
	}

	return nil
}

// UpdateProduct 更新云产品
func (s *cloudProductService) UpdateProduct(ctx context.Context, product *models.CloudProduct) error {
	ctx = WithContext(ctx)
	logger.Info("Updating cloud product", zap.Uint("id", product.ID))

	// 检查产品是否存在
	existingProduct, err := s.repo.GetByID(ctx, product.ID)
	if err != nil {
		logger.Error("Failed to check product existence", err, zap.Uint("id", product.ID))
		return NewServiceError(ErrCodeDatabase, "更新云产品失败", err)
	}

	if existingProduct == nil {
		return NewServiceError(ErrCodeNotFound, "云产品不存在", nil)
	}

	// 如果更改了服务商，检查服务商是否存在
	if product.CloudProviderID != existingProduct.CloudProviderID {
		provider, err := s.providerRepo.GetByID(ctx, product.CloudProviderID)
		if err != nil {
			logger.Error("Failed to check provider existence", err, zap.Uint("providerId", product.CloudProviderID))
			return NewServiceError(ErrCodeDatabase, "更新云产品失败", err)
		}

		if provider == nil {
			return NewServiceError(ErrCodeNotFound, "云服务商不存在", nil)
		}
	}

	// 如果更改了代码或服务商，检查新代码是否已存在
	if product.Code != existingProduct.Code || product.CloudProviderID != existingProduct.CloudProviderID {
		codeCheck, err := s.repo.GetByCode(ctx, product.CloudProviderID, product.Code)
		if err != nil {
			logger.Error("Failed to check product code", err, zap.String("code", product.Code))
			return NewServiceError(ErrCodeDatabase, "更新云产品失败", err)
		}

		if codeCheck != nil && codeCheck.ID != product.ID {
			return NewServiceError(ErrCodeDuplicate, "该服务商下云产品代码已存在", nil)
		}
	}

	if err := s.repo.Update(ctx, product); err != nil {
		logger.Error("Failed to update cloud product", err)
		return NewServiceError(ErrCodeDatabase, "更新云产品失败", err)
	}

	return nil
}

// DeleteProduct 删除云产品
func (s *cloudProductService) DeleteProduct(ctx context.Context, id uint) error {
	ctx = WithContext(ctx)
	logger.Info("Deleting cloud product", zap.Uint("id", id))

	// 检查产品是否存在
	existingProduct, err := s.repo.GetByID(ctx, id)
	if err != nil {
		logger.Error("Failed to check product existence", err, zap.Uint("id", id))
		return NewServiceError(ErrCodeDatabase, "删除云产品失败", err)
	}

	if existingProduct == nil {
		return NewServiceError(ErrCodeNotFound, "云产品不存在", nil)
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		logger.Error("Failed to delete cloud product", err)
		return NewServiceError(ErrCodeDatabase, "删除云产品失败", err)
	}

	return nil
}
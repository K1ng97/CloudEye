package service

import (
	"context"

	"github.com/yourusername/cloud-eye/internal/models"
	"github.com/yourusername/cloud-eye/internal/pkg/logger"
	"github.com/yourusername/cloud-eye/internal/repository"
	"go.uber.org/zap"
)

// CloudProviderService 云服务商服务接口
type CloudProviderService interface {
	Service
	GetAllProviders(ctx context.Context) ([]models.CloudProvider, error)
	GetProviderByID(ctx context.Context, id uint) (*models.CloudProvider, error)
	GetProviderByCode(ctx context.Context, code string) (*models.CloudProvider, error)
	CreateProvider(ctx context.Context, provider *models.CloudProvider) error
	UpdateProvider(ctx context.Context, provider *models.CloudProvider) error
	DeleteProvider(ctx context.Context, id uint) error
}

// cloudProviderService 云服务商服务实现
type cloudProviderService struct {
	BaseService
	repo repository.CloudProviderRepository
}

// NewCloudProviderService 创建云服务商服务
func NewCloudProviderService(repo repository.CloudProviderRepository) CloudProviderService {
	return &cloudProviderService{
		repo: repo,
	}
}

// GetAllProviders 获取所有云服务商
func (s *cloudProviderService) GetAllProviders(ctx context.Context) ([]models.CloudProvider, error) {
	ctx = WithContext(ctx)
	logger.Info("Getting all cloud providers")

	providers, err := s.repo.GetAll(ctx)
	if err != nil {
		logger.Error("Failed to get all cloud providers", err)
		return nil, NewServiceError(ErrCodeDatabase, "获取云服务商列表失败", err)
	}

	return providers, nil
}

// GetProviderByID 根据ID获取云服务商
func (s *cloudProviderService) GetProviderByID(ctx context.Context, id uint) (*models.CloudProvider, error) {
	ctx = WithContext(ctx)
	logger.Info("Getting cloud provider by ID", zap.Uint("id", id))

	provider, err := s.repo.GetByID(ctx, id)
	if err != nil {
		logger.Error("Failed to get cloud provider by ID", err, zap.Uint("id", id))
		return nil, NewServiceError(ErrCodeDatabase, "获取云服务商详情失败", err)
	}

	if provider == nil {
		return nil, NewServiceError(ErrCodeNotFound, "云服务商不存在", nil)
	}

	return provider, nil
}

// GetProviderByCode 根据代码获取云服务商
func (s *cloudProviderService) GetProviderByCode(ctx context.Context, code string) (*models.CloudProvider, error) {
	ctx = WithContext(ctx)
	logger.Info("Getting cloud provider by code", zap.String("code", code))

	provider, err := s.repo.GetByCode(ctx, code)
	if err != nil {
		logger.Error("Failed to get cloud provider by code", err, zap.String("code", code))
		return nil, NewServiceError(ErrCodeDatabase, "获取云服务商详情失败", err)
	}

	if provider == nil {
		return nil, NewServiceError(ErrCodeNotFound, "云服务商不存在", nil)
	}

	return provider, nil
}

// CreateProvider 创建云服务商
func (s *cloudProviderService) CreateProvider(ctx context.Context, provider *models.CloudProvider) error {
	ctx = WithContext(ctx)
	logger.Info("Creating cloud provider", zap.String("name", provider.Name), zap.String("code", provider.Code))

	// 检查代码是否已存在
	existingProvider, err := s.repo.GetByCode(ctx, provider.Code)
	if err != nil {
		logger.Error("Failed to check cloud provider code", err, zap.String("code", provider.Code))
		return NewServiceError(ErrCodeDatabase, "创建云服务商失败", err)
	}

	if existingProvider != nil {
		return NewServiceError(ErrCodeDuplicate, "云服务商代码已存在", nil)
	}

	if err := s.repo.Create(ctx, provider); err != nil {
		logger.Error("Failed to create cloud provider", err)
		return NewServiceError(ErrCodeDatabase, "创建云服务商失败", err)
	}

	return nil
}

// UpdateProvider 更新云服务商
func (s *cloudProviderService) UpdateProvider(ctx context.Context, provider *models.CloudProvider) error {
	ctx = WithContext(ctx)
	logger.Info("Updating cloud provider", zap.Uint("id", provider.ID))

	// 检查是否存在
	existingProvider, err := s.repo.GetByID(ctx, provider.ID)
	if err != nil {
		logger.Error("Failed to check cloud provider existence", err, zap.Uint("id", provider.ID))
		return NewServiceError(ErrCodeDatabase, "更新云服务商失败", err)
	}

	if existingProvider == nil {
		return NewServiceError(ErrCodeNotFound, "云服务商不存在", nil)
	}

	// 如果更改了代码，检查新代码是否已存在
	if provider.Code != existingProvider.Code {
		codeCheck, err := s.repo.GetByCode(ctx, provider.Code)
		if err != nil {
			logger.Error("Failed to check cloud provider code", err, zap.String("code", provider.Code))
			return NewServiceError(ErrCodeDatabase, "更新云服务商失败", err)
		}

		if codeCheck != nil && codeCheck.ID != provider.ID {
			return NewServiceError(ErrCodeDuplicate, "云服务商代码已存在", nil)
		}
	}

	if err := s.repo.Update(ctx, provider); err != nil {
		logger.Error("Failed to update cloud provider", err)
		return NewServiceError(ErrCodeDatabase, "更新云服务商失败", err)
	}

	return nil
}

// DeleteProvider 删除云服务商
func (s *cloudProviderService) DeleteProvider(ctx context.Context, id uint) error {
	ctx = WithContext(ctx)
	logger.Info("Deleting cloud provider", zap.Uint("id", id))

	// 检查是否存在
	existingProvider, err := s.repo.GetByID(ctx, id)
	if err != nil {
		logger.Error("Failed to check cloud provider existence", err, zap.Uint("id", id))
		return NewServiceError(ErrCodeDatabase, "删除云服务商失败", err)
	}

	if existingProvider == nil {
		return NewServiceError(ErrCodeNotFound, "云服务商不存在", nil)
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		logger.Error("Failed to delete cloud provider", err)
		return NewServiceError(ErrCodeDatabase, "删除云服务商失败", err)
	}

	return nil
}
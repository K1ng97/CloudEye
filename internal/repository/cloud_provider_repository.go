package repository

import (
	"context"
	"errors"

	"github.com/yourusername/cloud-eye/internal/models"
	"github.com/yourusername/cloud-eye/internal/pkg/logger"
	"gorm.io/gorm"
)

// CloudProviderRepository 云服务商仓库接口
type CloudProviderRepository interface {
	Repository
	GetAll(ctx context.Context) ([]models.CloudProvider, error)
	GetByID(ctx context.Context, id uint) (*models.CloudProvider, error)
	GetByCode(ctx context.Context, code string) (*models.CloudProvider, error)
	Create(ctx context.Context, provider *models.CloudProvider) error
	Update(ctx context.Context, provider *models.CloudProvider) error
	Delete(ctx context.Context, id uint) error
}

// cloudProviderRepository 云服务商仓库实现
type cloudProviderRepository struct {
	BaseRepository
}

// NewCloudProviderRepository 创建云服务商仓库
func NewCloudProviderRepository(db *gorm.DB) CloudProviderRepository {
	return &cloudProviderRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

// GetAll 获取所有云服务商
func (r *cloudProviderRepository) GetAll(ctx context.Context) ([]models.CloudProvider, error) {
	var providers []models.CloudProvider
	err := r.DB.WithContext(ctx).Find(&providers).Error
	if err != nil {
		logger.Error("Failed to get all cloud providers", err)
		return nil, err
	}
	return providers, nil
}

// GetByID 根据ID获取云服务商
func (r *cloudProviderRepository) GetByID(ctx context.Context, id uint) (*models.CloudProvider, error) {
	var provider models.CloudProvider
	err := r.DB.WithContext(ctx).First(&provider, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		logger.Error("Failed to get cloud provider by ID", err)
		return nil, err
	}
	return &provider, nil
}

// GetByCode 根据代码获取云服务商
func (r *cloudProviderRepository) GetByCode(ctx context.Context, code string) (*models.CloudProvider, error) {
	var provider models.CloudProvider
	err := r.DB.WithContext(ctx).Where("code = ?", code).First(&provider).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		logger.Error("Failed to get cloud provider by code", err)
		return nil, err
	}
	return &provider, nil
}

// Create 创建云服务商
func (r *cloudProviderRepository) Create(ctx context.Context, provider *models.CloudProvider) error {
	err := r.DB.WithContext(ctx).Create(provider).Error
	if err != nil {
		logger.Error("Failed to create cloud provider", err)
		return err
	}
	return nil
}

// Update 更新云服务商
func (r *cloudProviderRepository) Update(ctx context.Context, provider *models.CloudProvider) error {
	err := r.DB.WithContext(ctx).Save(provider).Error
	if err != nil {
		logger.Error("Failed to update cloud provider", err)
		return err
	}
	return nil
}

// Delete 删除云服务商
func (r *cloudProviderRepository) Delete(ctx context.Context, id uint) error {
	err := r.DB.WithContext(ctx).Delete(&models.CloudProvider{}, id).Error
	if err != nil {
		logger.Error("Failed to delete cloud provider", err)
		return err
	}
	return nil
}
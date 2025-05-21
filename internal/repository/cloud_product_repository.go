package repository

import (
	"context"
	"errors"

	"github.com/yourusername/cloud-eye/internal/models"
	"github.com/yourusername/cloud-eye/internal/pkg/logger"
	"gorm.io/gorm"
)

// CloudProductRepository 云产品仓库接口
type CloudProductRepository interface {
	Repository
	GetAll(ctx context.Context) ([]models.CloudProduct, error)
	GetByID(ctx context.Context, id uint) (*models.CloudProduct, error)
	GetByProviderID(ctx context.Context, providerID uint) ([]models.CloudProduct, error)
	GetByProviderCode(ctx context.Context, providerCode string) ([]models.CloudProduct, error)
	GetByCode(ctx context.Context, providerID uint, code string) (*models.CloudProduct, error)
	Create(ctx context.Context, product *models.CloudProduct) error
	Update(ctx context.Context, product *models.CloudProduct) error
	Delete(ctx context.Context, id uint) error
}

// cloudProductRepository 云产品仓库实现
type cloudProductRepository struct {
	BaseRepository
}

// NewCloudProductRepository 创建云产品仓库
func NewCloudProductRepository(db *gorm.DB) CloudProductRepository {
	return &cloudProductRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

// GetAll 获取所有云产品
func (r *cloudProductRepository) GetAll(ctx context.Context) ([]models.CloudProduct, error) {
	var products []models.CloudProduct
	err := r.DB.WithContext(ctx).Find(&products).Error
	if err != nil {
		logger.Error("Failed to get all cloud products", err)
		return nil, err
	}
	return products, nil
}

// GetByID 根据ID获取云产品
func (r *cloudProductRepository) GetByID(ctx context.Context, id uint) (*models.CloudProduct, error) {
	var product models.CloudProduct
	err := r.DB.WithContext(ctx).First(&product, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		logger.Error("Failed to get cloud product by ID", err)
		return nil, err
	}
	return &product, nil
}

// GetByProviderID 根据云服务商ID获取云产品
func (r *cloudProductRepository) GetByProviderID(ctx context.Context, providerID uint) ([]models.CloudProduct, error) {
	var products []models.CloudProduct
	err := r.DB.WithContext(ctx).Where("cloud_provider_id = ?", providerID).Find(&products).Error
	if err != nil {
		logger.Error("Failed to get cloud products by provider ID", err)
		return nil, err
	}
	return products, nil
}

// GetByProviderCode 根据云服务商代码获取云产品
func (r *cloudProductRepository) GetByProviderCode(ctx context.Context, providerCode string) ([]models.CloudProduct, error) {
	var products []models.CloudProduct
	err := r.DB.WithContext(ctx).
		Joins("JOIN cloud_providers ON cloud_products.cloud_provider_id = cloud_providers.id").
		Where("cloud_providers.code = ?", providerCode).
		Find(&products).Error
	if err != nil {
		logger.Error("Failed to get cloud products by provider code", err)
		return nil, err
	}
	return products, nil
}

// GetByCode 根据代码和云服务商ID获取云产品
func (r *cloudProductRepository) GetByCode(ctx context.Context, providerID uint, code string) (*models.CloudProduct, error) {
	var product models.CloudProduct
	err := r.DB.WithContext(ctx).
		Where("cloud_provider_id = ? AND code = ?", providerID, code).
		First(&product).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		logger.Error("Failed to get cloud product by code", err)
		return nil, err
	}
	return &product, nil
}

// Create 创建云产品
func (r *cloudProductRepository) Create(ctx context.Context, product *models.CloudProduct) error {
	err := r.DB.WithContext(ctx).Create(product).Error
	if err != nil {
		logger.Error("Failed to create cloud product", err)
		return err
	}
	return nil
}

// Update 更新云产品
func (r *cloudProductRepository) Update(ctx context.Context, product *models.CloudProduct) error {
	err := r.DB.WithContext(ctx).Save(product).Error
	if err != nil {
		logger.Error("Failed to update cloud product", err)
		return err
	}
	return nil
}

// Delete 删除云产品
func (r *cloudProductRepository) Delete(ctx context.Context, id uint) error {
	err := r.DB.WithContext(ctx).Delete(&models.CloudProduct{}, id).Error
	if err != nil {
		logger.Error("Failed to delete cloud product", err)
		return err
	}
	return nil
}
package repository

import (
	"context"

	"github.com/yourusername/cloud-eye/internal/models"
	"gorm.io/gorm"
)

// Repository 定义了所有仓库的通用接口
type Repository interface {
	Close() error
}

// BaseRepository 基础仓库实现，提供通用的数据库操作方法
type BaseRepository struct {
	DB *gorm.DB
}

// NewBaseRepository 创建一个新的基础仓库
func NewBaseRepository(db *gorm.DB) BaseRepository {
	return BaseRepository{DB: db}
}

// Close 关闭数据库连接
func (r *BaseRepository) Close() error {
	sqlDB, err := r.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// Transaction 在事务中执行操作
func (r *BaseRepository) Transaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	return r.DB.WithContext(ctx).Transaction(fn)
}

// Paginate 分页查询辅助函数
func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page <= 0 {
			page = 1
		}

		if pageSize <= 0 {
			pageSize = 10
		}

		if pageSize > 100 {
			pageSize = 100
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

// PageResult 分页结果
type PageResult struct {
	Total    int64       `json:"total"`     // 总记录数
	Page     int         `json:"page"`      // 当前页码
	PageSize int         `json:"page_size"` // 每页大小
	Data     interface{} `json:"data"`      // 数据列表
}
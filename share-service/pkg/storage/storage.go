package storage

import (
	"context"

	"github.com/playground/share-service/pkg/models"
)

// Storage 定义了存储层的接口
type Storage interface {
	// CreateShare 创建新的分享
	CreateShare(ctx context.Context, share *models.Share) error

	// GetShare 通过 shareId 获取分享
	GetShare(ctx context.Context, shareId string) (*models.Share, error)

	// IncrementViews 增加分享的访问次数
	IncrementViews(ctx context.Context, shareId string) error

	// DeleteExpiredShares 删除过期的分享
	DeleteExpiredShares(ctx context.Context) error

	// Close 关闭存储连接
	Close(ctx context.Context) error
}

// Cache 定义了缓存层的接口
type Cache interface {
	// SetShare 将分享存入缓存
	SetShare(ctx context.Context, share *models.Share) error

	// GetShare 从缓存获取分享
	GetShare(ctx context.Context, shareId string) (*models.Share, error)

	// DeleteShare 从缓存删除分享
	DeleteShare(ctx context.Context, shareId string) error

	// Close 关闭缓存连接
	Close() error
}

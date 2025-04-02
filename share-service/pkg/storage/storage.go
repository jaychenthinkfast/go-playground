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

	// IncrementViews 增加分享的访问次数，返回更新后的计数
	IncrementViews(ctx context.Context, shareId string) (int64, error)

	// DeleteExpiredShares 删除过期的分享
	DeleteExpiredShares(ctx context.Context) error

	// Close 关闭存储连接
	Close(ctx context.Context) error
}

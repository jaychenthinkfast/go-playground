package redis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/playground/share-service/pkg/models"
)

type RedisCache struct {
	client *redis.Client
}

const (
	shareKeyPrefix = "share:"
	defaultTTL     = 24 * time.Hour
)

// NewRedisCache 创建新的 Redis 缓存实例
func NewRedisCache(ctx context.Context, uri string) (*RedisCache, error) {
	opt, err := redis.ParseURL(uri)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(opt)

	// 测试连接
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &RedisCache{
		client: client,
	}, nil
}

// SetShare 实现 Cache 接口
func (c *RedisCache) SetShare(ctx context.Context, share *models.Share) error {
	data, err := json.Marshal(share)
	if err != nil {
		return err
	}

	key := shareKeyPrefix + share.ShareID
	ttl := defaultTTL
	if share.ExpiresAt != nil {
		ttl = time.Until(*share.ExpiresAt)
		if ttl < 0 {
			return c.client.Del(ctx, key).Err()
		}
	}

	return c.client.Set(ctx, key, data, ttl).Err()
}

// GetShare 实现 Cache 接口
func (c *RedisCache) GetShare(ctx context.Context, shareId string) (*models.Share, error) {
	data, err := c.client.Get(ctx, shareKeyPrefix+shareId).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var share models.Share
	if err := json.Unmarshal(data, &share); err != nil {
		return nil, err
	}

	return &share, nil
}

// DeleteShare 实现 Cache 接口
func (c *RedisCache) DeleteShare(ctx context.Context, shareId string) error {
	return c.client.Del(ctx, shareKeyPrefix+shareId).Err()
}

// Close 实现 Cache 接口
func (c *RedisCache) Close() error {
	return c.client.Close()
}

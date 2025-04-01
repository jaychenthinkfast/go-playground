package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/playground/share-service/pkg/models"
)

type RedisCache struct {
	client *redis.Client
}

// NewRedisCache 创建一个新的 Redis 缓存实例
func NewRedisCache(redisURI string) (*RedisCache, error) {
	opt, err := redis.ParseURL(redisURI)
	if err != nil {
		return nil, fmt.Errorf("parse redis uri: %w", err)
	}

	client := redis.NewClient(opt)

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("ping redis: %w", err)
	}

	return &RedisCache{
		client: client,
	}, nil
}

// SetShare 将分享存储到缓存中
func (c *RedisCache) SetShare(ctx context.Context, share *models.Share) error {
	data, err := json.Marshal(share)
	if err != nil {
		return fmt.Errorf("marshal share: %w", err)
	}

	key := fmt.Sprintf("share:%s", share.ShareID)
	// 设置过期时间与分享过期时间一致
	var ttl time.Duration
	if share.ExpiresAt != nil {
		ttl = time.Until(*share.ExpiresAt)
		if ttl <= 0 {
			// 如果已经过期，则删除
			return c.DeleteShare(ctx, share.ShareID)
		}
	} else {
		ttl = 7 * 24 * time.Hour // 默认7天过期
	}

	if err := c.client.Set(ctx, key, data, ttl).Err(); err != nil {
		return fmt.Errorf("set share: %w", err)
	}

	return nil
}

// GetShare 从缓存中获取分享
func (c *RedisCache) GetShare(ctx context.Context, shareID string) (*models.Share, error) {
	key := fmt.Sprintf("share:%s", shareID)
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // 缓存未命中
		}
		return nil, fmt.Errorf("get share: %w", err)
	}

	var share models.Share
	if err := json.Unmarshal(data, &share); err != nil {
		return nil, fmt.Errorf("unmarshal share: %w", err)
	}

	return &share, nil
}

// DeleteShare 从缓存中删除分享
func (c *RedisCache) DeleteShare(ctx context.Context, shareID string) error {
	key := fmt.Sprintf("share:%s", shareID)
	if err := c.client.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("delete share: %w", err)
	}
	return nil
}

// IncrementViews 增加分享的访问次数
func (c *RedisCache) IncrementViews(ctx context.Context, shareID string) error {
	key := fmt.Sprintf("views:%s", shareID)
	if err := c.client.Incr(ctx, key).Err(); err != nil {
		return fmt.Errorf("increment views: %w", err)
	}
	return nil
}

// GetViews 获取分享的访问次数
func (c *RedisCache) GetViews(ctx context.Context, shareID string) (int64, error) {
	key := fmt.Sprintf("views:%s", shareID)
	views, err := c.client.Get(ctx, key).Int64()
	if err != nil {
		if err == redis.Nil {
			return 0, nil
		}
		return 0, fmt.Errorf("get views: %w", err)
	}
	return views, nil
}

// IsRateLimited 检查是否超过访问频率限制
func (c *RedisCache) IsRateLimited(ctx context.Context, ip string, endpoint string) (bool, error) {
	key := fmt.Sprintf("rate:%s:%s", ip, endpoint)

	// 使用 MULTI/EXEC 保证原子性
	pipe := c.client.Pipeline()
	incr := pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, 60*time.Second) // 60秒的时间窗口

	_, err := pipe.Exec(ctx)
	if err != nil {
		return false, fmt.Errorf("check rate limit: %w", err)
	}

	// 获取当前计数
	count := incr.Val()
	return count > 60, nil // 每分钟最多60次请求
}

// SetRunResult 存储代码运行结果
func (c *RedisCache) SetRunResult(ctx context.Context, taskID string, result *models.RunResult) error {
	data, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("marshal result: %w", err)
	}

	key := fmt.Sprintf("result:%s", taskID)
	if err := c.client.Set(ctx, key, data, 5*time.Minute).Err(); err != nil {
		return fmt.Errorf("set result: %w", err)
	}

	return nil
}

// GetRunResult 获取代码运行结果
func (c *RedisCache) GetRunResult(ctx context.Context, taskID string) (*models.RunResult, error) {
	key := fmt.Sprintf("result:%s", taskID)
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, fmt.Errorf("get result: %w", err)
	}

	var result models.RunResult
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("unmarshal result: %w", err)
	}

	return &result, nil
}

// Close 关闭 Redis 连接
func (c *RedisCache) Close() error {
	return c.client.Close()
}

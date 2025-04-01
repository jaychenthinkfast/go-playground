package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/playground/share-service/pkg/models"
	"github.com/playground/share-service/pkg/storage"
)

type Handler struct {
	storage storage.Storage
	cache   storage.Cache
}

func NewHandler(storage storage.Storage, cache storage.Cache) *Handler {
	return &Handler{
		storage: storage,
		cache:   cache,
	}
}

// HealthCheck 处理健康检查请求
func (h *Handler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"time":   time.Now(),
	})
}

// CreateShare 处理创建分享请求
func (h *Handler) CreateShare(c *gin.Context) {
	var req models.CreateShareRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 生成唯一ID
	shareId := uuid.New().String()[:8]

	// 创建分享对象
	share := &models.Share{
		ShareID:     shareId,
		Code:        req.Code,
		Language:    "go", // 目前只支持 Go
		Version:     req.Version,
		Title:       req.Title,
		Description: req.Description,
		Author:      req.Author,
		CreatedAt:   time.Now(),
		Views:       0,
	}

	// 处理过期时间
	if req.ExpiresIn != "" {
		duration, err := time.ParseDuration(req.ExpiresIn)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid expires_in format"})
			return
		}
		expiresAt := time.Now().Add(duration)
		share.ExpiresAt = &expiresAt
	}

	// 保存到存储
	if err := h.storage.CreateShare(c.Request.Context(), share); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create share"})
		return
	}

	// 保存到缓存
	if err := h.cache.SetShare(c.Request.Context(), share); err != nil {
		// 缓存错误不影响主流程，只记录日志
		fmt.Printf("failed to cache share: %v\n", err)
	}

	// 构建响应
	resp := &models.CreateShareResponse{
		ShareID:   shareId,
		URL:       fmt.Sprintf("/share/%s", shareId),
		ExpiresAt: share.ExpiresAt,
	}

	c.JSON(http.StatusCreated, resp)
}

// GetShare 处理获取分享请求
func (h *Handler) GetShare(c *gin.Context) {
	shareId := c.Param("id")

	// 先从缓存获取
	share, err := h.cache.GetShare(c.Request.Context(), shareId)
	if err != nil {
		fmt.Printf("failed to get share from cache: %v\n", err)
	}

	// 缓存未命中，从存储获取
	if share == nil {
		share, err = h.storage.GetShare(c.Request.Context(), shareId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get share"})
			return
		}
		if share == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "share not found"})
			return
		}

		// 更新缓存
		if err := h.cache.SetShare(c.Request.Context(), share); err != nil {
			fmt.Printf("failed to cache share: %v\n", err)
		}
	}

	// 检查是否过期
	if share.ExpiresAt != nil && share.ExpiresAt.Before(time.Now()) {
		c.JSON(http.StatusNotFound, gin.H{"error": "share has expired"})
		return
	}

	// 构建响应
	resp := &models.GetShareResponse{
		Code:        share.Code,
		Version:     share.Version,
		Title:       share.Title,
		Description: share.Description,
		Author:      share.Author,
		CreatedAt:   share.CreatedAt,
		ExpiresAt:   share.ExpiresAt,
		Views:       share.Views,
	}

	c.JSON(http.StatusOK, resp)
}

// IncrementViews 处理增加访问次数请求
func (h *Handler) IncrementViews(c *gin.Context) {
	shareId := c.Param("id")

	// 增加存储中的访问次数
	if err := h.storage.IncrementViews(c.Request.Context(), shareId); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to increment views"})
		return
	}

	// 同时增加缓存中的访问次数
	if err := h.cache.IncrementViews(c.Request.Context(), shareId); err != nil {
		fmt.Printf("failed to increment views in cache: %v\n", err)
	}

	c.Status(http.StatusOK)
}

// ExecuteCode 处理代码执行请求
func (h *Handler) ExecuteCode(c *gin.Context) {
	var req struct {
		Code    string `json:"code" binding:"required"`
		Version string `json:"version" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查是否超过频率限制
	ip := c.ClientIP()
	limited, err := h.cache.IsRateLimited(c.Request.Context(), ip, "execute")
	if err != nil {
		fmt.Printf("failed to check rate limit: %v\n", err)
	}
	if limited {
		c.JSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded"})
		return
	}

	// 生成唯一任务ID
	taskID := uuid.New().String()

	// 这里会有实际的代码执行逻辑，可能是异步的
	// 为了演示，我们创建一个模拟的结果
	result := &models.RunResult{
		Output:    "Hello, World!",
		ExitCode:  0,
		Duration:  100,         // 假设耗时100ms
		Memory:    1024 * 1024, // 假设使用1MB内存
		CreatedAt: time.Now().Unix(),
	}

	// 存储结果到缓存
	if err := h.cache.SetRunResult(c.Request.Context(), taskID, result); err != nil {
		fmt.Printf("failed to cache run result: %v\n", err)
	}

	c.JSON(http.StatusOK, gin.H{
		"task_id": taskID,
		"result":  result,
	})
}

// GetRunResult 获取代码执行结果
func (h *Handler) GetRunResult(c *gin.Context) {
	taskID := c.Param("taskId")

	// 从缓存获取结果
	result, err := h.cache.GetRunResult(c.Request.Context(), taskID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get run result"})
		return
	}
	if result == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "run result not found or expired"})
		return
	}

	c.JSON(http.StatusOK, result)
}

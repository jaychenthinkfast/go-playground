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

	// 增加访问次数
	if err := h.storage.IncrementViews(c.Request.Context(), shareId); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to increment views"})
		return
	}

	// 从缓存删除，强制下次重新加载
	if err := h.cache.DeleteShare(c.Request.Context(), shareId); err != nil {
		fmt.Printf("failed to delete share from cache: %v\n", err)
	}

	c.Status(http.StatusOK)
}

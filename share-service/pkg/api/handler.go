package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"bytes"
	"crypto/sha256"
	"encoding/hex"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/playground/share-service/pkg/models"
	"github.com/playground/share-service/pkg/storage"
)

type Handler struct {
	storage storage.Storage
}

func NewHandler(storage storage.Storage) *Handler {
	return &Handler{
		storage: storage,
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

	// 从数据库获取分享
	share, err := h.storage.GetShare(c.Request.Context(), shareId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get share"})
		return
	}
	if share == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "share not found"})
		return
	}

	// 检查是否过期
	if share.ExpiresAt != nil && share.ExpiresAt.Before(time.Now()) {
		c.JSON(http.StatusNotFound, gin.H{"error": "share has expired"})
		return
	}

	// 检查请求头中的Cookie，看是否已经访问过
	viewedCookie, err := c.Cookie(fmt.Sprintf("viewed_%s", shareId))
	alreadyViewed := err == nil && viewedCookie == "true"

	// 仅当未被此客户端查看时才增加计数
	if !alreadyViewed {
		updatedViews, err := h.storage.IncrementViews(c.Request.Context(), shareId)
		if err != nil {
			fmt.Printf("failed to increment views: %v\n", err)
			// 继续处理请求，即使计数失败也不影响获取分享内容
		} else {
			// 增加计数成功，更新计数
			share.Views = updatedViews

			// 设置Cookie标记已访问，有效期24小时
			c.SetCookie(
				fmt.Sprintf("viewed_%s", shareId), // Cookie名称
				"true",                            // Cookie值
				86400,                             // 过期时间（秒）
				"/",                               // 路径
				"",                                // 域名
				false,                             // 仅HTTPS
				true,                              // HTTP Only
			)
		}
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

// IncrementViews 处理手动增加访问次数请求
func (h *Handler) IncrementViews(c *gin.Context) {
	shareId := c.Param("id")

	// 增加存储中的访问次数
	updatedViews, err := h.storage.IncrementViews(c.Request.Context(), shareId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to increment views"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"views": updatedViews})
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

	fmt.Printf("收到代码执行请求, 版本: %s, 代码长度: %d\n", req.Version, len(req.Code))

	// 生成唯一任务ID
	taskID := uuid.New().String()

	// 验证版本格式，确保版本格式正确
	var normalizedVersion string
	switch req.Version {
	case "go1.22", "1.22", "go1.22.0", "1.22.0":
		normalizedVersion = "go1.22"
	case "go1.23", "1.23", "go1.23.0", "1.23.0":
		normalizedVersion = "go1.23"
	case "go1.24", "1.24", "go1.24.0", "1.24.0":
		normalizedVersion = "go1.24"
	default:
		fmt.Printf("不支持的 Go 版本: %s\n", req.Version)

		// 返回错误响应
		result := &models.RunResult{
			Output:    "",
			Error:     "Unsupported Go version. Available versions: go1.24, go1.23, go1.22",
			ExitCode:  1,
			Duration:  0,
			Memory:    0,
			CreatedAt: time.Now().Unix(),
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"task_id": taskID,
			"result":  result,
			"error":   "Unsupported version",
		})
		return
	}

	// 构建请求转发到后端执行服务
	var backendURL string
	switch normalizedVersion {
	case "go1.22":
		backendURL = "http://backend-go122:3001/api/run"
	case "go1.23":
		backendURL = "http://backend-go123:3001/api/run"
	case "go1.24":
		backendURL = "http://backend-go124:3001/api/run"
	}

	// 准备发送到后端的请求
	backendReq, err := json.Marshal(map[string]interface{}{
		"code":     req.Code,
		"version":  normalizedVersion,
		"language": "go",
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to prepare backend request"})
		return
	}

	// 调用后端执行服务
	fmt.Printf("转发代码执行请求到后端服务: %s，版本: %s\n", backendURL, normalizedVersion)
	resp, err := http.Post(backendURL, "application/json", bytes.NewBuffer(backendReq))
	if err != nil {
		fmt.Printf("调用后端服务失败: %v\n", err)
		// 后端服务不可用，返回模拟结果以便测试
		mockResult := &models.RunResult{
			Output:    "Hello, World! (mock result - backend service unavailable)",
			Error:     fmt.Sprintf("后端服务不可用: %v", err),
			ExitCode:  0,
			Duration:  100,
			Memory:    1024 * 1024,
			CreatedAt: time.Now().Unix(),
		}

		c.JSON(http.StatusOK, gin.H{
			"task_id": taskID,
			"result":  mockResult,
			"mocked":  true,
		})
		return
	}
	defer resp.Body.Close()

	// 记录状态码
	fmt.Printf("后端服务响应状态码: %d\n", resp.StatusCode)

	// 读取响应内容（不管状态码如何）
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("读取后端服务响应失败: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read backend response"})
		return
	}
	fmt.Printf("后端服务响应内容: %s\n", string(respBody))

	// 检查响应状态码，非200状态码视为错误
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("后端服务返回非200状态码: %d，响应内容: %s\n", resp.StatusCode, string(respBody))
		result := &models.RunResult{
			Output:    "",
			Error:     fmt.Sprintf("后端服务错误: 状态码 %d, 响应: %s", resp.StatusCode, string(respBody)),
			ExitCode:  1,
			Duration:  0,
			Memory:    0,
			CreatedAt: time.Now().Unix(),
		}

		c.JSON(http.StatusOK, gin.H{
			"task_id": taskID,
			"result":  result,
			"error":   "Backend service error",
		})
		return
	}

	// 解析后端响应
	var backendResp map[string]interface{}
	if err := json.Unmarshal(respBody, &backendResp); err != nil {
		fmt.Printf("解析后端响应失败: %v, 响应内容: %s\n", err, string(respBody))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to decode backend response"})
		return
	}

	// 转换为我们的运行结果格式
	output := ""
	if out, ok := backendResp["output"].(string); ok {
		output = out
	}

	errMsg := ""
	if errOutput, ok := backendResp["error"].(string); ok {
		errMsg = errOutput
	}

	exitCode := 0
	if code, ok := backendResp["exitCode"].(float64); ok {
		exitCode = int(code)
	}

	duration := int64(100)
	if dur, ok := backendResp["duration"].(float64); ok {
		duration = int64(dur)
	}

	memory := int64(1024 * 1024)
	if mem, ok := backendResp["memory"].(float64); ok {
		memory = int64(mem)
	}

	// 创建结果对象
	result := &models.RunResult{
		Output:    output,
		Error:     errMsg,
		ExitCode:  exitCode,
		Duration:  duration,
		Memory:    memory,
		CreatedAt: time.Now().Unix(),
	}

	fmt.Printf("代码执行结果: 退出码=%d, 输出长度=%d, 错误长度=%d\n",
		exitCode, len(output), len(errMsg))

	c.JSON(http.StatusOK, gin.H{
		"task_id": taskID,
		"result":  result,
	})
}

// 计算代码的哈希值，用于缓存键
func calculateHash(code string) string {
	h := sha256.New()
	h.Write([]byte(code))
	return hex.EncodeToString(h.Sum(nil))
}

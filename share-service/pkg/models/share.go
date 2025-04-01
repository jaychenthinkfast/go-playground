package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Share 代表一个分享的代码片段
type Share struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	ShareID     string             `bson:"shareId" json:"shareId"`
	Code        string             `bson:"code" json:"code"`
	Language    string             `bson:"language" json:"language"`
	Version     string             `bson:"version" json:"version"`
	Title       string             `bson:"title,omitempty" json:"title,omitempty"`
	Description string             `bson:"description,omitempty" json:"description,omitempty"`
	Author      string             `bson:"author,omitempty" json:"author,omitempty"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	ExpiresAt   *time.Time         `bson:"expires_at,omitempty" json:"expires_at,omitempty"`
	Views       int64              `bson:"views" json:"views"`
	LastViewed  *time.Time         `bson:"last_viewed,omitempty" json:"last_viewed,omitempty"`
}

// CreateShareRequest 代表创建分享的请求
type CreateShareRequest struct {
	Code        string `json:"code" binding:"required"`
	Version     string `json:"version" binding:"required"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Author      string `json:"author,omitempty"`
	ExpiresIn   string `json:"expires_in,omitempty"` // 例如: "24h", "7d"
}

// CreateShareResponse 代表创建分享的响应
type CreateShareResponse struct {
	ShareID   string     `json:"shareId"`
	URL       string     `json:"url"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}

// GetShareResponse 代表获取分享的响应
type GetShareResponse struct {
	Code        string     `json:"code"`
	Version     string     `json:"version"`
	Title       string     `json:"title,omitempty"`
	Description string     `json:"description,omitempty"`
	Author      string     `json:"author,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
	Views       int64      `json:"views"`
}

# 构建阶段
FROM golang:1.22-alpine AS builder

# 设置工作目录
WORKDIR /app

# 安装必要的系统依赖
RUN apk add --no-cache git gcc musl-dev

# 设置 Go 环境变量
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOPROXY=https://goproxy.cn,direct

# 复制 Go 模块文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download && \
    go mod verify

# 复制源代码
COPY . .

# 构建应用
RUN go build -ldflags="-s -w" -o share-server ./cmd/server

# 运行阶段
FROM alpine:latest

# 安装必要的运行时依赖
RUN apk add --no-cache ca-certificates tzdata wget

# 创建非 root 用户
RUN adduser -D -g '' appuser

# 创建应用目录
RUN mkdir -p /app && chown -R appuser:appuser /app
WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder --chown=appuser:appuser /app/share-server .

# 使用非 root 用户
USER appuser

# 设置健康检查
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD wget -qO- http://localhost:${PORT:-3002}/health || exit 1

# 暴露端口
EXPOSE ${PORT:-3002}

# 运行应用
CMD ["./share-server"]

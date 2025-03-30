# 第一阶段：构建Go应用
FROM golang:1.24-alpine AS builder

# 设置工作目录
WORKDIR /app

# 设置Go环境变量
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# 复制go mod文件
COPY go.mod ./
COPY go.sum ./

# 下载依赖
RUN go mod download && go mod verify

# 复制源代码
COPY . .

# 构建应用
RUN go build -o /go-playground -ldflags="-s -w" ./cmd/server

# 第二阶段：创建最小化运行镜像
FROM alpine:latest

# 设置工作目录
WORKDIR /

# 安装必要的软件包
RUN apk add --no-cache ca-certificates tzdata wget

# 添加非root用户运行应用
RUN adduser -D -u 1000 gouser
USER gouser

# 从构建阶段复制二进制文件
COPY --from=builder --chown=gouser:gouser /go-playground /go-playground

# 确保临时目录存在并可写
USER root
RUN mkdir -p /tmp && chmod 777 /tmp
USER gouser

# 设置环境变量
ENV TZ=UTC \
    PORT=3001

# 暴露端口
EXPOSE 3001

# 设置健康检查
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD wget -qO- http://localhost:3001/api/health || exit 1

# 运行应用
ENTRYPOINT ["/go-playground"] 
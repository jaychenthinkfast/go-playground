ARG GO_VERSION=1.24
ARG GO_MOD_VERSION=1.24
FROM golang:${GO_VERSION}-alpine

# 设置工作目录
WORKDIR /app

# 设置Go环境变量
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOPROXY=https://goproxy.cn,direct

# 安装必要的软件包
RUN apk add --no-cache ca-certificates tzdata wget git

# 复制源代码
COPY . .

# 根据GO_VERSION参数修改go.mod文件中的版本
RUN rm  -rf go.mod  && go mod init go-playground

# 下载依赖
RUN go mod tidy && go mod download && go mod verify


# 构建应用
RUN go build -o /go-playground -ldflags="-s -w" ./cmd/server

# 确保临时目录存在并可写
RUN mkdir -p /tmp && chmod 777 /tmp

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
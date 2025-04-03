# Share Service 开发指南

本文档介绍如何在开发环境中运行和热更新 Share Service。

## 开发环境设置

Share Service 的开发环境使用 Docker Compose 和 Air 热重载工具，可以在代码变更时自动重新构建和重启服务。

### 环境版本

- Go 版本: 1.24
- 操作系统: Alpine Linux (容器内)
- 热重载工具: Air (github.com/air-verse/air)

### 目录结构

```
share-service/
├── cmd/           # 命令行入口
│   └── server/    # 服务器入口
├── pkg/           # 包目录
│   ├── api/       # API 处理程序
│   ├── models/    # 数据模型
│   ├── storage/   # 存储层
│   └── utils/     # 工具函数
├── docker/        # Docker 相关文件
├── .air.toml      # Air 热重载配置
├── Dockerfile     # 生产环境 Dockerfile
├── Dockerfile.dev # 开发环境 Dockerfile
├── dev.sh         # 开发环境启动脚本
├── go.mod         # Go 模块文件
└── go.sum         # Go 依赖校验
```

## 开发环境启动

### 方法一：使用脚本启动

```bash
# 进入项目根目录
cd share-service

# 启动开发环境
./dev.sh
```

### 方法二：手动启动

```bash
# 进入项目根目录
cd <项目根目录>

# 启动开发环境中的 share-service-dev 和 mongo 服务
docker-compose -f docker-compose.dev.yml up -d share-service-dev mongo

# 查看日志
docker logs -f go-playground-share-dev
```

## 热更新说明

1. 开发环境中，源代码目录被挂载到容器中的 `/app` 目录
2. Air 工具会监控源代码变化，自动重新构建和重启服务
3. 当您修改代码后，Air 将自动检测变化并重新构建应用

## 访问服务

开发环境中的服务将在以下地址可用：

- Share Service: http://localhost:3002
- Health Check: http://localhost:3002/health

## 与前端集成

### 开发环境中的服务名称

在开发环境中，Share Service 的服务名为 `share-service-dev`，而不是生产环境中的 `share-service`。前端应用需要正确配置代理以访问开发版本的服务。

### 前端代理配置

前端开发环境通过 Vue CLI 的代理配置与 Share Service 通信，配置位于 `frontend/vue.config.js` 文件中：

```javascript
// frontend/vue.config.js 中的代理配置
devServer: {
  proxy: {
    '/api/share': {
      target: 'http://share-service-dev:3002',  // 注意：开发环境中使用 share-service-dev
      changeOrigin: true,
      ws: true
    },
    '/api/execute': {
      target: 'http://share-service-dev:3002',  // 注意：开发环境中使用 share-service-dev
      changeOrigin: true,
      ws: true
    },
    // ...其他代理配置
  }
}
```

请确保在进行开发时，`vue.config.js` 中的代理配置正确指向了 `share-service-dev` 而不是 `share-service`，否则分享功能将无法正常工作。

## 开发提示

1. 查看 Air 的日志输出，了解构建和重启状态
2. 如果遇到任何问题，可以查看容器日志：`docker logs -f go-playground-share-dev`
3. 如需进入容器内部调试：`docker exec -it go-playground-share-dev sh`

## 常见问题

### Air 没有检测到文件变化

确保您的文件扩展名包含在 `.air.toml` 的 `include_ext` 配置中。

### 构建失败

检查容器日志以获取更多信息，可能是代码错误或依赖问题。

### 分享功能在开发环境中失败

如果在开发环境中分享功能返回"分享失败，请稍后重试"，请检查 `frontend/vue.config.js` 中的代理配置是否正确指向了 `share-service-dev`，而不是 `share-service`。 
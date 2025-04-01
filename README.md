# Go Playground

一个现代化的在线 Go 代码运行和分享平台。支持多个 Go 版本（1.22、1.23、1.24），提供代码运行、格式化和分享功能。

## 功能特性

- 🚀 支持多个 Go 版本（1.22、1.23、1.24）
- 💻 在线编辑和运行 Go 代码
- 🎨 代码格式化
- 📤 代码分享功能
- 🔄 实时运行结果
- 📱 响应式设计，支持移动端访问

## 技术栈

### 前端
- React + TypeScript
- Vite 构建工具
- TailwindCSS 样式框架
- Monaco Editor 代码编辑器

### 后端
- Go (支持 1.22、1.23、1.24 版本)
- Gin Web 框架
- MongoDB 数据库
- Redis 缓存

### 基础设施
- Docker + Docker Compose
- Nginx 反向代理
- 容器化部署

## 快速开始

### 前置要求

- Docker 24.0.0 或更高版本
- Docker Compose v2.24.0 或更高版本
- 至少 4GB 可用内存
- 至少 10GB 可用磁盘空间

### 部署步骤

1. 克隆项目
```bash
git clone <repository-url>
cd playground
```

2. 启动服务
```bash
docker compose up --build
```

3. 访问应用
打开浏览器访问 http://localhost:3003

## 项目结构

```
.
├── frontend/                # 前端代码
├── backend/                 # Go 后端代码
├── share-service/          # 分享服务
├── docker/                 # Docker 配置文件
│   ├── mongo/             # MongoDB 配置
│   └── redis/             # Redis 配置
└── docker-compose.yml     # Docker Compose 配置
```

## 开发指南

### 本地开发环境设置

1. 前端开发
```bash
cd frontend
npm install
npm run dev
```

2. 后端开发
```bash
cd backend
go mod download
go run cmd/server/main.go
```

3. 分享服务开发
```bash
cd share-service
go mod download
go run cmd/server/main.go
```

### 环境变量配置

- `MONGO_URI`: MongoDB 连接字符串
- `REDIS_URI`: Redis 连接字符串
- `PORT`: 服务端口号
- `ENV`: 运行环境 (development/production)

## API 文档

### 代码运行 API
- POST `/api/run` - 运行代码
- POST `/api/format` - 格式化代码

### 分享 API
- POST `/api/share` - 创建分享
- GET `/api/share/:id` - 获取分享
- POST `/api/share/:id/view` - 增加分享查看次数

## 贡献指南

1. Fork 项目
2. 创建特性分支 (`git checkout -b feature/amazing-feature`)
3. 提交改动 (`git commit -m 'Add amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 创建 Pull Request

## 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情

## 联系方式

如有问题或建议，请提交 Issue 或 Pull Request。 
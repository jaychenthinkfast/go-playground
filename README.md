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
- MongoDB 数据库（持久化存储）
  - 存储分享的代码内容
  - 存储分享的元数据（标题、描述、创建时间等）
  - 存储分享的统计信息（访问次数）
- Redis 缓存
  - 缓存热门分享的代码内容，提高访问速度
  - 使用 TTL 机制自动清理过期的分享
  - 限制分享的访问频率（Rate Limiting）
  - 临时存储运行结果

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
│   ├── src/                # 源代码
│   ├── public/             # 静态资源
│   └── nginx.conf          # Nginx 配置
├── backend/                 # Go 后端代码（多版本）
│   ├── cmd/                # 入口程序
│   ├── pkg/                # 包目录
│   │   ├── api/           # API 处理
│   │   ├── compiler/      # 代码编译和执行
│   │   └── models/        # 数据模型
├── share-service/           # 分享服务
│   ├── cmd/                # 入口程序
│   │   └── server/        # 服务器入口
│   ├── pkg/                # 包目录
│   │   ├── api/           # API 处理
│   │   ├── models/        # 数据模型
│   │   └── storage/       # 存储实现
│   │       ├── mongo/     # MongoDB 实现
│   │       └── redis/     # Redis 缓存实现
├── docker/                  # Docker 配置文件
│   ├── mongo/              # MongoDB 配置
│   └── redis/              # Redis 配置
└── docker-compose.yml       # Docker Compose 配置
```

## 分享服务架构

Share Service 是一个独立的微服务，负责代码分享功能。它使用了双存储架构：

1. **MongoDB**：作为主要存储，永久保存所有分享数据
   - 存储分享的完整信息，包括代码内容、元数据和统计信息
   - 提供数据持久化，确保分享不会丢失

2. **Redis**：作为缓存层，提高访问性能
   - 缓存热门分享，减轻 MongoDB 负担
   - 实现访问频率控制，防止 API 滥用
   - 跟踪分享访问统计，优化热点数据处理

### 主要组件

1. **Handler（API 处理器）**
   - 处理 HTTP 请求，包括创建、获取分享和增加访问计数
   - 管理代码执行请求和结果获取
   - 实现频率限制检查

2. **Storage（存储接口）**
   - 定义统一的存储层接口
   - MongoDB 实现提供持久化存储
   - Redis 实现提供高速缓存

3. **Models（数据模型）**
   - Share：分享模型，包含代码内容、元数据和访问统计
   - RunResult：代码执行结果模型，包含输出、错误信息和性能指标

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

- `MONGO_URI`: MongoDB 连接字符串，用于持久化存储
- `REDIS_URI`: Redis 连接字符串，格式为 `redis://redis:6379/0`
  - redis: Redis 服务器地址
  - 6379: Redis 默认端口
  - 0: 数据库编号
- `PORT`: 服务端口号
- `ENV`: 运行环境 (development/production)

## Redis 配置说明

Redis 在项目中主要用于以下场景：

1. **代码分享缓存**
   - 键格式：`share:{shareId}`
   - 值：序列化的分享对象（JSON）
   - TTL：与分享过期时间一致，无过期时间则默认 7 天
   - 实现：当请求分享时，先检查缓存，缓存未命中再查询 MongoDB

2. **访问频率限制**
   - 键格式：`rate:{ip}:{endpoint}`
   - 值：访问计数器（使用 INCR 命令原子递增）
   - TTL：60秒滑动窗口
   - 实现：使用 Pipeline 执行 INCR 和 EXPIRE 命令，保证原子性
   - 限制：每个 IP 每分钟对同一端点最多 60 次请求

3. **运行结果缓存**
   - 键格式：`result:{taskId}`
   - 值：序列化的运行结果对象（JSON）
   - TTL：5分钟
   - 实现：异步任务完成后，结果存入 Redis，客户端通过 taskId 查询

4. **分享访问统计**
   - 键格式：`views:{shareId}`
   - 值：访问计数（使用 INCR 命令原子递增）
   - 无 TTL，长期保存
   - 实现：每次访问分享时增加计数，定期同步到 MongoDB

### Redis 实现细节

```go
// 缓存实现 - 主要方法
// 存储分享到缓存
func (c *RedisCache) SetShare(ctx context.Context, share *models.Share) error {
    // 序列化分享对象
    // 设置 TTL 与分享过期时间一致
    // 保存到 Redis
}

// 从缓存获取分享
func (c *RedisCache) GetShare(ctx context.Context, shareId string) (*models.Share, error) {
    // 从 Redis 获取数据
    // 反序列化为分享对象
    // 处理缓存未命中情况
}

// 增加分享访问次数
func (c *RedisCache) IncrementViews(ctx context.Context, shareId string) error {
    // 原子递增访问计数
}

// 检查访问频率限制
func (c *RedisCache) IsRateLimited(ctx context.Context, ip, endpoint string) (bool, error) {
    // 使用 Pipeline 保证原子性
    // 递增计数并设置 TTL
    // 返回是否超过限制
}

// 存储代码运行结果
func (c *RedisCache) SetRunResult(ctx context.Context, taskId string, result *models.RunResult) error {
    // 序列化结果对象
    // 设置 5 分钟 TTL
    // 保存到 Redis
}
```

## API 文档

### 代码运行 API
- POST `/api/run` - 运行代码
- POST `/api/format` - 格式化代码

### 分享 API
- POST `/api/share` - 创建分享
- GET `/api/share/:id` - 获取分享
- POST `/api/share/:id/view` - 增加分享查看次数
- POST `/api/execute` - 执行代码（带频率限制）
- GET `/api/result/:taskId` - 获取代码执行结果

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
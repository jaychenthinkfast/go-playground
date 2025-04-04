# Go Playground

一个现代化的在线 Go 代码运行和分享平台。支持多个 Go 版本（1.22、1.23、1.24），提供代码运行、格式化和分享功能。

## 功能特性

- 🚀 支持多个 Go 版本（1.22、1.23、1.24）
- 💻 在线编辑和运行 Go 代码
- 🎨 代码格式化
- 📤 代码分享功能（带有过期选项）
- 🔄 实时运行结果
- 📱 响应式设计，支持移动端访问
- 📊 分享浏览计数统计

## 技术栈

### 前端
- Vue.js 3 + JavaScript
- Vue CLI 构建工具
- Monaco Editor 代码编辑器
- Axios HTTP 客户端

### 后端
- Go (支持 1.22、1.23、1.24 版本)
- Gin Web 框架
- MongoDB 数据库（持久化存储）
  - 存储分享的代码内容
  - 存储分享的元数据（标题、描述、创建时间等）
  - 存储分享的统计信息（访问次数）

### 基础设施
- Docker + Docker Compose
- Nginx 反向代理
- 容器化部署
- 健康检查机制
- 资源限制和重启策略

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
docker compose up --build --remove-orphans
```

3. 访问应用
打开浏览器访问 http://localhost:3003

> 注意：本项目的Docker配置已适配国内网络环境，使用了阿里云镜像源以加速构建和依赖安装。

### 开发模式启动

如果您想在开发模式下启动服务，可以使用：

```bash
docker compose -f docker-compose.dev.yml up --build
```

这将启动开发模式，提供热重载和更详细的日志输出。开发模式下，前端服务会在端口 3003 上运行Vue开发服务器，提供实时热更新功能。

> **重要提示**：在开发模式中，服务名称和通信方式与生产环境有所不同。特别是，分享服务在开发环境中的服务名为 `share-service-dev`，而不是 `share-service`。确保 `frontend/vue.config.js` 中的代理配置正确指向各开发服务。

## 项目结构

```
.
├── frontend/                # 前端代码
│   ├── src/                 # 源代码
│   │   ├── components/      # UI组件
│   │   ├── router/          # 路由配置
│   │   ├── views/           # 页面视图
│   │   ├── assets/          # 静态资源
│   │   ├── App.vue          # 应用入口组件
│   │   └── main.js          # 应用入口文件
│   ├── public/              # 静态资源
│   ├── package.json         # 项目依赖
│   ├── vue.config.js        # Vue 配置文件
│   ├── Dockerfile           # 生产环境 Docker 配置
│   ├── Dockerfile.dev       # 开发环境 Docker 配置
│   └── nginx.conf           # Nginx 配置
├── backend/                 # Go 后端代码（多版本）
│   ├── cmd/                 # 命令行入口
│   │   └── server/          # 服务器入口
│   ├── pkg/                 # 包目录
│   │   ├── runner/          # 代码运行服务
│   │   └── sandbox/         # 安全沙箱
│   ├── app/                 # 应用代码
│   ├── .air.toml            # Air 热更新配置
│   ├── Dockerfile           # 生产环境 Docker 配置
│   ├── Dockerfile.dev       # 开发环境 Docker 配置
│   ├── go.mod               # Go 模块文件
│   └── go.sum               # Go 依赖校验
├── share-service/           # 分享服务
│   ├── cmd/                 # 命令行入口
│   │   └── server/          # 服务器入口
│   ├── pkg/                 # 包目录
│   │   ├── api/             # API 处理
│   │   ├── models/          # 数据模型
│   │   ├── storage/         # 存储接口
│   │   │   └── mongo/       # MongoDB 实现
│   │   └── utils/           # 工具函数
│   ├── .air.toml            # Air 热更新配置
│   ├── Dockerfile           # 生产环境 Docker 配置
│   ├── Dockerfile.dev       # 开发环境 Docker 配置
│   ├── Dockerfile.mongo     # MongoDB Docker 配置
│   ├── go.mod               # Go 模块文件
│   └── go.sum               # Go 依赖校验
├── docker/                  # Docker 相关配置文件
├── docker-compose.yml       # 生产环境 Docker Compose 配置
└── docker-compose.dev.yml   # 开发环境 Docker Compose 配置
```

## 容器服务与端口

项目包含多个容器服务，各自负责不同功能并映射到主机的不同端口：

1. **前端服务 (frontend)**: 端口 3003
   - 生产环境：Nginx 服务运行在容器的 80 端口，映射到主机的 3003 端口
   - 开发环境：Vue 开发服务器运行在容器的 3003 端口，映射到主机的 3003 端口

2. **分享服务 (share-service)**: 端口 3002
   - 处理代码分享、获取和浏览统计功能
   - 作为前端和后端服务的协调者

3. **后端服务**:
   - Go 1.24 (backend-go124): 内部端口 3001
   - Go 1.23 (backend-go123): 内部端口 3001
   - Go 1.22 (backend-go122): 内部端口 3001
   - 各版本在独立容器中运行，通过内部网络通信

4. **MongoDB**: 内部端口 27017
   - 数据持久化存储

所有服务通过名为 `playground-network` 的Docker网络互相通信，保证了环境的隔离性和安全性。

## 安全与性能优化

项目采用了多种安全措施和性能优化策略：

1. **容器安全设置**
   - 使用 `no-new-privileges:true` 限制权限提升
   - 只保留必要的容器能力（capabilities）
   - 服务间严格网络隔离

2. **资源管理**
   - 每个服务设置了 CPU 和内存限制
   - 防止单个服务占用过多资源
   - 优化了容器重启策略

3. **健康检查**
   - 所有服务配置了健康检查
   - 自动检测服务状态并在需要时重启
   - 服务依赖关系明确定义

## 分享服务架构

Share Service 是一个独立的微服务，负责代码分享功能：

### 核心存储

1. **MongoDB**：作为主要存储，永久保存所有分享数据
   - 存储分享的完整信息，包括代码内容、元数据和统计信息
   - 提供数据持久化，确保分享不会丢失
   - 支持设置分享过期时间，自动清理过期内容

### 主要组件

1. **Handler（API 处理器）**
   - 处理 HTTP 请求，包括创建、获取分享和增加访问计数
   - 管理代码执行请求和结果获取
   - 提供健康检查接口，便于服务监控

2. **Storage（存储接口）**
   - 定义统一的存储层接口
   - MongoDB 实现提供持久化存储
   - 支持原子性操作，确保数据一致性

3. **Models（数据模型）**
   - Share：分享模型，包含代码内容、元数据和访问统计
   - RunResult：代码执行结果模型，包含输出、错误信息和性能指标
   - API请求/响应模型：标准化接口数据交换格式

### 执行流程

1. **代码分享**：
   - 用户编写代码并设置分享选项（标题、描述、过期时间等）
   - 前端调用 `/api/share` 接口创建分享
   - 服务端生成唯一ID，保存分享内容，并返回分享链接

2. **访问分享**：
   - 通过分享链接访问代码分享页面
   - 前端调用 `/api/share/:id` 获取分享内容
   - 服务端验证分享是否存在和过期状态，返回内容并增加访问计数

3. **代码执行**：
   - 用户点击运行按钮执行分享的代码
   - 前端调用 `/api/execute` 接口发送代码和版本信息
   - 服务端将请求转发到对应版本的后端服务执行，并返回结果

## 开发指南

### 本地开发环境设置

1. 前端开发
```bash
cd frontend
npm install
npm run serve
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

### 配置前端与后端服务的连接

前端开发环境通过 Vue CLI 的代理配置与各服务通信，配置位于 `frontend/vue.config.js` 文件中：

```javascript
// frontend/vue.config.js
module.exports = {
  devServer: {
    port: 3003,
    proxy: {
      '/api/share': {
        target: 'http://share-service-dev:3002',  // 开发环境中使用 share-service-dev
        changeOrigin: true,
        ws: true
      },
      '/api/execute': {
        target: 'http://share-service-dev:3002',  // 开发环境中使用 share-service-dev
        changeOrigin: true,
        ws: true
      },
      // ...其他代理配置
    }
  }
}
```

> **注意**：如果分享功能在开发环境中失败，检查上述配置中的服务名是否正确设置为 `share-service-dev`。

### 构建与部署

1. 构建镜像：
```bash
docker compose build
```

2. 生产环境部署：
```bash
docker compose up -d 
```

### 代码风格与提交规范

- 遵循 Go 标准代码风格
- 使用 `gofmt` 格式化代码
- 提交前运行测试确保功能正常
- 提交信息格式：`<类型>: <描述>`，例如 `feat: 添加用户验证功能`

## 浏览计数机制

分享服务对每个分享页面的浏览次数进行计数，采用基于 Cookie 的去重机制：

1. **访问追踪**：
   - 当用户首次访问分享页面时，服务器会增加浏览计数
   - 同时在用户浏览器中设置名为 `viewed_{shareId}` 的 Cookie，有效期 24 小时
   - 用户在 Cookie 有效期内重复访问同一分享页面，浏览计数不会增加

2. **计数实现**：
   - `GetShare` API 处理函数在返回分享内容的同时负责浏览计数
   - 使用 MongoDB 的 `FindOneAndUpdate` 操作原子性地增加计数并返回更新后的值
   - 保证计数准确性，避免并发访问时的竞态条件

3. **手动增加计数**：
   - 提供独立的 `/api/share/:id/view` API 端点，用于手动增加计数
   - 此端点不受 Cookie 去重机制限制，可用于特殊场景

## 代码执行流程

1. **请求处理**：
   - 前端发送包含代码和Go版本的请求到 `/api/execute` 端点
   - 服务根据版本将请求路由到对应的后端服务

2. **版本支持**：
   - 支持多种Go版本格式：`go1.24`、`1.24`、`go1.24.0`等
   - 自动标准化版本格式，确保正确路由

3. **安全限制**：
   - 代码执行在隔离容器中进行
   - 限制执行时间和资源使用
   - 防止恶意代码执行

## API 文档

### 健康检查
- GET `/health` - 服务健康状态检查

### 代码运行 API
- POST `/api/run` - 运行代码
  ```json
  {
    "code": "package main\n\nimport \"fmt\"\n\nfunc main() {\n\tfmt.Println(\"Hello, World!\")\n}",
    "version": "go1.24"
  }
  ```

- POST `/api/format` - 格式化代码
  ```json
  {
    "code": "package main\nimport \"fmt\"\nfunc main(){\nfmt.Println(\"Hello\")\n}"
  }
  ```

### 分享 API
- POST `/api/share` - 创建分享
  ```json
  {
    "code": "package main\n\nimport \"fmt\"\n\nfunc main() {\n\tfmt.Println(\"Hello, World!\")\n}",
    "version": "go1.24",
    "title": "Hello World Example",
    "description": "A simple hello world program",
    "author": "GoExpert",
    "expires_in": "48h"
  }
  ```

- GET `/api/share/:id` - 获取分享（自动处理浏览计数）
- POST `/api/share/:id/view` - 手动增加分享查看次数
- POST `/api/execute` - 执行代码
  ```json
  {
    "code": "package main\n\nimport \"fmt\"\n\nfunc main() {\n\tfmt.Println(\"Hello, World!\")\n}",
    "version": "go1.24"
  }
  ```

## 进阶功能

### 分享过期设置
创建分享时可指定 `expires_in` 参数，采用 Go 的时间格式（如 `24h`、`7d`），到期后分享将无法访问。

### 自定义分享元数据
支持为分享添加标题、描述和作者信息，使分享内容更加丰富和易于理解。

### 多版本支持
同一份代码可以在不同的 Go 版本中运行，便于测试新特性或检查兼容性问题。

## 常见问题

1. **浏览计数异常**
   问：为什么同一份分享刷新页面后计数不增加？
   答：系统使用 Cookie 机制确保 24 小时内同一浏览器访问同一分享只计数一次，这是为了避免重复计数。

2. **代码执行失败**
   问：代码执行返回错误如何排查？
   答：检查代码是否有编译错误，以及是否使用了不支持的特性或包。API 响应会包含详细的错误信息。

3. **分享功能在开发环境中失败**
   问：为什么在开发环境中进行代码分享时提示"分享失败，请稍后重试"？
   答：检查 `frontend/vue.config.js` 文件中的代理配置是否正确指向 `share-service-dev`，而不是 `share-service`。开发环境与生产环境的服务名称不同。

4. **Air 热更新不生效**
   问：修改代码后 Air 没有自动重新构建和重启服务怎么办？
   答：确保修改的文件扩展名包含在 `.air.toml` 的 `include_ext` 配置中，并且该文件没有被 `exclude_dir` 或 `exclude_file` 配置排除。也可以尝试重启 Air 服务。

5. **容器启动失败**
   问：Docker 容器启动时出现健康检查失败的错误怎么处理？
   答：检查容器日志 (`docker logs <container-name>`) 分析具体原因。常见问题包括端口冲突、依赖服务未就绪、配置错误等。增加健康检查的超时时间和重试次数也有助于解决服务启动较慢的问题。

## 开发环境热更新

本项目在开发环境中使用 Air 工具实现热更新功能，使开发过程更加高效：

1. **前端热更新**：
   - Vue CLI 开发服务器自动监控前端代码变化
   - 修改 Vue 组件后立即在浏览器中反映变化，无需手动刷新

2. **后端热更新**：
   - Air 工具监控 Go 代码变化
   - 自动重新编译和重启服务
   - 配置位于各服务目录下的 `.air.toml` 文件中

3. **如何使用**：
   - 启动开发环境: `docker compose -f docker-compose.dev.yml up`
   - 修改代码后，观察容器日志以确认热更新是否成功
   - 如需调试热更新问题: `docker logs -f <container-name>`

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
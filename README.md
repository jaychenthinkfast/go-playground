# Go Playground

![Go Logo](https://golang.org/lib/godoc/images/go-logo-blue.svg)

Go Playground是一个在线服务，允许用户编写、运行Go代码而无需本地安装Go环境。系统在服务器端编译和执行代码，并将结果返回给用户。

## 1. 系统概述

Go Playground 提供以下功能：

- 在线编写和执行 Go 代码
- 支持多个 Go 版本（1.24、1.23 和开发分支）
- 代码格式化
- 代码分享功能
- 提供常用示例代码
- 安全的沙盒环境执行

该项目基于 Go 官方 Playground (https://go.dev/play/) 的简化实现，旨在提供类似的功能和用户体验。

## 2. 系统架构

### 2.1 架构图

```
+-------------------+        +-------------------+
|                   |        |                   |
|  Frontend (Vue.js)|  <---> |  Backend (Go)     |
|                   |  HTTP  |                   |
+-------------------+        +-------------------+
                                     |
                              +------v------+
                              |             |
                              | Go Sandbox  |
                              |             |
                              +-------------+
```

### 2.2 前端组件
- 代码编辑器 - 基于原生文本域实现的编辑器
- 版本选择器 (Go 1.24/1.23/dev分支)
- 运行/格式化/分享按钮
- 示例代码选择 - 提供多种代码示例

### 2.3 后端服务
- 代码接收器 - 接收前端提交的代码
- 代码验证与编译服务 - 验证代码安全性并编译
- 沙盒执行环境 - 安全隔离的代码执行环境
- 结果返回服务 - 返回执行结果给前端

## 3. 系统安全

Go Playground 实现了多层安全机制：

- 沙盒隔离执行环境
- 资源限制
  - 内存使用上限（50MB）
  - CPU时间限制（5秒）
  - 执行时间限制
- 网络访问限制
  - 禁止网络连接（net.Dial, net.Listen）
  - 限制文件系统操作（os.Open, os.Create, os.Remove）
  - 限制系统调用（syscall）
- 代码静态分析，防止危险操作

## 4. 功能特性

- **代码编译与执行**：在安全的沙盒环境中执行用户提交的Go代码
- **格式化功能**：使用 gofmt 自动格式化代码
- **结果实时展示**：实时显示代码执行结果
- **代码分享功能**：生成可分享的链接
- **预设示例代码**：提供多个常用示例（Hello World、Conway's Game of Life、Fibonacci等）
- **多版本Go支持**：通过显示版本信息模拟不同Go版本支持

### 4.1 版本切换功能详解

Go Playground支持在不同Go版本间切换，实现方式如下：

- 前端通过下拉菜单选择目标Go版本
- 后端根据选择的版本提供相应的版本信息
- 系统显示当前使用的Go版本信息，方便用户验证

这使得用户可以:
- 了解不同Go版本的特点
- 查看代码在不同版本环境下的运行情况
- 学习Go语言版本变化

## 5. 限制条件

为了确保系统安全和资源合理使用，Go Playground 有以下限制：

- 只能使用部分标准库
- 无外部网络访问
- 时间固定为2009-11-10 23:00:00 UTC（为确保输出确定性）
- 执行时间、CPU和内存使用量有限制
- 禁止执行危险系统调用

## 6. 技术栈选型

- **前端**: Vue.js 3 - 轻量级渐进式JavaScript框架
- **后端**: Go语言 - 高性能网络服务开发语言
- **容器化**: Docker - 应用容器引擎
- **代理服务器**: Nginx - 高性能Web服务器
- **API通信**: RESTful API - 标准HTTP接口
- **开发工具**: Webpack, NPM, Go Modules

## 7. 项目结构

```
go-playground/
│
├── frontend/                # Vue.js前端项目
│   ├── public/              # 静态资源
│   │   └── index.html       # 主HTML页面
│   ├── src/                 # 源代码
│   │   ├── components/      # Vue组件
│   │   │   └── Playground.vue # 主要组件
│   │   ├── router/          # 路由配置
│   │   ├── App.vue          # 主应用组件
│   │   └── main.js          # 入口文件
│   ├── Dockerfile           # 前端Docker配置
│   ├── nginx.conf           # Nginx配置
│   ├── package.json         # 依赖管理
│   └── vue.config.js        # Vue配置
│
├── backend/                 # Go后端项目
│   ├── cmd/                 # 命令行入口
│   │   └── server/          # HTTP服务器
│   │       └── main.go      # 主程序入口
│   ├── pkg/                 # 包
│   │   ├── runner/          # 代码运行器
│   │   │   └── runner.go    # 运行器实现
│   │   └── sandbox/         # 沙盒实现
│   │       └── sandbox.go   # 沙盒隔离环境
│   ├── Dockerfile           # 后端Docker配置
│   └── go.mod               # Go模块定义
│
├── docker-compose.yml       # Docker Compose配置
└── README.md                # 项目文档
```

## 8. 如何运行

### 8.1 环境要求
- Node.js 16+ (前端开发)
- Go 1.18+ (后端开发)
- Docker & Docker Compose (容器化部署)

### 8.2 开发环境

#### 前端开发
```bash
# 进入前端目录
cd frontend

# 安装依赖
npm install

# 启动开发服务器
npm run serve
```

开发服务器将在 http://localhost:3003 上运行。

#### 后端开发
```bash
# 进入后端目录
cd backend

# 获取依赖
go mod download

# 运行服务器
go run cmd/server/main.go
```

后端服务器将在 http://localhost:3001 上运行。

### 8.3 使用Docker Compose启动完整应用
```bash
# 在项目根目录执行
docker-compose up -d

# 查看运行状态
docker-compose ps

# 查看日志
docker-compose logs -f
```

启动后可以通过 http://localhost:3003 访问应用。

### 8.4 快速使用指南

1. 访问Web界面 http://localhost:3003
2. 从下拉菜单选择Go版本 (1.24/1.23/dev)
3. 在编辑器中编写Go代码或选择预设的代码示例
4. 点击"Run"按钮执行代码
5. 查看输出结果，包括版本信息和代码执行结果
6. 可选：点击"Format"格式化代码
7. 可选：点击"Share"分享代码（功能开发中）

## 9. API文档

### 9.1 健康检查
- **URL**: `/api/health`
- **方法**: `GET`
- **响应**: `200 OK` 表示服务正常运行

### 9.2 运行代码
- **URL**: `/api/run`
- **方法**: `POST`
- **请求体**:
  ```json
  {
    "code": "package main\n\nimport \"fmt\"\n\nfunc main() {\n\tfmt.Println(\"Hello, 世界\")\n}",
    "version": "go1.24"
  }
  ```
- **响应**:
  ```json
  {
    "output": "Go 1.24 - Released February 2024\nRequested Go version: go1.24\nActual Go version: go version go1.24.1 darwin/amd64\n\nHello, 世界\n",
    "error": ""
  }
  ```

### 9.3 格式化代码
- **URL**: `/api/format`
- **方法**: `POST`
- **请求体**:
  ```json
  {
    "code": "package main\n\nimport \"fmt\"\n\nfunc main(){\nfmt.Println(\"Hello, 世界\")\n}"
  }
  ```
- **响应**:
  ```json
  {
    "formattedCode": "package main\n\nimport \"fmt\"\n\nfunc main() {\n\tfmt.Println(\"Hello, 世界\")\n}\n",
    "error": ""
  }
  ```

## 10. Docker部署详解

本项目使用Docker进行容器化部署，主要包括以下几个方面：

### 10.1 前端容器
- 基于Node.js Alpine构建Vue应用
- 使用Nginx提供静态文件服务
- 配置反向代理将API请求转发到后端
- 提供静态资源缓存和gzip压缩

### 10.2 后端容器
- 基于Golang Alpine构建Go应用
- 多阶段构建减小镜像体积
- 配置必要的安全限制
- 使用非root用户运行服务

### 10.3 容器网络
- 使用Docker网络实现前后端通信
- 暴露必要的端口供外部访问
- 实现服务隔离

### 10.4 持久化和配置
- 环境变量配置容器行为
- 支持临时文件目录共享
- 健康检查确保服务可用性

## 11. 版本实现说明

在当前实现中，版本切换并不是通过实际运行不同版本的Go编译器实现的，而是通过以下方式实现的：

1. 前端允许用户选择不同的Go版本（1.24、1.23、dev）
2. 后端记录用户请求的版本，并在响应中包含:
   - 版本的描述信息（如发布日期）
   - 用户请求的版本信息
   - 实际执行代码的Go版本

这种方式更适合教学和演示目的，让用户了解不同版本的特点。在生产环境中，完整实现应该使用不同版本的Go编译器或容器。

## 12. 故障排除

### 前端无法连接后端
- 检查后端服务是否正常运行（http://localhost:3001/api/health）
- 确认端口配置正确（默认后端3001，前端3003）
- 检查CORS设置是否正确

### 代码运行失败
- 检查代码语法是否正确
- 确认没有使用受限的包或API
- 检查日志获取更多错误信息

### Docker相关问题
- 确保Docker和Docker Compose已正确安装
- 检查端口冲突
- 查看容器日志排查问题

## 13. 贡献指南

欢迎提交Pull Request或Issue来改进项目。贡献步骤：

1. Fork本项目
2. 创建您的特性分支 (`git checkout -b feature/amazing-feature`)
3. 提交您的更改 (`git commit -m 'Add some amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 创建一个Pull Request

### 代码规范
- 前端代码遵循Vue.js风格指南
- 后端代码遵循Go官方代码规范
- 所有API需要有完整的文档
- 主要功能需要有测试覆盖

## 14. 许可证

本项目采用MIT许可证 - 详见 [LICENSE](LICENSE) 文件。

## 15. 致谢

- [Go 团队](https://golang.org/project/)提供的原始 Playground 实现
- [Vue.js](https://vuejs.org/) 团队提供的优秀前端框架
- 所有为本项目做出贡献的开发者 
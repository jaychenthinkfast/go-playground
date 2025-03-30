.PHONY: help run-frontend run-backend run-all docker-build docker-run docker-stop clean test test-frontend test-backend lint dev-backend

# 默认目标: 显示帮助信息
help:
	@echo "Go Playground - Makefile命令列表"
	@echo "----------------------------------------"
	@echo "make run-frontend  - 启动前端开发服务器"
	@echo "make run-backend   - 启动后端服务器"
	@echo "make run-all       - 同时启动前端和后端"
	@echo "make dev-backend   - 使用air启动后端（热重载）"
	@echo "make docker-build  - 构建Docker镜像"
	@echo "make docker-run    - 使用Docker Compose启动所有服务"
	@echo "make docker-stop   - 停止Docker Compose服务"
	@echo "make clean         - 清理构建文件"
	@echo "make test          - 运行所有测试"
	@echo "make test-frontend - 运行前端测试"
	@echo "make test-backend  - 运行后端测试"
	@echo "make lint          - 运行代码风格检查"
	@echo "----------------------------------------"

# 前端开发
run-frontend:
	@echo "启动前端开发服务器，在http://localhost:3003..."
	@cd frontend && npm install && npm run serve

# 后端开发
run-backend:
	@echo "启动后端服务器，在http://localhost:3001..."
	@cd backend && go mod download && go run cmd/server/main.go

# 使用air热重载后端
dev-backend:
	@echo "使用air启动后端（热重载），在http://localhost:3001..."
	@if ! command -v air > /dev/null; then \
		echo "安装air..."; \
		go install github.com/cosmtrek/air@latest; \
	fi
	@cd backend && air

# 同时启动前端和后端（使用子进程）
run-all:
	@echo "同时启动前端和后端服务..."
	@make run-backend & \
	make run-frontend

# Docker相关操作
docker-build:
	@echo "构建Docker镜像..."
	@docker-compose build

docker-run:
	@echo "启动Docker Compose服务..."
	@docker-compose up -d
	@echo "服务已启动，访问http://localhost:3003"

docker-stop:
	@echo "停止Docker Compose服务..."
	@docker-compose down
	@echo "服务已停止"

# 清理构建文件
clean:
	@echo "清理构建文件..."
	@rm -rf frontend/dist
	@rm -rf frontend/node_modules
	@echo "清理完成"

# 测试相关操作
test: test-frontend test-backend
	@echo "所有测试完成"

test-frontend:
	@echo "运行前端测试..."
	@cd frontend && npm test

test-backend:
	@echo "运行后端测试..."
	@cd backend && go test -v ./...

# 代码质量检查
lint:
	@echo "运行代码风格检查..."
	@cd frontend && npm run lint
	@cd backend && go vet ./...
	@echo "代码检查完成"

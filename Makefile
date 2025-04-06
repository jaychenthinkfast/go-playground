.PHONY: build build-dev up up-dev down down-dev restart restart-dev logs logs-dev ps clean help

# 默认目标
.DEFAULT_GOAL := help

# 生产环境目标
build: ## 构建生产环境镜像
	docker compose build

up: ## 启动生产环境服务（后台运行）
	docker compose up --build -d

down: ## 停止生产环境服务
	docker compose down

restart: down up ## 重启生产环境服务

logs: ## 查看生产环境服务日志
	docker compose logs -f

# 开发环境目标
build-dev: ## 构建开发环境镜像
	docker compose -f docker-compose.dev.yml build

up-dev: ## 启动开发环境服务（后台运行）
	docker compose -f docker-compose.dev.yml up --build -d

down-dev: ## 停止开发环境服务
	docker compose -f docker-compose.dev.yml down

restart-dev: down-dev up-dev ## 重启开发环境服务

logs-dev: ## 查看开发环境服务日志
	docker compose -f docker-compose.dev.yml logs -f

# 通用命令
ps: ## 显示当前运行的容器
	docker compose ps

clean: ## 清理Docker缓存和未使用的卷、网络
	docker system prune -f
	docker volume prune -f
	docker network prune -f

help: ## 显示帮助信息
	@echo "Go Playground - Makefile帮助"
	@echo ""
	@echo "使用方法:"
	@echo "  make [target]"
	@echo ""
	@echo "可用命令:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}' 
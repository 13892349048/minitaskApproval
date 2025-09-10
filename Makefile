# ================================================
# TaskFlow Makefile
# ================================================

.PHONY: help build test clean docker-build docker-run dev prod migrate lint fmt

# 默认目标
.DEFAULT_GOAL := help

# 变量定义
APP_NAME := taskflow
VERSION := $(shell git describe --tags --always --dirty)
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
GO_VERSION := $(shell go version | awk '{print $$3}')

# 构建标志
LDFLAGS := -ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.GoVersion=$(GO_VERSION)"

## help: 显示帮助信息
help:
	@echo "TaskFlow 项目管理系统"
	@echo ""
	@echo "可用命令:"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

## build: 构建应用程序
build:
	@echo "🔨 构建应用程序..."
	@CGO_ENABLED=0 go build $(LDFLAGS) -o bin/api ./cmd/api
	@CGO_ENABLED=0 go build $(LDFLAGS) -o bin/migrate ./cmd/migrate
	@echo "✅ 构建完成"

## test: 运行测试
test:
	@echo "🧪 运行测试..."
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "✅ 测试完成，覆盖率报告: coverage.html"

## test-coverage: 检查测试覆盖率
test-coverage: test
	@COVERAGE=$$(go tool cover -func=coverage.out | grep total | awk '{print substr($$3, 1, length($$3)-1)}'); \
	echo "📊 测试覆盖率: $$COVERAGE%"; \
	if [ $$(echo "$$COVERAGE < 70" | bc -l) -eq 1 ]; then \
		echo "❌ 测试覆盖率低于70%"; \
		exit 1; \
	else \
		echo "✅ 测试覆盖率达标"; \
	fi

## lint: 代码检查
lint:
	@echo "🔍 运行代码检查..."
	@go vet ./...
	@golint ./...
	@staticcheck ./...
	@echo "✅ 代码检查完成"

## fmt: 格式化代码
fmt:
	@echo "🎨 格式化代码..."
	@go fmt ./...
	@goimports -w .
	@echo "✅ 代码格式化完成"

## clean: 清理构建文件
clean:
	@echo "🧹 清理构建文件..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html
	@docker system prune -f
	@echo "✅ 清理完成"

## dev: 启动开发环境
dev:
	@echo "🚀 启动开发环境..."
	@docker-compose -f docker-compose.dev.yml up --build

## dev-down: 停止开发环境
dev-down:
	@echo "⏹️  停止开发环境..."
	@docker-compose -f docker-compose.dev.yml down -v

## prod: 启动生产环境
prod:
	@echo "🚀 启动生产环境..."
	@docker-compose -f docker-compose.prod.yml up -d --build

## prod-down: 停止生产环境
prod-down:
	@echo "⏹️  停止生产环境..."
	@docker-compose -f docker-compose.prod.yml down

## migrate: 运行数据库迁移
migrate:
	@echo "🗃️  运行数据库迁移..."
	@./scripts/migrate.sh migrate

## migrate-reset: 重置数据库
migrate-reset:
	@echo "⚠️  重置数据库..."
	@./scripts/migrate.sh reset

## migrate-status: 检查迁移状态
migrate-status:
	@echo "📊 检查迁移状态..."
	@./scripts/migrate.sh status

## validate-models: 验证GORM模型
validate-models:
	@echo "🔍 验证GORM模型..."
	@go run cmd/migrate/main.go -cmd=validate

## docker-build: 构建Docker镜像
docker-build:
	@echo "🐳 构建Docker镜像..."
	@docker build -t $(APP_NAME):$(VERSION) .
	@docker build -t $(APP_NAME):latest .
	@echo "✅ Docker镜像构建完成"

## docker-run: 运行Docker容器
docker-run: docker-build
	@echo "🐳 运行Docker容器..."
	@docker run -p 8080:8080 --name $(APP_NAME) $(APP_NAME):latest

## install-tools: 安装开发工具
install-tools:
	@echo "🛠️  安装开发工具..."
	@go install golang.org/x/lint/golint@latest
	@go install honnef.co/go/tools/cmd/staticcheck@latest
	@go install golang.org/x/tools/cmd/goimports@latest
	@go install github.com/air-verse/air@latest
	@echo "✅ 开发工具安装完成"

## logs: 查看应用日志
logs:
	@echo "📋 查看应用日志..."
	@docker-compose -f docker-compose.dev.yml logs -f app

## logs-prod: 查看生产环境日志
logs-prod:
	@echo "📋 查看生产环境日志..."
	@docker-compose -f docker-compose.prod.yml logs -f app

## health: 检查应用健康状态
health:
	@echo "💓 检查应用健康状态..."
	@curl -f http://localhost:8080/health || echo "❌ 应用不健康"

## version: 显示版本信息
version:
	@echo "📦 版本信息:"
	@echo "  应用版本: $(VERSION)"
	@echo "  构建时间: $(BUILD_TIME)"
	@echo "  Go版本: $(GO_VERSION)"

## security-scan: 安全扫描
security-scan:
	@echo "🔒 运行安全扫描..."
	@gosec ./...
	@echo "✅ 安全扫描完成"

## benchmark: 性能基准测试
benchmark:
	@echo "⚡ 运行性能基准测试..."
	@go test -bench=. -benchmem ./...
	@echo "✅ 基准测试完成"

## all: 完整构建流程
all: clean fmt lint test build docker-build
	@echo "🎉 完整构建流程完成!"

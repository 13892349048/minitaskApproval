# TaskFlow - 多层级项目管理系统

[![CI/CD Pipeline](https://github.com/your-username/taskflow/actions/workflows/ci.yml/badge.svg)](https://github.com/your-username/taskflow/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/your-username/taskflow)](https://goreportcard.com/report/github.com/your-username/taskflow)
[![Coverage Status](https://coveralls.io/repos/github/your-username/taskflow/badge.svg?branch=main)](https://coveralls.io/github/your-username/taskflow?branch=main)

> 基于DDD架构的企业级多层级项目管理系统，支持复杂的任务流程、权限控制和团队协作。

## 🚀 功能特性

### 核心功能
- 📋 **多层级项目管理** - 支持主项目、子项目、临时项目的层级结构
- 🎯 **灵活任务系统** - 单次执行任务和重复任务，支持复杂的审批流程
- 👥 **多人协作** - 任务参与者、负责人、审批人的角色分工
- ⏰ **智能调度** - 重复任务自动调度和延期申请管理
- 📊 **统计分析** - 项目进度、用户工作负载、任务完成率统计

### 技术特性
- 🏗️ **DDD架构** - 领域驱动设计，清晰的业务边界
- 🔐 **RBAC+ABAC权限** - 基于角色和属性的双重权限控制
- 📁 **文件管理** - 分片上传、断点续传、文件关联
- 🔍 **全文搜索** - MySQL全文索引，快速搜索任务和项目
- 🎭 **状态机** - 任务和执行状态的严格流转控制
- 📨 **事件驱动** - 领域事件和消息传递机制

## 🛠️ 技术栈

### 后端技术
- **语言**: Go 1.21+
- **框架**: Gin 1.9+
- **数据库**: MySQL 8.0+, Redis 7.0+
- **ORM**: GORM 1.25+
- **认证**: JWT-Go 4.5+
- **配置**: Viper 1.16+
- **日志**: Zap 1.24+ + Lumberjack
- **测试**: Testify 1.8+

### 基础设施
- **容器化**: Docker, Docker Compose
- **CI/CD**: GitHub Actions
- **代理**: Nginx
- **监控**: Prometheus, Grafana (计划中)
- **链路追踪**: Jaeger (计划中)

## 📦 快速开始

### 环境要求
- Go 1.21+
- Docker & Docker Compose
- MySQL 8.0+
- Redis 7.0+

### 1. 克隆项目
```bash
git clone https://github.com/your-username/taskflow.git
cd taskflow
```

### 2. 配置环境
```bash
# 复制环境配置文件
cp env.example .env

# 编辑环境配置
vim .env
```

### 3. 启动开发环境
```bash
# 使用Make命令（推荐）
make dev

# 或使用Docker Compose
docker-compose -f docker-compose.dev.yml up --build
```

### 4. 运行数据库迁移
```bash
# 使用Make命令
make migrate

# 或使用脚本
./scripts/migrate.sh migrate
```

### 5. 验证部署
```bash
# 健康检查
curl http://localhost:8080/health

# 版本信息
curl http://localhost:8080/version
```

## 🔧 开发指南

### 项目结构
```
taskflow/
├── cmd/                    # 应用程序入口
│   ├── api/               # Web API服务
│   └── migrate/           # 数据库迁移工具
├── internal/              # 私有应用代码
│   ├── domain/           # 领域层
│   ├── application/      # 应用服务层
│   ├── infrastructure/   # 基础设施层
│   └── interfaces/       # 接口层
├── pkg/                   # 公共包
├── configs/              # 配置文件
├── scripts/              # 脚本文件
├── .github/workflows/    # CI/CD工作流
└── docs/                 # 文档
```

### 常用命令
```bash
# 开发相关
make help          # 查看所有可用命令
make build         # 构建应用程序
make test          # 运行测试
make test-coverage # 检查测试覆盖率
make lint          # 代码检查
make fmt           # 格式化代码

# Docker相关
make docker-build  # 构建Docker镜像
make docker-run    # 运行Docker容器
make dev           # 启动开发环境
make prod          # 启动生产环境

# 数据库相关
make migrate       # 运行数据库迁移
make migrate-reset # 重置数据库
make validate-models # 验证GORM模型

# 部署相关
./scripts/deploy.sh dev --build --migrate    # 部署到开发环境
./scripts/deploy.sh prod --build --migrate   # 部署到生产环境
```

### 开发工具安装
```bash
# 安装开发工具
make install-tools

# 或手动安装
go install golang.org/x/lint/golint@latest
go install honnef.co/go/tools/cmd/staticcheck@latest
go install golang.org/x/tools/cmd/goimports@latest
go install github.com/air-verse/air@latest
```

## 📚 API文档

### 认证接口
- `POST /api/v1/auth/login` - 用户登录
- `POST /api/v1/auth/register` - 用户注册
- `POST /api/v1/auth/refresh` - 刷新Token

### 项目管理
- `GET /api/v1/projects` - 获取项目列表
- `POST /api/v1/projects` - 创建项目
- `GET /api/v1/projects/:id` - 获取项目详情
- `PUT /api/v1/projects/:id` - 更新项目
- `DELETE /api/v1/projects/:id` - 删除项目

### 任务管理
- `GET /api/v1/tasks` - 获取任务列表
- `POST /api/v1/tasks` - 创建任务
- `GET /api/v1/tasks/:id` - 获取任务详情
- `PUT /api/v1/tasks/:id` - 更新任务
- `POST /api/v1/tasks/:id/submit` - 提交任务
- `POST /api/v1/tasks/:id/approve` - 审批任务

详细API文档请参考：[API文档](docs/api.md)

## 🧪 测试

### 运行测试
```bash
# 运行所有测试
make test

# 运行测试并生成覆盖率报告
make test-coverage

# 运行基准测试
make benchmark

# 运行安全扫描
make security-scan
```

### 测试覆盖率要求
- 单元测试覆盖率 >= 70%
- 集成测试覆盖核心业务流程
- E2E测试覆盖主要用户场景

## 🚀 部署

### 开发环境部署
```bash
# 使用部署脚本
./scripts/deploy.sh dev --build --migrate

# 或使用Make命令
make dev
```

### 生产环境部署
```bash
# 1. 配置环境变量
cp env.example .env
vim .env

# 2. 部署服务
./scripts/deploy.sh prod --build --migrate

# 3. 健康检查
./scripts/deploy.sh prod --health
```

### Docker部署
```bash
# 构建镜像
docker build -t taskflow:latest .

# 运行容器
docker run -p 8080:8080 --env-file .env taskflow:latest
```

## 📊 监控和日志

### 健康检查
```bash
# 应用健康检查
curl http://localhost:8080/health

# 数据库连接检查
curl http://localhost:8080/health/db
```

### 日志查看
```bash
# 查看应用日志
make logs

# 查看生产环境日志
make logs-prod

# 使用部署脚本查看日志
./scripts/deploy.sh dev --logs
```

## 🤝 贡献指南

1. Fork 项目
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 创建 Pull Request

### 代码规范
- 遵循 Go 官方代码规范
- 使用 `gofmt` 格式化代码
- 通过 `golint` 和 `go vet` 检查
- 添加适当的单元测试
- 更新相关文档

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情

## 📞 联系我们

- 项目维护者: [Your Name](mailto:your.email@example.com)
- 项目地址: [https://github.com/your-username/taskflow](https://github.com/your-username/taskflow)
- 问题反馈: [Issues](https://github.com/your-username/taskflow/issues)

## 🙏 致谢

感谢所有为这个项目做出贡献的开发者！

---

⭐ 如果这个项目对你有帮助，请给它一个星标！

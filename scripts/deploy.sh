#!/bin/bash

# ================================================
# TaskFlow 部署脚本
# ================================================

set -e

# 配置
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
DEPLOY_ENV="${DEPLOY_ENV:-development}"
DOCKER_REGISTRY="${DOCKER_REGISTRY:-ghcr.io}"
IMAGE_NAME="${IMAGE_NAME:-taskflow}"
VERSION="${VERSION:-latest}"

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 日志函数
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 显示帮助信息
show_help() {
    echo "TaskFlow 部署脚本"
    echo ""
    echo "用法: $0 [环境] [选项]"
    echo ""
    echo "环境:"
    echo "  dev         部署到开发环境"
    echo "  prod        部署到生产环境"
    echo "  staging     部署到预发布环境"
    echo ""
    echo "选项:"
    echo "  --build     重新构建镜像"
    echo "  --migrate   运行数据库迁移"
    echo "  --rollback  回滚到上一个版本"
    echo "  --health    检查服务健康状态"
    echo "  --logs      查看服务日志"
    echo "  --stop      停止服务"
    echo "  --help      显示帮助信息"
    echo ""
    echo "环境变量:"
    echo "  DEPLOY_ENV       部署环境 (development/staging/production)"
    echo "  DOCKER_REGISTRY  Docker镜像仓库"
    echo "  IMAGE_NAME       镜像名称"
    echo "  VERSION          版本标签"
    echo ""
    echo "示例:"
    echo "  $0 dev --build"
    echo "  $0 prod --migrate"
    echo "  DEPLOY_ENV=staging $0 staging --build --migrate"
}

# 检查依赖
check_dependencies() {
    log_info "检查依赖..."
    
    local deps=("docker" "docker-compose" "curl")
    for dep in "${deps[@]}"; do
        if ! command -v "$dep" &> /dev/null; then
            log_error "缺少依赖: $dep"
            exit 1
        fi
    done
    
    log_success "依赖检查通过"
}

# 检查环境文件
check_env_file() {
    local env_file="$PROJECT_ROOT/.env"
    
    if [ ! -f "$env_file" ]; then
        log_warning "环境文件不存在: $env_file"
        log_info "请复制 env.example 为 .env 并配置相应的值"
        
        if [ -f "$PROJECT_ROOT/env.example" ]; then
            cp "$PROJECT_ROOT/env.example" "$env_file"
            log_info "已创建环境文件模板: $env_file"
            log_warning "请编辑 $env_file 并填入实际配置值"
            return 1
        fi
    fi
    
    log_success "环境文件检查通过"
    return 0
}

# 构建镜像
build_image() {
    log_info "构建Docker镜像..."
    
    cd "$PROJECT_ROOT"
    
    # 构建镜像
    docker build -t "${IMAGE_NAME}:${VERSION}" .
    docker build -t "${IMAGE_NAME}:latest" .
    
    log_success "镜像构建完成: ${IMAGE_NAME}:${VERSION}"
}

# 运行数据库迁移
run_migration() {
    log_info "运行数据库迁移..."
    
    local compose_file
    case "$DEPLOY_ENV" in
        "development")
            compose_file="docker-compose.dev.yml"
            ;;
        "production")
            compose_file="docker-compose.prod.yml"
            ;;
        *)
            compose_file="docker-compose.yml"
            ;;
    esac
    
    cd "$PROJECT_ROOT"
    
    # 确保数据库服务正在运行
    docker-compose -f "$compose_file" up -d mysql
    
    # 等待数据库就绪
    log_info "等待数据库就绪..."
    sleep 10
    
    # 运行迁移
    docker-compose -f "$compose_file" run --rm migrator
    
    log_success "数据库迁移完成"
}

# 部署服务
deploy_services() {
    log_info "部署服务到 $DEPLOY_ENV 环境..."
    
    local compose_file
    case "$DEPLOY_ENV" in
        "development")
            compose_file="docker-compose.dev.yml"
            ;;
        "production")
            compose_file="docker-compose.prod.yml"
            ;;
        *)
            log_error "未知的部署环境: $DEPLOY_ENV"
            exit 1
            ;;
    esac
    
    cd "$PROJECT_ROOT"
    
    # 部署服务
    docker-compose -f "$compose_file" up -d
    
    log_success "服务部署完成"
}

# 健康检查
health_check() {
    log_info "检查服务健康状态..."
    
    local max_attempts=30
    local attempt=1
    
    while [ $attempt -le $max_attempts ]; do
        if curl -f -s http://localhost:8080/health > /dev/null 2>&1; then
            log_success "服务健康检查通过"
            return 0
        fi
        
        log_info "等待服务启动... ($attempt/$max_attempts)"
        sleep 2
        ((attempt++))
    done
    
    log_error "服务健康检查失败"
    return 1
}

# 查看日志
show_logs() {
    log_info "显示服务日志..."
    
    local compose_file
    case "$DEPLOY_ENV" in
        "development")
            compose_file="docker-compose.dev.yml"
            ;;
        "production")
            compose_file="docker-compose.prod.yml"
            ;;
        *)
            compose_file="docker-compose.yml"
            ;;
    esac
    
    cd "$PROJECT_ROOT"
    docker-compose -f "$compose_file" logs -f --tail=100
}

# 停止服务
stop_services() {
    log_info "停止服务..."
    
    local compose_file
    case "$DEPLOY_ENV" in
        "development")
            compose_file="docker-compose.dev.yml"
            ;;
        "production")
            compose_file="docker-compose.prod.yml"
            ;;
        *)
            compose_file="docker-compose.yml"
            ;;
    esac
    
    cd "$PROJECT_ROOT"
    docker-compose -f "$compose_file" down
    
    log_success "服务已停止"
}

# 回滚部署
rollback_deployment() {
    log_warning "回滚功能待实现"
    log_info "请手动回滚到之前的版本"
}

# 主函数
main() {
    local environment=""
    local build_image_flag=false
    local migrate_flag=false
    local rollback_flag=false
    local health_flag=false
    local logs_flag=false
    local stop_flag=false
    
    # 解析参数
    while [[ $# -gt 0 ]]; do
        case $1 in
            "dev"|"development")
                environment="development"
                shift
                ;;
            "prod"|"production")
                environment="production"
                shift
                ;;
            "staging")
                environment="staging"
                shift
                ;;
            "--build")
                build_image_flag=true
                shift
                ;;
            "--migrate")
                migrate_flag=true
                shift
                ;;
            "--rollback")
                rollback_flag=true
                shift
                ;;
            "--health")
                health_flag=true
                shift
                ;;
            "--logs")
                logs_flag=true
                shift
                ;;
            "--stop")
                stop_flag=true
                shift
                ;;
            "--help"|"-h")
                show_help
                exit 0
                ;;
            *)
                log_error "未知参数: $1"
                show_help
                exit 1
                ;;
        esac
    done
    
    # 设置部署环境
    if [ -n "$environment" ]; then
        DEPLOY_ENV="$environment"
    fi
    
    log_info "部署环境: $DEPLOY_ENV"
    log_info "镜像: ${DOCKER_REGISTRY}/${IMAGE_NAME}:${VERSION}"
    
    # 检查依赖
    check_dependencies
    
    # 处理特殊操作
    if [ "$stop_flag" = true ]; then
        stop_services
        exit 0
    fi
    
    if [ "$logs_flag" = true ]; then
        show_logs
        exit 0
    fi
    
    if [ "$health_flag" = true ]; then
        health_check
        exit 0
    fi
    
    if [ "$rollback_flag" = true ]; then
        rollback_deployment
        exit 0
    fi
    
    # 检查环境文件
    if ! check_env_file; then
        exit 1
    fi
    
    # 构建镜像
    if [ "$build_image_flag" = true ]; then
        build_image
    fi
    
    # 运行迁移
    if [ "$migrate_flag" = true ]; then
        run_migration
    fi
    
    # 部署服务
    deploy_services
    
    # 健康检查
    if ! health_check; then
        log_error "部署后健康检查失败，请检查日志"
        show_logs
        exit 1
    fi
    
    log_success "🎉 部署完成！"
    log_info "服务地址: http://localhost:8080"
    log_info "健康检查: http://localhost:8080/health"
    log_info "查看日志: $0 $DEPLOY_ENV --logs"
}

# 执行主函数
main "$@"

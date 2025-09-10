#!/bin/bash

# ================================================
# TaskFlow 数据库迁移脚本
# ================================================

set -e

# 配置
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-3306}"
DB_USER="${DB_USER:-root}"
DB_PASSWORD="${DB_PASSWORD:-Aa13892349048!}"
DB_NAME="${DB_NAME:-taskflow}"

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
MIGRATIONS_DIR="${SCRIPT_DIR}/migrations"

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

# 检查MySQL连接
check_mysql_connection() {
    log_info "检查MySQL连接..."
    
    if ! command -v mysql &> /dev/null; then
        log_error "MySQL客户端未安装"
        exit 1
    fi
    
    if ! mysql -h"${DB_HOST}" -P"${DB_PORT}" -u"${DB_USER}" -p"${DB_PASSWORD}" -e "SELECT 1;" &> /dev/null; then
        log_error "无法连接到MySQL服务器"
        log_error "请检查连接参数: ${DB_USER}@${DB_HOST}:${DB_PORT}"
        exit 1
    fi
    
    log_success "MySQL连接正常"
}

# 创建数据库
create_database() {
    log_info "创建数据库 ${DB_NAME}..."
    
    mysql -h"${DB_HOST}" -P"${DB_PORT}" -u"${DB_USER}" -p"${DB_PASSWORD}" -e "
        CREATE DATABASE IF NOT EXISTS \`${DB_NAME}\` 
        CHARACTER SET utf8mb4 
        COLLATE utf8mb4_unicode_ci;
    " 2>/dev/null
    
    log_success "数据库 ${DB_NAME} 创建成功"
}

# 执行迁移文件
execute_migration() {
    local migration_file="$1"
    local filename=$(basename "$migration_file")
    
    log_info "执行迁移: ${filename}"
    
    if mysql -h"${DB_HOST}" -P"${DB_PORT}" -u"${DB_USER}" -p"${DB_PASSWORD}" "${DB_NAME}" < "$migration_file"; then
        log_success "迁移 ${filename} 执行成功"
        return 0
    else
        log_error "迁移 ${filename} 执行失败"
        return 1
    fi
}

# 创建迁移状态表
create_migration_table() {
    log_info "创建迁移状态表..."
    
    mysql -h"${DB_HOST}" -P"${DB_PORT}" -u"${DB_USER}" -p"${DB_PASSWORD}" "${DB_NAME}" -e "
        CREATE TABLE IF NOT EXISTS schema_migrations (
            version VARCHAR(255) PRIMARY KEY,
            executed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            checksum VARCHAR(32)
        ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='数据库迁移记录表';
    " 2>/dev/null
}

# 检查迁移是否已执行
is_migration_executed() {
    local migration_file="$1"
    local version=$(basename "$migration_file" .sql)
    
    local count=$(mysql -h"${DB_HOST}" -P"${DB_PORT}" -u"${DB_USER}" -p"${DB_PASSWORD}" "${DB_NAME}" -e "
        SELECT COUNT(*) FROM schema_migrations WHERE version = '${version}';
    " 2>/dev/null | tail -n 1)
    
    [ "$count" -gt 0 ]
}

# 记录迁移执行
record_migration() {
    local migration_file="$1"
    local version=$(basename "$migration_file" .sql)
    local checksum=$(md5sum "$migration_file" | cut -d' ' -f1)
    
    mysql -h"${DB_HOST}" -P"${DB_PORT}" -u"${DB_USER}" -p"${DB_PASSWORD}" "${DB_NAME}" -e "
        INSERT INTO schema_migrations (version, checksum) VALUES ('${version}', '${checksum}');
    " 2>/dev/null
}

# 运行所有迁移
run_migrations() {
    log_info "开始执行数据库迁移..."
    
    if [ ! -d "$MIGRATIONS_DIR" ]; then
        log_error "迁移目录不存在: ${MIGRATIONS_DIR}"
        exit 1
    fi
    
    # 创建迁移状态表
    create_migration_table
    
    # 按文件名排序执行迁移
    local migration_files=($(find "$MIGRATIONS_DIR" -name "*.sql" | sort))
    
    if [ ${#migration_files[@]} -eq 0 ]; then
        log_warning "没有找到迁移文件"
        return 0
    fi
    
    local success_count=0
    local total_count=${#migration_files[@]}
    local skipped_count=0
    
    for migration_file in "${migration_files[@]}"; do
        local filename=$(basename "$migration_file")
        
        if is_migration_executed "$migration_file"; then
            log_info "跳过已执行的迁移: ${filename}"
            ((skipped_count++))
            ((success_count++))
            continue
        fi
        
        if execute_migration "$migration_file"; then
            record_migration "$migration_file"
            ((success_count++))
        else
            log_error "迁移过程中出现错误，停止执行"
            break
        fi
    done
    
    log_info "迁移执行完成: ${success_count}/${total_count} 成功 (跳过: ${skipped_count})"
    
    if [ $success_count -eq $total_count ]; then
        log_success "所有迁移执行成功！"
        return 0
    else
        log_error "部分迁移执行失败"
        return 1
    fi
}

# 重置数据库
reset_database() {
    log_warning "警告：这将删除数据库 ${DB_NAME} 及其所有数据！"
    read -p "确定要继续吗？(yes/no): " confirm
    
    if [ "$confirm" != "yes" ]; then
        log_info "操作已取消"
        return 0
    fi
    
    log_info "删除数据库 ${DB_NAME}..."
    mysql -h"${DB_HOST}" -P"${DB_PORT}" -u"${DB_USER}" -p"${DB_PASSWORD}" -e "DROP DATABASE IF EXISTS \`${DB_NAME}\`;" 2>/dev/null
    
    log_info "重新创建数据库..."
    create_database
    
    log_info "执行迁移..."
    run_migrations
}

# 检查数据库状态
check_database_status() {
    log_info "检查数据库状态..."
    
    # 检查数据库是否存在
    local db_exists=$(mysql -h"${DB_HOST}" -P"${DB_PORT}" -u"${DB_USER}" -p"${DB_PASSWORD}" -e "SHOW DATABASES LIKE '${DB_NAME}';" 2>/dev/null | grep -c "${DB_NAME}" || true)
    
    if [ "$db_exists" -eq 0 ]; then
        log_warning "数据库 ${DB_NAME} 不存在"
        return 1
    fi
    
    # 检查表数量
    local table_count=$(mysql -h"${DB_HOST}" -P"${DB_PORT}" -u"${DB_USER}" -p"${DB_PASSWORD}" "${DB_NAME}" -e "SHOW TABLES;" 2>/dev/null | wc -l)
    table_count=$((table_count - 1)) # 减去表头行
    
    log_info "数据库 ${DB_NAME} 存在，包含 ${table_count} 个表"
    
    # 检查用户数量
    local user_count=$(mysql -h"${DB_HOST}" -P"${DB_PORT}" -u"${DB_USER}" -p"${DB_PASSWORD}" "${DB_NAME}" -e "SELECT COUNT(*) FROM users;" 2>/dev/null | tail -n 1)
    log_info "用户数量: ${user_count}"
    
    return 0
}

# 显示帮助信息
show_help() {
    echo "TaskFlow 数据库迁移工具"
    echo ""
    echo "用法: $0 [命令]"
    echo ""
    echo "命令:"
    echo "  migrate     执行数据库迁移（默认）"
    echo "  reset       重置数据库（删除并重建）"
    echo "  status      检查数据库状态"
    echo "  help        显示帮助信息"
    echo ""
    echo "环境变量:"
    echo "  DB_HOST     数据库主机 (默认: localhost)"
    echo "  DB_PORT     数据库端口 (默认: 3306)"
    echo "  DB_USER     数据库用户 (默认: root)"
    echo "  DB_PASSWORD 数据库密码 (默认: Aa13892349048!)"
    echo "  DB_NAME     数据库名称 (默认: taskflow)"
    echo ""
    echo "示例:"
    echo "  $0 migrate"
    echo "  DB_NAME=taskflow_test $0 reset"
}

# 主函数
main() {
    local command="${1:-migrate}"
    
    case "$command" in
        "migrate")
            check_mysql_connection
            create_database
            run_migrations
            ;;
        "reset")
            check_mysql_connection
            reset_database
            ;;
        "status")
            check_mysql_connection
            check_database_status
            ;;
        "help"|"-h"|"--help")
            show_help
            ;;
        *)
            log_error "未知命令: $command"
            show_help
            exit 1
            ;;
    esac
}

# 执行主函数
main "$@"

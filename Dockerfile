# ================================================
# 多阶段构建 - 生产级Dockerfile
# ================================================

# 构建阶段
FROM golang:1.21-alpine AS builder

# 设置必要的环境变量
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

# 安装必要的工具
RUN apk add --no-cache \
    git \
    ca-certificates \
    tzdata \
    make

# 设置工作目录
WORKDIR /app

# 复制go mod文件并下载依赖（利用Docker缓存）
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# 复制源代码
COPY . .

# 构建应用程序
RUN go build -a -installsuffix cgo -ldflags '-w -s -extldflags "-static"' -o bin/api ./cmd/api
RUN go build -a -installsuffix cgo -ldflags '-w -s -extldflags "-static"' -o bin/migrate ./cmd/migrate

# 最终运行阶段
FROM alpine:latest

# 安装运行时依赖
RUN apk --no-cache add ca-certificates tzdata curl

# 设置时区
ENV TZ=Asia/Shanghai

# 创建非root用户
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# 创建应用目录
RUN mkdir -p /app && chown -R appuser:appgroup /app

# 复制应用程序和配置文件
COPY --from=builder /app/bin/api /app/api
COPY --from=builder /app/bin/migrate /app/migrate
COPY --from=builder /app/configs /app/configs
COPY --from=builder /app/scripts /app/scripts

# 设置权限
RUN chmod +x /app/api /app/migrate /app/scripts/migrate.sh && \
    chown -R appuser:appgroup /app

# 设置工作目录
WORKDIR /app

# 暴露端口
EXPOSE 8080

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8080/health || exit 1

# 使用非root用户运行（安全最佳实践）
USER appuser:appgroup

# 运行应用程序
ENTRYPOINT ["/app/api"]
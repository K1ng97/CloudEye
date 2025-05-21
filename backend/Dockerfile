# 构建阶段
FROM golang:1.21-alpine AS builder

# 设置工作目录
WORKDIR /app

# 安装必要的系统依赖
RUN apk add --no-cache git

# 复制go.mod和go.sum文件并下载依赖
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -a -o cloud-eye-server .

# 运行阶段
FROM alpine:3.18

# 添加必要的CA证书
RUN apk --no-cache add ca-certificates tzdata

# 设置时区
ENV TZ=Asia/Shanghai

# 创建非root用户
RUN adduser -D -g '' appuser

# 创建应用目录
WORKDIR /app

# 创建日志和上传目录
RUN mkdir -p /app/logs /app/uploads/import /app/uploads/export && \
    chown -R appuser:appuser /app

# 从构建阶段复制可执行文件
COPY --from=builder /app/cloud-eye-server .
COPY --from=builder /app/configs /app/configs

# 使用非root用户运行
USER appuser

# 暴露端口
EXPOSE 8080

# 启动应用
CMD ["./cloud-eye-server"]
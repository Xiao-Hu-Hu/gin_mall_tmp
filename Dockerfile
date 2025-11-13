# 使用与你本地一致的 Go 版本
FROM golang:1.24-alpine AS builder

# 设置工作目录
WORKDIR /app

# 复制依赖文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd

# 使用更小的运行时镜像
FROM alpine:latest

# 安装 CA 证书（用于 HTTPS 请求）
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# 从构建阶段复制可执行文件
COPY --from=builder /app/main .
# 复制 Docker 配置文件
COPY --from=builder /app/conf/config.docker.ini ./conf/config.ini
# 复制静态文件目录
COPY --from=builder /app/static ./static

# 暴露端口
EXPOSE 3000

WORKDIR /root/

# 运行应用
CMD ["./main"]
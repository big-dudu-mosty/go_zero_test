# ============ 阶段 1：编译阶段 ============
FROM golang:1.24-alpine AS builder

# 安装必要的编译工具
RUN apk add --no-cache git

# 设置工作目录
WORKDIR /build

# 复制 go.mod 和 go.sum（利用 Docker 缓存）
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制所有代码
COPY . .

# 进入 user-rpc 目录并编译
WORKDIR /build/user-rpc
RUN go mod tidy && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o user-rpc .

# ============ 阶段 2：运行阶段 ============
FROM alpine:latest

# 安装必要的运行时依赖
RUN apk --no-cache add ca-certificates tzdata

# 设置工作目录
WORKDIR /app

# 从构建阶段复制可执行文件
COPY --from=builder /build/user-rpc/user-rpc .

# 复制配置文件
COPY --from=builder /build/user-rpc/etc ./etc

# 暴露 RPC 端口
EXPOSE 8080

# 启动命令
CMD ["./user-rpc", "-f", "etc/user.yaml"]

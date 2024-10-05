# 使用多阶段构建
FROM golang:1.22-alpine AS builder

# 设置工作目录
WORKDIR /app

# 复制go.mod和go.sum文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# 使用轻量级的alpine镜像作为最终镜像
FROM alpine:latest

# 安装必要的CA证书
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# 从builder阶段复制编译好的二进制文件
COPY --from=builder /app/main .

# 复制config.yaml文件
COPY config.yaml .

# 暴露端口
EXPOSE 8000

# 运行应用
CMD ["./main"]

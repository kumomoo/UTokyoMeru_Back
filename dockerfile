# 使用 Go 1.22.3 的官方镜像
FROM golang:1.22.3 AS builder

# 设置工作目录
WORKDIR /app

# 将本地的 Go 模块缓存复制到容器中
COPY go.mod ./
COPY go.sum ./

# 下载依赖项
RUN go mod download

# 将项目文件复制到容器中
COPY . .

# 设置交叉编译的目标平台和架构，例如 Linux 平台，amd64 架构
ENV GOOS=linux
ENV GOARCH=amd64
ENV CGO_ENABLED=0

# 编译 Go 程序，生成可执行文件
RUN go build -o /app/myapp

# 使用一个更小的基础镜像来打包最终的二进制文件
FROM alpine:latest

# 设置工作目录
WORKDIR /root/

# 从构建阶段复制生成的可执行文件
COPY --from=builder /app/myapp .

# 运行程序
CMD ["./myapp"]
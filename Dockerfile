# 构建阶段
FROM golang:1.25-alpine AS builder

# 设置工作目录
WORKDIR /app

# 复制go.mod和go.sum文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制所有源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -o duoduoyishan main.go

# 运行阶段
FROM alpine:latest

# 设置工作目录
WORKDIR /app

# 复制构建好的二进制文件
COPY --from=builder /app/duoduoyishan ./

# 复制配置文件
COPY config/config.yaml ./config/

# 创建必要的目录
RUN mkdir -p logs uploads static

# 暴露端口
EXPOSE 8080

# 设置环境变量
ENV GIN_MODE=release

# 启动命令
CMD ["./duoduoyishan"]

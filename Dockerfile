FROM golang:1.19-alpine AS builder

# 设置工作目录
WORKDIR /app

# 安装必要的工具
RUN apk add --no-cache git

# 复制 go.mod 和 go.sum 文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 编译应用
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o rainmcp ./cmd/rainmcp

# 使用 nginx 镜像作为最终镜像
FROM nginx:1.23-alpine

# 安装 ca-certificates 以支持 HTTPS 和 supervisor 来管理进程
RUN apk --no-cache add ca-certificates tzdata supervisor

# 设置工作目录
WORKDIR /app

# 从 builder 阶段复制编译好的二进制文件
COPY --from=builder /app/rainmcp .

# 设置时区为亚洲/上海
ENV TZ=Asia/Shanghai

# 设置环境变量默认值
ENV RAINBOND_HOST=127.0.0.1:8080
ENV RAINBOND_API=https://rainbond-api.example.com
ENV RAINBOND_TOKEN=""

# 创建 Nginx 配置文件，添加 CORS 支持
RUN mkdir -p /etc/nginx/conf.d
COPY nginx.conf /etc/nginx/conf.d/default.conf

# 创建 supervisor 配置文件
COPY supervisord.conf /etc/supervisor/conf.d/supervisord.conf

# 暴露 Nginx 端口
EXPOSE 80

# 使用 supervisor 启动 Nginx 和 Rainbond MCP 服务
CMD ["/usr/bin/supervisord", "-c", "/etc/supervisor/conf.d/supervisord.conf"]

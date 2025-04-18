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

# 使用轻量级的 alpine 镜像作为最终镜像
FROM alpine:3.17

# 安装 ca-certificates 以支持 HTTPS
RUN apk --no-cache add ca-certificates tzdata

# 设置工作目录
WORKDIR /app

# 从 builder 阶段复制编译好的二进制文件
COPY --from=builder /app/rainmcp .

# 设置时区为亚洲/上海
ENV TZ=Asia/Shanghai

# 暴露端口（根据应用需要调整）
EXPOSE 8080

# 设置环境变量默认值
ENV RAINBOND_HOST=0.0.0.0:8080
ENV RAINBOND_API=https://rainbond-api.example.com
ENV RAINBOND_TOKEN=""

# 运行应用
# 添加调试输出
ENTRYPOINT ["sh", "-c", "echo 'Starting Rainbond MCP Server...' && ./rainmcp"]

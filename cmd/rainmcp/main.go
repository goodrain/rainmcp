package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"rainmcp/pkg/services"
	"rainmcp/pkg/transport"

	"github.com/ThinkInAIXYZ/go-mcp/protocol"
	"github.com/ThinkInAIXYZ/go-mcp/server"
)

func main() {
	// 从环境变量获取配置
	host := getEnv("RAINBOND_HOST", "localhost:8080")
	rainbondAPI := getEnv("RAINBOND_API", "https://rainbond-api.example.com")
	rainbondToken := getEnv("RAINBOND_TOKEN", "")

	// 创建SSE服务器传输
	transportServer, err := transport.NewSSEServerTransport(host)
	if err != nil {
		log.Fatalf("Failed to create SSE server transport: %v", err)
	}

	// 创建MCP服务器
	mcpServer, err := server.NewServer(
		transportServer,
		// 设置服务器信息
		server.WithServerInfo(protocol.Implementation{
			Name:    "Rainbond MCP Server",
			Version: "1.0.0",
		}),
	)
	if err != nil {
		log.Fatalf("Failed to create MCP server: %v", err)
	}

	// 初始化服务
	serviceManager := services.NewManager(rainbondAPI, rainbondToken)

	// 注册所有工具
	registerTools(mcpServer, serviceManager)

	// 设置优雅关闭
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// 启动服务器
	go func() {
		fmt.Printf("Starting Rainbond MCP server on http://%s\n", host)
		fmt.Printf("SSE endpoint: http://%s/sse\n", host)
		fmt.Printf("Message endpoint: http://%s/message\n", host)
		if err := mcpServer.Run(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v\n", err)
		}
	}()

	// 等待关闭信号
	<-sigChan
	fmt.Println("Shutting down server...")

	// 创建一个带超时的上下文用于关闭
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 关闭服务器
	if err := mcpServer.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v\n", err)
	}

	fmt.Println("Server gracefully stopped")
}

// 注册所有工具
func registerTools(mcpServer *server.Server, serviceManager *services.Manager) {
	// 注册团队相关工具
	services.RegisterTeamTools(mcpServer, serviceManager)
	
	// 注册集群相关工具
	services.RegisterRegionTools(mcpServer, serviceManager)
	
	// 注册应用相关工具
	services.RegisterAppTools(mcpServer, serviceManager)
	
	// 注册组件相关工具
	services.RegisterComponentTools(mcpServer, serviceManager)
}

// 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

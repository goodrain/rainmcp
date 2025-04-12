package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ThinkInAIXYZ/go-mcp/protocol"
	"github.com/ThinkInAIXYZ/go-mcp/server"
	"github.com/ThinkInAIXYZ/go-mcp/transport"
)

// 定义计算器请求结构体
type CalculatorRequest struct {
	Operation string  `json:"operation" description:"要执行的运算（加法、减法、乘法、除法）" enum:"add,subtract,multiply,divide"`
	X         float64 `json:"x" description:"第一个数字"`
	Y         float64 `json:"y" description:"第二个数字"`
}

func main() {
	// 创建SSE服务器传输
	transportServer, err := transport.NewSSEServerTransport("localhost:8080")
	if err != nil {
		log.Fatalf("Failed to create SSE server transport: %v", err)
	}

	// 创建MCP服务器
	mcpServer, err := server.NewServer(
		transportServer,
		// 设置服务器信息
		server.WithServerInfo(protocol.Implementation{
			Name:    "Calculator Demo",
			Version: "1.0.0",
		}),
	)
	if err != nil {
		log.Fatalf("Failed to create MCP server: %v", err)
	}

	// 创建计算器工具
	calculatorTool, err := protocol.NewTool("calculate", "执行基本的算术运算，支持加减乘除四种运算", CalculatorRequest{})
	if err != nil {
		log.Fatalf("Failed to create calculator tool: %v", err)
		return
	}

	// 注册计算器工具
	mcpServer.RegisterTool(calculatorTool, func(request *protocol.CallToolRequest) (*protocol.CallToolResult, error) {
		// 解析请求参数
		req := new(CalculatorRequest)
		if err := protocol.VerifyAndUnmarshal(request.RawArguments, req); err != nil {
			return nil, fmt.Errorf("无效的计算器请求: %v", err)
		}

		// 执行计算
		var result float64
		switch req.Operation {
		case "add":
			result = req.X + req.Y
		case "subtract":
			result = req.X - req.Y
		case "multiply":
			result = req.X * req.Y
		case "divide":
			if req.Y == 0 {
				return nil, errors.New("除数不能为零")
			}
			result = req.X / req.Y
		default:
			return nil, fmt.Errorf("未知的运算类型: %s", req.Operation)
		}

		// 返回结果
		text := fmt.Sprintf("%.2f", result)
		return &protocol.CallToolResult{
			Content: []protocol.Content{
				protocol.TextContent{
					Type: "text",
					Text: text,
				},
			},
		}, nil
	})

	// 设置优雅关闭
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// 启动服务器
	go func() {
		fmt.Println("Starting MCP server on http://localhost:8080")
		fmt.Println("SSE endpoint: http://localhost:8080/sse")
		fmt.Println("Message endpoint: http://localhost:8080/message")
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

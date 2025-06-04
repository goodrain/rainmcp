package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"rainmcp/pkg/models"
	"syscall"
	"time"

	"rainmcp/pkg/logger"
	"rainmcp/pkg/services"

	"github.com/ThinkInAIXYZ/go-mcp/protocol"
	"github.com/ThinkInAIXYZ/go-mcp/server"
	"github.com/ThinkInAIXYZ/go-mcp/transport"
)

func main() {
	logger.Info("[启动] 开始启动Rainbond MCP服务器...")

	// 从环境变量获取配置
	host := getEnv("RAINBOND_HOST", ":8080") // 使用0.0.0.0允许从任何IP访问，适合Docker环境
	logger.Info("[配置] RAINBOND_HOST = %s", host)

	rainbondAPI := getEnv("RAINBOND_API", "https://rainbond-api.example.com")
	logger.Info("[配置] RAINBOND_API = %s", rainbondAPI)

	// 创建SSE服务器传输
	logger.Info("[初始化] 创建SSE服务器传输...")
	messageEndpointURL := "/message"
	rainTokenKey := "rainbond_token"
	paramKeysOpt := transport.WithSSEServerTransportAndHandlerOptionCopyParamKeys([]string{rainTokenKey})
	transportServer, mcpHandler, err := transport.NewSSEServerTransportAndHandler(messageEndpointURL, paramKeysOpt)

	if err != nil {
		logger.Fatal("[错误] 创建SSE服务器传输失败: %v", err)
	}
	router := http.NewServeMux()
	router.HandleFunc("/sse", mcpHandler.HandleSSE().ServeHTTP)
	router.HandleFunc(messageEndpointURL, func(w http.ResponseWriter, r *http.Request) {
		rainToken := r.URL.Query().Get(rainTokenKey)
		if rainToken == "" {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusBadRequest)
			if _, e := w.Write([]byte("lack rainbond_token")); e != nil {
				fmt.Printf("writeError: %+v", e)
			}
			return
		}
		r = r.WithContext(setRainTokenToCtx(r.Context(), rainToken))
		mcpHandler.HandleMessage().ServeHTTP(w, r)
	})
	httpServer := &http.Server{
		Addr:        host,
		Handler:     router,
		IdleTimeout: time.Minute,
	}
	logger.Info("[初始化] SSE服务器传输创建成功")

	// 创建MCP服务器
	logger.Info("[初始化] 创建MCP服务器...")
	mcpServer, err := server.NewServer(
		transportServer,
		// 设置服务器信息
		server.WithServerInfo(protocol.Implementation{
			Name:    "Rainbond MCP Server",
			Version: "1.0.0",
		}),
	)
	if err != nil {
		logger.Fatal("[错误] 创建MCP服务器失败: %v", err)
	}
	logger.Info("[初始化] MCP服务器创建成功")

	// 初始化服务
	logger.Info("[初始化] 创建服务管理器...")
	serviceManager := services.NewManager(rainbondAPI)
	logger.Info("[初始化] 服务管理器创建成功")

	// 注册所有工具
	logger.Info("[初始化] 注册所有工具...")
	registerTools(mcpServer, serviceManager)
	logger.Info("[初始化] 所有工具注册完成")

	// 设置优雅关闭
	logger.Info("[初始化] 设置信号处理...")
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// 启动服务器
	logger.Info("[启动] 开始启动服务器...")

	go func() {
		if err = httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal("[错误] 服务器错误: %v\n", err)
		}
	}()
	go func() {
		logger.Info("[启动] 服务器开始运行...")
		if err := mcpServer.Run(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("[错误] 服务器错误: %v\n", err)
		}
	}()

	// 等待关闭信号
	logger.Info("[信息] 服务器已启动，等待关闭信号...")
	<-sigChan
	logger.Info("[关闭] 正在关闭服务器...")

	// 创建一个带超时的上下文用于关闭
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 关闭服务器
	logger.Info("[关闭] 正在优雅关闭服务器...")
	if err := mcpServer.Shutdown(ctx); err != nil {
		logger.Fatal("[错误] 服务器关闭失败: %v\n", err)
	}

	logger.Info("[关闭] 服务器已优雅关闭")
}

// 注册所有工具
func registerTools(mcpServer *server.Server, serviceManager *services.Manager) {
	// 检查服务管理器
	if serviceManager == nil {
		logger.Error("[错误] 服务管理器为空，无法注册工具")
		return
	}

	// 注册团队相关工具
	services.RegisterTeamTools(mcpServer, serviceManager)

	// 注册集群相关工具
	services.RegisterRegionTools(mcpServer, serviceManager)

	// 注册应用相关工具
	services.RegisterAppTools(mcpServer, serviceManager)

	// 注册组件相关工具
	services.RegisterComponentTools(mcpServer, serviceManager)

	logger.Info("[工具] 所有工具注册完成")
}

// getEnv 获取环境变量值，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func setRainTokenToCtx(ctx context.Context, rainToken string) context.Context {
	return context.WithValue(ctx, models.RainTokenKey{}, rainToken)
}

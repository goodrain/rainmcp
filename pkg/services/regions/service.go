package regions

import (
	"encoding/json"
	"fmt"
	"rainmcp/pkg/api"
	"rainmcp/pkg/utils"

	"github.com/ThinkInAIXYZ/go-mcp/protocol"
	"github.com/ThinkInAIXYZ/go-mcp/server"
)

// Service 处理集群相关的API请求
type Service struct {
	client *api.Client
}

// NewService 创建一个新的集群服务
func NewService(client *api.Client) *Service {
	utils.Debug("创建新的集群服务")
	return &Service{
		client: client,
	}
}

// GetBaseURL 获取API基础URL
func (s *Service) GetBaseURL() string {
	if s == nil || s.client == nil {
		utils.Error("获取API基础URL失败: 服务或客户端为空")
		return ""
	}
	return s.client.BaseURL
}

// RegisterTools 注册集群相关的工具
func RegisterTools(mcpServer *server.Server, service *Service) {
	utils.Info("注册集群相关工具...")
	
	// 注册获取集群列表工具
	regionsListTool, err := protocol.NewTool(
		"rainbond_regions",
		"获取Rainbond平台中的集群列表",
		struct{}{}, // 无需参数
	)
	if err != nil {
		utils.Error("创建集群列表工具失败: %v", err)
		return
	}

	utils.Debug("成功创建集群列表工具，正在注册...")
	mcpServer.RegisterTool(regionsListTool, service.handleRegionsList)
	utils.Info("集群相关工具注册完成")
}

// handleRegionsList 处理获取集群列表的请求
func (s *Service) handleRegionsList(request *protocol.CallToolRequest) (*protocol.CallToolResult, error) {
	utils.Info("获取集群列表")
	
	// 检查API客户端是否正确初始化
	if s.client == nil {
		utils.Error("API客户端未初始化")
		return nil, fmt.Errorf("无法获取集群列表: API客户端未初始化")
	}
	
	utils.Debug("使用API地址: %s", s.client.BaseURL)
	
	// 调用Rainbond API获取集群列表
	utils.Debug("调用API获取集群列表: /openapi/v1/regions")
	resp, err := s.client.Get("/openapi/v1/regions")
	if err != nil {
		utils.Error("获取集群列表失败: %v", err)
		return nil, fmt.Errorf("获取集群列表失败: %v", err)
	}

	utils.Debug("获取到API响应: %s", string(resp))

	// 直接使用API响应作为结果
	utils.Info("成功获取集群列表数据")
	
	// 将响应转换为可读性更好的JSON格式
	var jsonData interface{}
	if err := json.Unmarshal(resp, &jsonData); err != nil {
		utils.Warn("响应不是有效的JSON格式，直接返回原始数据")
		// 如果不是JSON，直接返回原始字符串
		jsonData = string(resp)
	}
	
	// 将结果转换为格式化的JSON字符串
	utils.Debug("格式化响应数据...")
	resultJSON, err := json.MarshalIndent(jsonData, "", "  ")
	if err != nil {
		utils.Error("格式化响应数据失败: %v", err)
		// 如果格式化失败，直接返回原始数据
		resultJSON = resp
	}

	utils.Debug("成功序列化集群列表，准备返回结果")
	// 返回结果
	return &protocol.CallToolResult{
		Content: []protocol.Content{
			protocol.TextContent{
				Type: "text",
				Text: string(resultJSON),
			},
		},
	}, nil
}

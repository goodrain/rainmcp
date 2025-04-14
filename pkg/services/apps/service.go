package apps

import (
	"encoding/json"
	"fmt"
	"rainmcp/pkg/api"
	"rainmcp/pkg/models"
	"rainmcp/pkg/utils"

	"github.com/ThinkInAIXYZ/go-mcp/protocol"
	"github.com/ThinkInAIXYZ/go-mcp/server"
)

// Service 处理应用相关的API请求
type Service struct {
	client *api.Client
}

// NewService 创建一个新的应用服务
func NewService(client *api.Client) *Service {
	utils.Debug("创建新的应用服务")
	return &Service{
		client: client,
	}
}

// RegisterTools 注册应用相关的工具
func RegisterTools(mcpServer *server.Server, service *Service) {
	utils.Info("注册应用相关工具...")

	// 注册获取应用列表工具
	appsListTool, err := protocol.NewTool(
		"rainbond_apps",
		"获取Rainbond平台中的应用列表",
		models.AppsRequest{},
	)
	if err != nil {
		utils.Error("创建应用列表工具失败: %v", err)
		return
	}
	mcpServer.RegisterTool(appsListTool, service.handleAppsList)

	// 注册创建应用工具
	createAppTool, err := protocol.NewTool(
		"rainbond_create_app",
		"在Rainbond平台中创建应用",
		models.CreateAppRequest{},
	)
	if err != nil {
		utils.Error("创建应用创建工具失败: %v", err)
		return
	}
	mcpServer.RegisterTool(createAppTool, service.handleCreateApp)

	utils.Info("应用相关工具注册完成")
}

// handleAppsList 处理获取应用列表的请求
func (s *Service) handleAppsList(request *protocol.CallToolRequest) (*protocol.CallToolResult, error) {
	// 解析请求参数
	req := new(models.AppsRequest)
	if err := protocol.VerifyAndUnmarshal(request.RawArguments, req); err != nil {
		utils.Error("解析应用列表请求失败: %v", err)
		return nil, fmt.Errorf("无效的应用列表请求: %v", err)
	}

	// 构建API路径 - 根据Rainbond OpenAPI文档
	path := fmt.Sprintf("/openapi/v1/teams/%s/regions/%s/apps", req.TeamName, req.RegionName)

	utils.Info("获取应用列表")

	// 调用Rainbond API获取应用列表
	resp, err := s.client.Get(path)
	if err != nil {
		utils.Error("获取应用列表失败: %v", err)
		return nil, fmt.Errorf("获取应用列表失败: %v", err)
	}

	// 直接使用原始响应
	var result interface{}
	if err := json.Unmarshal(resp, &result); err != nil {
		utils.Error("解析应用列表响应失败: %v", err)
		return nil, fmt.Errorf("解析应用列表失败: %v", err)
	}

	// 将结果转换为JSON字符串
	resultJSON, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		utils.Error("序列化应用列表结果失败: %v", err)
		return nil, fmt.Errorf("序列化应用列表失败: %v", err)
	}

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

// handleCreateApp 处理创建应用的请求
func (s *Service) handleCreateApp(request *protocol.CallToolRequest) (*protocol.CallToolResult, error) {
	// 解析请求参数
	req := new(models.CreateAppRequest)
	if err := protocol.VerifyAndUnmarshal(request.RawArguments, req); err != nil {
		utils.Error("解析创建应用请求失败: %v", err)
		return nil, fmt.Errorf("无效的创建应用请求: %v", err)
	}

	// 构建API路径 - 根据Rainbond OpenAPI文档
	path := fmt.Sprintf("/openapi/v1/teams/%s/regions/%s/apps", req.TeamName, req.RegionName)

	utils.Info("创建应用: %s, 应用名称: %s", path, req.AppName)

	// 调用Rainbond API创建应用
	resp, err := s.client.Post(path, req)
	if err != nil {
		utils.Error("创建应用失败: %v", err)
		return nil, fmt.Errorf("创建应用失败: %v", err)
	}

	// 解析响应
	var result map[string]interface{}
	if err := json.Unmarshal(resp, &result); err != nil {
		utils.Error("解析创建应用响应失败: %v", err)
		return nil, fmt.Errorf("解析创建应用响应失败: %v", err)
	}

	// 将结果转换为JSON字符串
	resultJSON, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		utils.Error("序列化创建应用结果失败: %v", err)
		return nil, fmt.Errorf("序列化创建应用结果失败: %v", err)
	}

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

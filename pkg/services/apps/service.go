package apps

import (
	"encoding/json"
	"fmt"
	"rainmcp/pkg/api"
	"rainmcp/pkg/models"

	"github.com/ThinkInAIXYZ/go-mcp/protocol"
	"github.com/ThinkInAIXYZ/go-mcp/server"
)

// Service 处理应用相关的API请求
type Service struct {
	client *api.Client
}

// NewService 创建一个新的应用服务
func NewService(client *api.Client) *Service {
	return &Service{
		client: client,
	}
}

// RegisterTools 注册应用相关的工具
func RegisterTools(mcpServer *server.Server, service *Service) {
	// 注册获取应用列表工具
	appsListTool, err := protocol.NewTool(
		"rainbond_apps",
		"获取Rainbond平台中的应用列表",
		models.AppsRequest{},
	)
	if err != nil {
		fmt.Printf("Failed to create apps list tool: %v\n", err)
		return
	}

	mcpServer.RegisterTool(appsListTool, service.handleAppsList)
}

// handleAppsList 处理获取应用列表的请求
func (s *Service) handleAppsList(request *protocol.CallToolRequest) (*protocol.CallToolResult, error) {
	// 解析请求参数
	req := new(models.AppsRequest)
	if err := protocol.VerifyAndUnmarshal(request.RawArguments, req); err != nil {
		return nil, fmt.Errorf("无效的应用列表请求: %v", err)
	}

	// 构建API路径
	path := fmt.Sprintf("/console/teams/%s/regions/%s/apps", req.TeamName, req.RegionName)
	
	// 调用Rainbond API获取应用列表
	resp, err := s.client.Get(path)
	if err != nil {
		return nil, fmt.Errorf("获取应用列表失败: %v", err)
	}

	// 解析响应
	var apps []models.App
	if err := json.Unmarshal(resp, &apps); err != nil {
		return nil, fmt.Errorf("解析应用列表失败: %v", err)
	}

	// 将应用列表转换为JSON字符串
	appsJSON, err := json.MarshalIndent(apps, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("序列化应用列表失败: %v", err)
	}

	// 返回结果
	return &protocol.CallToolResult{
		Content: []protocol.Content{
			protocol.TextContent{
				Type: "text",
				Text: string(appsJSON),
			},
		},
	}, nil
}

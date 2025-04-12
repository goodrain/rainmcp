package teams

import (
	"encoding/json"
	"fmt"
	"rainmcp/pkg/api"
	"rainmcp/pkg/models"

	"github.com/ThinkInAIXYZ/go-mcp/protocol"
	"github.com/ThinkInAIXYZ/go-mcp/server"
)

// Service 处理团队相关的API请求
type Service struct {
	client *api.Client
}

// NewService 创建一个新的团队服务
func NewService(client *api.Client) *Service {
	return &Service{
		client: client,
	}
}

// RegisterTools 注册团队相关的工具
func RegisterTools(mcpServer *server.Server, service *Service) {
	// 注册获取团队列表工具
	teamsListTool, err := protocol.NewTool(
		"rainbond_teams",
		"获取Rainbond平台中的团队列表",
		struct{}{}, // 无需参数
	)
	if err != nil {
		fmt.Printf("Failed to create teams list tool: %v\n", err)
		return
	}

	mcpServer.RegisterTool(teamsListTool, service.handleTeamsList)
}

// handleTeamsList 处理获取团队列表的请求
func (s *Service) handleTeamsList(request *protocol.CallToolRequest) (*protocol.CallToolResult, error) {
	// 调用Rainbond API获取团队列表
	resp, err := s.client.Get("/console/teams/all_teams")
	if err != nil {
		return nil, fmt.Errorf("获取团队列表失败: %v", err)
	}

	// 解析响应
	var teams []models.Team
	if err := json.Unmarshal(resp, &teams); err != nil {
		return nil, fmt.Errorf("解析团队列表失败: %v", err)
	}

	// 将团队列表转换为JSON字符串
	teamsJSON, err := json.MarshalIndent(teams, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("序列化团队列表失败: %v", err)
	}

	// 返回结果
	return &protocol.CallToolResult{
		Content: []protocol.Content{
			protocol.TextContent{
				Type: "text",
				Text: string(teamsJSON),
			},
		},
	}, nil
}

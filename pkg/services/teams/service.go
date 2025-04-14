package teams

import (
	"encoding/json"
	"fmt"
	"rainmcp/pkg/api"
	"rainmcp/pkg/utils"

	"github.com/ThinkInAIXYZ/go-mcp/protocol"
	"github.com/ThinkInAIXYZ/go-mcp/server"
)

// Service 处理团队相关的API请求
type Service struct {
	client *api.Client
}

// NewService 创建一个新的团队服务
func NewService(client *api.Client) *Service {
	utils.Debug("创建新的团队服务")
	return &Service{
		client: client,
	}
}

// RegisterTools 注册团队相关的工具
func RegisterTools(mcpServer *server.Server, service *Service) {
	utils.Info("注册团队相关工具...")
	
	// 注册获取团队列表工具
	teamsListTool, err := protocol.NewTool(
		"rainbond_teams",
		"获取Rainbond平台中的团队列表",
		struct{}{}, // 无需参数
	)
	if err != nil {
		utils.Error("创建团队列表工具失败: %v", err)
		return
	}

	mcpServer.RegisterTool(teamsListTool, service.handleTeamsList)
	
	utils.Info("团队相关工具注册完成")
}

// handleTeamsList 处理获取团队列表的请求
func (s *Service) handleTeamsList(request *protocol.CallToolRequest) (*protocol.CallToolResult, error) {
	utils.Info("获取团队列表")
	
	// 调用Rainbond API获取团队列表
	resp, err := s.client.Get("/openapi/v1/teams")
	if err != nil {
		utils.Error("获取团队列表失败: %v", err)
		return nil, fmt.Errorf("获取团队列表失败: %v", err)
	}

	// 直接使用API响应作为结果
	utils.Info("成功获取团队列表数据")
	
	// 将响应转换为可读性更好的JSON格式
	var jsonData interface{}
	if err := json.Unmarshal(resp, &jsonData); err != nil {
		utils.Warn("响应不是有效的JSON格式，直接返回原始数据")
		// 如果不是JSON，直接返回原始字符串
		jsonData = string(resp)
	}
	
	// 将结果转换为格式化的JSON字符串
	resultJSON, err := json.MarshalIndent(jsonData, "", "  ")
	if err != nil {
		utils.Error("格式化响应数据失败: %v", err)
		// 如果格式化失败，直接返回原始数据
		resultJSON = resp
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

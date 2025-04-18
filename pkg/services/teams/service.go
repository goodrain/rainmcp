package teams

import (
	"encoding/json"
	"fmt"
	"rainmcp/pkg/api"
	"rainmcp/pkg/logger"
	"rainmcp/pkg/models"
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
	logger.Debug("创建新的团队服务")
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
		logger.Error("创建团队列表工具失败: %v", err)
		return
	}
	mcpServer.RegisterTool(teamsListTool, service.handleTeamsList)
}

// handleTeamsList 处理获取团队列表的请求
func (s *Service) handleTeamsList(request *protocol.CallToolRequest) (*protocol.CallToolResult, error) {
	logger.Info("获取团队列表")

	// 调用Rainbond API获取团队列表
	resp, err := s.client.Get("/openapi/v1/teams")
	if err != nil {
		logger.Error("获取团队列表失败: %v", err)
		return nil, fmt.Errorf("获取团队列表失败: %v", err)
	}

	// 直接使用API响应作为结果
	logger.Info("成功获取团队列表数据")

	// 使用TeamsResponse结构体解析响应
	var teamsResp models.TeamsResponse
	if err := json.Unmarshal(resp, &teamsResp); err != nil {
		logger.Warn("解析团队列表响应失败: %v", err)
		logger.Debug("原始响应数据: %s", string(resp))

		// 如果解析失败，尝试解析为通用JSON
		var jsonData interface{}
		if err := json.Unmarshal(resp, &jsonData); err != nil {
			logger.Warn("响应不是有效的JSON格式，直接返回原始数据")
			// 如果不是JSON，直接返回原始字符串
			jsonData = string(resp)
		}

		// 将结果转换为格式化的JSON字符串
		resultJSON, err := json.MarshalIndent(jsonData, "", "  ")
		if err != nil {
			logger.Error("格式化响应数据失败: %v", err)
			// 如果格式化失败，直接返回原始数据
			resultJSON = resp
		}

		return &protocol.CallToolResult{
			Content: []protocol.Content{
				protocol.TextContent{
					Type: "text",
					Text: string(resultJSON),
				},
			},
		}, nil
	}

	// 成功解析为TeamsResponse结构体
	logger.Info("成功解析团队列表数据，共有 %d 个团队", len(teamsResp.Tenants))

	// 将结果转换为包含字段描述的JSON字符串
	resultJSON, err := utils.MarshalJSONWithDescription(teamsResp)
	if err != nil {
		logger.Error("格式化响应数据失败: %v", err)
		// 如果格式化失败，尝试使用标准JSON序列化
		resultJSON, err = json.MarshalIndent(teamsResp, "", "  ")
		if err != nil {
			logger.Error("标准JSON格式化也失败: %v", err)
			// 如果标准格式化也失败，直接返回原始数据
			resultJSON = resp
		}
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

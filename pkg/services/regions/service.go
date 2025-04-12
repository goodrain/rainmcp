package regions

import (
	"encoding/json"
	"fmt"
	"rainmcp/pkg/api"
	"rainmcp/pkg/models"

	"github.com/ThinkInAIXYZ/go-mcp/protocol"
	"github.com/ThinkInAIXYZ/go-mcp/server"
)

// Service 处理集群相关的API请求
type Service struct {
	client *api.Client
}

// NewService 创建一个新的集群服务
func NewService(client *api.Client) *Service {
	return &Service{
		client: client,
	}
}

// RegisterTools 注册集群相关的工具
func RegisterTools(mcpServer *server.Server, service *Service) {
	// 注册获取集群列表工具
	regionsListTool, err := protocol.NewTool(
		"rainbond_regions",
		"获取Rainbond平台中的集群列表",
		struct{}{}, // 无需参数
	)
	if err != nil {
		fmt.Printf("Failed to create regions list tool: %v\n", err)
		return
	}

	mcpServer.RegisterTool(regionsListTool, service.handleRegionsList)
}

// handleRegionsList 处理获取集群列表的请求
func (s *Service) handleRegionsList(request *protocol.CallToolRequest) (*protocol.CallToolResult, error) {
	// 调用Rainbond API获取集群列表
	resp, err := s.client.Get("/console/regions")
	if err != nil {
		return nil, fmt.Errorf("获取集群列表失败: %v", err)
	}

	// 解析响应
	var regions []models.Region
	if err := json.Unmarshal(resp, &regions); err != nil {
		return nil, fmt.Errorf("解析集群列表失败: %v", err)
	}

	// 将集群列表转换为JSON字符串
	regionsJSON, err := json.MarshalIndent(regions, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("序列化集群列表失败: %v", err)
	}

	// 返回结果
	return &protocol.CallToolResult{
		Content: []protocol.Content{
			protocol.TextContent{
				Type: "text",
				Text: string(regionsJSON),
			},
		},
	}, nil
}

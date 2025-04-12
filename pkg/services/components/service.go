package components

import (
	"encoding/json"
	"fmt"
	"rainmcp/pkg/api"
	"rainmcp/pkg/models"

	"github.com/ThinkInAIXYZ/go-mcp/protocol"
	"github.com/ThinkInAIXYZ/go-mcp/server"
)

// Service 处理组件相关的API请求
type Service struct {
	client *api.Client
}

// NewService 创建一个新的组件服务
func NewService(client *api.Client) *Service {
	return &Service{
		client: client,
	}
}

// RegisterTools 注册组件相关的工具
func RegisterTools(mcpServer *server.Server, service *Service) {
	// 注册构建组件工具
	buildServiceTool, err := protocol.NewTool(
		"rainbond_build_service",
		"在Rainbond平台中构建组件",
		models.ComponentRequest{},
	)
	if err != nil {
		fmt.Printf("Failed to create build service tool: %v\n", err)
		return
	}

	mcpServer.RegisterTool(buildServiceTool, service.handleBuildService)
}

// handleBuildService 处理构建组件的请求
func (s *Service) handleBuildService(request *protocol.CallToolRequest) (*protocol.CallToolResult, error) {
	// 解析请求参数
	req := new(models.ComponentRequest)
	if err := protocol.VerifyAndUnmarshal(request.RawArguments, req); err != nil {
		return nil, fmt.Errorf("无效的组件构建请求: %v", err)
	}

	// 构建API路径
	path := fmt.Sprintf("/console/teams/%s/regions/%s/apps/%s/services/%s/build", 
		req.TeamName, req.RegionName, req.AppID, req.ServiceID)
	
	// 准备构建参数
	buildParams := map[string]interface{}{
		"is_deploy":      req.IsDeploy,
		"service_id":     req.ServiceID,
		"build_version":  req.BuildVersion,
	}
	
	// 调用Rainbond API构建组件
	resp, err := s.client.Post(path, buildParams)
	if err != nil {
		return nil, fmt.Errorf("构建组件失败: %v", err)
	}

	// 解析响应
	var result map[string]interface{}
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("解析构建结果失败: %v", err)
	}

	// 将构建结果转换为JSON字符串
	resultJSON, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("序列化构建结果失败: %v", err)
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

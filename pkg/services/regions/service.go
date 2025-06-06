package regions

import (
	"context"
	"encoding/json"
	"fmt"
	"rainmcp/pkg/api"
	"rainmcp/pkg/logger"
	"rainmcp/pkg/models"
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
	logger.Debug("创建新的集群服务")
	return &Service{
		client: client,
	}
}

// GetBaseURL 获取API基础URL
func (s *Service) GetBaseURL() string {
	if s == nil || s.client == nil {
		logger.Error("获取API基础URL失败: 服务或客户端为空")
		return ""
	}
	return s.client.BaseURL
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
		logger.Error("创建集群列表工具失败: %v", err)
		return
	}

	logger.Debug("成功创建集群列表工具，正在注册...")
	mcpServer.RegisterTool(regionsListTool, service.handleRegionsList)
}

// handleRegionsList 处理获取集群列表的请求
func (s *Service) handleRegionsList(ctx context.Context, request *protocol.CallToolRequest) (*protocol.CallToolResult, error) {
	logger.Info("获取集群列表")
	rainTokenValue := ctx.Value(models.RainTokenKey{})
	rainToken := rainTokenValue.(string)
	s.client.Token = rainToken
	// 检查API客户端是否正确初始化
	if s.client == nil {
		logger.Error("API客户端未初始化")
		return nil, fmt.Errorf("无法获取集群列表: API客户端未初始化")
	}

	logger.Debug("使用API地址: %s", s.client.BaseURL)

	// 调用Rainbond API获取集群列表
	logger.Debug("调用API获取集群列表: /openapi/v1/mcp/regions")
	resp, err := s.client.Get("/openapi/v1/mcp/regions")
	if err != nil {
		logger.Error("获取集群列表失败: %v", err)
		return nil, fmt.Errorf("获取集群列表失败: %v", err)
	}

	logger.Debug("获取到API响应: %s", string(resp))

	// 直接使用API响应作为结果
	logger.Info("成功获取集群列表数据")

	// 使用RegionsResponse结构体解析响应
	var regionsResp models.RegionsResponse
	if err := json.Unmarshal(resp, &regionsResp); err != nil {
		logger.Warn("解析集群列表响应失败: %v", err)
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
				&protocol.TextContent{
					Type: "text",
					Text: utils.FormatJSON(resultJSON),
				},
			},
		}, nil
	}

	// 成功解析为RegionsResponse结构体
	var regionCount int
	if len(regionsResp.Data.List) > 0 {
		regionCount = len(regionsResp.Data.List)
	} else {
		regionCount = len(regionsResp.Data.List)
	}
	logger.Info("成功解析集群列表数据，共有 %d 个集群", regionCount)

	// 将结果转换为包含字段描述的JSON字符串
	logger.Debug("格式化响应数据，包含字段描述...")
	resultJSON, err := utils.MarshalJSONWithDescription(regionsResp)
	if err != nil {
		logger.Error("带描述的格式化响应数据失败: %v", err)
		// 如果格式化失败，尝试使用标准JSON序列化
		resultJSON, err = json.MarshalIndent(regionsResp, "", "  ")
		if err != nil {
			logger.Error("标准JSON格式化也失败: %v", err)
			// 如果标准格式化也失败，直接返回原始数据
			resultJSON = resp
		}
	}

	logger.Debug("成功序列化集群列表，准备返回结果")
	// 返回结果
	return &protocol.CallToolResult{
		Content: []protocol.Content{
			&protocol.TextContent{
				Type: "text",
				Text: utils.FormatJSON(resultJSON),
			},
		},
	}, nil
}

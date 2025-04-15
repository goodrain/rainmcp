package components

import (
	"encoding/json"
	"fmt"
	"rainmcp/pkg/api"
	"rainmcp/pkg/models"
	"rainmcp/pkg/utils"
	"strings"

	"github.com/ThinkInAIXYZ/go-mcp/protocol"
	"github.com/ThinkInAIXYZ/go-mcp/server"
)

// Service 处理组件相关的API请求
type Service struct {
	client *api.Client
}

// NewService 创建一个新的组件服务
func NewService(client *api.Client) *Service {
	utils.Debug("创建新的组件服务")
	return &Service{
		client: client,
	}
}

// RegisterTools 注册组件相关的工具
func RegisterTools(mcpServer *server.Server, service *Service) {
	utils.Info("注册组件相关工具...")

	// 注册获取组件详情工具
	componentDetailTool, err := protocol.NewTool(
		"rainbond_get_component_detail",
		"获取Rainbond平台中的组件详情",
		models.ComponentDetailRequest{},
	)
	if err != nil {
		utils.Error("创建组件详情工具失败: %v", err)
		return
	}
	mcpServer.RegisterTool(componentDetailTool, service.handleGetComponentDetail)

	// 注册创建组件工具
	createComponentTool, err := protocol.NewTool(
		"rainbond_create_component",
		"在Rainbond平台中创建组件",
		models.CreateComponentRequest{},
	)
	if err != nil {
		utils.Error("创建组件创建工具失败: %v", err)
		return
	}
	mcpServer.RegisterTool(createComponentTool, service.handleCreateComponent)

	// 注册构建组件工具
	buildServiceTool, err := protocol.NewTool(
		"rainbond_build_service",
		"在Rainbond平台中构建组件",
		models.BuildComponentRequest{},
	)
	if err != nil {
		utils.Error("创建组件构建工具失败: %v", err)
		return
	}
	mcpServer.RegisterTool(buildServiceTool, service.handleBuildService)

	utils.Info("组件相关工具注册完成")
}

// handleGetComponentDetail 处理获取组件详情的请求
func (s *Service) handleGetComponentDetail(request *protocol.CallToolRequest) (*protocol.CallToolResult, error) {
	// 解析请求参数
	req := new(models.ComponentDetailRequest)
	if err := protocol.VerifyAndUnmarshal(request.RawArguments, req); err != nil {
		// 记录原始错误
		utils.Error("解析获取组件详情请求失败: %v", err)

		// 尝试解析原始请求数据
		var rawData map[string]interface{}
		var detailedErrMsg string
		if jsonErr := json.Unmarshal(request.RawArguments, &rawData); jsonErr == nil {
			// 检查必填字段
			requiredFields := []string{"team_name", "region_name", "app_id", "service_id"}
			var missingFields []string

			for _, field := range requiredFields {
				if _, exists := rawData[field]; !exists {
					missingFields = append(missingFields, field)
				}
			}

			// 构建详细错误信息
			if len(missingFields) > 0 {
				detailedErrMsg = fmt.Sprintf("缺少必填字段: %s", strings.Join(missingFields, ", "))
			} else {
				detailedErrMsg = fmt.Sprintf("请求参数验证失败: %v", err)
			}
		} else {
			detailedErrMsg = fmt.Sprintf("解析JSON数据失败: %v", jsonErr)
		}

		// 返回带有详细错误信息的响应
		return &protocol.CallToolResult{
			Content: []protocol.Content{
				protocol.TextContent{
					Type: "text",
					Text: detailedErrMsg,
				},
			},
			IsError: true,
		}, nil
	}

	// 参数校验
	if req.TeamName == "" {
		errMsg := "缺少必填字段: team_name"
		utils.Error(errMsg)
		return &protocol.CallToolResult{
			Content: []protocol.Content{
				protocol.TextContent{
					Type: "text",
					Text: errMsg,
				},
			},
			IsError: true,
		}, nil
	}

	if req.RegionName == "" {
		errMsg := "缺少必填字段: region_name"
		utils.Error(errMsg)
		return &protocol.CallToolResult{
			Content: []protocol.Content{
				protocol.TextContent{
					Type: "text",
					Text: errMsg,
				},
			},
			IsError: true,
		}, nil
	}

	if req.AppID == "" {
		errMsg := "缺少必填字段: app_id"
		utils.Error(errMsg)
		return &protocol.CallToolResult{
			Content: []protocol.Content{
				protocol.TextContent{
					Type: "text",
					Text: errMsg,
				},
			},
			IsError: true,
		}, nil
	}

	if req.ServiceID == "" {
		errMsg := "缺少必填字段: service_id"
		utils.Error(errMsg)
		return &protocol.CallToolResult{
			Content: []protocol.Content{
				protocol.TextContent{
					Type: "text",
					Text: errMsg,
				},
			},
			IsError: true,
		}, nil
	}

	// 构建API路径 - 根据Rainbond OpenAPI文档
	path := fmt.Sprintf("/openapi/v1/teams/%s/regions/%s/apps/%s/services/%s",
		req.TeamName, req.RegionName, req.AppID, req.ServiceID)

	utils.Info("获取组件详情: %s", path)

	// 调用Rainbond API获取组件详情
	resp, err := s.client.Get(path)
	if err != nil {
		errMsg := fmt.Sprintf("获取组件详情失败: %v", err)
		utils.Error(errMsg)
		return &protocol.CallToolResult{
			Content: []protocol.Content{
				protocol.TextContent{
					Type: "text",
					Text: errMsg,
				},
			},
			IsError: true,
		}, nil
	}

	// 解析响应
	var detailResp models.ComponentDetailResponse
	if err := json.Unmarshal(resp, &detailResp); err != nil {
		errMsg := fmt.Sprintf("解析组件详情响应失败: %v", err)
		utils.Error(errMsg)
		return &protocol.CallToolResult{
			Content: []protocol.Content{
				protocol.TextContent{
					Type: "text",
					Text: errMsg,
				},
			},
			IsError: true,
		}, nil
	}

	// 将结果转换为JSON字符串
	resultJSON, err := json.MarshalIndent(detailResp, "", "  ")
	if err != nil {
		errMsg := fmt.Sprintf("序列化组件详情结果失败: %v", err)
		utils.Error(errMsg)
		return &protocol.CallToolResult{
			Content: []protocol.Content{
				protocol.TextContent{
					Type: "text",
					Text: errMsg,
				},
			},
			IsError: true,
		}, nil
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

// handleCreateComponent 处理创建组件的请求
func (s *Service) handleCreateComponent(request *protocol.CallToolRequest) (*protocol.CallToolResult, error) {
	// 解析请求参数
	req := new(models.CreateComponentRequest)
	if err := protocol.VerifyAndUnmarshal(request.RawArguments, &req); err != nil {
		// 记录原始错误
		utils.Error("解析创建组件请求失败: %v", err)

		// 尝试解析原始请求数据
		var rawData map[string]interface{}
		var detailedErrMsg string
		if jsonErr := json.Unmarshal(request.RawArguments, &rawData); jsonErr == nil {
			// 检查必填字段
			requiredFields := []string{"team_name", "region_name", "group_id", "service_cname", "k8s_component_name", "image", "is_deploy"}
			var missingFields []string

			for _, field := range requiredFields {
				if _, exists := rawData[field]; !exists {
					missingFields = append(missingFields, field)
				}
			}

			// 检查字段类型
			var typeErrors []string
			if groupIDVal, ok := rawData["group_id"]; ok {
				switch groupIDVal.(type) {
				case float64: // JSON 解析数字会变成float64
					// 正确类型
				case string:
					typeErrors = append(typeErrors, "group_id 应为整数类型，但提供的是字符串")
				default:
					typeErrors = append(typeErrors, "group_id 类型错误，应为整数")
				}
			}

			if isDeployVal, ok := rawData["is_deploy"]; ok {
				switch isDeployVal.(type) {
				case bool:
					// 正确类型
				case string:
					typeErrors = append(typeErrors, "is_deploy 应为布尔类型，但提供的是字符串")
				default:
					typeErrors = append(typeErrors, "is_deploy 类型错误，应为布尔值(true/false)")
				}
			}

			// 构建详细错误信息
			if len(missingFields) > 0 {
				detailedErrMsg = fmt.Sprintf("缺少必填字段: %s", strings.Join(missingFields, ", "))
			} else if len(typeErrors) > 0 {
				detailedErrMsg = fmt.Sprintf("字段类型错误: %s", strings.Join(typeErrors, "; "))
			} else {
				detailedErrMsg = fmt.Sprintf("请求参数验证失败: %v", err)
			}
		} else {
			detailedErrMsg = fmt.Sprintf("解析JSON数据失败: %v", jsonErr)
		}

		// 返回带有详细错误信息的响应
		return &protocol.CallToolResult{
			Content: []protocol.Content{
				protocol.TextContent{
					Type: "text",
					Text: detailedErrMsg,
				},
			},
			IsError: true,
		}, nil
	}

	// 构建API路径 - 根据Rainbond OpenAPI文档
	path := fmt.Sprintf("/openapi/v1/teams/%s/regions/%s/apps/%d/services",
		req.TeamName, req.RegionName, req.GroupID)

	utils.Info("创建组件: %s", path)

	// 调用Rainbond API创建组件
	resp, err := s.client.Post(path, req)
	if err != nil {
		errMsg := fmt.Sprintf("创建组件失败: %v", err)
		utils.Error(errMsg)
		
		// 返回错误响应
		return &protocol.CallToolResult{
			Content: []protocol.Content{
				protocol.TextContent{
					Type: "text",
					Text: errMsg,
				},
			},
			IsError: true,
		}, nil
	}

	// 解析响应
	var result map[string]interface{}
	if err := json.Unmarshal(resp, &result); err != nil {
		errMsg := fmt.Sprintf("解析创建组件响应失败: %v", err)
		utils.Error(errMsg)
		
		// 返回错误响应
		return &protocol.CallToolResult{
			Content: []protocol.Content{
				protocol.TextContent{
					Type: "text",
					Text: errMsg,
				},
			},
			IsError: true,
		}, nil
	}

	// 将结果转换为JSON字符串
	resultJSON, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		errMsg := fmt.Sprintf("序列化创建组件结果失败: %v", err)
		utils.Error(errMsg)
		
		// 返回错误响应
		return &protocol.CallToolResult{
			Content: []protocol.Content{
				protocol.TextContent{
					Type: "text",
					Text: errMsg,
				},
			},
			IsError: true,
		}, nil
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

// handleBuildService 处理构建组件的请求
func (s *Service) handleBuildService(request *protocol.CallToolRequest) (*protocol.CallToolResult, error) {
	// 解析请求参数
	req := new(models.BuildComponentRequest)
	if err := protocol.VerifyAndUnmarshal(request.RawArguments, req); err != nil {
		utils.Error("解析构建组件请求失败: %v", err)
		return nil, fmt.Errorf("无效的组件构建请求: %v", err)
	}

	// 构建API路径 - 根据Rainbond OpenAPI文档
	path := fmt.Sprintf("/openapi/v1/teams/%s/regions/%s/apps/%s/services/%s/build",
		req.TeamName, req.RegionName, req.AppID, req.ServiceID)

	utils.Info("构建组件: %s", path)

	// 准备构建参数
	buildParams := map[string]interface{}{
		"is_deploy":  req.IsDeploy,
		"service_id": req.ServiceID,
	}

	// 如果有构建版本，添加到参数中
	if req.BuildVersion != "" {
		buildParams["build_version"] = req.BuildVersion
	}

	// 调用Rainbond API构建组件
	resp, err := s.client.Post(path, buildParams)
	if err != nil {
		utils.Error("构建组件失败: %v", err)
		return nil, fmt.Errorf("构建组件失败: %v", err)
	}

	// 解析响应
	var result map[string]interface{}
	if err := json.Unmarshal(resp, &result); err != nil {
		utils.Error("解析构建结果失败: %v", err)
		return nil, fmt.Errorf("解析构建结果失败: %v", err)
	}

	// 将构建结果转换为JSON字符串
	resultJSON, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		utils.Error("序列化构建结果失败: %v", err)
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

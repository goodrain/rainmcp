package components

import (
	"encoding/json"
	"fmt"
	"rainmcp/pkg/api"
	"rainmcp/pkg/logger"
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
	logger.Debug("创建新的组件服务")
	return &Service{
		client: client,
	}
}

// RegisterTools 注册组件相关的工具
func RegisterTools(mcpServer *server.Server, service *Service) {
	// 注册获取组件详情工具
	componentDetailTool, err := protocol.NewTool(
		"rainbond_get_component_detail",
		"获取Rainbond平台中的组件详情",
		models.ComponentDetailRequest{},
	)
	if err != nil {
		logger.Error("创建组件详情工具失败: %v", err)
		return
	}
	mcpServer.RegisterTool(componentDetailTool, service.handleGetComponentDetail)

	// 注册基于镜像创建组件工具
	createComponentTool, err := protocol.NewTool(
		"rainbond_create_image_component",
		"在Rainbond平台中基于镜像创建组件",
		models.CreateImageComponentRequest{},
	)
	if err != nil {
		logger.Error("创建组件创建工具失败: %v", err)
		return
	}
	mcpServer.RegisterTool(createComponentTool, service.handleCreateImageComponent)

	// 注册基于源码创建组件工具
	createCodeComponentTool, err := protocol.NewTool(
		"rainbond_create_code_component",
		"在Rainbond平台中基于源码创建组件",
		models.CreateCodeComponentRequest{},
	)
	if err != nil {
		logger.Error("创建基于源码的组件创建工具失败: %v", err)
		return
	}
	mcpServer.RegisterTool(createCodeComponentTool, service.handleCreateCodeComponent)

	// 注册获取组件端口列表工具
	listPortsTool, err := protocol.NewTool(
		"rainbond_list_component_ports",
		"获取Rainbond平台中组件的端口列表",
		models.ListPortsRequest{},
	)
	if err != nil {
		logger.Error("创建组件端口列表工具失败: %v", err)
		return
	}
	mcpServer.RegisterTool(listPortsTool, service.handleListComponentPorts)

	// 注册添加组件端口工具
	addPortTool, err := protocol.NewTool(
		"rainbond_add_component_port",
		"在Rainbond平台中添加组件端口",
		models.AddPortRequest{},
	)
	if err != nil {
		logger.Error("创建添加组件端口工具失败: %v", err)
		return
	}
	mcpServer.RegisterTool(addPortTool, service.handleAddComponentPort)

	// 注册更新组件端口工具
	updatePortTool, err := protocol.NewTool(
		"rainbond_update_component_port",
		"在Rainbond平台中更新组件端口",
		models.UpdatePortRequest{},
	)
	if err != nil {
		logger.Error("创建更新组件端口工具失败: %v", err)
		return
	}
	mcpServer.RegisterTool(updatePortTool, service.handleUpdateComponentPort)

	// 注册删除组件端口工具
	deletePortTool, err := protocol.NewTool(
		"rainbond_delete_component_port",
		"在Rainbond平台中删除组件端口",
		models.DeletePortRequest{},
	)
	if err != nil {
		logger.Error("创建删除组件端口工具失败: %v", err)
		return
	}
	mcpServer.RegisterTool(deletePortTool, service.handleDeleteComponentPort)

	// 注册构建组件工具
	buildServiceTool, err := protocol.NewTool(
		"rainbond_build_component",
		"在Rainbond平台中构建组件",
		models.BuildComponentRequest{},
	)
	if err != nil {
		logger.Error("创建组件构建工具失败: %v", err)
		return
	}
	mcpServer.RegisterTool(buildServiceTool, service.handleBuildService)

	// 注册获取应用下组件列表工具
	listComponentsTool, err := protocol.NewTool(
		"rainbond_list_components",
		"获取Rainbond平台中应用下的组件列表",
		models.ListComponentsRequest{},
	)
	if err != nil {
		logger.Error("创建组件列表工具失败: %v", err)
		return
	}
	mcpServer.RegisterTool(listComponentsTool, service.handleListComponents)
}

// handleListComponents 处理获取应用下组件列表的请求
func (s *Service) handleListComponents(request *protocol.CallToolRequest) (*protocol.CallToolResult, error) {
	// 解析请求参数
	req := new(models.ListComponentsRequest)
	if err := protocol.VerifyAndUnmarshal(request.RawArguments, req); err != nil {
		// 记录原始错误
		logger.Error("解析获取组件列表请求失败: %v", err)

		// 尝试解析原始请求数据
		var rawData map[string]interface{}
		var detailedErrMsg string
		if jsonErr := json.Unmarshal(request.RawArguments, &rawData); jsonErr == nil {
			// 检查必填字段
			requiredFields := []string{"team_name", "region_name", "app_id"}
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
	if req.TenantName == "" {
		errMsg := "缺少必填字段: team_name"
		logger.Error(errMsg)
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
		logger.Error(errMsg)
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
		logger.Error(errMsg)
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
	path := fmt.Sprintf("/openapi/v1/teams/%s/regions/%s/apps/%s/services",
		req.TenantName, req.RegionName, req.AppID)

	logger.Info("获取应用下组件列表: %s", path)

	// 调用Rainbond API获取组件列表
	resp, err := s.client.Get(path)
	if err != nil {
		errMsg := fmt.Sprintf("获取组件列表失败: %v", err)
		logger.Error(errMsg)
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
	var components []models.ComponentBaseInfo
	if err := json.Unmarshal(resp, &components); err != nil {
		errMsg := fmt.Sprintf("解析组件列表响应失败: %v", err)
		logger.Error(errMsg)
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
	resultJSON, err := json.MarshalIndent(components, "", "  ")
	if err != nil {
		errMsg := fmt.Sprintf("序列化组件列表结果失败: %v", err)
		logger.Error(errMsg)
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

// handleGetComponentDetail 处理获取组件详情的请求
func (s *Service) handleGetComponentDetail(request *protocol.CallToolRequest) (*protocol.CallToolResult, error) {
	// 解析请求参数
	req := new(models.ComponentDetailRequest)
	if err := protocol.VerifyAndUnmarshal(request.RawArguments, req); err != nil {
		// 记录原始错误
		logger.Error("解析获取组件详情请求失败: %v", err)

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
		logger.Error(errMsg)
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
		logger.Error(errMsg)
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
		logger.Error(errMsg)
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
		logger.Error(errMsg)
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

	logger.Info("获取组件详情: %s", path)

	// 调用Rainbond API获取组件详情
	resp, err := s.client.Get(path)
	if err != nil {
		errMsg := fmt.Sprintf("获取组件详情失败: %v", err)
		logger.Error(errMsg)
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
		logger.Error(errMsg)
		return &protocol.CallToolResult{
			Content: []protocol.Content{
				protocol.TextContent{
					Type: "text",
					Text: errMsg,
				},
			},
		}, nil
	}

	// 根据 service_source 字段区分展示源码组件和镜像组件的相关字段
	var result map[string]interface{}

	// 将组件详情转换为 map
	detailJSON, err := json.Marshal(detailResp)
	if err != nil {
		logger.Error("序列化组件详情失败: %v", err)
		return &protocol.CallToolResult{
			Content: []protocol.Content{
				protocol.TextContent{
					Type: "text",
					Text: fmt.Sprintf("序列化组件详情失败: %v", err),
				},
			},
			IsError: true,
		}, nil
	}

	if err := json.Unmarshal(detailJSON, &result); err != nil {
		logger.Error("反序列化组件详情失败: %v", err)
		return &protocol.CallToolResult{
			Content: []protocol.Content{
				protocol.TextContent{
					Type: "text",
					Text: fmt.Sprintf("反序列化组件详情失败: %v", err),
				},
			},
			IsError: true,
		}, nil
	}

	// 创建一个新的结果对象，包含基本字段
	formattedResult := map[string]interface{}{
		"基本信息": map[string]interface{}{
			"组件ID":  result["service_id"],
			"组件名称":  result["service_cname"],
			"组件英文名": result["k8s_component_name"],
			"所属集群":  result["service_region"],
			"运行状态":  result["status"],
			"内存配额":  fmt.Sprintf("%dMB", result["min_memory"]),
			"CPU配额": fmt.Sprintf("%d毫核", result["min_cpu"]),
			"伸缩方式":  result["extend_method"],
		},
	}

	// 根据 service_source 字段添加不同的字段
	serviceSource, _ := result["service_source"].(string)
	if serviceSource == "source_code" {
		// 源码组件特有字段
		formattedResult["源码信息"] = map[string]interface{}{
			"仓库地址": result["git_url"],
			"分支版本": result["code_version"],
		}
	} else if serviceSource == "docker_image" {
		// 镜像组件特有字段
		formattedResult["镜像信息"] = map[string]interface{}{
			"镜像地址":     result["image"],
			"启动命令":     result["cmd"],
			"Docker命令": result["docker_cmd"],
		}
	}

	// 添加访问地址信息
	accessInfos, ok := result["access_infos"].([]interface{})
	if ok && len(accessInfos) > 0 {
		addresses := make([]string, 0)
		for _, info := range accessInfos {
			if infoMap, ok := info.(map[string]interface{}); ok {
				if address, exists := infoMap["access_url"].(string); exists && address != "" {
					addresses = append(addresses, address)
				}
			}
		}
		
		if len(addresses) > 0 {
			formattedResult["访问信息"] = map[string]interface{}{
				"访问地址": addresses,
			}
		}
	}

	// 添加其他通用信息
	formattedResult["其他信息"] = map[string]interface{}{
		"架构": result["arch"],
	}

	// 将格式化的结果转换为JSON
	formattedJSON, err := json.MarshalIndent(formattedResult, "", "  ")
	if err != nil {
		logger.Error("格式化结果转换为JSON失败: %v", err)
		// 如果格式化失败，尝试使用原始数据
		formattedJSON, _ = json.MarshalIndent(detailResp, "", "  ")
	}

	// 返回结果
	return &protocol.CallToolResult{
		Content: []protocol.Content{
			protocol.TextContent{
				Type: "text",
				Text: string(formattedJSON),
			},
		},
	}, nil
}

// handleCreateImageComponent 处理基于镜像创建组件的请求
func (s *Service) handleCreateImageComponent(request *protocol.CallToolRequest) (*protocol.CallToolResult, error) {
	// 解析请求参数
	req := new(models.CreateImageComponentRequest)
	if err := protocol.VerifyAndUnmarshal(request.RawArguments, req); err != nil {
		// 记录原始错误
		logger.Error("解析基于镜像创建组件请求失败: %v", err)

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

	logger.Info("创建组件: %s", path)

	// 准备请求数据，始终将 is_deploy 设置为 true
	requestData := map[string]interface{}{
		"service_cname":      req.ServiceCName,
		"k8s_component_name": req.K8sComponentName,
		"image":              req.Image,
		"is_deploy":          true,
	}

	// 添加可选字段
	if req.DockerCmd != "" {
		requestData["docker_cmd"] = req.DockerCmd
	}

	if req.UserName != "" {
		requestData["user_name"] = req.UserName
	}

	if req.Password != "" {
		requestData["password"] = req.Password
	}

	// 调用Rainbond API创建组件
	resp, err := s.client.Post(path, requestData)
	if err != nil {
		errMsg := fmt.Sprintf("创建组件失败: %v", err)
		logger.Error(errMsg)

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
		logger.Error(errMsg)

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

	// 格式化结果
	resultJSON, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		errMsg := fmt.Sprintf("序列化创建组件结果失败: %v", err)
		logger.Error(errMsg)
		return nil, fmt.Errorf("序列化创建组件结果失败: %v", err)
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

// handleCreateCodeComponent 处理基于源码创建组件的请求
func (s *Service) handleCreateCodeComponent(request *protocol.CallToolRequest) (*protocol.CallToolResult, error) {
	// 解析请求参数
	req := new(models.CreateCodeComponentRequest)
	if err := protocol.VerifyAndUnmarshal(request.RawArguments, req); err != nil {
		// 记录原始错误
		logger.Error("解析基于源码创建组件请求失败: %v", err)

		// 尝试解析原始请求数据
		var rawData map[string]interface{}
		var detailedErrMsg string
		if jsonErr := json.Unmarshal(request.RawArguments, &rawData); jsonErr == nil {
			// 检查必填字段
			requiredFields := []string{"team_name", "region_name", "app_id", "service_cname", "k8s_component_name", "repo_url", "branch"}
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

	// 构建 API 路径 - 根据Rainbond OpenAPI文档
	path := fmt.Sprintf("/openapi/v1/teams/%s/regions/%s/apps/%s/code-services",
		req.TeamName, req.RegionName, req.AppID)

	logger.Info("基于源码创建组件: %s", path)

	// 准备请求数据
	requestData := map[string]interface{}{
		"service_cname":      req.ServiceCName,
		"k8s_component_name": req.K8sComponentName,
		"repo_url":           req.RepoURL,
		"branch":             req.Branch,
	}

	// 添加可选字段
	if req.Username != "" {
		requestData["username"] = req.Username
	}

	if req.Password != "" {
		requestData["password"] = req.Password
	}

	// 始终设置 is_deploy 为 true
	requestData["is_deploy"] = true

	// 调用Rainbond API创建组件
	resp, err := s.client.Post(path, requestData)
	if err != nil {
		errMsg := fmt.Sprintf("创建组件失败: %v", err)
		logger.Error(errMsg)

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
		logger.Error(errMsg)

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

	// 格式化结果
	resultJSON, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		errMsg := fmt.Sprintf("序列化创建组件结果失败: %v", err)
		logger.Error(errMsg)
		return nil, fmt.Errorf("序列化创建组件结果失败: %v", err)
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

// handleListComponentPorts 处理获取组件端口列表的请求
func (s *Service) handleListComponentPorts(request *protocol.CallToolRequest) (*protocol.CallToolResult, error) {
	// 解析请求参数
	req := new(models.ListPortsRequest)
	if err := protocol.VerifyAndUnmarshal(request.RawArguments, req); err != nil {
		logger.Error("解析获取组件端口列表请求失败: %v", err)
		return nil, fmt.Errorf("无效的获取组件端口列表请求: %v", err)
	}

	// 构建API路径
	apiPath := fmt.Sprintf("/openapi/v1/teams/%s/regions/%s/apps/%s/services/%s/ports",
		req.TenantID, req.RegionName, req.AppID, req.ServiceID)

	// 发送请求
	resp, err := s.client.Get(apiPath)
	if err != nil {
		errMsg := fmt.Sprintf("获取组件端口列表失败: %v", err)
		logger.Error(errMsg)

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
	logger.Debug("原始响应数据: %s", string(resp))

	// 使用PortListResponse结构体解析响应
	var portListResp models.PortListResponse
	if err := json.Unmarshal(resp, &portListResp); err != nil {
		logger.Warn("解析组件端口列表响应失败: %v", err)

		// 如果解析失败，尝试解析为通用JSON
		var result interface{}
		if err := json.Unmarshal(resp, &result); err != nil {
			errMsg := fmt.Sprintf("解析组件端口列表响应失败: %v", err)
			logger.Error(errMsg)
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

		// 将结果转换为格式化的JSON字符串
		resultJSON, err := json.MarshalIndent(result, "", "  ")
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

	// 成功解析为PortListResponse结构体
	logger.Info("成功解析组件端口列表响应，共有 %d 个端口", len(portListResp.Data.List))

	// 将结果转换为包含字段描述的JSON字符串
	resultJSON, err := utils.MarshalJSONWithDescription(portListResp)
	if err != nil {
		logger.Error("带描述的格式化响应数据失败: %v", err)
		// 如果格式化失败，尝试使用标准JSON序列化
		resultJSON, err = json.MarshalIndent(portListResp, "", "  ")
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

// handleAddComponentPort 处理添加组件端口的请求
func (s *Service) handleAddComponentPort(request *protocol.CallToolRequest) (*protocol.CallToolResult, error) {
	// 解析请求参数
	req := new(models.AddPortRequest)
	if err := protocol.VerifyAndUnmarshal(request.RawArguments, req); err != nil {
		logger.Error("解析添加组件端口请求失败: %v", err)
		return nil, fmt.Errorf("无效的添加组件端口请求: %v", err)
	}

	// 构建API路径
	apiPath := fmt.Sprintf("/openapi/v1/teams/%s/regions/%s/apps/%s/services/%s/ports",
		req.TenantID, req.RegionName, req.AppID, req.ServiceID)

	// 验证协议类型
	if req.Protocol != "tcp" && req.Protocol != "udp" && req.Protocol != "http" {
		errMsg := fmt.Sprintf("不支持的协议类型: %s，只支持 tcp/udp/http", req.Protocol)
		logger.Error(errMsg)
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

	// 准备请求体
	requestBody := map[string]interface{}{
		"port":             req.Port,
		"protocol":         req.Protocol,
		"is_inner_service": true,
		"is_outer_service": req.IsOuterService,
	}

	// 发送请求
	resp, err := s.client.Post(apiPath, requestBody)
	if err != nil {
		errMsg := fmt.Sprintf("添加组件端口失败: %v", err)
		logger.Error(errMsg)

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
	logger.Debug("原始响应数据: %s", string(resp))

	// 使用PortResponse结构体解析响应
	var portResp models.PortResponse
	if err := json.Unmarshal(resp, &portResp); err != nil {
		logger.Warn("解析添加组件端口响应失败: %v", err)

		// 如果解析失败，尝试解析为通用JSON
		var result interface{}
		if err := json.Unmarshal(resp, &result); err != nil {
			errMsg := fmt.Sprintf("解析添加组件端口响应失败: %v", err)
			logger.Error(errMsg)
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

		// 将结果转换为格式化的JSON字符串
		resultJSON, err := json.MarshalIndent(result, "", "  ")
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

	// 成功解析为PortResponse结构体
	logger.Info("成功解析添加端口响应，端口号: %d", portResp.Data.Bean.ContainerPort)

	// 将结果转换为包含字段描述的JSON字符串
	resultJSON, err := utils.MarshalJSONWithDescription(portResp)
	if err != nil {
		logger.Error("带描述的格式化响应数据失败: %v", err)
		// 如果格式化失败，尝试使用标准JSON序列化
		resultJSON, err = json.MarshalIndent(portResp, "", "  ")
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

// handleUpdateComponentPort 处理更新组件端口的请求
func (s *Service) handleUpdateComponentPort(request *protocol.CallToolRequest) (*protocol.CallToolResult, error) {
	// 解析请求参数
	req := new(models.UpdatePortRequest)
	if err := protocol.VerifyAndUnmarshal(request.RawArguments, req); err != nil {
		logger.Error("解析更新组件端口请求失败: %v", err)
		return nil, fmt.Errorf("无效的更新组件端口请求: %v", err)
	}

	// 构建API路径
	apiPath := fmt.Sprintf("/openapi/v1/teams/%s/regions/%s/apps/%s/services/%s/ports/%d",
		req.TenantID, req.RegionName, req.AppID, req.ServiceID, req.Port)

	// 准备请求体
	requestBody := map[string]interface{}{
		"action": req.Action,
	}

	// 根据不同的操作类型添加相应的参数
	switch req.Action {
	case "open_outer", "close_outer", "open_inner", "close_inner":
		// 这些操作只需要action参数，不需要额外的参数
		logger.Debug("执行端口操作: %s", req.Action)
	case "change_protocol":
		if req.Protocol == "" {
			errMsg := "更改协议时必须指定协议类型"
			logger.Error(errMsg)
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

		// 验证协议类型
		if req.Protocol != "tcp" && req.Protocol != "udp" && req.Protocol != "http" {
			errMsg := fmt.Sprintf("不支持的协议类型: %s，只支持 tcp/udp/http", req.Protocol)
			logger.Error(errMsg)
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

		requestBody["protocol"] = req.Protocol
		logger.Debug("更改端口协议为: %s", req.Protocol)
	default:
		errMsg := fmt.Sprintf("不支持的操作类型: %s", req.Action)
		logger.Error(errMsg)
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

	// 发送请求
	resp, err := s.client.Put(apiPath, requestBody)
	if err != nil {
		errMsg := fmt.Sprintf("更新组件端口失败: %v", err)
		logger.Error(errMsg)

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
	logger.Debug("原始响应数据: %s", string(resp))

	// 使用PortResponse结构体解析响应
	var portResp models.PortResponse
	if err := json.Unmarshal(resp, &portResp); err != nil {
		logger.Warn("解析更新组件端口响应失败: %v", err)

		// 如果解析失败，尝试解析为通用JSON
		var result interface{}
		if err := json.Unmarshal(resp, &result); err != nil {
			errMsg := fmt.Sprintf("解析更新组件端口响应失败: %v", err)
			logger.Error(errMsg)
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

		// 将结果转换为格式化的JSON字符串
		resultJSON, err := json.MarshalIndent(result, "", "  ")
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

	// 成功解析为PortResponse结构体
	logger.Info("成功解析更新端口响应，端口号: %d", portResp.Data.Bean.ContainerPort)

	// 将结果转换为包含字段描述的JSON字符串
	resultJSON, err := utils.MarshalJSONWithDescription(portResp)
	if err != nil {
		logger.Error("带描述的格式化响应数据失败: %v", err)
		// 如果格式化失败，尝试使用标准JSON序列化
		resultJSON, err = json.MarshalIndent(portResp, "", "  ")
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

// handleDeleteComponentPort 处理删除组件端口的请求
func (s *Service) handleDeleteComponentPort(request *protocol.CallToolRequest) (*protocol.CallToolResult, error) {
	// 解析请求参数
	req := new(models.DeletePortRequest)
	if err := protocol.VerifyAndUnmarshal(request.RawArguments, req); err != nil {
		logger.Error("解析删除组件端口请求失败: %v", err)
		return nil, fmt.Errorf("无效的删除组件端口请求: %v", err)
	}

	// 构建API路径
	apiPath := fmt.Sprintf("/openapi/v1/teams/%s/regions/%s/apps/%s/services/%s/ports/%d",
		req.TenantID, req.RegionName, req.AppID, req.ServiceID, req.Port)

	// 发送请求
	resp, err := s.client.Delete(apiPath)
	if err != nil {
		errMsg := fmt.Sprintf("删除组件端口失败: %v", err)
		logger.Error(errMsg)

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
	logger.Debug("原始响应数据: %s", string(resp))

	// 尝试解析为通用JSON
	var result interface{}
	if err := json.Unmarshal(resp, &result); err != nil {
		errMsg := fmt.Sprintf("解析删除组件端口响应失败: %v", err)
		logger.Error(errMsg)
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

	logger.Info("成功解析删除端口响应")

	// 格式化结果
	resultJSON, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		errMsg := fmt.Sprintf("序列化删除组件端口结果失败: %v", err)
		logger.Error(errMsg)
		return nil, fmt.Errorf("序列化删除组件端口结果失败: %v", err)
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
		logger.Error("解析构建组件请求失败: %v", err)
		return nil, fmt.Errorf("无效的组件构建请求: %v", err)
	}

	// 构建API路径 - 根据Rainbond OpenAPI文档
	path := fmt.Sprintf("/openapi/v1/teams/%s/regions/%s/apps/%s/services/%s/build",
		req.TeamName, req.RegionName, req.AppID, req.ServiceID)

	logger.Info("构建组件: %s", path)

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
		logger.Error("构建组件失败: %v", err)
		return nil, fmt.Errorf("构建组件失败: %v", err)
	}

	// 解析响应
	var result map[string]interface{}
	if err := json.Unmarshal(resp, &result); err != nil {
		logger.Error("解析构建结果失败: %v", err)
		return nil, fmt.Errorf("解析构建结果失败: %v", err)
	}

	// 将构建结果转换为JSON字符串
	resultJSON, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		logger.Error("序列化构建结果失败: %v", err)
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

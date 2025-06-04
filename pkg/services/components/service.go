package components

import (
	"context"
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

	//// 注册更新组件端口工具
	//updatePortTool, err := protocol.NewTool(
	//	"rainbond_update_component_port",
	//	"在Rainbond平台中更新组件端口",
	//	models.UpdatePortRequest{},
	//)
	//if err != nil {
	//	logger.Error("创建更新组件端口工具失败: %v", err)
	//	return
	//}
	//mcpServer.RegisterTool(updatePortTool, service.handleUpdateComponentPort)
	//
	//// 注册删除组件端口工具
	//deletePortTool, err := protocol.NewTool(
	//	"rainbond_delete_component_port",
	//	"在Rainbond平台中删除组件端口",
	//	models.DeletePortRequest{},
	//)
	//if err != nil {
	//	logger.Error("创建删除组件端口工具失败: %v", err)
	//	return
	//}
	//mcpServer.RegisterTool(deletePortTool, service.handleDeleteComponentPort)
	//
	//// 注册构建组件工具
	//buildServiceTool, err := protocol.NewTool(
	//	"rainbond_build_component",
	//	"在Rainbond平台中构建组件",
	//	models.BuildComponentRequest{},
	//)
	//if err != nil {
	//	logger.Error("创建组件构建工具失败: %v", err)
	//	return
	//}
	//mcpServer.RegisterTool(buildServiceTool, service.handleBuildService)

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
func (service *Service) handleListComponents(ctx context.Context, request *protocol.CallToolRequest) (*protocol.CallToolResult, error) {
	// 解析请求参数
	rainTokenValue := ctx.Value(models.RainTokenKey{})
	rainToken := rainTokenValue.(string)
	service.client.Token = rainToken

	req := new(models.ListComponentsRequest)
	if err := protocol.VerifyAndUnmarshal(request.RawArguments, req); err != nil {
		// 记录原始错误
		logger.Error("解析获取组件列表请求失败: %v", err)

		// 尝试解析原始请求数据
		var rawData map[string]interface{}
		var detailedErrMsg string
		if jsonErr := json.Unmarshal(request.RawArguments, &rawData); jsonErr == nil {
			// 检查必填字段
			requiredFields := []string{"team_alias", "app_id"}
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
				&protocol.TextContent{
					Type: "text",
					Text: detailedErrMsg,
				},
			},
			IsError: true,
		}, nil
	}

	// 参数校验
	if req.TeamAlias == "" {
		errMsg := "缺少必填字段: team_alias"
		logger.Error(errMsg)
		return &protocol.CallToolResult{
			Content: []protocol.Content{
				&protocol.TextContent{
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
				&protocol.TextContent{
					Type: "text",
					Text: errMsg,
				},
			},
			IsError: true,
		}, nil
	}

	// 构建API路径 - 根据Rainbond OpenAPI文档
	path := fmt.Sprintf("/openapi/v1/mcp/teams/%s/apps/%s/components",
		req.TeamAlias, req.AppID)

	logger.Info("获取应用下组件列表: %s", path)

	// 调用Rainbond API获取组件列表
	resp, err := service.client.Get(path)
	if err != nil {
		errMsg := fmt.Sprintf("获取组件列表失败: %v", err)
		logger.Error(errMsg)
		return &protocol.CallToolResult{
			Content: []protocol.Content{
				&protocol.TextContent{
					Type: "text",
					Text: errMsg,
				},
			},
			IsError: true,
		}, nil
	}

	// 使用ComponentListResponse结构体解析响应
	var componentsResp models.ComponentListResponse
	logger.Debug("原始响应数据: %s", string(resp))

	if err := json.Unmarshal(resp, &componentsResp); err != nil {
		logger.Warn("解析组件列表响应失败: %v", err)

		// 如果解析失败，尝试解析为通用JSON
		var result interface{}
		if err := json.Unmarshal(resp, &result); err != nil {
			errMsg := fmt.Sprintf("解析组件列表响应失败: %v", err)
			logger.Error(errMsg)
			return &protocol.CallToolResult{
				Content: []protocol.Content{
					&protocol.TextContent{
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
				&protocol.TextContent{
					Type: "text",
					Text: string(resultJSON),
				},
			},
		}, nil
	}

	// 成功解析为ComponentListResponse结构体
	logger.Info("成功解析组件列表数据，共有 %d 个组件", len(componentsResp.Data.List))

	// 将结果转换为包含字段描述的JSON字符串
	resultJSON, err := utils.MarshalJSONWithDescription(componentsResp)
	if err != nil {
		logger.Error("格式化响应数据失败: %v", err)
		// 如果格式化失败，尝试使用标准JSON序列化
		resultJSON, err = json.MarshalIndent(componentsResp, "", "  ")
		if err != nil {
			logger.Error("标准JSON格式化也失败: %v", err)
			// 如果标准格式化也失败，直接返回原始数据
			resultJSON = resp
		}
	}

	// 返回结果
	return &protocol.CallToolResult{
		Content: []protocol.Content{
			&protocol.TextContent{
				Type: "text",
				Text: string(resultJSON),
			},
		},
	}, nil
}

// handleGetComponentDetail 处理获取组件详情的请求
func (service *Service) handleGetComponentDetail(ctx context.Context, request *protocol.CallToolRequest) (*protocol.CallToolResult, error) {
	// 解析请求参数
	rainTokenValue := ctx.Value(models.RainTokenKey{})
	rainToken := rainTokenValue.(string)
	service.client.Token = rainToken

	req := new(models.ComponentDetailRequest)
	if err := protocol.VerifyAndUnmarshal(request.RawArguments, req); err != nil {
		// 记录原始错误
		logger.Error("解析获取组件详情请求失败: %v", err)

		// 尝试解析原始请求数据
		var rawData map[string]interface{}
		var detailedErrMsg string
		if jsonErr := json.Unmarshal(request.RawArguments, &rawData); jsonErr == nil {
			// 检查必填字段
			requiredFields := []string{"team_alias", "app_id", "service_id"}
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
				&protocol.TextContent{
					Type: "text",
					Text: detailedErrMsg,
				},
			},
			IsError: true,
		}, nil
	}

	// 参数校验
	if req.TeamAlias == "" {
		errMsg := "缺少必填字段: team_alias"
		logger.Error(errMsg)
		return &protocol.CallToolResult{
			Content: []protocol.Content{
				&protocol.TextContent{
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
				&protocol.TextContent{
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
				&protocol.TextContent{
					Type: "text",
					Text: errMsg,
				},
			},
			IsError: true,
		}, nil
	}

	// 构建API路径 - 根据Rainbond OpenAPI文档
	path := fmt.Sprintf("/openapi/v1/mcp/teams/%s/apps/%s/components/%s",
		req.TeamAlias, req.AppID, req.ServiceID)

	logger.Info("获取组件详情: %s", path)

	// 调用Rainbond API获取组件详情
	resp, err := service.client.Get(path)
	if err != nil {
		errMsg := fmt.Sprintf("获取组件详情失败: %v", err)
		logger.Error(errMsg)
		return &protocol.CallToolResult{
			Content: []protocol.Content{
				&protocol.TextContent{
					Type: "text",
					Text: errMsg,
				},
			},
			IsError: true,
		}, nil
	}

	// 使用NewComponentDetailResponse结构体解析响应
	var detailResp models.NewComponentDetailResponse
	logger.Debug("原始响应数据: %s", string(resp))

	if err := json.Unmarshal(resp, &detailResp); err != nil {
		logger.Warn("解析组件详情响应失败: %v", err)

		// 如果解析失败，尝试解析为通用JSON
		var result interface{}
		if err := json.Unmarshal(resp, &result); err != nil {
			errMsg := fmt.Sprintf("解析组件详情响应失败: %v", err)
			logger.Error(errMsg)
			return &protocol.CallToolResult{
				Content: []protocol.Content{
					&protocol.TextContent{
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
				&protocol.TextContent{
					Type: "text",
					Text: string(resultJSON),
				},
			},
		}, nil
	}

	// 成功解析为NewComponentDetailResponse结构体
	logger.Info("成功解析组件详情数据，组件名称: %s", detailResp.Data.Bean.ServiceCName)

	// 创建格式化的结果对象
	formattedResult := map[string]interface{}{
		"基本信息": map[string]interface{}{
			"组件ID":  detailResp.Data.Bean.ServiceID,
			"组件名称":  detailResp.Data.Bean.ServiceCName,
			"组件别名":  detailResp.Data.Bean.ServiceAlias,
			"运行状态":  detailResp.Data.Bean.StatusCN,
			"更新时间":  detailResp.Data.Bean.UpdateTime,
			"内存配额":  fmt.Sprintf("%dMB", detailResp.Data.Bean.MinMemory),
			"CPU配额": fmt.Sprintf("%d毫核", detailResp.Data.Bean.MinCPU),
		},
	}

	// 添加端口信息
	if len(detailResp.Data.Bean.Ports) > 0 {
		ports := make([]map[string]interface{}, 0)
		for _, port := range detailResp.Data.Bean.Ports {
			portInfo := map[string]interface{}{
				"端口号":  port.ContainerPort,
				"协议":   port.Protocol,
				"对外服务": port.IsOuterService,
				"对内服务": port.IsInnerService,
			}
			if len(port.AccessUrls) > 0 {
				portInfo["访问地址"] = port.AccessUrls
			}
			ports = append(ports, portInfo)
		}
		formattedResult["端口信息"] = ports
	}

	// 添加环境变量信息
	if len(detailResp.Data.Bean.Envs) > 0 {
		envs := make([]map[string]interface{}, 0)
		for _, env := range detailResp.Data.Bean.Envs {
			envInfo := map[string]interface{}{
				"变量名": env.AttrName,
				"变量值": env.AttrValue,
				"作用域": env.Scope,
				"可更改": env.IsChange,
			}
			if env.Name != "" {
				envInfo["显示名称"] = env.Name
			}
			envs = append(envs, envInfo)
		}
		formattedResult["环境变量"] = envs
	}

	// 添加存储卷信息
	if len(detailResp.Data.Bean.Volumes) > 0 {
		volumes := make([]map[string]interface{}, 0)
		for _, volume := range detailResp.Data.Bean.Volumes {
			volumeInfo := map[string]interface{}{
				"存储卷名称": volume.VolumeName,
				"挂载路径":  volume.VolumePath,
				"存储容量":  fmt.Sprintf("%dGB", volume.VolumeCapacity),
			}
			volumes = append(volumes, volumeInfo)
		}
		formattedResult["存储卷"] = volumes
	}

	// 将格式化的结果转换为JSON
	resultJSON, err := json.MarshalIndent(formattedResult, "", "  ")
	if err != nil {
		logger.Error("格式化结果转换为JSON失败: %v", err)
		// 如果格式化失败，尝试使用包含字段描述的JSON序列化
		resultJSON, err = utils.MarshalJSONWithDescription(detailResp)
		if err != nil {
			logger.Error("带描述的格式化也失败: %v", err)
			// 如果格式化失败，尝试使用标准JSON序列化
			resultJSON, err = json.MarshalIndent(detailResp, "", "  ")
			if err != nil {
				logger.Error("标准JSON格式化也失败: %v", err)
				// 如果标准格式化也失败，直接返回原始数据
				resultJSON = resp
			}
		}
	}

	// 返回结果
	return &protocol.CallToolResult{
		Content: []protocol.Content{
			&protocol.TextContent{
				Type: "text",
				Text: string(resultJSON),
			},
		},
	}, nil
}

// handleCreateCodeComponent 处理基于源码创建组件的请求
func (service *Service) handleCreateCodeComponent(ctx context.Context, request *protocol.CallToolRequest) (*protocol.CallToolResult, error) {
	rainTokenValue := ctx.Value(models.RainTokenKey{})
	rainToken := rainTokenValue.(string)
	service.client.Token = rainToken
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
				&protocol.TextContent{
					Type: "text",
					Text: detailedErrMsg,
				},
			},
			IsError: true,
		}, nil
	}

	// 构建 API 路径 - 根据Rainbond OpenAPI文档
	path := fmt.Sprintf("/openapi/v1/mcp/teams/%s/apps/%s/components/create",
		req.TeamAlias, req.AppID)

	logger.Info("基于源码创建组件: %s", path)

	// 准备请求数据
	requestData := map[string]interface{}{
		"service_cname": req.ServiceCName,
		"repo_url":      req.RepoURL,
		"branch":        req.Branch,
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
	resp, err := service.client.Post(path, requestData)
	if err != nil {
		errMsg := fmt.Sprintf("创建组件失败: %v", err)
		logger.Error(errMsg)

		// 返回错误响应
		return &protocol.CallToolResult{
			Content: []protocol.Content{
				&protocol.TextContent{
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
				&protocol.TextContent{
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
			&protocol.TextContent{
				Type: "text",
				Text: string(resultJSON),
			},
		},
	}, nil
}

// handleListComponentPorts 处理获取组件端口列表的请求
func (service *Service) handleListComponentPorts(ctx context.Context, request *protocol.CallToolRequest) (*protocol.CallToolResult, error) {
	rainTokenValue := ctx.Value(models.RainTokenKey{})
	rainToken := rainTokenValue.(string)
	service.client.Token = rainToken
	// 解析请求参数
	req := new(models.ListPortsRequest)
	if err := protocol.VerifyAndUnmarshal(request.RawArguments, req); err != nil {
		// 记录原始错误
		logger.Error("解析获取组件端口列表请求失败: %v", err)

		// 尝试解析原始请求数据
		var rawData map[string]interface{}
		var detailedErrMsg string
		if jsonErr := json.Unmarshal(request.RawArguments, &rawData); jsonErr == nil {
			// 检查必填字段
			requiredFields := []string{"team_alias", "app_id", "service_id"}
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
				&protocol.TextContent{
					Type: "text",
					Text: detailedErrMsg,
				},
			},
			IsError: true,
		}, nil
	}

	// 参数校验
	if req.TeamAlias == "" {
		errMsg := "缺少必填字段: team_alias"
		logger.Error(errMsg)
		return &protocol.CallToolResult{
			Content: []protocol.Content{
				&protocol.TextContent{
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
				&protocol.TextContent{
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
				&protocol.TextContent{
					Type: "text",
					Text: errMsg,
				},
			},
			IsError: true,
		}, nil
	}

	// 构建API路径
	apiPath := fmt.Sprintf("/openapi/v1/mcp/teams/%s/apps/%s/components/%s/ports",
		req.TeamAlias, req.AppID, req.ServiceID)

	logger.Info("获取组件端口列表: %s", apiPath)

	// 发送请求
	resp, err := service.client.Get(apiPath)
	if err != nil {
		errMsg := fmt.Sprintf("获取组件端口列表失败: %v", err)
		logger.Error(errMsg)

		// 返回错误响应
		return &protocol.CallToolResult{
			Content: []protocol.Content{
				&protocol.TextContent{
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
					&protocol.TextContent{
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
				&protocol.TextContent{
					Type: "text",
					Text: string(resultJSON),
				},
			},
		}, nil
	}

	// 成功解析为PortListResponse结构体
	logger.Info("成功解析组件端口列表响应，共有 %d 个端口", len(portListResp.Data.List))

	// 创建格式化的结果对象
	formattedResult := map[string]interface{}{
		"端口列表": make([]map[string]interface{}, 0),
	}

	// 添加端口信息
	if len(portListResp.Data.List) > 0 {
		ports := make([]map[string]interface{}, 0)
		for _, port := range portListResp.Data.List {
			portInfo := map[string]interface{}{
				"端口号":  port.Port,
				"协议":   port.Protocol,
				"对外服务": port.IsOuterService,
				"对内服务": port.IsInnerService,
			}
			ports = append(ports, portInfo)
		}
		formattedResult["端口列表"] = ports
	}

	// 将结果转换为包含字段描述的JSON字符串
	resultJSON, err := json.MarshalIndent(formattedResult, "", "  ")
	if err != nil {
		logger.Error("格式化结果转换为JSON失败: %v", err)
		// 如果格式化失败，尝试使用包含字段描述的JSON序列化
		resultJSON, err = utils.MarshalJSONWithDescription(portListResp)
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
	}

	// 返回结果
	return &protocol.CallToolResult{
		Content: []protocol.Content{
			&protocol.TextContent{
				Type: "text",
				Text: string(resultJSON),
			},
		},
	}, nil
}

// handleAddComponentPort 处理添加组件端口的请求
func (service *Service) handleAddComponentPort(ctx context.Context, request *protocol.CallToolRequest) (*protocol.CallToolResult, error) {
	rainTokenValue := ctx.Value(models.RainTokenKey{})
	rainToken := rainTokenValue.(string)
	service.client.Token = rainToken
	// 解析请求参数
	req := new(models.AddPortRequest)
	if err := protocol.VerifyAndUnmarshal(request.RawArguments, req); err != nil {
		logger.Error("解析添加组件端口请求失败: %v", err)
		return nil, fmt.Errorf("无效的添加组件端口请求: %v", err)
	}

	// 构建API路径
	apiPath := fmt.Sprintf("/openapi/v1/mcp/teams/%s/apps/%s/components/%s/ports",
		req.TeamAlias, req.AppID, req.ServiceID)

	// 验证协议类型
	if req.Protocol != "tcp" && req.Protocol != "udp" && req.Protocol != "http" {
		errMsg := fmt.Sprintf("不支持的协议类型: %s，只支持 tcp/udp/http", req.Protocol)
		logger.Error(errMsg)
		return &protocol.CallToolResult{
			Content: []protocol.Content{
				&protocol.TextContent{
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
		"is_outer_service": req.IsOuterService,
	}

	// 发送请求
	resp, err := service.client.Post(apiPath, requestBody)
	if err != nil {
		errMsg := fmt.Sprintf("添加组件端口失败: %v", err)
		logger.Error(errMsg)

		// 返回错误响应
		return &protocol.CallToolResult{
			Content: []protocol.Content{
				&protocol.TextContent{
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
					&protocol.TextContent{
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
				&protocol.TextContent{
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
			&protocol.TextContent{
				Type: "text",
				Text: string(resultJSON),
			},
		},
	}, nil
}

//// handleUpdateComponentPort 处理更新组件端口的请求
//func (service *Service) handleUpdateComponentPort(ctx context.Context, request *protocol.CallToolRequest) (*protocol.CallToolResult, error) {
//	// 解析请求参数
//	req := new(models.UpdatePortRequest)
//	if err := protocol.VerifyAndUnmarshal(request.RawArguments, req); err != nil {
//		logger.Error("解析更新组件端口请求失败: %v", err)
//		return nil, fmt.Errorf("无效的更新组件端口请求: %v", err)
//	}
//
//	// 构建API路径
//	apiPath := fmt.Sprintf("/openapi/v1/teams/%s/regions/%s/apps/%s/services/%s/ports/%d",
//		req.TenantID, req.RegionName, req.AppID, req.ServiceID, req.Port)
//
//	// 准备请求体
//	requestBody := map[string]interface{}{
//		"action": req.Action,
//	}
//
//	// 根据不同的操作类型添加相应的参数
//	switch req.Action {
//	case "open_outer", "close_outer", "open_inner", "close_inner":
//		// 这些操作只需要action参数，不需要额外的参数
//		logger.Debug("执行端口操作: %s", req.Action)
//	case "change_protocol":
//		if req.Protocol == "" {
//			errMsg := "更改协议时必须指定协议类型"
//			logger.Error(errMsg)
//			return &protocol.CallToolResult{
//				Content: []protocol.Content{
//					&protocol.TextContent{
//						Type: "text",
//						Text: errMsg,
//					},
//				},
//				IsError: true,
//			}, nil
//		}
//
//		// 验证协议类型
//		if req.Protocol != "tcp" && req.Protocol != "udp" && req.Protocol != "http" {
//			errMsg := fmt.Sprintf("不支持的协议类型: %s，只支持 tcp/udp/http", req.Protocol)
//			logger.Error(errMsg)
//			return &protocol.CallToolResult{
//				Content: []protocol.Content{
//					&protocol.TextContent{
//						Type: "text",
//						Text: errMsg,
//					},
//				},
//				IsError: true,
//			}, nil
//		}
//
//		requestBody["protocol"] = req.Protocol
//		logger.Debug("更改端口协议为: %s", req.Protocol)
//	default:
//		errMsg := fmt.Sprintf("不支持的操作类型: %s", req.Action)
//		logger.Error(errMsg)
//		return &protocol.CallToolResult{
//			Content: []protocol.Content{
//				&protocol.TextContent{
//					Type: "text",
//					Text: errMsg,
//				},
//			},
//			IsError: true,
//		}, nil
//	}
//	rainTokenValue := ctx.Value(models.RainTokenKey{})
//	rainToken := rainTokenValue.(string)
//	service.client.Token = rainToken
//	// 发送请求
//	resp, err := service.client.Put(apiPath, requestBody)
//	if err != nil {
//		errMsg := fmt.Sprintf("更新组件端口失败: %v", err)
//		logger.Error(errMsg)
//
//		// 返回错误响应
//		return &protocol.CallToolResult{
//			Content: []protocol.Content{
//				&protocol.TextContent{
//					Type: "text",
//					Text: errMsg,
//				},
//			},
//			IsError: true,
//		}, nil
//	}
//
//	// 解析响应
//	logger.Debug("原始响应数据: %s", string(resp))
//
//	// 使用PortResponse结构体解析响应
//	var portResp models.PortResponse
//	if err := json.Unmarshal(resp, &portResp); err != nil {
//		logger.Warn("解析更新组件端口响应失败: %v", err)
//
//		// 如果解析失败，尝试解析为通用JSON
//		var result interface{}
//		if err := json.Unmarshal(resp, &result); err != nil {
//			errMsg := fmt.Sprintf("解析更新组件端口响应失败: %v", err)
//			logger.Error(errMsg)
//			return &protocol.CallToolResult{
//				Content: []protocol.Content{
//					&protocol.TextContent{
//						Type: "text",
//						Text: errMsg,
//					},
//				},
//				IsError: true,
//			}, nil
//		}
//
//		// 将结果转换为格式化的JSON字符串
//		resultJSON, err := json.MarshalIndent(result, "", "  ")
//		if err != nil {
//			logger.Error("格式化响应数据失败: %v", err)
//			// 如果格式化失败，直接返回原始数据
//			resultJSON = resp
//		}
//
//		return &protocol.CallToolResult{
//			Content: []protocol.Content{
//				&protocol.TextContent{
//					Type: "text",
//					Text: string(resultJSON),
//				},
//			},
//		}, nil
//	}
//
//	// 成功解析为PortResponse结构体
//	logger.Info("成功解析更新端口响应，端口号: %d", portResp.Data.Bean.ContainerPort)
//
//	// 将结果转换为包含字段描述的JSON字符串
//	resultJSON, err := utils.MarshalJSONWithDescription(portResp)
//	if err != nil {
//		logger.Error("带描述的格式化响应数据失败: %v", err)
//		// 如果格式化失败，尝试使用标准JSON序列化
//		resultJSON, err = json.MarshalIndent(portResp, "", "  ")
//		if err != nil {
//			logger.Error("标准JSON格式化也失败: %v", err)
//			// 如果标准格式化也失败，直接返回原始数据
//			resultJSON = resp
//		}
//	}
//
//	// 返回结果
//	return &protocol.CallToolResult{
//		Content: []protocol.Content{
//			&protocol.TextContent{
//				Type: "text",
//				Text: string(resultJSON),
//			},
//		},
//	}, nil
//}
//
//// handleDeleteComponentPort 处理删除组件端口的请求
//func (service *Service) handleDeleteComponentPort(ctx context.Context, request *protocol.CallToolRequest) (*protocol.CallToolResult, error) {
//	// 解析请求参数
//	req := new(models.DeletePortRequest)
//	if err := protocol.VerifyAndUnmarshal(request.RawArguments, req); err != nil {
//		logger.Error("解析删除组件端口请求失败: %v", err)
//		return nil, fmt.Errorf("无效的删除组件端口请求: %v", err)
//	}
//
//	// 构建API路径
//	apiPath := fmt.Sprintf("/openapi/v1/teams/%s/regions/%s/apps/%s/services/%s/ports/%d",
//		req.TenantID, req.RegionName, req.AppID, req.ServiceID, req.Port)
//
//	// 发送请求
//	resp, err := service.client.Delete(apiPath)
//	if err != nil {
//		errMsg := fmt.Sprintf("删除组件端口失败: %v", err)
//		logger.Error(errMsg)
//
//		// 返回错误响应
//		return &protocol.CallToolResult{
//			Content: []protocol.Content{
//				&protocol.TextContent{
//					Type: "text",
//					Text: errMsg,
//				},
//			},
//			IsError: true,
//		}, nil
//	}
//
//	// 解析响应
//	logger.Debug("原始响应数据: %s", string(resp))
//
//	// 尝试解析为通用JSON
//	var result interface{}
//	if err := json.Unmarshal(resp, &result); err != nil {
//		errMsg := fmt.Sprintf("解析删除组件端口响应失败: %v", err)
//		logger.Error(errMsg)
//		return &protocol.CallToolResult{
//			Content: []protocol.Content{
//				&protocol.TextContent{
//					Type: "text",
//					Text: errMsg,
//				},
//			},
//			IsError: true,
//		}, nil
//	}
//
//	logger.Info("成功解析删除端口响应")
//
//	// 格式化结果
//	resultJSON, err := json.MarshalIndent(result, "", "  ")
//	if err != nil {
//		errMsg := fmt.Sprintf("序列化删除组件端口结果失败: %v", err)
//		logger.Error(errMsg)
//		return nil, fmt.Errorf("序列化删除组件端口结果失败: %v", err)
//	}
//
//	// 返回结果
//	return &protocol.CallToolResult{
//		Content: []protocol.Content{
//			&protocol.TextContent{
//				Type: "text",
//				Text: string(resultJSON),
//			},
//		},
//	}, nil
//}
//
//// handleBuildService 处理构建组件的请求
//func (service *Service) handleBuildService(ctx context.Context, request *protocol.CallToolRequest) (*protocol.CallToolResult, error) {
//	// 解析请求参数
//	req := new(models.BuildComponentRequest)
//	if err := protocol.VerifyAndUnmarshal(request.RawArguments, req); err != nil {
//		logger.Error("解析构建组件请求失败: %v", err)
//		return nil, fmt.Errorf("无效的组件构建请求: %v", err)
//	}
//
//	// 构建API路径 - 根据Rainbond OpenAPI文档
//	path := fmt.Sprintf("/openapi/v1/teams/%s/regions/%s/apps/%s/services/%s/build",
//		req.TeamAlias, req.RegionName, req.AppID, req.ServiceID)
//
//	logger.Info("构建组件: %s", path)
//
//	// 准备构建参数
//	buildParams := map[string]interface{}{
//		"is_deploy":  req.IsDeploy,
//		"service_id": req.ServiceID,
//	}
//
//	// 如果有构建版本，添加到参数中
//	if req.BuildVersion != "" {
//		buildParams["build_version"] = req.BuildVersion
//	}
//
//	// 调用Rainbond API构建组件
//	resp, err := service.client.Post(path, buildParams)
//	if err != nil {
//		logger.Error("构建组件失败: %v", err)
//		return nil, fmt.Errorf("构建组件失败: %v", err)
//	}
//
//	// 解析响应
//	var result map[string]interface{}
//	if err := json.Unmarshal(resp, &result); err != nil {
//		logger.Error("解析构建结果失败: %v", err)
//		return nil, fmt.Errorf("解析构建结果失败: %v", err)
//	}
//
//	// 将构建结果转换为JSON字符串
//	resultJSON, err := json.MarshalIndent(result, "", "  ")
//	if err != nil {
//		logger.Error("序列化构建结果失败: %v", err)
//		return nil, fmt.Errorf("序列化构建结果失败: %v", err)
//	}
//
//	// 返回结果
//	return &protocol.CallToolResult{
//		Content: []protocol.Content{
//			&protocol.TextContent{
//				Type: "text",
//				Text: string(resultJSON),
//			},
//		},
//	}, nil
//}

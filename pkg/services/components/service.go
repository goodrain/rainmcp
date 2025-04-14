package components

import (
	"encoding/json"
	"fmt"
	"rainmcp/pkg/api"
	"rainmcp/pkg/models"
	"rainmcp/pkg/utils"

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
		struct {
			TeamName   string `json:"team_name" description:"团队名称"`
			RegionName string `json:"region_name" description:"集群名称"`
			AppID      string `json:"app_id" description:"应用ID"`
			ServiceID  string `json:"service_id" description:"组件ID"`
		}{},
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
		struct {
			TeamName         string `json:"team_name" description:"团队名称"`
			RegionName       string `json:"region_name" description:"集群名称"`
			GroupID          int    `json:"group_id" description:"应用ID"`
			ServiceCName     string `json:"service_cname" description:"组件名称"`
			K8sComponentName string `json:"k8s_component_name" description:"Kubernetes中的组件名称"`
			Image            string `json:"image" description:"镜像地址"`
			DockerCmd        string `json:"docker_cmd,omitempty" description:"启动命令"`
			UserName         string `json:"user_name,omitempty" description:"镜像仓库用户名"`
			Password         string `json:"password,omitempty" description:"镜像仓库密码"`
			IsDeploy         bool   `json:"is_deploy" description:"是否立即部署"`
		}{},
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
	type ComponentDetailRequest struct {
		TeamName   string `json:"team_name"`
		RegionName string `json:"region_name"`
		AppID      string `json:"app_id"`
		ServiceID  string `json:"service_id"`
	}

	req := new(ComponentDetailRequest)
	if err := protocol.VerifyAndUnmarshal(request.RawArguments, req); err != nil {
		utils.Error("解析获取组件详情请求失败: %v", err)
		return nil, fmt.Errorf("无效的组件详情请求: %v", err)
	}

	// 构建API路径 - 根据Rainbond OpenAPI文档
	path := fmt.Sprintf("/openapi/v1/teams/%s/regions/%s/apps/%s/services/%s",
		req.TeamName, req.RegionName, req.AppID, req.ServiceID)

	utils.Info("获取组件详情: %s", path)

	// 调用Rainbond API获取组件详情
	resp, err := s.client.Get(path)
	if err != nil {
		utils.Error("获取组件详情失败: %v", err)
		return nil, fmt.Errorf("获取组件详情失败: %v", err)
	}

	// 解析响应
	var detailResp models.ComponentDetailResponse
	if err := json.Unmarshal(resp, &detailResp); err != nil {
		utils.Error("解析组件详情响应失败: %v", err)
		return nil, fmt.Errorf("解析组件详情失败: %v", err)
	}

	// 将结果转换为JSON字符串
	resultJSON, err := json.MarshalIndent(detailResp, "", "  ")
	if err != nil {
		utils.Error("序列化组件详情结果失败: %v", err)
		return nil, fmt.Errorf("序列化组件详情失败: %v", err)
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
	req := struct {
		TeamName         string `json:"team_name"`
		RegionName       string `json:"region_name"`
		GroupID          int    `json:"group_id"`
		ServiceCName     string `json:"service_cname"`
		K8sComponentName string `json:"k8s_component_name"`
		Image            string `json:"image"`
		DockerCmd        string `json:"docker_cmd,omitempty"`
		UserName         string `json:"user_name,omitempty"`
		Password         string `json:"password,omitempty"`
		IsDeploy         bool   `json:"is_deploy"`
	}{}
	if err := protocol.VerifyAndUnmarshal(request.RawArguments, &req); err != nil {
		utils.Error("解析创建组件请求失败: %v", err)
		return nil, fmt.Errorf("无效的创建组件请求: %v", err)
	}

	// 构建API路径 - 根据Rainbond OpenAPI文档
	path := fmt.Sprintf("/openapi/v1/teams/%s/regions/%s/apps/%d/services",
		req.TeamName, req.RegionName, req.GroupID)

	utils.Info("创建组件: %s", path)

	// 调用Rainbond API创建组件
	resp, err := s.client.Post(path, req)
	if err != nil {
		utils.Error("创建组件失败: %v", err)
		return nil, fmt.Errorf("创建组件失败: %v", err)
	}

	// 解析响应
	var result map[string]interface{}
	if err := json.Unmarshal(resp, &result); err != nil {
		utils.Error("解析创建组件响应失败: %v", err)
		return nil, fmt.Errorf("解析创建组件响应失败: %v", err)
	}

	// 将结果转换为JSON字符串
	resultJSON, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		utils.Error("序列化创建组件结果失败: %v", err)
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

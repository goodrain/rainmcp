package services

import (
	"rainmcp/pkg/api"
	"rainmcp/pkg/logger"
	"rainmcp/pkg/services/apps"
	"rainmcp/pkg/services/components"
	"rainmcp/pkg/services/regions"
	"rainmcp/pkg/services/teams"

	"github.com/ThinkInAIXYZ/go-mcp/server"
)

// Manager 管理所有Rainbond服务
type Manager struct {
	APIClient        *api.Client
	TeamService      *teams.Service
	RegionService    *regions.Service
	AppService       *apps.Service
	ComponentService *components.Service
}

// NewManager 创建一个新的服务管理器
func NewManager(apiURL string) *Manager {
	logger.Info("[Manager] 创建新的服务管理器: apiURL=%s", apiURL)

	// 验证参数
	if apiURL == "" {
		logger.Warn("[Manager] 警告: Rainbond API URL为空")
		apiURL = "https://rainbond-api.example.com" // 设置一个默认值以避免错误
	}

	logger.Info("[Manager] 创建 API 客户端...")
	client := api.NewClient(apiURL)

	logger.Info("[Manager] 初始化各个服务...")
	manager := &Manager{
		APIClient:        client,
		TeamService:      teams.NewService(client),
		RegionService:    regions.NewService(client),
		AppService:       apps.NewService(client),
		ComponentService: components.NewService(client),
	}

	logger.Info("[Manager] 服务管理器初始化完成")
	return manager
}

// RegisterTeamTools 注册团队相关工具
func RegisterTeamTools(mcpServer *server.Server, manager *Manager) {
	logger.Info("[Manager] 注册团队相关工具...")

	// 验证参数
	if manager == nil {
		logger.Error("[Manager] 错误: 服务管理器为空")
		return
	}

	if manager.TeamService == nil {
		logger.Error("[Manager] 错误: 团队服务为空")
		return
	}

	teams.RegisterTools(mcpServer, manager.TeamService)
	logger.Info("[Manager] 团队相关工具注册完成")
}

// RegisterRegionTools 注册集群相关工具
func RegisterRegionTools(mcpServer *server.Server, manager *Manager) {
	logger.Info("[Manager] 注册集群相关工具...")

	// 验证参数
	if manager == nil {
		logger.Error("[Manager] 错误: 服务管理器为空")
		return
	}

	if manager.RegionService == nil {
		logger.Error("[Manager] 错误: 集群服务为空")
		return
	}

	logger.Info("[Manager] RegionService.client.BaseURL = %s", manager.RegionService.GetBaseURL())
	regions.RegisterTools(mcpServer, manager.RegionService)
	logger.Info("[Manager] 集群相关工具注册完成")
}

// RegisterAppTools 注册应用相关工具
func RegisterAppTools(mcpServer *server.Server, manager *Manager) {
	logger.Info("[Manager] 注册应用相关工具...")

	// 验证参数
	if manager == nil {
		logger.Error("[Manager] 错误: 服务管理器为空")
		return
	}

	if manager.AppService == nil {
		logger.Error("[Manager] 错误: 应用服务为空")
		return
	}

	apps.RegisterTools(mcpServer, manager.AppService)
	logger.Info("[Manager] 应用相关工具注册完成")
}

// RegisterComponentTools 注册组件相关工具
func RegisterComponentTools(mcpServer *server.Server, manager *Manager) {
	logger.Info("[Manager] 注册组件相关工具...")

	// 验证参数
	if manager == nil {
		logger.Error("[Manager] 错误: 服务管理器为空")
		return
	}

	if manager.ComponentService == nil {
		logger.Error("[Manager] 错误: 组件服务为空")
		return
	}

	components.RegisterTools(mcpServer, manager.ComponentService)
	logger.Info("[Manager] 组件相关工具注册完成")
}

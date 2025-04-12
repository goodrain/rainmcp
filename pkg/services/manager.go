package services

import (
	"rainmcp/pkg/api"
	"rainmcp/pkg/services/apps"
	"rainmcp/pkg/services/components"
	"rainmcp/pkg/services/regions"
	"rainmcp/pkg/services/teams"

	"github.com/ThinkInAIXYZ/go-mcp/server"
)

// Manager 管理所有Rainbond服务
type Manager struct {
	APIClient      *api.Client
	TeamService    *teams.Service
	RegionService  *regions.Service
	AppService     *apps.Service
	ComponentService *components.Service
}

// NewManager 创建一个新的服务管理器
func NewManager(apiURL, token string) *Manager {
	client := api.NewClient(apiURL, token)
	
	return &Manager{
		APIClient:      client,
		TeamService:    teams.NewService(client),
		RegionService:  regions.NewService(client),
		AppService:     apps.NewService(client),
		ComponentService: components.NewService(client),
	}
}

// RegisterTeamTools 注册团队相关工具
func RegisterTeamTools(mcpServer *server.Server, manager *Manager) {
	teams.RegisterTools(mcpServer, manager.TeamService)
}

// RegisterRegionTools 注册集群相关工具
func RegisterRegionTools(mcpServer *server.Server, manager *Manager) {
	regions.RegisterTools(mcpServer, manager.RegionService)
}

// RegisterAppTools 注册应用相关工具
func RegisterAppTools(mcpServer *server.Server, manager *Manager) {
	apps.RegisterTools(mcpServer, manager.AppService)
}

// RegisterComponentTools 注册组件相关工具
func RegisterComponentTools(mcpServer *server.Server, manager *Manager) {
	components.RegisterTools(mcpServer, manager.ComponentService)
}

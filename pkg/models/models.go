package models

import "time"

// 通用响应结构
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

// 分页信息
type PageInfo struct {
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
	Total    int `json:"total"`
}

// 团队相关模型
// ===============

// Team 表示Rainbond平台中的团队
type Team struct {
	TeamName   string `json:"team_name"`
	TeamAlias  string `json:"team_alias"`
	TeamID     string `json:"team_id"`
	CreateTime string `json:"create_time"`
	Region     string `json:"region"`
	Role       string `json:"role"`
	Enterprise string `json:"enterprise"`
	Useable    int    `json:"useable"`
}

// TeamsResponse 获取团队列表的响应
type TeamsResponse struct {
	Response
	Data []Team `json:"data"`
}

// 集群相关模型
// ===============

// Region 表示Rainbond平台中的集群
type Region struct {
	ID          string `json:"ID"`
	RegionID    string `json:"region_id"`
	RegionName  string `json:"region_name"`
	RegionAlias string `json:"region_alias"`
	Status      string `json:"status"`
	Desc        string `json:"desc"`
	URL         string `json:"url"`
	WSURL       string `json:"wsurl"`
	HTTPURL     string `json:"httpurl"`
	TCPDomain   string `json:"tcpdomain"`
	Scope       string `json:"scope"`
	SSL         bool   `json:"ssl"`
}

// RegionsResponse 获取集群列表的响应
type RegionsResponse struct {
	Response
	Data []Region `json:"data"`
}

// 应用相关模型
// ===============

// App 表示Rainbond平台中的应用
type App struct {
	ID          string    `json:"ID"`
	GroupName   string    `json:"group_name"`
	UpdateTime  time.Time `json:"update_time"`
	CreateTime  time.Time `json:"create_time"`
	Region      string    `json:"region"`
	RegionAlias string    `json:"region_alias"`
	Status      string    `json:"status"`
	TeamName    string    `json:"tenant_name"`
}

// AppsRequest 表示获取应用列表的请求参数
type AppsRequest struct {
	TeamName   string `json:"team_name" description:"团队名称"`
	RegionName string `json:"region_name" description:"集群名称"`
}

// AppsResponse 获取应用列表的响应
type AppsResponse struct {
	Response
	Data []App `json:"data"`
}

// CreateAppRequest 创建应用的请求参数
type CreateAppRequest struct {
	TeamName   string `json:"team_name" description:"团队名称"`
	RegionName string `json:"region_name" description:"集群名称"`
	AppName    string `json:"app_name" description:"应用名称"` // 修改为app_name字段
	Note       string `json:"note,omitempty" description:"应用描述"`
}

// CreateAppResponse 创建应用的响应
type CreateAppResponse struct {
	Response
	Data App `json:"data"`
}

// 组件相关模型
// ===============

// Component 表示Rainbond平台中的组件
type Component struct {
	ID             string    `json:"ID"`
	ServiceAlias   string    `json:"service_alias"`
	ServiceID      string    `json:"service_id"`
	ServiceCName   string    `json:"service_cname"`
	ServiceType    string    `json:"service_type"`
	ServiceRegion  string    `json:"service_region"`
	DeployVersion  string    `json:"deploy_version"`
	Version        string    `json:"version"`
	CreateTime     time.Time `json:"create_time"`
	UpdateTime     time.Time `json:"update_time"`
	CurStatus      string    `json:"cur_status"`
	Status         string    `json:"status"`
	ContainerMem   int       `json:"container_memory"`
	ContainerCPU   int       `json:"container_cpu"`
	Replicas       int       `json:"replicas"`
	TeamID         string    `json:"tenant_id"`
	TeamName       string    `json:"tenant_name"`
	AppID          string    `json:"group_id"`
	AppName        string    `json:"group_name"`
	ServiceOrigin  string    `json:"service_origin"`
	MemoryWarn     string    `json:"memory_warn"`
	Image          string    `json:"image"`
	K8sComponentID string    `json:"k8s_component_id"`
}

// ComponentDetailResponse 获取组件详情的响应
type ComponentDetailResponse struct {
	Response
	Data Component `json:"data"`
}

// CreateImageComponentRequest 创建组件的请求参数
type CreateImageComponentRequest struct {
	// 路径参数（不包含在请求体中）
	TeamName         string `json:"team_name" description:"团队名称"`
	RegionName       string `json:"region_name" description:"集群名称"`
	GroupID          int    `json:"group_id" description:"应用ID"`
	ServiceCName     string `json:"service_cname" description:"组件名称"`
	K8sComponentName string `json:"k8s_component_name" description:"组件英文名称"`
	Image            string `json:"image" description:"镜像地址"`
	DockerCmd        string `json:"docker_cmd,omitempty" description:"启动命令"`
	UserName         string `json:"user_name,omitempty" description:"镜像仓库用户名"`
	Password         string `json:"password,omitempty" description:"镜像仓库密码"`
	IsDeploy         bool   `json:"is_deploy" description:"是否立即部署"`
}

// BuildComponentRequest 表示构建组件的请求参数
type BuildComponentRequest struct {
	TeamName     string `json:"team_name" description:"团队名称"`
	RegionName   string `json:"region_name" description:"集群名称"`
	AppID        string `json:"app_id" description:"应用ID"`
	ServiceID    string `json:"service_id" description:"组件ID"`
	IsDeploy     bool   `json:"is_deploy" description:"是否部署"`
	BuildVersion string `json:"build_version,omitempty" description:"构建版本"`
}

// ComponentDetailRequest 获取组件详情的请求参数
type ComponentDetailRequest struct {
	TeamName   string `json:"team_name" description:"团队名称"`
	RegionName string `json:"region_name" description:"集群名称"`
	AppID      string `json:"app_id" description:"应用ID"`
	ServiceID  string `json:"service_id" description:"组件ID"`
}

// ListComponentsRequest 获取应用下组件列表的请求参数
type ListComponentsRequest struct {
	TeamName   string `json:"team_name" description:"团队名称"`
	RegionName string `json:"region_name" description:"集群名称"`
	AppID      string `json:"app_id" description:"应用ID"`
}

// ComponentBaseInfo 组件基本信息
type ComponentBaseInfo struct {
	Status               string   `json:"status"`
	AccessInfos          []string `json:"access_infos"`
	ServiceID            string   `json:"service_id"`
	TenantID             string   `json:"tenant_id"`
	ServiceKey           string   `json:"service_key"`
	ServiceAlias         string   `json:"service_alias"`
	ServiceCName         string   `json:"service_cname"`
	ServiceRegion        string   `json:"service_region"`
	Desc                 string   `json:"desc"`
	Category             string   `json:"category"`
	Version              string   `json:"version"`
	UpdateVersion        int      `json:"update_version"`
	Image                string   `json:"image"`
	Cmd                  string   `json:"cmd"`
	MinNode              int      `json:"min_node"`
	MinCPU               int      `json:"min_cpu"`
	ContainerGPU         int      `json:"container_gpu"`
	MinMemory            int      `json:"min_memory"`
	ExtendMethod         string   `json:"extend_method"`
	CodeFrom             string   `json:"code_from"`
	GitURL               string   `json:"git_url"`
	GitProjectID         int      `json:"git_project_id"`
	CodeVersion          string   `json:"code_version"`
	ServiceType          string   `json:"service_type"`
	Creater              int      `json:"creater"`
	Language             string   `json:"language"`
	TotalMemory          int      `json:"total_memory"`
	IsService            bool     `json:"is_service"`
	ServiceOrigin        string   `json:"service_origin"`
	TenantServiceGroupID int      `json:"tenant_service_group_id"`
	OpenWebhooks         bool     `json:"open_webhooks"`
	ServiceSource        string   `json:"service_source"`
	CreateStatus         string   `json:"create_status"`
	CheckUUID            string   `json:"check_uuid"`
	CheckEventID         string   `json:"check_event_id"`
	DockerCmd            string   `json:"docker_cmd"`
	ServerType           string   `json:"server_type"`
	IsUpgrate            bool     `json:"is_upgrate"`
	BuildUpgrade         bool     `json:"build_upgrade"`
	OauthServiceID       int      `json:"oauth_service_id"`
	K8sComponentName     string   `json:"k8s_component_name"`
}

package models

import "time"

// 通用响应结构
type Response struct {
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

// TeamRegion 表示团队关联的集群信息
type TeamRegion struct {
	RegionID    string `json:"region_id"`
	RegionName  string `json:"region_name"`
	RegionAlias string `json:"region_alias"`
	Status      string `json:"status"`
}

// Team 表示Rainbond平台中的团队
type Team struct {
	ID          int          `json:"ID"`
	TenantID    string       `json:"tenant_id"`
	TenantName  string       `json:"tenant_name" description:"团队英文名称"`
	TenantAlias string       `json:"tenant_alias" description:"团队中文名称"`
	CreateTime  string       `json:"create_time"`
	Creater     string       `json:"creater"`
	Regions     []TeamRegion `json:"regions"`
}

// TeamsResponse 表示获取团队列表的响应
type TeamsResponse struct {
	Tenants  []Team `json:"tenants"`
	Total    int    `json:"total"`
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
	Response
}

// 集群相关模型
// ===============

// Region 表示Rainbond平台中的集群
type Region struct {
	ID          string `json:"ID,omitempty"`
	RegionID    string `json:"region_id"`
	RegionName  string `json:"region_name"`
	RegionAlias string `json:"region_alias"`
	Status      string `json:"status"`
	Desc        string `json:"desc"`
	URL         string `json:"url,omitempty"`
	WSURL       string `json:"wsurl,omitempty"`
	HTTPURL     string `json:"httpurl,omitempty"`
	TCPDomain   string `json:"tcpdomain"`
	HTTPDomain  string `json:"httpdomain,omitempty"`
	Scope       string `json:"scope,omitempty"`
	SSL         bool   `json:"ssl,omitempty"`
	WSL         bool   `json:"wsL,omitempty"`
}

// RegionsResponse 表示获取集群列表的响应
type RegionsResponse struct {
	Regions []Region `json:"regions,omitempty"`
	Data    []Region `json:"data,omitempty"`
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
	TeamName    string `json:"tenant_name" description:"团队名称"`
	TenantAlias string `json:"tenant_alias" description:"团队别名"`
	RegionName  string `json:"region_name" description:"集群名称"`
}

// AppItem 表示应用列表中的单个应用项
type AppItem struct {
	AppID      int    `json:"app_id"`
	TenantID   string `json:"tenant_id"`
	GroupName  string `json:"group_name" description:"应用中文名称"`
	RegionName string `json:"region_name"`
	CreateTime string `json:"create_time"`
	K8sApp     string `json:"k8s_app" description:"应用英文名称"`
}

// AppListData 应用列表响应中的数据部分
type AppListData struct {
	List []AppItem `json:"list"`
}

// AppsResponse 获取应用列表的响应
type AppsResponse struct {
	Msg     string      `json:"msg"`
	MsgShow string      `json:"msg_show"`
	Data    AppListData `json:"data"`
}

// CreateAppRequest 创建应用的请求参数
type CreateAppRequest struct {
	TeamName   string `json:"tenant_name" description:"团队名称"`
	RegionName string `json:"region_name" description:"集群名称"`
	AppName    string `json:"app_name" description:"应用名称"` // 修改为app_name字段
	Note       string `json:"note,omitempty" description:"应用描述"`
}

// CreateAppResponse 创建应用的响应
type CreateAppResponse struct {
	Response
	Data App `json:"data"`
}

// ComponentDetailResponse 获取组件详情的响应
type ComponentDetailResponse struct {
	// 基本信息
	ServiceID     string `json:"service_id" description:"组件ID"`
	TenantID      string `json:"tenant_id" description:"团队ID"`
	ServiceAlias  string `json:"service_alias" description:"组件别名"`
	ServiceCName  string `json:"service_cname" description:"组件中文名"`
	ServiceRegion string `json:"service_region" description:"所属集群"`
	Image         string `json:"image" description:"镜像"`
	Cmd           string `json:"cmd" description:"启动命令"`
	MinNode       int    `json:"min_node" description:"最小节点数"`
	MinCPU        int    `json:"min_cpu" description:"CPU配额"`
	ContainerGPU  int    `json:"container_gpu" description:"GPU配额"`
	MinMemory     int    `json:"min_memory" description:"内存配额"`
	ExtendMethod  string `json:"extend_method" description:"伸缩方式"`

	// 源码相关字段
	CodeFrom     string `json:"code_from" description:"源码来源"`
	GitURL       string `json:"git_url" description:"Git仓库地址"`
	GitProjectID int    `json:"git_project_id" description:"Git项目ID"`
	CodeVersion  string `json:"code_version" description:"代码版本"`

	// 状态相关字段
	ServiceType          string        `json:"service_type" description:"服务类型"`
	Creater              int           `json:"creater" description:"创建者ID"`
	Language             interface{}   `json:"language" description:"编程语言"`
	TotalMemory          int           `json:"total_memory" description:"总内存"`
	IsService            bool          `json:"is_service" description:"是否为服务"`
	ServiceOrigin        string        `json:"service_origin" description:"服务来源"`
	TenantServiceGroupID int           `json:"tenant_service_group_id" description:"服务组ID"`
	OpenWebhooks         bool          `json:"open_webhooks" description:"是否开启Webhooks"`
	ServiceSource        string        `json:"service_source" description:"服务源类型，区分源码或镜像"`
	CreateStatus         string        `json:"create_status" description:"创建状态"`
	Status               string        `json:"status" description:"运行状态"`
	AccessInfos          []interface{} `json:"access_infos" description:"访问信息"`

	// Docker相关字段
	DockerCmd  interface{} `json:"docker_cmd" description:"Docker命令"`
	ServerType string      `json:"server_type" description:"服务器类型"`

	// Kubernetes相关字段
	K8sComponentName string `json:"k8s_component_name" description:"Kubernetes组件名称"`
	Arch             string `json:"arch" description:"架构"`
}

// CreateImageComponentRequest 创建组件的请求参数
type CreateImageComponentRequest struct {
	// 路径参数（不包含在请求体中）
	TeamName         string `json:"tenant_name" description:"团队名称"`
	RegionName       string `json:"region_name" description:"集群名称"`
	GroupID          int    `json:"group_id" description:"应用ID"`
	ServiceCName     string `json:"service_cname" description:"组件名称"`
	K8sComponentName string `json:"k8s_component_name" description:"组件英文名称"`
	Image            string `json:"image" description:"镜像地址"`
	DockerCmd        string `json:"docker_cmd,omitempty" description:"启动命令"`
	UserName         string `json:"user_name,omitempty" description:"镜像仓库用户名"`
	Password         string `json:"password,omitempty" description:"镜像仓库密码"`
	// is_deploy 参数在服务器端始终设置为 true
}

// CreateCodeComponentRequest 基于源码创建组件的请求参数
type CreateCodeComponentRequest struct {
	// 路径参数（不包含在请求体中）
	TeamName         string `json:"tenant_name" description:"团队名称"`
	RegionName       string `json:"region_name" description:"集群名称"`
	AppID            string `json:"app_id" description:"应用ID"`
	ServiceCName     string `json:"service_cname" description:"组件名称"`
	K8sComponentName string `json:"k8s_component_name" description:"组件英文名称"`
	RepoURL          string `json:"repo_url" description:"代码仓库地址"`
	Branch           string `json:"branch" description:"分支名称"`
	Username         string `json:"username,omitempty" description:"仓库用户名"`
	Password         string `json:"password,omitempty" description:"仓库密码"`
	// is_deploy 参数在服务器端始终设置为 true
}

// ComponentPort 表示组件端口信息
type ComponentPort struct {
	ID             int    `json:"ID"`
	TenantID       string `json:"tenant_id"`
	ServiceID      string `json:"service_id"`
	ContainerPort  int    `json:"container_port" description:"容器端口"`
	Protocol       string `json:"protocol" description:"协议类型，可选值：tcp/udp/http"`
	IsInnerService bool   `json:"is_inner_service" description:"是否开启对内服务"`
	IsOuterService bool   `json:"is_outer_service" description:"是否开启对外服务"`
	K8sServiceName string `json:"k8s_service_name" description:"Kubernetes服务名称"`
}

// PortResponseData 端口操作响应中的数据部分
type PortResponseData struct {
	Bean ComponentPort `json:"bean"`
}

// PortResponse 端口操作的响应
type PortResponse struct {
	Msg     string           `json:"msg"`
	MsgShow string           `json:"msg_show"`
	Data    PortResponseData `json:"data"`
}

// PortListData 端口列表响应中的数据部分
type PortListData struct {
	List []ComponentPort `json:"list"`
}

// PortListResponse 获取端口列表的响应
type PortListResponse struct {
	Data PortListData `json:"data"`
}

// AddPortRequest 表示添加组件端口的请求参数
type AddPortRequest struct {
	TenantID       string `json:"tenant_id" description:"团队ID"`
	RegionName     string `json:"region_name" description:"集群名称"`
	AppID          string `json:"app_id" description:"应用ID"`
	ServiceID      string `json:"service_id" description:"组件ID"`
	Port           int    `json:"port" description:"端口号"`
	Protocol       string `json:"protocol" description:"协议类型，可选值：tcp/udp/http"`
	IsOuterService bool   `json:"is_outer_service" description:"是否开启对外服务"`
}

// UpdatePortRequest 表示更新组件端口的请求参数
// Action字段可选值及含义：
// - open_outer: 打开对外服务
// - close_outer: 关闭对外服务
// - open_inner: 打开对内服务
// - close_inner: 关闭对内服务
// - change_protocol: 更改端口协议，需要提供protocol参数
// - change_port_alias: 更改端口别名，需要提供port_alias和k8s_service_name参数
type UpdatePortRequest struct {
	TenantID   string `json:"tenant_id" description:"团队ID"`
	RegionName string `json:"region_name" description:"集群名称"`
	AppID      string `json:"app_id" description:"应用ID"`
	ServiceID  string `json:"service_id" description:"组件ID"`
	Port       int    `json:"port" description:"端口号"`
	Action     string `json:"action" description:"操作类型，可选值：open_outer/close_outer/open_inner/close_inner/change_protocol"`
	Protocol   string `json:"protocol,omitempty" description:"协议类型，当action为change_protocol时使用，可选值：tcp/udp/http"`
}

// ListPortsRequest 表示获取组件端口列表的请求参数
type ListPortsRequest struct {
	TenantID   string `json:"tenant_id" description:"团队ID"`
	RegionName string `json:"region_name" description:"集群名称"`
	AppID      string `json:"app_id" description:"应用ID"`
	ServiceID  string `json:"service_id" description:"组件ID"`
}

// DeletePortRequest 表示删除组件端口的请求参数
type DeletePortRequest struct {
	TenantID   string `json:"tenant_id" description:"团队ID"`
	RegionName string `json:"region_name" description:"集群名称"`
	AppID      string `json:"app_id" description:"应用ID"`
	ServiceID  string `json:"service_id" description:"组件ID"`
	Port       int    `json:"port" description:"端口号"`
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
	TenantName string `json:"tenant_name" description:"团队英文名称"`
	RegionName string `json:"region_name" description:"集群名称"`
	AppID      string `json:"app_id" description:"应用ID"`
}

// ComponentBaseInfo 组件基本信息
type ComponentBaseInfo struct {
	Status           string `json:"status"`
	ServiceID        string `json:"service_id"`
	ServiceCName     string `json:"service_cname" description:"组件中文名称"`
	K8sComponentName string `json:"k8s_component_name" description:"组件英文名称"`
}

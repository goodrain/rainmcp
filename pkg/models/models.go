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

// TeamRegionInfo 表示团队关联的区域信息
type TeamRegionInfo struct {
	RegionName  string `json:"region_name"`
	RegionAlias string `json:"region_alias"`
}

// Team 表示Rainbond平台中的团队
type Team struct {
	TeamAlias  string           `json:"team_alias" description:"团队中文名称"`
	CreateTime string           `json:"create_time"`
	OwnerName  string           `json:"owner_name"`
	RegionList []TeamRegionInfo `json:"region_list"`
}

// TeamsData 表示团队数据结构
type TeamsData struct {
	Bean interface{} `json:"bean"`
	List []Team      `json:"list"`
}

// TeamsResponse 表示获取团队列表的响应
type TeamsResponse struct {
	Code    int       `json:"code"`
	Msg     string    `json:"msg"`
	MsgShow string    `json:"msg_show"`
	Data    TeamsData `json:"data"`
}

// TeamRegion 表示团队关联的集群信息(保留旧结构以兼容)
type TeamRegion struct {
	RegionID    string `json:"region_id"`
	RegionName  string `json:"region_name"`
	RegionAlias string `json:"region_alias"`
	Status      string `json:"status"`
}

// LegacyTeam 表示旧版本的团队结构(保留以兼容)
type LegacyTeam struct {
	ID          int          `json:"ID"`
	TenantID    string       `json:"tenant_id"`
	TenantName  string       `json:"tenant_name" description:"团队英文名称"`
	TenantAlias string       `json:"tenant_alias" description:"团队中文名称"`
	CreateTime  string       `json:"create_time"`
	Creater     string       `json:"creater"`
	Regions     []TeamRegion `json:"regions"`
}

// 集群相关模型
// ===============

// RegionInfo 表示集群基本信息（新版本API响应）
type RegionInfo struct {
	RegionName  string `json:"region_name"`
	RegionAlias string `json:"region_alias"`
	Status      string `json:"status"`
	Desc        string `json:"desc"`
}

// RegionsData 表示集群数据结构
type RegionsData struct {
	Bean interface{}  `json:"bean"`
	List []RegionInfo `json:"list"`
}

// RegionsResponse 表示获取集群列表的响应（新版本）
type RegionsResponse struct {
	Code    int         `json:"code"`
	Msg     string      `json:"msg"`
	MsgShow string      `json:"msg_show"`
	Data    RegionsData `json:"data"`
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
	TeamAlias   string    `json:"team_alias"`
}

// AppsRequest 表示获取应用列表的请求参数
type AppsRequest struct {
	TeamAlias  string `json:"team_alias" description:"团队别名"`
	RegionName string `json:"region_name" description:"集群名称"`
}

// AppItem 表示应用列表中的单个应用项（新版本）
type AppItem struct {
	GroupID     int     `json:"group_id" description:"应用ID"`
	GroupName   string  `json:"group_name" description:"应用名称"`
	Description *string `json:"description" description:"应用描述"`
	UpdateTime  string  `json:"update_time" description:"更新时间"`
	CreateTime  string  `json:"create_time" description:"创建时间"`
}

// AppListData 应用列表响应中的数据部分
type AppListData struct {
	Bean interface{} `json:"bean"`
	List []AppItem   `json:"list"`
}

// AppsResponse 获取应用列表的响应（新版本）
type AppsResponse struct {
	Code    int         `json:"code"`
	Msg     string      `json:"msg"`
	MsgShow string      `json:"msg_show"`
	Data    AppListData `json:"data"`
}

// LegacyAppItem 表示应用列表中的单个应用项（旧版本，保留以兼容）
type LegacyAppItem struct {
	AppID      int    `json:"app_id"`
	TenantID   string `json:"tenant_id"`
	GroupName  string `json:"group_name" description:"应用中文名称"`
	RegionName string `json:"region_name"`
	CreateTime string `json:"create_time"`
	K8sApp     string `json:"k8s_app" description:"应用英文名称"`
}

// LegacyAppListData 应用列表响应中的数据部分（旧版本）
type LegacyAppListData struct {
	List []LegacyAppItem `json:"list"`
}

// LegacyAppsResponse 获取应用列表的响应（旧版本，保留以兼容）
type LegacyAppsResponse struct {
	Msg     string            `json:"msg"`
	MsgShow string            `json:"msg_show"`
	Data    LegacyAppListData `json:"data"`
}

// CreateAppRequest 创建应用的请求参数
type CreateAppRequest struct {
	TeamAlias  string `json:"team_alias" description:"团队别名"`
	RegionName string `json:"region_name" description:"集群名称"`
	AppName    string `json:"app_name" description:"应用名称"` // 修改为app_name字段
}

// CreateAppResponse 创建应用的响应
type CreateAppResponse struct {
	Response
	Data App `json:"data"`
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

// PortInfo 端口信息（新版本API响应）
type PortInfo struct {
	Port           int    `json:"port" description:"端口号"`
	Protocol       string `json:"protocol" description:"协议类型"`
	IsOuterService bool   `json:"is_outer_service" description:"是否开启对外服务"`
	IsInnerService bool   `json:"is_inner_service" description:"是否开启对内服务"`
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

// PortListData 端口列表响应中的数据部分（新版本）
type PortListData struct {
	Bean interface{} `json:"bean"`
	List []PortInfo  `json:"list"`
}

// PortListResponse 获取端口列表的响应（新版本）
type PortListResponse struct {
	Code    int          `json:"code"`
	Msg     string       `json:"msg"`
	MsgShow string       `json:"msg_show"`
	Data    PortListData `json:"data"`
}

// AddPortRequest 表示添加组件端口的请求参数
type AddPortRequest struct {
	TeamAlias      string `json:"team_alias" description:"团队别名"`
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
	TeamAlias string `json:"team_alias" description:"团队别名"`
	AppID     string `json:"app_id" description:"应用ID"`
	ServiceID string `json:"service_id" description:"组件ID"`
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
	TeamAlias    string `json:"team_alias" description:"团队名称"`
	RegionName   string `json:"region_name" description:"集群名称"`
	AppID        string `json:"app_id" description:"应用ID"`
	ServiceID    string `json:"service_id" description:"组件ID"`
	IsDeploy     bool   `json:"is_deploy" description:"是否部署"`
	BuildVersion string `json:"build_version,omitempty" description:"构建版本"`
}

// ComponentDetailRequest 获取组件详情的请求参数
type ComponentDetailRequest struct {
	TeamAlias string `json:"team_alias" description:"团队别名"`
	AppID     string `json:"app_id" description:"应用ID"`
	ServiceID string `json:"service_id" description:"组件ID"`
}

// ListComponentsRequest 获取应用下组件列表的请求参数
type ListComponentsRequest struct {
	TeamAlias string `json:"team_alias" description:"团队别名"`
	AppID     string `json:"app_id" description:"应用ID"`
}

// ComponentBaseInfo 组件基本信息
type ComponentBaseInfo struct {
	Status           string `json:"status"`
	ServiceID        string `json:"service_id"`
	ServiceCName     string `json:"service_cname" description:"组件中文名称"`
	K8sComponentName string `json:"k8s_component_name" description:"组件英文名称"`
}

// ComponentInfo 组件信息（新版本API响应）
type ComponentInfo struct {
	ServiceID    string `json:"service_id" description:"组件ID"`
	ServiceCName string `json:"service_cname" description:"组件中文名称"`
	UpdateTime   string `json:"update_time" description:"更新时间"`
	Status       string `json:"status" description:"组件状态"`
}

// ComponentListData 组件列表响应中的数据部分
type ComponentListData struct {
	Bean interface{}     `json:"bean"`
	List []ComponentInfo `json:"list"`
}

// ComponentListResponse 获取组件列表的响应（新版本）
type ComponentListResponse struct {
	Code    int               `json:"code"`
	Msg     string            `json:"msg"`
	MsgShow string            `json:"msg_show"`
	Data    ComponentListData `json:"data"`
}

// ComponentPortInfo 组件端口信息（用于组件详情）
type ComponentPortInfo struct {
	ContainerPort  int      `json:"container_port" description:"容器端口"`
	Protocol       string   `json:"protocol" description:"协议类型"`
	IsOuterService bool     `json:"is_outer_service" description:"是否开启对外服务"`
	IsInnerService bool     `json:"is_inner_service" description:"是否开启对内服务"`
	AccessUrls     []string `json:"access_urls" description:"访问地址列表"`
}

// ComponentEnv 组件环境变量信息
type ComponentEnv struct {
	AttrName  string `json:"attr_name" description:"环境变量名"`
	AttrValue string `json:"attr_value" description:"环境变量值"`
	Name      string `json:"name" description:"显示名称"`
	Scope     string `json:"scope" description:"作用域"`
	IsChange  bool   `json:"is_change" description:"是否可更改"`
}

// ComponentVolume 组件存储卷信息
type ComponentVolume struct {
	VolumeName     string `json:"volume_name" description:"存储卷名称"`
	VolumePath     string `json:"volume_path" description:"挂载路径"`
	VolumeCapacity int    `json:"volume_capacity" description:"存储容量(GB)"`
}

// ComponentDetailInfo 组件详情信息（新版本API响应）
type ComponentDetailInfo struct {
	ServiceID    string              `json:"service_id" description:"组件ID"`
	ServiceCName string              `json:"service_cname" description:"组件中文名"`
	ServiceAlias string              `json:"service_alias" description:"组件别名"`
	UpdateTime   string              `json:"update_time" description:"更新时间"`
	MinMemory    int                 `json:"min_memory" description:"内存配额(MB)"`
	MinCPU       int                 `json:"min_cpu" description:"CPU配额(毫核)"`
	StatusCN     string              `json:"status_cn" description:"状态中文"`
	Ports        []ComponentPortInfo `json:"ports" description:"端口列表"`
	Envs         []ComponentEnv      `json:"envs" description:"环境变量列表"`
	Volumes      []ComponentVolume   `json:"volumes" description:"存储卷列表"`
}

// NewComponentDetailResponse 获取组件详情的响应（新版本）
type NewComponentDetailResponse struct {
	Code    int    `json:"code"`
	Msg     string `json:"msg"`
	MsgShow string `json:"msg_show"`
	Data    struct {
		Bean ComponentDetailInfo `json:"bean"`
		List []interface{}       `json:"list"`
	} `json:"data"`
}

// CreateCodeComponentRequest 基于源码创建组件的请求参数
type CreateCodeComponentRequest struct {
	TeamAlias    string `json:"team_alias" description:"团队名称"`
	AppID        string `json:"app_id" description:"应用ID"`
	ServiceCName string `json:"service_cname" description:"组件名称"`
	RepoURL      string `json:"repo_url" description:"代码仓库地址"`
	Branch       string `json:"branch" description:"分支名称"`
	Username     string `json:"username,omitempty" description:"仓库用户名"`
	Password     string `json:"password,omitempty" description:"仓库密码"`
}

type RainTokenKey struct{}

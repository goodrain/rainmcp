package models

// Team 表示Rainbond平台中的团队
type Team struct {
	TeamName  string `json:"team_name"`
	TeamAlias string `json:"team_alias"`
	TeamID    string `json:"team_id"`
	CreateTime string `json:"create_time"`
	Region    string `json:"region"`
	Role      string `json:"role"`
}

// Region 表示Rainbond平台中的集群
type Region struct {
	RegionID   string `json:"region_id"`
	RegionName string `json:"region_name"`
	RegionAlias string `json:"region_alias"`
	Status     string `json:"status"`
	Desc       string `json:"desc"`
	Url        string `json:"url"`
	WSUrl      string `json:"wsurl"`
	HTTPUrl    string `json:"httpurl"`
	TCPDomain  string `json:"tcpdomain"`
}

// App 表示Rainbond平台中的应用
type App struct {
	ID          string `json:"ID"`
	GroupName   string `json:"group_name"`
	UpdateTime  string `json:"update_time"`
	CreateTime  string `json:"create_time"`
	Region      string `json:"region"`
	RegionAlias string `json:"region_alias"`
}

// AppsRequest 表示获取应用列表的请求参数
type AppsRequest struct {
	TeamName   string `json:"team_name" description:"团队名称"`
	RegionName string `json:"region_name" description:"集群名称"`
}

// ComponentRequest 表示构建组件的请求参数
type ComponentRequest struct {
	TeamName     string `json:"team_name" description:"团队名称"`
	RegionName   string `json:"region_name" description:"集群名称"`
	AppID        string `json:"app_id" description:"应用ID"`
	ServiceID    string `json:"service_id" description:"组件ID"`
	IsDeploy     bool   `json:"is_deploy" description:"是否部署"`
	BuildVersion string `json:"build_version" description:"构建版本"`
}

# Rainbond MCP Server

这是一个基于MCP协议的Rainbond OpenAPI服务器，提供了与Rainbond平台交互的能力，支持通过自然语言对话方式管理云原生应用。

## 功能特点

- 通过SSE方式暴露MCP服务，支持实时交互
- 支持环境变量配置Rainbond API地址和访问令牌
- 提供丰富的Rainbond平台管理功能：
  - **团队管理**：获取团队列表 (rainbond_teams)
  - **集群管理**：获取集群列表 (rainbond_regions)
  - **应用管理**：获取应用列表 (rainbond_apps)
  - **组件管理**：
    - 获取组件列表 (rainbond_list_components)
    - 获取组件详情 (rainbond_get_component_detail)
    - 创建镜像组件 (rainbond_create_image_component)
    - 创建源码组件 (rainbond_create_code_component)
    - 构建组件 (rainbond_build_service)
  - **端口管理**：
    - 获取组件端口列表 (rainbond_list_component_ports)
    - 添加组件端口 (rainbond_add_component_port)
    - 更新组件端口 (rainbond_update_component_port)
    - 删除组件端口 (rainbond_delete_component_port)
- 实现了完整的错误处理和优雅关闭机制
- 支持Docker容器化部署

## 项目结构

```
rainmcp/
├── cmd/
│   └── rainmcp/
│       └── main.go           # 程序入口点
├── pkg/
│   ├── api/
│   │   └── client.go         # Rainbond API客户端
│   ├── models/
│   │   └── models.go         # 数据模型定义
│   ├── services/
│   │   ├── manager.go        # 服务管理器
│   │   ├── teams/            # 团队相关服务
│   │   ├── regions/          # 集群相关服务
│   │   ├── apps/             # 应用相关服务
│   │   └── components/       # 组件相关服务
│   ├── transport/
│   │   └── sse.go            # SSE传输层
│   └── utils/                # 工具函数
└── go.mod                    # Go模块定义
```

## 使用方法

### 环境变量配置

- `RAINBOND_HOST`: MCP服务器监听地址，默认为 "localhost:8080"
- `RAINBOND_API`: Rainbond API地址，例如 "https://api.rainbond.com"
- `RAINBOND_TOKEN`: Rainbond API访问令牌

### 构建和运行

#### 本地运行

```bash
# 构建
go build -o rainmcp ./cmd/rainmcp

# 运行
export RAINBOND_API="https://your-rainbond-api.com"
export RAINBOND_TOKEN="your-token"
./rainmcp
```

#### Docker 部署

```bash
# 构建 Docker 镜像
docker build -t rainmcp:latest .

# 运行容器
docker run -d --name rainmcp \
  -p 8080:8080 \
  -e RAINBOND_API="https://your-rainbond-api.com" \
  -e RAINBOND_TOKEN="your-token" \
  rainmcp:latest
```

## API说明

### 团队管理

#### 获取团队列表

工具名称: `rainbond_teams`  
描述: 获取Rainbond平台中的团队列表  
参数: 无

### 集群管理

#### 获取集群列表

工具名称: `rainbond_regions`  
描述: 获取Rainbond平台中的集群列表  
参数: 无

### 应用管理

#### 获取应用列表

工具名称: `rainbond_apps`  
描述: 获取Rainbond平台中的应用列表  
参数:
- `team_name`: 团队名称
- `region_name`: 集群名称

### 组件管理

#### 获取组件列表

工具名称: `rainbond_list_components`  
描述: 获取应用下的组件列表  
参数:
- `team_name`: 团队名称
- `region_name`: 集群名称
- `group_id`: 应用ID

#### 获取组件详情

工具名称: `rainbond_get_component_detail`  
描述: 获取组件的详细信息  
参数:
- `team_name`: 团队名称
- `region_name`: 集群名称
- `service_id`: 组件ID

#### 创建镜像组件

工具名称: `rainbond_create_image_component`  
描述: 基于镜像创建新组件  
参数:
- `team_name`: 团队名称
- `region_name`: 集群名称
- `group_id`: 应用ID
- `service_cname`: 组件中文名
- `k8s_component_name`: 组件英文名
- `image`: 镜像地址
- `is_deploy`: 是否立即部署

#### 创建源码组件

工具名称: `rainbond_create_code_component`  
描述: 基于源码创建新组件  
参数:
- `team_name`: 团队名称
- `region_name`: 集群名称
- `group_id`: 应用ID
- `service_cname`: 组件中文名
- `k8s_component_name`: 组件英文名
- `code_type`: 代码类型
- `git_url`: Git仓库地址
- `code_version`: 代码版本
- `is_deploy`: 是否立即部署

#### 构建组件

工具名称: `rainbond_build_service`  
描述: 在Rainbond平台中构建组件  
参数:
- `team_name`: 团队名称
- `region_name`: 集群名称
- `app_id`: 应用ID
- `service_id`: 组件ID
- `is_deploy`: 是否部署
- `build_version`: 构建版本

### 端口管理

#### 获取组件端口列表

工具名称: `rainbond_list_component_ports`  
描述: 获取组件的端口列表  
参数:
- `team_name`: 团队名称
- `region_name`: 集群名称
- `service_id`: 组件ID

#### 添加组件端口

工具名称: `rainbond_add_component_port`  
描述: 为组件添加新的端口  
参数:
- `team_name`: 团队名称
- `region_name`: 集群名称
- `service_id`: 组件ID
- `port`: 端口号
- `protocol`: 协议类型（tcp/udp/http）
- `is_outer_service`: 是否开启对外服务

#### 更新组件端口

工具名称: `rainbond_update_component_port`  
描述: 更新组件的端口配置  
参数:
- `team_name`: 团队名称
- `region_name`: 集群名称
- `service_id`: 组件ID
- `port`: 端口号
- `protocol`: 协议类型（tcp/udp/http）
- `is_outer_service`: 是否开启对外服务
- `action`: 操作类型（如开启对外服务、关闭对外服务等）

#### 删除组件端口

工具名称: `rainbond_delete_component_port`  
描述: 删除组件的端口  
参数:
- `team_name`: 团队名称
- `region_name`: 集群名称
- `service_id`: 组件ID
- `port`: 端口号

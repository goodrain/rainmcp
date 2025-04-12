# Rainbond MCP Server

这是一个基于MCP协议的Rainbond OpenAPI服务器，提供了与Rainbond平台交互的能力。

## 功能特点

- 通过SSE方式暴露MCP服务
- 支持环境变量配置Rainbond API地址和访问令牌
- 提供以下Rainbond相关功能：
  - 获取团队列表 (rainbond_teams)
  - 获取集群列表 (rainbond_regions)
  - 获取应用列表 (rainbond_apps)
  - 构建组件 (rainbond_build_service)
- 实现了完整的错误处理和优雅关闭机制

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

```bash
# 构建
go build -o rainmcp ./cmd/rainmcp

# 运行
export RAINBOND_API="https://your-rainbond-api.com"
export RAINBOND_TOKEN="your-token"
./rainmcp
```

## API说明

### 获取团队列表

工具名称: `rainbond_teams`
描述: 获取Rainbond平台中的团队列表
参数: 无

### 获取集群列表

工具名称: `rainbond_regions`
描述: 获取Rainbond平台中的集群列表
参数: 无

### 获取应用列表

工具名称: `rainbond_apps`
描述: 获取Rainbond平台中的应用列表
参数:
- `team_name`: 团队名称
- `region_name`: 集群名称

### 构建组件

工具名称: `rainbond_build_service`
描述: 在Rainbond平台中构建组件
参数:
- `team_name`: 团队名称
- `region_name`: 集群名称
- `app_id`: 应用ID
- `service_id`: 组件ID
- `is_deploy`: 是否部署
- `build_version`: 构建版本

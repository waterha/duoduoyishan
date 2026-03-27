# 多多益善 (DuoDuoYiShan)

一个基于 Go 语言开发的现代化社交平台，支持实时聊天、好友管理、社区互动等功能。

## 项目简介

多多益善是一个功能完整的社交平台，采用前后端分离架构，提供实时通信、好友管理、社区互动等核心功能。

## 主要功能

- **用户系统**：注册登录、JWT认证、信息管理
- **即时通讯**：私聊、群聊、实时消息推送(WebSocket)
- **好友管理**：添加好友、处理请求、好友列表
- **社区管理**：创建社区、加入退出、成员管理
- **其他功能**：文件上传、在线状态、请求限流

## 技术栈

- **后端**：Go 1.25 + Gin + GORM + MySQL + Redis + WebSocket
- **前端**：HTML5 + CSS3 + JavaScript
- **部署**：Docker + Docker Compose

## 快速开始

### 使用 Docker Compose（推荐）

1. 克隆项目
```bash
git clone https://gitee.com/firefishvick/duoduoyishan.git
cd duoduoyishan
```

2. 启动服务
```bash
docker-compose up -d
```

3. 访问应用
```
http://localhost:8080
```

### 本地运行

1. 安装依赖
```bash
go mod download
```

2. 配置数据库和Redis
编辑 `config/config.yaml` 文件修改连接信息

3. 运行项目
```bash
go run main.go
```

4. 访问应用
```
http://localhost:8080
```

## 项目结构

```
duoduoyishan/
├── cache/              # Redis 缓存
├── config/             # 配置文件
├── controller/         # 控制器
├── database/           # 数据库连接
├── middleware/         # 中间件
├── models/             # 数据模型
├── router/             # 路由配置
├── service/            # 业务逻辑
├── static/             # 静态文件
├── utils/              # 工具函数
├── websocket_own/      # WebSocket 处理
├── Dockerfile          # Docker 构建文件
├── docker-compose.yml  # Docker Compose 配置
└── main.go             # 程序入口
```

## 部署

### 生产环境配置

1. 修改 `config/config.yaml` 中的配置：
   - 将 `server.mode` 改为 `release`
   - 修改 `jwt.secret` 为强密码
   - 配置正确的数据库和Redis连接信息

2. 使用 Docker Compose 部署
```bash
docker-compose up -d
```

## 常见问题

- **Redis连接失败**：确保Redis服务已启动，检查配置是否正确
- **数据库连接失败**：确保MySQL服务已启动，检查连接信息
- **WebSocket连接失败**：检查JWT token是否有效
- **文件上传失败**：确保uploads目录存在且有写入权限

## 许可证

本项目采用 MIT 许可证。

## 联系方式

如有问题或建议，请提交 Issue 或联系项目维护者。

---

**多多益善** - 简约社交，连接你我
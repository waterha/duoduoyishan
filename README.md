# 多多益善 (DuoDuoYiShan)

一个基于 Go 语言开发的现代化社交平台，支持实时聊天、好友管理、社区互动等功能。

## 项目简介

多多益善是一个功能完整的社交平台，采用前后端分离架构，提供实时通信、好友管理、社区互动等核心功能。项目使用 Gin 框架构建 RESTful API，结合 WebSocket 实现实时消息推送，使用 Redis 缓存提升性能，MySQL 作为数据存储。

## 主要功能

### 用户系统
- 用户注册与登录
- JWT 身份认证
- 用户信息管理（昵称、头像、个性签名等）
- 密码修改
- 用户搜索

### 即时通讯
- 私聊功能
- 群聊功能（社区）
- 实时消息推送（WebSocket）
- 消息历史记录
- 消息撤回
- 已读/未读状态
- 支持多种消息类型（文本、图片、文件、语音、视频）

### 好友管理
- 发送好友请求
- 处理好友请求（同意/拒绝）
- 好友列表管理
- 删除好友
- 好友备注

### 社区管理
- 创建社区
- 加入/退出社区
- 社区列表浏览
- 我的社区管理
- 社区成员管理

### 其他功能
- 文件上传
- 用户在线状态
- 请求限流
- 日志记录
- CORS 跨域支持

## 技术栈

### 后端
- **语言**: Go 1.25
- **框架**: Gin v1.12.0
- **ORM**: GORM v1.31.1
- **数据库**: MySQL 8.0
- **缓存**: Redis 7.0
- **WebSocket**: Gorilla WebSocket v1.5.0
- **认证**: JWT (golang-jwt/jwt v5.0.0)
- **配置管理**: Viper v1.17.0
- **日志**: Logrus v1.9.3
- **密码加密**: bcrypt (golang.org/x/crypto)

### 前端
- **HTML5/CSS3/JavaScript**
- **Font Awesome** 图标库
- **Google Fonts** 字体

### 部署
- **Docker** 容器化部署
- **Docker Compose** 编排管理

## 项目结构

```
duoduoyishan/
├── cache/              # Redis 缓存层
│   └── redis.go
├── config/             # 配置文件
│   ├── config.go
│   └── config.yaml
├── controller/         # 控制器层
│   ├── community_controller.go
│   ├── friend_controller.go
│   ├── message_controller.go
│   ├── user_controller.go
│   └── ws_controller.go
├── database/           # 数据库连接
│   └── mysql.go
├── logs/               # 日志文件
│   └── app.log
├── middleware/         # 中间件
│   ├── auth.go         # JWT 认证
│   ├── cors.go         # 跨域处理
│   ├── logger.go       # 日志中间件
│   └── ratelimit.go    # 限流中间件
├── models/             # 数据模型
│   ├── chat_room.go
│   ├── community.go
│   ├── friend.go
│   ├── message.go
│   └── user.go
├── router/             # 路由配置
│   └── router.go
├── service/            # 业务逻辑层
│   ├── community_service.go
│   ├── friend_service.go
│   ├── message_service.go
│   └── user_service.go
├── static/             # 静态文件
│   ├── css/
│   │   └── style.css
│   ├── js/
│   │   └── app.js
│   └── index.html
├── utils/              # 工具函数
│   ├── encrypt.go      # 加密工具
│   ├── jwt.go          # JWT 工具
│   ├── logger.go       # 日志工具
│   ├── response.go     # 响应封装
│   └── upload.go       # 文件上传
├── websocket_own/      # WebSocket 处理
│   ├── client.go
│   └── hub.go
├── .env                # 环境变量
├── Dockerfile          # Docker 构建文件
├── docker-compose.yml  # Docker Compose 配置
├── go.mod              # Go 模块依赖
├── go.sum              # 依赖版本锁定
└── main.go             # 程序入口
```

## 快速开始

### 环境要求
- Go 1.25 或更高版本
- MySQL 8.0
- Redis 7.0
- Docker 和 Docker Compose（可选）

### 方式一：使用 Docker Compose（推荐）

1. 克隆项目
```bash
git clone <repository-url>
cd duoduoyishan
```

2. 使用 Docker Compose 启动所有服务
```bash
docker-compose up -d
```

3. 访问应用
```
http://localhost:8080
```

### 方式二：本地运行

1. 安装依赖
```bash
go mod download
```

2. 配置数据库和 Redis

编辑 `config/config.yaml` 文件，修改数据库和 Redis 连接信息：
```yaml
database:
  host: "localhost"
  port: "3306"
  username: "root"
  password: "your_password"
  database: "duoduoyishan"

redis:
  host: "localhost"
  port: "6379"
  password: ""
  db: 0
```

3. 启动 MySQL 和 Redis

确保 MySQL 和 Redis 服务已启动。

4. 运行项目
```bash
go run main.go
```

5. 访问应用
```
http://localhost:8080
```

## 配置说明

### 配置文件 (config/config.yaml)

```yaml
server:
  port: "8080"          # 服务端口
  mode: "debug"         # 运行模式: debug/release

database:
  host: "localhost"     # 数据库主机
  port: "3306"          # 数据库端口
  username: "root"      # 数据库用户名
  password: "123456"    # 数据库密码
  database: "duoduoyishan"  # 数据库名
  charset: "utf8mb4"    # 字符集
  max_idle_conns: 10    # 最大空闲连接数
  max_open_conns: 100   # 最大打开连接数

redis:
  host: "localhost"     # Redis 主机
  port: "6379"          # Redis 端口
  password: ""          # Redis 密码
  db: 0                 # Redis 数据库
  pool_size: 10         # 连接池大小

jwt:
  secret: "your-256-bit-secret-key-change-in-production"  # JWT 密钥（生产环境请修改）
  expire_time: 86400    # Token 过期时间（秒），默认 24 小时

upload:
  max_size: 10485760    # 最大上传文件大小（10MB）
  save_path: "./uploads"  # 文件保存路径
  allow_exts:           # 允许的文件扩展名
    - ".jpg"
    - ".jpeg"
    - ".png"
    - ".gif"
    - ".mp4"
    - ".mp3"
    - ".pdf"
    - ".doc"
    - ".docx"

log:
  level: "info"         # 日志级别
  filename: "./logs/app.log"  # 日志文件路径
  max_size: 100         # 日志文件最大大小（MB）
  max_backups: 30       # 保留的旧日志文件数量
  max_age: 30           # 日志文件保留天数
```

## API 接口文档

### 认证接口

#### 用户注册
```
POST /api/auth/register
Content-Type: application/json

{
  "username": "testuser",
  "password": "password123",
  "email": "test@example.com"
}
```

#### 用户登录
```
POST /api/auth/login
Content-Type: application/json

{
  "username": "testuser",
  "password": "password123"
}

Response:
{
  "code": 200,
  "message": "success",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "username": "testuser",
      "nickname": "测试用户",
      "avatar": "",
      "email": "test@example.com",
      "status": 1
    }
  }
}
```

#### 退出登录
```
POST /api/auth/logout
Authorization: Bearer <token>
```

### 用户接口（需要认证）

#### 获取用户信息
```
GET /api/user/info
Authorization: Bearer <token>
```

#### 更新用户信息
```
PUT /api/user/info
Authorization: Bearer <token>
Content-Type: application/json

{
  "nickname": "新昵称",
  "signature": "这是我的个性签名",
  "gender": 1
}
```

#### 修改密码
```
PUT /api/user/password
Authorization: Bearer <token>
Content-Type: application/json

{
  "old_password": "oldpassword",
  "new_password": "newpassword"
}
```

#### 搜索用户
```
GET /api/user/search?keyword=用户名&page=1&page_size=20
Authorization: Bearer <token>
```

### 好友接口（需要认证）

#### 发送好友请求
```
POST /api/friend/request
Authorization: Bearer <token>
Content-Type: application/json

{
  "to_user_id": 2,
  "message": "你好，我想添加你为好友"
}
```

#### 处理好友请求
```
PUT /api/friend/request/:id
Authorization: Bearer <token>
Content-Type: application/json

{
  "status": 1  // 1:同意 2:拒绝 3:忽略
}
```

#### 获取好友列表
```
GET /api/friend/list
Authorization: Bearer <token>
```

#### 删除好友
```
DELETE /api/friend/:id
Authorization: Bearer <token>
```

#### 获取好友请求
```
GET /api/friend/requests
Authorization: Bearer <token>
```

### 社区接口（需要认证）

#### 创建社区
```
POST /api/community/create
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "技术交流群",
  "description": "讨论各种技术话题",
  "category": "技术",
  "max_members": 200
}
```

#### 加入社区
```
POST /api/community/join/:id
Authorization: Bearer <token>
```

#### 退出社区
```
POST /api/community/quit/:id
Authorization: Bearer <token>
```

#### 获取社区列表
```
GET /api/community/list?page=1&page_size=20
Authorization: Bearer <token>
```

#### 获取社区详情
```
GET /api/community/detail/:id
Authorization: Bearer <token>
```

#### 获取我的社区
```
GET /api/community/my
Authorization: Bearer <token>
```

### 消息接口（需要认证）

#### 获取聊天历史
```
GET /api/message/history?to_type=1&to_id=2&page=1&page_size=50
Authorization: Bearer <token>
```

#### 获取未读消息数
```
GET /api/message/unread
Authorization: Bearer <token>
```

#### 标记消息已读
```
PUT /api/message/read/:id
Authorization: Bearer <token>
```

#### 撤回消息
```
PUT /api/message/recall/:id
Authorization: Bearer <token>
```

### WebSocket 接口

#### 连接 WebSocket
```
GET /api/ws?token=<token>&room_id=<room_id>
```

#### 获取房间在线人数
```
GET /api/ws/online?room_id=<room_id>
Authorization: Bearer <token>
```

## WebSocket 消息格式

### 发送消息
```json
{
  "type": "message",
  "to_type": 1,
  "to_id": 2,
  "msg_type": 1,
  "content": "你好"
}
```

### 接收消息
```json
{
  "type": "message",
  "msg_id": "msg_123456",
  "from_user_id": 1,
  "from_username": "testuser",
  "from_nickname": "测试用户",
  "from_avatar": "",
  "to_type": 1,
  "to_id": 2,
  "msg_type": 1,
  "content": "你好",
  "created_at": "2024-01-01T12:00:00Z"
}
```

### 用户状态通知
```json
{
  "type": "user_status",
  "user_id": 1,
  "online": true,
  "time": "2024-01-01T12:00:00Z"
}
```

## 数据库表结构

### users（用户表）
- id: 用户ID
- username: 用户名
- password: 密码（加密存储）
- nickname: 昵称
- avatar: 头像
- email: 邮箱
- phone: 手机号
- gender: 性别（0:未知 1:男 2:女）
- birthday: 生日
- signature: 个性签名
- status: 状态（1:在线 2:离线 3:隐身）
- last_login_at: 最后登录时间
- last_login_ip: 最后登录IP
- created_at: 创建时间
- updated_at: 更新时间

### messages（消息表）
- id: 消息ID
- msg_id: 消息唯一标识
- from_user_id: 发送者ID
- to_type: 接收类型（1:私聊 2:群聊）
- to_id: 接收者ID（用户ID或社区ID）
- msg_type: 消息类型（1:文本 2:图片 3:文件 4:语音 5:视频 6:系统消息）
- content: 消息内容
- media_url: 媒体文件URL
- media_size: 媒体文件大小
- duration: 语音/视频时长（秒）
- status: 状态（1:未读 2:已读 3:撤回）
- created_at: 创建时间

### communities（社区表）
- id: 社区ID
- name: 社区名称
- description: 社区描述
- avatar: 社区头像
- creator_id: 创建者ID
- category: 社区分类
- member_count: 成员数量
- max_members: 最大成员数
- status: 状态（1:正常 2:封禁）
- created_at: 创建时间
- updated_at: 更新时间

### community_members（社区成员表）
- id: 成员ID
- community_id: 社区ID
- user_id: 用户ID
- role: 角色（1:成员 2:管理员 3:群主）
- nickname: 群昵称
- join_time: 加入时间
- last_read_at: 最后阅读时间
- created_at: 创建时间
- updated_at: 更新时间

### friends（好友表）
- id: 好友关系ID
- user_id: 用户ID
- friend_id: 好友ID
- remark: 备注
- status: 状态（1:正常 2:拉黑）
- created_at: 创建时间
- updated_at: 更新时间

### friend_requests（好友请求表）
- id: 请求ID
- from_user_id: 发送者ID
- to_user_id: 接收者ID
- message: 请求消息
- status: 状态（0:待处理 1:同意 2:拒绝 3:忽略）
- created_at: 创建时间
- updated_at: 更新时间

## 开发说明

### 添加新功能

1. 在 `models/` 中定义数据模型
2. 在 `service/` 中实现业务逻辑
3. 在 `controller/` 中创建控制器
4. 在 `router/` 中注册路由
5. 如需认证，在路由组中使用 `middleware.JWTAuth()`

### 运行测试

```bash
go test ./...
```

### 代码格式化

```bash
go fmt ./...
```

### 依赖管理

```bash
# 添加新依赖
go get <package-name>

# 更新依赖
go get -u <package-name>

# 整理依赖
go mod tidy
```

## 部署

### Docker 部署

1. 构建镜像
```bash
docker build -t duoduoyishan:latest .
```

2. 运行容器
```bash
docker run -d -p 8080:8080 duoduoyishan:latest
```

### Docker Compose 部署

```bash
docker-compose up -d
```

### 生产环境配置

1. 修改 `config/config.yaml` 中的配置：
   - 将 `server.mode` 改为 `release`
   - 修改 `jwt.secret` 为强密码
   - 配置正确的数据库和 Redis 连接信息
   - 调整日志级别和文件路径

2. 使用环境变量覆盖配置（可选）

## 常见问题

### Redis 连接失败
确保 Redis 服务已启动，并检查 `config/config.yaml` 中的 Redis 配置是否正确。

### 数据库连接失败
确保 MySQL 服务已启动，并检查数据库用户名、密码和数据库名是否正确。

### WebSocket 连接失败
检查 JWT token 是否有效，确保在连接时传递正确的 token 参数。

### 文件上传失败
检查 `uploads/` 目录是否存在，并确保有写入权限。同时检查 `config.yaml` 中的上传配置。

## 贡献指南

欢迎提交 Issue 和 Pull Request！

1. Fork 本项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

## 许可证

本项目采用 MIT 许可证。

## 联系方式

如有问题或建议，请提交 Issue 或联系项目维护者。

## 更新日志

### v1.0.0 (2024-01-01)
- 初始版本发布
- 实现用户注册、登录功能
- 实现私聊和群聊功能
- 实现好友管理功能
- 实现社区管理功能
- 实现 WebSocket 实时通信
- 添加文件上传功能
- 添加日志记录和限流功能

## 致谢

感谢所有为本项目做出贡献的开发者！

---

**多多益善** - 简约社交，连接你我

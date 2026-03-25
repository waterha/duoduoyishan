package controller

import (
	"duoduoyishan/cache"
	"duoduoyishan/service"
	"duoduoyishan/utils"
	"duoduoyishan/websocket_own"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // 生产环境需要限制
	},
}

type WSController struct {
	hub *websocket_own.Hub
}

func NewWSController(hub *websocket_own.Hub) *WSController {
	return &WSController{hub: hub}
}

// WebSocket连接处理
func (ctrl *WSController) Connect(c *gin.Context) {
	// 从URL参数获取token
	token := c.Query("token")
	utils.Logger.Infof("WebSocket连接请求，token: %s", token)
	if token == "" {
		// 尝试从Header获取
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) == 2 && parts[0] == "Bearer" {
				token = parts[1]
				utils.Logger.Infof("从Header获取token: %s", token)
			}
		}
	}

	if token == "" {
		utils.Logger.Errorf("WebSocket连接失败：未提供认证令牌")
		c.JSON(401, gin.H{"code": 401, "message": "未提供认证令牌"})
		return
	}

	// 解析token
	claims, err := utils.ParseToken(token)
	if err != nil {
		utils.Logger.Errorf("WebSocket连接失败：无效的认证令牌，错误: %v", err)
		c.JSON(401, gin.H{"code": 401, "message": "无效的认证令牌"})
		return
	}
	utils.Logger.Infof("解析token成功，用户ID: %d", claims.UserID)

	// 检查session是否有效
	sessionUserID, err := cache.GetUserSession(token)
	if err != nil {
		utils.Logger.Errorf("WebSocket连接失败：会话已过期，错误: %v", err)
		c.JSON(401, gin.H{"code": 401, "message": "会话已过期"})
		return
	}
	if sessionUserID != claims.UserID {
		utils.Logger.Errorf("WebSocket连接失败：会话用户ID不匹配，sessionUserID: %d, claims.UserID: %d", sessionUserID, claims.UserID)
		c.JSON(401, gin.H{"code": 401, "message": "会话已过期"})
		return
	}
	utils.Logger.Infof("会话验证成功，用户ID: %d", sessionUserID)

	userID := claims.UserID

	// 获取房间ID（可选）
	roomID := c.Query("room_id")
	if roomID == "" {
		roomID = "global" // 默认全局房间
	}

	// 升级为WebSocket连接
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		utils.Logger.Errorf("WebSocket升级失败: %v", err)
		return
	}

	// 创建客户端
	client := &websocket_own.Client{
		Hub:            ctrl.hub,
		Conn:           conn,
		Send:           make(chan []byte, 256),
		UserID:         userID,
		RoomID:         roomID,
		MessageService: service.NewMessageService(),
	}

	// 注册客户端
	client.Hub.Register <- client

	// 更新用户在线状态
	cache.SetUserOnline(userID, true)

	// 启动读写协程
	go client.Write()
	go client.Read()

	utils.Logger.Infof("用户%d建立WebSocket连接，房间:%s", userID, roomID)
}

// 获取房间在线人数
func (ctrl *WSController) GetRoomOnlineCount(c *gin.Context) {
	roomID := c.Query("room_id")
	if roomID == "" {
		utils.ValidationError(c, "房间ID不能为空")
		return
	}

	count := ctrl.hub.GetRoomOnlineCount(roomID)

	utils.Success(c, gin.H{
		"room_id":      roomID,
		"online_count": count,
	})
}

package websocket_own

import (
	"duoduoyishan/cache"
	"duoduoyishan/service"
	"encoding/json"
	"log"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

// UintToString 辅助函数：uint转string
func UintToString(u uint) string {
	return strconv.FormatUint(uint64(u), 10)
}

type Client struct {
	Hub            *Hub
	Conn           *websocket.Conn
	Send           chan []byte
	UserID         uint
	RoomID         string
	MessageService *service.MessageService
}

type Message struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

// 读取消息
func (c *Client) Read() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	// 设置读取超时
	c.Conn.SetReadLimit(512 * 1024) // 512KB
	c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		var message Message
		err := c.Conn.ReadJSON(&message)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket读取错误: %v", err)
			}
			break
		}

		// 处理不同类型的消息
		c.handleMessage(message)
	}
}

// 写入消息
func (c *Client) Write() {
	ticker := time.NewTicker(30 * time.Second)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				// 通道已关闭
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// 将队列中的消息一起发送
			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			// 发送心跳
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// 处理消息
func (c *Client) handleMessage(msg Message) {
	switch msg.Type {
	case "ping":
		// 心跳响应
		c.Send <- []byte(`{"type":"pong"}`)

	case "message":
		// 聊天消息
		c.handleChatMessage(msg.Data)

	case "read":
		// 消息已读
		c.handleReadMessage(msg.Data)

	case "typing":
		// 正在输入
		c.handleTyping(msg.Data)
	}
}

// 处理聊天消息
func (c *Client) handleChatMessage(data []byte) {
	var chatMsg struct {
		ToType  int    `json:"to_type"`
		ToID    uint   `json:"to_id"`
		MsgType int    `json:"msg_type"`
		Content string `json:"content"`
	}

	if err := json.Unmarshal(data, &chatMsg); err != nil {
		return
	}

	// 调用service保存消息
	message, err := c.MessageService.SendMessage(c.UserID, chatMsg.ToType, chatMsg.ToID, chatMsg.MsgType, chatMsg.Content)
	if err != nil {
		log.Printf("保存消息失败: %v", err)
		return
	}

	// 构建消息响应
	response := map[string]interface{}{
		"type": "message",
		"data": map[string]interface{}{
			"id":           message.ID,
			"msg_id":       message.MsgID,
			"from_user_id": message.FromUserID,
			"to_type":      message.ToType,
			"to_id":        message.ToID,
			"msg_type":     message.MsgType,
			"content":      message.Content,
			"created_at":   message.CreatedAt,
			"status":       message.Status,
		},
	}

	msgBytes, _ := json.Marshal(response)

	// 转发消息
	if chatMsg.ToType == 1 { // 私聊
		c.Hub.SendToUser(chatMsg.ToID, msgBytes)
		// 发送给发送者
		c.Send <- msgBytes
	} else { // 群聊
		roomID := "group:" + UintToString(chatMsg.ToID)
		c.Hub.SendToRoom(roomID, msgBytes)
	}

	// 增加未读计数
	if chatMsg.ToType == 1 {
		cache.IncrUnreadCount(chatMsg.ToID, "user:"+UintToString(c.UserID))
	} else {
		// 群聊的未读计数需要单独处理
	}
}

// 处理消息已读
func (c *Client) handleReadMessage(data []byte) {
	var readMsg struct {
		RoomID string `json:"room_id"`
		MsgID  uint   `json:"msg_id"`
	}

	if err := json.Unmarshal(data, &readMsg); err != nil {
		return
	}

	// 清除未读计数
	cache.ClearUnreadCount(c.UserID, readMsg.RoomID)
}

// 处理正在输入
func (c *Client) handleTyping(data []byte) {
	var typingMsg struct {
		ToID   uint `json:"to_id"`
		ToType int  `json:"to_type"`
	}

	if err := json.Unmarshal(data, &typingMsg); err != nil {
		return
	}

	// 广播给对方
	msg := map[string]interface{}{
		"type":    "typing",
		"user_id": c.UserID,
		"to_id":   typingMsg.ToID,
	}
	msgBytes, _ := json.Marshal(msg)

	// 私聊时只发给对方
	if typingMsg.ToType == 1 {
		c.Hub.SendToUser(typingMsg.ToID, msgBytes)
	} else {
		// 群聊时发给房间所有人
		roomID := "group:" + UintToString(typingMsg.ToID)
		c.Hub.SendToRoom(roomID, msgBytes)
	}
}

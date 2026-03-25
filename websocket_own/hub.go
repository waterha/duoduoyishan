package websocket_own

import (
	"encoding/json"
	"sync"
	"time"
)

type Hub struct {
	// 所有连接的客户端
	clients map[*Client]bool

	// 客户端注册通道
	Register chan *Client

	// 客户端注销通道
	Unregister chan *Client

	// 广播消息通道
	Broadcast chan []byte

	// 房间管理
	rooms map[string]map[*Client]bool

	// 读写锁
	mu sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan []byte),
		rooms:      make(map[string]map[*Client]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.registerClient(client)

		case client := <-h.Unregister:
			h.unregisterClient(client)

		case message := <-h.Broadcast:
			h.broadcastMessage(message)
		}
	}
}

// 注册客户端
func (h *Hub) registerClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.clients[client] = true

	// 加入房间
	if _, ok := h.rooms[client.RoomID]; !ok {
		h.rooms[client.RoomID] = make(map[*Client]bool)
	}
	h.rooms[client.RoomID][client] = true

	// 通知用户上线
	h.notifyUserStatus(client.UserID, true)
}

// 注销客户端
func (h *Hub) unregisterClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, ok := h.clients[client]; ok {
		delete(h.clients, client)
		close(client.Send)

		// 从房间移除
		if room, ok := h.rooms[client.RoomID]; ok {
			delete(room, client)
		}

		// 通知用户下线
		h.notifyUserStatus(client.UserID, false)
	}
}

// 广播消息给所有客户端
func (h *Hub) broadcastMessage(message []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for client := range h.clients {
		select {
		case client.Send <- message:
		default:
			close(client.Send)
			delete(h.clients, client)
		}
	}
}

// 发送消息到指定房间
func (h *Hub) SendToRoom(roomID string, message []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if clients, ok := h.rooms[roomID]; ok {
		for client := range clients {
			select {
			case client.Send <- message:
			default:
				close(client.Send)
				delete(clients, client)
			}
		}
	}
}

// 发送消息给指定用户
func (h *Hub) SendToUser(userID uint, message []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for client := range h.clients {
		if client.UserID == userID {
			select {
			case client.Send <- message:
			default:
				close(client.Send)
				delete(h.clients, client)
			}
			break
		}
	}
}

// 通知用户在线状态
func (h *Hub) notifyUserStatus(userID uint, online bool) {
	statusMsg := map[string]interface{}{
		"type":    "user_status",
		"user_id": userID,
		"online":  online,
		"time":    time.Now(),
	}

	msgBytes, _ := json.Marshal(statusMsg)

	// 广播给所有在线用户
	h.mu.RLock()
	defer h.mu.RUnlock()

	for client := range h.clients {
		if client.UserID != userID {
			select {
			case client.Send <- msgBytes:
			default:
			}
		}
	}
}

// 获取房间在线人数
func (h *Hub) GetRoomOnlineCount(roomID string) int {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if room, ok := h.rooms[roomID]; ok {
		return len(room)
	}
	return 0
}

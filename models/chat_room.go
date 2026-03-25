package models

import "time"

type ChatRoom struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	RoomType    int       `gorm:"not null" json:"room_type"` // 1:私聊 2:群聊
	RoomKey     string    `gorm:"type:varchar(100);uniqueIndex" json:"room_key"`
	Name        string    `gorm:"type:varchar(100)" json:"name"`
	Avatar      string    `gorm:"type:varchar(255)" json:"avatar"`
	LastMsgID   uint      `json:"last_msg_id"`
	LastMsg     string    `gorm:"type:varchar(500)" json:"last_msg"`
	LastMsgTime time.Time `json:"last_msg_time"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type RoomUser struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	RoomID      uint      `gorm:"index:idx_room_user,unique" json:"room_id"`
	UserID      uint      `gorm:"index:idx_room_user,unique" json:"user_id"`
	UnreadCount int       `gorm:"default:0" json:"unread_count"`
	LastReadID  uint      `json:"last_read_id"`
	LastReadAt  time.Time `json:"last_read_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

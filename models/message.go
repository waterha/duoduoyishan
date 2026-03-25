package models

import "time"

type Message struct {
	ID         uint      `gorm:"primarykey" json:"id"`
	MsgID      string    `gorm:"type:varchar(50);uniqueIndex" json:"msg_id"`
	FromUserID uint      `gorm:"index;not null" json:"from_user_id"`
	ToType     int       `gorm:"not null" json:"to_type"` // 1:私聊 2:群聊
	ToID       uint      `gorm:"index;not null" json:"to_id"`
	MsgType    int       `gorm:"default:1" json:"msg_type"` // 1:文本 2:图片 3:文件 4:语音 5:视频 6:系统消息
	Content    string    `gorm:"type:text" json:"content"`
	MediaURL   string    `gorm:"type:varchar(500)" json:"media_url"`
	MediaSize  int64     `json:"media_size"`
	Duration   int       `json:"duration"`                // 语音/视频时长(秒)
	Status     int       `gorm:"default:1" json:"status"` // 1:未读 2:已读 3:撤回
	CreatedAt  time.Time `json:"created_at"`
}

type MessageRead struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	MessageID uint      `gorm:"index" json:"message_id"`
	UserID    uint      `gorm:"index" json:"user_id"`
	ReadAt    time.Time `json:"read_at"`
}

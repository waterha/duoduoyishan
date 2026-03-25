package models

import "time"

type Friend struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	UserID    uint      `gorm:"index:idx_user_friend,unique;not null" json:"user_id"`
	FriendID  uint      `gorm:"index:idx_user_friend,unique;not null" json:"friend_id"`
	Remark    string    `gorm:"type:varchar(50)" json:"remark"`
	Status    int       `gorm:"default:1" json:"status"` // 1:正常 2:拉黑
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type FriendRequest struct {
	ID         uint      `gorm:"primarykey" json:"id"`
	FromUserID uint      `gorm:"index;not null" json:"from_user_id"`
	ToUserID   uint      `gorm:"index;not null" json:"to_user_id"`
	Message    string    `gorm:"type:varchar(200)" json:"message"`
	Status     int       `gorm:"default:0" json:"status"` // 0:待处理 1:同意 2:拒绝 3:忽略
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

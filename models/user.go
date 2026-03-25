package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	Username    string         `gorm:"type:varchar(50);uniqueIndex;not null" json:"username"`
	Password    string         `gorm:"type:varchar(255);not null" json:"-"`
	Nickname    string         `gorm:"type:varchar(50)" json:"nickname"`
	Avatar      string         `gorm:"type:varchar(255)" json:"avatar"`
	Email       string         `gorm:"type:varchar(100);uniqueIndex" json:"email"`
	Phone       string         `gorm:"type:varchar(20)" json:"phone"`
	Gender      int            `gorm:"default:0" json:"gender"` // 0:未知 1:男 2:女
	Birthday    *time.Time     `json:"birthday"`
	Signature   string         `gorm:"type:varchar(200)" json:"signature"`
	Status      int            `gorm:"default:2" json:"status"` // 1:在线 2:离线 3:隐身
	LastLoginAt *time.Time     `json:"last_login_at"`
	LastLoginIP string         `gorm:"type:varchar(50)" json:"last_login_ip"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (u *User) TableName() string {
	return "users"
}

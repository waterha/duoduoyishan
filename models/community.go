package models

import "time"

type Community struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	Name        string    `gorm:"type:varchar(100);not null" json:"name"`
	Description string    `gorm:"type:varchar(500)" json:"description"`
	Avatar      string    `gorm:"type:varchar(255)" json:"avatar"`
	CreatorID   uint      `gorm:"not null" json:"creator_id"`
	Category    string    `gorm:"type:varchar(50)" json:"category"`
	MemberCount int       `gorm:"default:1" json:"member_count"`
	MaxMembers  int       `gorm:"default:200" json:"max_members"`
	Status      int       `gorm:"default:1" json:"status"` // 1:正常 2:封禁
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CommunityMember struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	CommunityID uint      `gorm:"index:idx_community_user,unique" json:"community_id"`
	UserID      uint      `gorm:"index:idx_community_user,unique" json:"user_id"`
	Role        int       `gorm:"default:1" json:"role"` // 1:成员 2:管理员 3:群主
	Nickname    string    `gorm:"type:varchar(50)" json:"nickname"`
	JoinTime    time.Time `json:"join_time"`
	LastReadAt  time.Time `json:"last_read_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
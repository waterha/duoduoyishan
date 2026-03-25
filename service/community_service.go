package service

import (
	"duoduoyishan/database"
	"duoduoyishan/models"
	"errors"
	"time"

	"gorm.io/gorm"
)

type CommunityService struct{}

func NewCommunityService() *CommunityService {
	return &CommunityService{}
}

// 创建社区
func (s *CommunityService) CreateCommunity(creatorID uint, name, description, category string) (*models.Community, error) {
	// 检查社区名是否已存在
	var count int64
	database.DB.Model(&models.Community{}).Where("name = ?", name).Count(&count)
	if count > 0 {
		return nil, errors.New("社区名已存在")
	}

	community := &models.Community{
		Name:        name,
		Description: description,
		CreatorID:   creatorID,
		Category:    category,
		MemberCount: 1,
		MaxMembers:  200,
		Status:      1,
	}

	tx := database.DB.Begin()

	if err := tx.Create(community).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// 添加创建者为群主
	member := &models.CommunityMember{
		CommunityID: community.ID,
		UserID:      creatorID,
		Role:        3, // 群主
		Nickname:    "",
		JoinTime:    time.Now(),
		LastReadAt:  time.Now(),
	}

	if err := tx.Create(member).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()
	return community, nil
}

// 加入社区
func (s *CommunityService) JoinCommunity(userID, communityID uint) error {
	var community models.Community
	if err := database.DB.First(&community, communityID).Error; err != nil {
		return errors.New("社区不存在")
	}

	if community.MemberCount >= community.MaxMembers {
		return errors.New("社区人数已满")
	}

	// 检查是否已加入
	var count int64
	database.DB.Model(&models.CommunityMember{}).
		Where("community_id = ? AND user_id = ?", communityID, userID).
		Count(&count)
	if count > 0 {
		return errors.New("已加入该社区")
	}

	tx := database.DB.Begin()

	member := &models.CommunityMember{
		CommunityID: communityID,
		UserID:      userID,
		Role:        1,
		JoinTime:    time.Now(),
		LastReadAt:  time.Now(),
	}

	if err := tx.Create(member).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 更新社区人数
	if err := tx.Model(&community).Update("member_count", community.MemberCount+1).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// 退出社区
func (s *CommunityService) QuitCommunity(userID, communityID uint) error {
	var member models.CommunityMember
	if err := database.DB.Where("community_id = ? AND user_id = ?", communityID, userID).First(&member).Error; err != nil {
		return errors.New("未加入该社区")
	}

	// 群主不能退出
	if member.Role == 3 {
		return errors.New("群主不能退出社区")
	}

	tx := database.DB.Begin()

	if err := tx.Delete(&member).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 更新社区人数
	if err := tx.Model(&models.Community{}).Where("id = ?", communityID).
		Update("member_count", gorm.Expr("member_count - ?", 1)).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// 获取社区列表
func (s *CommunityService) GetCommunityList(category string, page, pageSize int) ([]models.Community, int64, error) {
	var communities []models.Community
	var total int64

	db := database.DB.Model(&models.Community{}).Where("status = 1")

	if category != "" {
		db = db.Where("category = ?", category)
	}

	db.Count(&total)

	offset := (page - 1) * pageSize
	err := db.Offset(offset).Limit(pageSize).Order("created_at desc").Find(&communities).Error

	return communities, total, err
}

// 获取社区详情
func (s *CommunityService) GetCommunityDetail(communityID uint) (*models.Community, []models.CommunityMember, error) {
	var community models.Community
	if err := database.DB.First(&community, communityID).Error; err != nil {
		return nil, nil, err
	}

	var members []models.CommunityMember
	if err := database.DB.Where("community_id = ?", communityID).Order("role desc, created_at").Find(&members).Error; err != nil {
		return nil, nil, err
	}

	return &community, members, nil
}

// 获取用户加入的社区
func (s *CommunityService) GetUserCommunities(userID uint) ([]models.Community, error) {
	var communities []models.Community

	err := database.DB.Table("communities").
		Select("communities.*").
		Joins("join community_members on communities.id = community_members.community_id").
		Where("community_members.user_id = ?", userID).
		Find(&communities).Error

	return communities, err
}

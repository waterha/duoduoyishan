package service

import (
	"duoduoyishan/cache"
	"duoduoyishan/database"
	"duoduoyishan/models"
	"errors"
)

type FriendService struct{}

func NewFriendService() *FriendService {
	return &FriendService{}
}

// 发送好友请求
func (s *FriendService) SendFriendRequest(fromUserID, toUserID uint, message string) error {
	if fromUserID == toUserID {
		return errors.New("不能添加自己为好友")
	}

	// 检查是否已经是好友
	var count int64
	database.DB.Model(&models.Friend{}).Where("user_id = ? AND friend_id = ?", fromUserID, toUserID).Count(&count)
	if count > 0 {
		return errors.New("已经是好友关系")
	}

	// 检查是否已存在待处理的请求
	database.DB.Model(&models.FriendRequest{}).
		Where("from_user_id = ? AND to_user_id = ? AND status = 0", fromUserID, toUserID).
		Count(&count)
	if count > 0 {
		return errors.New("已发送过好友请求")
	}

	request := &models.FriendRequest{
		FromUserID: fromUserID,
		ToUserID:   toUserID,
		Message:    message,
		Status:     0,
	}

	return database.DB.Create(request).Error
}

// 处理好友请求
func (s *FriendService) HandleFriendRequest(requestID uint, status int) error {
	var request models.FriendRequest
	if err := database.DB.First(&request, requestID).Error; err != nil {
		return err
	}

	if request.Status != 0 {
		return errors.New("请求已处理")
	}

	// 更新请求状态
	request.Status = status

	if err := database.DB.Save(&request).Error; err != nil {
		return err
	}

	// 如果同意，添加好友关系
	if status == 1 {
		tx := database.DB.Begin()

		// 双向添加好友
		friend1 := &models.Friend{
			UserID:   request.FromUserID,
			FriendID: request.ToUserID,
		}

		friend2 := &models.Friend{
			UserID:   request.ToUserID,
			FriendID: request.FromUserID,
		}

		if err := tx.Create(friend1).Error; err != nil {
			tx.Rollback()
			return err
		}

		if err := tx.Create(friend2).Error; err != nil {
			tx.Rollback()
			return err
		}

		return tx.Commit().Error
	}

	return nil
}

// 获取好友列表
func (s *FriendService) GetFriends(userID uint) ([]map[string]interface{}, error) {
	var friends []models.Friend
	if err := database.DB.Where("user_id = ? AND status = 1", userID).Find(&friends).Error; err != nil {
		return nil, err
	}

	var result []map[string]interface{}
	for _, friend := range friends {
		var user models.User
		if err := database.DB.First(&user, friend.FriendID).Error; err != nil {
			continue
		}

		// 获取在线状态
		status := 2
		if cache.IsUserOnline(user.ID) {
			status = 1
		}

		result = append(result, map[string]interface{}{
			"id":       user.ID,
			"username": user.Username,
			"nickname": user.Nickname,
			"avatar":   user.Avatar,
			"remark":   friend.Remark,
			"status":   status,
			"email":    user.Email,
		})
	}

	return result, nil
}

// 删除好友
func (s *FriendService) DeleteFriend(userID, friendID uint) error {
	tx := database.DB.Begin()

	// 删除双向好友关系
	if err := tx.Where("user_id = ? AND friend_id = ?", userID, friendID).Delete(&models.Friend{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Where("user_id = ? AND friend_id = ?", friendID, userID).Delete(&models.Friend{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// 获取好友请求列表
func (s *FriendService) GetFriendRequests(userID uint) ([]models.FriendRequest, error) {
	var requests []models.FriendRequest
	err := database.DB.Where("to_user_id = ? AND status = 0", userID).
		Order("created_at desc").
		Find(&requests).Error
	return requests, err
}

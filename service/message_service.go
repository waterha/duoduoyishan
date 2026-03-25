package service

import (
	"duoduoyishan/database"
	"duoduoyishan/models"
	"errors"
	"time"
	"github.com/google/uuid"
)

type MessageService struct{}

func NewMessageService() *MessageService {
	return &MessageService{}
}

// 发送消息
func (s *MessageService) SendMessage(fromUserID uint, toType int, toID uint, msgType int, content string) (*models.Message, error) {
	// 生成唯一消息ID
	msgID := uuid.New().String()

	message := &models.Message{
		MsgID:      msgID,
		FromUserID: fromUserID,
		ToType:     toType,
		ToID:       toID,
		MsgType:    msgType,
		Content:    content,
		Status:     1, // 未读
	}

	if err := database.DB.Create(message).Error; err != nil {
		return nil, err
	}

	return message, nil
}

// 获取聊天记录
func (s *MessageService) GetChatHistory(userID uint, toType int, toID uint, page, pageSize int) ([]models.Message, int64, error) {
	var messages []models.Message
	var total int64

	db := database.DB.Model(&models.Message{})

	// 根据聊天类型构建查询条件
	if toType == 1 { // 私聊
		db = db.Where("(from_user_id = ? AND to_id = ? AND to_type = 1) OR (from_user_id = ? AND to_id = ? AND to_type = 1)",
			userID, toID, toID, userID)
	} else if toType == 2 { // 群聊
		db = db.Where("to_id = ? AND to_type = 2", toID)
	} else {
		return nil, 0, errors.New("无效的聊天类型")
	}

	// 计算总数
	db.Count(&total)

	// 分页查询
	offset := (page - 1) * pageSize
	err := db.Order("created_at desc").Offset(offset).Limit(pageSize).Find(&messages).Error

	return messages, total, err
}

// 标记消息已读
func (s *MessageService) MarkMessageRead(messageID uint, userID uint) error {
	// 检查消息是否存在
	var message models.Message
	if err := database.DB.First(&message, messageID).Error; err != nil {
		return err
	}

	// 检查是否是消息的接收者
	if message.ToType == 1 && message.ToID != userID {
		return errors.New("无权限标记此消息")
	}

	// 标记消息已读
	message.Status = 2
	if err := database.DB.Save(&message).Error; err != nil {
		return err
	}

	// 记录已读状态
	messageRead := &models.MessageRead{
		MessageID: messageID,
		UserID:    userID,
		ReadAt:    time.Now(),
	}

	return database.DB.Create(messageRead).Error
}

// 获取未读消息数
func (s *MessageService) GetUnreadCount(userID uint) (int64, error) {
	var count int64
	err := database.DB.Model(&models.Message{}).Where("to_id = ? AND to_type = 1 AND status = 1", userID).Count(&count).Error
	return count, err
}

// 撤回消息
func (s *MessageService) RecallMessage(messageID uint, userID uint) error {
	// 检查消息是否存在
	var message models.Message
	if err := database.DB.First(&message, messageID).Error; err != nil {
		return err
	}

	// 检查是否是消息的发送者
	if message.FromUserID != userID {
		return errors.New("无权限撤回此消息")
	}

	// 撤回消息
	message.Status = 3
	return database.DB.Save(&message).Error
}

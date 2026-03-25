package controller

import (
	"duoduoyishan/service"
	"duoduoyishan/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MessageController struct {
	messageService *service.MessageService
}

func NewMessageController() *MessageController {
	return &MessageController{
		messageService: service.NewMessageService(),
	}
}

// @Summary 获取聊天记录
// @Tags 消息管理
// @Security BearerAuth
// @Param to_type query int true "聊天类型 1:私聊 2:群聊"
// @Param to_id query int true "目标ID"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(50)
// @Success 200 {object} utils.Response
// @Router /message/history [get]
func (ctrl *MessageController) GetChatHistory(c *gin.Context) {
	userID := c.GetUint("userID")
	
	toType, _ := strconv.Atoi(c.Query("to_type"))
	toID, _ := strconv.ParseUint(c.Query("to_id"), 10, 32)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "50"))
	
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 50
	}
	
	messages, total, err := ctrl.messageService.GetChatHistory(userID, toType, uint(toID), page, pageSize)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	
	utils.Success(c, gin.H{
		"messages": messages,
		"total":    total,
		"page":     page,
		"page_size": pageSize,
	})
}

// @Summary 获取未读消息数
// @Tags 消息管理
// @Security BearerAuth
// @Success 200 {object} utils.Response
// @Router /message/unread [get]
func (ctrl *MessageController) GetUnreadCount(c *gin.Context) {
	userID := c.GetUint("userID")
	
	count, err := ctrl.messageService.GetUnreadCount(userID)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	
	utils.Success(c, gin.H{
		"unread_count": count,
	})
}

// @Summary 标记消息已读
// @Tags 消息管理
// @Security BearerAuth
// @Param id path int true "消息ID"
// @Success 200 {object} utils.Response
// @Router /message/read/{id} [put]
func (ctrl *MessageController) MarkMessageRead(c *gin.Context) {
	userID := c.GetUint("userID")
	
	messageID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ValidationError(c, "无效的消息ID")
		return
	}
	
	if err := ctrl.messageService.MarkMessageRead(uint(messageID), userID); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	
	utils.Success(c, gin.H{"message": "标记成功"})
}

// @Summary 撤回消息
// @Tags 消息管理
// @Security BearerAuth
// @Param id path int true "消息ID"
// @Success 200 {object} utils.Response
// @Router /message/recall/{id} [put]
func (ctrl *MessageController) RecallMessage(c *gin.Context) {
	userID := c.GetUint("userID")
	
	messageID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ValidationError(c, "无效的消息ID")
		return
	}
	
	if err := ctrl.messageService.RecallMessage(uint(messageID), userID); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	
	utils.Success(c, gin.H{"message": "撤回成功"})
}

package controller

import (
	"duoduoyishan/service"
	"duoduoyishan/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type FriendController struct {
	friendService *service.FriendService
}

func NewFriendController() *FriendController {
	return &FriendController{
		friendService: service.NewFriendService(),
	}
}

// @Summary 发送好友请求
// @Tags 好友管理
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body friendRequest true "好友请求"
// @Success 200 {object} utils.Response
// @Router /friend/request [post]
func (ctrl *FriendController) SendFriendRequest(c *gin.Context) {
	userID := c.GetUint("userID")

	var req struct {
		ToUserID uint   `json:"to_user_id" binding:"required"`
		Message  string `json:"message"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, "参数错误："+err.Error())
		return
	}

	if err := ctrl.friendService.SendFriendRequest(userID, req.ToUserID, req.Message); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, gin.H{"message": "请求已发送"})
}

// @Summary 处理好友请求
// @Tags 好友管理
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "请求ID"
// @Param request body handleFriendRequest true "处理信息"
// @Success 200 {object} utils.Response
// @Router /friend/request/{id} [put]
func (ctrl *FriendController) HandleFriendRequest(c *gin.Context) {
	requestID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ValidationError(c, "无效的请求ID")
		return
	}

	var req struct {
		Status int `json:"status" binding:"required,oneof=1 2"` // 1:同意 2:拒绝
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, "参数错误："+err.Error())
		return
	}

	if err := ctrl.friendService.HandleFriendRequest(uint(requestID), req.Status); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, gin.H{"message": "操作成功"})
}

// @Summary 获取好友列表
// @Tags 好友管理
// @Security BearerAuth
// @Success 200 {object} utils.Response
// @Router /friend/list [get]
func (ctrl *FriendController) GetFriends(c *gin.Context) {
	userID := c.GetUint("userID")

	friends, err := ctrl.friendService.GetFriends(userID)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, gin.H{"friends": friends})
}

// @Summary 删除好友
// @Tags 好友管理
// @Security BearerAuth
// @Param id path int true "好友ID"
// @Success 200 {object} utils.Response
// @Router /friend/{id} [delete]
func (ctrl *FriendController) DeleteFriend(c *gin.Context) {
	userID := c.GetUint("userID")

	friendID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ValidationError(c, "无效的好友ID")
		return
	}

	if err := ctrl.friendService.DeleteFriend(userID, uint(friendID)); err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, gin.H{"message": "删除成功"})
}

// @Summary 获取好友请求列表
// @Tags 好友管理
// @Security BearerAuth
// @Success 200 {object} utils.Response
// @Router /friend/requests [get]
func (ctrl *FriendController) GetFriendRequests(c *gin.Context) {
	userID := c.GetUint("userID")

	requests, err := ctrl.friendService.GetFriendRequests(userID)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, gin.H{"requests": requests})
}

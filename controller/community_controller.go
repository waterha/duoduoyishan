package controller

import (
	"duoduoyishan/service"
	"duoduoyishan/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CommunityController struct {
	communityService *service.CommunityService
}

func NewCommunityController() *CommunityController {
	return &CommunityController{
		communityService: service.NewCommunityService(),
	}
}

// @Summary 创建社区
// @Tags 社区管理
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body createCommunityRequest true "社区信息"
// @Success 200 {object} utils.Response
// @Router /community/create [post]
func (ctrl *CommunityController) CreateCommunity(c *gin.Context) {
	userID := c.GetUint("userID")

	var req struct {
		Name        string `json:"name" binding:"required,min=2,max=100"`
		Description string `json:"description" binding:"max=500"`
		Category    string `json:"category" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, "参数错误："+err.Error())
		return
	}

	community, err := ctrl.communityService.CreateCommunity(userID, req.Name, req.Description, req.Category)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, community)
}

// @Summary 加入社区
// @Tags 社区管理
// @Security BearerAuth
// @Param id path int true "社区ID"
// @Success 200 {object} utils.Response
// @Router /community/join/{id} [post]
func (ctrl *CommunityController) JoinCommunity(c *gin.Context) {
	userID := c.GetUint("userID")

	communityID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ValidationError(c, "无效的社区ID")
		return
	}

	if err := ctrl.communityService.JoinCommunity(userID, uint(communityID)); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, gin.H{"message": "加入成功"})
}

// @Summary 退出社区
// @Tags 社区管理
// @Security BearerAuth
// @Param id path int true "社区ID"
// @Success 200 {object} utils.Response
// @Router /community/quit/{id} [post]
func (ctrl *CommunityController) QuitCommunity(c *gin.Context) {
	userID := c.GetUint("userID")

	communityID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ValidationError(c, "无效的社区ID")
		return
	}

	if err := ctrl.communityService.QuitCommunity(userID, uint(communityID)); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, gin.H{"message": "退出成功"})
}

// @Summary 获取社区列表
// @Tags 社区管理
// @Param category query string false "社区分类"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} utils.Response
// @Router /community/list [get]
func (ctrl *CommunityController) GetCommunityList(c *gin.Context) {
	category := c.Query("category")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	communities, total, err := ctrl.communityService.GetCommunityList(category, page, pageSize)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, gin.H{
		"list":      communities,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// @Summary 获取社区详情
// @Tags 社区管理
// @Param id path int true "社区ID"
// @Success 200 {object} utils.Response
// @Router /community/detail/{id} [get]
func (ctrl *CommunityController) GetCommunityDetail(c *gin.Context) {
	communityID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ValidationError(c, "无效的社区ID")
		return
	}

	community, members, err := ctrl.communityService.GetCommunityDetail(uint(communityID))
	if err != nil {
		utils.NotFound(c, "社区不存在")
		return
	}

	utils.Success(c, gin.H{
		"community":    community,
		"members":      members,
		"member_count": len(members),
	})
}

// @Summary 获取用户加入的社区
// @Tags 社区管理
// @Security BearerAuth
// @Success 200 {object} utils.Response
// @Router /community/my [get]
func (ctrl *CommunityController) GetMyCommunities(c *gin.Context) {
	userID := c.GetUint("userID")

	communities, err := ctrl.communityService.GetUserCommunities(userID)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, gin.H{"communities": communities})
}

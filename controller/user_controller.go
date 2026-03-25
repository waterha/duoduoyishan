package controller

import (
	"duoduoyishan/service"
	"duoduoyishan/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService *service.UserService
}

func NewUserController() *UserController {
	return &UserController{
		userService: service.NewUserService(),
	}
}

// @Summary 用户注册
// @Tags 用户认证
// @Accept json
// @Produce json
// @Param request body registerRequest true "注册信息"
// @Success 200 {object} utils.Response
// @Router /auth/register [post]
func (ctrl *UserController) Register(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required,min=3,max=50"`
		Password string `json:"password" binding:"required,min=6,max=50"`
		Email    string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, "参数错误："+err.Error())
		return
	}

	user, err := ctrl.userService.Register(req.Username, req.Password, req.Email)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, gin.H{
		"id":         user.ID,
		"username":   user.Username,
		"nickname":   user.Nickname,
		"email":      user.Email,
		"created_at": user.CreatedAt,
	})
}

// @Summary 用户登录
// @Tags 用户认证
// @Accept json
// @Produce json
// @Param request body loginRequest true "登录信息"
// @Success 200 {object} utils.Response
// @Router /auth/login [post]
func (ctrl *UserController) Login(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, "参数错误："+err.Error())
		return
	}

	token, user, err := ctrl.userService.Login(req.Username, req.Password)
	if err != nil {
		utils.Error(c, http.StatusUnauthorized, err.Error())
		return
	}

	utils.Success(c, gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"nickname": user.Nickname,
			"avatar":   user.Avatar,
			"email":    user.Email,
			"status":   user.Status,
		},
	})
}

// @Summary 退出登录
// @Tags 用户认证
// @Security BearerAuth
// @Success 200 {object} utils.Response
// @Router /auth/logout [post]
func (ctrl *UserController) Logout(c *gin.Context) {
	token := c.GetString("token")
	userID := c.GetUint("userID")

	if err := ctrl.userService.Logout(token, userID); err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, nil)
}

// @Summary 获取当前用户信息
// @Tags 用户管理
// @Security BearerAuth
// @Success 200 {object} utils.Response
// @Router /user/info [get]
func (ctrl *UserController) GetUserInfo(c *gin.Context) {
	userID := c.GetUint("userID")

	user, err := ctrl.userService.GetUserByID(userID)
	if err != nil {
		utils.NotFound(c, "用户不存在")
		return
	}

	utils.Success(c, gin.H{
		"id":         user.ID,
		"username":   user.Username,
		"nickname":   user.Nickname,
		"avatar":     user.Avatar,
		"email":      user.Email,
		"phone":      user.Phone,
		"gender":     user.Gender,
		"birthday":   user.Birthday,
		"signature":  user.Signature,
		"status":     user.Status,
		"created_at": user.CreatedAt,
	})
}

// @Summary 更新用户信息
// @Tags 用户管理
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body updateUserRequest true "用户信息"
// @Success 200 {object} utils.Response
// @Router /user/info [put]
func (ctrl *UserController) UpdateUserInfo(c *gin.Context) {
	userID := c.GetUint("userID")

	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, "参数错误："+err.Error())
		return
	}

	user, err := ctrl.userService.UpdateUserInfo(userID, req)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, user)
}

// @Summary 修改密码
// @Tags 用户管理
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body changePasswordRequest true "密码信息"
// @Success 200 {object} utils.Response
// @Router /user/password [put]
func (ctrl *UserController) ChangePassword(c *gin.Context) {
	userID := c.GetUint("userID")

	var req struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, "参数错误："+err.Error())
		return
	}

	if err := ctrl.userService.ChangePassword(userID, req.OldPassword, req.NewPassword); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, gin.H{"message": "密码修改成功"})
}

// @Summary 搜索用户
// @Tags 用户管理
// @Security BearerAuth
// @Param keyword query string false "搜索关键词"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} utils.Response
// @Router /user/search [get]
func (ctrl *UserController) SearchUsers(c *gin.Context) {
	keyword := c.Query("keyword")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	users, total, err := ctrl.userService.SearchUsers(keyword, page, pageSize)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, gin.H{
		"list":      users,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

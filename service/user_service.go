package service

import (
	"duoduoyishan/cache"
	"duoduoyishan/database"
	"duoduoyishan/models"
	"duoduoyishan/utils"
	"errors"
	"time"
)

type UserService struct{}

func NewUserService() *UserService {
	return &UserService{}
}

// 用户注册
func (s *UserService) Register(username, password, email string) (*models.User, error) {
	// 检查用户名是否已存在
	var count int64
	if err := database.DB.Model(&models.User{}).Where("username = ?", username).Count(&count).Error; err != nil {
		utils.Logger.Errorf("检查用户名是否存在失败: %v", err)
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("用户名已存在")
	}

	// 检查邮箱是否已存在
	if err := database.DB.Model(&models.User{}).Where("email = ?", email).Count(&count).Error; err != nil {
		utils.Logger.Errorf("检查邮箱是否存在失败: %v", err)
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("邮箱已注册")
	}

	// 加密密码
	hashedPassword, err := utils.EncryptPassword(password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Username: username,
		Password: hashedPassword,
		Email:    email,
		Nickname: username,
		Status:   2, // 离线
	}

	if err := database.DB.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

// 用户登录
func (s *UserService) Login(username, password string) (string, *models.User, error) {
	var user models.User
	if err := database.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return "", nil, errors.New("用户不存在")
	}

	if !utils.CheckPassword(password, user.Password) {
		return "", nil, errors.New("密码错误")
	}

	// 生成JWT token
	token, err := utils.GenerateToken(user.ID, user.Username)
	if err != nil {
		utils.Logger.Errorf("生成JWT token失败: %v", err)
		return "", nil, err
	}

	// 保存session到Redis
	if err := cache.SetUserSession(token, user.ID); err != nil {
		utils.Logger.Errorf("保存session到Redis失败: %v", err)
		return "", nil, err
	}

	// 更新登录信息
	updates := map[string]interface{}{
		"last_login_at": time.Now(),
		"status":        1,
	}
	if err := database.DB.Model(&user).Updates(updates).Error; err != nil {
		utils.Logger.Errorf("更新登录信息失败: %v", err)
		// 这里不返回错误，因为登录已经成功，只是更新登录信息失败
	}

	// 缓存用户在线状态
	cache.SetUserOnline(user.ID, true)

	return token, &user, nil
}

// 退出登录
func (s *UserService) Logout(token string, userID uint) error {
	// 删除session
	if err := cache.DeleteUserSession(token); err != nil {
		utils.Logger.Errorf("删除session失败: %v", err)
		return err
	}

	// 更新状态
	if err := database.DB.Model(&models.User{}).Where("id = ?", userID).Update("status", 2).Error; err != nil {
		utils.Logger.Errorf("更新用户状态失败: %v", err)
		// 这里不返回错误，因为删除session已经成功
	}

	// 清除在线状态
	cache.SetUserOnline(userID, false)

	return nil
}

// 获取用户信息
func (s *UserService) GetUserByID(userID uint) (*models.User, error) {
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		utils.Logger.Errorf("获取用户信息失败: %v", err)
		return nil, err
	}

	// 获取在线状态
	user.Status = 2
	if cache.IsUserOnline(user.ID) {
		user.Status = 1
	}

	return &user, nil
}

// 更新用户信息
func (s *UserService) UpdateUserInfo(userID uint, updates map[string]interface{}) (*models.User, error) {
	// 过滤可更新字段
	allowedFields := []string{"nickname", "avatar", "phone", "gender", "birthday", "signature"}

	data := make(map[string]interface{})
	for _, field := range allowedFields {
		if val, ok := updates[field]; ok {
			data[field] = val
		}
	}

	if len(data) > 0 {
		if err := database.DB.Model(&models.User{}).Where("id = ?", userID).Updates(data).Error; err != nil {
			utils.Logger.Errorf("更新用户信息失败: %v", err)
			return nil, err
		}
	}

	return s.GetUserByID(userID)
}

// 修改密码
func (s *UserService) ChangePassword(userID uint, oldPassword, newPassword string) error {
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		return err
	}

	if !utils.CheckPassword(oldPassword, user.Password) {
		return errors.New("原密码错误")
	}

	hashedPassword, err := utils.EncryptPassword(newPassword)
	if err != nil {
		return err
	}

	return database.DB.Model(&user).Update("password", hashedPassword).Error
}

// 搜索用户
func (s *UserService) SearchUsers(keyword string, page, pageSize int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	db := database.DB.Model(&models.User{})

	if keyword != "" {
		db = db.Where("username LIKE ? OR nickname LIKE ? OR email LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	db.Count(&total)

	offset := (page - 1) * pageSize
	err := db.Offset(offset).Limit(pageSize).Find(&users).Error

	return users, total, err
}

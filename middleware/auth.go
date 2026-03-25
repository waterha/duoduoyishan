package middleware

import (
	"duoduoyishan/cache"
	"duoduoyishan/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从Header获取token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.Unauthorized(c, "未提供认证令牌")
			c.Abort()
			return
		}

		// 处理Bearer token
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			utils.Unauthorized(c, "认证格式错误")
			c.Abort()
			return
		}

		// 解析token
		claims, err := utils.ParseToken(parts[1])
		if err != nil {
			utils.Unauthorized(c, "无效的认证令牌")
			c.Abort()
			return
		}

		// 检查session是否有效
		userID, err := cache.GetUserSession(parts[1])
		if err != nil || userID != claims.UserID {
			utils.Unauthorized(c, "会话已过期")
			c.Abort()
			return
		}

		// 将用户信息存入上下文
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("token", parts[1])

		// 更新用户在线状态
		cache.RefreshUserOnline(claims.UserID)

		c.Next()
	}
}

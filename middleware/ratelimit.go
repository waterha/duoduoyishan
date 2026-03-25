package middleware

import (
	"duoduoyishan/cache"
	"duoduoyishan/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// 简单IP限流中间件
func RateLimit(limit int, duration time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		key := "rate_limit:" + ip

		// 获取当前计数
		count, err := cache.RedisClient.Get(c, key).Int()
		if err != nil {
			// Key不存在，设置初始值
			cache.RedisClient.Set(c, key, 1, duration)
			c.Next()
			return
		}

		if count >= limit {
			utils.Error(c, http.StatusTooManyRequests, "请求过于频繁，请稍后再试")
			c.Abort()
			return
		}

		// 增加计数
		cache.RedisClient.Incr(c, key)
		c.Next()
	}
}

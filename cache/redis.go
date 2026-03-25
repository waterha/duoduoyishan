package cache

import (
	"context"
	"duoduoyishan/config"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"strconv"
	"time"
)

var RedisClient *redis.Client
var ctx = context.Background()

func InitRedis() error {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%s", config.GlobalConfig.Redis.Host, config.GlobalConfig.Redis.Port),
		Password:     config.GlobalConfig.Redis.Password,
		DB:           config.GlobalConfig.Redis.DB,
		PoolSize:     config.GlobalConfig.Redis.PoolSize,
		MinIdleConns: 5,
		IdleTimeout:  5 * time.Minute,
	})

	// 测试连接
	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("Redis连接失败: %v", err)
	}

	log.Println("Redis连接成功")
	return nil
}

// 辅助函数：uint转string
func UintToString(u uint) string {
	return strconv.FormatUint(uint64(u), 10)
}

// 用户在线状态
func SetUserOnline(userID uint, online bool) {
	key := "user:online:" + UintToString(userID)
	if online {
		RedisClient.Set(ctx, key, 1, 30*time.Minute)
	} else {
		RedisClient.Del(ctx, key)
	}
}

func IsUserOnline(userID uint) bool {
	key := "user:online:" + UintToString(userID)
	val, err := RedisClient.Get(ctx, key).Result()
	return err == nil && val == "1"
}

// 更新用户在线状态（心跳）
func RefreshUserOnline(userID uint) {
	key := "user:online:" + UintToString(userID)
	RedisClient.Expire(ctx, key, 30*time.Minute)
}

// 获取用户所有在线好友
func GetUserOnlineFriends(userID uint, friendIDs []uint) []uint {
	var onlineFriends []uint
	for _, friendID := range friendIDs {
		if IsUserOnline(friendID) {
			onlineFriends = append(onlineFriends, friendID)
		}
	}
	return onlineFriends
}

// 缓存用户会话
func SetUserSession(token string, userID uint) error {
	key := "session:" + token
	return RedisClient.Set(ctx, key, userID, config.GlobalConfig.JWT.ExpireTime).Err()
}

func GetUserSession(token string) (uint, error) {
	key := "session:" + token
	val, err := RedisClient.Get(ctx, key).Uint64()
	if err != nil {
		return 0, err
	}
	return uint(val), nil
}

func DeleteUserSession(token string) error {
	key := "session:" + token
	return RedisClient.Del(ctx, key).Err()
}

// 未读消息计数
func GetUnreadCount(userID uint, roomID string) int {
	key := "unread:" + UintToString(userID) + ":" + roomID
	val, err := RedisClient.Get(ctx, key).Int()
	if err != nil {
		return 0
	}
	return val
}

func IncrUnreadCount(userID uint, roomID string) {
	key := "unread:" + UintToString(userID) + ":" + roomID
	RedisClient.Incr(ctx, key)
	RedisClient.Expire(ctx, key, 7*24*time.Hour)
}

func ClearUnreadCount(userID uint, roomID string) {
	key := "unread:" + UintToString(userID) + ":" + roomID
	RedisClient.Del(ctx, key)
}

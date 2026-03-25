// main.go
package main

import (
	"duoduoyishan/cache"
	"duoduoyishan/config"
	"duoduoyishan/database"
	"duoduoyishan/router"
	"duoduoyishan/utils"
	"duoduoyishan/websocket_own"
	"log"
)

func main() {
	// 1. 初始化配置
	if err := config.InitConfig(); err != nil {
		log.Fatal("初始化配置失败:", err)
	}

	// 2. 初始化日志
	if err := utils.InitLogger(); err != nil {
		log.Fatal("初始化日志失败:", err)
	}

	// 3. 初始化MySQL
	if err := database.InitMySQL(); err != nil {
		utils.Logger.Fatal("初始化MySQL失败:", err)
	}

	// 4. 自动迁移数据库表
	if err := database.AutoMigrate(); err != nil {
		utils.Logger.Fatal("数据库迁移失败:", err)
	}

	// 5. 初始化Redis
	if err := cache.InitRedis(); err != nil {
		utils.Logger.Fatal("初始化Redis失败:", err)
	}

	// 6. 创建WebSocket hub并启动
	hub := websocket_own.NewHub()
	go hub.Run()

	// 7. 初始化路由
	r := router.InitRouter(hub)

	// 8. 启动服务器
	utils.Logger.Infof("服务器启动在端口 %s", config.GlobalConfig.Server.Port)
	if err := r.Run(":" + config.GlobalConfig.Server.Port); err != nil {
		utils.Logger.Fatal("启动服务器失败:", err)
	}
}

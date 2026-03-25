package router

import (
	"duoduoyishan/controller"
	"duoduoyishan/middleware"
	"duoduoyishan/websocket_own"

	"github.com/gin-gonic/gin"
)

func InitRouter(hub *websocket_own.Hub) *gin.Engine {
	// 设置运行模式
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()

	// 使用中间件
	r.Use(gin.Recovery())
	r.Use(middleware.Logger())
	r.Use(middleware.Cors())

	// 静态文件服务
	r.Static("/static", "./static")
	r.Static("/uploads", "./uploads")

	// 首页
	r.GET("/", func(c *gin.Context) {
		c.File("./static/index.html")
	})

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// 初始化控制器
	userController := controller.NewUserController()
	friendController := controller.NewFriendController()
	communityController := controller.NewCommunityController()
	messageController := controller.NewMessageController()
	wsController := controller.NewWSController(hub)

	// API路由组
	api := r.Group("/api")
	{
		// 公开接口
		auth := api.Group("/auth")
		{
			auth.POST("/register", userController.Register)
			auth.POST("/login", userController.Login)
		}

		// 需要认证的接口
		authorized := api.Group("/")
		authorized.Use(middleware.JWTAuth())
		{
			// 认证相关
			auth := authorized.Group("/auth")
			{
				auth.POST("/logout", userController.Logout)
			}

			// 用户相关
			user := authorized.Group("/user")
			{
				user.GET("/info", userController.GetUserInfo)
				user.PUT("/info", userController.UpdateUserInfo)
				user.PUT("/password", userController.ChangePassword)
				user.GET("/search", userController.SearchUsers)
			}

			// 好友相关
			friend := authorized.Group("/friend")
			{
				friend.POST("/request", friendController.SendFriendRequest)
				friend.PUT("/request/:id", friendController.HandleFriendRequest)
				friend.GET("/list", friendController.GetFriends)
				friend.DELETE("/:id", friendController.DeleteFriend)
				friend.GET("/requests", friendController.GetFriendRequests)
			}

			// 社区相关
			community := authorized.Group("/community")
			{
				community.POST("/create", communityController.CreateCommunity)
				community.POST("/join/:id", communityController.JoinCommunity)
				community.POST("/quit/:id", communityController.QuitCommunity)
				community.GET("/list", communityController.GetCommunityList)
				community.GET("/detail/:id", communityController.GetCommunityDetail)
				community.GET("/my", communityController.GetMyCommunities)
			}

			// 消息相关
			message := authorized.Group("/message")
			{
				message.GET("/history", messageController.GetChatHistory)
				message.GET("/unread", messageController.GetUnreadCount)
				message.PUT("/read/:id", messageController.MarkMessageRead)
				message.PUT("/recall/:id", messageController.RecallMessage)
			}

		// WebSocket
		api.GET("/ws", wsController.Connect)
		api.GET("/ws/online", wsController.GetRoomOnlineCount)
		}
	}

	return r
}

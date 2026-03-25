package database

import (
	"duoduoyishan/config"
	"duoduoyishan/models"
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitMySQL() error {
	// 构建DSN
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		config.GlobalConfig.Database.Username,
		config.GlobalConfig.Database.Password,
		config.GlobalConfig.Database.Host,
		config.GlobalConfig.Database.Port,
		config.GlobalConfig.Database.Database,
		config.GlobalConfig.Database.Charset,
	)
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NowFunc: func() time.Time {
			return time.Now().Local()
		},
	})
	if err != nil {
		return fmt.Errorf("连接数据库失败: %v", err)
	}

	// 获取通用数据库对象
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("获取数据库对象失败: %v", err)
	}

	// 连接池配置
	sqlDB.SetMaxIdleConns(config.GlobalConfig.Database.MaxIdleConns)
	sqlDB.SetMaxOpenConns(config.GlobalConfig.Database.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("MySQL连接成功")
	return nil
}

// 自动迁移表结构
func AutoMigrate() error {
	err := DB.AutoMigrate(
		&models.User{},
		&models.Community{},
		&models.CommunityMember{},
		&models.Message{},
		&models.MessageRead{},
		&models.ChatRoom{},
		&models.RoomUser{},
		&models.Friend{},
		&models.FriendRequest{},
	)
	if err != nil {
		return fmt.Errorf("自动迁移失败: %v", err)
	}

	log.Println("数据库表迁移完成")
	return nil
}

package utils

import (
	"duoduoyishan/config"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger

func InitLogger() error {
	Logger = logrus.New()

	// 设置日志级别
	level, err := logrus.ParseLevel(config.GlobalConfig.Log.Level)
	if err != nil {
		level = logrus.InfoLevel
	}
	Logger.SetLevel(level)

	// 设置输出格式
	Logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
		PrettyPrint:     false,
	})

	// 创建日志目录
	logDir := filepath.Dir(config.GlobalConfig.Log.Filename)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return err
	}

	// 配置日志轮转
	lumberjackLogger := &lumberjack.Logger{
		Filename:   config.GlobalConfig.Log.Filename,
		MaxSize:    config.GlobalConfig.Log.MaxSize,
		MaxBackups: config.GlobalConfig.Log.MaxBackups,
		MaxAge:     config.GlobalConfig.Log.MaxAge,
		Compress:   true,
	}

	// 同时输出到文件和控制台
	multiWriter := io.MultiWriter(os.Stdout, lumberjackLogger)
	Logger.SetOutput(multiWriter)

	return nil
}

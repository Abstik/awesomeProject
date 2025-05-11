package utils

import (
	"os"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Logger 通用日志，以天为分割，中间件和正常程序处理都用这一个Logger
var Logger *logrus.Logger

func InitLogger() {
	Logger = logrus.New()
	Logger.SetLevel(logrus.InfoLevel)

	// 确保logs目录存在
	if err := os.MkdirAll("logs", 0755); err != nil {
		Logger.Errorf("无法创建日志目录: %v", err)
		return
	}

	logFile := filepath.Join("logs", time.Now().Format("2006-01-02")+".log")

	// 测试日志文件是否可写
	if _, err := os.OpenFile(logFile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666); err != nil {
		Logger.Errorf("无法打开日志文件: %v", err)
		return
	}

	Logger.SetOutput(&lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    10,
		MaxBackups: 0,
		MaxAge:     0,
		Compress:   false,
	})

	Logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	// 测试日志记录
	Logger.Info("日志初始化成功")
}

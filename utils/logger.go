package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Logger *zap.Logger
var currentLogDate string

// 初始化Logger
func InitLogger() {
	ensureLogsDir()
	cleanupOldLogs(30)
	setLoggerForToday()
	Logger.Info("日志初始化成功")
}

func ensureLogsDir() {
	if err := os.MkdirAll("logs", 0755); err != nil {
		fmt.Printf("创建日志目录失败: %v\n", err)
	}
}

// 删除30天前的日志
func cleanupOldLogs(days int) {
	files, _ := os.ReadDir("logs")
	expire := time.Now().AddDate(0, 0, -days)
	for _, f := range files {
		name := f.Name()
		if len(name) == len("2006-01-02.log") {
			dateStr := name[:10]
			t, err := time.Parse("2006-01-02", dateStr)
			if err == nil && t.Before(expire) {
				os.Remove(filepath.Join("logs", name))
			}
		}
	}
}

// 每天切换日志文件
func setLoggerForToday() {
	today := time.Now().Format("2006-01-02")
	if currentLogDate == today && Logger != nil {
		return
	}

	// 日志文件路径
	logDir := "logs"
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		_ = os.MkdirAll(logDir, 0755) // 确保目录存在
	}
	filename := filepath.Join(logDir, today+".log")

	// lumberjack日志滚动配置
	writer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   filename,
		MaxSize:    50,    // 单文件最大50MB
		MaxBackups: 0,     // 不保留历史备份
		MaxAge:     7,     // 日志保留天数
		Compress:   false, // 不压缩
	})

	// 日志编码器配置
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 区分环境，生产用JSON，开发用Console
	var encoder zapcore.Encoder
	if os.Getenv("GO_ENV") == "debug" {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	}

	// 日志级别，默认Info
	level := zapcore.InfoLevel

	// 创建核心Core
	core := zapcore.NewCore(encoder, writer, level)

	// 开发环境额外输出到终端
	if os.Getenv("GO_ENV") == "debug" {
		consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
		consoleCore := zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), zapcore.DebugLevel)
		core = zapcore.NewTee(core, consoleCore)
	}

	// 替换全局Logger
	Logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	currentLogDate = today
}

// 跨天时确保日志切换
func EnsureZapLoggerReady() {
	setLoggerForToday()
}

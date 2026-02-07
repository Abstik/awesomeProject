package utils

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Logger *zap.Logger

// InitLogger 初始化 Logger，由 lumberjack 统一管理日志轮转
func InitLogger() {
	// lumberjack 负责按大小轮转和自动清理
	writer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "logs/app.log",
		MaxSize:    50,   // 单文件最大 50MB，超过后自动轮转
		MaxBackups: 30,   // 最多保留 30 个历史文件
		MaxAge:     30,   // 保留 30 天
		Compress:   true, // 压缩历史日志
	})

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

	// 区分环境：生产用 JSON，开发用 Console
	var encoder zapcore.Encoder
	if os.Getenv("GO_ENV") == "debug" {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	}

	core := zapcore.NewCore(encoder, writer, zapcore.InfoLevel)

	// 开发环境额外输出到终端
	if os.Getenv("GO_ENV") == "debug" {
		consoleCore := zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderConfig),
			zapcore.Lock(os.Stdout),
			zapcore.DebugLevel,
		)
		core = zapcore.NewTee(core, consoleCore)
	}

	Logger = zap.New(core, zap.AddCaller())
	zap.ReplaceGlobals(Logger)
	Logger.Info("日志初始化成功")
}

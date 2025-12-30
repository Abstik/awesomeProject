package dao

import (
	"os"
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

// 开启数据库连接
func InitDatabaseConnector() error {
	dsn := "root:325523@tcp(127.0.0.1:3306)/xiyoumobile_data?charset=utf8mb4&parseTime=True&loc=Local"
	// todo 部署到服务器后，使用以下配置
	//dsn := fmt.Sprintf("root:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True",
	//	os.Getenv("DB_PASSWORD"), // 密码
	//	os.Getenv("DB_HOST"),     // 主机（必须改为 "mysql"）
	//	os.Getenv("DB_PORT"),
	//	os.Getenv("MYSQL_DATABASE"))

	// 根据环境变量设置GORM日志等级
	var logLevel logger.LogLevel
	switch os.Getenv("GO_ENV") {
	case "debug":
		logLevel = logger.Info
	case "release":
		logLevel = logger.Error
	default:
		logLevel = logger.Warn
	}

	// 用zap.Logger包装成gorm logger接口
	zapLogger := logger.New(
		zap.NewStdLog(zap.L()), // zap.L() 是全局Logger
		logger.Config{
			SlowThreshold:             2 * time.Second, // 慢查询阈值
			LogLevel:                  logLevel,        // 日志等级
			IgnoreRecordNotFoundError: true,            // 忽略record not found错误日志
			Colorful:                  false,           // 关闭终端颜色
		},
	)

	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: zapLogger,
	})
	return err
}

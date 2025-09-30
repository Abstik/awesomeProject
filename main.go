package main

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-contrib/cors"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.uber.org/zap"

	"awesomeProject/dao"
	"awesomeProject/middleware"
	"awesomeProject/router"
	"awesomeProject/utils"
)

func main() {
	// 加载上级目录的.env文件，读取环境变量(仅在开发环境生效)
	cwd, _ := os.Getwd()
	parentEnvPath := filepath.Join(cwd, "..", ".env")
	if err := godotenv.Load(parentEnvPath); err != nil && !(os.Getenv("GO_ENV") == "release") {
		log.Printf("未加载 %s 文件，使用系统环境变量，错误: %v\n", parentEnvPath, err)
	}

	// 设置gin的模式
	env := os.Getenv("GO_ENV")
	switch env {
	case "debug", "release", "test":
		gin.SetMode(env)
	default:
		gin.SetMode(gin.DebugMode)
	}

	// 初始化日志
	utils.InitLogger()
	// 替换全局zap全局Logger
	zap.ReplaceGlobals(utils.Logger)

	// 初始化MySQL连接
	err := dao.InitDatabaseConnector()
	if err != nil {
		panic("db init failed error: " + err.Error())
	}

	// 初始化Gin
	r := gin.New()
	// 使用 zap 替代gin日志的中间件
	r.Use(
		ginzap.Ginzap(zap.L(), time.RFC3339, true), // 访问日志
		ginzap.RecoveryWithZap(zap.L(), true),      // panic恢复日志
	)
	// 设置日志中间件
	r.Use(middleware.ZapLoggerMiddleware())
	r.Use(cors.New(cors.Config{
		// 允许的域名（前端地址）
		AllowOrigins: []string{"*"}, // 允许所有源
		// 允许的请求方法
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		// 允许的请求头
		AllowHeaders: []string{"Content-Type", "Authorization"},
		// 允许携带认证信息
		AllowCredentials: true,
	}))

	// 配置静态资源访问
	r.Static("/res", "./res")
	r.Static("/videos", "./videos")

	// 设置路由
	r = router.SetupRouter(r)

	// 启动服务
	if err := r.Run(":8080"); err != nil {
		panic("gin run err, error is " + err.Error())
	}
}

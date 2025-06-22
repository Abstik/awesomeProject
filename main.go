package main

import (
	"github.com/gin-gonic/gin"

	"awesomeProject/dao"
	"awesomeProject/middleware"
	"awesomeProject/router"
	"awesomeProject/utils"
)

func main() {
	// 初始化日志
	utils.InitLogger()

	// 初始化MySQL连接
	err := dao.InitDatabaseConnector()
	if err != nil {
		panic("db init failed error: " + err.Error())
	}

	// 初始化Gin
	r := gin.Default()

	// 配置静态资源访问
	r.Static("/res", "./res")

	// 设置日志中间件
	r.Use(middleware.SetLoggerMiddleware())

	// 设置路由
	r = router.SetupRouter(r)

	// 启动服务
	if err := r.Run(":8080"); err != nil {
		panic("gin run err, error is " + err.Error())
	}
}

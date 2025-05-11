package main

import (
	"github.com/gin-gonic/gin"

	"awesomeProject/dao"
	"awesomeProject/middlewire"
	"awesomeProject/router"
	"awesomeProject/utils"
)

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

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
	r.Use(middlewire.SetLoggerMiddleware())
	r = router.SetupRouter(r)
	if err := r.Run(":8080"); err != nil {
		panic("gin run err, error is " + err.Error())
	}
}

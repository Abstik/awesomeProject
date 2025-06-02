package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"awesomeProject/handler"
	"awesomeProject/middleware"
)

func SetupRouter(router *gin.Engine) *gin.Engine {
	SetupRegisterRouter(router)
	return router
}

func SetupRegisterRouter(r *gin.Engine) {
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

	apiGroup := r.Group("/api")
	{
		// 生成验证码
		apiGroup.GET("/captcha", handler.GetCaptcha)
		// 校验验证码
		//apiGroup.GET("/verifyCaptcha", handler.VerifyCaptcha)
	}

	{
		// 用户登录
		apiGroup.POST("/login", handler.Login)
		// 根据用户名修改个人信息
		apiGroup.POST("/changeinfo", middleware.JWTAuthMiddleware(), handler.ChangeMemberInfo)
		// 根据用户名查询个人信息
		apiGroup.GET("/userinfo", middleware.JWTAuthMiddleware(), handler.GetMemberByName)
	}

	{
		// 批量获取成员信息
		apiGroup.GET("/members", handler.GetMemberList)
		// 添加成员
		apiGroup.POST("/members", middleware.JWTAuthMiddleware(), middleware.IsAdminAuthMiddleware(), handler.Register)
		// 查询毕业生的所有年份
		apiGroup.GET("/years", handler.GetYears)
		// 删除成员
		apiGroup.DELETE("/members", middleware.JWTAuthMiddleware(), middleware.IsAdminAuthMiddleware(), handler.DeleteMember)
	}

	{
		// 获取活动数据
		apiGroup.GET("/activities/list", handler.GetActivityList)
		// 查询活动（根据id）
		apiGroup.GET("/activity", handler.GetActivityById)
		// 上传活动
		apiGroup.POST("/activity", middleware.JWTAuthMiddleware(), middleware.IsAdminAuthMiddleware(), handler.AddActivity)
		// 修改活动
		apiGroup.PUT("/activity", middleware.JWTAuthMiddleware(), middleware.IsAdminAuthMiddleware(), handler.UpdateActivity)
		// 删除活动
		apiGroup.DELETE("/activity", middleware.JWTAuthMiddleware(), middleware.IsAdminAuthMiddleware(), handler.DeleteActivity)
	}

	{
		// 增加实验室方向
		apiGroup.POST("/team", middleware.JWTAuthMiddleware(), middleware.IsAdminAuthMiddleware(), handler.AddTeam)
		// 查询实验室方向（游客）
		apiGroup.GET("/team", handler.GetTeams)
		// 查询实验室方向（用户和管理员）
		apiGroup.GET("/team/allinfo", middleware.JWTAuthMiddleware(), handler.GetTeams)
		// 修改实验室方向
		apiGroup.PUT("/team", middleware.JWTAuthMiddleware(), middleware.IsAdminAuthMiddleware(), handler.UpdateTeam)
	}

	{
		// 上传捐款信息
		apiGroup.POST("/donation/list", middleware.JWTAuthMiddleware(), middleware.IsAdminAuthMiddleware(), handler.AddDonations)
		// 查询捐款信息
		apiGroup.GET("/donation/list", middleware.JWTAuthMiddleware(), handler.GetDonations)
		// 删除捐款信息
		apiGroup.DELETE("/donation", middleware.JWTAuthMiddleware(), middleware.IsAdminAuthMiddleware(), handler.DeleteDonation)
	}

	{
		// 上传图片附带水印
		apiGroup.POST("/uploadImgWithWaterMark", middleware.JWTAuthMiddleware(), handler.UploadImgWithWaterMark)
		// 删除图片
		apiGroup.DELETE("/deleteImg", middleware.JWTAuthMiddleware(), handler.DeleteImg)
	}

	{
		// 上传qq纳新群
		apiGroup.POST("/contact", handler.UpdateContact)
		// 获取联系方式
		apiGroup.GET("/contact", handler.ContactWithUs)
	}

	{
		// 上传培养计划
		apiGroup.POST("/trainplan", middleware.JWTAuthMiddleware(), middleware.IsAdminAuthMiddleware(), handler.UpdateTrainPlan)
		// 获取培养计划
		apiGroup.GET("/trainplan", handler.GetTrainPlan)
	}
}

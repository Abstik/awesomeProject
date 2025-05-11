package router

import (
	"github.com/gin-gonic/gin"

	"awesomeProject/handler"
)

func SetupRouter(router *gin.Engine) *gin.Engine {
	SetupRegisterRouter(router)
	return router
}

func SetupRegisterRouter(r *gin.Engine) {
	apiGroup := r.Group("/api")
	{
		// 用户注册
		apiGroup.POST("/register", handler.Register)
		// 修改个人信息
		apiGroup.POST("/changeInfo", handler.ChangeMemberInfo)
		// 查询个人信息
		apiGroup.GET("/userInfo", handler.GetMemberByName)
	}
	{
		// 生成验证码
		apiGroup.GET("/Captcha", handler.GetCaptcha)
		// 校验验证码
		apiGroup.GET("/VerifyCaptcha", handler.VerifyCaptcha)
	}
	{
		// 查询活动
		apiGroup.GET("/activity", handler.GetActivityById)
		// 上传活动
		apiGroup.POST("/activity", handler.AddActivity)
		// 获取活动数据
		apiGroup.GET("/activities/list", handler.GetActivityList)
	}
	{
		// 获取成员信息
		apiGroup.GET("/members", handler.GetMemberList)
		// 添加成员
		apiGroup.POST("/members", handler.AddMember)
	}
	{
		// 增加实验室方向
		apiGroup.POST("/team", handler.AddTeam)
		// 查询实验室方向
		apiGroup.GET("/team", handler.GetTeams)
	}
	{
		// 上传捐款信息
		apiGroup.POST("/donation", handler.AddDonations)
	}
	{
		// 上传图片附带水印
		apiGroup.POST("/uploadImgWithWaterMark", handler.UploadImgWithWaterMark)
	}
}

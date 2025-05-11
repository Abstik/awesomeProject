package handler

import (
	"github.com/gin-gonic/gin"

	"awesomeProject/model"
	"awesomeProject/service"
	"awesomeProject/utils"
)

func Register(c *gin.Context) {
	var memReq *model.MemberRequest
	// todo 处理下验证码相关问题
	err := c.ShouldBindJSON(&memReq)
	if err != nil {
		utils.Logger.Errorf("handler.Register format error, error is %s", err.Error())
		utils.BuildErrorResponse(c, 500, "Register format error, error is "+err.Error())
		return
	}
	err = service.Register(memReq)
	if err != nil {
		utils.BuildErrorResponse(c, 500, "Register Failed error is "+err.Error())
		return
	}
	utils.BuildSuccessResponse(c, "注册成功")
}

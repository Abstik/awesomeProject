package handler

import (
	"github.com/gin-gonic/gin"

	"awesomeProject/model"
	"awesomeProject/service"
	"awesomeProject/utils"
)

func AddDonations(c *gin.Context) {
	var req model.AddDonationsReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BuildErrorResponse(c, 500, "AddDonations parse failed err is: "+err.Error())
		return
	}

	// 调用 service 层处理
	if err := service.AddDonations(req); err != nil {
		utils.BuildErrorResponse(c, 500, "AddDonations insert failed err is: "+err.Error())
		return
	}

	utils.BuildSuccessResponse(c, "添加成功")
}

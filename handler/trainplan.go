package handler

import (
	"github.com/gin-gonic/gin"

	"awesomeProject/dao"
	"awesomeProject/model"
	"awesomeProject/utils"
)

func UpdateTrainPlan(c *gin.Context) {
	var req model.TrainPlan
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BuildErrorResponse(c, 500, "UpdateTrainPlan parse failed err is: "+err.Error())
		return
	}

	if req.Content == nil {
		utils.BuildErrorResponse(c, 400, "UpdateTrainPlan content is nil")
	}

	err := dao.UpdateTrainPlan(req)
	if err != nil {
		utils.BuildErrorResponse(c, 500, "UpdateTrainPlan failed err is: "+err.Error())
		return
	}
	utils.BuildSuccessResponse(c, "更新成功")
}

func GetTrainPlan(c *gin.Context) {
	trainPlan, err := dao.GetTrainPlan()
	if err != nil {
		utils.BuildErrorResponse(c, 500, "GetTrainPlan failed err is: "+err.Error())
		return
	}
	utils.BuildSuccessResponse(c, trainPlan)
}

package handler

import (
	"errors"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"awesomeProject/dao"
	"awesomeProject/model"
	"awesomeProject/utils"
)

func UpdateTrainPlan(c *gin.Context) {
	var req model.TrainPlanPO
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BuildErrorResponse(c, 400, "参数格式有误")
		return
	}

	err := dao.UpdateTrainPlan(&req)
	if err != nil {
		utils.BuildServerError(c, "更新培养计划失败", err)
		return
	}
	utils.BuildSuccessResponse(c, "更新成功")
}

func GetTrainPlan(c *gin.Context) {
	trainPlan, err := dao.GetTrainPlan()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		utils.BuildServerError(c, "查询培养计划失败", err)
		return
	}
	utils.BuildSuccessResponse(c, trainPlan)
}

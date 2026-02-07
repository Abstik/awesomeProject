package handler

import (
	"errors"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"awesomeProject/dao"
	"awesomeProject/model"
	"awesomeProject/utils"
)

func UpdateIntroduction(c *gin.Context) {
	var introduction *model.IntroductionPO
	err := c.ShouldBindJSON(&introduction)
	if err != nil {
		utils.BuildErrorResponse(c, 400, "参数格式有误")
		return
	}

	if err := dao.UpdateIntroduction(introduction); err != nil {
		utils.BuildServerError(c, "更新简介失败", err)
		return
	}
	utils.BuildSuccessResponse(c, "更新成功")
}

func GetIntroduction(c *gin.Context) {
	introduction, err := dao.GetIntroduction()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		utils.BuildServerError(c, "查询简介失败", err)
		return
	}
	utils.BuildSuccessResponse(c, introduction)
}

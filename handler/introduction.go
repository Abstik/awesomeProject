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
		utils.BuildErrorResponse(c, 400, "UpdateIntroduction format error, error is "+err.Error())
	}

	introduction.Id = 1
	if err := dao.UpdateIntroduction(introduction); err != nil {
		utils.BuildErrorResponse(c, 500, "UpdateIntroduction failed err is: "+err.Error())
	}
	utils.BuildSuccessResponse(c, "更新成功")
}

func GetIntroduction(c *gin.Context) {
	introduction, err := dao.GetIntroduction()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		utils.BuildErrorResponse(c, 500, "查询失败")
		return
	}
	utils.BuildSuccessResponse(c, introduction)
}

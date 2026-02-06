package handler

import (
	"errors"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"awesomeProject/dao"
	"awesomeProject/model"
	"awesomeProject/utils"
)

func ContactWithUs(c *gin.Context) {
	contact, err := dao.GetContact()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		utils.BuildErrorResponse(c, 500, "查询失败")
		return
	}
	utils.BuildSuccessResponse(c, contact)
}

func UpdateContact(c *gin.Context) {
	var contact *model.ContactPO
	err := c.ShouldBindJSON(&contact)
	if err != nil {
		utils.BuildErrorResponse(c, 400, "UpdateContact format error, error is "+err.Error())
	}

	if err := dao.UpdateContact(contact); err != nil {
		utils.BuildErrorResponse(c, 500, "UpdateContact failed err is: "+err.Error())
	}
	utils.BuildSuccessResponse(c, "更新成功")
}

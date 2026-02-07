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
		utils.BuildServerError(c, "查询联系方式失败", err)
		return
	}
	utils.BuildSuccessResponse(c, contact)
}

func UpdateContact(c *gin.Context) {
	var contact *model.ContactPO
	err := c.ShouldBindJSON(&contact)
	if err != nil {
		utils.BuildErrorResponse(c, 400, "参数格式有误")
		return
	}

	if err := dao.UpdateContact(contact); err != nil {
		utils.BuildServerError(c, "更新联系方式失败", err)
		return
	}
	utils.BuildSuccessResponse(c, "更新成功")
}

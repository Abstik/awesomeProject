package handler

import (
	"github.com/gin-gonic/gin"

	"awesomeProject/dao"
	"awesomeProject/model"
	"awesomeProject/utils"
)

func ContactWithUs(c *gin.Context) {
	contact, err := dao.GetContact()
	if err != nil {
		utils.BuildErrorResponse(c, 500, "GetContact failed err is: "+err.Error())
	}
	utils.BuildSuccessResponse(c, contact)
}

func UpdateContact(c *gin.Context) {
	var contact *model.ContactPO
	err := c.ShouldBindJSON(&contact)
	if err != nil {
		utils.BuildErrorResponse(c, 400, "AddContact format error, error is "+err.Error())
	}
	if err := dao.UpdateContact(contact); err != nil {
		utils.BuildErrorResponse(c, 500, "AddContact failed err is: "+err.Error())
	}
	utils.BuildSuccessResponse(c, "更新成功")
}

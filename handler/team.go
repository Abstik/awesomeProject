package handler

import (
	"errors"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"awesomeProject/model"
	"awesomeProject/service"
	"awesomeProject/utils"
)

func AddTeam(c *gin.Context) {
	req := model.AddTeamReq{}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BuildErrorResponse(c, 400, "参数格式有误")
		return
	}

	// 调用 service 层
	err := service.AddTeam(req)
	if err != nil {
		utils.BuildServerError(c, "添加团队失败", err)
		return
	}

	utils.BuildSuccessResponse(c, "添加成功")
}

func GetTeams(c *gin.Context) {
	// 获取查询参数（都是可选参数）
	name := c.Query("name")
	isExistStr := c.Query("isExist")

	var isExist *bool
	if isExistStr != "" {
		isExistValue := isExistStr == "true"
		isExist = &isExistValue
	}

	// 用户校验
	var err error
	var teams []model.TeamVO
	_, ok := c.Get("userID")
	if !ok {
		// 如果是游客
		teams, err = service.GetTeams(name, isExist, false)
	} else {
		// 如果是用户或管理员
		teams, err = service.GetTeams(name, isExist, true)
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		utils.BuildServerError(c, "查询团队失败", err)
		return
	}
	utils.BuildSuccessResponse(c, teams)
	return
}

func UpdateTeam(c *gin.Context) {
	req := model.AddTeamReq{}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BuildErrorResponse(c, 400, "参数格式有误")
		return
	}

	err := service.UpdateTeam(req)
	if err != nil {
		utils.BuildServerError(c, "更新团队失败", err)
		return
	}
	utils.BuildSuccessResponse(c, "更新成功")
}

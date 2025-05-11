package handler

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"awesomeProject/model"
	"awesomeProject/service"
	"awesomeProject/utils"
)

func GetMemberList(c *gin.Context) {
	team := c.Query("team")
	isGraduate := c.Query("isGraduate")
	var teamQuery *string
	var isGraduateQuery *int
	if team != "" {
		teamQuery = &team
	}
	if isGraduate != "" {
		isGraduateInt, err := strconv.Atoi(isGraduate)
		if err != nil {
			utils.BuildErrorResponse(c, 500, "GetMemberList parse isGraduate failed err is: "+err.Error())
			return
		}
		isGraduateQuery = &isGraduateInt
	}
	res, err := service.GetMemberList(teamQuery, isGraduateQuery)
	if err != nil {
		utils.BuildErrorResponse(c, 500, "GetMemberList failed err is: "+err.Error())
		return
	}
	utils.BuildSuccessResponse(c, res)
}

func AddMember(c *gin.Context) {
	// 这功能暂时不做 不知道和注册的区别在哪儿
}

func ChangeMemberInfo(c *gin.Context) {
	// 定义请求体结构
	var req model.UpdateMemberRequest

	// 解析 JSON 请求
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BuildErrorResponse(c, 400, "Invalid JSON data")
		return
	}

	// 检查必填参数 UID
	if req.UID == 0 {
		utils.BuildErrorResponse(c, 400, "uid is required")
		return
	}

	// 调用 Service 层更新用户信息
	err := service.UpdateMember(req)
	if err != nil {
		utils.BuildErrorResponse(c, 500, err.Error())
		return
	}

	// 返回成功响应
	utils.BuildSuccessResponse(c, gin.H{
		"message": "User updated successfully",
	})
}

func GetMemberByName(c *gin.Context) {
	// 获取 userName 参数
	userName := c.Query("userName")
	if userName == "" {
		utils.BuildErrorResponse(c, 400, "userName is required")
		return
	}

	// 调用 Service 层获取用户数据
	member, err := service.GetMemberByUsername(userName)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.BuildErrorResponse(c, 404, "User not found")
		} else {
			utils.BuildErrorResponse(c, 500, "Failed to query user")
		}
		return
	}

	// 返回成功响应
	utils.BuildSuccessResponse(c, member)
}

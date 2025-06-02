package handler

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"awesomeProject/dao"
	"awesomeProject/model"
	"awesomeProject/service"
	"awesomeProject/utils"
)

// 添加成员
func Register(c *gin.Context) {
	var memReq *model.MemberRequest
	err := c.ShouldBindJSON(&memReq)
	if err != nil {
		utils.BuildErrorResponse(c, 400, "Register format error, error is "+err.Error())
		return
	}
	// 参数校验
	if memReq.Username == nil || memReq.Name == nil || memReq.Year == nil || memReq.Team == nil {
		utils.BuildErrorResponse(c, 400, "Register format error")
		return
	}

	err = service.Register(memReq)
	if err != nil {
		fmt.Printf("handler.Register format error, error is %s", err.Error())
		utils.BuildErrorResponse(c, 500, "Register Failed error is "+err.Error())
		return
	}
	utils.BuildSuccessResponse(c, "注册成功")
}

// 用户登录
func Login(c *gin.Context) {
	var memReq *model.MemberRequest
	err := c.ShouldBindJSON(&memReq)
	if err != nil {
		utils.BuildErrorResponse(c, 400, "Register format error, error is "+err.Error())
		return
	}
	// 参数校验
	if memReq.Username == nil || memReq.Password == nil || memReq.CaptchaID == nil || memReq.CaptchaData == nil {
		utils.BuildErrorResponse(c, 400, "Register format error")
	}

	// 校验验证码
	if ok := VerifyCaptcha(*memReq.CaptchaID, *memReq.CaptchaData); !ok {
		utils.BuildErrorResponse(c, 400, "Captcha verification failed or expired")
		return
	}

	data, err := service.Login(memReq)
	if err != nil {
		utils.BuildErrorResponse(c, 500, "Register Failed error is "+err.Error())
		return
	}
	utils.BuildSuccessResponse(c, data)
}

// 获取用户列表
func GetMemberList(c *gin.Context) {
	teamStr := c.Query("team")
	isGraduateStr := c.Query("isGraduate")
	pageSizeStr := c.Query("pageSize")
	pageNumStr := c.Query("pageNum")
	yearStr := c.Query("year")

	var pageSize, pageNum, isGraduate, year *int
	var team *string
	if isGraduateStr != "" {
		var tmp int
		if isGraduateStr == "true" {
			tmp = 1
			isGraduate = &tmp
		} else if isGraduateStr == "false" {
			tmp = 0
			isGraduate = &tmp
		} else {
			utils.BuildErrorResponse(c, 400, "isGraduate format error")
			return
		}
	}
	if pageSizeStr != "" {
		pageSizeInt, _ := strconv.Atoi(pageSizeStr)
		pageSize = &pageSizeInt
	}
	if pageNumStr != "" {
		pageNumInt, _ := strconv.Atoi(pageNumStr)
		pageNum = &pageNumInt
	}
	if teamStr != "" {
		team = &teamStr
	}
	if yearStr != "" {
		yearInt, _ := strconv.Atoi(yearStr)
		year = &yearInt
	}

	res, total, err := service.GetMemberList(team, isGraduate, pageSize, pageNum, year)
	if err != nil {
		utils.BuildErrorResponse(c, 500, "GetMemberList failed err is: "+err.Error())
		return
	}

	utils.BuildSuccessResponse(c, gin.H{
		"data":  res,
		"total": total,
	})
}

func AddMember(c *gin.Context) {
	// 这功能暂时不做 不知道和注册的区别在哪儿
}

// 修改用户信息
func ChangeMemberInfo(c *gin.Context) {
	// 定义请求体结构
	var req model.UpdateMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BuildErrorResponse(c, 400, "Invalid JSON data")
		return
	}

	// 检查必填参数Username
	if req.Username == nil {
		utils.BuildErrorResponse(c, 400, "username is required")
		return
	}

	status, ok := c.Get("status")
	if !ok {
		utils.BuildErrorResponse(c, 500, "status is not found")
		return
	}
	statusInt := status.(int)

	// 调用 Service 层更新用户信息
	err := service.UpdateMember(req, statusInt)
	if err != nil {
		utils.BuildErrorResponse(c, 500, err.Error())
		return
	}

	// 返回成功响应
	utils.BuildSuccessResponse(c, "修改成功")
}

// 根据用户名获取用户信息
func GetMemberByName(c *gin.Context) {
	userName := c.Query("username")
	if userName == "" {
		utils.BuildErrorResponse(c, 400, "username is not found")
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

func GetYears(c *gin.Context) {
	years, err := service.GetYears()
	if err != nil {
		utils.BuildErrorResponse(c, 500, "GetYears failed err is: "+err.Error())
		return
	}
	utils.BuildSuccessResponse(c, years)
	return
}

func DeleteMember(c *gin.Context) {
	uidStr := c.Query("uid")
	if uidStr == "" {
		utils.BuildErrorResponse(c, 400, "uidStr is not found")
		return
	}

	uid, _ := strconv.ParseInt(uidStr, 10, 64)
	err := dao.DeleteMember(uid)
	if err != nil {
		utils.BuildErrorResponse(c, 500, "DeleteMember failed err is: "+err.Error())
		return
	}
	utils.BuildSuccessResponse(c, "删除成功")
}

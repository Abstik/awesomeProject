package handler

import (
	"errors"
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
		utils.BuildErrorResponse(c, 400, "参数格式有误")
		return
	}
	// 参数校验
	if memReq.Username == nil || memReq.Name == nil || memReq.Year == nil || memReq.Team == nil {
		utils.BuildErrorResponse(c, 400, "参数格式有误")
		return
	}

	err = service.Register(memReq)
	if err != nil {
		if err.Error() == "用户名已存在" {
			utils.BuildErrorResponse(c, 400, "用户名已存在")
			return
		}
		utils.BuildErrorResponse(c, 500, "服务器繁忙")
		return
	}
	utils.BuildSuccessResponse(c, "注册成功")
}

// 用户登录
func Login(c *gin.Context) {
	var memReq *model.MemberRequest
	err := c.ShouldBindJSON(&memReq)
	if err != nil {
		utils.BuildErrorResponse(c, 400, "参数格式有误")
		return
	}
	// 参数校验
	if memReq.Username == nil || memReq.Password == nil || memReq.CaptchaID == nil || memReq.CaptchaData == nil {
		utils.BuildErrorResponse(c, 400, "参数格式有误")
	}

	// 校验验证码
	if ok := VerifyCaptcha(*memReq.CaptchaID, *memReq.CaptchaData); !ok {
		utils.BuildErrorResponse(c, 400, "验证码有误")
		return
	}

	data, err := service.Login(memReq)
	if err != nil {
		// 如果是用户名不存在
		if err.Error() == "用户名不存在" || err.Error() == "密码错误" {
			utils.BuildErrorResponse(c, 400, "用户名或密码错误")
			return
		}
		utils.BuildErrorResponse(c, 500, "服务器繁忙")
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
			utils.BuildErrorResponse(c, 400, "参数格式有误")
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
		utils.BuildErrorResponse(c, 500, "服务器繁忙")
		return
	}

	utils.BuildSuccessResponse(c, gin.H{
		"data":  res,
		"total": total,
	})
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
	var err error
	if statusInt == 0 {
		err = service.AdminUpdateMember(req)
	} else {
		err = service.UserUpdateMember(req)
	}
	if err != nil {
		utils.BuildErrorResponse(c, 500, err.Error())
		return
	}

	// 返回成功响应
	utils.BuildSuccessResponse(c, "修改成功")
}

// 根据用户名获取用户信息
func GetMemberByUserName(c *gin.Context) {
	userName := c.Query("username")
	if userName == "" {
		utils.BuildErrorResponse(c, 400, "username is not found")
		return
	}

	// 调用 Service 层获取用户数据
	member, err := service.GetMemberByUsername(userName)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		utils.BuildErrorResponse(c, 500, "查询失败")
		return
	}

	// 返回成功响应
	utils.BuildSuccessResponse(c, member)
}

func GetMemberByName(c *gin.Context) {
	name := c.Query("name")
	if name == "" {
		utils.BuildErrorResponse(c, 400, "name is not found")
		return
	}

	// 调用 Service 层获取用户数据
	members, err := service.GetMemberByName(name)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		utils.BuildErrorResponse(c, 500, "查询失败")
		return
	}

	// 成功响应
	utils.BuildSuccessResponse(c, members)
}

func GetYears(c *gin.Context) {
	years, err := service.GetYears()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		utils.BuildErrorResponse(c, 500, "查询失败")
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

func ResetPassword(c *gin.Context) {
	username := c.PostForm("username")
	user, err := dao.GetMemberByUsername(username)
	if err != nil {
		utils.BuildErrorResponse(c, 400, err.Error())
	}
	if user == nil || user.Status == nil || user.Username == nil {
		utils.BuildErrorResponse(c, 500, "此用户状态不合法")
	}
	if *user.Status == 0 {
		utils.BuildErrorResponse(c, 400, "无法修改管理员账号")
	}

	password := utils.EncryptPassword(*user.Username + "123")

	err = dao.ResetPassword(username, password)
	if err != nil {
		utils.BuildErrorResponse(c, 500, err.Error())
	}
	utils.BuildSuccessResponse(c, "重置成功")
}
